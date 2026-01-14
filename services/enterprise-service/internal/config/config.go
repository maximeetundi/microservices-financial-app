package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment string
	Port        string
	MongoDBURI  string
	DBName      string
	JWTSecret   string
	
	// MinIO Config
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	PublicURL      string
	
	// Kafka Config
	KafkaBrokers []string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8097"),
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://admin:secure_password@localhost:27017"),
		DBName:      getEnv("DB_NAME", "enterprise_service"),
		JWTSecret:   getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:    getEnv("MINIO_BUCKET", "finance-app-assets"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		PublicURL:      getEnv("PUBLIC_URL", getEnv("MINIO_PUBLIC_URL", "http://localhost:9000")),
		
		KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

