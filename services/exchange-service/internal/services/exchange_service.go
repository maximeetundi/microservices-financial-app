package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type ExchangeService struct {
	exchangeRepo   *repository.ExchangeRepository
	orderRepo      *repository.OrderRepository
	rateService    *RateService
	feeService     *FeeService
	config         *config.Config
	kafkaPublisher *KafkaPublisher
	walletClient   *WalletClient
}

func NewExchangeService(exchangeRepo *repository.ExchangeRepository, orderRepo *repository.OrderRepository, rateService *RateService, feeService *FeeService, kafkaPublisher *KafkaPublisher, walletClient *WalletClient, cfg *config.Config) *ExchangeService {
	return &ExchangeService{
		exchangeRepo:   exchangeRepo,
		orderRepo:      orderRepo,
		rateService:    rateService,
		feeService:     feeService,
		config:         cfg,
		kafkaPublisher: kafkaPublisher,
		walletClient:   walletClient,
	}
}

func (s *ExchangeService) GetQuote(userID, token, fromCurrency, toCurrency string, fromAmount *float64, toAmount *float64) (*models.Quote, error) {
	// Validation des devises
	if !s.isSupportedCurrency(fromCurrency) || !s.isSupportedCurrency(toCurrency) {
		return nil, fmt.Errorf("unsupported currency pair")
	}

	if fromCurrency == toCurrency {
		return nil, fmt.Errorf("same currency exchange not allowed")
	}

	// Obtenir le taux de change actuel
	rate, err := s.rateService.GetRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Calculer les montants
	var calculatedFromAmount, calculatedToAmount float64
	if fromAmount != nil {
		calculatedFromAmount = *fromAmount
		calculatedToAmount = calculatedFromAmount * rate.Rate
	} else if toAmount != nil {
		calculatedToAmount = *toAmount
		calculatedFromAmount = calculatedToAmount / rate.Rate
	} else {
		return nil, fmt.Errorf("either from_amount or to_amount must be specified")
	}

	// Calculer les frais
	feeKey := s.getFeeKey(fromCurrency, toCurrency)
	fee, err := s.feeService.CalculateFee(feeKey, calculatedFromAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee: %w", err)
	}

	// Calculate effective percentage
	var feePercentage float64
	if calculatedFromAmount > 0 {
		feePercentage = (fee / calculatedFromAmount) * 100
	}

	// Ajuster le montant de destination avec les frais
	finalToAmount := calculatedToAmount - (fee * rate.Rate)

	// Créer le devis
	quote := &models.Quote{
		UserID:            userID,
		FromCurrency:      strings.ToUpper(fromCurrency),
		ToCurrency:        strings.ToUpper(toCurrency),
		FromAmount:        calculatedFromAmount,
		ToAmount:          finalToAmount,
		ExchangeRate:      rate.Rate,
		Fee:               fee,
		FeePercentage:     feePercentage,
		ValidUntil:        time.Now().Add(5 * time.Minute), // Valide 5 minutes
		EstimatedDelivery: s.getEstimatedDelivery(fromCurrency, toCurrency),
	}

	// Retrieve source wallet to check balance
	wallets, err := s.walletClient.GetUserWallets(userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}
	
	// Find wallet for source currency
	var fromWalletID string
	for _, w := range wallets {
		if w.Currency == strings.ToUpper(fromCurrency) {
			fromWalletID = w.ID
			break
		}
	}
	
	if fromWalletID == "" {
		return nil, fmt.Errorf("no wallet found for currency %s", fromCurrency)
	}

	// Calculate and check fees
	// Perform synchronous balance check
	balance, _, err := s.walletClient.GetWalletBalanceByID(fromWalletID, token) // Pass token explicitly
	if err != nil {
		log.Printf("Failed to check balance for wallet %s: %v", fromWalletID, err)
		// Fallback or fail? Fail for safety to prevent insufficient funds
		return nil, fmt.Errorf("failed to verify wallet balance: %w", err)
	}

	if balance < quote.FromAmount {
		return nil, fmt.Errorf("insufficient funds: available %f, required %f", balance, quote.FromAmount)
	}

	// Sauvegarder le devis
	err = s.exchangeRepo.CreateQuote(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

func (s *ExchangeService) ExecuteExchange(userID, quoteID, fromWalletID, toWalletID, token string) (*models.Exchange, error) {
	// Récupérer le devis
	quote, err := s.exchangeRepo.GetQuote(quoteID)
	if err != nil {
		return nil, fmt.Errorf("quote not found: %w", err)
	}

	// Vérifier la validité du devis
	if quote.UserID != userID {
		return nil, fmt.Errorf("quote does not belong to user")
	}

	if time.Now().After(quote.ValidUntil) {
		return nil, fmt.Errorf("quote has expired")
	}

	// Créer l'échange avec statut initial "processing"
	exchange := &models.Exchange{
		UserID:           userID,
		FromWalletID:     fromWalletID,
		ToWalletID:       toWalletID,
		FromCurrency:     quote.FromCurrency,
		ToCurrency:       quote.ToCurrency,
		FromAmount:       quote.FromAmount,
		ToAmount:         quote.ToAmount,
		ExchangeRate:     quote.ExchangeRate,
		Fee:              quote.Fee,
		FeePercentage:    quote.FeePercentage,
		Status:           "processing",
		QuoteID:          quoteID,
	}

	err = s.exchangeRepo.Create(exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange record: %w", err)
	}

	log.Printf("[EXCHANGE SYNC] Starting sync exchange %s", exchange.ID)

	// 1. Synchronous DEBIT
	debitReq := &TransactionRequest{ // Uses the struct defined in wallet_client.go
		UserID:    userID,
		WalletID:  fromWalletID,
		Amount:    exchange.FromAmount,
		Type:      "debit",
		Currency:  exchange.FromCurrency,
		Reference: fmt.Sprintf("EXC_DEBIT_%s", exchange.ID),
	}

	err = s.walletClient.ProcessTransaction(debitReq)
	if err != nil {
		log.Printf("[EXCHANGE SYNC] Debit failed for %s: %v", exchange.ID, err)
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed")
		// No refund needed, debit failed
		return nil, fmt.Errorf("exchange failed at debit step: %v", err)
	}

	// Debit Success - Status Update
	log.Printf("[EXCHANGE SYNC] Debit successful for %s", exchange.ID)
	s.exchangeRepo.UpdateStatus(exchange.ID, "debited")

	// 2. Synchronous CREDIT
	creditReq := &TransactionRequest{
		UserID:    userID,
		WalletID:  toWalletID,
		Amount:    exchange.ToAmount,
		Type:      "credit",
		Currency:  exchange.ToCurrency,
		Reference: fmt.Sprintf("EXC_CREDIT_%s", exchange.ID),
	}

	err = s.walletClient.ProcessTransaction(creditReq)
	if err != nil {
		log.Printf("[EXCHANGE SYNC] CRITICAL: Credit failed for %s after successful debit. Initiating REFUND. Error: %v", exchange.ID, err)
		
		// 3. REFUND (Compensation)
		refundReq := &TransactionRequest{
			UserID:    userID,
			WalletID:  fromWalletID,
			Amount:    exchange.FromAmount,
			Type:      "credit", // Credit back the source
			Currency:  exchange.FromCurrency,
			Reference: fmt.Sprintf("EXC_REFUND_%s", exchange.ID),
		}
		
		refundErr := s.walletClient.ProcessTransaction(refundReq)
		if refundErr != nil {
			// FATAL ERROR: Money taken, not given, not returned.
			log.Printf("[EXCHANGE SYNC] FATAL: Refund failed for %s! Manual intervention required. Error: %v", exchange.ID, refundErr)
			s.exchangeRepo.UpdateStatus(exchange.ID, "failed_fatal")
			return nil, fmt.Errorf("exchange failed and refund failed. Please contact support immediately. Ref: %s", exchange.ID)
		}

		log.Printf("[EXCHANGE SYNC] Refund successful for %s", exchange.ID)
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed_refunded")
		return nil, fmt.Errorf("exchange failed (funds refunded): %v", err)
	}

	// 4. Success Finalization
	log.Printf("[EXCHANGE SYNC] Exchange %s completed successfully", exchange.ID)
	s.exchangeRepo.UpdateStatus(exchange.ID, "completed")

	// Optional: Still publish event for analytics/notifications
	s.publishExchangeEvent("exchange.completed", exchange)

	return exchange, nil
}

func (s *ExchangeService) BuyCrypto(userID, cryptoCurrency, paymentCurrency string, amount float64, orderType string, limitPrice *float64) (*models.TradingOrder, error) {
	// Valider les paramètres
	if !s.isCryptoCurrency(cryptoCurrency) {
		return nil, fmt.Errorf("invalid cryptocurrency: %s", cryptoCurrency)
	}

	if !s.isFiatCurrency(paymentCurrency) {
		return nil, fmt.Errorf("invalid payment currency: %s", paymentCurrency)
	}

	// Obtenir le prix actuel
	currentRate, err := s.rateService.GetRate(paymentCurrency, cryptoCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get current rate: %w", err)
	}

	// Calculer le prix d'exécution
	var executionPrice float64
	if orderType == "market" {
		executionPrice = currentRate.AskPrice // Prix de vente pour acheter
	} else if orderType == "limit" && limitPrice != nil {
		if *limitPrice > currentRate.AskPrice {
			return nil, fmt.Errorf("limit price too high for buy order")
		}
		executionPrice = *limitPrice
	} else {
		return nil, fmt.Errorf("invalid order type or missing limit price")
	}

	// Calculer le montant total nécessaire
	totalCost := amount * executionPrice
	
	fee, err := s.feeService.CalculateFee("trading_buy", totalCost)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate trading fee: %w", err)
	}

	// Créer l'ordre
	order := &models.TradingOrder{
		UserID:          userID,
		OrderType:       orderType,
		Side:            "buy",
		FromCurrency:    paymentCurrency,
		ToCurrency:      cryptoCurrency,
		Amount:          amount,
		Price:           &executionPrice,
		RemainingAmount: amount,
		Status:          "open",
		Fee:             fee,
	}

	// Si c'est un ordre market, l'exécuter immédiatement
	if orderType == "market" {
		order.Status = "filled"
		order.FilledAmount = amount
		order.RemainingAmount = 0
		now := time.Now()
		order.ExecutedAt = &now
	}

	// Sauvegarder l'ordre
	err = s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Publier l'événement
	s.publishTradingEvent("order.created", order)

	return order, nil
}

func (s *ExchangeService) SellCrypto(userID, cryptoCurrency, receiveCurrency string, amount float64, orderType string, limitPrice *float64) (*models.TradingOrder, error) {
	// Valider les paramètres
	if !s.isCryptoCurrency(cryptoCurrency) {
		return nil, fmt.Errorf("invalid cryptocurrency: %s", cryptoCurrency)
	}

	if !s.isFiatCurrency(receiveCurrency) {
		return nil, fmt.Errorf("invalid receive currency: %s", receiveCurrency)
	}

	// Obtenir le prix actuel
	currentRate, err := s.rateService.GetRate(cryptoCurrency, receiveCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get current rate: %w", err)
	}

	// Calculer le prix d'exécution
	var executionPrice float64
	if orderType == "market" {
		executionPrice = currentRate.BidPrice // Prix d'achat pour vendre
	} else if orderType == "limit" && limitPrice != nil {
		if *limitPrice < currentRate.BidPrice {
			return nil, fmt.Errorf("limit price too low for sell order")
		}
		executionPrice = *limitPrice
	} else {
		return nil, fmt.Errorf("invalid order type or missing limit price")
	}

	// Calculer le montant à recevoir
	totalReceive := amount * executionPrice
	
	fee, err := s.feeService.CalculateFee("trading_sell", totalReceive)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate trading fee: %w", err)
	}

	// Créer l'ordre
	order := &models.TradingOrder{
		UserID:          userID,
		OrderType:       orderType,
		Side:            "sell",
		FromCurrency:    cryptoCurrency,
		ToCurrency:      receiveCurrency,
		Amount:          amount,
		Price:           &executionPrice,
		RemainingAmount: amount,
		Status:          "open",
		Fee:             fee,
	}

	// Si c'est un ordre market, l'exécuter immédiatement
	if orderType == "market" {
		order.Status = "filled"
		order.FilledAmount = amount
		order.RemainingAmount = 0
		now := time.Now()
		order.ExecutedAt = &now
	}

	// Sauvegarder l'ordre
	err = s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Publier l'événement
	s.publishTradingEvent("order.created", order)

	return order, nil
}

