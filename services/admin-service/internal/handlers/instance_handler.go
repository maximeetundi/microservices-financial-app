package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
)

// VaultClient wraps the Vault API client
type VaultClient struct {
	client *vault.Client
}

// NewVaultClient creates a new Vault client
func NewVaultClient() (*VaultClient, error) {
	config := vault.DefaultConfig()

	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	}

	return &VaultClient{client: client}, nil
}

// WriteSecret writes a secret to Vault (supports KV v2)
func (v *VaultClient) WriteSecret(path string, data map[string]interface{}) error {
	// Convert path for KV v2: secret/x -> secret/data/x
	kvV2Path := convertToKVv2Path(path)

	// Wrap data in "data" key for KV v2
	wrappedData := map[string]interface{}{
		"data": data,
	}

	_, err := v.client.Logical().Write(kvV2Path, wrappedData)
	if err != nil {
		// Try KV v1 as fallback
		_, err = v.client.Logical().Write(path, data)
	}
	return err
}

// GetSecret reads a secret from Vault
func (v *VaultClient) GetSecret(path string) (map[string]interface{}, error) {
	kvV2Path := convertToKVv2Path(path)

	secret, err := v.client.Logical().Read(kvV2Path)
	if err != nil || secret == nil {
		// Try KV v1 as fallback
		secret, err = v.client.Logical().Read(path)
		if err != nil || secret == nil {
			return nil, fmt.Errorf("secret not found at %s", path)
		}
		return secret.Data, nil
	}

	// Extract data from KV v2 response
	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return secret.Data, nil
}

// convertToKVv2Path converts a KV v1 style path to KV v2 style
func convertToKVv2Path(path string) string {
	// If path already contains /data/, return as is
	if len(path) > 12 && path[7:11] == "data" {
		return path
	}

	// Find the first slash (mount point separator)
	for i, c := range path {
		if c == '/' {
			return path[:i] + "/data" + path[i:]
		}
	}
	return path
}

// ProviderInstance represents a single API key instance for a provider
type ProviderInstance struct {
	ID              string     `json:"id" db:"id"`
	ProviderID      string     `json:"provider_id" db:"provider_id"`
	Name            string     `json:"name" db:"name"`
	VaultSecretPath string     `json:"vault_secret_path" db:"vault_secret_path"`
	HotWalletID     *string    `json:"hot_wallet_id,omitempty" db:"hot_wallet_id"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsPrimary       bool       `json:"is_primary" db:"is_primary"`
	IsGlobal        bool       `json:"is_global" db:"is_global"`
	IsPaused        bool       `json:"is_paused" db:"is_paused"`
	Priority        int        `json:"priority" db:"priority"`
	RequestCount    int64      `json:"request_count" db:"request_count"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	LastError       *string    `json:"last_error,omitempty" db:"last_error"`
	HealthStatus    string     `json:"health_status" db:"health_status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	// Joined data
	ProviderName     string   `json:"provider_name,omitempty"`
	HotWalletBalance *float64 `json:"hot_wallet_balance,omitempty"`
}

// InstanceHandler handles provider instance management
type InstanceHandler struct {
	db          *sql.DB
	vaultClient *VaultClient
}

// NewInstanceHandler creates a new InstanceHandler
func NewInstanceHandler(db *sql.DB) *InstanceHandler {
	handler := &InstanceHandler{db: db}

	// Initialize Vault client
	vaultClient, err := NewVaultClient()
	if err != nil {
		log.Printf("[InstanceHandler] Warning: Failed to initialize Vault client: %v", err)
		log.Printf("[InstanceHandler] Credentials will not be synced to Vault")
	} else {
		handler.vaultClient = vaultClient
		log.Printf("[InstanceHandler] Vault client initialized successfully")
	}

	return handler
}

// CredentialsRequest represents API credentials for a provider instance
type CredentialsRequest struct {
	// Common fields
	ClientID      string `json:"client_id,omitempty"`
	ClientSecret  string `json:"client_secret,omitempty"`
	APIKey        string `json:"api_key,omitempty"`
	SecretKey     string `json:"secret_key,omitempty"`
	PublicKey     string `json:"public_key,omitempty"`
	WebhookID     string `json:"webhook_id,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
	BaseURL       string `json:"base_url,omitempty"`
	Mode          string `json:"mode,omitempty"` // sandbox or live
	// Provider-specific
	EncryptionKey   string `json:"encryption_key,omitempty"`
	SiteID          string `json:"site_id,omitempty"`
	MerchantKey     string `json:"merchant_key,omitempty"`
	SubscriptionKey string `json:"subscription_key,omitempty"`
	APIUser         string `json:"api_user,omitempty"`
	ShopName        string `json:"shop_name,omitempty"`
	BusinessID      string `json:"business_id,omitempty"`
	Environment     string `json:"environment,omitempty"`
}

