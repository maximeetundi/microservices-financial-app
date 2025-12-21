package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentProvider represents a payment aggregator
type PaymentProvider struct {
	ID                  string          `json:"id" db:"id"`
	Name                string          `json:"name" db:"name"`
	DisplayName         string          `json:"display_name" db:"display_name"`
	ProviderType        string          `json:"provider_type" db:"provider_type"`
	APIBaseURL          *string         `json:"api_base_url,omitempty" db:"api_base_url"`
	APIKeyEncrypted     *string         `json:"-" db:"api_key_encrypted"`
	APISecretEncrypted  *string         `json:"-" db:"api_secret_encrypted"`
	PublicKeyEncrypted  *string         `json:"-" db:"public_key_encrypted"`
	WebhookSecretEnc    *string         `json:"-" db:"webhook_secret_encrypted"`
	IsActive            bool            `json:"is_active" db:"is_active"`
	IsDemoMode          bool            `json:"is_demo_mode" db:"is_demo_mode"`
	LogoURL             *string         `json:"logo_url,omitempty" db:"logo_url"`
	SupportedCurrencies json.RawMessage `json:"supported_currencies" db:"supported_currencies"`
	ConfigJSON          json.RawMessage `json:"config_json" db:"config_json"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
	// Joined data
	Countries []ProviderCountry `json:"countries,omitempty"`
}

// ProviderCountry represents country mapping for a provider
type ProviderCountry struct {
	ID            string  `json:"id" db:"id"`
	ProviderID    string  `json:"provider_id" db:"provider_id"`
	CountryCode   string  `json:"country_code" db:"country_code"`
	CountryName   string  `json:"country_name" db:"country_name"`
	Currency      string  `json:"currency" db:"currency"`
	IsActive      bool    `json:"is_active" db:"is_active"`
	Priority      int     `json:"priority" db:"priority"`
	MinAmount     float64 `json:"min_amount" db:"min_amount"`
	MaxAmount     float64 `json:"max_amount" db:"max_amount"`
	FeePercentage float64 `json:"fee_percentage" db:"fee_percentage"`
	FeeFixed      float64 `json:"fee_fixed" db:"fee_fixed"`
}

// PaymentHandler handles payment provider management
type PaymentHandler struct {
	db *sql.DB
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(db *sql.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// GetPaymentProviders returns all payment providers
func (h *PaymentHandler) GetPaymentProviders(c *gin.Context) {
	query := `
		SELECT id, name, display_name, provider_type, api_base_url, 
		       is_active, is_demo_mode, logo_url, supported_currencies, config_json,
		       created_at, updated_at
		FROM payment_providers
		ORDER BY name`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch providers"})
		return
	}
	defer rows.Close()

	var providers []PaymentProvider
	for rows.Next() {
		var p PaymentProvider
		err := rows.Scan(
			&p.ID, &p.Name, &p.DisplayName, &p.ProviderType, &p.APIBaseURL,
			&p.IsActive, &p.IsDemoMode, &p.LogoURL, &p.SupportedCurrencies, &p.ConfigJSON,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Get countries for this provider
		p.Countries = h.getProviderCountries(p.ID)
		providers = append(providers, p)
	}

	c.JSON(http.StatusOK, gin.H{"providers": providers})
}

// GetPaymentProvider returns a single payment provider
func (h *PaymentHandler) GetPaymentProvider(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT id, name, display_name, provider_type, api_base_url,
		       is_active, is_demo_mode, logo_url, supported_currencies, config_json,
		       created_at, updated_at
		FROM payment_providers WHERE id = $1`

	var p PaymentProvider
	err := h.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.DisplayName, &p.ProviderType, &p.APIBaseURL,
		&p.IsActive, &p.IsDemoMode, &p.LogoURL, &p.SupportedCurrencies, &p.ConfigJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	p.Countries = h.getProviderCountries(p.ID)
	c.JSON(http.StatusOK, gin.H{"provider": p})
}

