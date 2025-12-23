package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/repository"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
)

type CardService struct {
	cardRepo        *repository.CardRepository
	transactionRepo *repository.CardTransactionRepository
	cardIssuer      *CardIssuerService
	config          *config.Config
	mqChannel       *amqp.Channel
}

func NewCardService(
	cardRepo *repository.CardRepository,
	transactionRepo *repository.CardTransactionRepository,
	cardIssuer *CardIssuerService,
	mqChannel *amqp.Channel,
	cfg *config.Config,
) *CardService {
	return &CardService{
		cardRepo:        cardRepo,
		transactionRepo: transactionRepo,
		cardIssuer:      cardIssuer,
		config:          cfg,
		mqChannel:       mqChannel,
	}
}

func (s *CardService) CreateCard(userID string, req *models.CreateCardRequest) (*models.Card, error) {
	// Validation des limites de cartes par utilisateur
	userCards, err := s.cardRepo.GetUserCards(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user cards: %w", err)
	}

	if len(userCards) >= s.config.MaxCardsPerUser {
		return nil, fmt.Errorf("maximum number of cards reached")
	}

	// Générer les détails de la carte
	cardNumber, cvv, err := s.generateCardDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to generate card details: %w", err)
	}

	// Chiffrer les données sensibles
	cardNumberFull, err := s.encryptCardNumber(cardNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt card number: %w", err)
	}

	cvvEncrypted, err := s.encryptCVV(cvv)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt CVV: %w", err)
	}

	// Définir les limites par défaut
	limits := s.getDefaultLimits(req.CardType, req.Currency)
	if req.DailyLimit != nil {
		limits.DailyLimit = *req.DailyLimit
	}
	if req.MonthlyLimit != nil {
		limits.MonthlyLimit = *req.MonthlyLimit
	}

	// Générer un ID unique
	cardID, err := s.generateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate card ID: %w", err)
	}

	// Créer la carte
	card := &models.Card{
		ID:                cardID,
		UserID:            userID,
		CardNumber:        s.maskCardNumber(cardNumber),
		CardNumberFull:    cardNumberFull,
		CardType:          req.CardType,
		CardCategory:      req.CardCategory,
		Currency:          strings.ToUpper(req.Currency),
		Balance:           req.InitialAmount,
		AvailableBalance:  req.InitialAmount,
		CardholderName:    req.CardholderName,
		ExpiryMonth:       int(time.Now().Month()),
		ExpiryYear:        time.Now().Year() + 3, // 3 ans de validité
		CVV:               cvvEncrypted,
		Status:            "inactive",
		IsVirtual:         req.CardType == "virtual",
		IsActive:          false,
		ExpiresAt:         time.Now().AddDate(3, 0, 0), // 3 ans
		
		// Limites
		DailyLimit:        limits.DailyLimit,
		MonthlyLimit:      limits.MonthlyLimit,
		SingleTxLimit:     limits.SingleTxLimit,
		ATMDailyLimit:     limits.ATMDailyLimit,
		OnlineTxLimit:     limits.OnlineTxLimit,
		
		// Paramètres par défaut
		AllowATM:          req.AllowATM != nil && *req.AllowATM,
		AllowOnline:       req.AllowOnline == nil || *req.AllowOnline, // true par défaut
		AllowInternational: req.AllowInternational != nil && *req.AllowInternational,
		AllowContactless:  req.AllowContactless == nil || *req.AllowContactless, // true par défaut
		
		// Issuer
		IssuerID:          s.config.DefaultCardIssuer,
	}

	// Créer la carte chez l'émetteur (Marqeta, etc.)
	externalCardID, err := s.cardIssuer.CreateCard(card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card with issuer: %w", err)
	}
	card.ExternalCardID = &externalCardID

	// Sauvegarder en base
	err = s.cardRepo.Create(card)
	if err != nil {
		// Tenter de supprimer chez l'émetteur en cas d'erreur
		if externalCardID != "" {
			s.cardIssuer.CancelCard(externalCardID)
		}
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	// Si montant initial > 0, charger la carte
	if req.InitialAmount > 0 && req.SourceWalletID != nil {
		s.loadCardFromWallet(card.ID, *req.SourceWalletID, req.InitialAmount, userID)
	}

	// Publier l'événement
	s.publishCardEvent("card.created", card)

	// Supprimer les données sensibles pour la réponse
	card.CardNumberFull = ""
	card.CVV = ""

	return card, nil
}