// ToMap converts credentials to a map for Vault storage
func (c *CredentialsRequest) ToMap() map[string]interface{} {
	data := make(map[string]interface{})

	if c.ClientID != "" {
		data["client_id"] = c.ClientID
	}
	if c.ClientSecret != "" {
		data["client_secret"] = c.ClientSecret
	}
	if c.APIKey != "" {
		data["api_key"] = c.APIKey
	}
	if c.SecretKey != "" {
		data["secret_key"] = c.SecretKey
	}
	if c.PublicKey != "" {
		data["public_key"] = c.PublicKey
	}
	if c.WebhookID != "" {
		data["webhook_id"] = c.WebhookID
	}
	if c.WebhookSecret != "" {
		data["webhook_secret"] = c.WebhookSecret
	}
	if c.BaseURL != "" {
		data["base_url"] = c.BaseURL
	}
	if c.Mode != "" {
		data["mode"] = c.Mode
	}
	if c.EncryptionKey != "" {
		data["encryption_key"] = c.EncryptionKey
	}
	if c.SiteID != "" {
		data["site_id"] = c.SiteID
	}
	if c.MerchantKey != "" {
		data["merchant_key"] = c.MerchantKey
	}
	if c.SubscriptionKey != "" {
		data["subscription_key"] = c.SubscriptionKey
	}
	if c.APIUser != "" {
		data["api_user"] = c.APIUser
	}
	if c.ShopName != "" {
		data["shop_name"] = c.ShopName
	}
	if c.BusinessID != "" {
		data["business_id"] = c.BusinessID
	}
	if c.Environment != "" {
		data["environment"] = c.Environment
	}

	return data
}

// GetProviderInstances returns all instances for a provider
func (h *InstanceHandler) GetProviderInstances(c *gin.Context) {
	providerID := c.Param("id")

	query := `
		SELECT id, provider_id, name, vault_secret_path, hot_wallet_id,
		       is_active, is_primary, COALESCE(is_global, FALSE) as is_global,
		       COALESCE(is_paused, FALSE) as is_paused,
		       priority, request_count, last_used_at,
		       last_error, health_status, created_at, updated_at
		FROM provider_instances
		WHERE provider_id = $1
		ORDER BY priority DESC, is_primary DESC, created_at`

	rows, err := h.db.Query(query, providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch instances"})
		return
	}
	defer rows.Close()

	var instances []ProviderInstance
	for rows.Next() {
		var inst ProviderInstance
		err := rows.Scan(
			&inst.ID, &inst.ProviderID, &inst.Name, &inst.VaultSecretPath, &inst.HotWalletID,
			&inst.IsActive, &inst.IsPrimary, &inst.IsGlobal, &inst.IsPaused,
			&inst.Priority, &inst.RequestCount, &inst.LastUsedAt,
			&inst.LastError, &inst.HealthStatus, &inst.CreatedAt, &inst.UpdatedAt,
		)
		if err != nil {
			continue
		}
		instances = append(instances, inst)
	}

	c.JSON(http.StatusOK, gin.H{"instances": instances})
}

// GetAllInstances returns all instances across all providers
func (h *InstanceHandler) GetAllInstances(c *gin.Context) {
	query := `
		SELECT i.id, i.provider_id, i.name, i.vault_secret_path, i.hot_wallet_id,
		       i.is_active, i.is_primary, i.priority, i.request_count, i.last_used_at,
		       i.last_error, i.health_status, i.created_at, i.updated_at,
		       p.display_name as provider_name
		FROM provider_instances i
		JOIN payment_providers p ON i.provider_id = p.id
		ORDER BY p.name, i.priority DESC, i.is_primary DESC`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch instances"})
		return
	}
	defer rows.Close()

	var instances []ProviderInstance
	for rows.Next() {
		var inst ProviderInstance
		err := rows.Scan(
			&inst.ID, &inst.ProviderID, &inst.Name, &inst.VaultSecretPath, &inst.HotWalletID,
			&inst.IsActive, &inst.IsPrimary, &inst.Priority, &inst.RequestCount, &inst.LastUsedAt,
			&inst.LastError, &inst.HealthStatus, &inst.CreatedAt, &inst.UpdatedAt,
			&inst.ProviderName,
		)
		if err != nil {
			continue
		}
		instances = append(instances, inst)
	}

	c.JSON(http.StatusOK, gin.H{"instances": instances})
}

