package providers

import (
	"log"
	"os"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/common/secrets"
)

// Config holds all provider configurations
// Priority order: 1. Vault, 2. Environment variables
type Config struct {
	Flutterwave FlutterwaveConfig
	Thunes      ThunesConfig
	Stripe      StripeConfig
	PayPal      PayPalConfig
	MTNMomo     MTNMomoConfig
	OrangeMoney OrangeMoneyConfig
	Pesapal     PesapalConfig
	Chipper     ChipperConfig
	CryptoRails CryptoRailsConfig
	CinetPay    CinetPayConfig
	Wave        WaveConfig
	Lygos       LygosConfig
	YellowCard  YellowCardConfig
	Demo        DemoConfig
}

// LoadConfig loads provider configuration prioritizing Vault over environment variables
// This is the recommended configuration loading method for production
func LoadConfig() *Config {
	// Try to initialize Vault client
	vaultClient, err := secrets.NewVaultClient()
	if err != nil {
		log.Printf("[Config] Vault unavailable, falling back to environment variables: %v", err)
		return LoadConfigFromEnv()
	}

	// Seed default secrets if needed (for fresh installations)
	if err := vaultClient.SeedDefaultAggregatorSecrets(); err != nil {
		log.Printf("[Config] Warning: Failed to seed default secrets: %v", err)
	}

	return LoadConfigFromVault(vaultClient)
}