func (s *CardService) CreateVirtualCard(userID string, req *models.CreateVirtualCardRequest) (*models.Card, error) {
	// Convertir en CreateCardRequest
	createReq := &models.CreateCardRequest{
		CardType:       "virtual",
		CardCategory:   "personal",
		Currency:       req.Currency,
		CardholderName: req.CardholderName,
		InitialAmount:  req.InitialAmount,
		SourceWalletID: &req.SourceWalletID,
	}

	// Définir la validité
	validityMonths := 24 // 2 ans par défaut
	if req.ValidityMonths != nil {
		validityMonths = *req.ValidityMonths
	}

	card, err := s.CreateCard(userID, createReq)
	if err != nil {
		return nil, err
	}

	// Mettre à jour la date d'expiration
	card.ExpiresAt = time.Now().AddDate(0, validityMonths, 0)
	s.cardRepo.Update(card)

	return card, nil
}

func (s *CardService) OrderPhysicalCard(userID string, req *models.OrderPhysicalCardRequest) (*models.Card, error) {
	// Convertir en CreateCardRequest
	createReq := &models.CreateCardRequest{
		CardType:       "prepaid",
		CardCategory:   "personal",
		Currency:       req.Currency,
		CardholderName: req.CardholderName,
		InitialAmount:  req.InitialAmount,
		SourceWalletID: &req.SourceWalletID,
	}

	card, err := s.CreateCard(userID, createReq)
	if err != nil {
		return nil, err
	}

	// Ajouter les informations de livraison
	shippingAddressJSON, _ := json.Marshal(req.ShippingAddress)
	shippingAddress := string(shippingAddressJSON)
	card.ShippingAddress = &shippingAddress
	card.ShippingStatus = stringPtr("processing")
	card.IsVirtual = false

	// Calculer les frais de livraison
	shippingFee := s.config.CardFees["shipping"]
	if req.ExpressShipping {
		shippingFee = s.config.CardFees["express_shipping"]
	}
	_ = shippingFee // TODO: use this for billing

	// Organiser la production et livraison
	shippingMap := map[string]string{
		"street":      req.ShippingAddress.AddressLine1,
		"city":        req.ShippingAddress.City,
		"postal_code": req.ShippingAddress.PostalCode,
		"country":     req.ShippingAddress.Country,
	}
	externalID := ""
	if card.ExternalCardID != nil {
		externalID = *card.ExternalCardID
	}
	trackingNumber, err := s.cardIssuer.OrderPhysicalCard(externalID, shippingMap)
	if err != nil {
		return nil, fmt.Errorf("failed to order physical card: %w", err)
	}
	
	card.TrackingNumber = &trackingNumber
	card.ShippingStatus = stringPtr("shipped")
	now := time.Now()
	card.ShippedAt = &now

	// Sauvegarder les modifications
	s.cardRepo.Update(card)

	// Publier l'événement
	s.publishCardEvent("card.shipped", card)

	return card, nil
}

func (s *CardService) LoadCard(userID, cardID string, req *models.LoadCardRequest) error {
	// Vérifier la propriété de la carte
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	if card.UserID != userID {
		return fmt.Errorf("card does not belong to user")
	}

	if !card.IsActive || card.Status != "active" {
		return fmt.Errorf("card is not active")
	}

	// Vérifier les limites de chargement
	if req.Amount < s.config.MinLoadAmount {
		return fmt.Errorf("amount below minimum load amount")
	}

	if req.Amount > s.config.MaxLoadAmount {
		return fmt.Errorf("amount exceeds maximum load amount")
	}

	// TODO: Vérifier les fonds du portefeuille source via Wallet Service

	// Calculer les frais de chargement
	loadFee := s.calculateLoadFee(req.Amount, card.Currency)

	// Créer la transaction de chargement
	transaction := &models.CardTransaction{
		CardID:          cardID,
		UserID:          userID,
		TransactionType: "load",
		Amount:          req.Amount,
		Currency:        card.Currency,
		Fee:             loadFee,
		Status:          "pending",
	}

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return fmt.Errorf("failed to create load transaction: %w", err)
	}

	// Traiter le chargement
	go s.processCardLoad(cardID, req.Amount, loadFee, transaction.ID)

	return nil
}

