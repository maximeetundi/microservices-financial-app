package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
)

// AdminClient handles communication with the Admin Service
type AdminClient struct {
	baseURL    string
	httpClient *http.Client
	cacheTTL   time.Duration

	// In-memory cache
	cacheMu sync.RWMutex
	cache   map[string]cacheEntry
}

type cacheEntry struct {
	Data      []models.AggregatorForFrontend
	ExpiresAt time.Time
}

// NewAdminClient creates a new AdminClient
func NewAdminClient(cfg *config.Config) *AdminClient {
	ttl := time.Duration(cfg.AggregatorCacheTTL) * time.Second
	disableCache := strings.TrimSpace(strings.ToLower(os.Getenv("DISABLE_ADMIN_CACHE")))
	if disableCache == "1" || disableCache == "true" || disableCache == "yes" {
		ttl = 0
	}
	// If ttl is 0, caching is disabled (useful during development)
	if ttl < 0 {
		ttl = 0
	}

	return &AdminClient{
		baseURL: cfg.AdminServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cacheTTL: ttl,
		cache:    make(map[string]cacheEntry),
	}
}

// PaymentMethodResponse matches the JSON structure from Admin Service
type PaymentMethodResponse struct {
	Countries      []string                 `json:"countries"`
	PaymentMethods []map[string]interface{} `json:"payment_methods"`
}

// GetPaymentMethods fetches payment methods from admin-service with caching
func (c *AdminClient) GetPaymentMethods(country string) ([]models.AggregatorForFrontend, error) {
	// 1. Check Cache (unless disabled)
	if c.cacheTTL > 0 {
		c.cacheMu.RLock()
		entry, found := c.cache[country]
		c.cacheMu.RUnlock()

		if found && time.Now().Before(entry.ExpiresAt) {
			log.Printf("[AdminClient] ⚡ Cache HIT for country: %s", country)
			return entry.Data, nil
		}
		log.Printf("[AdminClient] 🔄 Cache MISS for country: %s. Fetching from Admin Service...", country)
	} else {
		log.Printf("[AdminClient] ♻️ Cache disabled. Fetching payment methods from Admin Service for country: %s", country)
	}

	url := fmt.Sprintf("%s/api/v1/admin/payment-methods?country=%s", c.baseURL, country)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call admin service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("admin service returned status: %d", resp.StatusCode)
	}

	// 3. Decode Response
	var result PaymentMethodResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 4. Map to Transfer Service Model
	aggregators := c.mapToAggregators(result.PaymentMethods)

	// 5. Update Cache
	if c.cacheTTL > 0 {
		c.cacheMu.Lock()
		c.cache[country] = cacheEntry{
			Data:      aggregators,
			ExpiresAt: time.Now().Add(c.cacheTTL),
		}
		c.cacheMu.Unlock()
	}

	return aggregators, nil
}

// mapToAggregators converts Admin Service loose JSON to Transfer Service strict structs
func (c *AdminClient) mapToAggregators(rawMethods []map[string]interface{}) []models.AggregatorForFrontend {
	var aggregators []models.AggregatorForFrontend

	for _, m := range rawMethods {
		// Helper to safely get string
		getString := func(key string) string {
			if v, ok := m[key].(string); ok {
				return v
			}
			return ""
		}

		// Helper to safely get float
		getFloat := func(key string) float64 {
			if v, ok := m[key].(float64); ok {
				return v
			}
			return 0
		}

		// Helper to safely get bool
		getBool := func(key string) bool {
			if v, ok := m[key].(bool); ok {
				return v
			}
			return false
		}

		agg := models.AggregatorForFrontend{
			Code:            getString("name"),         // Mapping 'name' to 'Code' (e.g. mtn_momo)
			Name:            getString("display_name"), // Mapping 'display_name' to 'Name'
			LogoURL:         getString("logo_url"),
			DepositEnabled:  getBool("deposit_enabled"),
			WithdrawEnabled: getBool("withdraw_enabled"),
			MinAmount:       getFloat("min_amount"),
			MaxAmount:       getFloat("max_amount"),
			FeePercent:      getFloat("fee_percentage"),
			FeeFixed:        getFloat("fee_fixed"),
			FeeCurrency:     "XOF", // Default or extract from somewhere else if needed
			MaintenanceMode: false, // Admin service doesn't send this explicitly in same format yet
		}

		// Map 'provider_type' or other fields if needed
		// For now this covers the frontend needs

		aggregators = append(aggregators, agg)
	}

	return aggregators
}

