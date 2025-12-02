package services

import (
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type RateService struct {
	rateRepo *repository.RateRepository
	config   *config.Config
}

func NewRateService(rateRepo *repository.RateRepository, cfg *config.Config) *RateService {
	return &RateService{
		rateRepo: rateRepo,
		config:   cfg,
	}
}

func (s *RateService) GetRate(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	return s.rateRepo.GetRate(fromCurrency, toCurrency)
}

func (s *RateService) UpdateRate(rate *models.ExchangeRate) error {
	return s.rateRepo.SaveRate(rate)
}

func (s *RateService) GetAllRates() ([]*models.ExchangeRate, error) {
	return s.rateRepo.GetAllRates()
}

func (s *RateService) InvalidateCache(fromCurrency, toCurrency string) {
	s.rateRepo.InvalidateCache(fromCurrency, toCurrency)
}

// StartRateUpdater starts a background goroutine that periodically updates exchange rates
func (s *RateService) StartRateUpdater() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.config.RateUpdateInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			s.updateAllRates()
		}
	}()
}

func (s *RateService) updateAllRates() {
	// Update crypto rates (simulation - in real app, call external APIs)
	cryptoPairs := [][]string{
		{"BTC", "USD"}, {"ETH", "USD"}, {"BTC", "EUR"}, {"ETH", "EUR"},
		{"LTC", "USD"}, {"ADA", "USD"}, {"DOT", "USD"}, {"XRP", "USD"},
	}

	for _, pair := range cryptoPairs {
		rate := s.simulateCryptoRate(pair[0], pair[1])
		s.rateRepo.SaveRate(rate)
	}

	// Update fiat rates (simulation)
	fiatPairs := [][]string{
		{"USD", "EUR"}, {"USD", "GBP"}, {"USD", "JPY"}, {"USD", "CAD"},
		{"EUR", "GBP"}, {"EUR", "JPY"}, {"GBP", "JPY"}, {"USD", "AUD"},
	}

	for _, pair := range fiatPairs {
		rate := s.simulateFiatRate(pair[0], pair[1])
		s.rateRepo.SaveRate(rate)
	}
}

func (s *RateService) simulateCryptoRate(from, to string) *models.ExchangeRate {
	// Simulate realistic crypto rates (in real app, fetch from APIs like CoinGecko, Binance, etc.)
	baseRates := map[string]float64{
		"BTC": 43500.0,
		"ETH": 2450.0,
		"LTC": 72.0,
		"ADA": 0.52,
		"DOT": 7.8,
		"XRP": 0.63,
	}

	fiatRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.78,
	}

	cryptoPrice := baseRates[from]
	fiatMultiplier := fiatRates[to]

	if cryptoPrice == 0 || fiatMultiplier == 0 {
		return nil
	}

	// Add some volatility (+/- 2%)
	variation := 1.0 + (0.04*float64(time.Now().UnixNano()%1000)/1000.0 - 0.02)
	rate := cryptoPrice * fiatMultiplier * variation

	// Calculate spread (0.5% for crypto)
	spread := rate * 0.005
	bidPrice := rate - spread/2
	askPrice := rate + spread/2

	return &models.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		BidPrice:     bidPrice,
		AskPrice:     askPrice,
		Spread:       spread,
		Source:       "crypto_exchange_api",
		Volume24h:    1000000.0 + float64(time.Now().UnixNano()%500000),
		Change24h:    -0.1 + (0.2 * float64(time.Now().UnixNano()%1000)/1000), // +/- 10%
		LastUpdated:  time.Now(),
	}
}

func (s *RateService) simulateFiatRate(from, to string) *models.ExchangeRate {
	// Simulate realistic fiat rates
	rates := map[string]map[string]float64{
		"USD": {"EUR": 0.8456, "GBP": 0.7821, "JPY": 149.23, "CAD": 1.3567, "AUD": 1.5234},
		"EUR": {"USD": 1.1826, "GBP": 0.9248, "JPY": 176.45},
		"GBP": {"USD": 1.2787, "EUR": 1.0814, "JPY": 190.76},
		"JPY": {"USD": 0.0067, "EUR": 0.0057, "GBP": 0.0052},
	}

	baseRate, exists := rates[from][to]
	if !exists {
		return nil
	}

	// Add minimal volatility for fiat (+/- 0.1%)
	variation := 1.0 + (0.002*float64(time.Now().UnixNano()%1000)/1000.0 - 0.001)
	rate := baseRate * variation

	// Calculate spread (0.1% for fiat)
	spread := rate * 0.001
	bidPrice := rate - spread/2
	askPrice := rate + spread/2

	return &models.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		BidPrice:     bidPrice,
		AskPrice:     askPrice,
		Spread:       spread,
		Source:       "fiat_exchange_api",
		Volume24h:    5000000.0,
		Change24h:    -0.02 + (0.04 * float64(time.Now().UnixNano()%1000)/1000), // +/- 2%
		LastUpdated:  time.Now(),
	}
}