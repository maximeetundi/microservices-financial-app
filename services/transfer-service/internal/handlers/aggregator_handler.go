package handlers

import (
	"log"
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// AggregatorHandler handles aggregator configuration endpoints
type AggregatorHandler struct {
	repo *repository.AggregatorRepository
}

// NewAggregatorHandler creates a new aggregator handler
func NewAggregatorHandler(repo *repository.AggregatorRepository) *AggregatorHandler {
	return &AggregatorHandler{repo: repo}
}

// ==================== ADMIN ENDPOINTS ====================

// GetAllAggregators returns all aggregators for admin panel
// GET /admin/aggregators
func (h *AggregatorHandler) GetAllAggregators(c *gin.Context) {
	aggregators, err := h.repo.GetAll()
	if err != nil {
		log.Printf("[AggregatorHandler] Error getting aggregators: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aggregators": aggregators,
		"count":       len(aggregators),
	})
}

// GetAggregator returns a single aggregator by code
// GET /admin/aggregators/:code
func (h *AggregatorHandler) GetAggregator(c *gin.Context) {
	code := c.Param("code")

	aggregator, err := h.repo.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregator"})
		return
	}
	if aggregator == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aggregator not found"})
		return
	}

	c.JSON(http.StatusOK, aggregator)
}

// UpdateAggregator updates an aggregator's settings
// PATCH /admin/aggregators/:code
func (h *AggregatorHandler) UpdateAggregator(c *gin.Context) {
	code := c.Param("code")

	var req models.UpdateAggregatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if aggregator exists
	existing, err := h.repo.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aggregator not found"})
		return
	}

	// Update
	if err := h.repo.Update(code, &req); err != nil {
		log.Printf("[AggregatorHandler] Error updating %s: %v", code, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update aggregator"})
		return
	}

	// Return updated aggregator
	updated, _ := h.repo.GetByCode(code)
	log.Printf("[AggregatorHandler] ✅ Updated aggregator %s", code)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Aggregator updated successfully",
		"aggregator": updated,
	})
}

// EnableAggregator enables an aggregator
// POST /admin/aggregators/:code/enable
func (h *AggregatorHandler) EnableAggregator(c *gin.Context) {
	code := c.Param("code")

	if err := h.repo.SetEnabled(code, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable aggregator"})
		return
	}

	log.Printf("[AggregatorHandler] ✅ Enabled aggregator: %s", code)
	c.JSON(http.StatusOK, gin.H{"message": "Aggregator enabled", "code": code, "enabled": true})
}

// DisableAggregator disables an aggregator
// POST /admin/aggregators/:code/disable
func (h *AggregatorHandler) DisableAggregator(c *gin.Context) {
	code := c.Param("code")

	if err := h.repo.SetEnabled(code, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable aggregator"})
		return
	}

	log.Printf("[AggregatorHandler] ⛔ Disabled aggregator: %s", code)
	c.JSON(http.StatusOK, gin.H{"message": "Aggregator disabled", "code": code, "enabled": false})
}

// SetMaintenanceMode puts an aggregator in maintenance mode
// POST /admin/aggregators/:code/maintenance
func (h *AggregatorHandler) SetMaintenanceMode(c *gin.Context) {
	code := c.Param("code")

	var req struct {
		Enabled bool   `json:"enabled"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.SetMaintenanceMode(code, req.Enabled, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set maintenance mode"})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}
	log.Printf("[AggregatorHandler] 🔧 Maintenance mode %s for: %s", status, code)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Maintenance mode " + status,
		"code":        code,
		"maintenance": req.Enabled,
	})
}

// ToggleDeposit enables/disables deposits for an aggregator
// POST /admin/aggregators/:code/toggle-deposit
func (h *AggregatorHandler) ToggleDeposit(c *gin.Context) {
	code := c.Param("code")

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &models.UpdateAggregatorRequest{DepositEnabled: &req.Enabled}
	if err := h.repo.Update(code, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}
	log.Printf("[AggregatorHandler] 💰 Deposits %s for: %s", status, code)
	c.JSON(http.StatusOK, gin.H{"message": "Deposits " + status, "code": code, "deposit_enabled": req.Enabled})
}

// ToggleWithdraw enables/disables withdrawals for an aggregator
// POST /admin/aggregators/:code/toggle-withdraw
func (h *AggregatorHandler) ToggleWithdraw(c *gin.Context) {
	code := c.Param("code")

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &models.UpdateAggregatorRequest{WithdrawEnabled: &req.Enabled}
	if err := h.repo.Update(code, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}
	log.Printf("[AggregatorHandler] 💸 Withdrawals %s for: %s", status, code)
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawals " + status, "code": code, "withdraw_enabled": req.Enabled})
}

// ==================== PUBLIC ENDPOINTS (for frontend) ====================

// GetAvailableForDeposit returns aggregators available for deposits
// GET /aggregators/deposit?country=XX
func (h *AggregatorHandler) GetAvailableForDeposit(c *gin.Context) {
	country := c.DefaultQuery("country", "*")

	aggregators, err := h.repo.GetEnabledForDeposit(country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aggregators": aggregators,
		"count":       len(aggregators),
		"country":     country,
	})
}

// GetAvailableForWithdraw returns aggregators available for withdrawals
// GET /aggregators/withdraw?country=XX
func (h *AggregatorHandler) GetAvailableForWithdraw(c *gin.Context) {
	country := c.DefaultQuery("country", "*")

	aggregators, err := h.repo.GetEnabledForWithdraw(country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aggregators": aggregators,
		"count":       len(aggregators),
		"country":     country,
	})
}

// GetAllPublic returns all enabled aggregators for frontend display
// GET /aggregators
func (h *AggregatorHandler) GetAllPublic(c *gin.Context) {
	aggregators, err := h.repo.GetEnabled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch aggregators"})
		return
	}

	// Convert to frontend format
	var result []models.AggregatorForFrontend
	for _, a := range aggregators {
		result = append(result, models.AggregatorForFrontend{
			Code:            a.ProviderCode,
			Name:            a.ProviderName,
			LogoURL:         a.LogoURL,
			DepositEnabled:  a.DepositEnabled,
			WithdrawEnabled: a.WithdrawEnabled,
			MinAmount:       a.MinAmount,
			MaxAmount:       a.MaxAmount,
			FeePercent:      a.FeePercent,
			FeeFixed:        a.FeeFixed,
			FeeCurrency:     a.FeeCurrency,
			MaintenanceMode: a.MaintenanceMode,
			MaintenanceMsg:  a.MaintenanceMsg,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"aggregators": result,
		"count":       len(result),
	})
}

// RegisterRoutes registers all aggregator routes
func (h *AggregatorHandler) RegisterRoutes(router *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	// Public routes (for frontend)
	public := router.Group("/aggregators")
	{
		public.GET("", h.GetAllPublic)
		public.GET("/deposit", h.GetAvailableForDeposit)
		public.GET("/withdraw", h.GetAvailableForWithdraw)
	}

	// Admin routes (require admin auth middleware)
	admin := adminRouter.Group("/aggregators")
	{
		admin.GET("", h.GetAllAggregators)
		admin.GET("/:code", h.GetAggregator)
		admin.PATCH("/:code", h.UpdateAggregator)
		admin.POST("/:code/enable", h.EnableAggregator)
		admin.POST("/:code/disable", h.DisableAggregator)
		admin.POST("/:code/maintenance", h.SetMaintenanceMode)
		admin.POST("/:code/toggle-deposit", h.ToggleDeposit)
		admin.POST("/:code/toggle-withdraw", h.ToggleWithdraw)
	}
}
