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

type FiatExchangeService struct {
	exchangeRepo   *repository.ExchangeRepository
	rateService    *RateService
	feeService     *FeeService
	config         *config.Config
	kafkaPublisher *KafkaPublisher
}

func NewFiatExchangeService(exchangeRepo *repository.ExchangeRepository, rateService *RateService, feeService *FeeService, kafkaPublisher *KafkaPublisher, cfg *config.Config) *FiatExchangeService {
	return &FiatExchangeService{
		exchangeRepo:   exchangeRepo,
		rateService:    rateService,
		feeService:     feeService,
		config:         cfg,
		kafkaPublisher: kafkaPublisher,
	}
}

// Conversion entre monnaies fiduciaires (USD, EUR, GBP, CAD, AUD, JPY, CHF, etc.)
func (s *FiatExchangeService) ConvertFiat(userID, fromCurrency, toCurrency string, amount float64, fromWalletID, toWalletID string) (*models.Exchange, error) {
	// Validation des devises fiduciaires
	if !s.isFiatCurrency(fromCurrency) {
		return nil, fmt.Errorf("source currency %s is not a supported fiat currency", fromCurrency)
	}

	if !s.isFiatCurrency(toCurrency) {
		return nil, fmt.Errorf("target currency %s is not a supported fiat currency", toCurrency)
	}

	if fromCurrency == toCurrency {
		return nil, fmt.Errorf("cannot convert same currency")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Obtenir le taux de change fiat
	exchangeRate, err := s.getFiatExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Calculer le montant de destination
	convertedAmount := amount * exchangeRate.Rate

	// Calculer les frais (plus faibles pour fiat-to-fiat)
	fee, err := s.feeService.CalculateFee("exchange_fiat_to_fiat", amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee: %w", err)
	}
	
	// Calculate effective percentage for record keeping
	feePercentage := (fee / amount) * 100

	finalAmount := convertedAmount - (fee * exchangeRate.Rate)

	// Créer la transaction d'échange
	exchange := &models.Exchange{
		UserID:              userID,
		FromWalletID:        fromWalletID,
		ToWalletID:          toWalletID,
		FromCurrency:        strings.ToUpper(fromCurrency),
		ToCurrency:          strings.ToUpper(toCurrency),
		FromAmount:          amount,
		ToAmount:            finalAmount,
		ExchangeRate:        exchangeRate.Rate,
		Fee:                 fee,
		FeePercentage:       feePercentage,
		Status:              "pending",
	}

	// Sauvegarder l'échange
	err = s.exchangeRepo.Create(exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to create fiat exchange: %w", err)
	}

	// Traiter l'échange (plus rapide pour fiat-to-fiat)
	go s.processFiatExchange(exchange)

	return exchange, nil
}

// Obtenir un devis pour conversion fiat
func (s *FiatExchangeService) GetFiatQuote(userID, fromCurrency, toCurrency string, amount float64) (*models.Quote, error) {
	if !s.isFiatCurrency(fromCurrency) || !s.isFiatCurrency(toCurrency) {
		return nil, fmt.Errorf("both currencies must be fiat currencies")
	}

	if fromCurrency == toCurrency {
		return nil, fmt.Errorf("cannot quote same currency conversion")
	}

	// Obtenir le taux actuel
	rate, err := s.getFiatExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	// Calculer les montants
	convertedAmount := amount * rate.Rate
	feePercentage := s.CalculateFiatExchangeFee(fromCurrency, toCurrency, amount)
	fee := amount * feePercentage / 100
	finalAmount := convertedAmount - (fee * rate.Rate)

	quote := &models.Quote{
		UserID:            userID,
		FromCurrency:      strings.ToUpper(fromCurrency),
		ToCurrency:        strings.ToUpper(toCurrency),
		FromAmount:        amount,
		ToAmount:          finalAmount,
		ExchangeRate:      rate.Rate,
		Fee:               fee,
		FeePercentage:     feePercentage,
		ValidUntil:        time.Now().Add(2 * time.Minute), // Plus court pour fiat
		EstimatedDelivery: "Instant", // Fiat-to-fiat est instantané
	}

	err = s.exchangeRepo.CreateQuote(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

// Obtenir les taux de change fiat en temps réel
func (s *FiatExchangeService) getFiatExchangeRate(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	// Essayer d'abord le cache Redis
	rate, err := s.rateService.GetRate(fromCurrency, toCurrency)
	if err == nil && time.Since(rate.LastUpdated) < 5*time.Minute {
		return rate, nil
	}

	// Obtenir les taux en temps réel depuis une API externe
	freshRate, err := s.fetchLiveFiatRate(fromCurrency, toCurrency)
	if err != nil {
		// Fallback vers les taux en cache si l'API externe échoue
		if rate != nil {
			return rate, nil
		}
		return nil, fmt.Errorf("failed to get live rate and no cached rate available: %w", err)
	}

	// Mettre à jour le cache
	s.rateService.UpdateRate(freshRate)

	return freshRate, nil
}

// Récupérer les taux de change en temps réel (simulation)
func (s *FiatExchangeService) fetchLiveFiatRate(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	// Dans un vrai système, ici on ferait appel à une API comme:
	// - Fixer.io, CurrencyLayer, OpenExchangeRates, etc.
	
	// Taux de change simulés mais réalistes
	rates := map[string]map[string]float64{
		"USD": {
			"EUR": 0.8456,
			"GBP": 0.7821,
			"JPY": 149.23,
			"CAD": 1.3567,
			"AUD": 1.5234,
			"CHF": 0.8912,
			"SEK": 10.8765,
			"NOK": 10.9876,
			"XOF": 600.0,
			"XAF": 600.0,
			"NGN": 1500.0,
			"GHS": 12.5,
			"KES": 153.0,
			"ZAR": 18.5,
			"MAD": 10.0,
		},
		"EUR": {
			"USD": 1.1826,
			"GBP": 0.9248,
			"JPY": 176.45,
			"CAD": 1.6043,
			"AUD": 1.8012,
			"CHF": 1.0539,
			"XOF": 655.957,
			"XAF": 655.957,
			"NGN": 1780.0,
			"GHS": 14.8,
			"KES": 181.0,
			"ZAR": 21.9,
			"MAD": 10.85,
		},
		"GBP": {
			"USD": 1.2787,
			"EUR": 1.0814,
			"JPY": 190.76,
			"CAD": 1.7354,
			"AUD": 1.9467,
			"XOF": 765.0,
			"XAF": 765.0,
			"NGN": 1920.0,
		},
		"JPY": {
			"USD": 0.0067,
			"EUR": 0.0057,
			"GBP": 0.0052,
			"CAD": 0.0091,
			"AUD": 0.0102,
			"XOF": 4.02,
			"XAF": 4.02,
		},
		"CAD": {
			"USD": 0.7371,
			"EUR": 0.6232,
			"GBP": 0.5764,
			"JPY": 109.98,
			"AUD": 1.1228,
			"XOF": 442.0,
			"XAF": 442.0,
		},
		"AUD": {
			"USD": 0.6564,
			"EUR": 0.5551,
			"GBP": 0.5137,
			"JPY": 97.96,
			"CAD": 0.8906,
			"XOF": 394.0,
			"XAF": 394.0,
		},
		// African currencies
		"XOF": {
			"USD": 0.00167,
			"EUR": 0.001524,
			"GBP": 0.001307,
			"XAF": 1.0,
			"NGN": 2.5,
			"GHS": 0.0208,
			"KES": 0.255,
			"ZAR": 0.0308,
		},
		"XAF": {
			"USD": 0.00167,
			"EUR": 0.001524,
			"GBP": 0.001307,
			"XOF": 1.0,
			"NGN": 2.5,
			"GHS": 0.0208,
			"KES": 0.255,
			"ZAR": 0.0308,
		},
		"NGN": {
			"USD": 0.000667,
			"EUR": 0.000562,
			"XOF": 0.4,
			"XAF": 0.4,
			"GHS": 0.00833,
			"KES": 0.102,
			"ZAR": 0.0123,
		},
		"GHS": {
			"USD": 0.08,
			"EUR": 0.0676,
			"XOF": 48.0,
			"XAF": 48.0,
			"NGN": 120.0,
			"KES": 12.24,
			"ZAR": 1.48,
		},
		"KES": {
			"USD": 0.00654,
			"EUR": 0.00553,
			"XOF": 3.92,
			"XAF": 3.92,
			"NGN": 9.8,
			"GHS": 0.0817,
			"ZAR": 0.121,
		},
		"ZAR": {
			"USD": 0.054,
			"EUR": 0.0456,
			"XOF": 32.4,
			"XAF": 32.4,
			"NGN": 81.0,
			"GHS": 0.676,
			"KES": 8.27,
		},
		"MAD": {
			"USD": 0.1,
			"EUR": 0.092,
			"XOF": 60.0,
			"XAF": 60.0,
		},
	}

	fromRates, exists := rates[fromCurrency]
	if !exists {
		return nil, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}

	rate, exists := fromRates[toCurrency]
	if !exists {
		return nil, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	// Ajouter une petite variation aléatoire (+/- 0.1%)
	variation := (0.999 + (0.002 * float64(time.Now().UnixNano()%1000)/1000))
	rate = rate * variation

	// Calculer bid/ask spread (0.1% spread typique pour fiat)
	spread := rate * 0.001
	bidPrice := rate - spread/2
	askPrice := rate + spread/2

	exchangeRate := &models.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
		BidPrice:     bidPrice,
		AskPrice:     askPrice,
		Spread:       spread,
		Source:       "fiat_exchange_api",
		Volume24h:    1000000.0, // Volume simulé
		Change24h:    -0.05 + (0.1 * float64(time.Now().UnixNano()%1000)/1000), // +/- 5% max
		LastUpdated:  time.Now(),
	}

	return exchangeRate, nil
}

// Calculer les frais pour échange fiat (plus bas que crypto)
func (s *FiatExchangeService) CalculateFiatExchangeFee(fromCurrency, toCurrency string, amount float64) float64 {
	baseFee := 0.25 // 0.25% pour fiat-to-fiat (moins cher que crypto)

	// Frais préférentiels pour certaines paires populaires
	majorPairs := map[string]bool{
		"USD-EUR": true, "EUR-USD": true,
		"USD-GBP": true, "GBP-USD": true,
		"EUR-GBP": true, "GBP-EUR": true,
	}

	pair := fromCurrency + "-" + toCurrency
	if majorPairs[pair] {
		baseFee = 0.15 // 0.15% pour paires majeures
	}

	// Réductions selon le volume
	if amount >= 50000 {
		baseFee *= 0.5 // 50% de réduction pour gros volumes
	} else if amount >= 10000 {
		baseFee *= 0.7 // 30% de réduction
	} else if amount >= 1000 {
		baseFee *= 0.9 // 10% de réduction
	}

	return baseFee
}

// Vérifier si c'est une devise fiduciaire supportée
func (s *FiatExchangeService) isFiatCurrency(currency string) bool {
	fiatCurrencies := map[string]bool{
		"USD": true, // Dollar américain
		"EUR": true, // Euro
		"GBP": true, // Livre sterling
		"JPY": true, // Yen japonais
		"CAD": true, // Dollar canadien
		"AUD": true, // Dollar australien
		"CHF": true, // Franc suisse
		"SEK": true, // Couronne suédoise
		"NOK": true, // Couronne norvégienne
		"DKK": true, // Couronne danoise
		"PLN": true, // Zloty polonais
		"CZK": true, // Couronne tchèque
		"HUF": true, // Forint hongrois
		"SGD": true, // Dollar de Singapour
		"HKD": true, // Dollar de Hong Kong
		"NZD": true, // Dollar néo-zélandais
		"MXN": true, // Peso mexicain
		"BRL": true, // Real brésilien
		"RUB": true, // Rouble russe
		"CNY": true, // Yuan chinois
		"INR": true, // Roupie indienne
		"KRW": true, // Won sud-coréen
		"TRY": true, // Livre turque
		"ZAR": true, // Rand sud-africain
		// African currencies
		"XOF": true, // CFA Franc BCEAO (West Africa)
		"XAF": true, // CFA Franc BEAC (Central Africa)
		"NGN": true, // Nigerian Naira
		"GHS": true, // Ghanaian Cedi
		"KES": true, // Kenyan Shilling
		"MAD": true, // Moroccan Dirham
		"EGP": true, // Egyptian Pound
		"TZS": true, // Tanzanian Shilling
		"UGX": true, // Ugandan Shilling
	}

	return fiatCurrencies[strings.ToUpper(currency)]
}

// Obtenir toutes les paires fiat supportées
func (s *FiatExchangeService) GetSupportedFiatPairs() map[string][]string {
	return map[string][]string{
		"USD": {"EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "SEK", "NOK"},
		"EUR": {"USD", "GBP", "JPY", "CAD", "AUD", "CHF", "SEK", "NOK"},
		"GBP": {"USD", "EUR", "JPY", "CAD", "AUD", "CHF"},
		"JPY": {"USD", "EUR", "GBP", "CAD", "AUD"},
		"CAD": {"USD", "EUR", "GBP", "JPY", "AUD"},
		"AUD": {"USD", "EUR", "GBP", "JPY", "CAD"},
		"CHF": {"USD", "EUR", "GBP", "JPY"},
		"SEK": {"USD", "EUR", "NOK"},
		"NOK": {"USD", "EUR", "SEK"},
	}
}

// processFiatExchange handles the fiat exchange via Kafka payment events
func (s *FiatExchangeService) processFiatExchange(exchange *models.Exchange) {
	// 1. Initiate Debit via Kafka payment event
	debitReq := &messaging.PaymentRequestEvent{
		RequestID:    fmt.Sprintf("TX-FIAT-DEBIT-%s", exchange.ID),
		FromWalletID: exchange.FromWalletID,
		DebitAmount:  exchange.FromAmount,
		Currency:     exchange.FromCurrency,
		Type:         "fiat_exchange_debit",
		ReferenceID:  fmt.Sprintf("FIAT_EXCHANGE_DEBIT_%s", exchange.ID),
	}

	// Update status to indicate processing
	s.exchangeRepo.UpdateStatus(exchange.ID, "processing_debit")

	if s.kafkaPublisher != nil && s.kafkaPublisher.IsConnected() {
		if err := s.kafkaPublisher.PublishPaymentRequest(debitReq); err != nil {
			log.Printf("Failed to publish debit request for fiat exchange %s: %v", exchange.ID, err)
			s.exchangeRepo.UpdateStatus(exchange.ID, "failed")
		} else {
			log.Printf("[FIAT EXCHANGE] Successfully published debit request for %s", exchange.ID)
		}
	} else {
		log.Printf("Warning: Kafka publisher not available for fiat exchange %s, marking as failed", exchange.ID)
		s.exchangeRepo.UpdateStatus(exchange.ID, "failed")
	}
}

// CompleteFiatExchangeCredit continues the exchange after debit success
func (s *FiatExchangeService) CompleteFiatExchangeCredit(exchangeID string) {
	exchange, err := s.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		log.Printf("Failed to retrieve fiat exchange %s for credit step: %v", exchangeID, err)
		return
	}

	// 2. Initiate Credit via Kafka
	creditReq := &messaging.PaymentRequestEvent{
		RequestID:    fmt.Sprintf("TX-FIAT-CREDIT-%s", exchange.ID),
		ToWalletID:   exchange.ToWalletID,
		CreditAmount: exchange.ToAmount,
		Currency:     exchange.ToCurrency,
		Type:         "fiat_exchange_credit",
		ReferenceID:  fmt.Sprintf("FIAT_EXCHANGE_CREDIT_%s", exchange.ID),
	}

	s.exchangeRepo.UpdateStatus(exchange.ID, "processing_credit")

	if s.kafkaPublisher != nil && s.kafkaPublisher.IsConnected() {
		if err := s.kafkaPublisher.PublishPaymentRequest(creditReq); err != nil {
			log.Printf("Failed to publish credit request for fiat exchange %s: %v", exchange.ID, err)
			s.exchangeRepo.UpdateStatus(exchange.ID, "failed_partial")
		} else {
			log.Printf("[FIAT EXCHANGE] Successfully published credit request for %s", exchange.ID)
		}
	}
}

// FinalizeFiatExchange marks the exchange as completed
func (s *FiatExchangeService) FinalizeFiatExchange(exchangeID string) {
	s.exchangeRepo.UpdateStatus(exchangeID, "completed")

	exchange, _ := s.exchangeRepo.GetByID(exchangeID)
	if exchange != nil {
		now := time.Now()
		exchange.CompletedAt = &now
		exchange.Status = "completed"
		s.publishFiatExchangeEvent("fiat_exchange.completed", exchange)
	}
}

// FailFiatExchange marks the exchange as failed
func (s *FiatExchangeService) FailFiatExchange(exchangeID, reason string) {
	s.exchangeRepo.UpdateStatus(exchangeID, "failed")
	exchange, _ := s.exchangeRepo.GetByID(exchangeID)
	log.Printf("Fiat Exchange %s failed: %s. Details: %+v", exchangeID, reason, exchange)
}

func (s *FiatExchangeService) publishFiatExchangeEvent(eventType string, exchange *models.Exchange) {
	if s.kafkaPublisher == nil {
		return
	}

	event := map[string]interface{}{
		"type":             eventType,
		"exchange_id":      exchange.ID,
		"user_id":          exchange.UserID,
		"from_currency":    exchange.FromCurrency,
		"to_currency":      exchange.ToCurrency,
		"amount":           exchange.FromAmount,
		"converted_amount": exchange.ToAmount,
		"fee":              exchange.Fee,
		"exchange_rate":    exchange.ExchangeRate,
		"status":           exchange.Status,
		"timestamp":        time.Now(),
	}

	s.kafkaPublisher.PublishExchangeEvent(eventType, event)
}