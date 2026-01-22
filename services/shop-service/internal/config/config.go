package config

import "os"

type Config struct {
	Port             string
	MongoDBURI       string
	MongoDBName      string
	JWTSecret        string
	KafkaBrokers     string
	KafkaGroupID     string
	WalletServiceURL string
	ExchangeServiceURL string
	MinioEndpoint    string
	MinioAccessKey   string
	MinioSecretKey   string
	MinioBucket      string
	MinioUseSSL      bool
	MinioPublicURL   string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8098"),
		MongoDBURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017/shop_service"),
		MongoDBName:      getEnv("MONGODB_DATABASE", "shop_service"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "shop-service-group"),
		WalletServiceURL: getEnv("WALLET_SERVICE_URL", "http://localhost:8083"),
		ExchangeServiceURL: getEnv("EXCHANGE_SERVICE_URL", "http://localhost:8085"),
		MinioEndpoint:    getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:   getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:   getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioBucket:      getEnv("MINIO_BUCKET", "shop-assets"),
		MinioUseSSL:      getEnv("MINIO_USE_SSL", "false") == "true",
		MinioPublicURL:   getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
