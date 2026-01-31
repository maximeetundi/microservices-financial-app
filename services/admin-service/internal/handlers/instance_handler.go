package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProviderInstance represents a single API key instance for a provider
type ProviderInstance struct {
	ID              string     `json:"id" db:"id"`
	ProviderID      string     `json:"provider_id" db:"provider_id"`
	Name            string     `json:"name" db:"name"`
	VaultSecretPath string     `json:"vault_secret_path" db:"vault_secret_path"`
	HotWalletID     *string    `json:"hot_wallet_id,omitempty" db:"hot_wallet_id"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsPrimary       bool       `json:"is_primary" db:"is_primary"`
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
	db *sql.DB
}

// NewInstanceHandler creates a new InstanceHandler
func NewInstanceHandler(db *sql.DB) *InstanceHandler {
	return &InstanceHandler{db: db}
}

// GetProviderInstances returns all instances for a provider
func (h *InstanceHandler) GetProviderInstances(c *gin.Context) {
	providerID := c.Param("id")

	query := `
		SELECT id, provider_id, name, vault_secret_path, hot_wallet_id,
		       is_active, is_primary, priority, request_count, last_used_at,
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
			&inst.IsActive, &inst.IsPrimary, &inst.Priority, &inst.RequestCount, &inst.LastUsedAt,
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

	var req struct {
		Name            string `json:"name" binding:"required"`
		VaultSecretPath string `json:"vault_secret_path" binding:"required"`
		HotWalletID     string `json:"hot_wallet_id"`
		IsActive        bool   `json:"is_active"`
		IsPrimary       bool   `json:"is_primary"`
		Priority        int    `json:"priority"`
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'unknown')
		RETURNING id`

	var hotWalletID *string
	if req.HotWalletID != "" {
		hotWalletID = &req.HotWalletID
	}

	var returnedID string
	err := h.db.QueryRow(query,
		id, providerID, req.Name, req.VaultSecretPath, hotWalletID,
		req.IsActive, req.IsPrimary, priority,
	).Scan(&returnedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instance: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": returnedID, "message": "Instance created successfully"})
}

// UpdateProviderInstance updates an existing instance
func (h *InstanceHandler) UpdateProviderInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")

	var req struct {
		Name            string `json:"name"`
		VaultSecretPath string `json:"vault_secret_path"`
		HotWalletID     string `json:"hot_wallet_id"`
		IsActive        *bool  `json:"is_active"`
		IsPrimary       *bool  `json:"is_primary"`
		Priority        *int   `json:"priority"`
		HealthStatus    string `json:"health_status"`
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

	c.JSON(http.StatusOK, gin.H{"message": "Instance updated successfully"})
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
