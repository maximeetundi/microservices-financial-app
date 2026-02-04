package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// VaultConfig holds Vault connection details
type VaultConfig struct {
	Address string
	Token   string
}

// ProviderConfigLoader loads payment provider configs from Vault
type ProviderConfigLoader struct {
	vaultConfig VaultConfig
	httpClient  *http.Client
	cache       map[string]interface{}
	cacheMutex  sync.RWMutex
}

// NewProviderConfigLoader creates a new config loader
func NewProviderConfigLoader() *ProviderConfigLoader {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://vault:8200"
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		vaultToken = "dev-token-secure-2024"
	}

	return &ProviderConfigLoader{
		vaultConfig: VaultConfig{
			Address: vaultAddr,
			Token:   vaultToken,
		},
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      make(map[string]interface{}),
	}
}

// getSecret retrieves a secret from Vault
func (p *ProviderConfigLoader) getSecret(path string) (map[string]interface{}, error) {
	// Check cache first
	p.cacheMutex.RLock()
	if cached, ok := p.cache[path]; ok {
		p.cacheMutex.RUnlock()
		return cached.(map[string]interface{}), nil
	}
	p.cacheMutex.RUnlock()

	// Fetch from Vault
	url := fmt.Sprintf("%s/v1/secret/%s", p.vaultConfig.Address, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Vault-Token", p.vaultConfig.Token)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Cache it
	p.cacheMutex.Lock()
	p.cache[path] = result.Data
	p.cacheMutex.Unlock()

	return result.Data, nil
}

func (p *ProviderConfigLoader) getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// LoadFlutterwaveConfig loads Flutterwave configuration from Vault
func (p *ProviderConfigLoader) LoadFlutterwaveConfig(ctx context.Context) (FlutterwaveConfig, error) {
	data, err := p.getSecret("payment/flutterwave")
	if err != nil {
		return FlutterwaveConfig{}, err
	}

	return FlutterwaveConfig{
		PublicKey:   p.getString(data, "public_key"),
		SecretKey:   p.getString(data, "secret_key"),
		EncryptKey:  p.getString(data, "encryption_key"),
		CallbackURL: p.getString(data, "callback_url"),
		BaseURL:     p.getString(data, "base_url"),
	}, nil
}

// LoadCinetPayConfig loads CinetPay configuration from Vault
func (p *ProviderConfigLoader) LoadCinetPayConfig(ctx context.Context) (CinetPayConfig, error) {
	data, err := p.getSecret("payment/cinetpay")
	if err != nil {
		return CinetPayConfig{}, err
	}

	return CinetPayConfig{
		APIKey:  p.getString(data, "api_key"),
		SiteID:  p.getString(data, "site_id"),
		BaseURL: p.getString(data, "base_url"),
	}, nil
}

// LoadPaystackConfig loads Paystack configuration from Vault
func (p *ProviderConfigLoader) LoadPaystackConfig(ctx context.Context) (PaystackConfig, error) {
	data, err := p.getSecret("payment/paystack")
	if err != nil {
		return PaystackConfig{}, err
	}

	return PaystackConfig{
		PublicKey: p.getString(data, "public_key"),
		SecretKey: p.getString(data, "secret_key"),
		BaseURL:   p.getString(data, "base_url"),
	}, nil
}

// LoadStripeConfig loads Stripe configuration from Vault
func (p *ProviderConfigLoader) LoadStripeConfig(ctx context.Context) (StripeConfig, error) {
	data, err := p.getSecret("payment/stripe")
	if err != nil {
		return StripeConfig{}, err
	}

	return StripeConfig{
		PublishableKey: p.getString(data, "public_key"),
		SecretKey:      p.getString(data, "secret_key"),
		WebhookSecret:  p.getString(data, "webhook_secret"),
		BaseURL:        p.getString(data, "base_url"),
	}, nil
}

// LoadWaveConfig loads Wave configuration from Vault
func (p *ProviderConfigLoader) LoadWaveConfig(ctx context.Context) (WaveConfig, error) {
	data, err := p.getSecret("payment/wave")
	if err != nil {
		return WaveConfig{}, err
	}

	return WaveConfig{
		APIKey:  p.getString(data, "api_key"),
		BaseURL: p.getString(data, "base_url"),
	}, nil
}

// LoadPayPalConfig loads PayPal configuration from Vault
func (p *ProviderConfigLoader) LoadPayPalConfig(ctx context.Context) (PayPalConfig, error) {
	data, err := p.getSecret("payment/paypal")
	if err != nil {
		return PayPalConfig{}, err
	}

	return PayPalConfig{
		ClientID:     p.getString(data, "client_id"),
		ClientSecret: p.getString(data, "client_secret"),
		Mode:         p.getString(data, "mode"),
		BaseURL:      p.getString(data, "base_url"),
	}, nil
}

// LoadOrangeMoneyConfig loads Orange Money configuration from Vault
func (p *ProviderConfigLoader) LoadOrangeMoneyConfig(ctx context.Context) (OrangeMoneyConfig, error) {
	data, err := p.getSecret("payment/orange_money")
	if err != nil {
		return OrangeMoneyConfig{}, err
	}

	return OrangeMoneyConfig{
		ClientID:     p.getString(data, "client_id"),
		ClientSecret: p.getString(data, "client_secret"),
		MerchantKey:  p.getString(data, "merchant_key"),
		BaseURL:      p.getString(data, "base_url"),
	}, nil
}

// LoadMTNMoMoConfig loads MTN MoMo configuration from Vault
func (p *ProviderConfigLoader) LoadMTNMoMoConfig(ctx context.Context) (MTNMoMoConfig, error) {
	data, err := p.getSecret("payment/mtn_momo")
	if err != nil {
		return MTNMoMoConfig{}, err
	}

	return MTNMoMoConfig{
		APIUser:         p.getString(data, "api_user"),
		APIKey:          p.getString(data, "api_key"),
		SubscriptionKey: p.getString(data, "subscription_key"),
		BaseURL:         p.getString(data, "base_url"),
		Environment:     "sandbox", // Default to sandbox
	}, nil
}

// InitializeAllProviders loads ALL providers from Vault
func (p *ProviderConfigLoader) InitializeAllProviders(ctx context.Context) (map[string]CollectionProvider, error) {
	providers := make(map[string]CollectionProvider)

	// Flutterwave
	if config, err := p.LoadFlutterwaveConfig(ctx); err == nil && config.SecretKey != "" {
		providers["flutterwave"] = NewFlutterwaveCollectionProvider(config)
	}

	// CinetPay
	if config, err := p.LoadCinetPayConfig(ctx); err == nil && config.APIKey != "" {
		providers["cinetpay"] = NewCinetPayCollectionProvider(config)
	}

	// Paystack
	if config, err := p.LoadPaystackConfig(ctx); err == nil && config.SecretKey != "" {
		providers["paystack"] = NewPaystackCollectionProvider(config)
	}

	// Stripe
	if config, err := p.LoadStripeConfig(ctx); err == nil && config.SecretKey != "" {
		providers["stripe"] = NewStripeCollectionProvider(config)
	}

	// Wave
	if config, err := p.LoadWaveConfig(ctx); err == nil && config.APIKey != "" {
		providers["wave"] = NewWaveCollectionProvider(config)
	}

	// PayPal
	if config, err := p.LoadPayPalConfig(ctx); err == nil && config.ClientID != "" {
		providers["paypal"] = NewPayPalCollectionProvider(config)
	}

	// Orange Money
	if config, err := p.LoadOrangeMoneyConfig(ctx); err == nil && config.ClientID != "" {
		providers["orange_money"] = NewOrangeMoneyCollectionProvider(config)
	}

	// MTN MoMo
	if config, err := p.LoadMTNMoMoConfig(ctx); err == nil && config.APIUser != "" {
		providers["mtn_momo"] = NewMTNMoMoCollectionProvider(config)
	}

	return providers, nil
}
