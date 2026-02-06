package providers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crypto-bank/microservices-financial-app/services/common/secrets"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

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

// InstanceBasedProviderLoader loads providers from database instances
// and retrieves credentials from Vault when not stored in database
type InstanceBasedProviderLoader struct {
	instanceRepo *repository.AggregatorInstanceRepository
	vaultClient  *secrets.VaultClient
}

// NewInstanceBasedProviderLoader creates a new instance-based loader
func NewInstanceBasedProviderLoader(instanceRepo *repository.AggregatorInstanceRepository) *InstanceBasedProviderLoader {
	loader := &InstanceBasedProviderLoader{
		instanceRepo: instanceRepo,
	}

	// Initialize Vault client
	vaultClient, err := secrets.NewVaultClient()
	if err != nil {
		log.Printf("[InstanceLoader] Warning: Failed to initialize Vault client: %v", err)
		log.Printf("[InstanceLoader] Will fall back to database credentials only")
	} else {
		loader.vaultClient = vaultClient
		log.Printf("[InstanceLoader] Vault client initialized successfully")
	}

	return loader
}

// getCredentialsFromVault loads credentials from Vault for a given provider
func (l *InstanceBasedProviderLoader) getCredentialsFromVault(providerCode string, vaultPath string) (map[string]string, error) {
	if l.vaultClient == nil {
		log.Printf("[InstanceLoader] ⚠️ VAULT CLIENT NOT INITIALIZED - Cannot load credentials for %s", providerCode)
		return nil, fmt.Errorf("vault client not initialized")
	}

	// Use custom vault path if provided, otherwise use default
	path := vaultPath
	if path == "" {
		path = fmt.Sprintf("secret/aggregators/%s", providerCode)
	}

	// Convert path for KV v2: secret/aggregators/x -> secret/data/aggregators/x
	kvV2Path := convertToKVv2Path(path)
	log.Printf("[InstanceLoader] 🔑 Loading credentials from Vault:")
	log.Printf("[InstanceLoader]    Provider: %s", providerCode)
	log.Printf("[InstanceLoader]    Path (KV v1): %s", path)
	log.Printf("[InstanceLoader]    Path (KV v2): %s", kvV2Path)

	// Try KV v2 path first
	data, err := l.vaultClient.GetSecret(kvV2Path)
	if err != nil {
		// Fall back to KV v1 path
		log.Printf("[InstanceLoader] ⚠️ KV v2 path failed: %v", err)
		log.Printf("[InstanceLoader] 🔄 Trying KV v1 path: %s", path)
		data, err = l.vaultClient.GetSecret(path)
		if err != nil {
			log.Printf("[InstanceLoader] ❌ VAULT ERROR: Failed to get secret from both KV v1 and v2 paths")
			log.Printf("[InstanceLoader]    Error: %v", err)
			log.Printf("[InstanceLoader]    💡 TIP: Configure credentials in admin panel at /dashboard/aggregators")
			return nil, fmt.Errorf("failed to get secret from vault: %w", err)
		}
	}

	// Handle KV v2 nested data structure
	actualData := extractVaultData(data)

	// Convert map[string]interface{} to map[string]string
	creds := make(map[string]string)
	for key, value := range actualData {
		if str, ok := value.(string); ok {
			creds[key] = str
		}
	}

	// Log credential keys (NOT values for security)
	log.Printf("[InstanceLoader] ✅ Loaded %d credentials from Vault for %s", len(creds), providerCode)
	if len(creds) > 0 {
		keys := make([]string, 0, len(creds))
		for k := range creds {
			keys = append(keys, k)
		}
		log.Printf("[InstanceLoader]    Keys found: %v", keys)
	} else {
		log.Printf("[InstanceLoader] ⚠️ WARNING: Vault path exists but contains no credentials!")
	}

	return creds, nil
}

// convertToKVv2Path converts a KV v1 style path to KV v2 style
// secret/aggregators/paypal -> secret/data/aggregators/paypal
func convertToKVv2Path(path string) string {
	// If path already contains /data/, return as is
	if len(path) > 7 && path[7:11] == "data" {
		return path
	}

	// Find the first slash (mount point separator)
	for i, c := range path {
		if c == '/' {
			// Insert /data after the mount point
			return path[:i] + "/data" + path[i:]
		}
	}
	return path
}

// extractVaultData handles both KV v1 and KV v2 data structures
func extractVaultData(data map[string]interface{}) map[string]interface{} {
	// KV v2 stores actual data under a "data" key
	if nestedData, ok := data["data"]; ok {
		if nested, ok := nestedData.(map[string]interface{}); ok {
			return nested
		}
	}
	// KV v1 or already extracted data
	return data
}

