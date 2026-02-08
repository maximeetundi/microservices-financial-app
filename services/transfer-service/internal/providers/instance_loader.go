package providers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
)

// InstanceCredentialsClient interface for fetching instance credentials
// This interface is implemented by services.AdminClient to avoid import cycles
type InstanceCredentialsClient interface {
	GetBestInstanceWithCredentials(providerCode, country string, amount float64, currency string, operation string, isTestMode *bool) (*models.AggregatorInstanceWithDetails, error)
	GetInstanceByID(instanceID string) (*models.AggregatorInstanceWithDetails, error)
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// InstanceBasedProviderLoader loads providers from admin-service API
// Credentials are fetched via secure service-to-service communication
type InstanceBasedProviderLoader struct {
	credentialsClient InstanceCredentialsClient
}

// NewInstanceBasedProviderLoader creates a new instance-based loader
// The client parameter must implement InstanceCredentialsClient (e.g., services.AdminClient)
func NewInstanceBasedProviderLoader(client InstanceCredentialsClient) *InstanceBasedProviderLoader {
	log.Printf("[InstanceLoader] ✅ Initialized - credentials fetched from admin-service API")
	return &InstanceBasedProviderLoader{
		credentialsClient: client,
	}
}

func (l *InstanceBasedProviderLoader) CredentialsClient() InstanceCredentialsClient {
	return l.credentialsClient
}

// getCredentials returns credentials from database (api_credentials JSONB column)
func (l *InstanceBasedProviderLoader) getCredentials(instance *models.AggregatorInstanceWithDetails) (map[string]string, error) {
	log.Printf("[InstanceLoader] 🔍 Loading credentials for: %s (instance: %s)", instance.ProviderCode, instance.InstanceName)

	// Read credentials from database
	if len(instance.APICredentials) == 0 {
		log.Printf("[InstanceLoader] ❌ No credentials found in database for %s", instance.ProviderCode)
		return nil, fmt.Errorf("no credentials configured for %s - add them in admin panel", instance.ProviderCode)
	}

	// Filter out placeholders and empty values
	validCreds := make(map[string]string)
	for k, v := range instance.APICredentials {
		if v != "" && !isPlaceholder(v) {
			validCreds[k] = v
		}
	}

	if len(validCreds) == 0 {
		log.Printf("[InstanceLoader] ❌ All credentials are placeholders for %s", instance.ProviderCode)
		log.Printf("[InstanceLoader] 💡 Configure real credentials in Admin Panel -> Agrégateurs -> %s", instance.ProviderCode)
		return nil, fmt.Errorf("credentials are placeholders for %s - configure real API keys in admin panel", instance.ProviderCode)
	}

	// Log keys found (not values for security)
	keys := make([]string, 0, len(validCreds))
	for k := range validCreds {
		keys = append(keys, k)
	}
	log.Printf("[InstanceLoader] ✅ Loaded %d credentials from DB: %v", len(validCreds), keys)

	return validCreds, nil
}

// isPlaceholder checks if a credential value is a placeholder
func isPlaceholder(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return true
	}
	upper := strings.ToUpper(trimmed)
	placeholders := []string{
		"REPLACE_ME", "REPLACE_WITH", "YOUR_", "PLACEHOLDER",
		"xxx", "XXX", "test_", "TEST_", "dummy", "DUMMY",
		"FLWSECK_", "FLWPUBK_", "FLWSECKTEST", "FLWPUBKTEST",
		"SK_TEST_", "PK_TEST_",
		"SANDBOX_", "LIVE_",
	}
	for _, p := range placeholders {
		if len(upper) >= len(p) && (upper[:len(p)] == p || contains(upper, p)) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// LoadProviderFromInstance creates a provider instance from aggregator instance config
func (l *InstanceBasedProviderLoader) LoadProviderFromInstance(
	ctx context.Context,
	instance *models.AggregatorInstanceWithDetails,
) (CollectionProvider, error) {

	log.Printf("[InstanceLoader] ════════════════════════════════════════════════════")
	log.Printf("[InstanceLoader] 🚀 LOADING PROVIDER: %s", instance.ProviderCode)
	log.Printf("[InstanceLoader]    Instance: %s (%s)", instance.InstanceName, instance.ID)
	log.Printf("[InstanceLoader]    Enabled: %v | Test Mode: %v", instance.Enabled, instance.IsTestMode)
	log.Printf("[InstanceLoader] ════════════════════════════════════════════════════")

	switch instance.ProviderCode {
	case "flutterwave":
		p, err := l.loadFlutterwaveFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "cinetpay":
		p, err := l.loadCinetPayFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "paystack":
		p, err := l.loadPaystackFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "stripe":
		p, err := l.loadStripeFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "wave", "wave_money", "wave_ci", "wave_sn":
		p, err := l.loadWaveFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "paypal":
		p, err := l.loadPayPalFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "orange_money":
		p, err := l.loadOrangeMoneyFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "mtn_momo", "mtn_money":
		p, err := l.loadMTNMoMoFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "lygos":
		p, err := l.loadLygosFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "yellowcard":
		p, err := l.loadYellowCardFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "moov_money", "moov":
		p, err := l.loadMoovMoneyFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "fedapay":
		p, err := l.loadFedaPayFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "pawapay":
		p, err := l.loadPawapayFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "hubtel":
		p, err := l.loadHubtelFromInstance(instance)
		if err != nil {
			return l.fallbackToDemoIfTestMode(instance, err)
		}
		return p, nil
	case "demo":
		return l.loadDemoFromInstance(instance)
	default:
		log.Printf("[InstanceLoader] ⚠️ Unknown provider: %s", instance.ProviderCode)
		return nil, fmt.Errorf("unknown provider: %s", instance.ProviderCode)
	}
}

func (l *InstanceBasedProviderLoader) fallbackToDemoIfTestMode(instance *models.AggregatorInstanceWithDetails, cause error) (CollectionProvider, error) {
	if instance == nil {
		return nil, cause
	}

	// Default: enable demo fallback in dev/fresh installs so deposits work after reset.
	env := strings.TrimSpace(strings.ToLower(os.Getenv("DEMO_FALLBACK_ON_CREDENTIAL_ERROR")))
	allowFallback := env == "" || env == "1" || env == "true" || env == "yes"

	if instance.IsTestMode || (allowFallback && isCredentialRelatedError(cause)) {
		log.Printf("[InstanceLoader] 🧪 Falling back to Demo provider for %s (instance: %s): %v", instance.ProviderCode, instance.ID, cause)
		return NewDemoCollectionProvider(DemoConfig{SuccessRate: 0.95, DefaultFee: 1.5}), nil
	}

	return nil, cause
}

func isCredentialRelatedError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())

	// Broad matching because providers return many different strings.
	needles := []string{
		"credential",
		"api key",
		"client_id",
		"clientsecret",
		"client_secret",
		"invalid_client",
		"invalid authorization",
		"invalid authorization key",
		"invalid api key",
		"unauthorized",
		"token request failed with status 401",
	}
	for _, n := range needles {
		if strings.Contains(msg, n) {
			return true
		}
	}
	return false
}

