package handlers

import (
	"fmt"
	"net/http"
	"time"

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

	var req struct {
		SourceWalletID string `json:"source_wallet_id"`
	}
	c.ShouldBindJSON(&req) // Optional, default to enterprise default wallet

	// 1. Prepare
	run, err := h.service.PreparePayroll(c.Request.Context(), enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare payroll", "details": err.Error()})
		return
	}

	// 2. Execute with source wallet ID
	if err := h.service.ExecutePayroll(c.Request.Context(), run, req.SourceWalletID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute payroll", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll executed successfully", "run_id": run.ID})
}

// ListPayrollRuns returns payroll history
func (h *PayrollHandler) ListPayrollRuns(c *gin.Context) {
	enterpriseID := c.Param("id")
	yearStr := c.DefaultQuery("year", "2025")
	year := 2025
	if _, err := fmt.Sscan(yearStr, &year); err != nil {
		year = time.Now().Year()
	}

	runs, err := h.service.ListPayrollRuns(c.Request.Context(), enterpriseID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payroll runs", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, runs)
}