// getCredentials returns credentials from Vault using instance's vault path, or falls back to database credentials
func (l *InstanceBasedProviderLoader) getCredentials(instance *models.AggregatorInstanceWithDetails, vaultPath string) (map[string]string, error) {
	log.Printf("[InstanceLoader] 🔍 Getting credentials for provider: %s (instance: %s)", instance.ProviderCode, instance.InstanceName)
	log.Printf("[InstanceLoader]    Vault path: %s", vaultPath)
	log.Printf("[InstanceLoader]    DB credentials count: %d", len(instance.APICredentials))

	// Always try Vault first if we have a vault path (this is the recommended approach)
	if vaultPath != "" && l.vaultClient != nil {
		creds, err := l.getCredentialsFromVault(instance.ProviderCode, vaultPath)
		if err == nil && len(creds) > 0 {
			// Validate that credentials are not empty/placeholder
			validCreds := make(map[string]string)
			for k, v := range creds {
				if v != "" && !isPlaceholder(v) {
					validCreds[k] = v
				} else {
					log.Printf("[InstanceLoader] ⚠️ Skipping placeholder/empty credential: %s", k)
				}
			}
			if len(validCreds) > 0 {
				log.Printf("[InstanceLoader] ✅ Using %d valid credentials from Vault for %s", len(validCreds), instance.ProviderCode)
				return validCreds, nil
			}
			log.Printf("[InstanceLoader] ⚠️ Vault credentials exist but all are placeholders/empty")
		} else if err != nil {
			log.Printf("[InstanceLoader] ⚠️ Vault load failed for %s: %v", instance.ProviderCode, err)
		}
		log.Printf("[InstanceLoader] 🔄 Checking database fallback...")
	} else {
		if vaultPath == "" {
			log.Printf("[InstanceLoader] ⚠️ No vault path configured for instance %s", instance.InstanceName)
		}
		if l.vaultClient == nil {
			log.Printf("[InstanceLoader] ⚠️ Vault client not available")
		}
	}

	// Fall back to database credentials if Vault fails
	if len(instance.APICredentials) > 0 {
		log.Printf("[InstanceLoader] 📦 Checking database credentials...")
		// Verify credentials are not placeholders
		validCreds := make(map[string]string)
		for k, v := range instance.APICredentials {
			if v != "" && !isPlaceholder(v) {
				validCreds[k] = v
			} else {
				log.Printf("[InstanceLoader]    Skipping placeholder/empty DB credential: %s = '%s'", k, maskCredential(v))
			}
		}
		if len(validCreds) > 0 {
			log.Printf("[InstanceLoader] ✅ Using %d valid credentials from database for %s", len(validCreds), instance.ProviderCode)
			return validCreds, nil
		}
		log.Printf("[InstanceLoader] ❌ Database credentials exist but all are placeholders/empty")
	}

	log.Printf("[InstanceLoader] ❌ NO VALID CREDENTIALS FOUND for %s", instance.ProviderCode)
	log.Printf("[InstanceLoader] ═══════════════════════════════════════════════════════")
	log.Printf("[InstanceLoader] 💡 TO FIX THIS ISSUE:")
	log.Printf("[InstanceLoader]    1. Go to Admin Dashboard -> Agrégateurs -> %s", instance.ProviderCode)
	log.Printf("[InstanceLoader]    2. Click on the instance '%s'", instance.InstanceName)
	log.Printf("[InstanceLoader]    3. Go to 'Credentials API' tab")
	log.Printf("[InstanceLoader]    4. Enter your real API credentials from %s", instance.ProviderCode)
	log.Printf("[InstanceLoader]    5. Save and retry the transaction")
	log.Printf("[InstanceLoader] ═══════════════════════════════════════════════════════")

	return nil, fmt.Errorf("no valid credentials found for %s (vault path: %s) - configure in admin panel", instance.ProviderCode, vaultPath)
}

// maskCredential masks a credential for safe logging
func maskCredential(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "***" + value[len(value)-2:]
}

