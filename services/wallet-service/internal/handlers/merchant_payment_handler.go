package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
	"github.com/gin-gonic/gin"
)

// MerchantPaymentHandler handles merchant payment API endpoints
type MerchantPaymentHandler struct {
	merchantService *services.MerchantPaymentService
}

// NewMerchantPaymentHandler creates a new merchant payment handler
func NewMerchantPaymentHandler(merchantService *services.MerchantPaymentService) *MerchantPaymentHandler {
	return &MerchantPaymentHandler{
		merchantService: merchantService,
	}
}

// CreatePaymentRequest creates a new payment request with QR code
// POST /merchant/payments
func (h *MerchantPaymentHandler) CreatePaymentRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreatePaymentRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.merchantService.CreatePaymentRequest(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPaymentRequest gets a payment request by ID
// GET /payments/:id
func (h *MerchantPaymentHandler) GetPaymentRequest(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID required"})
		return
	}

	payment, err := h.merchantService.GetPaymentRequest(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentQRCode gets/regenerates the QR code for a payment
// GET /payments/:id/qr
func (h *MerchantPaymentHandler) GetPaymentQRCode(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID required"})
		return
	}

	qrCode, err := h.merchantService.RegenerateQRCode(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_id":    paymentID,
		"qr_code_base64": qrCode,
	})
}

// PayPaymentRequest processes a payment for a payment request
// POST /payments/:id/pay
func (h *MerchantPaymentHandler) PayPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	paymentID := c.Param("id")
	
	var req models.PayPaymentRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	req.PaymentRequestID = paymentID

	history, err := h.merchantService.PayPaymentRequest(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful",
		"payment": history,
	})
}

// GetMerchantPayments gets all payment requests for the authenticated merchant
// GET /merchant/payments
func (h *MerchantPaymentHandler) GetMerchantPayments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	payments, err := h.merchantService.GetMerchantPayments(userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetMerchantHistory gets payment history for the authenticated merchant
// GET /merchant/payments/history
func (h *MerchantPaymentHandler) GetMerchantHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	history, err := h.merchantService.GetPaymentHistory(userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"limit":   limit,
		"offset":  offset,
	})
}

// CancelPaymentRequest cancels a payment request
// DELETE /merchant/payments/:id
func (h *MerchantPaymentHandler) CancelPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID required"})
		return
	}

	err := h.merchantService.CancelPaymentRequest(userID.(string), paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment request cancelled"})
}

// QuickPaymentRequest creates a quick payment (simplified endpoint)
// POST /merchant/quick-pay
func (h *MerchantPaymentHandler) QuickPaymentRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		WalletID    string   `json:"wallet_id" binding:"required"`
		Amount      *float64 `json:"amount"`
		Currency    string   `json:"currency" binding:"required"`
		Description string   `json:"description"`
		NeverExpires bool    `json:"never_expires"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine type based on amount
	paymentType := models.PaymentTypeVariable
	if req.Amount != nil && *req.Amount > 0 {
		paymentType = models.PaymentTypeFixed
	}

	// Set expiration
	var expiresIn *int
	if req.NeverExpires {
		neverExpires := -1
		expiresIn = &neverExpires
	}

	createReq := &models.CreatePaymentRequestDTO{
		Type:             paymentType,
		WalletID:         req.WalletID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		Title:            "Quick Payment",
		Description:      req.Description,
		ExpiresInMinutes: expiresIn,
		Reusable:         req.NeverExpires, // If never expires, likely reusable
	}

	response, err := h.merchantService.CreatePaymentRequest(userID.(string), createReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ScanPayment gets payment details for scanning (public endpoint)
// GET /pay/:id
func (h *MerchantPaymentHandler) ScanPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment ID required"})
		return
	}

	payment, err := h.merchantService.GetPaymentRequest(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return public info only
	c.JSON(http.StatusOK, gin.H{
		"payment_id":   payment.ID,
		"merchant_id":  payment.MerchantID,
		"type":         payment.Type,
		"amount":       payment.Amount,
		"min_amount":   payment.MinAmount,
		"max_amount":   payment.MaxAmount,
		"currency":     payment.Currency,
		"title":        payment.Title,
		"description":  payment.Description,
		"status":       payment.Status,
		"expires_at":   payment.ExpiresAt,
		"never_expires": payment.NeverExpires,
		"is_expired":   payment.Status == models.PaymentStatusExpired,
		"is_paid":      payment.Status == models.PaymentStatusPaid && !payment.Reusable,
	})
}