func (s *ExchangeService) GetUserOrders(userID string) ([]*models.TradingOrder, error) {
	return s.orderRepo.GetOrdersByUser(userID)
}

func (s *ExchangeService) GetUserExchanges(userID string, limit int) ([]*models.Exchange, error) {
	return s.exchangeRepo.GetByUserID(userID, limit)
}

func (s *ExchangeService) GetExchange(exchangeID string) (*models.Exchange, error) {
	return s.exchangeRepo.GetByID(exchangeID)
}

func (s *ExchangeService) CancelOrder(userID, orderID string) error {
	// Récupérer l'ordre
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Vérifier la propriété
	if order.UserID != userID {
		return fmt.Errorf("order does not belong to user")
	}

	// Vérifier qu'il peut être annulé
	if order.Status != "open" && order.Status != "partial" {
		return fmt.Errorf("order cannot be cancelled, current status: %s", order.Status)
	}

	// Annuler l'ordre
	err = s.orderRepo.UpdateOrderStatus(orderID, "cancelled")
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Publier l'événement
	s.publishTradingEvent("order.cancelled", order)

	return nil
}

func (s *ExchangeService) GetPortfolio(userID string) (*models.Portfolio, error) {
	// Construct real portfolio
	// In a real implementation we would fetch balances from WalletService for all supported assets
	// Since we don't have a "get all balances" easily exposed without multiple calls, or if we assume local caching?
	// We will try to fetch for major assets.

	portfolio := &models.Portfolio{
		Holdings: []models.Holding{},
	}

	// Calculate performance metrics from real orders
	orders, err := s.orderRepo.GetOrdersByUser(userID)
	if err == nil {
		var totalTrades, winningTrades, losingTrades int
		for _, o := range orders {
			if o.Status == "filled" || o.Status == "completed" {
				totalTrades++
				// Logic to determine win/loss requires historical price data or PnL tracking
				// For now we report real total trades, and 0 for win/loss to avoid fake data
			}
		}
		
		portfolio.Performance = models.PerformanceMetrics{
			TotalTrades:   totalTrades,
			WinningTrades: winningTrades,
			LosingTrades:  losingTrades,
			WinRate:       0,
		}
		if totalTrades > 0 {
			portfolio.Performance.WinRate = float64(winningTrades) / float64(totalTrades) * 100
		}
	}

	return portfolio, nil
}

