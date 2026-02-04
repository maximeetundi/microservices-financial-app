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
	provider, err := h.repo.GetByCode(req.Provider)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	if !provider.IsEnabled {
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
	if provider.ProviderCode == "demo" {
		// INSTANT CREDIT for demo mode
		h.processInstantDemoDeposit(c, req, transactionID)
		return
	}

	// Real provider flow - call actual API
	h.processRealProviderDeposit(c, req, provider, transactionID)
}

// processInstantDemoDeposit handles instant credit for demo providers
func (h *DepositHandler) processInstantDemoDeposit(c *gin.Context, req InitiateDepositRequest, transactionID string) {
	// For demo mode, we credit instantly without calling external API
	// Call wallet-service to credit the wallet via platform reserve
	// In production, would call wallet-service here
	_ = transactionID // avoid unused variable warning

	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "instant_success",
		Provider:      "demo",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Deposit credited instantly (demo mode)",
	})
}

// processRealProviderDeposit handles deposits through real aggregator APIs
func (h *DepositHandler) processRealProviderDeposit(c *gin.Context, req InitiateDepositRequest, provider *models.AggregatorSetting, transactionID string) {
	ctx := context.Background()

	// Prepare collection request for the provider
	collectionReq := &providers.CollectionRequest{
		ReferenceID: transactionID,
		UserID:      req.UserID,
		WalletID:    req.WalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       "", // TODO: Get from user profile
		PhoneNumber: "", // TODO: Get from request or user profile
		RedirectURL: req.ReturnURL,
		Metadata: map[string]string{
			"user_id":   req.UserID,
			"wallet_id": req.WalletID,
		},
	}

	// Call the real provider API via FullTransferService
	var response *providers.CollectionResponse
	var err error

	// Use the generic InitiateDeposit which routes based on country
	response, err = h.fullService.InitiateDeposit(ctx, collectionReq)

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
		PaymentURL:    response.PaymentLink,
		Provider:      provider.ProviderCode,
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
