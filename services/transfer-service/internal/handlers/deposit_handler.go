package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	depositRepo      *repository.DepositRepository
	depositNumberRepo *repository.DepositNumberRepository
	instanceRepo     *repository.AggregatorInstanceRepository
	providerLoader   *providers.InstanceBasedProviderLoader
	fundMovement     *service.DepositFundMovementService
	walletServiceURL string
	webhookSecrets   map[string]string // provider -> webhook secret
}

// NewDepositHandler creates a new deposit handler
func NewDepositHandler(
	depositRepo *repository.DepositRepository,
	depositNumberRepo *repository.DepositNumberRepository,
	instanceRepo *repository.AggregatorInstanceRepository,
	providerLoader *providers.InstanceBasedProviderLoader,
	fundMovement *service.DepositFundMovementService,
	walletServiceURL string,
	webhookSecrets map[string]string,
) *DepositHandler {
	if webhookSecrets == nil {
		webhookSecrets = make(map[string]string)
	}
	return &DepositHandler{
		depositRepo:      depositRepo,
		depositNumberRepo: depositNumberRepo,
		instanceRepo:     instanceRepo,
		providerLoader:   providerLoader,
		fundMovement:     fundMovement,
		walletServiceURL: walletServiceURL,
		webhookSecrets:   webhookSecrets,
	}
}

// InitiateDepositRequest represents a deposit initiation request
type InitiateDepositRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	WalletID  string  `json:"wallet_id"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Currency  string  `json:"currency" binding:"required"`
	Provider  string  `json:"provider" binding:"required"`
	Country   string  `json:"country" binding:"required"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	DepositNumberID string `json:"deposit_number_id"`
	ReturnURL string  `json:"return_url"`
	CancelURL string  `json:"cancel_url"`
}

// InitiateDepositResponse represents the response from deposit initiation
type InitiateDepositResponse struct {
	TransactionID        string     `json:"transaction_id"`
	Status               string     `json:"status"`
	PaymentURL           string     `json:"payment_url,omitempty"`
	Provider             string     `json:"provider"`
	Amount               float64    `json:"amount"`
	Currency             string     `json:"currency"`
	Fee                  float64    `json:"fee"`
	Message              string     `json:"message,omitempty"`
	NewBalance           *float64   `json:"new_balance,omitempty"`
	AggregatorInstanceID string     `json:"aggregator_instance_id,omitempty"`
	HotWalletID          string     `json:"hot_wallet_id,omitempty"`
	ExpiresAt            string     `json:"expires_at,omitempty"`
	SDKConfig            *SDKConfig `json:"sdk_config,omitempty"`
}

// SDKConfig contains configuration for frontend SDK integration
type SDKConfig struct {
	PublicKey   string            `json:"public_key,omitempty"`
	Environment string            `json:"environment,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Country     string            `json:"country,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
}

// Default timeout for deposit transactions (2 hours)
const DefaultDepositTimeout = 2 * time.Hour

