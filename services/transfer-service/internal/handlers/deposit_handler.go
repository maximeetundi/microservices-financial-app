package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/providers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DepositHandler handles deposit collection flows
type DepositHandler struct {
	repo             *repository.AggregatorRepository
	fullService      *providers.FullTransferService
	walletServiceURL string // URL to wallet-service for webhook callbacks
}

// NewDepositHandler creates a new deposit handler
func NewDepositHandler(repo *repository.AggregatorRepository, fullService *providers.FullTransferService, walletServiceURL string) *DepositHandler {
	return &DepositHandler{
		repo:             repo,
		fullService:      fullService,
		walletServiceURL: walletServiceURL,
	}
}

// InitiateDepositRequest represents a deposit initiation request
type InitiateDepositRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	WalletID  string  `json:"wallet_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Currency  string  `json:"currency" binding:"required"`
	Provider  string  `json:"provider" binding:"required"` // demo, flutterwave, stripe, etc.
	Country   string  `json:"country" binding:"required"`  // For routing
	ReturnURL string  `json:"return_url"`                  // Callback URL for user after payment
}

// InitiateDepositResponse represents the response from deposit initiation
type InitiateDepositResponse struct {
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"` // instant_success, pending_payment, failed
	PaymentURL    string  `json:"payment_url,omitempty"`
	Provider      string  `json:"provider"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Message       string  `json:"message,omitempty"`
	// For demo/instant success
	NewBalance *float64 `json:"new_balance,omitempty"`
}

// InitiateDeposit initiates a deposit from an aggregator
// POST /api/v1/deposits/initiate
func (h *DepositHandler) InitiateDeposit(c *gin.Context) {
	var req InitiateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get provider configuration
	provider, err := h.repo.GetAggregatorByName(req.Provider)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	if !provider.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider is not active"})
		return
	}

	if !provider.DepositEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposits are not enabled for this provider"})
		return
	}

	// Generate transaction ID
	transactionID := fmt.Sprintf("dep_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Check if this is demo mode
	if provider.IsDemoMode {
		// INSTANT CREDIT for demo mode
		return h.processInstantDemoDeposit(c, req, transactionID)
	}

	// Real provider flow - call actual API
	return h.processRealProviderDeposit(c, req, provider, transactionID)
}

// processInstantDemoDeposit handles instant credit for demo providers
func (h *DepositHandler) processInstantDemoDeposit(c *gin.Context, req InitiateDepositRequest, transactionID string) {
	// For demo mode, we credit instantly without calling external API
	// Call wallet-service to credit the wallet via platform reserve
	callbackReq := map[string]interface{}{
		"user_id":       req.UserID,
		"wallet_id":     req.WalletID,
		"amount":        req.Amount,
		"currency":      req.Currency,
		"provider_ref":  transactionID,
		"provider_name": "demo",
	}

	// TODO: Make HTTP call to wallet-service ProcessPlatformDeposit endpoint
	// For now, just return success indicating instant credit
	// The wallet-service will handle the actual crediting

	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "instant_success",
		Provider:      "demo",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Deposit credited instantly (demo mode)",
		// NewBalance will be set by wallet-service response
	})

	// In production, make async HTTP call to wallet-service
	// go h.callWalletServiceDeposit(callbackReq)
}

// processRealProviderDeposit handles deposits through real aggregator APIs
func (h *DepositHandler) processRealProviderDeposit(c *gin.Context, req InitiateDepositRequest, provider *models.Aggregator, transactionID string) {
	ctx := context.Background()

	// Prepare collection request for the provider
	collectionReq := &providers.CollectionRequest{
		Amount:       req.Amount,
		Currency:     req.Currency,
		Country:      req.Country,
		Email:        "", // TODO: Get from user profile
		PhoneNumber:  "", // TODO: Get from request or user profile
		CustomerName: "", // TODO: Get from user profile
		Reference:    transactionID,
		RedirectURL:  req.ReturnURL,
		WebhookURL:   fmt.Sprintf("%s/webhooks/%s/deposit", h.walletServiceURL, provider.Name),
		Metadata: map[string]interface{}{
			"user_id":   req.UserID,
			"wallet_id": req.WalletID,
		},
	}

	// Call the real provider API via FullTransferService
	var response *providers.CollectionResponse
	var err error

	switch provider.Name {
	case "flutterwave":
		response, err = h.fullService.GetFlutterwaveCollection().InitiateCollection(ctx, collectionReq)
	case "stripe":
		response, err = h.fullService.GetStripeCollection().InitiateCollection(ctx, collectionReq)
	default:
		// Use the generic InitiateDeposit which routes based on country
		response, err = h.fullService.InitiateDeposit(ctx, collectionReq)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initiate payment",
			"details": err.Error(),
		})
		return
	}

	// Return payment URL for user to complete payment
	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "pending_payment",
		PaymentURL:    response.PaymentURL,
		Provider:      provider.Name,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Please complete payment on the provider's page",
	})

	// Store transaction in database as pending
	// Will be updated when webhook is received
}

// GetDepositStatus returns the status of a deposit
// GET /api/v1/deposits/:id/status
func (h *DepositHandler) GetDepositStatus(c *gin.Context) {
	transactionID := c.Param("id")

	// TODO: Query transaction from database
	// For now, return a mock response

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": transactionID,
		"status":         "pending",
		"message":        "Payment is being processed",
	})
}
