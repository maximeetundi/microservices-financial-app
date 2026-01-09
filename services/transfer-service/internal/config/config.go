package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment string
	Port        string
	DBUrl       string
	RabbitMQURL string
	WalletServiceURL string
	JWTSecret   string
	
	// Transfer limits
	DailyTransferLimits    map[string]float64
	MonthlyTransferLimits  map[string]float64
	SingleTransferLimits   map[string]float64
	
	// Transfer fees
	TransferFees map[string]float64
	
	// Mobile money providers
	MobileMoneyProviders map[string]MobileMoneyConfig
	
	// International transfer
	InternationalProviders map[string]InternationalConfig
	
	// Compliance settings
	ComplianceSettings ComplianceConfig
	
	// Security
	RateLimitRPS int
}

type MobileMoneyConfig struct {
	APIKey      string
	APISecret   string
	BaseURL     string
	CallbackURL string
	Countries   []string
}

type InternationalConfig struct {
	APIKey    string
	APISecret string
	BaseURL   string
	Countries []string
	Currencies []string
}

type ComplianceConfig struct {
	KYCRequired           bool
	AMLCheckRequired      bool
	MaxAmountWithoutKYC   float64
	SanctionsCheckEnabled bool
	PEPCheckEnabled       bool
}

func Load() *Config {
	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8084"),
		DBUrl:       getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:secure_password@localhost:5672/"),
		WalletServiceURL: getEnv("WALLET_SERVICE_URL", "http://wallet-service:8081"),
		JWTSecret:   getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		
		// Daily limits by KYC level
		DailyTransferLimits: map[string]float64{
			"level_1": 1000.0,
			"level_2": 10000.0,
			"level_3": 100000.0,
		},
		
		// Monthly limits by KYC level
		MonthlyTransferLimits: map[string]float64{
			"level_1": 5000.0,
			"level_2": 50000.0,
			"level_3": 500000.0,
		},
		
		// Single transaction limits
		SingleTransferLimits: map[string]float64{
			"USD": 10000.0,
			"EUR": 8500.0,
			"BTC": 1.0,
			"ETH": 10.0,
		},
		
		// Transfer fees (percentage)
		TransferFees: map[string]float64{
			"domestic":      0.5,  // 0.5%
			"international": 2.0,  // 2.0%
			"mobile_money":  1.0,  // 1.0%
			"crypto":        0.1,  // 0.1%
		},
		
		// Mobile money providers
		MobileMoneyProviders: map[string]MobileMoneyConfig{
			"mtn": {
				APIKey:      getEnv("MTN_API_KEY", ""),
				APISecret:   getEnv("MTN_API_SECRET", ""),
				BaseURL:     getEnv("MTN_BASE_URL", "https://api.mtn.com"),
				CallbackURL: getEnv("MTN_CALLBACK_URL", ""),
				Countries:   []string{"UG", "GH", "CI", "CM", "BJ", "GN", "RW", "ZM", "SS"},
			},
			"airtel": {
				APIKey:      getEnv("AIRTEL_API_KEY", ""),
				APISecret:   getEnv("AIRTEL_API_SECRET", ""),
				BaseURL:     getEnv("AIRTEL_BASE_URL", "https://api.airtel.com"),
				CallbackURL: getEnv("AIRTEL_CALLBACK_URL", ""),
				Countries:   []string{"UG", "KE", "TZ", "RW", "ZM", "MW", "MG", "TD", "NE", "CM"},
			},
			"mpesa": {
				APIKey:      getEnv("MPESA_API_KEY", ""),
				APISecret:   getEnv("MPESA_API_SECRET", ""),
				BaseURL:     getEnv("MPESA_BASE_URL", "https://api.safaricom.co.ke"),
				CallbackURL: getEnv("MPESA_CALLBACK_URL", ""),
				Countries:   []string{"KE"},
			},
		},
		
		// International transfer providers
		InternationalProviders: map[string]InternationalConfig{
			"wise": {
				APIKey:     getEnv("WISE_API_KEY", ""),
				APISecret:  getEnv("WISE_API_SECRET", ""),
				BaseURL:    getEnv("WISE_BASE_URL", "https://api.transferwise.com"),
				Countries:  []string{"US", "GB", "DE", "FR", "AU", "CA", "SG", "JP"},
				Currencies: []string{"USD", "EUR", "GBP", "AUD", "CAD", "SGD", "JPY"},
			},
			"remitly": {
				APIKey:     getEnv("REMITLY_API_KEY", ""),
				APISecret:  getEnv("REMITLY_API_SECRET", ""),
				BaseURL:    getEnv("REMITLY_BASE_URL", "https://api.remitly.com"),
				Countries:  []string{"US", "GB", "CA", "AU"},
				Currencies: []string{"USD", "EUR", "GBP", "CAD", "AUD"},
			},
		},
		
		// Compliance settings
		ComplianceSettings: ComplianceConfig{
			KYCRequired:           true,
			AMLCheckRequired:      true,
			MaxAmountWithoutKYC:   500.0, // $500 max without KYC
			SanctionsCheckEnabled: true,
			PEPCheckEnabled:       true,
		},
		
		// Security
		RateLimitRPS: rateLimitRPS,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}