// InitiateDeposit initiates a deposit using instance-based provider
// POST /api/v1/deposits/initiate
func (h *DepositHandler) InitiateDeposit(c *gin.Context) {
	var req InitiateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enforce that caller can only initiate deposits for themselves
	authUserID, _ := c.Get("user_id")
	if authUserIDStr := fmt.Sprintf("%v", authUserID); authUserIDStr == "" || authUserIDStr != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// If a deposit number is selected, load it and override phone
	if strings.TrimSpace(req.DepositNumberID) != "" {
		n, err := h.depositNumberRepo.GetByID(c.Request.Context(), req.DepositNumberID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit number"})
			return
		}
		if n == nil || n.UserID != req.UserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit number"})
			return
		}
		if !n.IsVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit number not verified"})
			return
		}
		req.Phone = n.Phone
		if req.Country == "" {
			req.Country = n.Country
		}
	}

	ctx := c.Request.Context()
	transactionID := fmt.Sprintf("dep_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Handle demo provider separately (instant credit)
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
		log.Printf("[DepositHandler] Provider not available: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Provider not available",
			"details": err.Error(),
		})
		return
	}

	// Select best hot wallet for this instance
	walletID, err := h.fundMovement.SelectBestWalletForDeposit(ctx, instance.ID, req.Currency, req.Amount)
	if err != nil {
		log.Printf("[DepositHandler] No available wallet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No available wallet for this currency",
			"details": err.Error(),
		})
		return
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(DefaultDepositTimeout)

	// Create deposit transaction record
	depositReq := &repository.CreateDepositRequest{
		ID:                   transactionID,
		UserID:               req.UserID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Fee:                  0, // Will be updated after provider response
		ProviderCode:         req.Provider,
		AggregatorInstanceID: instance.ID,
		HotWalletID:          walletID,
		UserEmail:            req.Email,
		UserPhone:            req.Phone,
		UserCountry:          req.Country,
		ReturnURL:            req.ReturnURL,
		CancelURL:            req.CancelURL,
		ExpiresAt:            expiresAt,
		Metadata: map[string]interface{}{
			"user_wallet_id": req.WalletID,
			"ip_address":     c.ClientIP(),
		},
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	_, err = h.depositRepo.Create(ctx, depositReq)
	if err != nil {
		log.Printf("[DepositHandler] Failed to create deposit record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deposit"})
		return
	}

	// Prepare collection request for provider
	collectionReq := &providers.CollectionRequest{
		ReferenceID: transactionID,
		UserID:      req.UserID,
		WalletID:    walletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Country:     req.Country,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		RedirectURL: req.ReturnURL,
		Metadata: map[string]string{
			"transaction_id": transactionID,
			"user_id":        req.UserID,
			"wallet_id":      walletID,
		},
	}

	// Initiate collection with provider
	response, err := provider.InitiateCollection(ctx, collectionReq)
	if err != nil {
		log.Printf("[DepositHandler] Failed to initiate payment with %s: %v", req.Provider, err)
		// Mark transaction as failed
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Provider error: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initiate payment",
			"details": err.Error(),
		})
		return
	}

	// Update transaction with provider reference and payment URL
	_ = h.depositRepo.UpdateProviderReference(ctx, transactionID, response.ProviderReference, response.PaymentLink)

	// Record usage in instance
	_ = h.providerLoader.RecordTransactionUsage(
		ctx,
		instance.ID,
		transactionID,
		req.Amount,
		req.Currency,
		"pending",
		response.ProviderReference,
	)

	// Build SDK config for frontend
	sdkConfig := h.buildSDKConfig(req.Provider, instance, req)

	log.Printf("[DepositHandler] ✅ Deposit initiated: %s | Provider: %s | Amount: %.2f %s",
		transactionID, req.Provider, req.Amount, req.Currency)

	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID:        transactionID,
		Status:               "pending",
		PaymentURL:           response.PaymentLink,
		Provider:             req.Provider,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Fee:                  response.Fee,
		Message:              response.Message,
		AggregatorInstanceID: instance.ID,
		HotWalletID:          walletID,
		ExpiresAt:            expiresAt.Format(time.RFC3339),
		SDKConfig:            sdkConfig,
	})
}

// buildSDKConfig builds SDK configuration for frontend integration
func (h *DepositHandler) buildSDKConfig(provider string, instance *models.AggregatorInstanceWithDetails, req InitiateDepositRequest) *SDKConfig {
	config := &SDKConfig{
		Environment: "test",
		Currency:    req.Currency,
		Country:     req.Country,
		Extra:       make(map[string]string),
	}

	if !instance.IsTestMode {
		config.Environment = "live"
	}

	// Add provider-specific SDK configuration
	switch provider {
	case "flutterwave":
		if pk, ok := instance.APICredentials["public_key"]; ok {
			config.PublicKey = pk
		}
		config.Extra["payment_plan"] = ""
		config.Extra["customizations_title"] = "Zekora - Recharge"

	case "stripe":
		if pk, ok := instance.APICredentials["public_key"]; ok {
			config.PublicKey = pk
		}

	case "paystack":
		if pk, ok := instance.APICredentials["public_key"]; ok {
			config.PublicKey = pk
		}
		config.Extra["channels"] = "card,bank,ussd,mobile_money"

	case "cinetpay":
		if siteID, ok := instance.APICredentials["site_id"]; ok {
			config.Extra["site_id"] = siteID
		}
		if apiKey, ok := instance.APICredentials["api_key"]; ok {
			config.PublicKey = apiKey
		}

	case "paypal":
		if clientID, ok := instance.APICredentials["client_id"]; ok {
			config.PublicKey = clientID
		}

	case "lygos":
		if shopName, ok := instance.APICredentials["shop_name"]; ok {
			config.Extra["shop_name"] = shopName
		}

	case "wave", "wave_ci":
		config.Extra["merchant_name"] = "Zekora Finance"

	case "orange_money":
		if merchantKey, ok := instance.APICredentials["merchant_key"]; ok {
			config.Extra["merchant_key"] = merchantKey
		}

	case "mtn_momo":
		config.Extra["target_environment"] = config.Environment
	}

	return config
}

type CreatePayPalOrderRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// CreatePayPalOrder creates a PayPal order server-side for the JS SDK Buttons flow
// POST /api/v1/deposits/paypal/create-order
func (h *DepositHandler) CreatePayPalOrder(c *gin.Context) {
	var req CreatePayPalOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	deposit, err := h.depositRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load deposit"})
		return
	}
	if deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if userIDStr := fmt.Sprintf("%v", userID); userIDStr == "" || userIDStr != deposit.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if strings.ToLower(deposit.ProviderCode) != "paypal" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider for PayPal order"})
		return
	}
	if deposit.AggregatorInstanceID == nil || *deposit.AggregatorInstanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing aggregator instance"})
		return
	}

	instance, err := h.providerLoader.CredentialsClient().GetInstanceByID(*deposit.AggregatorInstanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load provider instance"})
		return
	}

	creds := instance.APICredentials
	mode := strings.ToLower(strings.TrimSpace(creds["mode"]))
	if mode == "" {
		if instance.IsTestMode {
			mode = "sandbox"
		} else {
			mode = "live"
		}
	}

	config := providers.PayPalConfig{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		Mode:         mode,
		BaseURL:      creds["base_url"],
	}

	pp := providers.NewPayPalProvider(config)

	// Best-effort wallet id (used as custom_id on purchase_unit)
	walletID := ""
	if deposit.UserWalletID != nil {
		walletID = *deposit.UserWalletID
	}
	if deposit.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(deposit.Metadata, &metadata); err == nil {
			if v, ok := metadata["user_wallet_id"]; ok {
				walletID = fmt.Sprintf("%v", v)
			}
		}
	}

	order, err := pp.CreateOrder(ctx, deposit.Amount, deposit.Currency, walletID, "Wallet Deposit")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.depositRepo.UpdateProviderReference(ctx, deposit.ID, order.ID, "")

	c.JSON(http.StatusOK, gin.H{"order_id": order.ID})
}

type CapturePayPalOrderRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	OrderID       string `json:"order_id" binding:"required"`
}

