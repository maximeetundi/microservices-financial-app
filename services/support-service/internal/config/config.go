package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment    string
	DBUrl          string
	RedisURL       string
	KafkaBrokers   []string
	JWTSecret      string
	AdminJWTSecret string
	// MinIO Configuration
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
	MinIOPublicURL string

	AIServiceURL         string
	AIAPIKey             string
	MaxMessagesPerMinute int
	ConversationTimeout  int // hours
}

func Load() *Config {
	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		DBUrl:          getEnv("DB_URL", "postgres://user:password@localhost/crypto_bank_support?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		AdminJWTSecret: getEnv("ADMIN_JWT_SECRET", ""),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "support-attachments"),
		MinIOUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		MinIOPublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),

		AIServiceURL:         getEnv("AI_SERVICE_URL", "https://api.openai.com/v1"),
		AIAPIKey:             getEnv("AI_API_KEY", ""),
		MaxMessagesPerMinute: getEnvInt("MAX_MESSAGES_PER_MINUTE", 30),
		ConversationTimeout:  getEnvInt("CONVERSATION_TIMEOUT_HOURS", 24),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
