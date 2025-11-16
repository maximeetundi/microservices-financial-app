package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment    string
	Port          string
	JWTSecret     string
	DBUrl         string
	RedisURL      string
	RabbitMQURL   string
	
	// Service URLs
	AuthServiceURL     string
	UserServiceURL     string
	WalletServiceURL   string
	TransferServiceURL string
	ExchangeServiceURL string
	CardServiceURL     string
	NotificationServiceURL string
	
	// Security
	RateLimitRPS      int
	MaxRequestSize    int64
	
	// Crypto settings
	EncryptionKey     string
	BlockchainRPC     map[string]string
}

func Load() *Config {
	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))
	maxRequestSize, _ := strconv.ParseInt(getEnv("MAX_REQUEST_SIZE", "10485760"), 10, 64) // 10MB

	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		DBUrl:         getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://admin:secure_password@localhost:5672/"),
		
		// Service URLs
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
		UserServiceURL:     getEnv("USER_SERVICE_URL", "http://user-service:8082"),
		WalletServiceURL:   getEnv("WALLET_SERVICE_URL", "http://wallet-service:8083"),
		TransferServiceURL: getEnv("TRANSFER_SERVICE_URL", "http://transfer-service:8084"),
		ExchangeServiceURL: getEnv("EXCHANGE_SERVICE_URL", "http://exchange-service:8085"),
		CardServiceURL:     getEnv("CARD_SERVICE_URL", "http://card-service:8086"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8087"),
		
		// Security
		RateLimitRPS:   rateLimitRPS,
		MaxRequestSize: maxRequestSize,
		
		// Crypto
		EncryptionKey: getEnv("ENCRYPTION_KEY", "32-byte-encryption-key-for-crypto"),
		BlockchainRPC: map[string]string{
			"BTC": getEnv("BTC_RPC", "https://bitcoin-rpc.example.com"),
			"ETH": getEnv("ETH_RPC", "https://mainnet.infura.io/v3/YOUR-PROJECT-ID"),
			"BSC": getEnv("BSC_RPC", "https://bsc-dataseed.binance.org/"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}