// LoadConfigFromVault loads all configurations from HashiCorp Vault
func LoadConfigFromVault(vaultClient *secrets.VaultClient) *Config {
	log.Println("[Config] Loading configurations from Vault...")

	poolThreshold, _ := strconv.ParseFloat(getEnv("CRYPTO_POOL_THRESHOLD", "500"), 64)
	poolEnabled, _ := strconv.ParseBool(getEnv("CRYPTO_POOL_ENABLED", "true"))

	cfg := &Config{}

	// Load Stripe from Vault
	if stripeSecrets, err := vaultClient.GetStripeSecrets(); err == nil {
		cfg.Stripe = StripeConfig{
			SecretKey:      stripeSecrets.SecretKey,
			PublishableKey: stripeSecrets.PublishableKey,
			WebhookSecret:  stripeSecrets.WebhookSecret,
			BaseURL:        stripeSecrets.BaseURL,
		}
		log.Println("[Config] ✅ Loaded Stripe config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Stripe secrets not available in Vault: %v", err)
		cfg.Stripe = loadStripeFromEnv()
	}

	// Load Flutterwave from Vault
	if flwSecrets, err := vaultClient.GetFlutterwaveSecrets(); err == nil {
		cfg.Flutterwave = FlutterwaveConfig{
			SecretKey:   flwSecrets.SecretKey,
			PublicKey:   flwSecrets.PublicKey,
			EncryptKey:  flwSecrets.EncryptKey,
			BaseURL:     flwSecrets.BaseURL,
			CallbackURL: getEnv("FLUTTERWAVE_CALLBACK_URL", ""), // Callback URL from env (per-deployment)
		}
		log.Println("[Config] ✅ Loaded Flutterwave config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Flutterwave secrets not available in Vault: %v", err)
		cfg.Flutterwave = loadFlutterwaveFromEnv()
	}

	// Load PayPal from Vault
	if paypalSecrets, err := vaultClient.GetPayPalSecrets(); err == nil {
		cfg.PayPal = PayPalConfig{
			ClientID:     paypalSecrets.ClientID,
			ClientSecret: paypalSecrets.ClientSecret,
			WebhookID:    paypalSecrets.WebhookID,
			BaseURL:      paypalSecrets.BaseURL,
			Mode:         paypalSecrets.Mode,
		}
		log.Println("[Config] ✅ Loaded PayPal config from Vault")
	} else {
		log.Printf("[Config] ⚠️ PayPal secrets not available in Vault: %v", err)
	}

	// Load Thunes from Vault
	if thunesSecrets, err := vaultClient.GetThunesSecrets(); err == nil {
		cfg.Thunes = ThunesConfig{
			APIKey:      thunesSecrets.APIKey,
			APISecret:   thunesSecrets.APISecret,
			BaseURL:     thunesSecrets.BaseURL,
			CallbackURL: thunesSecrets.CallbackURL,
		}
		log.Println("[Config] ✅ Loaded Thunes config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Thunes secrets not available in Vault: %v", err)
		cfg.Thunes = loadThunesFromEnv()
	}

	// Load MTN MoMo from Vault
	if mtnSecrets, err := vaultClient.GetMTNMomoSecrets(); err == nil {
		cfg.MTNMomo = MTNMomoConfig{
			SubscriptionKey: mtnSecrets.SubscriptionKey,
			APIUser:         mtnSecrets.APIUser,
			APIKey:          mtnSecrets.APIKey,
			BaseURL:         mtnSecrets.BaseURL,
			Environment:     mtnSecrets.Environment,
			CallbackURL:     getEnv("MTN_CALLBACK_URL", ""),
		}
		log.Println("[Config] ✅ Loaded MTN MoMo config from Vault")
	} else {
		log.Printf("[Config] ⚠️ MTN MoMo secrets not available in Vault: %v", err)
	}

	// Load Orange Money from Vault
	if orangeSecrets, err := vaultClient.GetOrangeMoneySecrets(); err == nil {
		cfg.OrangeMoney = OrangeMoneyConfig{
			ClientID:     orangeSecrets.ClientID,
			ClientSecret: orangeSecrets.ClientSecret,
			MerchantKey:  orangeSecrets.MerchantKey,
			BaseURL:      orangeSecrets.BaseURL,
		}
		log.Println("[Config] ✅ Loaded Orange Money config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Orange Money secrets not available in Vault: %v", err)
	}

	// Load Pesapal from Vault
	if pesapalSecrets, err := vaultClient.GetPesapalSecrets(); err == nil {
		cfg.Pesapal = PesapalConfig{
			ConsumerKey:    pesapalSecrets.ConsumerKey,
			ConsumerSecret: pesapalSecrets.ConsumerSecret,
			BaseURL:        pesapalSecrets.BaseURL,
		}
		log.Println("[Config] ✅ Loaded Pesapal config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Pesapal secrets not available in Vault: %v", err)
	}

	// Load Chipper from Vault
	if chipperSecrets, err := vaultClient.GetChipperSecrets(); err == nil {
		cfg.Chipper = ChipperConfig{
			APIKey:    chipperSecrets.APIKey,
			APISecret: chipperSecrets.APISecret,
			BaseURL:   chipperSecrets.BaseURL,
		}
		log.Println("[Config] ✅ Loaded Chipper config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Chipper secrets not available in Vault: %v", err)
	}

	// Load CinetPay from Vault
	if cinetpaySecrets, err := vaultClient.GetCinetPaySecrets(); err == nil {
		cfg.CinetPay = CinetPayConfig{
			APIKey:    cinetpaySecrets.APIKey,
			SiteID:    cinetpaySecrets.SiteID,
			SecretKey: cinetpaySecrets.SecretKey,
			BaseURL:   cinetpaySecrets.BaseURL,
		}
		log.Println("[Config] ✅ Loaded CinetPay config from Vault")
	} else {
		log.Printf("[Config] ⚠️ CinetPay secrets not available in Vault: %v", err)
	}

	// Load Wave from Vault
	if waveSecrets, err := vaultClient.GetWaveSecrets(); err == nil {
		cfg.Wave = WaveConfig{
			APIKey:      waveSecrets.APIKey,
			MerchantID:  waveSecrets.MerchantID,
			WebhookKey:  waveSecrets.WebhookKey,
			BaseURL:     waveSecrets.BaseURL,
			Environment: waveSecrets.Environment,
		}
		log.Println("[Config] ✅ Loaded Wave config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Wave secrets not available in Vault: %v", err)
	}

	// Load Lygos from Vault
	if lygosSecrets, err := vaultClient.GetLygosSecrets(); err == nil {
		cfg.Lygos = LygosConfig{
			APIKey:     lygosSecrets.APIKey,
			ShopName:   lygosSecrets.ShopName,
			BaseURL:    lygosSecrets.BaseURL,
			WebhookURL: lygosSecrets.WebhookURL,
		}
		log.Println("[Config] ✅ Loaded Lygos config from Vault")
	} else {
		log.Printf("[Config] ⚠️ Lygos secrets not available in Vault: %v", err)
	}

	// Load YellowCard from Vault
	if ycSecrets, err := vaultClient.GetYellowCardSecrets(); err == nil {
		cfg.YellowCard = YellowCardConfig{
			APIKey:      ycSecrets.APIKey,
			SecretKey:   ycSecrets.SecretKey,
			BusinessID:  ycSecrets.BusinessID,
			BaseURL:     ycSecrets.BaseURL,
			WebhookURL:  ycSecrets.WebhookURL,
			Environment: ycSecrets.Environment,
		}
		log.Println("[Config] ✅ Loaded YellowCard config from Vault")
	} else {
		log.Printf("[Config] ⚠️ YellowCard secrets not available in Vault: %v", err)
	}

	// Load Demo config (no secrets needed)
	cfg.Demo = DemoConfig{
		SuccessRate: 0.95,
		DefaultFee:  1.5,
	}
	log.Println("[Config] ✅ Loaded Demo provider config")

	// Load CryptoRails from Vault (Circle, Binance)
	cfg.CryptoRails = loadCryptoRailsFromVault(vaultClient, poolThreshold, poolEnabled)

	log.Println("[Config] Configuration loading complete")
	return cfg
}

