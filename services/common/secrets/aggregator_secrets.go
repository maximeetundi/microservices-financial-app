package secrets

import (
	"fmt"
	"log"
)

// AggregatorSecrets contains all payment aggregator API credentials
// These are stored in Vault and seeded with placeholder values on first run
type AggregatorSecrets struct {
	Stripe      StripeSecrets      `json:"stripe"`
	Flutterwave FlutterwaveSecrets `json:"flutterwave"`
	PayPal      PayPalSecrets      `json:"paypal"`
	Thunes      ThunesSecrets      `json:"thunes"`
	Wise        WiseSecrets        `json:"wise"`
	Circle      CircleSecrets      `json:"circle"`
	Binance     BinanceSecrets     `json:"binance"`
	Pesapal     PesapalSecrets     `json:"pesapal"`
	Chipper     ChipperSecrets     `json:"chipper"`
	MTN         MTNMomoSecrets     `json:"mtn_momo"`
	Orange      OrangeMoneySecrets `json:"orange_money"`
}

// StripeSecrets - https://dashboard.stripe.com/apikeys
type StripeSecrets struct {
	SecretKey      string `json:"secret_key"`
	PublishableKey string `json:"publishable_key"`
	WebhookSecret  string `json:"webhook_secret"`
	BaseURL        string `json:"base_url"`
}

// FlutterwaveSecrets - https://dashboard.flutterwave.com/settings/api
type FlutterwaveSecrets struct {
	SecretKey   string `json:"secret_key"`
	PublicKey   string `json:"public_key"`
	EncryptKey  string `json:"encrypt_key"`
	BaseURL     string `json:"base_url"`
	WebhookHash string `json:"webhook_hash"`
}

// PayPalSecrets - https://developer.paypal.com/dashboard/applications
type PayPalSecrets struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	WebhookID    string `json:"webhook_id"`
	BaseURL      string `json:"base_url"` // https://api-m.sandbox.paypal.com or https://api-m.paypal.com
	Mode         string `json:"mode"`     // sandbox or live
}

// ThunesSecrets - https://portal.thunes.com
type ThunesSecrets struct {
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	BaseURL     string `json:"base_url"`
	CallbackURL string `json:"callback_url"`
}

// WiseSecrets - https://wise.com/developer/api
type WiseSecrets struct {
	APIToken      string `json:"api_token"`
	ProfileID     string `json:"profile_id"`
	WebhookSecret string `json:"webhook_secret"`
	BaseURL       string `json:"base_url"`
}

// CircleSecrets - https://developers.circle.com
type CircleSecrets struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

// BinanceSecrets - https://www.binance.com/en/my/settings/api-management
type BinanceSecrets struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	BaseURL   string `json:"base_url"`
}

// PesapalSecrets - https://developer.pesapal.com
type PesapalSecrets struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	BaseURL        string `json:"base_url"`
}

// ChipperSecrets - https://developers.chipper.cash
type ChipperSecrets struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	BaseURL   string `json:"base_url"`
}

// MTNMomoSecrets - https://momodeveloper.mtn.com
type MTNMomoSecrets struct {
	SubscriptionKey string `json:"subscription_key"`
	APIUser         string `json:"api_user"`
	APIKey          string `json:"api_key"`
	BaseURL         string `json:"base_url"`
	Environment     string `json:"environment"` // sandbox or production
}

// OrangeMoneySecrets - https://developer.orange.com/apis
type OrangeMoneySecrets struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	MerchantKey  string `json:"merchant_key"`
	BaseURL      string `json:"base_url"`
}

// VaultPaths for each aggregator
const (
	VaultPathAggregators = "secret/aggregators"
	VaultPathStripe      = "secret/aggregators/stripe"
	VaultPathFlutterwave = "secret/aggregators/flutterwave"
	VaultPathPayPal      = "secret/aggregators/paypal"
	VaultPathThunes      = "secret/aggregators/thunes"
	VaultPathWise        = "secret/aggregators/wise"
	VaultPathCircle      = "secret/aggregators/circle"
	VaultPathBinance     = "secret/aggregators/binance"
	VaultPathPesapal     = "secret/aggregators/pesapal"
	VaultPathChipper     = "secret/aggregators/chipper"
	VaultPathMTNMomo     = "secret/aggregators/mtn_momo"
	VaultPathOrangeMoney = "secret/aggregators/orange_money"
)