func (s *CardService) SetupAutoLoad(userID, cardID string, req *models.SetupAutoLoadRequest) error {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	if card.UserID != userID {
		return fmt.Errorf("card does not belong to user")
	}

	// Valider les paramètres
	if req.ReloadThreshold >= req.ReloadAmount {
		return fmt.Errorf("reload threshold must be less than reload amount")
	}

	if req.ReloadAmount > s.config.MaxAutoReloadAmount {
		return fmt.Errorf("reload amount exceeds maximum auto-reload amount")
	}

	// Mettre à jour la carte
	card.AutoReloadEnabled = true
	card.AutoReloadAmount = req.ReloadAmount
	card.AutoReloadThreshold = req.ReloadThreshold
	card.ReloadWalletID = &req.SourceWalletID

	err = s.cardRepo.Update(card)
	if err != nil {
		return fmt.Errorf("failed to setup auto-reload: %w", err)
	}

	// Publier l'événement
	s.publishCardEvent("card.auto_reload_enabled", card)

	return nil
}

func (s *CardService) SetCardPIN(userID, cardID, pin string) error {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	if card.UserID != userID {
		return fmt.Errorf("card does not belong to user")
	}

	// Valider le PIN
	if len(pin) != 4 || !s.isNumeric(pin) {
		return fmt.Errorf("PIN must be 4 digits")
	}

	// Hacher le PIN
	pinHash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	// Mettre à jour la carte
	card.PINHash = string(pinHash)
	err = s.cardRepo.Update(card)
	if err != nil {
		return fmt.Errorf("failed to set PIN: %w", err)
	}

	// Notifier l'émetteur de carte
	if card.ExternalCardID != nil {
		s.cardIssuer.UpdateCardPIN(*card.ExternalCardID, pin)
	}

	return nil
}

func (s *CardService) CreateGiftCard(userID string, req *models.CreateGiftCardRequest) (*models.GiftCard, error) {
	// Générer le code de gift card
	code, err := s.generateGiftCardCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate gift card code: %w", err)
	}

	// Calculer la date d'expiration
	validityDays := 365 // 1 an par défaut
	if req.ValidityDays != nil {
		validityDays = *req.ValidityDays
	}

	giftCard := &models.GiftCard{
		Code:           code,
		SenderID:       userID,
		RecipientEmail: req.RecipientEmail,
		RecipientPhone: req.RecipientPhone,
		Amount:         req.Amount,
		Currency:       strings.ToUpper(req.Currency),
		Message:        req.Message,
		Design:         req.Design,
		Status:         "pending",
		ExpiresAt:      time.Now().AddDate(0, 0, validityDays),
	}

	// TODO: Débiter le portefeuille source

	// Sauvegarder la gift card
	err = s.cardRepo.CreateGiftCard(giftCard)
	if err != nil {
		return nil, fmt.Errorf("failed to create gift card: %w", err)
	}

	// Publier l'événement
	s.publishGiftCardEvent("gift_card.created", giftCard)

	return giftCard, nil
}

func (s *CardService) RedeemGiftCard(userID, code, targetWalletID string) error {
	// Récupérer la gift card
	giftCard, err := s.cardRepo.GetGiftCardByCode(code)
	if err != nil {
		return fmt.Errorf("gift card not found or invalid")
	}

	// Vérifier la validité
	if giftCard.Status != "sent" {
		return fmt.Errorf("gift card is not available for redemption")
	}

	if time.Now().After(giftCard.ExpiresAt) {
		return fmt.Errorf("gift card has expired")
	}

	if giftCard.RedeemedBy != nil {
		return fmt.Errorf("gift card has already been redeemed")
	}

	// Marquer comme utilisée
	giftCard.Status = "redeemed"
	giftCard.RedeemedBy = &userID
	now := time.Now()
	giftCard.RedeemedAt = &now

	err = s.cardRepo.UpdateGiftCard(giftCard)
	if err != nil {
		return fmt.Errorf("failed to redeem gift card: %w", err)
	}

	// TODO: Créditer le portefeuille de destination

	// Publier l'événement
	s.publishGiftCardEvent("gift_card.redeemed", giftCard)

	return nil
}