// loadCryptoRailsFromVault loads crypto rails config from Vault
func loadCryptoRailsFromVault(vaultClient *secrets.VaultClient, poolThreshold float64, poolEnabled bool) CryptoRailsConfig {
	cfg := CryptoRailsConfig{
		UseInternalPoolThreshold: poolThreshold,
		InternalPoolEnabled:      poolEnabled,
		PreferredStablecoin:      getEnv("PREFERRED_STABLECOIN", "USDC"),
		EthereumWallet:           getEnv("ETH_WALLET_ADDRESS", ""),
		TronWallet:               getEnv("TRX_WALLET_ADDRESS", ""),
		PolygonWallet:            getEnv("MATIC_WALLET_ADDRESS", ""),
	}

	// Load Circle from Vault
	if circleSecrets, err := vaultClient.GetCircleSecrets(); err == nil {
		cfg.CircleAPIKey = circleSecrets.APIKey
		cfg.CircleBaseURL = circleSecrets.BaseURL
		log.Println("[Config] ✅ Loaded Circle config from Vault")
	} else {
		cfg.CircleAPIKey = getEnv("CIRCLE_API_KEY", "")
		cfg.CircleBaseURL = getEnv("CIRCLE_BASE_URL", "https://api.circle.com/v1")
	}

	// Load Binance from Vault
	if binanceSecrets, err := vaultClient.GetBinanceSecrets(); err == nil {
		cfg.BinanceAPIKey = binanceSecrets.APIKey
		cfg.BinanceAPISecret = binanceSecrets.APISecret
		cfg.BinanceBaseURL = binanceSecrets.BaseURL
		log.Println("[Config] ✅ Loaded Binance config from Vault")
	} else {
		cfg.BinanceAPIKey = getEnv("BINANCE_API_KEY", "")
		cfg.BinanceAPISecret = getEnv("BINANCE_API_SECRET", "")
		cfg.BinanceBaseURL = getEnv("BINANCE_BASE_URL", "https://api.binance.com")
	}

	return cfg
}