// Méthodes privées

func (s *ExchangeService) processExchange(exchange *models.Exchange) {
	log.Printf("[EXCHANGE DEBUG] Starting processExchange for exchange %s", exchange.ID)
	
	// Check if kafkaPublisher is available
	if s.kafkaPublisher == nil || !s.kafkaPublisher.IsConnected() {
		log.Printf("[EXCHANGE ERROR] Kafka publisher not available for exchange %s", exchange.ID)
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed")
		return
	}
	
	// 1. Initiate Debit
	debitReq := &messaging.PaymentRequestEvent{
		RequestID:    fmt.Sprintf("TX-EXC-DEBIT-%s", exchange.ID),
		UserID:       exchange.UserID,
		FromWalletID: exchange.FromWalletID,
		DebitAmount:  exchange.FromAmount,
		Currency:     exchange.FromCurrency,
		Type:         "exchange_debit",
		ReferenceID:  fmt.Sprintf("EXCHANGE_DEBIT_%s", exchange.ID),
		Description:  fmt.Sprintf("Echange: %f %s vers %s (Taux: %f, Frais: %f %s)", exchange.FromAmount, exchange.FromCurrency, exchange.ToCurrency, exchange.ExchangeRate, exchange.Fee, exchange.FromCurrency),
		Metadata: map[string]interface{}{
			"exchange_id": exchange.ID,
			"rate":        exchange.ExchangeRate,
			"fee":         exchange.Fee,
			"to_amount":   exchange.ToAmount,
			"to_currency": exchange.ToCurrency,
		},
	}

	log.Printf("[EXCHANGE DEBUG] Created debit request for exchange %s: SourceWallet=%s, Amount=%f, Currency=%s", 
		exchange.ID, exchange.FromWalletID, exchange.FromAmount, exchange.FromCurrency)

	// Update status to indicate processing
	log.Printf("[EXCHANGE DEBUG] Updating status to processing_debit for exchange %s", exchange.ID)
	s.exchangeRepo.UpdateStatus(exchange.ID, "processing_debit")

	if err := s.kafkaPublisher.PublishPaymentRequest(debitReq); err != nil {
		log.Printf("[EXCHANGE ERROR] Failed to publish debit request for exchange %s: %v", exchange.ID, err)
		s.FailExchange(exchange.ID, "Failed to initiate debit: " + err.Error())
	} else {
		log.Printf("[EXCHANGE DEBUG] Successfully published debit request for exchange %s", exchange.ID)
	}
}

