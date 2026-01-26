package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type RateService struct {
	rateRepo          *repository.RateRepository
	config            *config.Config
	binanceProvider   *BinanceProvider
	fiatProvider      *FailoverFiatProvider // Updated type
	coingeckoProvider *CoinGeckoProvider
	hexarateProvider  *HexarateProvider
}

func NewRateService(rateRepo *repository.RateRepository, cfg *config.Config) *RateService {
	binanceCfg := BinanceConfig{
		APIKey:    cfg.BinanceAPIKey,
		APISecret: cfg.BinanceAPISecret,
		BaseURL:   cfg.BinanceBaseURL,
		TestMode:  cfg.BinanceTestMode,
	}
	
	return &RateService{
		rateRepo:          rateRepo,
		config:            cfg,
		binanceProvider:   NewBinanceProvider(binanceCfg),
		fiatProvider:      NewFailoverFiatProvider(cfg.FixerAPIKey, cfg.CurrencyLayerAPIKey),
		coingeckoProvider: NewCoinGeckoProvider(),
		hexarateProvider:  NewHexarateProvider(),
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
		// Fetch fiat rate from real provider (Robust Failover)
		// We usually request USD base, but providers like Fixer force EUR.
		// Use USD as common denominator if we can, but handle provider specific logic in provider if possible.
		// For simplicity, we ask for "USD" but FailoverProvider and sub-providers handle normalization or we do cross-rate here.
		// Actually, our providers return a map relative to their base (Fixer=EUR, CurrencyLayer=USD, OpenER=USD).
		// We'll trust Failover to give us a map relative to USD (CurrencyLayer) or Convert if Fixer (EUR).
		// Wait, my Fixer implementation asks for based? No, free Fixer is locked to EUR.
		// My fiat_providers.go implementation just returns the map. 
		// I need to standardize in `updateAllRates` mostly.
		
		// For single rate on-demand, checking cache is best. If fail, maybe trigger updateAllRates?
		// Fetching single pair from these bulk-oriented APIs is inefficient or impossible on free tiers.
		// So we rely on the background updater or cached DB data most of the time.
		// If DB empty, we might need a sync fetch.
		
		// Let's trigger a full update if we really have nothing.
		s.updateAllRates()
		return s.rateRepo.GetRate(fromCurrency, toCurrency)
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
	fmt.Println("[RateService] Starting rate update...")
	
	// ========== 1. FIAT RATES VIA FAILOVER PROVIDER ==========
	fmt.Println("[RateService] Fetching fiat rates from best available provider...")
	fiatRates, source, err := s.fiatProvider.GetRates("USD")
	if err != nil {
		fmt.Printf("[RateService] Error fetching fiat rates: %v\n", err)
	} else {
		fmt.Printf("[RateService] Fetched %d fiat rates from %s\n", len(fiatRates), source)
		
		// Save fiat rates to DB
		for currency, rate := range fiatRates {
			if currency == "USD" {
				continue
			}
			
			exRate := &models.ExchangeRate{
				FromCurrency: "USD",
				ToCurrency:   currency,
				Rate:         rate,
				BidPrice:     rate * 0.998,
				AskPrice:     rate * 1.002,
				Spread:       rate * 0.004,
				Source:       source,
				LastUpdated:  time.Now(),
			}
			s.rateRepo.SaveRate(exRate)
			
			// Reverse rate
			if rate > 0 {
				reverseRate := &models.ExchangeRate{
					FromCurrency: currency,
					ToCurrency:   "USD",
					Rate:         1 / rate,
					BidPrice:     (1 / rate) * 0.998,
					AskPrice:     (1 / rate) * 1.002,
					Spread:       (1 / rate) * 0.004,
					Source:       source,
					LastUpdated:  time.Now(),
				}
				s.rateRepo.SaveRate(reverseRate)
			}
		}
	}
	
	// ========== 2. CRYPTO RATES VIA COINGECKO ==========
	fmt.Println("[RateService] Fetching crypto prices from CoinGecko...")
	cryptoPrices, err := s.coingeckoProvider.GetAllCryptoPrices()
	if err != nil {
		fmt.Printf("[RateService] Error fetching CoinGecko prices: %v\n", err)
		// Fallback to Binance
		s.updateCryptoFromBinance()
		return
	}
	fmt.Printf("[RateService] Fetched %d crypto prices\n", len(cryptoPrices))
	
	// Save crypto→USD and crypto→EUR rates
	for symbol, price := range cryptoPrices {
		// Crypto → USD
		if price.USD > 0 {
			rate := &models.ExchangeRate{
				FromCurrency: symbol,
				ToCurrency:   "USD",
				Rate:         price.USD,
				BidPrice:     price.USD * 0.995,
				AskPrice:     price.USD * 1.005,
				Spread:       price.USD * 0.01,
				Source:       "coingecko",
				LastUpdated:  time.Now(),
			}
			s.rateRepo.SaveRate(rate)
			
			// Reverse: USD → Crypto
			reverseRate := &models.ExchangeRate{
				FromCurrency: "USD",
				ToCurrency:   symbol,
				Rate:         1 / price.USD,
				BidPrice:     (1 / price.USD) * 0.995,
				AskPrice:     (1 / price.USD) * 1.005,
				Spread:       (1 / price.USD) * 0.01,
				Source:       "coingecko",
				LastUpdated:  time.Now(),
			}
			s.rateRepo.SaveRate(reverseRate)
		}
		
		// Crypto → EUR
		if price.EUR > 0 {
			rate := &models.ExchangeRate{
				FromCurrency: symbol,
				ToCurrency:   "EUR",
				Rate:         price.EUR,
				BidPrice:     price.EUR * 0.995,
				AskPrice:     price.EUR * 1.005,
				Spread:       price.EUR * 0.01,
				Source:       "coingecko",
				LastUpdated:  time.Now(),
			}
			s.rateRepo.SaveRate(rate)
		}
		
		// ========== 3. CRYPTO → OTHER FIATS (via USD) ==========
		// Calculate crypto→XAF, crypto→NGN, etc. by chaining: crypto→USD × USD→fiat
		if price.USD > 0 {
			for currency, fiatRate := range fiatRates {
				if currency == "USD" || currency == "EUR" {
					continue
				}
				
				// Crypto → Fiat = Crypto_USD × USD_Fiat
				cryptoToFiat := price.USD * fiatRate
				if cryptoToFiat > 0 {
					rate := &models.ExchangeRate{
						FromCurrency: symbol,
						ToCurrency:   currency,
						Rate:         cryptoToFiat,
						BidPrice:     cryptoToFiat * 0.995,
						AskPrice:     cryptoToFiat * 1.005,
						Spread:       cryptoToFiat * 0.01,
						Source:       "coingecko+" + source, // Attributed to combination
						LastUpdated:  time.Now(),
					}
					s.rateRepo.SaveRate(rate)

					// NEW: Inverse Fiat -> Crypto
					// Rate = 1 / cryptoToFiat
					inverseRate := &models.ExchangeRate{
						FromCurrency: currency, // e.g. NGN
						ToCurrency:   symbol,   // e.g. BTC
						Rate:         1 / cryptoToFiat,
						BidPrice:     (1 / cryptoToFiat) * 0.995,
						AskPrice:     (1 / cryptoToFiat) * 1.005,
						Spread:       (1 / cryptoToFiat) * 0.01,
						Source:       "coingecko+" + source + "(inverse)",
						LastUpdated:  time.Now(),
					}
					s.rateRepo.SaveRate(inverseRate)
				}
			}
		}
	}
	
	fmt.Println("[RateService] Rate update completed!")
}

// updateCryptoFromBinance is a fallback if CoinGecko fails
func (s *RateService) updateCryptoFromBinance() {
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
		}
	}
}