// GetUserCards returns all cards for a user
func (s *CardService) GetUserCards(userID string) ([]models.Card, error) {
	return s.cardRepo.GetUserCards(userID)
}

// GetCard returns a specific card by ID
func (s *CardService) GetCard(userID, cardID string) (*models.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, err
	}
	if card.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to card")
	}
	return card, nil
}

// UpdateCard updates card details
func (s *CardService) UpdateCard(userID, cardID string, req *models.UpdateCardRequest) (*models.Card, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	if req.Name != nil {
		card.CardholderName = *req.Name
	}
	
	if err := s.cardRepo.Update(card); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}
	
	s.publishCardEvent("card.updated", card)
	return card, nil
}

// DeleteCard soft deletes a card
func (s *CardService) DeleteCard(userID, cardID string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	card.Status = "deleted"
	card.IsActive = false
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	
	s.publishCardEvent("card.deleted", card)
	return nil
}

// ActivateCard activates a card
func (s *CardService) ActivateCard(userID, cardID string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	if card.Status == "blocked" {
		return fmt.Errorf("cannot activate blocked card")
	}
	
	card.Status = "active"
	card.IsActive = true
	now := time.Now()
	card.ActivatedAt = &now
	
	if card.ExternalCardID != nil {
		if err := s.cardIssuer.ActivateCard(*card.ExternalCardID); err != nil {
			return fmt.Errorf("failed to activate card with issuer: %w", err)
		}
	}
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}
	
	s.publishCardEvent("card.activated", card)
	return nil
}

// DeactivateCard deactivates a card
func (s *CardService) DeactivateCard(userID, cardID string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	card.Status = "inactive"
	card.IsActive = false
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}
	
	s.publishCardEvent("card.deactivated", card)
	return nil
}

// FreezeCard temporarily freezes a card
func (s *CardService) FreezeCard(userID, cardID, reason string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	card.Status = "frozen"
	card.FreezeReason = &reason
	now := time.Now()
	card.FrozenAt = &now
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to freeze card: %w", err)
	}
	
	s.publishCardEvent("card.frozen", card)
	return nil
}

// UnfreezeCard unfreezes a frozen card
func (s *CardService) UnfreezeCard(userID, cardID string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	if card.Status != "frozen" {
		return fmt.Errorf("card is not frozen")
	}
	
	card.Status = "active"
	card.IsActive = true
	card.FreezeReason = nil
	card.FrozenAt = nil
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to unfreeze card: %w", err)
	}
	
	s.publishCardEvent("card.unfrozen", card)
	return nil
}

// BlockCard permanently blocks a card
func (s *CardService) BlockCard(userID, cardID, reason string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	card.Status = "blocked"
	card.IsActive = false
	card.BlockReason = &reason
	now := time.Now()
	card.BlockedAt = &now
	
	if card.ExternalCardID != nil {
		if err := s.cardIssuer.BlockCard(*card.ExternalCardID, reason); err != nil {
			return fmt.Errorf("failed to block card with issuer: %w", err)
		}
	}
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to block card: %w", err)
	}
	
	s.publishCardEvent("card.blocked", card)
	return nil
}

// CancelAutoLoad cancels auto-load for a card
func (s *CardService) CancelAutoLoad(userID, cardID string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	card.AutoReloadEnabled = false
	card.AutoReloadAmount = 0
	card.AutoReloadThreshold = 0
	card.AutoReloadWalletID = nil
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to cancel auto-load: %w", err)
	}
	
	s.publishCardEvent("card.autoload_cancelled", card)
	return nil
}