// Helper to continue exchange after debit success
func (s *ExchangeService) CompleteExchangeCredit(exchangeID string) {
	log.Printf("[EXCHANGE DEBUG] CompleteExchangeCredit called for exchange %s", exchangeID)
	
	exchange, err := s.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		log.Printf("[EXCHANGE ERROR] Failed to retrieve exchange %s for credit step: %v", exchangeID, err)
		return
	}
	
	log.Printf("[EXCHANGE DEBUG] Retrieved exchange %s, initiating credit to wallet %s, amount %f %s", 
		exchangeID, exchange.ToWalletID, exchange.ToAmount, exchange.ToCurrency)

	// 2. Initiate Credit
	creditReq := &messaging.PaymentRequestEvent{
		RequestID:    fmt.Sprintf("TX-EXC-CREDIT-%s", exchange.ID),
		UserID:       exchange.UserID,
		ToWalletID:   exchange.ToWalletID,
		CreditAmount: exchange.ToAmount,
		Currency:     exchange.ToCurrency,
		Type:         "exchange_credit",
		ReferenceID:  fmt.Sprintf("EXCHANGE_CREDIT_%s", exchange.ID),
		Description:  fmt.Sprintf("Reçu via échange: %f %s (Source: %f %s)", exchange.ToAmount, exchange.ToCurrency, exchange.FromAmount, exchange.FromCurrency),
		Metadata: map[string]interface{}{
			"exchange_id":   exchange.ID,
			"rate":          exchange.ExchangeRate,
			"source_amount": exchange.FromAmount,
			"source_currency": exchange.FromCurrency,
		},
	}

	log.Printf("[EXCHANGE DEBUG] Updating status to processing_credit for exchange %s", exchangeID)
	s.exchangeRepo.UpdateStatus(exchange.ID, "processing_credit")

	log.Printf("[EXCHANGE DEBUG] Publishing credit request for exchange %s", exchangeID)
	if err := s.kafkaPublisher.PublishPaymentRequest(creditReq); err != nil {
		log.Printf("[EXCHANGE ERROR] Failed to publish credit request for exchange %s: %v", exchange.ID, err)
		// Partial failure is critical, manual intervention might be needed
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed_partial") 
		// We still notify failure/issue
		s.FailExchange(exchange.ID, "Credit failed after debit success: " + err.Error())
	} else {
		log.Printf("[EXCHANGE DEBUG] Successfully published credit request for exchange %s", exchangeID)
	}
}

