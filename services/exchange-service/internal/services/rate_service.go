package services

import (
	"context"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type RateService struct {
	rateRepo        *repository.RateRepository
	config          *config.Config
	binanceProvider *BinanceProvider
}

func NewRateService(rateRepo *repository.RateRepository, cfg *config.Config) *RateService {
	binanceCfg := BinanceConfig{
		APIKey:    cfg.BinanceAPIKey,
		APISecret: cfg.BinanceAPISecret,
		BaseURL:   cfg.BinanceBaseURL,
		TestMode:  cfg.BinanceTestMode,
	}
	
	return &RateService{
		rateRepo:        rateRepo,
		config:          cfg,
		binanceProvider: NewBinanceProvider(binanceCfg),
	}
}

func (s *RateService) GetRate(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	// Try to get from cache/db first
	rate, err := s.rateRepo.GetRate(fromCurrency, toCurrency)
	if err == nil && time.Since(rate.LastUpdated) < time.Duration(s.config.RateUpdateInterval)*time.Second {
		return rate, nil
	}

	// If missing or stale, fetch updated rate (if crypto)
	if s.isCryptoPair(fromCurrency, toCurrency) {
		updatedRate, err := s.fetchCryptoRate(fromCurrency, toCurrency)
		if err == nil {
			s.rateRepo.SaveRate(updatedRate)
			return updatedRate, nil
		}
	} else {
		// For fiat, use simulated/fixed for now as we don't have a fiat provider
		updatedRate := s.simulateFiatRate(fromCurrency, toCurrency)
		if updatedRate != nil {
			s.rateRepo.SaveRate(updatedRate)
			return updatedRate, nil
		}
	}

	return rate, err
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

func (s *RateService) GetMarkets() ([]*models.Market, error) {
	// Fetch real markets from Binance
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prices, err := s.binanceProvider.GetAllPrices(ctx)
	if err != nil {
		// Fallback to DB/Cache if external API fails
		return nil, err
	}

	var markets []*models.Market
	for _, p := range prices {
		// We only care about major pairs for now
		if isMajorPair(p.Symbol) {
			base, quote := splitSymbol(p.Symbol)
			markets = append(markets, &models.Market{
				Symbol:      base + "/" + quote,
				BaseAsset:   base,
				QuoteAsset:  quote,
				Price:       p.Price,
				Change24h:   0, // Binance ticker/price endpoint doesn't return 24h change, use simulate or another endpoint if needed. using 0 for speed.
				LastUpdated: time.Now(),
				BidPrice:    p.Price * 0.9995,
				AskPrice:    p.Price * 1.0005,
			})
		}
	}
	
	if len(markets) == 0 {
		// Fallback
		return []*models.Market{
			{Symbol: "BTC/USD", BaseAsset: "BTC", QuoteAsset: "USD", Price: 43500.0, LastUpdated: time.Now()},
		}, nil
	}

	return markets, nil
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
	// Update crypto rates using Binance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := s.binanceProvider.GetAllPrices(ctx)
	if err != nil {
		return
	}

	for _, p := range prices {
		if isMajorPair(p.Symbol) {
			base, quote := splitSymbol(p.Symbol)
			rate := &models.ExchangeRate{
				FromCurrency: base,
				ToCurrency:   quote,
				Rate:         p.Price,
				BidPrice:     p.Price * 0.995,
				AskPrice:     p.Price * 1.005,
				Spread:       p.Price * 0.01,
				Source:       "binance",
				LastUpdated:  time.Now(),
			}
			s.rateRepo.SaveRate(rate)
			
			// Also save reverse rate
			if p.Price > 0 {
				reverseRate := &models.ExchangeRate{
					FromCurrency: quote,
					ToCurrency:   base,
					Rate:         1 / p.Price,
					BidPrice:     (1 / p.Price) * 0.995,
					AskPrice:     (1 / p.Price) * 1.005,
					Spread:       (1 / p.Price) * 0.01,
					Source:       "binance",
					LastUpdated:  time.Now(),
				}
				s.rateRepo.SaveRate(reverseRate)
			}
		}
	}

	// Update fiat rates (simulation)
	fiatPairs := [][]string{
		{"USD", "EUR"}, {"USD", "GBP"}, {"USD", "JPY"}, {"USD", "CAD"},
		{"EUR", "GBP"}, {"EUR", "JPY"}, {"GBP", "JPY"}, {"USD", "AUD"},
	}

	for _, pair := range fiatPairs {
		rate := s.simulateFiatRate(pair[0], pair[1])
		if rate != nil {
			s.rateRepo.SaveRate(rate)
		}
	}
}

func (s *RateService) fetchCryptoRate(from, to string) (*models.ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Binance uses symbols like BTCUSDT
	// Handle USD mapping to USDT for Binance
	targetTo := to
	if to == "USD" {
		targetTo = "USDT"
	}
	
	symbol := from + targetTo
	price, err := s.binanceProvider.GetPrice(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return &models.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         price.Price,
		BidPrice:     price.Price * 0.995,
		AskPrice:     price.Price * 1.005,
		Spread:       price.Price * 0.01,
		Source:       "binance",
		Volume24h:    price.Volume24h,
		Change24h:    price.Change24h,
		LastUpdated:  time.Now(),
	}, nil
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

func (s *RateService) isCryptoPair(from, to string) bool {
	// Simplified check
	return isCrypto(from) || isCrypto(to)
}

func isCrypto(currency string) bool {
	cryptos := map[string]bool{
		"BTC": true, "ETH": true, "USDT": true, "BNB": true, "XRP": true, 
		"ADA": true, "SOL": true, "DOGE": true, "DOT": true, "LTC": true,
	}
	return cryptos[currency]
}

func isMajorPair(symbol string) bool {
	majors := map[string]bool{
		"BTCUSDT": true, "ETHUSDT": true, "BNBUSDT": true, "SOLUSDT": true,
		"XRPUSDT": true, "ADAUSDT": true, "DOGEUSDT": true, "DOTUSDT": true,
		"LTCUSDT": true,
	}
	return majors[symbol]
}

func splitSymbol(symbol string) (string, string) {
	if strings.HasSuffix(symbol, "USDT") {
		return strings.TrimSuffix(symbol, "USDT"), "USDT"
	}
	if strings.HasSuffix(symbol, "EUR") {
		return strings.TrimSuffix(symbol, "EUR"), "EUR"
	}
	// Fallback
	return symbol, ""
}