// GetCardLimits returns the current limits for a card
func (s *CardService) GetCardLimits(userID, cardID string) (*models.CardLimits, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	return &models.CardLimits{
		DailyLimit:    card.DailyLimit,
		MonthlyLimit:  card.MonthlyLimit,
		SingleTxLimit: card.SingleTxLimit,
		ATMDailyLimit: card.ATMDailyLimit,
		Currency:      card.Currency,
	}, nil
}

// UpdateCardLimits updates the limits for a card
func (s *CardService) UpdateCardLimits(userID, cardID string, req *models.UpdateCardLimitsRequest) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	if req.DailyLimit != nil {
		card.DailyLimit = *req.DailyLimit
	}
	if req.MonthlyLimit != nil {
		card.MonthlyLimit = *req.MonthlyLimit
	}
	if req.SingleTxLimit != nil {
		card.SingleTxLimit = *req.SingleTxLimit
	}
	if req.ATMDailyLimit != nil {
		card.ATMDailyLimit = *req.ATMDailyLimit
	}
	
	if err := s.cardRepo.Update(card); err != nil {
		return fmt.Errorf("failed to update card limits: %w", err)
	}
	
	s.publishCardEvent("card.limits_updated", card)
	return nil
}

// GetCardTransactions returns transactions for a card
func (s *CardService) GetCardTransactions(userID, cardID string, limit, offset int) ([]models.CardTransaction, error) {
	_, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	return s.transactionRepo.GetByCardID(cardID, limit, offset)
}

// GetCardBalance returns the current balance of a card
func (s *CardService) GetCardBalance(userID, cardID string) (float64, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return 0, err
	}
	return card.Balance, nil
}

// GetCardDetails returns detailed card information including sensitive data
func (s *CardService) GetCardDetails(userID, cardID string) (*models.Card, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	return card, nil
}

// ChangeCardPIN changes the PIN after verifying the current PIN
func (s *CardService) ChangeCardPIN(userID, cardID, currentPIN, newPIN string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	// Verify current PIN
	if err := bcrypt.CompareHashAndPassword([]byte(card.PINHash), []byte(currentPIN)); err != nil {
		return fmt.Errorf("invalid current PIN")
	}
	
	// Set new PIN
	return s.SetCardPIN(userID, cardID, newPIN)
}

// ResetCardPIN initiates a PIN reset process
func (s *CardService) ResetCardPIN(userID, cardID string) error {
	_, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	// TODO: Send PIN reset email/SMS
	s.publishCardEvent("card.pin_reset_requested", nil)
	return nil
}

// RegenerateVirtualCard regenerates a virtual card with new details
func (s *CardService) RegenerateVirtualCard(userID, cardID string) (*models.Card, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	if !card.IsVirtual {
		return nil, fmt.Errorf("can only regenerate virtual cards")
	}
	
	// Generate new card details
	newCardNumber, newCVV, err := s.generateCardDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new card details: %w", err)
	}
	
	card.CardNumberFull = newCardNumber
	card.CardNumber = s.maskCardNumber(newCardNumber)
	card.CVV = newCVV
	card.ExpiryMonth = int(time.Now().Month())
	card.ExpiryYear = time.Now().Year() + 3
	
	if err := s.cardRepo.Update(card); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}
	
	s.publishCardEvent("card.regenerated", card)
	return card, nil
}

// GenerateCardQR generates a QR code for virtual card payments
func (s *CardService) GenerateCardQR(userID, cardID string) (string, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return "", err
	}
	
	if !card.IsVirtual {
		return "", fmt.Errorf("QR codes only available for virtual cards")
	}
	
	// Generate QR code data (in production, this would be encrypted)
	qrData := fmt.Sprintf("ZEKORA:%s:%s", card.ID, card.CardNumber)
	return qrData, nil
}

// GetShippingStatus returns the shipping status for a physical card
func (s *CardService) GetShippingStatus(userID, cardID string) (*models.ShippingStatus, error) {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	if card.IsVirtual {
		return nil, fmt.Errorf("virtual cards don't have shipping status")
	}
	
	status := &models.ShippingStatus{
		Status:         "pending",
		TrackingNumber: card.TrackingNumber,
		ShippedAt:      card.ShippedAt,
		DeliveredAt:    card.DeliveredAt,
	}
	
	if card.ShippingStatus != nil {
		status.Status = *card.ShippingStatus
	}
	
	return status, nil
}