// SeedDefaultAggregatorSecrets seeds placeholder/test values in Vault at startup
// Admin should replace these with real credentials via Vault UI or CLI
func (v *VaultClient) SeedDefaultAggregatorSecrets() error {
	log.Println("[Vault] Seeding default aggregator secrets (placeholders for admin to configure)...")

	secrets := map[string]map[string]interface{}{
		VaultPathStripe: {
			"secret_key":      "sk_test_REPLACE_WITH_YOUR_STRIPE_SECRET_KEY",
			"publishable_key": "pk_test_REPLACE_WITH_YOUR_STRIPE_PUBLISHABLE_KEY",
			"webhook_secret":  "whsec_REPLACE_WITH_YOUR_STRIPE_WEBHOOK_SECRET",
			"base_url":        "https://api.stripe.com/v1",
		},
		VaultPathFlutterwave: {
			"secret_key":   "FLWSECK_TEST-REPLACE_WITH_YOUR_FLUTTERWAVE_SECRET",
			"public_key":   "FLWPUBK_TEST-REPLACE_WITH_YOUR_FLUTTERWAVE_PUBLIC",
			"encrypt_key":  "FLWENCRYPT-REPLACE_WITH_YOUR_FLUTTERWAVE_ENCRYPT",
			"base_url":     "https://api.flutterwave.com/v3",
			"webhook_hash": "REPLACE_WITH_YOUR_FLUTTERWAVE_WEBHOOK_HASH",
		},
		VaultPathPayPal: {
			"client_id":     "PAYPAL_CLIENT_ID_REPLACE_ME",
			"client_secret": "PAYPAL_CLIENT_SECRET_REPLACE_ME",
			"webhook_id":    "PAYPAL_WEBHOOK_ID_REPLACE_ME",
			"base_url":      "https://api-m.sandbox.paypal.com", // Sandbox by default
			"mode":          "sandbox",
		},
		VaultPathThunes: {
			"api_key":      "THUNES_API_KEY_REPLACE_ME",
			"api_secret":   "THUNES_API_SECRET_REPLACE_ME",
			"base_url":     "https://api.thunes.com/v2",
			"callback_url": "https://your-domain.com/webhooks/thunes",
		},
		VaultPathWise: {
			"api_token":      "WISE_API_TOKEN_REPLACE_ME",
			"profile_id":     "WISE_PROFILE_ID_REPLACE_ME",
			"webhook_secret": "WISE_WEBHOOK_SECRET_REPLACE_ME",
			"base_url":       "https://api.sandbox.transferwise.tech", // Sandbox by default
		},
		VaultPathCircle: {
			"api_key":  "CIRCLE_API_KEY_REPLACE_ME",
			"base_url": "https://api-sandbox.circle.com/v1", // Sandbox by default
		},
		VaultPathBinance: {
			"api_key":    "BINANCE_API_KEY_REPLACE_ME",
			"api_secret": "BINANCE_API_SECRET_REPLACE_ME",
			"base_url":   "https://testnet.binance.vision", // Testnet by default
		},
		VaultPathPesapal: {
			"consumer_key":    "PESAPAL_CONSUMER_KEY_REPLACE_ME",
			"consumer_secret": "PESAPAL_CONSUMER_SECRET_REPLACE_ME",
			"base_url":        "https://cybqa.pesapal.com/pesapalv3", // Sandbox
		},
		VaultPathChipper: {
			"api_key":    "CHIPPER_API_KEY_REPLACE_ME",
			"api_secret": "CHIPPER_API_SECRET_REPLACE_ME",
			"base_url":   "https://api.sandbox.chipper.cash/v1",
		},
		VaultPathMTNMomo: {
			"subscription_key": "MTN_SUBSCRIPTION_KEY_REPLACE_ME",
			"api_user":         "MTN_API_USER_REPLACE_ME",
			"api_key":          "MTN_API_KEY_REPLACE_ME",
			"base_url":         "https://sandbox.momodeveloper.mtn.com",
			"environment":      "sandbox",
		},
		VaultPathOrangeMoney: {
			"client_id":     "ORANGE_CLIENT_ID_REPLACE_ME",
			"client_secret": "ORANGE_CLIENT_SECRET_REPLACE_ME",
			"merchant_key":  "ORANGE_MERCHANT_KEY_REPLACE_ME",
			"base_url":      "https://api.orange.com/orange-money-webpay/dev/v1",
		},
	}

	for path, data := range secrets {
		// Check if secret already exists
		existing, err := v.GetSecret(path)
		if err == nil && existing != nil {
			// Secret exists, check if it's a placeholder
			if secretKey, ok := existing["secret_key"].(string); ok {
				if len(secretKey) > 0 && secretKey[0:7] != "REPLACE" && secretKey[0:3] != "sk_" {
					log.Printf("[Vault] Secret at %s already configured (not placeholder), skipping", path)
					continue
				}
			}
		}

		// Write or update the secret
		if err := v.WriteSecret(path, data); err != nil {
			log.Printf("[Vault] Warning: Failed to seed %s: %v", path, err)
		} else {
			log.Printf("[Vault] Seeded placeholder at %s (admin must configure real keys)", path)
		}
	}

	log.Println("[Vault] ✅ Aggregator secrets seeding complete")
	log.Println("[Vault] ⚠️  IMPORTANT: Replace placeholder values in Vault with real API keys before production use")
	return nil
}

