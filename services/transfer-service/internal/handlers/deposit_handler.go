package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/providers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DepositHandler handles deposit collection flows with instance-based providers
type DepositHandler struct {
	instanceRepo      *repository.AggregatorInstanceRepository
	providerLoader    *providers.InstanceBasedProviderLoader
	fundMovement      *service.DepositFundMovementService
	walletServiceURL  string
}

// NewDepositHandler creates a new deposit handler
func NewDepositHandler(
	instanceRepo *repository.AggregatorInstanceRepository,
	providerLoader *providers.InstanceBasedProviderLoader,
	fundMovement *service.DepositFundMovementService,
	walletServiceURL string,
) *DepositHandler {
	return &DepositHandler{
		instanceRepo:     instanceRepo,
		providerLoader:   providerLoader,
		fundMovement:     fundMovement,
		walletServiceURL: walletServiceURL,
	}
}

// InitiateDepositRequest represents a deposit initiation request
type InitiateDepositRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Currency  string  `json:"currency" binding:"required"`
	Provider  string  `json:"provider" binding:"required"` // demo, flutterwave, stripe, etc.
	Country   string  `json:"country" binding:"required"`  // For routing
	Email     string  `json:"email"`                       // For payment provider
	Phone     string  `json:"phone"`                       // For mobile money
	ReturnURL string  `json:"return_url"`                  // Callback URL
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
	NewBalance    *float64 `json:"new_balance,omitempty"`
}

// InitiateDeposit initiates a deposit using instance-based provider
// POST /api/v1/deposits/initiate
func (h *DepositHandler) InitiateDeposit(c *gin.Context) {
	var req InitiateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	transactionID := fmt.Sprintf("dep_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Handle demo provider separately
	if req.Provider == "demo" {
		h.processInstantDemoDeposit(c, req, transactionID)
		return
	}

	// Get best instance for this provider
	provider, instance, err := h.providerLoader.GetBestProviderForDeposit(
		ctx,
		req.Provider,
		req.Country,
		req.Amount,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Provider not available: %v", err)})
		return
	}

	// Select best hot wallet for this instance
	walletID, err := h.fundMovement.SelectBestWalletForDeposit(ctx, instance.ID, req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No available wallet: %v", err)})
		return
	}

	// Prepare collection request
	collectionReq := &providers.CollectionRequest{
		ReferenceID: transactionID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		FirstName:   "",
		LastName:    "",
		RedirectURL: req.ReturnURL,
	}

	// Initiate collection with provider
	response, err := provider.InitiateCollection(ctx, collectionReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to initiate payment: %v", err),
		})
		return
	}

	// Record transaction in instance
	_ = h.providerLoader.RecordTransactionUsage(
		ctx,
		instance.ID,
		transactionID,
		req.Amount,
		req.Currency,
		"pending",
		response.ProviderReference,
	)

	// TODO: Store transaction reference and wallet_id for webhook processing

	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "pending_payment",
		PaymentURL:    response.PaymentLink,
		Provider:      req.Provider,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       response.Message,
	})
}

// processInstantDemoDeposit handles instant credit for demo providers
func (h *DepositHandler) processInstantDemoDeposit(c *gin.Context, req InitiateDepositRequest, transactionID string) {
	// For demo mode: instant success
	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "instant_success",
		Provider:      "demo",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Demo deposit completed instantly",
	})
}

// ProcessDepositWebhook processes webhook from payment provider
// POST /api/v1/deposits/webhook/:provider
func (h *DepositHandler) ProcessDepositWebhook(c *gin.Context) {
	providerCode := c.Param("provider")
	ctx := c.Request.Context()

	// TODO: Verify webhook signature based on provider

	// Get transaction details from webhook body
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract transaction reference and status
	// This depends on each provider's webhook format
	transactionRef := "" // Extract from webhookData
	status := ""         // Extract from webhookData

	if status == "successful" || status == "completed" {
		// Get transaction details from DB
		// Find which instance and wallet was used
		// Process fund movement from hot wallet to user account
		
		userID := "" // Get from transaction record
		walletID := "" // Get from transaction record
		amount := 0.0  // Get from transaction record
		currency := "" // Get from transaction record
		
		err := h.fundMovement.ProcessDepositFromWallet(
			ctx,
			userID,
			walletID,
			amount,
			currency,
			transactionRef,
			providerCode,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Check if wallet needs recharge
		// instanceWalletID := "" // Get from transaction record
		// _ = h.fundMovement.CheckAndTriggerRecharge(ctx, instanceWalletID)
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

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