func (l *InstanceBasedProviderLoader) loadFlutterwaveFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Flutterwave credentials: %w", err)
	}

	config := FlutterwaveConfig{
		PublicKey:   creds["public_key"],
		SecretKey:   creds["secret_key"],
		EncryptKey:  creds["encryption_key"],
		CallbackURL: creds["callback_url"],
		BaseURL:     "https://api.flutterwave.com/v3",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewFlutterwaveCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadCinetPayFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing CinetPay provider...")

	creds, err := l.getCredentials(instance)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ CinetPay credential error: %v", err)
		return nil, fmt.Errorf("failed to get CinetPay credentials: %w", err)
	}

	// Validate required credentials
	apiKey := creds["api_key"]
	siteID := creds["site_id"]

	if apiKey == "" {
		log.Printf("[InstanceLoader] ❌ CinetPay: Missing 'api_key' credential")
		return nil, fmt.Errorf("CinetPay requires 'api_key' credential")
	}
	if siteID == "" {
		log.Printf("[InstanceLoader] ❌ CinetPay: Missing 'site_id' credential")
		return nil, fmt.Errorf("CinetPay requires 'site_id' credential")
	}

	log.Printf("[InstanceLoader] ✅ CinetPay credentials validated:")
	log.Printf("[InstanceLoader]    API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	log.Printf("[InstanceLoader]    Site ID: %s", siteID)

	config := CinetPayConfig{
		APIKey:  apiKey,
		SiteID:  siteID,
		BaseURL: "https://api-checkout.cinetpay.com/v2",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewCinetPayCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPaystackFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Paystack credentials: %w", err)
	}

	config := PaystackConfig{
		PublicKey: creds["public_key"],
		SecretKey: creds["secret_key"],
		BaseURL:   "https://api.paystack.co",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewPaystackCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadStripeFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Stripe credentials: %w", err)
	}

	config := StripeConfig{
		PublishableKey: creds["publishable_key"],
		SecretKey:      creds["secret_key"],
		WebhookSecret:  creds["webhook_secret"],
		BaseURL:        "https://api.stripe.com/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewStripeCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadWaveFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Wave provider...")

	creds, err := l.getCredentials(instance)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ Wave credential error: %v", err)
		return nil, fmt.Errorf("failed to get Wave credentials: %w", err)
	}

	apiKey := creds["api_key"]
	if apiKey == "" {
		log.Printf("[InstanceLoader] ❌ Wave: Missing 'api_key' credential")
		return nil, fmt.Errorf("Wave requires 'api_key' credential")
	}

	log.Printf("[InstanceLoader] ✅ Wave credentials validated (API Key: %s...)", apiKey[:min(8, len(apiKey))])

	config := WaveConfig{
		APIKey:  apiKey,
		BaseURL: "https://api.wave.com/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewWaveCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPayPalFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing PayPal provider...")

	creds, err := l.getCredentials(instance)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ PayPal credential error: %v", err)
		return nil, fmt.Errorf("failed to get PayPal credentials: %w", err)
	}

	clientID := creds["client_id"]
	clientSecret := creds["client_secret"]

	if clientID == "" {
		log.Printf("[InstanceLoader] ❌ PayPal: Missing 'client_id' credential")
		return nil, fmt.Errorf("PayPal requires 'client_id' credential")
	}
	if clientSecret == "" {
		log.Printf("[InstanceLoader] ❌ PayPal: Missing 'client_secret' credential")
		return nil, fmt.Errorf("PayPal requires 'client_secret' credential")
	}

	log.Printf("[InstanceLoader] ✅ PayPal credentials found:")
	log.Printf("[InstanceLoader]    Client ID: %s...%s", clientID[:min(8, len(clientID))], clientID[max(0, len(clientID)-4):])

	// Determine mode and base URL
	mode := creds["mode"]
	if mode == "" {
		if instance.IsTestMode {
			mode = "sandbox"
		} else {
			mode = "live"
		}
	}

	baseURL := creds["base_url"]
	if baseURL == "" {
		if mode == "sandbox" {
			baseURL = "https://api-m.sandbox.paypal.com"
		} else {
			baseURL = "https://api-m.paypal.com"
		}
	}

	config := PayPalConfig{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		Mode:         mode,
		BaseURL:      baseURL,
	}

	log.Printf("[InstanceLoader] PayPal config: mode=%s, baseURL=%s, clientID=%s...",
		config.Mode, config.BaseURL, truncate(config.ClientID, 10))

	return NewPayPalCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadOrangeMoneyFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Orange Money credentials: %w", err)
	}

	config := OrangeMoneyConfig{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		MerchantKey:  creds["merchant_key"],
		BaseURL:      "https://api.orange.com/orange-money-webpay",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewOrangeMoneyCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadMTNMoMoFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get MTN MoMo credentials: %w", err)
	}

	mode := strings.TrimSpace(strings.ToLower(creds["environment"]))
	if mode == "" {
		if instance.IsTestMode {
			mode = "sandbox"
		} else {
			mode = "production"
		}
	}
	if mode == "prod" || mode == "live" {
		mode = "production"
	}
	if mode != "production" {
		mode = "sandbox"
	}

	targetEnv := strings.TrimSpace(creds["target_environment"])
	if targetEnv == "" {
		// Backward compatible default: use mode.
		// In production, you should set target_environment to the operator environment required by MTN.
		targetEnv = mode
	}

	baseURL := creds["base_url"]
	if baseURL == "" {
		if mode == "sandbox" {
			baseURL = "https://sandbox.momodeveloper.mtn.com"
		} else {
			baseURL = "https://proxy.momoapi.mtn.com"
		}
	}

	config := MTNMomoConfig{
		APIUser:         creds["api_user"],
		APIKey:          creds["api_key"],
		SubscriptionKey: creds["subscription_key"],
		BaseURL:         baseURL,
		Mode:            mode,
		TargetEnvironment: targetEnv,
	}

	return NewMTNMoMoCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadLygosFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Lygos credentials: %w", err)
	}

	config := LygosConfig{
		APIKey:   creds["api_key"],
		ShopName: creds["shop_name"],
		BaseURL:  "https://api.lygosapp.com/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewLygosCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadYellowCardFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get YellowCard credentials: %w", err)
	}

	config := YellowCardConfig{
		APIKey:     creds["api_key"],
		SecretKey:  creds["secret_key"],
		BusinessID: creds["business_id"],
		BaseURL:    "https://api.yellowcard.io/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewYellowCardCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadMoovMoneyFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get Moov Money credentials: %w", err)
	}

	config := MoovMoneyConfig{
		APIKey:      creds["api_key"],
		MerchantKey: creds["merchant_key"],
		BaseURL:     "https://api.moov-africa.bj/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewMoovMoneyCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadFedaPayFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to get FedaPay credentials: %w", err)
	}

	config := FedaPayConfig{
		APIKey:    creds["api_key"],
		PublicKey: creds["public_key"],
		SecretKey: creds["secret_key"],
		BaseURL:   "https://api.fedapay.com/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewFedaPayCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPawapayFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Pawapay provider...")

	creds, err := l.getCredentials(instance)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ Pawapay credential error: %v", err)
		return nil, fmt.Errorf("failed to get Pawapay credentials: %w", err)
	}

	apiKey := creds["api_key"]
	if apiKey == "" {
		log.Printf("[InstanceLoader] ❌ Pawapay: Missing 'api_key' credential")
		return nil, fmt.Errorf("Pawapay requires 'api_key' credential")
	}

	log.Printf("[InstanceLoader] ✅ Pawapay credentials validated")

	config := PawapayConfig{
		APIKey:  apiKey,
		BaseURL: "https://api.sandbox.pawapay.cloud",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewPawapayCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadHubtelFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Hubtel provider...")

	creds, err := l.getCredentials(instance)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ Hubtel credential error: %v", err)
		return nil, fmt.Errorf("failed to get Hubtel credentials: %w", err)
	}

	clientID := creds["client_id"]
	clientSecret := creds["client_secret"]

	if clientID == "" {
		log.Printf("[InstanceLoader] ❌ Hubtel: Missing 'client_id' credential")
		return nil, fmt.Errorf("Hubtel requires 'client_id' credential")
	}
	if clientSecret == "" {
		log.Printf("[InstanceLoader] ❌ Hubtel: Missing 'client_secret' credential")
		return nil, fmt.Errorf("Hubtel requires 'client_secret' credential")
	}

	log.Printf("[InstanceLoader] ✅ Hubtel credentials validated")

	config := HubtelConfig{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		MerchantAccount: creds["merchant_key"],
		BaseURL:         "https://api.hubtel.com",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewHubtelCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadDemoFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := DemoConfig{
		SuccessRate: 0.95,
		DefaultFee:  1.5,
	}
	return NewDemoCollectionProvider(config), nil
}

