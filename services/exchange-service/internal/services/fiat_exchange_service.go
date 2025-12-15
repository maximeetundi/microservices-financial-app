package services

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
	"github.com/streadway/amqp"
)

type FiatExchangeService struct {
	exchangeRepo *repository.ExchangeRepository
	rateService  *RateService
	config       *config.Config
	mqChannel    *amqp.Channel
}

func NewFiatExchangeService(exchangeRepo *repository.ExchangeRepository, rateService *RateService, mqChannel *amqp.Channel, cfg *config.Config) *FiatExchangeService {
	return &FiatExchangeService{
		exchangeRepo: exchangeRepo,
		rateService:  rateService,
		config:       cfg,
		mqChannel:    mqChannel,
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
	feePercentage := s.calculateFiatExchangeFee(fromCurrency, toCurrency, amount)
	fee := amount * feePercentage / 100
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
	feePercentage := s.calculateFiatExchangeFee(fromCurrency, toCurrency, amount)
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
		},
		"EUR": {
			"USD": 1.1826,
			"GBP": 0.9248,
			"JPY": 176.45,
			"CAD": 1.6043,
			"AUD": 1.8012,
			"CHF": 1.0539,
		},
		"GBP": {
			"USD": 1.2787,
			"EUR": 1.0814,
			"JPY": 190.76,
			"CAD": 1.7354,
			"AUD": 1.9467,
		},
		"JPY": {
			"USD": 0.0067,
			"EUR": 0.0057,
			"GBP": 0.0052,
			"CAD": 0.0091,
			"AUD": 0.0102,
		},
		"CAD": {
			"USD": 0.7371,
			"EUR": 0.6232,
			"GBP": 0.5764,
			"JPY": 109.98,
			"AUD": 1.1228,
		},
		"AUD": {
			"USD": 0.6564,
			"EUR": 0.5551,
			"GBP": 0.5137,
			"JPY": 97.96,
			"CAD": 0.8906,
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
func (s *FiatExchangeService) calculateFiatExchangeFee(fromCurrency, toCurrency string, amount float64) float64 {
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

// Traitement de l'échange fiat (rapide)
func (s *FiatExchangeService) processFiatExchange(exchange *models.Exchange) {
	// Simuler un court délai de traitement
	time.Sleep(500 * time.Millisecond)

	// Dans un vrai système:
	// 1. Vérifier les soldes des portefeuilles
	// 2. Débiter le portefeuille source
	// 3. Créditer le portefeuille destination
	// 4. Enregistrer les transactions

	// Marquer comme complété
	s.exchangeRepo.UpdateStatus(exchange.ID, "completed")
	
	now := time.Now()
	exchange.CompletedAt = &now
	exchange.Status = "completed"

	// Publier l'événement
	s.publishFiatExchangeEvent("fiat_exchange.completed", exchange)
}

func (s *FiatExchangeService) publishFiatExchangeEvent(eventType string, exchange *models.Exchange) {
	if s.mqChannel == nil {
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

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		"fiat_exchange.events", // exchange
		eventType,              // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}