// ActivatePhysicalCard activates a physical card with activation code
func (s *CardService) ActivatePhysicalCard(userID, cardID, activationCode, lastFourDigits string) error {
	card, err := s.GetCard(userID, cardID)
	if err != nil {
		return err
	}
	
	if card.IsVirtual {
		return fmt.Errorf("cannot activate physical card: this is a virtual card")
	}
	
	// Verify last four digits match
	if len(card.CardNumber) >= 4 {
		last4 := card.CardNumber[len(card.CardNumber)-4:]
		if last4 != lastFourDigits {
			return fmt.Errorf("card verification failed")
		}
	}
	
	// Activate the card
	return s.ActivateCard(userID, cardID)
}

// ReportLostCard reports a card as lost and blocks it
func (s *CardService) ReportLostCard(userID, cardID string) error {
	return s.BlockCard(userID, cardID, "reported_lost")
}

// RequestReplacement requests a replacement for a lost/damaged card
func (s *CardService) RequestReplacement(userID, cardID, reason string, shippingAddress *models.ShippingAddress) (*models.Card, error) {
	oldCard, err := s.GetCard(userID, cardID)
	if err != nil {
		return nil, err
	}
	
	// Block the old card
	if err := s.BlockCard(userID, cardID, "replacement_requested"); err != nil {
		return nil, err
	}
	
	// Create a new card with the same settings
	req := &models.OrderPhysicalCardRequest{
		Currency:        oldCard.Currency,
		CardholderName:  oldCard.CardholderName,
		ShippingAddress: *shippingAddress,
	}
	
	newCard, err := s.OrderPhysicalCard(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to order replacement card: %w", err)
	}
	
	s.publishCardEvent("card.replacement_requested", newCard)
	return newCard, nil
}

// SendGiftCard sends a gift card to someone
func (s *CardService) SendGiftCard(userID, giftCardID, message string) error {
	giftCard, err := s.cardRepo.GetGiftCardByID(giftCardID)
	if err != nil {
		return fmt.Errorf("gift card not found: %w", err)
	}
	
	if giftCard.SenderID != userID {
		return fmt.Errorf("unauthorized")
	}
	
	if giftCard.Status != "active" {
		return fmt.Errorf("gift card not available for sending")
	}
	
	giftCard.Status = "sent"
	if err := s.cardRepo.UpdateGiftCard(giftCard); err != nil {
		return fmt.Errorf("failed to update gift card: %w", err)
	}
	
	// TODO: Send notification to recipient
	s.publishGiftCardEvent("gift_card.sent", giftCard)
	return nil
}

// GetUserGiftCards returns all gift cards for a user
func (s *CardService) GetUserGiftCards(userID string) ([]models.GiftCard, error) {
	return s.cardRepo.GetUserGiftCards(userID)
}

// Méthodes privées

func (s *CardService) generateCardDetails() (cardNumber, cvv string, err error) {
	// Générer un numéro de carte (format Luhn valide)
	cardNumber = s.generateCardNumber()
	
	// Générer le CVV
	cvvNum, err := rand.Int(rand.Reader, big.NewInt(900))
	if err != nil {
		return "", "", err
	}
	cvv = fmt.Sprintf("%03d", cvvNum.Int64()+100)

	return cardNumber, cvv, nil
}

func (s *CardService) generateCardNumber() string {
	// Générer un numéro de carte commençant par 4 (Visa) ou 5 (Mastercard)
	prefix := "4532" // Préfixe Visa test
	
	// Générer 12 chiffres aléatoires
	number := prefix
	for i := 0; i < 12; i++ {
		digit, _ := rand.Int(rand.Reader, big.NewInt(10))
		number += fmt.Sprintf("%d", digit.Int64())
	}

	// Calculer et ajouter le chiffre de contrôle Luhn
	checksum := s.calculateLuhnChecksum(number)
	return number + fmt.Sprintf("%d", checksum)
}

