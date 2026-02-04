package providers

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

// InstanceBasedProviderLoader loads providers from database instances
// instead of Vault, respecting admin panel configuration
type InstanceBasedProviderLoader struct {
	instanceRepo *repository.AggregatorInstanceRepository
}

// NewInstanceBasedProviderLoader creates a new instance-based loader
func NewInstanceBasedProviderLoader(instanceRepo *repository.AggregatorInstanceRepository) *InstanceBasedProviderLoader {
	return &InstanceBasedProviderLoader{
		instanceRepo: instanceRepo,
	}
}

// LoadProviderFromInstance creates a provider instance from aggregator instance config
func (l *InstanceBasedProviderLoader) LoadProviderFromInstance(
	ctx context.Context,
	instance *models.AggregatorInstanceWithDetails,
) (CollectionProvider, error) {

	switch instance.ProviderCode {
	case "flutterwave":
		return l.loadFlutterwaveFromInstance(instance)
	case "cinetpay":
		return l.loadCinetPayFromInstance(instance)
	case "paystack":
		return l.loadPaystackFromInstance(instance)
	case "stripe":
		return l.loadStripeFromInstance(instance)
	case "wave":
		return l.loadWaveFromInstance(instance)
	case "paypal":
		return l.loadPayPalFromInstance(instance)
	case "orange_money":
		return l.loadOrangeMoneyFromInstance(instance)
	case "mtn_momo":
		return l.loadMTNMoMoFromInstance(instance)
	default:
		return nil, fmt.Errorf("unknown provider: %s", instance.ProviderCode)
	}
}

func (l *InstanceBasedProviderLoader) loadFlutterwaveFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := FlutterwaveConfig{
		PublicKey:     instance.APICredentials["public_key"],
		SecretKey:     instance.APICredentials["secret_key"],
		EncryptionKey: instance.APICredentials["encryption_key"],
		WebhookSecret: instance.APICredentials["webhook_secret"],
		BaseURL:       "https://api.flutterwave.com/v3",
	}

	if instance.IsTestMode {
		// Use test URL if in test mode
		config.BaseURL = "https://api.flutterwave.com/v3"
	}

	return NewFlutterwaveCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadCinetPayFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := CinetPayConfig{
		APIKey:  instance.APICredentials["api_key"],
		SiteID:  instance.APICredentials["site_id"],
		BaseURL: "https://api-checkpoint.cinetpay.com/v2",
	}

	return NewCinetPayCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPaystackFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := PaystackConfig{
		PublicKey: instance.APICredentials["public_key"],
		SecretKey: instance.APICredentials["secret_key"],
		BaseURL:   "https://api.paystack.co",
	}

	return NewPaystackCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadStripeFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := StripeConfig{
		PublicKey:     instance.APICredentials["public_key"],
		SecretKey:     instance.APICredentials["secret_key"],
		WebhookSecret: instance.APICredentials["webhook_secret"],
		BaseURL:       "https://api.stripe.com/v1",
	}

	return NewStripeCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadWaveFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := WaveConfig{
		APIKey:  instance.APICredentials["api_key"],
		BaseURL: "https://api.wave.com/v1",
	}

	return NewWaveCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadPayPalFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	baseURL := "https://api-m.paypal.com"
	if instance.IsTestMode {
		baseURL = "https://api-m.sandbox.paypal.com"
	}

	config := PayPalConfig{
		ClientID:     instance.APICredentials["client_id"],
		ClientSecret: instance.APICredentials["client_secret"],
		Mode:         "sandbox",
		BaseURL:      baseURL,
	}

	if !instance.IsTestMode {
		config.Mode = "live"
	}

	return NewPayPalCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadOrangeMoneyFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	config := OrangeMoneyConfig{
		ClientID:     instance.APICredentials["client_id"],
		ClientSecret: instance.APICredentials["client_secret"],
		MerchantKey:  instance.APICredentials["merchant_key"],
		BaseURL:      "https://api.orange.com/orange-money-webpay",
	}

	return NewOrangeMoneyCollectionProvider(config), nil
}

func (l *InstanceBasedProviderLoader) loadMTNMoMoFromInstance(instance *models.AggregatorInstanceWithDetails) (CollectionProvider, error) {
	baseURL := "https://sandbox.momodeveloper.mtn.com"
	environment := "sandbox"

	if !instance.IsTestMode {
		baseURL = "https://proxy.momoapi.mtn.com"
		environment = "prod"
	}

	config := MTNMoMoConfig{
		APIUser:         instance.APICredentials["api_user"],
		APIKey:          instance.APICredentials["api_key"],
		SubscriptionKey: instance.APICredentials["subscription_key"],
		BaseURL:         baseURL,
		Environment:     environment,
	}

	return NewMTNMoMoCollectionProvider(config), nil
}

// GetBestProviderForDeposit selects the best instance for a deposit
// based on provider code, country, amount, and instance availability
func (l *InstanceBasedProviderLoader) GetBestProviderForDeposit(
	ctx context.Context,
	providerCode string,
	country string,
	amount float64,
) (CollectionProvider, *models.AggregatorInstanceWithDetails, error) {

	// Get best instance from DB
	instance, err := l.instanceRepo.GetBestInstanceForProvider(ctx, providerCode, country, amount)
	if err != nil {
		return nil, nil, fmt.Errorf("no available instance for provider %s: %w", providerCode, err)
	}

	if !instance.IsAvailable() {
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
