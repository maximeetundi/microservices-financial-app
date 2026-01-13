package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	service *services.PayrollService
}

func NewPayrollHandler(service *services.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// PreviewPayroll
func (h *PayrollHandler) PreviewPayroll(c *gin.Context) {
	enterpriseID := c.Param("id")
	run, err := h.service.PreparePayroll(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, run)
}

// ExecutePayroll
func (h *PayrollHandler) ExecutePayroll(c *gin.Context) {
	enterpriseID := c.Param("id")
	
	// 1. Prepare
	run, err := h.service.PreparePayroll(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare payroll", "details": err.Error()})
		return
	}

	// 2. Execute
	if err := h.service.ExecutePayroll(c.Request.Context(), run); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute payroll", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll executed successfully", "run_id": run.ID})
}