// isPlaceholder checks if a credential value is a placeholder
func isPlaceholder(value string) bool {
	placeholders := []string{
		"REPLACE_ME", "REPLACE_WITH", "YOUR_", "PLACEHOLDER",
		"xxx", "XXX", "test_", "TEST_", "dummy", "DUMMY",
	}
	for _, p := range placeholders {
		if len(value) >= len(p) && (value[:len(p)] == p || contains(value, p)) {
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
	log.Printf("[InstanceLoader]    Instance ID: %s", instance.ID)
	log.Printf("[InstanceLoader]    Instance Name: %s", instance.InstanceName)
	log.Printf("[InstanceLoader]    Enabled: %v", instance.Enabled)
	log.Printf("[InstanceLoader]    Is Test Mode: %v", instance.IsTestMode)
	log.Printf("[InstanceLoader] ════════════════════════════════════════════════════")

	// Use VaultSecretPath from instance if available, otherwise fall back to default
	vaultPath := instance.VaultSecretPath
	if vaultPath == "" {
		// Check environment variable as fallback
		envKey := fmt.Sprintf("%s_VAULT_PATH", instance.ProviderCode)
		vaultPath = os.Getenv(envKey)
		if vaultPath != "" {
			log.Printf("[InstanceLoader] 📍 Vault path from env var %s: %s", envKey, vaultPath)
		}
	}
	if vaultPath == "" {
		// Final fallback to default path
		vaultPath = fmt.Sprintf("secret/aggregators/%s", instance.ProviderCode)
		log.Printf("[InstanceLoader] 📍 Using default vault path: %s", vaultPath)
	} else {
		log.Printf("[InstanceLoader] 📍 Using configured vault path: %s", vaultPath)
	}

	switch instance.ProviderCode {
	case "flutterwave":
		return l.loadFlutterwaveFromInstance(instance, vaultPath)
	case "cinetpay":
		return l.loadCinetPayFromInstance(instance, vaultPath)
	case "paystack":
		return l.loadPaystackFromInstance(instance, vaultPath)
	case "stripe":
		return l.loadStripeFromInstance(instance, vaultPath)
	case "wave", "wave_money", "wave_ci", "wave_sn":
		return l.loadWaveFromInstance(instance, vaultPath)
	case "paypal":
		return l.loadPayPalFromInstance(instance, vaultPath)
	case "orange_money":
		return l.loadOrangeMoneyFromInstance(instance, vaultPath)
	case "mtn_momo", "mtn_money":
		return l.loadMTNMoMoFromInstance(instance, vaultPath)
	case "lygos":
		return l.loadLygosFromInstance(instance, vaultPath)
	case "yellowcard":
		return l.loadYellowCardFromInstance(instance, vaultPath)
	case "moov_money", "moov":
		return l.loadMoovMoneyFromInstance(instance, vaultPath)
	case "fedapay":
		return l.loadFedaPayFromInstance(instance, vaultPath)
	case "pawapay":
		return l.loadPawapayFromInstance(instance, vaultPath)
	case "hubtel":
		return l.loadHubtelFromInstance(instance, vaultPath)
	case "demo":
		return l.loadDemoFromInstance(instance)
	default:
		log.Printf("[InstanceLoader] ⚠️ Unknown provider code: %s - attempting generic load", instance.ProviderCode)
		return nil, fmt.Errorf("unknown provider: %s", instance.ProviderCode)
	}
}

func (l *InstanceBasedProviderLoader) loadFlutterwaveFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadCinetPayFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing CinetPay provider...")

	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadPaystackFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadStripeFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadWaveFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Wave provider...")

	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadPayPalFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing PayPal provider...")

	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadOrangeMoneyFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadMTNMoMoFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get MTN MoMo credentials: %w", err)
	}

	environment := creds["environment"]
	if environment == "" {
		if instance.IsTestMode {
			environment = "sandbox"
		} else {
			environment = "prod"
		}
	}

	baseURL := creds["base_url"]
	if baseURL == "" {
		if environment == "sandbox" {
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
		Environment:     environment,
	}

	return NewMTNMoMoCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadLygosFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadYellowCardFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadMoovMoneyFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadFedaPayFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadPawapayFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Pawapay provider...")

	creds, err := l.getCredentials(instance, vaultPath)
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

func (l *InstanceBasedProviderLoader) loadHubtelFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	log.Printf("[InstanceLoader] 🔧 Initializing Hubtel provider...")

	creds, err := l.getCredentials(instance, vaultPath)
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
func (l *InstanceBasedProviderLoader) GetBestProviderForDeposit(
	ctx context.Context,
	providerCode string,
	country string,
	amount float64,
) (CollectionProvider, *models.AggregatorInstanceWithDetails, error) {

	// Get best instance from DB (already excludes paused instances)
	instance, err := l.instanceRepo.GetBestInstanceForProvider(ctx, providerCode, country, amount)
	if err != nil {
		return nil, nil, fmt.Errorf("no available instance for provider %s: %w", providerCode, err)
	}

	// Double-check availability status and provide user-friendly messages
	if instance.AvailabilityStatus == models.WalletInstancePaused || instance.IsPaused {
		reason := "Service temporairement indisponible"
		if instance.PauseReason != nil && *instance.PauseReason != "" {
			reason = *instance.PauseReason
		}
		return nil, nil, fmt.Errorf("instance_paused: %s", reason)
	}

	if instance.AvailabilityStatus != models.WalletAvailable && instance.AvailabilityStatus != "" {
		return nil, nil, fmt.Errorf("provider %s instance not available: %s", providerCode, instance.AvailabilityStatus)
	}

	// Load provider from this instance
	provider, err := l.LoadProviderFromInstance(ctx, instance)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load provider from instance: %w", err)
	}

	return provider, instance, nil
}

// RecordTransactionUsage records usage for an instance after a transaction
func (l *InstanceBasedProviderLoader) RecordTransactionUsage(
	ctx context.Context,
	instanceID string,
	transactionID string,
	amount float64,
	currency string,
	status string,
	providerRef string,
) error {
	// Increment usage statistics
	if err := l.instanceRepo.IncrementUsage(ctx, instanceID, amount); err != nil {
		return fmt.Errorf("increment usage: %w", err)
	}

	// Log the transaction
	if err := l.instanceRepo.LogTransaction(ctx, instanceID, transactionID, "deposit", amount, currency, status, providerRef); err != nil {
		return fmt.Errorf("log transaction: %w", err)
	}

	return nil
}

// truncate truncates a string for logging purposes
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
