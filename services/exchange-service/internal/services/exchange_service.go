package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
	"github.com/streadway/amqp"
)

type ExchangeService struct {
	exchangeRepo *repository.ExchangeRepository
	orderRepo    *repository.OrderRepository
	rateService  *RateService
	config       *config.Config
	mqChannel    *amqp.Channel
	walletClient *WalletClient
}

func NewExchangeService(exchangeRepo *repository.ExchangeRepository, orderRepo *repository.OrderRepository, rateService *RateService, mqChannel *amqp.Channel, walletClient *WalletClient, cfg *config.Config) *ExchangeService {
	return &ExchangeService{
		exchangeRepo: exchangeRepo,
		orderRepo:    orderRepo,
		rateService:  rateService,
		config:       cfg,
		mqChannel:    mqChannel,
		walletClient: walletClient,
	}
}

func (s *ExchangeService) GetQuote(userID, fromCurrency, toCurrency string, fromAmount *float64, toAmount *float64) (*models.Quote, error) {
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
	feePercentage := s.calculateFeePercentage(fromCurrency, toCurrency, calculatedFromAmount)
	fee := calculatedFromAmount * feePercentage / 100

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

	// Sauvegarder le devis
	err = s.exchangeRepo.CreateQuote(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

func (s *ExchangeService) ExecuteExchange(userID, quoteID, fromWalletID, toWalletID string) (*models.Exchange, error) {
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

	// Calculate and check fees
	// In a real system, we would lock funds here or check balance again
	// We rely on the async processor to do the actual transfer and fail if insufficient funds

	// Créer l'échange
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
		Status:           "pending",
		QuoteID:          quoteID,
	}

	// Sauvegarder l'échange
	err = s.exchangeRepo.Create(exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err)
	}

	// Traiter l'échange de manière asynchrone
	go s.processExchange(exchange)

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
	fee := s.calculateTradingFee(totalCost, "buy")

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
	fee := s.calculateTradingFee(totalReceive, "sell")

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
	// Remove artificial delay
	// time.Sleep(2 * time.Second)

	// Execute transaction via Wallet Service
	// Debit Sender
	debitReq := &TransactionRequest{
		UserID:    exchange.UserID,
		WalletID:  exchange.FromWalletID,
		Amount:    exchange.FromAmount,
		Type:      "debit",
		Currency:  exchange.FromCurrency,
		Reference: "EXCHANGE_DEBIT_" + exchange.ID,
	}
	
	if err := s.walletClient.ProcessTransaction(debitReq); err != nil {
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed")
		return
	}

	// Credit Receiver
	creditReq := &TransactionRequest{
		UserID:    exchange.UserID,
		WalletID:  exchange.ToWalletID,
		Amount:    exchange.ToAmount,
		Type:      "credit",
		Currency:  exchange.ToCurrency,
		Reference: "EXCHANGE_CREDIT_" + exchange.ID,
	}

	if err := s.walletClient.ProcessTransaction(creditReq); err != nil {
		// Ideally we should rollback debit here
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed_partial") // Manual intervention needed
		return
	}

	// Marquer comme terminé
	s.exchangeRepo.UpdateStatus(exchange.ID, "completed")
	
	now := time.Now()
	exchange.CompletedAt = &now
	exchange.Status = "completed"

	// Publier l'événement
	s.publishExchangeEvent("exchange.completed", exchange)
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
	supportedCurrencies := []string{
		"USD", "EUR", "GBP", "CAD", "AUD", "JPY", "CHF",
		"BTC", "ETH", "USDT", "USDC", "BNB", "ADA", "XRP", "DOT", "LTC",
	}
	
	currency = strings.ToUpper(currency)
	for _, supported := range supportedCurrencies {
		if currency == supported {
			return true
		}
	}
	return false
}

func (s *ExchangeService) isCryptoCurrency(currency string) bool {
	cryptoCurrencies := []string{
		"BTC", "ETH", "USDT", "USDC", "BNB", "ADA", "XRP", "DOT", "LTC", "LINK", "BCH", "XLM",
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
	fiatCurrencies := []string{
		"USD", "EUR", "GBP", "CAD", "AUD", "JPY", "CHF", "SEK", "NOK", "DKK",
	}
	
	currency = strings.ToUpper(currency)
	for _, fiat := range fiatCurrencies {
		if currency == fiat {
			return true
		}
	}
	return false
}

func (s *ExchangeService) publishExchangeEvent(eventType string, exchange *models.Exchange) {
	if s.mqChannel == nil {
		return
	}

	event := map[string]interface{}{
		"type":        eventType,
		"exchange_id": exchange.ID,
		"user_id":     exchange.UserID,
		"from_currency": exchange.FromCurrency,
		"to_currency": exchange.ToCurrency,
		"amount":      exchange.FromAmount,
		"status":      exchange.Status,
		"timestamp":   time.Now(),
	}

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		"exchange.events",
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}

func (s *ExchangeService) publishTradingEvent(eventType string, order *models.TradingOrder) {
	if s.mqChannel == nil {
		return
	}

	event := map[string]interface{}{
		"type":     eventType,
		"order_id": order.ID,
		"user_id":  order.UserID,
		"side":     order.Side,
		"amount":   order.Amount,
		"price":    order.Price,
		"status":   order.Status,
		"timestamp": time.Now(),
	}

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		"trading.events",
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}