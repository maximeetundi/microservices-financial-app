package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	if ttl == 0 {
		ttl = 30 * time.Minute // Increased to 30 mins to reduce load
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
	// 1. Check Cache
	c.cacheMu.RLock()
	entry, found := c.cache[country]
	c.cacheMu.RUnlock()

	if found && time.Now().Before(entry.ExpiresAt) {
		log.Printf("[AdminClient] ⚡ Cache HIT for country: %s", country)
		return entry.Data, nil
	}

	// 2. Cache Miss - Fetch from Admin Service
	log.Printf("[AdminClient] 🔄 Cache MISS for country: %s. Fetching from Admin Service...", country)

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
	c.cacheMu.Lock()
	c.cache[country] = cacheEntry{
		Data:      aggregators,
		ExpiresAt: time.Now().Add(c.cacheTTL),
	}
	c.cacheMu.Unlock()

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
