package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment        string
	DBUrl              string
	RedisURL           string
	RabbitMQURL        string
	JWTSecret          string
	WalletServiceURL   string
	ExchangeFees       map[string]float64
	RateUpdateInterval int
	
	// Binance configuration
	BinanceAPIKey    string
	BinanceAPISecret string
	BinanceBaseURL   string
	BinanceTestMode  bool
}

func Load() *Config {
	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		DBUrl:            getEnv("DATABASE_URL", "postgres://user:password@localhost/crypto_bank_exchange?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		WalletServiceURL: getEnv("WALLET_SERVICE_URL", "http://localhost:8084"),
		ExchangeFees: map[string]float64{
			"crypto_to_crypto": getEnvFloat("CRYPTO_TO_CRYPTO_FEE", 0.5),  // 0.5%
			"crypto_to_fiat":   getEnvFloat("CRYPTO_TO_FIAT_FEE", 0.75),   // 0.75%
			"fiat_to_crypto":   getEnvFloat("FIAT_TO_CRYPTO_FEE", 0.75),   // 0.75%
			"fiat_to_fiat":     getEnvFloat("FIAT_TO_FIAT_FEE", 0.25),     // 0.25%
		},
		RateUpdateInterval: getEnvInt("RATE_UPDATE_INTERVAL", 30), // 30 seconds
		
		// Binance API - Get from https://www.binance.com/en/my/settings/api-management
		BinanceAPIKey:    getEnv("BINANCE_API_KEY", ""),
		BinanceAPISecret: getEnv("BINANCE_API_SECRET", ""),
		BinanceBaseURL:   getEnv("BINANCE_BASE_URL", "https://api.binance.com"),
		BinanceTestMode:  getEnvBool("BINANCE_TEST_MODE", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}