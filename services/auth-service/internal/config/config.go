package config

	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment string
	Port        string
	DBUrl       string
	RedisURL    string
	RabbitMQURL string
	KafkaBrokers []string
	JWTSecret   string
	
	// JWT settings
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	
	// Email settings
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	
	// SMS settings
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string
	
	// Security settings
	PasswordMinLength int
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
	
	// Rate limiting
	RateLimitRPS int
	
	// Verification settings
	EmailVerificationExpiry time.Duration
	PhoneVerificationExpiry time.Duration
	
	// 2FA settings
	TOTPIssuer string

	// Minio/S3 storage settings
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	MinioPublicURL string
}

func Load() *Config {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	passwordMinLength, _ := strconv.Atoi(getEnv("PASSWORD_MIN_LENGTH", "8"))
	maxLoginAttempts, _ := strconv.Atoi(getEnv("MAX_LOGIN_ATTEMPTS", "5"))
	rateLimitRPS, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPS", "100"))

	accessTokenExpiry, _ := time.ParseDuration(getEnv("ACCESS_TOKEN_EXPIRY", "15m"))
	refreshTokenExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "24h"))
	lockoutDuration, _ := time.ParseDuration(getEnv("LOCKOUT_DURATION", "15m"))
	emailVerificationExpiry, _ := time.ParseDuration(getEnv("EMAIL_VERIFICATION_EXPIRY", "24h"))
	phoneVerificationExpiry, _ := time.ParseDuration(getEnv("PHONE_VERIFICATION_EXPIRY", "5m"))

	// Minio SSL setting
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8081"),
		DBUrl:       getEnv("DB_URL", "postgres://admin:secure_password@localhost:5432/crypto_bank?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:secure_password@localhost:5672/"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		JWTSecret:   getEnv("JWT_SECRET", "ultra_secure_jwt_secret_key_2024"),

		// JWT settings
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,

		// Email settings
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@zekora.com"),

		// SMS settings
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),

		// Security settings
		PasswordMinLength: passwordMinLength,
		MaxLoginAttempts:  maxLoginAttempts,
		LockoutDuration:   lockoutDuration,

		// Rate limiting
		RateLimitRPS: rateLimitRPS,

		// Verification settings
		EmailVerificationExpiry: emailVerificationExpiry,
		PhoneVerificationExpiry: phoneVerificationExpiry,

		// 2FA settings
		TOTPIssuer: getEnv("TOTP_ISSUER", "Zekora"),

		// Minio/S3 storage settings
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioBucket:    getEnv("MINIO_BUCKET", "kyc-documents"),
		MinioUseSSL:    minioUseSSL,
		MinioPublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}