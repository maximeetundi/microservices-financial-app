package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string
	DBUrl       string
	RabbitMQURL string
	
	// Email configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	
	// SMS configuration
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string
	
	// Push notification configuration
	FCMServerKey string
	APNSKeyID    string
	APNSTeamID   string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8087"),
		DBUrl:       getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:secure_password@localhost:5672/"),
		
		// Email
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@cryptobank.app"),
		
		// SMS
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
		
		// Push
		FCMServerKey: getEnv("FCM_SERVER_KEY", ""),
		APNSKeyID:    getEnv("APNS_KEY_ID", ""),
		APNSTeamID:   getEnv("APNS_TEAM_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
