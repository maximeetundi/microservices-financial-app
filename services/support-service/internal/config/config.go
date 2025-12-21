package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment      string
	DBUrl            string
	RedisURL         string
	RabbitMQURL      string
	JWTSecret        string
	AIServiceURL     string
	AIAPIKey         string
	MaxMessagesPerMinute int
	ConversationTimeout  int // hours
}

func Load() *Config {
	return &Config{
		Environment:          getEnv("ENVIRONMENT", "development"),
		DBUrl:                getEnv("DB_URL", "postgres://user:password@localhost/crypto_bank_support?sslmode=disable"),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL:          getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key"),
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