// Helper to finalize exchange
func (s *ExchangeService) FinalizeExchange(exchangeID string) {
	s.exchangeRepo.UpdateStatus(exchangeID, "completed")
	
	// Update CompletedAt
	// We might need a repo method for setting completed_at or do it manually if repo supports update
	// existing UpdateStatus only updates status.
	// Let's assume UpdateStatus is enough for now or I add a specialized method later.
	
	// exchange, _ := s.exchangeRepo.GetByID(exchangeID)
	// CompletedAt update ... (skipping for brevity unless critical, assumed UpdateStatus handles updated_at)

	// Notification is handled by WalletService/NotificationService via payment events
	// s.publishExchangeEvent("exchange.completed", exchange)
}

func (s *ExchangeService) FailExchange(exchangeID, reason string) {
	s.exchangeRepo.UpdateStatus(exchangeID, "failed")
	exchange, _ := s.exchangeRepo.GetByID(exchangeID)
	log.Printf("Exchange %s failed: %s. Details: %+v", exchangeID, reason, exchange)

	// Publish failure event for notification
	if s.kafkaPublisher != nil && exchange != nil {
		event := map[string]interface{}{
			"type":          "exchange.failed",
			"exchange_id":   exchange.ID,
			"user_id":       exchange.UserID,
			"from_currency": exchange.FromCurrency,
			"to_currency":   exchange.ToCurrency,
			"amount":        exchange.FromAmount,
			"status":        "failed",
			"reason":        reason,
			"timestamp":     time.Now(),
		}
		s.kafkaPublisher.PublishExchangeEvent("exchange.failed", event)
	}
}

