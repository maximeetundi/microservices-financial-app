package config

import (
	"os"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/common/secrets"
)

type Config struct {
	Environment  string
	Port         string
	DBUrl        string
	RedisURL     string
	KafkaBrokers string
	KafkaGroupID string
	JWTSecret    string
	BaseURL      string

	// Tatum settings
	TatumAPIKey  string
	TatumBaseURL string

	// Transaction settings
	MaxTransactionAmount map[string]float64
	MinConfirmations     map[string]int

	// Fee settings
	NetworkFees map[string]float64

	// Crypto settings
	EncryptionKey string
	BlockchainRPC map[string]string
	CryptoAPIKeys map[string]string

	// Security
	RateLimitRPS int
}

func Load() *Config {
	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))

	return &Config{
		Environment:  getEnv("ENVIRONMENT", "development"),
		Port:         getEnv("PORT", "8083"),
		DBUrl:        getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "wallet-service-group"),
		JWTSecret:    getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),
		BaseURL:      getEnv("BASE_URL", "https://api.app.tech-afm.com/wallet-service"),

		// Tatum settings
		TatumAPIKey:  getEnv("TATUM_API_KEY", ""),
		TatumBaseURL: getEnv("TATUM_BASE_URL", "https://api.tatum.io/v3"), // Default to v3 API

		// Crypto settings
		EncryptionKey: getEnv("ENCRYPTION_KEY", "32-byte-encryption-key-for-crypto"),
		BlockchainRPC: map[string]string{
			"BTC": getEnv("BTC_RPC", "https://bitcoin-rpc.example.com"),
			"ETH": getEnv("ETH_RPC", "https://mainnet.infura.io/v3/YOUR-PROJECT-ID"),
			"BSC": getEnv("BSC_RPC", "https://bsc-dataseed.binance.org/"),
		},
		CryptoAPIKeys: map[string]string{
			"INFURA":  getEnv("INFURA_API_KEY", ""),
			"ALCHEMY": getEnv("ALCHEMY_API_KEY", ""),
			"MORALIS": getEnv("MORALIS_API_KEY", ""),
		},

		// Security
		RateLimitRPS: rateLimitRPS,

		// Transaction limits
		MaxTransactionAmount: map[string]float64{
			"BTC": 10.0,
			"ETH": 100.0,
			"USD": 50000.0,
			"EUR": 45000.0,
		},
		MinConfirmations: map[string]int{
			"BTC": 3,
			"ETH": 12,
			"BSC": 15,
		},

		// Network fees (in respective currencies)
		NetworkFees: map[string]float64{
			"BTC": 0.0001,
			"ETH": 0.002,
			"BSC": 0.0005,
		},
	}

	// Vault Integration
	if vaultClient, err := secrets.NewVaultClient(); err == nil {
		if secretData, err := vaultClient.GetSecret("secret/wallet-service"); err == nil {
			if val, ok := secretData["tatum_api_key"].(string); ok && val != "" {
				cfg.TatumAPIKey = val
			}
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