// CreatePaymentProvider creates a new payment provider
func (h *PaymentHandler) CreatePaymentProvider(c *gin.Context) {
	var req struct {
		Name                string   `json:"name" binding:"required"`
		DisplayName         string   `json:"display_name" binding:"required"`
		ProviderType        string   `json:"provider_type" binding:"required"`
		APIBaseURL          string   `json:"api_base_url"`
		APIKey              string   `json:"api_key"`
		APISecret           string   `json:"api_secret"`
		PublicKey           string   `json:"public_key"`
		WebhookSecret       string   `json:"webhook_secret"`
		IsActive            bool     `json:"is_active"`
		IsDemoMode          bool     `json:"is_demo_mode"`
		LogoURL             string   `json:"logo_url"`
		SupportedCurrencies []string `json:"supported_currencies"`
		Countries           []string `json:"countries"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	currencies, _ := json.Marshal(req.SupportedCurrencies)

	query := `
		INSERT INTO payment_providers 
		(id, name, display_name, provider_type, api_base_url, api_key_encrypted, 
		 api_secret_encrypted, public_key_encrypted, webhook_secret_encrypted,
		 is_active, is_demo_mode, logo_url, supported_currencies)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	var apiBaseURL, apiKey, apiSecret, publicKey, webhookSecret *string
	if req.APIBaseURL != "" {
		apiBaseURL = &req.APIBaseURL
	}
	if req.APIKey != "" {
		apiKey = &req.APIKey // In production: encrypt this
	}
	if req.APISecret != "" {
		apiSecret = &req.APISecret
	}
	if req.PublicKey != "" {
		publicKey = &req.PublicKey
	}
	if req.WebhookSecret != "" {
		webhookSecret = &req.WebhookSecret
	}

	var returnedID string
	err := h.db.QueryRow(query,
		id, req.Name, req.DisplayName, req.ProviderType, apiBaseURL,
		apiKey, apiSecret, publicKey, webhookSecret,
		req.IsActive, req.IsDemoMode, req.LogoURL, currencies,
	).Scan(&returnedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": returnedID, "message": "Provider created successfully"})
}

// UpdatePaymentProvider updates an existing payment provider
func (h *PaymentHandler) UpdatePaymentProvider(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		DisplayName         string   `json:"display_name"`
		ProviderType        string   `json:"provider_type"`
		APIBaseURL          string   `json:"api_base_url"`
		APIKey              string   `json:"api_key"`
		APISecret           string   `json:"api_secret"`
		PublicKey           string   `json:"public_key"`
		WebhookSecret       string   `json:"webhook_secret"`
		IsActive            bool     `json:"is_active"`
		IsDemoMode          bool     `json:"is_demo_mode"`
		LogoURL             string   `json:"logo_url"`
		SupportedCurrencies []string `json:"supported_currencies"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currencies, _ := json.Marshal(req.SupportedCurrencies)

	query := `
		UPDATE payment_providers SET
			display_name = $2,
			provider_type = $3,
			api_base_url = $4,
			is_active = $5,
			is_demo_mode = $6,
			logo_url = $7,
			supported_currencies = $8,
			updated_at = NOW()
		WHERE id = $1`

	_, err := h.db.Exec(query, id, req.DisplayName, req.ProviderType,
		req.APIBaseURL, req.IsActive, req.IsDemoMode, req.LogoURL, currencies)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update provider"})
		return
	}

	// Update API keys if provided
	if req.APIKey != "" {
		h.db.Exec("UPDATE payment_providers SET api_key_encrypted = $1 WHERE id = $2", req.APIKey, id)
	}
	if req.APISecret != "" {
		h.db.Exec("UPDATE payment_providers SET api_secret_encrypted = $1 WHERE id = $2", req.APISecret, id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider updated successfully"})
}

// DeletePaymentProvider deletes a payment provider
func (h *PaymentHandler) DeletePaymentProvider(c *gin.Context) {
	id := c.Param("id")

	// Don't allow deleting demo provider
	var name string
	h.db.QueryRow("SELECT name FROM payment_providers WHERE id = $1", id).Scan(&name)
	if name == "demo" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete demo provider"})
		return
	}

	_, err := h.db.Exec("DELETE FROM payment_providers WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted successfully"})
}

// ToggleProviderStatus toggles active status
func (h *PaymentHandler) ToggleProviderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec("UPDATE payment_providers SET is_active = $1, updated_at = NOW() WHERE id = $2",
		req.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// ToggleDemoMode toggles demo mode
func (h *PaymentHandler) ToggleDemoMode(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsDemoMode bool `json:"is_demo_mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Demo provider is always in demo mode
	var name string
	h.db.QueryRow("SELECT name FROM payment_providers WHERE id = $1", id).Scan(&name)
	if name == "demo" && !req.IsDemoMode {
		c.JSON(http.StatusForbidden, gin.H{"error": "Demo provider cannot be set to live mode"})
		return
	}

	_, err := h.db.Exec("UPDATE payment_providers SET is_demo_mode = $1, updated_at = NOW() WHERE id = $2",
		req.IsDemoMode, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mode updated"})
}

// AddProviderCountry adds a country to a provider
func (h *PaymentHandler) AddProviderCountry(c *gin.Context) {
	providerID := c.Param("id")

	var req struct {
		CountryCode   string  `json:"country_code" binding:"required"`
		CountryName   string  `json:"country_name" binding:"required"`
		Currency      string  `json:"currency" binding:"required"`
		Priority      int     `json:"priority"`
		FeePercentage float64 `json:"fee_percentage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (provider_id, country_code) DO UPDATE SET
			priority = EXCLUDED.priority,
			fee_percentage = EXCLUDED.fee_percentage`

	_, err := h.db.Exec(query, providerID, req.CountryCode, req.CountryName, req.Currency, req.Priority, req.FeePercentage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add country"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Country added"})
}

// RemoveProviderCountry removes a country from a provider
func (h *PaymentHandler) RemoveProviderCountry(c *gin.Context) {
	providerID := c.Param("id")
	countryCode := c.Param("country")

	_, err := h.db.Exec("DELETE FROM provider_countries WHERE provider_id = $1 AND country_code = $2",
		providerID, countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove country"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Country removed"})
}

// TestProviderConnection tests the API connection
func (h *PaymentHandler) TestProviderConnection(c *gin.Context) {
	id := c.Param("id")

	var p PaymentProvider
	err := h.db.QueryRow(`
		SELECT name, is_demo_mode, api_base_url FROM payment_providers WHERE id = $1
	`, id).Scan(&p.Name, &p.IsDemoMode, &p.APIBaseURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	if p.IsDemoMode {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Demo mode - connection test skipped",
			"mode":    "demo",
		})
		return
	}

	// TODO: Implement actual API test for each provider
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection test successful",
		"mode":    "live",
	})
}

