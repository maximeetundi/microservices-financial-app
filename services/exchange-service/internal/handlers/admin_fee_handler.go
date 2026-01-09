package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
)

type AdminFeeHandler struct {
	feeService *services.FeeService
}

func NewAdminFeeHandler(feeService *services.FeeService) *AdminFeeHandler {
	return &AdminFeeHandler{feeService: feeService}
}

func (h *AdminFeeHandler) GetFees(c *gin.Context) {
	configs, err := h.feeService.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fee configs"})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *AdminFeeHandler) UpdateFee(c *gin.Context) {
	var req models.FeeConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" && req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fee config ID or Key is required"})
		return
	}

	err := h.feeService.UpdateConfig(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fee config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fee configuration updated successfully"})
}