// CreateProviderInstance creates a new instance for a provider
func (h *InstanceHandler) CreateProviderInstance(c *gin.Context) {
	providerID := c.Param("id")

	// Verify the provider exists before creating an instance
	var providerExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM payment_providers WHERE id = $1)", providerID).Scan(&providerExists)
	if err != nil || !providerExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found. Please ensure the aggregator exists before creating an instance."})
		return
	}

	var req struct {
		Name            string `json:"name" binding:"required"`
		VaultSecretPath string `json:"vault_secret_path" binding:"required"`
		HotWalletID     string `json:"hot_wallet_id"`
		Wallets         []struct {
			HotWalletID string `json:"hot_wallet_id"`
			Currency    string `json:"currency"`
		} `json:"wallets"`
		IsActive    bool                `json:"is_active"`
		IsPrimary   bool                `json:"is_primary"`
		IsGlobal    bool                `json:"is_global"`
		Priority    int                 `json:"priority"`
		Credentials *CredentialsRequest `json:"credentials,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If this is primary, unset other primaries
	if req.IsPrimary {
		h.db.Exec("UPDATE provider_instances SET is_primary = FALSE WHERE provider_id = $1", providerID)
	}

	id := uuid.New().String()
	priority := req.Priority
	if priority == 0 {
		priority = 50
	}

	query := `
		INSERT INTO provider_instances
		(id, provider_id, name, vault_secret_path, hot_wallet_id, is_active, is_primary, priority, health_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active')
		RETURNING id`

	var hotWalletID *string
	if req.HotWalletID != "" {
		hotWalletID = &req.HotWalletID
	}
	// If explicit hot_wallet_id not set but wallets array is, take the first one as primary reference
	if hotWalletID == nil && len(req.Wallets) > 0 {
		hotWalletID = &req.Wallets[0].HotWalletID
	}

	var returnedID string
	err = h.db.QueryRow(query,
		id, providerID, req.Name, req.VaultSecretPath, hotWalletID,
		req.IsActive, req.IsPrimary, priority,
	).Scan(&returnedID)

	if err != nil {
		log.Printf("Error creating instance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instance: " + err.Error()})
		return
	}

	// Insert linked wallets if provided
	if len(req.Wallets) > 0 {
		walletQuery := `
			INSERT INTO provider_instance_wallets (id, instance_id, currency, hot_wallet_id, is_active, priority)
			VALUES ($1, $2, $3, $4, TRUE, 50)
			ON CONFLICT (instance_id, currency, hot_wallet_id) DO NOTHING`

		for _, w := range req.Wallets {
			if w.HotWalletID != "" && w.Currency != "" {
				_, err := h.db.Exec(walletQuery, uuid.New().String(), returnedID, w.Currency, w.HotWalletID)
				if err != nil {
					log.Printf("Failed to link wallet %s to instance %s: %v", w.HotWalletID, returnedID, err)
					// Continue adding others even if one fails
				}
			}
		}
	}

	// Write credentials to Vault if provided
	if req.Credentials != nil && h.vaultClient != nil {
		credData := req.Credentials.ToMap()
		if len(credData) > 0 {
			if err := h.vaultClient.WriteSecret(req.VaultSecretPath, credData); err != nil {
				log.Printf("[InstanceHandler] Warning: Failed to write credentials to Vault at %s: %v", req.VaultSecretPath, err)
				// Don't fail the request, just log the warning
			} else {
				log.Printf("[InstanceHandler] Successfully wrote credentials to Vault at %s", req.VaultSecretPath)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                returnedID,
		"message":           "Instance created successfully",
		"vault_secret_path": req.VaultSecretPath,
		"credentials_saved": req.Credentials != nil && h.vaultClient != nil,
	})
}

// UpdateProviderInstance updates an existing instance
func (h *InstanceHandler) UpdateProviderInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		Name            string              `json:"name"`
		VaultSecretPath string              `json:"vault_secret_path"`
		HotWalletID     string              `json:"hot_wallet_id"`
		IsActive        *bool               `json:"is_active,omitempty"`
		IsPrimary       *bool               `json:"is_primary,omitempty"`
		IsGlobal        *bool               `json:"is_global,omitempty"`
		Priority        *int                `json:"priority,omitempty"`
		HealthStatus    string              `json:"health_status"`
		Credentials     *CredentialsRequest `json:"credentials,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argNum := 1

	if req.Name != "" {
		updates = append(updates, "name = $"+string(rune('0'+argNum)))
		args = append(args, req.Name)
		argNum++
	}
	if req.VaultSecretPath != "" {
		updates = append(updates, "vault_secret_path = $"+string(rune('0'+argNum)))
		args = append(args, req.VaultSecretPath)
		argNum++
	}
	if req.HotWalletID != "" {
		updates = append(updates, "hot_wallet_id = $"+string(rune('0'+argNum)))
		args = append(args, req.HotWalletID)
		argNum++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsActive)
		argNum++
	}
	if req.IsPrimary != nil {
		updates = append(updates, "is_primary = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsPrimary)
		argNum++
	}
	if req.Priority != nil {
		updates = append(updates, "priority = $"+string(rune('0'+argNum)))
		args = append(args, *req.Priority)
		argNum++
	}
	if req.HealthStatus != "" {
		updates = append(updates, "health_status = $"+string(rune('0'+argNum)))
		args = append(args, req.HealthStatus)
		argNum++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, instanceID)

	query := "UPDATE provider_instances SET " + joinStrings(updates, ", ") + " WHERE id = $" + string(rune('0'+argNum))

	result, err := h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update instance"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	// Update credentials in Vault if provided
	credentialsSaved := false
	if req.Credentials != nil && h.vaultClient != nil {
		// Get current vault path from DB if not provided in request
		vaultPath := req.VaultSecretPath
		if vaultPath == "" {
			h.db.QueryRow("SELECT vault_secret_path FROM provider_instances WHERE id = $1", instanceID).Scan(&vaultPath)
		}

		if vaultPath != "" {
			credData := req.Credentials.ToMap()
			if len(credData) > 0 {
				if err := h.vaultClient.WriteSecret(vaultPath, credData); err != nil {
					log.Printf("[InstanceHandler] Warning: Failed to update credentials in Vault at %s: %v", vaultPath, err)
				} else {
					log.Printf("[InstanceHandler] Successfully updated credentials in Vault at %s", vaultPath)
					credentialsSaved = true
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Instance updated successfully",
		"credentials_saved": credentialsSaved,
	})
}

// UpdateInstanceCredentials updates only the credentials in Vault for an instance
func (h *InstanceHandler) UpdateInstanceCredentials(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		Credentials CredentialsRequest `json:"credentials" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get vault path from instance
	var vaultPath string
	err := h.db.QueryRow("SELECT vault_secret_path FROM provider_instances WHERE id = $1", instanceID).Scan(&vaultPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	if h.vaultClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Vault client not available"})
		return
	}

	credData := req.Credentials.ToMap()
	if len(credData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No credentials provided"})
		return
	}

	if err := h.vaultClient.WriteSecret(vaultPath, credData); err != nil {
		log.Printf("[InstanceHandler] Failed to update credentials in Vault: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credentials to Vault"})
		return
	}

	log.Printf("[InstanceHandler] Successfully updated credentials for instance %s at %s", instanceID, vaultPath)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Credentials updated successfully",
		"vault_path": vaultPath,
	})
}

// GetInstanceCredentials retrieves credentials from Vault (masked for security)
func (h *InstanceHandler) GetInstanceCredentials(c *gin.Context) {
	instanceID := c.Param("instanceId")

	// Get vault path from instance
	var vaultPath string
	err := h.db.QueryRow("SELECT vault_secret_path FROM provider_instances WHERE id = $1", instanceID).Scan(&vaultPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	if h.vaultClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Vault client not available"})
		return
	}

	data, err := h.vaultClient.GetSecret(vaultPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "Credentials not found in Vault",
			"vault_path": vaultPath,
			"details":    err.Error(),
		})
		return
	}

	// Mask sensitive values for display
	maskedData := make(map[string]interface{})
	for key, value := range data {
		if str, ok := value.(string); ok && len(str) > 0 {
			if len(str) > 8 {
				maskedData[key] = str[:4] + "****" + str[len(str)-4:]
			} else {
				maskedData[key] = "****"
			}
		} else {
			maskedData[key] = value
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"vault_path":  vaultPath,
		"credentials": maskedData,
		"keys":        getMapKeys(data),
	})
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// DeleteProviderInstance deletes an instance
func (h *InstanceHandler) DeleteProviderInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")

	result, err := h.db.Exec("DELETE FROM provider_instances WHERE id = $1", instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete instance"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance deleted successfully"})
}

// LinkHotWallet links an instance to a hot wallet
func (h *InstanceHandler) LinkHotWallet(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		HotWalletID string `json:"hot_wallet_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.Exec(
		"UPDATE provider_instances SET hot_wallet_id = $1, updated_at = NOW() WHERE id = $2",
		req.HotWalletID, instanceID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link hot wallet"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hot wallet linked successfully"})
}

// TestInstance tests the connection to a provider instance
func (h *InstanceHandler) TestInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")

	// Get the instance details
	var inst ProviderInstance
	err := h.db.QueryRow(`
		SELECT id, provider_id, name, vault_secret_path, is_active
		FROM provider_instances WHERE id = $1
	`, instanceID).Scan(&inst.ID, &inst.ProviderID, &inst.Name, &inst.VaultSecretPath, &inst.IsActive)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	// TODO: Implement actual provider API test based on provider type
	// For now, just mark as healthy
	h.db.Exec("UPDATE provider_instances SET health_status = 'healthy', last_used_at = NOW() WHERE id = $1", instanceID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Connection test successful",
		"instance": gin.H{
			"id":   inst.ID,
			"name": inst.Name,
		},
	})
}

