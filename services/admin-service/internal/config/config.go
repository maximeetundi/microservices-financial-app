package config

import (
	"os"
)

type Config struct {
	// Server
	Port        string
	Environment string
	
	// Database (separate admin database)
	AdminDBURL  string
	MainDBURL   string  // Read-only access to main database
	
	// Redis
	RedisURL    string
	
	// RabbitMQ
	RabbitMQURL string
	
	// JWT (separate from user JWT)
	AdminJWTSecret     string
	AdminJWTExpiration string
	
	// Services URLs (for direct API calls when needed)
	AuthServiceURL         string
	WalletServiceURL       string
	TransferServiceURL     string
	CardServiceURL         string
	ExchangeServiceURL     string
	NotificationServiceURL string

	// Minio/S3 storage (for generating presigned URLs)
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	MinioPublicURL string // Public URL for replacing internal Docker address
}

func Load() *Config {
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	return &Config{
		Port:        getEnv("PORT", "8088"),
		Environment: getEnv("ENVIRONMENT", "development"),
		
		AdminDBURL:  getEnv("ADMIN_DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank_admin?sslmode=disable"),
		MainDBURL:   getEnv("MAIN_DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:secure_password@localhost:5672/"),
		
		AdminJWTSecret:     getEnv("ADMIN_JWT_SECRET", "admin_ultra_secure_jwt_secret_2024"),
		AdminJWTExpiration: getEnv("ADMIN_JWT_EXPIRATION", "8h"),
		
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		WalletServiceURL:       getEnv("WALLET_SERVICE_URL", "http://localhost:8083"),
		TransferServiceURL:     getEnv("TRANSFER_SERVICE_URL", "http://localhost:8084"),
		CardServiceURL:         getEnv("CARD_SERVICE_URL", "http://localhost:8086"),
		ExchangeServiceURL:     getEnv("EXCHANGE_SERVICE_URL", "http://localhost:8085"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8087"),

		// Minio configuration
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioBucket:    getEnv("MINIO_BUCKET", "kyc-documents"),
		MinioUseSSL:    minioUseSSL,
		MinioPublicURL: getEnv("MINIO_PUBLIC_URL", "https://minio.maximeetundi.store"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

