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
		return nil, fmt.Errorf("vault client not initialized")
	}

	// Use custom vault path if provided, otherwise use default
	path := vaultPath
	if path == "" {
		path = fmt.Sprintf("secret/aggregators/%s", providerCode)
	}

	// Convert path for KV v2: secret/aggregators/x -> secret/data/aggregators/x
	kvV2Path := convertToKVv2Path(path)
	log.Printf("[InstanceLoader] Loading credentials from Vault path: %s (KV v2: %s)", path, kvV2Path)

	// Try KV v2 path first
	data, err := l.vaultClient.GetSecret(kvV2Path)
	if err != nil {
		// Fall back to KV v1 path
		log.Printf("[InstanceLoader] KV v2 path failed, trying KV v1 path: %s", path)
		data, err = l.vaultClient.GetSecret(path)
		if err != nil {
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

	log.Printf("[InstanceLoader] Loaded %d credentials from Vault for %s", len(creds), providerCode)
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
	// Always try Vault first if we have a vault path (this is the recommended approach)
	if vaultPath != "" && l.vaultClient != nil {
		creds, err := l.getCredentialsFromVault(instance.ProviderCode, vaultPath)
		if err == nil && len(creds) > 0 {
			log.Printf("[InstanceLoader] Loaded credentials from Vault for %s (path: %s)", instance.ProviderCode, vaultPath)
			return creds, nil
		}
		log.Printf("[InstanceLoader] Failed to load from Vault for %s: %v, checking database fallback", instance.ProviderCode, err)
	}

	// Fall back to database credentials if Vault fails
	if len(instance.APICredentials) > 0 {
		// Verify credentials are not placeholders
		hasRealCreds := false
		for _, v := range instance.APICredentials {
			if v != "" && !isPlaceholder(v) {
				hasRealCreds = true
				break
			}
		}
		if hasRealCreds {
			log.Printf("[InstanceLoader] Using credentials from database for %s", instance.ProviderCode)
			return instance.APICredentials, nil
		}
	}

	return nil, fmt.Errorf("no valid credentials found for %s (vault path: %s)", instance.ProviderCode, vaultPath)
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

	// Use VaultSecretPath from instance if available, otherwise fall back to default
	vaultPath := instance.VaultSecretPath
	if vaultPath == "" {
		// Check environment variable as fallback
		vaultPath = os.Getenv(fmt.Sprintf("%s_VAULT_PATH", instance.ProviderCode))
	}
	if vaultPath == "" {
		// Final fallback to default path
		vaultPath = fmt.Sprintf("secret/aggregators/%s", instance.ProviderCode)
	}

	log.Printf("[InstanceLoader] Instance %s (%s): using vault path: %s", instance.ID, instance.ProviderCode, vaultPath)

	switch instance.ProviderCode {
	case "flutterwave":
		return l.loadFlutterwaveFromInstance(instance, vaultPath)
	case "cinetpay":
		return l.loadCinetPayFromInstance(instance, vaultPath)
	case "paystack":
		return l.loadPaystackFromInstance(instance, vaultPath)
	case "stripe":
		return l.loadStripeFromInstance(instance, vaultPath)
	case "wave":
		return l.loadWaveFromInstance(instance, vaultPath)
	case "paypal":
		return l.loadPayPalFromInstance(instance, vaultPath)
	case "orange_money":
		return l.loadOrangeMoneyFromInstance(instance, vaultPath)
	case "mtn_momo":
		return l.loadMTNMoMoFromInstance(instance, vaultPath)
	case "lygos":
		return l.loadLygosFromInstance(instance, vaultPath)
	case "yellowcard":
		return l.loadYellowCardFromInstance(instance, vaultPath)
	case "moov_money":
		return l.loadMoovMoneyFromInstance(instance, vaultPath)
	case "fedapay":
		return l.loadFedaPayFromInstance(instance, vaultPath)
	case "demo":
		return l.loadDemoFromInstance(instance)
	default:
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
	creds, err := l.getCredentials(instance, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get CinetPay credentials: %w", err)
	}

	config := CinetPayConfig{
		APIKey:  creds["api_key"],
		SiteID:  creds["site_id"],
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
	creds, err := l.getCredentials(instance, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get Wave credentials: %w", err)
	}

	config := WaveConfig{
		APIKey:  creds["api_key"],
		BaseURL: "https://api.wave.com/v1",
	}

	if baseURL := creds["base_url"]; baseURL != "" {
		config.BaseURL = baseURL
	}

	return NewWaveCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPayPalFromInstance(instance *models.AggregatorInstanceWithDetails, vaultPath string) (CollectionProvider, error) {
	creds, err := l.getCredentials(instance, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get PayPal credentials: %w", err)
	}

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