// LoadConfigFromEnv loads configuration from environment variables only
// Used as fallback when Vault is unavailable
func LoadConfigFromEnv() *Config {
	log.Println("[Config] Loading configurations from environment variables...")

	poolThreshold, _ := strconv.ParseFloat(getEnv("CRYPTO_POOL_THRESHOLD", "500"), 64)
	poolEnabled, _ := strconv.ParseBool(getEnv("CRYPTO_POOL_ENABLED", "true"))

	return &Config{
		Flutterwave: loadFlutterwaveFromEnv(),
		Thunes:      loadThunesFromEnv(),
		Stripe:      loadStripeFromEnv(),
		CryptoRails: CryptoRailsConfig{
			CircleAPIKey:             getEnv("CIRCLE_API_KEY", ""),
			CircleBaseURL:            getEnv("CIRCLE_BASE_URL", "https://api.circle.com/v1"),
			BinanceAPIKey:            getEnv("BINANCE_API_KEY", ""),
			BinanceAPISecret:         getEnv("BINANCE_API_SECRET", ""),
			BinanceBaseURL:           getEnv("BINANCE_BASE_URL", "https://api.binance.com"),
			UseInternalPoolThreshold: poolThreshold,
			InternalPoolEnabled:      poolEnabled,
			PreferredStablecoin:      getEnv("PREFERRED_STABLECOIN", "USDC"),
			EthereumWallet:           getEnv("ETH_WALLET_ADDRESS", ""),
			TronWallet:               getEnv("TRX_WALLET_ADDRESS", ""),
			PolygonWallet:            getEnv("MATIC_WALLET_ADDRESS", ""),
		},
	}
}

func loadStripeFromEnv() StripeConfig {
	return StripeConfig{
		SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		BaseURL:        getEnv("STRIPE_BASE_URL", "https://api.stripe.com/v1"),
	}
}

func loadFlutterwaveFromEnv() FlutterwaveConfig {
	return FlutterwaveConfig{
		SecretKey:   getEnv("FLUTTERWAVE_SECRET_KEY", ""),
		PublicKey:   getEnv("FLUTTERWAVE_PUBLIC_KEY", ""),
		EncryptKey:  getEnv("FLUTTERWAVE_ENCRYPT_KEY", ""),
		BaseURL:     getEnv("FLUTTERWAVE_BASE_URL", "https://api.flutterwave.com/v3"),
		CallbackURL: getEnv("FLUTTERWAVE_CALLBACK_URL", ""),
	}
}

func loadThunesFromEnv() ThunesConfig {
	return ThunesConfig{
		APIKey:      getEnv("THUNES_API_KEY", ""),
		APISecret:   getEnv("THUNES_API_SECRET", ""),
		BaseURL:     getEnv("THUNES_BASE_URL", "https://api.thunes.com/v2"),
		CallbackURL: getEnv("THUNES_CALLBACK_URL", ""),
	}
}