// CapturePayPalOrder captures a PayPal order server-side and completes the deposit
// POST /api/v1/deposits/paypal/capture
func (h *DepositHandler) CapturePayPalOrder(c *gin.Context) {
	var req CapturePayPalOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	deposit, err := h.depositRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load deposit"})
		return
	}
	if deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if userIDStr := fmt.Sprintf("%v", userID); userIDStr == "" || userIDStr != deposit.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if strings.ToLower(deposit.ProviderCode) != "paypal" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider for PayPal capture"})
		return
	}
	if deposit.AggregatorInstanceID == nil || *deposit.AggregatorInstanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing aggregator instance"})
		return
	}

	instance, err := h.providerLoader.CredentialsClient().GetInstanceByID(*deposit.AggregatorInstanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load provider instance"})
		return
	}

	creds := instance.APICredentials
	mode := strings.ToLower(strings.TrimSpace(creds["mode"]))
	if mode == "" {
		if instance.IsTestMode {
			mode = "sandbox"
		} else {
			mode = "live"
		}
	}

	config := providers.PayPalConfig{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		Mode:         mode,
		BaseURL:      creds["base_url"],
	}

	pp := providers.NewPayPalProvider(config)
	order, err := pp.CaptureOrder(ctx, req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToUpper(order.Status) == "COMPLETED" {
		h.completeDeposit(ctx, deposit, 0)
		c.JSON(http.StatusOK, gin.H{"status": "completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "pending", "order_status": order.Status})
}

// processInstantDemoDeposit handles instant credit for demo providers
func (h *DepositHandler) processInstantDemoDeposit(c *gin.Context, req InitiateDepositRequest, transactionID string) {
	ctx := c.Request.Context()

	// Create deposit record
	depositReq := &repository.CreateDepositRequest{
		ID:           transactionID,
		UserID:       req.UserID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Fee:          0,
		ProviderCode: "demo",
		UserEmail:    req.Email,
		UserPhone:    req.Phone,
		UserCountry:  req.Country,
		ExpiresAt:    time.Now().Add(1 * time.Minute),
		Metadata: map[string]interface{}{
			"user_wallet_id": req.WalletID,
			"demo_mode":      true,
		},
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	_, err := h.depositRepo.Create(ctx, depositReq)
	if err != nil {
		log.Printf("[DepositHandler] Failed to create demo deposit: %v", err)
	}

	// Credit user wallet directly
	newBalance, err := h.creditUserWallet(ctx, req.UserID, req.WalletID, req.Amount, req.Currency, transactionID)
	if err != nil {
		log.Printf("[DepositHandler] Failed to credit wallet in demo mode: %v", err)
		_ = h.depositRepo.MarkFailed(ctx, transactionID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to credit wallet"})
		return
	}

	// Mark as completed
	_ = h.depositRepo.MarkWalletCredited(ctx, transactionID, req.WalletID)

	log.Printf("[DepositHandler] ✅ Demo deposit completed: %s | Amount: %.2f %s", transactionID, req.Amount, req.Currency)

	c.JSON(http.StatusOK, InitiateDepositResponse{
		TransactionID: transactionID,
		Status:        "completed",
		Provider:      "demo",
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       "Dépôt effectué avec succès (mode démo)",
		NewBalance:    &newBalance,
	})
}

// GetDepositStatus returns the status of a deposit
// GET /api/v1/deposits/:id/status
func (h *DepositHandler) GetDepositStatus(c *gin.Context) {
	transactionID := c.Param("id")
	ctx := c.Request.Context()

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deposit"})
		return
	}

	if deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id":  deposit.ID,
		"status":          deposit.Status,
		"status_message":  deposit.StatusMessage,
		"amount":          deposit.Amount,
		"currency":        deposit.Currency,
		"fee":             deposit.Fee,
		"provider":        deposit.ProviderCode,
		"wallet_credited": deposit.WalletCredited,
		"completed_at":    deposit.CompletedAt,
		"failed_at":       deposit.FailedAt,
		"cancelled_at":    deposit.CancelledAt,
		"created_at":      deposit.CreatedAt,
	})
}

// CancelDeposit cancels a pending deposit
// POST /api/v1/deposits/:id/cancel
func (h *DepositHandler) CancelDeposit(c *gin.Context) {
	transactionID := c.Param("id")
	ctx := c.Request.Context()

	// Get current status
	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit not found"})
		return
	}

	// Can only cancel pending transactions
	if deposit.Status != repository.DepositStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Cannot cancel deposit",
			"reason": fmt.Sprintf("Deposit is already %s", deposit.Status),
		})
		return
	}

	// Mark as cancelled
	err = h.depositRepo.MarkCancelled(ctx, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel deposit"})
		return
	}

	log.Printf("[DepositHandler] 🚫 Deposit cancelled: %s", transactionID)

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": transactionID,
		"status":         "cancelled",
		"message":        "Dépôt annulé avec succès",
	})
}

// GetUserDeposits returns all deposits for a user
// GET /api/v1/deposits/user/:user_id
func (h *DepositHandler) GetUserDeposits(c *gin.Context) {
	userID := c.Param("user_id")
	limit := 20
	offset := 0

	ctx := c.Request.Context()

	deposits, err := h.depositRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deposits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deposits": deposits,
		"count":    len(deposits),
	})
}

// ==================== WEBHOOK HANDLERS ====================

// HandleWebhook routes webhooks to the appropriate provider handler
// POST /api/v1/deposits/webhook/:provider
func (h *DepositHandler) HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")
	ctx := c.Request.Context()

	// Read raw body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[Webhook] Failed to read body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	// Restore body for further processing
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	log.Printf("[Webhook] Received from %s: %s", provider, string(body))

	switch provider {
	case "flutterwave":
		h.handleFlutterwaveWebhook(ctx, c, body)
	case "stripe":
		h.handleStripeWebhook(ctx, c, body)
	case "paystack":
		h.handlePaystackWebhook(ctx, c, body)
	case "cinetpay":
		h.handleCinetPayWebhook(ctx, c, body)
	case "lygos":
		h.handleLygosWebhook(ctx, c, body)
	case "orange_money":
		h.handleOrangeMoneyWebhook(ctx, c, body)
	case "mtn_momo":
		h.handleMTNMoMoWebhook(ctx, c, body)
	case "wave", "wave_ci":
		h.handleWaveWebhook(ctx, c, body)
	case "paypal":
		h.handlePayPalWebhook(ctx, c, body)
	case "fedapay":
		h.handleFedaPayWebhook(ctx, c, body)
	case "moov_money":
		h.handleMoovMoneyWebhook(ctx, c, body)
	case "yellowcard":
		h.handleYellowCardWebhook(ctx, c, body)
	default:
		log.Printf("[Webhook] Unknown provider: %s", provider)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unknown provider: %s", provider)})
		return
	}
}