// GetStripeSecrets retrieves Stripe configuration from Vault
func (v *VaultClient) GetStripeSecrets() (*StripeSecrets, error) {
	data, err := v.GetSecret(VaultPathStripe)
	if err != nil {
		return nil, fmt.Errorf("failed to get Stripe secrets: %w", err)
	}
	return &StripeSecrets{
		SecretKey:      getStringValue(data, "secret_key"),
		PublishableKey: getStringValue(data, "publishable_key"),
		WebhookSecret:  getStringValue(data, "webhook_secret"),
		BaseURL:        getStringValue(data, "base_url"),
	}, nil
}

// GetFlutterwaveSecrets retrieves Flutterwave configuration from Vault
func (v *VaultClient) GetFlutterwaveSecrets() (*FlutterwaveSecrets, error) {
	data, err := v.GetSecret(VaultPathFlutterwave)
	if err != nil {
		return nil, fmt.Errorf("failed to get Flutterwave secrets: %w", err)
	}
	return &FlutterwaveSecrets{
		SecretKey:   getStringValue(data, "secret_key"),
		PublicKey:   getStringValue(data, "public_key"),
		EncryptKey:  getStringValue(data, "encrypt_key"),
		BaseURL:     getStringValue(data, "base_url"),
		WebhookHash: getStringValue(data, "webhook_hash"),
	}, nil
}

// GetPayPalSecrets retrieves PayPal configuration from Vault
func (v *VaultClient) GetPayPalSecrets() (*PayPalSecrets, error) {
	data, err := v.GetSecret(VaultPathPayPal)
	if err != nil {
		return nil, fmt.Errorf("failed to get PayPal secrets: %w", err)
	}
	return &PayPalSecrets{
		ClientID:     getStringValue(data, "client_id"),
		ClientSecret: getStringValue(data, "client_secret"),
		WebhookID:    getStringValue(data, "webhook_id"),
		BaseURL:      getStringValue(data, "base_url"),
		Mode:         getStringValue(data, "mode"),
	}, nil
}

// GetThunesSecrets retrieves Thunes configuration from Vault
func (v *VaultClient) GetThunesSecrets() (*ThunesSecrets, error) {
	data, err := v.GetSecret(VaultPathThunes)
	if err != nil {
		return nil, fmt.Errorf("failed to get Thunes secrets: %w", err)
	}
	return &ThunesSecrets{
		APIKey:      getStringValue(data, "api_key"),
		APISecret:   getStringValue(data, "api_secret"),
		BaseURL:     getStringValue(data, "base_url"),
		CallbackURL: getStringValue(data, "callback_url"),
	}, nil
}

