package providers

import "os"

// Config holds all provider configurations
type Config struct {
	Flutterwave FlutterwaveConfig
	Thunes      ThunesConfig
	Stripe      StripeConfig
}

// LoadConfig loads provider configuration from environment variables
func LoadConfig() *Config {
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