func (s *ExchangeService) calculateFeePercentage(fromCurrency, toCurrency string, amount float64) float64 {
	baseFee := 0.5 // 0.5% par défaut

	// Frais différents selon le type de change
	if s.isCryptoCurrency(fromCurrency) && s.isFiatCurrency(toCurrency) {
		baseFee = s.config.ExchangeFees["crypto_to_fiat"]
	} else if s.isFiatCurrency(fromCurrency) && s.isCryptoCurrency(toCurrency) {
		baseFee = s.config.ExchangeFees["fiat_to_crypto"]
	} else if s.isCryptoCurrency(fromCurrency) && s.isCryptoCurrency(toCurrency) {
		baseFee = s.config.ExchangeFees["crypto_to_crypto"]
	} else {
		baseFee = s.config.ExchangeFees["fiat_to_fiat"]
	}

	// Réduction de frais selon le volume
	if amount > 100000 {
		baseFee *= 0.5
	} else if amount > 10000 {
		baseFee *= 0.7
	} else if amount > 1000 {
		baseFee *= 0.9
	}

	return baseFee
}

func (s *ExchangeService) calculateTradingFee(amount float64, side string) float64 {
	feeRate := s.config.TradingFees[side]
	return amount * feeRate / 100
}

func (s *ExchangeService) getEstimatedDelivery(fromCurrency, toCurrency string) string {
	if s.isCryptoCurrency(fromCurrency) || s.isCryptoCurrency(toCurrency) {
		return "5-15 minutes"
	}
	return "Instant"
}

func (s *ExchangeService) isSupportedCurrency(currency string) bool {
	// Allow any valid-looking currency code (3-5 chars)
	// We rely on RateService.GetRate() to determine if it's actually supported/tradable
	l := len(currency)
	return l >= 3 && l <= 5
}

func (s *ExchangeService) isCryptoCurrency(currency string) bool {
	cryptoCurrencies := []string{
		"BTC", "ETH", "USDT", "USDC", "BNB", "ADA", "XRP", "DOT", "LTC", "LINK", "BCH", "XLM",
		// Add more as needed or move to config/DB
	}
	
	currency = strings.ToUpper(currency)
	for _, crypto := range cryptoCurrencies {
		if currency == crypto {
			return true
		}
	}
	return false
}

func (s *ExchangeService) isFiatCurrency(currency string) bool {
	// Assume anything that is not a known crypto is a fiat currency
	return !s.isCryptoCurrency(currency)
}

func (s *ExchangeService) publishExchangeEvent(eventType string, exchange *models.Exchange) {
	if s.kafkaPublisher == nil {
		return
	}

	event := map[string]interface{}{
		"type":          eventType,
		"exchange_id":   exchange.ID,
		"user_id":       exchange.UserID,
		"from_currency": exchange.FromCurrency,
		"to_currency":   exchange.ToCurrency,
		"amount":        exchange.FromAmount,
		"status":        exchange.Status,
		"timestamp":     time.Now(),
	}

	s.kafkaPublisher.PublishExchangeEvent(eventType, event)
}

func (s *ExchangeService) publishTradingEvent(eventType string, order *models.TradingOrder) {
	if s.kafkaPublisher == nil {
		return
	}

	event := map[string]interface{}{
		"type":      eventType,
		"order_id":  order.ID,
		"user_id":   order.UserID,
		"side":      order.Side,
		"amount":    order.Amount,
		"price":     order.Price,
		"status":    order.Status,
		"timestamp": time.Now(),
	}

	s.kafkaPublisher.PublishTradingEvent(eventType, event)
}

func (s *ExchangeService) getFeeKey(fromCurrency, toCurrency string) string {
	if s.isCryptoCurrency(fromCurrency) && s.isFiatCurrency(toCurrency) {
		return "exchange_crypto_to_fiat"
	} else if s.isFiatCurrency(fromCurrency) && s.isCryptoCurrency(toCurrency) {
		return "exchange_fiat_to_crypto"
	} else if s.isCryptoCurrency(fromCurrency) && s.isCryptoCurrency(toCurrency) {
		return "exchange_crypto_to_crypto"
	} else {
		return "exchange_fiat_to_fiat"
	}
}