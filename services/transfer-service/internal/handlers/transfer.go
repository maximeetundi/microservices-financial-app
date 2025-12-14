package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/services"
)

type TransferHandler struct {
	transferService      *services.TransferService
	mobileMoneyService   *services.MobileMoneyService
	internationalService *services.InternationalTransferService
	complianceService    *services.ComplianceService
}

func NewTransferHandler(
	transferService *services.TransferService,
	mobileMoneyService *services.MobileMoneyService,
	internationalService *services.InternationalTransferService,
	complianceService *services.ComplianceService,
) *TransferHandler {
	return &TransferHandler{
		transferService:      transferService,
		mobileMoneyService:   mobileMoneyService,
		internationalService: internationalService,
		complianceService:    complianceService,
	}
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer, err := h.transferService.CreateTransfer(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Run compliance check
	complianceResult, _ := h.complianceService.CheckTransfer(transfer)
	
	c.JSON(http.StatusCreated, gin.H{
		"transfer":   transfer,
		"compliance": complianceResult,
	})
}

func (h *TransferHandler) GetTransfer(c *gin.Context) {
	transferID := c.Param("transfer_id")
	
	transfer, err := h.transferService.GetTransfer(transferID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transfer not found"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (h *TransferHandler) GetTransferHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transfers, err := h.transferService.GetTransferHistory(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transfers": transfers,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *TransferHandler) CancelTransfer(c *gin.Context) {
	transferID := c.Param("transfer_id")
	
	if err := h.transferService.CancelTransfer(transferID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer cancelled successfully"})
}

func (h *TransferHandler) CreateInternationalTransfer(c *gin.Context) {
	var req models.InternationalTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer, err := h.internationalService.CreateTransfer(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

// Mobile Money handlers
func (h *TransferHandler) SendMobileMoney(c *gin.Context) {
	var req models.MobileMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.mobileMoneyService.Send(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransferHandler) ReceiveMobileMoney(c *gin.Context) {
	var req models.MobileMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.mobileMoneyService.Receive(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransferHandler) GetMobileProviders(c *gin.Context) {
	providers := h.mobileMoneyService.GetProviders()
	c.JSON(http.StatusOK, gin.H{"providers": providers})
}

// Bulk transfer handlers
func (h *TransferHandler) CreateBulkTransfer(c *gin.Context) {
	var req models.BulkTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"batch_id": "batch_" + strconv.FormatInt(int64(len(req.Transfers)), 10),
		"status":   "pending",
		"count":    len(req.Transfers),
	})
}

func (h *TransferHandler) GetBulkTransferStatus(c *gin.Context) {
	batchID := c.Param("batch_id")
	c.JSON(http.StatusOK, gin.H{
		"batch_id": batchID,
		"status":   "processing",
	})
}

func (h *TransferHandler) ApproveBulkTransfer(c *gin.Context) {
	batchID := c.Param("batch_id")
	c.JSON(http.StatusOK, gin.H{
		"batch_id": batchID,
		"status":   "approved",
	})
}

// Webhook handlers
func (h *TransferHandler) HandleMobileMoneyCallback(c *gin.Context) {
	var callback map[string]interface{}
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *TransferHandler) HandleBankCallback(c *gin.Context) {
	var callback map[string]interface{}
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
