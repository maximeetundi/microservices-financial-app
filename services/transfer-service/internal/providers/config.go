package providers

import (
	"os"
	"strconv"
)

// Config holds all provider configurations
type Config struct {
	Flutterwave FlutterwaveConfig
	Thunes      ThunesConfig
	Stripe      StripeConfig
	CryptoRails CryptoRailsConfig
}

// LoadConfig loads provider configuration from environment variables
func LoadConfig() *Config {
	poolThreshold, _ := strconv.ParseFloat(getEnv("CRYPTO_POOL_THRESHOLD", "500"), 64)
	poolEnabled, _ := strconv.ParseBool(getEnv("CRYPTO_POOL_ENABLED", "true"))
	
	return &Config{
		Flutterwave: FlutterwaveConfig{
			// Flutterwave API Keys - Get from https://dashboard.flutterwave.com/settings/api
			SecretKey:   getEnv("FLUTTERWAVE_SECRET_KEY", ""),
			PublicKey:   getEnv("FLUTTERWAVE_PUBLIC_KEY", ""),
			EncryptKey:  getEnv("FLUTTERWAVE_ENCRYPT_KEY", ""),
			BaseURL:     getEnv("FLUTTERWAVE_BASE_URL", "https://api.flutterwave.com/v3"),
			CallbackURL: getEnv("FLUTTERWAVE_CALLBACK_URL", ""),
		},
		Thunes: ThunesConfig{
			// Thunes API Keys - Get from https://portal.thunes.com
			APIKey:      getEnv("THUNES_API_KEY", ""),
			APISecret:   getEnv("THUNES_API_SECRET", ""),
			BaseURL:     getEnv("THUNES_BASE_URL", "https://api.thunes.com/v2"),
			CallbackURL: getEnv("THUNES_CALLBACK_URL", ""),
		},
		Stripe: StripeConfig{
			// Stripe API Keys - Get from https://dashboard.stripe.com/apikeys
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			BaseURL:        getEnv("STRIPE_BASE_URL", "https://api.stripe.com/v1"),
		},
		CryptoRails: CryptoRailsConfig{
			// Circle USDC API - Get from https://developers.circle.com
			CircleAPIKey:  getEnv("CIRCLE_API_KEY", ""),
			CircleBaseURL: getEnv("CIRCLE_BASE_URL", "https://api.circle.com/v1"),
			
			// Binance API - Get from https://www.binance.com/en/my/settings/api-management
			BinanceAPIKey:    getEnv("BINANCE_API_KEY", ""),
			BinanceAPISecret: getEnv("BINANCE_API_SECRET", ""),
			BinanceBaseURL:   getEnv("BINANCE_BASE_URL", "https://api.binance.com"),
			
			// Internal pool settings
			UseInternalPoolThreshold: poolThreshold,
			InternalPoolEnabled:      poolEnabled,
			PreferredStablecoin:      getEnv("PREFERRED_STABLECOIN", "USDC"),
			
			// Wallet addresses for receiving crypto
			EthereumWallet: getEnv("ETH_WALLET_ADDRESS", ""),
			TronWallet:     getEnv("TRX_WALLET_ADDRESS", ""),
			PolygonWallet:  getEnv("MATIC_WALLET_ADDRESS", ""),
		},
	}
}

// InitializeRouter creates and configures the zone router with all providers
func InitializeRouter(cfg *Config) *ZoneRouter {
	router := NewZoneRouter()
	
	// Register Flutterwave (Africa)
	if cfg.Flutterwave.SecretKey != "" {
		flutterwave := NewFlutterwaveProvider(cfg.Flutterwave)
		router.RegisterProvider(flutterwave)
	}
	
	// Register Thunes (Global)
	if cfg.Thunes.APIKey != "" {
		thunes := NewThunesProvider(cfg.Thunes)
		router.RegisterProvider(thunes)
	}
	
	// Register Stripe (Europe, North America)
	if cfg.Stripe.SecretKey != "" {
		stripe := NewStripeProvider(cfg.Stripe)
		router.RegisterProvider(stripe)
	}
	
	return router
}

// InitializeCryptoRails creates the crypto rails provider
func InitializeCryptoRails(cfg *Config) *CryptoRailsProvider {
	return NewCryptoRailsProvider(cfg.CryptoRails)
}

// InitializeOrchestrator creates the full transfer orchestrator
func InitializeOrchestrator(cfg *Config) *TransferOrchestrator {
	cryptoRails := InitializeCryptoRails(cfg)
	zoneRouter := InitializeRouter(cfg)
	return NewTransferOrchestrator(cryptoRails, zoneRouter)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
