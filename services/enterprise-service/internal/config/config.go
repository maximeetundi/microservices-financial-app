package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string
	MongoDBURI  string
	DBName      string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8097"),
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://admin:secure_password@localhost:27017"),
		DBName:      getEnv("DB_NAME", "enterprise_service"),
		JWTSecret:   getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