// GetBestProviderForDeposit selects the best instance for a deposit
// based on provider code, country, amount, and instance availability
// Fetches instance with credentials from admin-service via internal API
func (l *InstanceBasedProviderLoader) GetBestProviderForDeposit(
	ctx context.Context,
	providerCode string,
	country string,
	amount float64,
	isTestMode *bool,
) (CollectionProvider, *models.AggregatorInstanceWithDetails, error) {

	log.Printf("[InstanceLoader] 🔍 Getting best instance for provider=%s, country=%s, amount=%.2f",
		providerCode, country, amount)

	// Get best instance with credentials from admin-service API
	instance, err := l.credentialsClient.GetBestInstanceWithCredentials(providerCode, country, amount, "XOF", "deposit", isTestMode)
	if err != nil {
		log.Printf("[InstanceLoader] ❌ Failed to get instance from admin-service: %v", err)
		return nil, nil, fmt.Errorf("no available instance for provider %s: %w", providerCode, err)
	}

	// Check if instance is paused
	if instance.IsPaused {
		reason := "Service temporairement indisponible"
		if instance.PauseReason != nil && *instance.PauseReason != "" {
			reason = *instance.PauseReason
		}
		log.Printf("[InstanceLoader] ⚠️ Instance %s is paused: %s", instance.ID, reason)
		return nil, nil, fmt.Errorf("instance_paused: %s", reason)
	}

	// Check if instance is enabled
	if !instance.Enabled {
		log.Printf("[InstanceLoader] ⚠️ Instance %s is not enabled", instance.ID)
		return nil, nil, fmt.Errorf("provider %s instance not enabled", providerCode)
	}

	log.Printf("[InstanceLoader] ✅ Got instance %s (%s) with %d credentials",
		instance.ID, instance.InstanceName, len(instance.APICredentials))

	// Load provider from this instance
	provider, err := l.LoadProviderFromInstance(ctx, instance)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load provider from instance: %w", err)
	}

	return provider, instance, nil
}

// RecordTransactionUsage is now handled by admin-service via API
// This is a no-op placeholder for backward compatibility
func (l *InstanceBasedProviderLoader) RecordTransactionUsage(
	ctx context.Context,
	instanceID string,
	transactionID string,
	amount float64,
	currency string,
	status string,
	providerRef string,
) error {
	// Usage tracking is handled by admin-service when it serves the instance
	// The admin-service increments request_count and last_used_at automatically
	log.Printf("[InstanceLoader] Transaction recorded: instance=%s, tx=%s, amount=%.2f %s, status=%s",
		instanceID, transactionID, amount, currency, status)
	return nil
}

// truncate truncates a string for logging purposes
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
