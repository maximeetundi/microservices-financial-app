package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/providers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SimpleDepositHandler handles deposit requests using provider configurations
type SimpleDepositHandler struct {
	providerConfig *providers.Config
}

// NewSimpleDepositHandler creates a new simple deposit handler
func NewSimpleDepositHandler(providerConfig *providers.Config) *SimpleDepositHandler {
	return &SimpleDepositHandler{
		providerConfig: providerConfig,
	}
}

// SimpleDepositRequest represents a deposit initiation request
type SimpleDepositRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Currency  string  `json:"currency" binding:"required"`
	Provider  string  `json:"provider" binding:"required"` // flutterwave, stripe, cinetpay, etc.
	Country   string  `json:"country"`                     // For routing
	Email     string  `json:"email"`                       // For payment provider
	Phone     string  `json:"phone"`                       // For mobile money
	ReturnURL string  `json:"return_url"`                  // Callback URL
}

// SimpleDepositResponse represents the response from deposit initiation
type SimpleDepositResponse struct {
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"` // pending_payment, instant_success, failed
	PaymentURL    string  `json:"payment_url,omitempty"`
	Provider      string  `json:"provider"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Message       string  `json:"message,omitempty"`
}

// InitiateDeposit handles deposit initiation requests
// POST /api/v1/deposits/initiate
func (h *SimpleDepositHandler) InitiateDeposit(c *gin.Context) {
	var req SimpleDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	transactionID := fmt.Sprintf("dep_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Handle demo provider - instant success
	if req.Provider == "demo" || req.Provider == "Mode Test" || req.Provider == "Mode Démo" {
		c.JSON(http.StatusOK, SimpleDepositResponse{
			TransactionID: transactionID,
			Status:        "instant_success",
			Provider:      "demo",
			Amount:        req.Amount,
			Currency:      req.Currency,
			Message:       "Demo deposit - instant success",
		})
		return
	}

	// Get provider based on req.Provider
	switch req.Provider {
	case "flutterwave", "Flutterwave":
		paymentURL, err := h.initiateFlutterwaveDeposit(ctx, req, transactionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to initiate payment",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, SimpleDepositResponse{
			TransactionID: transactionID,
			Status:        "pending_payment",
			PaymentURL:    paymentURL,
			Provider:      req.Provider,
			Amount:        req.Amount,
			Currency:      req.Currency,
			Message:       "Redirect to payment page",
		})

	case "cinetpay", "CinetPay":
		paymentURL, err := h.initiateCinetPayDeposit(ctx, req, transactionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to initiate payment",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, SimpleDepositResponse{
			TransactionID: transactionID,
			Status:        "pending_payment",
			PaymentURL:    paymentURL,
			Provider:      req.Provider,
			Amount:        req.Amount,
			Currency:      req.Currency,
			Message:       "Redirect to payment page",
		})

	case "stripe", "Stripe":
		paymentURL, err := h.initiateStripeDeposit(ctx, req, transactionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to initiate payment",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, SimpleDepositResponse{
			TransactionID: transactionID,
			Status:        "pending_payment",
			PaymentURL:    paymentURL,
			Provider:      req.Provider,
			Amount:        req.Amount,
			Currency:      req.Currency,
			Message:       "Redirect to payment page",
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Unsupported provider",
			"details": fmt.Sprintf("Provider '%s' is not supported for deposits", req.Provider),
		})
	}
}

// initiateFlutterwaveDeposit creates payment via Flutterwave
func (h *SimpleDepositHandler) initiateFlutterwaveDeposit(ctx context.Context, req SimpleDepositRequest, txRef string) (string, error) {
	cfg := h.providerConfig.Flutterwave
	if cfg.SecretKey == "" {
		return "", fmt.Errorf("flutterwave not configured")
	}

	// Use the collection provider wrapper
	collectionProvider := providers.NewFlutterwaveCollectionProvider(cfg)

	collectionReq := &providers.CollectionRequest{
		ReferenceID: txRef,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		RedirectURL: req.ReturnURL,
	}

	response, err := collectionProvider.InitiateCollection(ctx, collectionReq)
	if err != nil {
		return "", err
	}

	return response.PaymentLink, nil
}

// initiateCinetPayDeposit creates payment via CinetPay
func (h *SimpleDepositHandler) initiateCinetPayDeposit(ctx context.Context, req SimpleDepositRequest, txRef string) (string, error) {
	cfg := h.providerConfig.CinetPay
	if cfg.APIKey == "" {
		return "", fmt.Errorf("cinetpay not configured")
	}

	// Use the collection provider wrapper
	collectionProvider := providers.NewCinetPayCollectionProvider(cfg)

	collectionReq := &providers.CollectionRequest{
		ReferenceID: txRef,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		RedirectURL: req.ReturnURL,
	}

	response, err := collectionProvider.InitiateCollection(ctx, collectionReq)
	if err != nil {
		return "", err
	}

	return response.PaymentLink, nil
}

// initiateStripeDeposit creates payment via Stripe
func (h *SimpleDepositHandler) initiateStripeDeposit(ctx context.Context, req SimpleDepositRequest, txRef string) (string, error) {
	cfg := h.providerConfig.Stripe
	if cfg.SecretKey == "" {
		return "", fmt.Errorf("stripe not configured")
	}

	// Use the collection provider wrapper
	collectionProvider := providers.NewStripeCollectionProvider(cfg)

	collectionReq := &providers.CollectionRequest{
		ReferenceID: txRef,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		RedirectURL: req.ReturnURL,
	}

	response, err := collectionProvider.InitiateCollection(ctx, collectionReq)
	if err != nil {
		return "", err
	}

	return response.PaymentLink, nil
}

// GetDepositStatus returns the status of a deposit
// GET /api/v1/deposits/:id/status
func (h *SimpleDepositHandler) GetDepositStatus(c *gin.Context) {
	transactionID := c.Param("id")

	// TODO: Query transaction from database
	c.JSON(http.StatusOK, gin.H{
		"transaction_id": transactionID,
		"status":         "pending",
		"message":        "Payment is being processed",
	})
}

// HandleDepositWebhook handles webhooks from payment providers
// POST /api/v1/deposits/webhook/:provider
func (h *SimpleDepositHandler) HandleDepositWebhook(c *gin.Context) {
	providerCode := c.Param("provider")

	// TODO: Validate webhook signature and process
	c.JSON(http.StatusOK, gin.H{
		"status":   "received",
		"provider": providerCode,
	})
}