// GetWiseSecrets retrieves Wise configuration from Vault
func (v *VaultClient) GetWiseSecrets() (*WiseSecrets, error) {
	data, err := v.GetSecret(VaultPathWise)
	if err != nil {
		return nil, fmt.Errorf("failed to get Wise secrets: %w", err)
	}
	return &WiseSecrets{
		APIToken:      getStringValue(data, "api_token"),
		ProfileID:     getStringValue(data, "profile_id"),
		WebhookSecret: getStringValue(data, "webhook_secret"),
		BaseURL:       getStringValue(data, "base_url"),
	}, nil
}

// GetCircleSecrets retrieves Circle configuration from Vault
func (v *VaultClient) GetCircleSecrets() (*CircleSecrets, error) {
	data, err := v.GetSecret(VaultPathCircle)
	if err != nil {
		return nil, fmt.Errorf("failed to get Circle secrets: %w", err)
	}
	return &CircleSecrets{
		APIKey:  getStringValue(data, "api_key"),
		BaseURL: getStringValue(data, "base_url"),
	}, nil
}

// GetBinanceSecrets retrieves Binance configuration from Vault
func (v *VaultClient) GetBinanceSecrets() (*BinanceSecrets, error) {
	data, err := v.GetSecret(VaultPathBinance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Binance secrets: %w", err)
	}
	return &BinanceSecrets{
		APIKey:    getStringValue(data, "api_key"),
		APISecret: getStringValue(data, "api_secret"),
		BaseURL:   getStringValue(data, "base_url"),
	}, nil
}

// GetPesapalSecrets retrieves Pesapal configuration from Vault
func (v *VaultClient) GetPesapalSecrets() (*PesapalSecrets, error) {
	data, err := v.GetSecret(VaultPathPesapal)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pesapal secrets: %w", err)
	}
	return &PesapalSecrets{
		ConsumerKey:    getStringValue(data, "consumer_key"),
		ConsumerSecret: getStringValue(data, "consumer_secret"),
		BaseURL:        getStringValue(data, "base_url"),
	}, nil
}

// GetChipperSecrets retrieves Chipper Cash configuration from Vault
func (v *VaultClient) GetChipperSecrets() (*ChipperSecrets, error) {
	data, err := v.GetSecret(VaultPathChipper)
	if err != nil {
		return nil, fmt.Errorf("failed to get Chipper secrets: %w", err)
	}
	return &ChipperSecrets{
		APIKey:    getStringValue(data, "api_key"),
		APISecret: getStringValue(data, "api_secret"),
		BaseURL:   getStringValue(data, "base_url"),
	}, nil
}

// GetMTNMomoSecrets retrieves MTN Mobile Money configuration from Vault
func (v *VaultClient) GetMTNMomoSecrets() (*MTNMomoSecrets, error) {
	data, err := v.GetSecret(VaultPathMTNMomo)
	if err != nil {
		return nil, fmt.Errorf("failed to get MTN MoMo secrets: %w", err)
	}
	return &MTNMomoSecrets{
		SubscriptionKey: getStringValue(data, "subscription_key"),
		APIUser:         getStringValue(data, "api_user"),
		APIKey:          getStringValue(data, "api_key"),
		BaseURL:         getStringValue(data, "base_url"),
		Environment:     getStringValue(data, "environment"),
	}, nil
}

// GetOrangeMoneySecrets retrieves Orange Money configuration from Vault
func (v *VaultClient) GetOrangeMoneySecrets() (*OrangeMoneySecrets, error) {
	data, err := v.GetSecret(VaultPathOrangeMoney)
	if err != nil {
		return nil, fmt.Errorf("failed to get Orange Money secrets: %w", err)
	}
	return &OrangeMoneySecrets{
		ClientID:     getStringValue(data, "client_id"),
		ClientSecret: getStringValue(data, "client_secret"),
		MerchantKey:  getStringValue(data, "merchant_key"),
		BaseURL:      getStringValue(data, "base_url"),
	}, nil
}

// Helper function to safely extract string values from map
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