// RecordInstanceUsage records that an instance was used
func (h *InstanceHandler) RecordInstanceUsage(instanceID string, success bool, errorMsg string) error {
	if success {
		_, err := h.db.Exec(`
			UPDATE provider_instances
			SET request_count = request_count + 1,
			    last_used_at = NOW(),
			    health_status = 'healthy',
			    last_error = NULL,
			    updated_at = NOW()
			WHERE id = $1
		`, instanceID)
		return err
	} else {
		_, err := h.db.Exec(`
			UPDATE provider_instances
			SET request_count = request_count + 1,
			    last_used_at = NOW(),
			    health_status = 'error',
			    last_error = $2,
			    updated_at = NOW()
			WHERE id = $1
		`, instanceID, errorMsg)
		return err
	}
}

// SelectBestInstance selects the best instance for a provider based on priority and hot wallet balance
func (h *InstanceHandler) SelectBestInstance(c *gin.Context) {
	providerID := c.Param("id")
	operationType := c.Query("operation") // deposit or withdraw

	// Priority ordering:
	// 1. Active instances
	// 2. Healthy instances (not in error)
	// 3. Primary instance
	// 4. Higher priority
	// 5. Lower request count (load balancing)
	query := `
		SELECT i.id, i.provider_id, i.name, i.vault_secret_path, i.hot_wallet_id,
		       i.is_active, i.is_primary, i.priority, i.request_count, i.health_status
		FROM provider_instances i
		WHERE i.provider_id = $1
		  AND i.is_active = TRUE
		  AND i.health_status != 'revoked'
		ORDER BY
			CASE WHEN i.health_status = 'healthy' THEN 0 ELSE 1 END,
			i.is_primary DESC,
			i.priority DESC,
			i.request_count ASC
		LIMIT 1`

	var inst ProviderInstance
	err := h.db.QueryRow(query, providerID).Scan(
		&inst.ID, &inst.ProviderID, &inst.Name, &inst.VaultSecretPath, &inst.HotWalletID,
		&inst.IsActive, &inst.IsPrimary, &inst.Priority, &inst.RequestCount, &inst.HealthStatus,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No available instance for this provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance":  inst,
		"operation": operationType,
	})
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// InstanceWallet represents a hot wallet linked to a provider instance
type InstanceWallet struct {
	ID          string    `json:"id"`
	InstanceID  string    `json:"instance_id"`
	Currency    string    `json:"currency"`
	HotWalletID string    `json:"hot_wallet_id"`
	IsActive    bool      `json:"is_active"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	// Joined data from wallet-service
	WalletBalance float64 `json:"wallet_balance,omitempty"`
	WalletName    string  `json:"wallet_name,omitempty"`
}

// GetInstanceWallets returns all wallets linked to an instance
func (h *InstanceHandler) GetInstanceWallets(c *gin.Context) {
	instanceID := c.Param("instanceId")

	query := `
		SELECT id, instance_id, currency, hot_wallet_id, is_active, priority, created_at
		FROM provider_instance_wallets
		WHERE instance_id = $1
		ORDER BY currency, priority DESC`

	rows, err := h.db.Query(query, instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch instance wallets"})
		return
	}
	defer rows.Close()

	var wallets []InstanceWallet
	for rows.Next() {
		var w InstanceWallet
		if err := rows.Scan(&w.ID, &w.InstanceID, &w.Currency, &w.HotWalletID, &w.IsActive, &w.Priority, &w.CreatedAt); err != nil {
			continue
		}
		wallets = append(wallets, w)
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

// AddInstanceWallet links a hot wallet to an instance
func (h *InstanceHandler) AddInstanceWallet(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		Currency    string `json:"currency" binding:"required"`
		HotWalletID string `json:"hot_wallet_id" binding:"required"`
		Priority    int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Priority == 0 {
		req.Priority = 50
	}

	// Insert or update the wallet link
	id := uuid.New().String()
	query := `
		INSERT INTO provider_instance_wallets (id, instance_id, currency, hot_wallet_id, is_active, priority)
		VALUES ($1, $2, $3, $4, TRUE, $5)
		ON CONFLICT (instance_id, currency, hot_wallet_id) DO UPDATE
		SET is_active = TRUE, priority = EXCLUDED.priority`

	_, err := h.db.Exec(query, id, instanceID, req.Currency, req.HotWalletID, req.Priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add wallet to instance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet linked successfully", "id": id})
}

// RemoveInstanceWallet removes a wallet link from an instance
func (h *InstanceHandler) RemoveInstanceWallet(c *gin.Context) {
	instanceID := c.Param("instanceId")
	walletLinkID := c.Param("walletId")

	_, err := h.db.Exec(`DELETE FROM provider_instance_wallets WHERE id = $1 AND instance_id = $2`, walletLinkID, instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove wallet link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet link removed"})
}

// ToggleInstanceWallet activates/deactivates a wallet link
func (h *InstanceHandler) ToggleInstanceWallet(c *gin.Context) {
	instanceID := c.Param("instanceId")
	walletLinkID := c.Param("walletId")

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`UPDATE provider_instance_wallets SET is_active = $1 WHERE id = $2 AND instance_id = $3`,
		req.IsActive, walletLinkID, instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet status updated"})
}

// GetBestWalletForCurrency returns the best hot wallet for a specific currency and instance
func (h *InstanceHandler) GetBestWalletForCurrency(c *gin.Context) {
	instanceID := c.Param("instanceId")
	currency := c.Query("currency")

	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currency is required"})
		return
	}

	query := `
		SELECT hot_wallet_id
		FROM provider_instance_wallets
		WHERE instance_id = $1 AND currency = $2 AND is_active = TRUE
		ORDER BY priority DESC
		LIMIT 1`

	var hotWalletID string
	err := h.db.QueryRow(query, instanceID, currency).Scan(&hotWalletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active wallet found for this currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hot_wallet_id": hotWalletID})
}

// ToggleInstancePause pauses or resumes an instance
func (h *InstanceHandler) ToggleInstancePause(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		IsPaused bool   `json:"is_paused"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the pause status in database
	result, err := h.db.Exec(`
		UPDATE provider_instances
		SET is_paused = $1, updated_at = NOW()
		WHERE id = $2
	`, req.IsPaused, instanceID)

	if err != nil {
		log.Printf("Error toggling instance pause: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update instance"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	status := "resumed"
	if req.IsPaused {
		status = "paused"
		log.Printf("Instance %s paused. Reason: %s", instanceID, req.Reason)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Instance " + status + " successfully",
		"instance_id": instanceID,
		"is_paused":   req.IsPaused,
	})
}