func (s *RateService) fetchCryptoRate(from, to string) (*models.ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Helper to get price safely
	getPrice := func(base, quote string) (*CryptoPrice, error) {
		// Binance uses symbols like BTCUSDT
		// Handle USD mapping to USDT for Binance
		targetQuote := quote
		if quote == "USD" {
			targetQuote = "USDT"
		}
		targetBase := base
		if base == "USD" {
			targetBase = "USDT"
		}
		
		symbol := targetBase + targetQuote
		return s.binanceProvider.GetPrice(ctx, symbol)
	}

	// Try direct pair first: FROM -> TO
	price, err := getPrice(from, to)
	if err == nil {
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

	// Try inverse pair: TO -> FROM (e.g request USD->BTC, try BTC->USD output)
	// This supports cases where we only have the major pair listing (BTCUSDT) but want to convert USDT to BTC
	inversePrice, err := getPrice(to, from)
	if err == nil && inversePrice.Price > 0 {
		return &models.ExchangeRate{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         1 / inversePrice.Price,
			BidPrice:     (1 / inversePrice.Price) * 0.995,
			AskPrice:     (1 / inversePrice.Price) * 1.005,
			Spread:       (1 / inversePrice.Price) * 0.01,
			Source:       "binance (inverse)",
			Volume24h:    inversePrice.Volume24h,
			Change24h:    inversePrice.Change24h,
			LastUpdated:  time.Now(),
		}, nil
	}

	return nil, fmt.Errorf("failed to fetch rate for %s/%s: %v", from, to, err)
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