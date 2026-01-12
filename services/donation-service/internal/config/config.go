package config

import (
	"os"
)

type Config struct {
	Port          string
	MongoDBURI    string
	MongoDBName   string
	KafkaBrokers  string
	KafkaGroupID  string
	JWTSecret     string
	WalletServiceURL string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8096"),
		MongoDBURI:    getEnv("MONGODB_URI", "mongodb://admin:secure_password@mongodb:27017/donation_service?authSource=admin"),
		MongoDBName:   getEnv("MONGODB_DATABASE", "donation_service"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "donation-service-group"),
		JWTSecret:     getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		WalletServiceURL: getEnv("WALLET_SERVICE_URL", "http://wallet-service:8083"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
