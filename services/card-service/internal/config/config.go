package config

import (
	"os"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/models"
)

type Config struct {
	Environment      string
	Port             string
	DBUrl            string
	KafkaBrokers     string
	KafkaGroupID     string
	JWTSecret        string
	WalletServiceURL string

	// Card issuer settings
	CardIssuer        CardIssuerConfig
	DefaultCardIssuer string

	// Card limits
	MaxCardsPerUser     int
	DefaultDailyLimit   float64
	DefaultMonthlyLimit float64
	MinLoadAmount       float64
	MaxLoadAmount       float64
	MaxAutoReloadAmount float64

	// Fees
	CardFees      map[string]float64
	DefaultLimits map[string]*models.CardLimits

	// Security
	RateLimitRPS int
}

type CardIssuerConfig struct {
	Provider      string
	APIKey        string
	APISecret     string
	BaseURL       string
	WebhookSecret string
}

func Load() *Config {
	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))
	maxCardsPerUser, _ := strconv.Atoi(getEnv("MAX_CARDS_PER_USER", "10"))
	dailyLimit, _ := strconv.ParseFloat(getEnv("DEFAULT_DAILY_LIMIT", "1000"), 64)
	monthlyLimit, _ := strconv.ParseFloat(getEnv("DEFAULT_MONTHLY_LIMIT", "10000"), 64)

	return &Config{
		Environment:       getEnv("ENVIRONMENT", "development"),
		Port:              getEnv("PORT", "8086"),
		DBUrl:             getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaGroupID:      getEnv("KAFKA_GROUP_ID", "card-service-group"),
		JWTSecret:         getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		WalletServiceURL:  getEnv("WALLET_SERVICE_URL", "http://wallet-service:8083"),
		DefaultCardIssuer: getEnv("DEFAULT_CARD_ISSUER", "marqeta"),

		CardIssuer: CardIssuerConfig{
			Provider:      getEnv("CARD_ISSUER_PROVIDER", "marqeta"),
			APIKey:        getEnv("CARD_ISSUER_API_KEY", ""),
			APISecret:     getEnv("CARD_ISSUER_API_SECRET", ""),
			BaseURL:       getEnv("CARD_ISSUER_BASE_URL", "https://sandbox-api.marqeta.com/v3"),
			WebhookSecret: getEnv("CARD_ISSUER_WEBHOOK_SECRET", ""),
		},

		MaxCardsPerUser:     maxCardsPerUser,
		DefaultDailyLimit:   dailyLimit,
		DefaultMonthlyLimit: monthlyLimit,
		MinLoadAmount:       10.0,
		MaxLoadAmount:       10000.0,
		MaxAutoReloadAmount: 5000.0,
		RateLimitRPS:        rateLimitRPS,

		CardFees: map[string]float64{
			"load_fee_percentage": 0.5,
			"load_fee_minimum":    0.50,
			"shipping":            5.00,
			"express_shipping":    15.00,
			"replacement":         10.00,
		},

		DefaultLimits: map[string]*models.CardLimits{
			"prepaid": {
				DailyLimit:    1000,
				MonthlyLimit:  5000,
				SingleTxLimit: 500,
				ATMDailyLimit: 300,
			},
			"virtual": {
				DailyLimit:    2000,
				MonthlyLimit:  10000,
				SingleTxLimit: 1000,
				ATMDailyLimit: 0, // No ATM for virtual
			},
			"gift": {
				DailyLimit:    500,
				MonthlyLimit:  500,
				SingleTxLimit: 500,
				ATMDailyLimit: 0,
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
