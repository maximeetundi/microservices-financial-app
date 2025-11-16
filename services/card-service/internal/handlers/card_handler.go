package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/card-service/internal/models"
	"github.com/crypto-bank/card-service/internal/services"
)

type CardHandler struct {
	cardService  *services.CardService
	walletClient *services.WalletClient
}

func NewCardHandler(cardService *services.CardService, walletClient *services.WalletClient) *CardHandler {
	return &CardHandler{
		cardService:  cardService,
		walletClient: walletClient,
	}
}

func (h *CardHandler) GetUserCards(c *gin.Context) {
	userID, _ := c.Get("user_id")

	cards, err := h.cardService.GetUserCards(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	kycLevel, _ := c.Get("kyc_level")

	var req models.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier le niveau KYC requis pour les cartes prépayées
	if req.CardType == "prepaid" && kycLevel.(int) < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "KYC level 2 required for prepaid cards",
			"required_kyc_level": 2,
			"current_kyc_level": kycLevel,
		})
		return
	}

	card, err := h.cardService.CreateCard(userID.(string), &req)
	if err != nil {
		if err.Error() == "maximum number of cards reached" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Card created successfully",
		"card":    card,
	})
}

func (h *CardHandler) GetCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	card, err := h.cardService.GetCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": card})
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.UpdateCard(userID.(string), cardID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card updated successfully"})
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.DeleteCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully"})
}

func (h *CardHandler) ActivateCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.ActivateCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card activated successfully"})
}

func (h *CardHandler) DeactivateCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.DeactivateCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deactivated successfully"})
}

func (h *CardHandler) FreezeCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.FreezeCard(userID.(string), cardID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card frozen successfully"})
}

func (h *CardHandler) UnfreezeCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.UnfreezeCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card unfrozen successfully"})
}

func (h *CardHandler) BlockCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.BlockCard(userID.(string), cardID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card blocked successfully"})
}

func (h *CardHandler) LoadCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req models.LoadCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.LoadCard(userID.(string), cardID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card load initiated successfully"})
}

func (h *CardHandler) SetupAutoLoad(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req models.SetupAutoLoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.SetupAutoLoad(userID.(string), cardID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-reload setup successfully"})
}

func (h *CardHandler) CancelAutoLoad(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.CancelAutoLoad(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-reload cancelled successfully"})
}

func (h *CardHandler) GetCardLimits(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	limits, err := h.cardService.GetCardLimits(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"limits": limits})
}

func (h *CardHandler) UpdateCardLimits(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req models.UpdateCardLimitsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.UpdateCardLimits(userID.(string), cardID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card limits updated successfully"})
}

func (h *CardHandler) GetCardTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	txType := c.Query("type")
	status := c.Query("status")

	transactions, err := h.cardService.GetCardTransactions(userID.(string), cardID, limit, offset, txType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *CardHandler) GetCardBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	balance, err := h.cardService.GetCardBalance(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (h *CardHandler) GetCardDetails(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	details, err := h.cardService.GetCardDetails(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": details})
}

func (h *CardHandler) SetCardPIN(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req models.SetCardPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.SetCardPIN(userID.(string), cardID, req.PIN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN set successfully"})
}

func (h *CardHandler) ChangeCardPIN(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req models.ChangeCardPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.ChangeCardPIN(userID.(string), cardID, req.CurrentPIN, req.NewPIN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN changed successfully"})
}

func (h *CardHandler) ResetCardPIN(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.ResetCardPIN(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN reset initiated. Check your email for instructions."})
}

// Virtual Card specific handlers
func (h *CardHandler) CreateVirtualCard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateVirtualCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.CreateVirtualCard(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create virtual card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Virtual card created successfully",
		"card":    card,
	})
}

func (h *CardHandler) RegenerateVirtualCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	newCard, err := h.cardService.RegenerateVirtualCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Virtual card regenerated successfully",
		"card":    newCard,
	})
}

func (h *CardHandler) GetCardQR(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	qrCode, err := h.cardService.GenerateCardQR(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": qrCode})
}

// Physical Card specific handlers
func (h *CardHandler) OrderPhysicalCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	kycLevel, _ := c.Get("kyc_level")

	// Vérifier le niveau KYC requis pour les cartes physiques
	if kycLevel.(int) < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "KYC level 2 required for physical cards",
			"required_kyc_level": 2,
			"current_kyc_level": kycLevel,
		})
		return
	}

	var req models.OrderPhysicalCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.cardService.OrderPhysicalCard(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to order physical card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Physical card order placed successfully",
		"card":    card,
	})
}