func (s *CardService) calculateLuhnChecksum(number string) int {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		
		sum += digit
		alternate = !alternate
	}

	return (10 - (sum % 10)) % 10
}

func (s *CardService) generateGiftCardCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16

	code := make([]byte, length)
	for i := range code {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[charIndex.Int64()]
	}

	// Format: XXXX-XXXX-XXXX-XXXX
	codeStr := string(code)
	return fmt.Sprintf("%s-%s-%s-%s", codeStr[0:4], codeStr[4:8], codeStr[8:12], codeStr[12:16]), nil
}

func (s *CardService) maskCardNumber(cardNumber string) string {
	if len(cardNumber) < 8 {
		return cardNumber
	}
	
	masked := "****-****-****-" + cardNumber[len(cardNumber)-4:]
	return masked
}

func (s *CardService) encryptCardNumber(cardNumber string) (string, error) {
	// Dans un vrai système, utiliser un chiffrement symétrique fort
	// Pour la démo, on retourne le numéro hashé
	hash, err := bcrypt.GenerateFromPassword([]byte(cardNumber), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *CardService) encryptCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *CardService) getDefaultLimits(cardType, currency string) *models.CardLimits {
	limits := s.config.DefaultLimits[cardType]
	if limits == nil {
		// Limites par défaut si non configurées
		return &models.CardLimits{
			DailyLimit:    1000,
			MonthlyLimit:  5000,
			SingleTxLimit: 500,
			ATMDailyLimit: 300,
			OnlineTxLimit: 1000,
		}
	}
	return limits
}

func (s *CardService) calculateLoadFee(amount float64, currency string) float64 {
	feeRate := s.config.CardFees["load_fee_percentage"]
	minimumFee := s.config.CardFees["load_fee_minimum"]
	
	fee := amount * feeRate / 100
	if fee < minimumFee {
		fee = minimumFee
	}
	
	return fee
}

func (s *CardService) processCardLoad(cardID string, amount, fee float64, transactionID string) {
	// Simuler le traitement
	time.Sleep(2 * time.Second)

	// Mettre à jour le solde de la carte
	s.cardRepo.UpdateBalance(cardID, amount)
	
	// Marquer la transaction comme terminée
	s.transactionRepo.UpdateStatus(transactionID, "completed")

	// Publier l'événement
	event := map[string]interface{}{
		"type":           "card.loaded",
		"card_id":        cardID,
		"amount":         amount,
		"fee":            fee,
		"transaction_id": transactionID,
		"timestamp":      time.Now(),
	}

	s.publishEvent("card.events", event)
}

func (s *CardService) isNumeric(str string) bool {
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (s *CardService) publishCardEvent(eventType string, card *models.Card) {
	event := map[string]interface{}{
		"type":     eventType,
		"card_id":  card.ID,
		"user_id":  card.UserID,
		"currency": card.Currency,
		"status":   card.Status,
		"timestamp": time.Now(),
	}

	s.publishEvent("card.events", event)
}

func (s *CardService) publishGiftCardEvent(eventType string, giftCard *models.GiftCard) {
	event := map[string]interface{}{
		"type":        eventType,
		"gift_card_id": giftCard.ID,
		"sender_id":   giftCard.SenderID,
		"amount":      giftCard.Amount,
		"currency":    giftCard.Currency,
		"status":      giftCard.Status,
		"timestamp":   time.Now(),
	}

	s.publishEvent("gift_card.events", event)
}

func (s *CardService) publishEvent(exchange string, event map[string]interface{}) {
	if s.mqChannel == nil {
		return
	}

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		exchange,        // exchange
		event["type"].(string), // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}

func stringPtr(s string) *string {
	return &s
}

func (s *CardService) loadCardFromWallet(cardID, walletID string, amount float64, userID string) {
	// TODO: Intégrer avec le Wallet Service pour débiter le portefeuille
	// et charger la carte
	loadReq := &models.LoadCardRequest{
		Amount:         amount,
		SourceWalletID: walletID,
		Description:    "Initial card load",
	}
	
	
	s.LoadCard(userID, cardID, loadReq)
}

func (s *CardService) generateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}