package handlers

import (
	"database/sql"
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// AdminInstanceHandler handles admin operations for aggregator instances
type AdminInstanceHandler struct {
	instanceRepo *repository.AggregatorInstanceRepository
	aggRepo      *repository.AggregatorRepository
}

func NewAdminInstanceHandler(
	instanceRepo *repository.AggregatorInstanceRepository,
	aggRepo *repository.AggregatorRepository,
) *AdminInstanceHandler {
	return &AdminInstanceHandler{
		instanceRepo: instanceRepo,
		aggRepo:      aggRepo,
	}
}

// ListInstances returns all instances with their wallets
// GET /api/v1/admin/instances
func (h *AdminInstanceHandler) ListInstances(c *gin.Context) {
	aggregatorID := c.Query("aggregator_id")
	_ = aggregatorID // Will be used for filtering

	// TODO: Get all instances
	// For now return empty array
	c.JSON(http.StatusOK, gin.H{"instances": []interface{}{}})
}

// GetInstance returns a specific instance with details
// GET /api/v1/admin/instances/:id
func (h *AdminInstanceHandler) GetInstance(c *gin.Context) {
	instanceID := c.Param("id")
	ctx := c.Request.Context()

	instance, err := h.instanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// CreateInstance creates a new aggregator instance
// POST /api/v1/admin/instances
func (h *AdminInstanceHandler) CreateInstance(c *gin.Context) {
	var req models.CreateAggregatorInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	adminID := c.GetString("admin_id") // From auth middleware

	instance, err := h.instanceRepo.Create(ctx, &req, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// LinkWallet links a hot wallet to an instance
// POST /api/v1/admin/instances/:id/wallets
func (h *AdminInstanceHandler) LinkWallet(c *gin.Context) {
	instanceID := c.Param("id")

	var req models.CreateInstanceWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.InstanceID = instanceID

	// TODO: Implement wallet linking
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet linked successfully"})
}

// UnlinkWallet removes a wallet from an instance
// DELETE /api/v1/admin/instances/:id/wallets/:wallet_id
func (h *AdminInstanceHandler) UnlinkWallet(c *gin.Context) {
	// TODO: Implement wallet unlinking
	c.JSON(http.StatusOK, gin.H{"message": "Wallet unlinked successfully"})
}

// UpdateInstanceWallet updates wallet configuration
// PUT /api/v1/admin/instances/:id/wallets/:wallet_id
func (h *AdminInstanceHandler) UpdateInstanceWallet(c *gin.Context) {
	// TODO: Implement wallet update
	c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
}

// GetInstanceStats returns statistics for an instance
// GET /api/v1/admin/instances/:id/stats
func (h *AdminInstanceHandler) GetInstanceStats(c *gin.Context) {
	instanceID := c.Param("id")

	// TODO: Get stats from aggregator_instance_transactions
	c.JSON(http.StatusOK, gin.H{
		"instance_id":        instanceID,
		"total_transactions": 0,
		"total_volume":       0,
		"total_deposits":     0,
		"total_withdrawals":  0,
	})
}