func (h *CardHandler) GetShippingStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	status, err := h.cardService.GetShippingStatus(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shipping_status": status})
}

func (h *CardHandler) ActivatePhysicalCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req struct {
		ActivationCode string `json:"activation_code" binding:"required"`
		LastFourDigits string `json:"last_four_digits" binding:"required,len=4"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.ActivatePhysicalCard(userID.(string), cardID, req.ActivationCode, req.LastFourDigits)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Physical card activated successfully"})
}

func (h *CardHandler) ReportLostCard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	err := h.cardService.ReportLostCard(userID.(string), cardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card reported as lost. A replacement will be issued."})
}

func (h *CardHandler) RequestReplacement(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cardID := c.Param("card_id")

	var req struct {
		Reason          string `json:"reason" binding:"required"`
		ShippingAddress *models.ShippingAddress `json:"shipping_address,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCard, err := h.cardService.RequestReplacement(userID.(string), cardID, req.Reason, req.ShippingAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Replacement card ordered successfully",
		"new_card": newCard,
	})
}

// Gift Card handlers
func (h *CardHandler) CreateGiftCard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	giftCard, err := h.cardService.CreateGiftCard(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gift card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Gift card created successfully",
		"gift_card": giftCard,
	})
}

func (h *CardHandler) SendGiftCard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		GiftCardID string `json:"gift_card_id" binding:"required"`
		Message    string `json:"message,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.SendGiftCard(userID.(string), req.GiftCardID, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gift card sent successfully"})
}

func (h *CardHandler) RedeemGiftCard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.RedeemGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cardService.RedeemGiftCard(userID.(string), req.Code, req.TargetWalletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gift card redeemed successfully"})
}

func (h *CardHandler) GetGiftCards(c *gin.Context) {
	userID, _ := c.Get("user_id")

	giftCards, err := h.cardService.GetUserGiftCards(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gift cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gift_cards": giftCards})
}

// Webhook handlers
func (h *CardHandler) HandleMarqetaTransaction(c *gin.Context) {
	var webhook map[string]interface{}
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Traiter les notifications de transaction Marqeta
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

func (h *CardHandler) HandleMarqetaAuth(c *gin.Context) {
	var webhook map[string]interface{}
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Traiter les demandes d'autorisation Marqeta
	c.JSON(http.StatusOK, gin.H{"message": "Authorization processed"})
}

func (h *CardHandler) HandleIssuerCallback(c *gin.Context) {
	var callback map[string]interface{}
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Traiter les callbacks de l'émetteur de cartes
	c.JSON(http.StatusOK, gin.H{"message": "Callback processed"})
}

// Public endpoints
func (h *CardHandler) GetSupportedCurrencies(c *gin.Context) {
	currencies := []string{"USD", "EUR", "GBP", "CAD", "AUD", "BTC", "ETH"}
	c.JSON(http.StatusOK, gin.H{"currencies": currencies})
}

func (h *CardHandler) GetCardFees(c *gin.Context) {
	fees := models.CardFees{
		IssuanceFee:      5.00,
		MonthlyFee:       2.00,
		ATMWithdrawalFee: 2.50,
		ForeignTxFee:     3.00,
		ReloadFee:        1.00,
		ReplacementFee:   10.00,
		CurrencyFees: map[string]float64{
			"USD": 0.00,
			"EUR": 0.50,
			"GBP": 0.50,
			"BTC": 1.00,
			"ETH": 1.00,
		},
	}

	c.JSON(http.StatusOK, gin.H{"fees": fees})
}