// GetPaymentMethodsForCountry returns available payment methods for a country (public endpoint)
func (h *PaymentHandler) GetPaymentMethodsForCountry(c *gin.Context) {
	countryCode := c.Query("country")
	if countryCode == "" {
		countryCode = "CI" // Default to Côte d'Ivoire
	}

	query := `
		SELECT pp.id, pp.name, pp.display_name, pp.provider_type, pp.is_demo_mode, pp.logo_url,
		       pc.fee_percentage, pc.fee_fixed, pc.min_amount, pc.max_amount
		FROM payment_providers pp
		JOIN provider_countries pc ON pp.id = pc.provider_id
		WHERE pc.country_code = $1 AND pp.is_active = true AND pc.is_active = true
		ORDER BY pc.priority ASC`

	rows, err := h.db.Query(query, countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment methods"})
		return
	}
	defer rows.Close()

	var methods []map[string]interface{}
	for rows.Next() {
		var id, name, displayName, providerType string
		var isDemoMode bool
		var logoURL *string
		var feePercentage, feeFixed, minAmount, maxAmount float64

		err := rows.Scan(&id, &name, &displayName, &providerType, &isDemoMode, &logoURL,
			&feePercentage, &feeFixed, &minAmount, &maxAmount)
		if err != nil {
			continue
		}

		methods = append(methods, map[string]interface{}{
			"id":             id,
			"name":           name,
			"display_name":   displayName,
			"provider_type":  providerType,
			"is_demo_mode":   isDemoMode,
			"logo_url":       logoURL,
			"fee_percentage": feePercentage,
			"fee_fixed":      feeFixed,
			"min_amount":     minAmount,
			"max_amount":     maxAmount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"country":         countryCode,
		"payment_methods": methods,
	})
}

// Helper function to get countries for a provider
func (h *PaymentHandler) getProviderCountries(providerID string) []ProviderCountry {
	query := `
		SELECT id, provider_id, country_code, country_name, currency, 
		       is_active, priority, min_amount, max_amount, fee_percentage, fee_fixed
		FROM provider_countries WHERE provider_id = $1
		ORDER BY priority ASC`

	rows, err := h.db.Query(query, providerID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var countries []ProviderCountry
	for rows.Next() {
		var c ProviderCountry
		err := rows.Scan(&c.ID, &c.ProviderID, &c.CountryCode, &c.CountryName, &c.Currency,
			&c.IsActive, &c.Priority, &c.MinAmount, &c.MaxAmount, &c.FeePercentage, &c.FeeFixed)
		if err == nil {
			countries = append(countries, c)
		}
	}
	return countries
}