// ==================== INTERNAL API FOR CREDENTIALS ====================

// InstanceResponse represents the response from admin-service internal API
type InstanceResponse struct {
	Instance InstanceData `json:"instance"`
}

// InstanceData contains instance details with credentials
type InstanceData struct {
	ID             string            `json:"id"`
	InstanceName   string            `json:"instance_name"`
	ProviderCode   string            `json:"provider_code"`
	ProviderName   string            `json:"provider_name"`
	IsActive       bool              `json:"is_active"`
	IsPrimary      bool              `json:"is_primary"`
	IsGlobal       bool              `json:"is_global"`
	IsPaused       bool              `json:"is_paused"`
	PauseReason    *string           `json:"pause_reason"`
	Priority       int               `json:"priority"`
	HealthStatus   string            `json:"health_status"`
	IsTestMode     bool              `json:"is_test_mode"`
	DepositEnabled bool              `json:"deposit_enabled"`
	WithdrawEnabled bool             `json:"withdraw_enabled"`
	APICredentials map[string]string `json:"api_credentials"`
}

// GetBestInstanceWithCredentials fetches the best instance with full API credentials
// from admin-service internal API (service-to-service communication)
func (c *AdminClient) GetBestInstanceWithCredentials(providerCode, country string, amount float64, currency string, operation string) (*models.AggregatorInstanceWithDetails, error) {
	log.Printf("[AdminClient] 🔐 Fetching instance with credentials: provider=%s, country=%s", providerCode, country)
	if operation == "" {
		operation = "deposit"
	}

	// Build request
	reqBody := map[string]interface{}{
		"provider_code": providerCode,
		"country":       country,
		"amount":        amount,
		"currency":      currency,
		"operation":     operation,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call admin-service internal API
	url := fmt.Sprintf("%s/api/v1/internal/instances/best", c.baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[AdminClient] ❌ Failed to call admin-service: %v", err)
		return nil, fmt.Errorf("failed to call admin service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		log.Printf("[AdminClient] ❌ Admin service returned status %d: %v", resp.StatusCode, errResp)
		return nil, fmt.Errorf("admin service error: %v", errResp["error"])
	}

	// Decode response
	var result InstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	inst := result.Instance
	log.Printf("[AdminClient] ✅ Got instance %s (%s) with %d credentials",
		inst.ID, inst.InstanceName, len(inst.APICredentials))

	// Map to transfer-service model
	instance := &models.AggregatorInstanceWithDetails{
		ID:             inst.ID,
		InstanceName:   inst.InstanceName,
		ProviderCode:   inst.ProviderCode,
		ProviderName:   inst.ProviderName,
		Enabled:        inst.IsActive,
		IsPaused:       inst.IsPaused,
		IsGlobal:       inst.IsGlobal,
		Priority:       inst.Priority,
		IsTestMode:     inst.IsTestMode,
		PauseReason:    inst.PauseReason,
		APICredentials: inst.APICredentials,
		DepositEnabled: inst.DepositEnabled,
		WithdrawEnabled: inst.WithdrawEnabled,
	}

	return instance, nil
}

// GetInstanceByID fetches a specific instance with credentials by ID
func (c *AdminClient) GetInstanceByID(instanceID string) (*models.AggregatorInstanceWithDetails, error) {
	log.Printf("[AdminClient] 🔐 Fetching instance by ID: %s", instanceID)

	url := fmt.Sprintf("%s/api/v1/internal/instances/%s", c.baseURL, instanceID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call admin service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instance not found: %s", instanceID)
	}

	var result InstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	inst := result.Instance

	return &models.AggregatorInstanceWithDetails{
		ID:             inst.ID,
		InstanceName:   inst.InstanceName,
		ProviderCode:   inst.ProviderCode,
		ProviderName:   inst.ProviderName,
		Enabled:        inst.IsActive,
		IsPaused:       inst.IsPaused,
		IsGlobal:       inst.IsGlobal,
		Priority:       inst.Priority,
		IsTestMode:     inst.IsTestMode,
		PauseReason:    inst.PauseReason,
		APICredentials: inst.APICredentials,
	}, nil
}