// InitializeRouter creates and configures the zone router with all providers
func InitializeRouter(cfg *Config) *ZoneRouter {
	router := NewZoneRouter()

	// Register Stripe (Europe, North America)
	if cfg.Stripe.SecretKey != "" && !isPlaceholder(cfg.Stripe.SecretKey) {
		stripe := NewStripeProvider(cfg.Stripe)
		router.RegisterProvider(stripe)
		log.Println("[Router] Registered Stripe provider")
	}

	// Register Flutterwave (Africa)
	if cfg.Flutterwave.SecretKey != "" && !isPlaceholder(cfg.Flutterwave.SecretKey) {
		flutterwave := NewFlutterwaveProvider(cfg.Flutterwave)
		router.RegisterProvider(flutterwave)
		log.Println("[Router] Registered Flutterwave provider")
	}

	// Register CinetPay (Africa)
	if cfg.CinetPay.APIKey != "" && !isPlaceholder(cfg.CinetPay.APIKey) {
		cinetpay := NewCinetPayProvider(cfg.CinetPay)
		router.RegisterProvider(cinetpay)
		log.Println("[Router] Registered CinetPay provider")
	}

	// Register Wave (Africa) - Unified
	if cfg.Wave.APIKey != "" && !isPlaceholder(cfg.Wave.APIKey) {
		waveCfg := cfg.Wave
		waveCfg.Name = "wave_money"
		router.RegisterProvider(NewWaveProvider(waveCfg))
		log.Println("[Router] Registered Wave provider (Unified: wave_money)")
	}

	// Register Lygos (Africa)
	if cfg.Lygos.APIKey != "" && !isPlaceholder(cfg.Lygos.APIKey) {
		lygos := NewLygosProvider(cfg.Lygos)
		router.RegisterProvider(lygos)
		log.Println("[Router] Registered Lygos provider")
	}

	// Register YellowCard (Africa)
	if cfg.YellowCard.APIKey != "" && !isPlaceholder(cfg.YellowCard.APIKey) {
		yc := NewYellowCardProvider(cfg.YellowCard)
		router.RegisterProvider(yc)
		log.Println("[Router] Registered YellowCard provider")
	}

	// Register PayPal (Global)
	if cfg.PayPal.ClientID != "" && !isPlaceholder(cfg.PayPal.ClientID) {
		paypal := NewPayPalProvider(cfg.PayPal)
		router.RegisterProvider(paypal)
		log.Println("[Router] Registered PayPal provider")
	}

	// Register Thunes (Global)
	if cfg.Thunes.APIKey != "" && !isPlaceholder(cfg.Thunes.APIKey) {
		thunes := NewThunesProvider(cfg.Thunes)
		router.RegisterProvider(thunes)
		log.Println("[Router] Registered Thunes provider")
	}

	// Register MTN MoMo (Africa) - Unified
	if cfg.MTNMomo.SubscriptionKey != "" && !isPlaceholder(cfg.MTNMomo.SubscriptionKey) {
		mtnCfg := cfg.MTNMomo
		mtnCfg.Name = "mtn_money"
		router.RegisterProvider(NewMTNMomoProvider(mtnCfg))
		log.Println("[Router] Registered MTN MoMo provider (Unified: mtn_money)")
	}

	// Register Orange Money (Africa) - Unified
	if cfg.OrangeMoney.ClientID != "" && !isPlaceholder(cfg.OrangeMoney.ClientID) {
		omCfg := cfg.OrangeMoney
		omCfg.Name = "orange_money"
		router.RegisterProvider(NewOrangeMoneyProvider(omCfg))
		log.Println("[Router] Registered Orange Money provider (Unified: orange_money)")
	}

	// Register Pesapal (East Africa)
	if cfg.Pesapal.ConsumerKey != "" && !isPlaceholder(cfg.Pesapal.ConsumerKey) {
		pesapal := NewPesapalProvider(cfg.Pesapal)
		router.RegisterProvider(pesapal)
		log.Println("[Router] Registered Pesapal provider")
	}

	// Register Chipper Cash (Africa)
	if cfg.Chipper.APIKey != "" && !isPlaceholder(cfg.Chipper.APIKey) {
		chipper := NewChipperProvider(cfg.Chipper)
		router.RegisterProvider(chipper)
		log.Println("[Router] Registered Chipper Cash provider")
	}

	// Register Demo Provider
	if cfg.Demo.DefaultFee > 0 {
		demo := NewDemoProvider(cfg.Demo)
		router.RegisterProvider(demo)
		log.Println("[Router] Registered Demo provider")
	}

	return router
}

// InitializeCryptoRails creates the crypto rails provider
func InitializeCryptoRails(cfg *Config) *CryptoRailsProvider {
	return NewCryptoRailsProvider(cfg.CryptoRails)
}

// InitializeOrchestrator creates the full transfer orchestrator
func InitializeOrchestrator(cfg *Config) *TransferOrchestrator {
	cryptoRails := InitializeCryptoRails(cfg)
	zoneRouter := InitializeRouter(cfg)
	return NewTransferOrchestrator(cryptoRails, zoneRouter)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