// handleFlutterwaveWebhook handles Flutterwave webhooks
func (h *DepositHandler) handleFlutterwaveWebhook(ctx context.Context, c *gin.Context, body []byte) {
	// Verify signature
	signature := c.GetHeader("verif-hash")
	secret := h.webhookSecrets["flutterwave"]

	if secret != "" && signature != secret {
		log.Printf("[Flutterwave Webhook] Invalid signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var payload struct {
		Event string `json:"event"`
		Data  struct {
			ID            int             `json:"id"`
			TxRef         string          `json:"tx_ref"`
			FlwRef        string          `json:"flw_ref"`
			Amount        float64         `json:"amount"`
			Currency      string          `json:"currency"`
			ChargedAmount float64         `json:"charged_amount"`
			AppFee        float64         `json:"app_fee"`
			Status        string          `json:"status"`
			Customer      json.RawMessage `json:"customer"`
			Meta          json.RawMessage `json:"meta"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[Flutterwave Webhook] Parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Data.TxRef
	providerRef := fmt.Sprintf("%d", payload.Data.ID)

	// Save webhook data
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	// Get deposit record
	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		log.Printf("[Flutterwave Webhook] Deposit not found: %s", transactionID)
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "deposit not found"})
		return
	}

	// Update provider reference
	_ = h.depositRepo.UpdateProviderReference(ctx, transactionID, providerRef, "")

	// Process based on status
	switch payload.Data.Status {
	case "successful":
		h.completeDeposit(ctx, deposit, payload.Data.AppFee)
	case "failed":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, "Payment failed at Flutterwave")
	}

	log.Printf("[Flutterwave Webhook] Processed: %s -> %s", transactionID, payload.Data.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleStripeWebhook handles Stripe webhooks
func (h *DepositHandler) handleStripeWebhook(ctx context.Context, c *gin.Context, body []byte) {
	signature := c.GetHeader("Stripe-Signature")
	secret := h.webhookSecrets["stripe"]

	// Verify signature (simplified - in production use Stripe SDK)
	if secret != "" && !verifyStripeSignature(body, signature, secret) {
		log.Printf("[Stripe Webhook] Invalid signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var payload struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID            string            `json:"id"`
				PaymentIntent string            `json:"payment_intent"`
				AmountTotal   int               `json:"amount_total"`
				Currency      string            `json:"currency"`
				PaymentStatus string            `json:"payment_status"`
				Metadata      map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Data.Object.Metadata["reference_id"]
	if transactionID == "" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch payload.Type {
	case "checkout.session.completed", "payment_intent.succeeded":
		if payload.Data.Object.PaymentStatus == "paid" {
			h.completeDeposit(ctx, deposit, 0)
		}
	case "payment_intent.payment_failed":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, "Payment failed at Stripe")
	}

	log.Printf("[Stripe Webhook] Processed: %s -> %s", transactionID, payload.Type)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handlePaystackWebhook handles Paystack webhooks
func (h *DepositHandler) handlePaystackWebhook(ctx context.Context, c *gin.Context, body []byte) {
	signature := c.GetHeader("X-Paystack-Signature")
	secret := h.webhookSecrets["paystack"]

	if secret != "" && !verifyHMACSHA512(body, signature, secret) {
		log.Printf("[Paystack Webhook] Invalid signature")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var payload struct {
		Event string `json:"event"`
		Data  struct {
			Reference string            `json:"reference"`
			Amount    int               `json:"amount"`
			Currency  string            `json:"currency"`
			Status    string            `json:"status"`
			Fees      int               `json:"fees"`
			Metadata  map[string]string `json:"metadata"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Data.Reference
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch payload.Event {
	case "charge.success":
		fee := float64(payload.Data.Fees) / 100
		h.completeDeposit(ctx, deposit, fee)
	case "charge.failed":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, "Payment failed at Paystack")
	}

	log.Printf("[Paystack Webhook] Processed: %s -> %s", transactionID, payload.Event)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleCinetPayWebhook handles CinetPay webhooks
func (h *DepositHandler) handleCinetPayWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		CpmTransID     string `json:"cpm_trans_id"`
		CpmSiteID      string `json:"cpm_site_id"`
		CpmAmount      string `json:"cpm_amount"`
		CpmCurrency    string `json:"cpm_currency"`
		CpmPaymentDate string `json:"cpm_payment_date"`
		CpmPaymentTime string `json:"cpm_payment_time"`
		CpmResult      string `json:"cpm_result"`
		CpmTransStatus string `json:"cpm_trans_status"`
		CpmDesignation string `json:"cpm_designation"`
		CpmCustom      string `json:"cpm_custom"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		// Try form data
		payload.CpmTransID = c.PostForm("cpm_trans_id")
		payload.CpmResult = c.PostForm("cpm_result")
		payload.CpmTransStatus = c.PostForm("cpm_trans_status")
		payload.CpmCustom = c.PostForm("cpm_custom")
	}

	transactionID := payload.CpmCustom
	if transactionID == "" {
		transactionID = payload.CpmTransID
	}

	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch payload.CpmResult {
	case "00", "SUCCESS":
		h.completeDeposit(ctx, deposit, 0)
	default:
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("CinetPay error: %s", payload.CpmResult))
	}

	log.Printf("[CinetPay Webhook] Processed: %s -> %s", transactionID, payload.CpmResult)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleLygosWebhook handles Lygos webhooks
func (h *DepositHandler) handleLygosWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		TransactionID string            `json:"transaction_id"`
		Reference     string            `json:"reference"`
		Status        string            `json:"status"`
		Amount        float64           `json:"amount"`
		Currency      string            `json:"currency"`
		Fee           float64           `json:"fee"`
		Metadata      map[string]string `json:"metadata"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Reference
	if transactionID == "" {
		transactionID = payload.Metadata["reference_id"]
	}

	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToLower(payload.Status) {
	case "successful", "completed", "success":
		h.completeDeposit(ctx, deposit, payload.Fee)
	case "failed", "error":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, "Payment failed at Lygos")
	}

	log.Printf("[Lygos Webhook] Processed: %s -> %s", transactionID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleOrangeMoneyWebhook handles Orange Money webhooks
func (h *DepositHandler) handleOrangeMoneyWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		Status     string `json:"status"`
		NotifToken string `json:"notif_token"`
		TxnID      string `json:"txnid"`
		Message    string `json:"message"`
		Amount     string `json:"amount"`
		Currency   string `json:"currency"`
		OrderID    string `json:"order_id"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.OrderID
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToUpper(payload.Status) {
	case "SUCCESS", "SUCCESSFUL":
		h.completeDeposit(ctx, deposit, 0)
	case "FAILED", "CANCELLED":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Orange Money: %s", payload.Message))
	}

	log.Printf("[Orange Money Webhook] Processed: %s -> %s", transactionID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleMTNMoMoWebhook handles MTN Mobile Money webhooks
func (h *DepositHandler) handleMTNMoMoWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		ReferenceID      string `json:"referenceId"`
		ExternalID       string `json:"externalId"`
		FinancialTransID string `json:"financialTransactionId"`
		Status           string `json:"status"`
		Amount           string `json:"amount"`
		Currency         string `json:"currency"`
		Reason           string `json:"reason"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.ExternalID
	if transactionID == "" {
		transactionID = payload.ReferenceID
	}

	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToUpper(payload.Status) {
	case "SUCCESSFUL":
		h.completeDeposit(ctx, deposit, 0)
	case "FAILED":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("MTN MoMo: %s", payload.Reason))
	}

	log.Printf("[MTN MoMo Webhook] Processed: %s -> %s", transactionID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleWaveWebhook handles Wave webhooks
func (h *DepositHandler) handleWaveWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		ID           string `json:"id"`
		Type         string `json:"type"`
		ClientRef    string `json:"client_reference"`
		Amount       string `json:"amount"`
		Currency     string `json:"currency"`
		Status       string `json:"status"`
		ErrorMessage string `json:"error_message"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.ClientRef
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToLower(payload.Status) {
	case "succeeded", "completed":
		h.completeDeposit(ctx, deposit, 0)
	case "failed":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Wave: %s", payload.ErrorMessage))
	}

	log.Printf("[Wave Webhook] Processed: %s -> %s", transactionID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handlePayPalWebhook handles PayPal webhooks
func (h *DepositHandler) handlePayPalWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID            string `json:"id"`
			Status        string `json:"status"`
			CustomID      string `json:"custom_id"`
			PurchaseUnits []struct {
				ReferenceID string `json:"reference_id"`
				Amount      struct {
					Value    string `json:"value"`
					Currency string `json:"currency_code"`
				} `json:"amount"`
			} `json:"purchase_units"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Resource.CustomID
	if transactionID == "" && len(payload.Resource.PurchaseUnits) > 0 {
		transactionID = payload.Resource.PurchaseUnits[0].ReferenceID
	}

	if transactionID == "" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch payload.EventType {
	case "CHECKOUT.ORDER.APPROVED", "PAYMENT.CAPTURE.COMPLETED":
		h.completeDeposit(ctx, deposit, 0)
	case "PAYMENT.CAPTURE.DENIED", "CHECKOUT.ORDER.CANCELLED":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, "PayPal payment denied or cancelled")
	}

	log.Printf("[PayPal Webhook] Processed: %s -> %s", transactionID, payload.EventType)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleFedaPayWebhook handles FedaPay webhooks
func (h *DepositHandler) handleFedaPayWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		Event string `json:"event"`
		Data  struct {
			ID          int    `json:"id"`
			Reference   string `json:"reference"`
			Description string `json:"description"`
			Status      string `json:"status"`
			Amount      int    `json:"amount"`
			Currency    struct {
				ISO string `json:"iso"`
			} `json:"currency"`
		} `json:"entity"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Data.Reference
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch payload.Event {
	case "transaction.approved", "transaction.completed":
		h.completeDeposit(ctx, deposit, 0)
	case "transaction.declined", "transaction.canceled":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("FedaPay: %s", payload.Data.Status))
	}

	log.Printf("[FedaPay Webhook] Processed: %s -> %s", transactionID, payload.Event)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleMoovMoneyWebhook handles Moov Money webhooks
func (h *DepositHandler) handleMoovMoneyWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		Status        string `json:"status"`
		TransactionID string `json:"transaction-id"`
		Message       string `json:"message"`
		ReferenceID   string `json:"referenceid"`
		Amount        string `json:"amount"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.ReferenceID
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToLower(payload.Status) {
	case "0", "success", "successful":
		h.completeDeposit(ctx, deposit, 0)
	default:
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Moov Money: %s", payload.Message))
	}

	log.Printf("[Moov Money Webhook] Processed: %s -> %s", transactionID, payload.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// handleYellowCardWebhook handles YellowCard webhooks
func (h *DepositHandler) handleYellowCardWebhook(ctx context.Context, c *gin.Context, body []byte) {
	var payload struct {
		Event string `json:"event"`
		Data  struct {
			ID            string  `json:"id"`
			SequenceID    string  `json:"sequenceId"`
			Status        string  `json:"status"`
			Amount        float64 `json:"amount"`
			Currency      string  `json:"currency"`
			ChannelID     string  `json:"channelId"`
			Reason        string  `json:"reason"`
			CustomerEmail string  `json:"customerEmail"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	transactionID := payload.Data.SequenceID
	_ = h.depositRepo.SaveWebhookData(ctx, transactionID, payload, true)

	deposit, err := h.depositRepo.GetByID(ctx, transactionID)
	if err != nil || deposit == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	switch strings.ToLower(payload.Data.Status) {
	case "completed", "successful":
		h.completeDeposit(ctx, deposit, 0)
	case "failed", "cancelled":
		_ = h.depositRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("YellowCard: %s", payload.Data.Reason))
	}

	log.Printf("[YellowCard Webhook] Processed: %s -> %s", transactionID, payload.Data.Status)
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// ==================== HELPER FUNCTIONS ====================

// completeDeposit completes a deposit by crediting the user's wallet
func (h *DepositHandler) completeDeposit(ctx context.Context, deposit *repository.DepositTransaction, fee float64) {
	// Get user wallet ID from metadata
	var metadata map[string]interface{}
	if deposit.Metadata != nil {
		_ = json.Unmarshal(deposit.Metadata, &metadata)
	}

	userWalletID := ""
	if wid, ok := metadata["user_wallet_id"].(string); ok {
		userWalletID = wid
	}

	// Credit user wallet
	_, err := h.creditUserWallet(ctx, deposit.UserID, userWalletID, deposit.Amount-fee, deposit.Currency, deposit.ID)
	if err != nil {
		log.Printf("[DepositHandler] Failed to credit wallet: %v", err)
		_ = h.depositRepo.MarkFailed(ctx, deposit.ID, fmt.Sprintf("Failed to credit wallet: %v", err))
		return
	}

	// Mark deposit as completed
	_ = h.depositRepo.MarkWalletCredited(ctx, deposit.ID, userWalletID)

	log.Printf("[DepositHandler] ✅ Deposit completed: %s | Amount: %.2f %s | User: %s",
		deposit.ID, deposit.Amount, deposit.Currency, deposit.UserID)
}

// creditUserWallet credits the user's wallet via wallet-service
func (h *DepositHandler) creditUserWallet(ctx context.Context, userID, walletID string, amount float64, currency, transactionRef string) (float64, error) {
	// Call wallet-service to credit the user's wallet
	reqBody := map[string]interface{}{
		"user_id":         userID,
		"wallet_id":       walletID,
		"amount":          amount,
		"currency":        currency,
		"transaction_ref": transactionRef,
		"type":            "deposit",
		"description":     "Dépôt via agrégateur de paiement",
	}

	body, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/api/v1/wallets/credit", h.walletServiceURL)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success    bool    `json:"success"`
		NewBalance float64 `json:"new_balance"`
		Error      string  `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode wallet response: %w", err)
	}

	if !result.Success {
		return 0, fmt.Errorf("wallet credit failed: %s", result.Error)
	}

	return result.NewBalance, nil
}

// ==================== SIGNATURE VERIFICATION ====================

// verifyHMACSHA512 verifies HMAC-SHA512 signature (used by Paystack)
func verifyHMACSHA512(body []byte, signature, secret string) bool {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// verifyHMACSHA256 verifies HMAC-SHA256 signature
func verifyHMACSHA256(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// verifyStripeSignature verifies Stripe webhook signature (simplified)
func verifyStripeSignature(body []byte, signature, secret string) bool {
	// In production, use the official Stripe Go library for signature verification
	// This is a simplified check
	if signature == "" || secret == "" {
		return false
	}
	// Parse signature header (format: t=timestamp,v1=signature)
	parts := strings.Split(signature, ",")
	var timestamp, sig string
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			if kv[0] == "t" {
				timestamp = kv[1]
			} else if kv[0] == "v1" {
				sig = kv[1]
			}
		}
	}
	if timestamp == "" || sig == "" {
		return false
	}
	// Compute expected signature
	payload := fmt.Sprintf("%s.%s", timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sig), []byte(expected))
}

// ==================== BACKGROUND JOBS ====================

// ExpireOldTransactions marks old pending transactions as expired
// This should be called periodically (e.g., every 5 minutes)
func (h *DepositHandler) ExpireOldTransactions(ctx context.Context) (int, error) {
	count, err := h.depositRepo.ExpirePendingTransactions(ctx)
	if err != nil {
		log.Printf("[DepositHandler] Failed to expire transactions: %v", err)
		return 0, err
	}
	if count > 0 {
		log.Printf("[DepositHandler] ⏰ Expired %d pending transactions", count)
	}
	return count, nil
}

// RegisterRoutes registers all deposit routes
func (h *DepositHandler) RegisterRoutes(router *gin.RouterGroup) {
	deposits := router.Group("/deposits")
	{
		deposits.POST("/initiate", h.InitiateDeposit)
		deposits.GET("/:id/status", h.GetDepositStatus)
		deposits.POST("/:id/cancel", h.CancelDeposit)
		deposits.GET("/user/:user_id", h.GetUserDeposits)

		deposits.POST("/paypal/create-order", h.CreatePayPalOrder)
		deposits.POST("/paypal/capture", h.CapturePayPalOrder)

		// Webhook endpoints
		deposits.POST("/webhook/:provider", h.HandleWebhook)
	}
}
