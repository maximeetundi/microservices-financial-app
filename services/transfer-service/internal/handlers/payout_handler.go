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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/providers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/service"
)

// PayoutHandler handles external withdrawal/payout operations
type PayoutHandler struct {
	payoutRepo       *repository.PayoutRepository
	instanceRepo     *repository.AggregatorInstanceRepository
	providerLoader   *providers.InstanceBasedProviderLoader
	fundMovement     *service.WithdrawalFundMovementService
	walletServiceURL string
	webhookSecrets   map[string]string
}

// NewPayoutHandler creates a new payout handler
func NewPayoutHandler(
	payoutRepo *repository.PayoutRepository,
	instanceRepo *repository.AggregatorInstanceRepository,
	providerLoader *providers.InstanceBasedProviderLoader,
	fundMovement *service.WithdrawalFundMovementService,
	walletServiceURL string,
	webhookSecrets map[string]string,
) *PayoutHandler {
	return &PayoutHandler{
		payoutRepo:       payoutRepo,
		instanceRepo:     instanceRepo,
		providerLoader:   providerLoader,
		fundMovement:     fundMovement,
		walletServiceURL: walletServiceURL,
		webhookSecrets:   webhookSecrets,
	}
}

// ==================== REQUEST/RESPONSE TYPES ====================

// InitiatePayoutRequest represents a payout initiation request
type InitiatePayoutRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	WalletID       string  `json:"wallet_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required"`
	Provider       string  `json:"provider" binding:"required"` // flutterwave, stripe, paypal, mtn_momo, orange_money, etc.
	Country        string  `json:"country" binding:"required"`
	IsTestMode     *bool   `json:"is_test_mode,omitempty"`
	PayoutMethod   string  `json:"payout_method" binding:"required"` // mobile_money, bank_transfer, paypal, card
	RecipientName  string  `json:"recipient_name" binding:"required"`
	RecipientEmail string  `json:"recipient_email,omitempty"`
	RecipientPhone string  `json:"recipient_phone,omitempty"`
	BankCode       string  `json:"bank_code,omitempty"`
	BankName       string  `json:"bank_name,omitempty"`
	AccountNumber  string  `json:"account_number,omitempty"`
	IBAN           string  `json:"iban,omitempty"`
	SwiftCode      string  `json:"swift_code,omitempty"`
	RoutingNumber  string  `json:"routing_number,omitempty"`
	MobileOperator string  `json:"mobile_operator,omitempty"`
	MobileNumber   string  `json:"mobile_number,omitempty"`
	PayPalEmail    string  `json:"paypal_email,omitempty"`
	Narration      string  `json:"narration,omitempty"`
	PIN            string  `json:"pin" binding:"required"` // User PIN for security
}

// InitiatePayoutResponse represents the response for payout initiation
type InitiatePayoutResponse struct {
	TransactionID        string   `json:"transaction_id"`
	Status               string   `json:"status"`
	Provider             string   `json:"provider"`
	PayoutMethod         string   `json:"payout_method"`
	Amount               float64  `json:"amount"`
	Currency             string   `json:"currency"`
	Fee                  float64  `json:"fee"`
	AmountReceived       float64  `json:"amount_received"`
	RecipientName        string   `json:"recipient_name"`
	RecipientAccount     string   `json:"recipient_account"` // Masked account/number
	Message              string   `json:"message,omitempty"`
	AggregatorInstanceID string   `json:"aggregator_instance_id,omitempty"`
	HotWalletID          string   `json:"hot_wallet_id,omitempty"`
	EstimatedDelivery    string   `json:"estimated_delivery,omitempty"`
	NewBalance           *float64 `json:"new_balance,omitempty"`
}

// PayoutQuoteRequest represents a quote request
type PayoutQuoteRequest struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
	Provider     string  `json:"provider" binding:"required"`
	Country      string  `json:"country" binding:"required"`
	IsTestMode   *bool   `json:"is_test_mode,omitempty"`
	PayoutMethod string  `json:"payout_method" binding:"required"`
}

// PayoutQuoteResponse represents a quote response
type PayoutQuoteResponse struct {
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Fee              float64 `json:"fee"`
	FeeType          string  `json:"fee_type"` // flat, percentage, hybrid
	AmountReceived   float64 `json:"amount_received"`
	ExchangeRate     float64 `json:"exchange_rate,omitempty"`
	ReceivedCurrency string  `json:"received_currency,omitempty"`
	EstimatedMinutes int     `json:"estimated_minutes"`
	Provider         string  `json:"provider"`
	PayoutMethod     string  `json:"payout_method"`
	MinAmount        float64 `json:"min_amount"`
	MaxAmount        float64 `json:"max_amount"`
}

// PayoutStatusResponse represents payout status
type PayoutStatusResponse struct {
	TransactionID     string  `json:"transaction_id"`
	Status            string  `json:"status"` // pending, processing, completed, failed, cancelled
	ProviderReference string  `json:"provider_reference,omitempty"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Fee               float64 `json:"fee"`
	AmountReceived    float64 `json:"amount_received"`
	RecipientName     string  `json:"recipient_name"`
	RecipientAccount  string  `json:"recipient_account"`
	PayoutMethod      string  `json:"payout_method"`
	Provider          string  `json:"provider"`
	StatusMessage     string  `json:"status_message,omitempty"`
	CreatedAt         string  `json:"created_at"`
	CompletedAt       string  `json:"completed_at,omitempty"`
}

const DefaultPayoutTimeout = 24 * time.Hour

// ==================== MAIN ENDPOINTS ====================

// GetPayoutQuote returns fee estimation for a payout
// POST /api/v1/payouts/quote
func (h *PayoutHandler) GetPayoutQuote(c *gin.Context) {
	var req PayoutQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Get best instance for this provider
	instance, err := h.providerLoader.CredentialsClient().GetBestInstanceWithCredentials(req.Provider, req.Country, req.Amount, req.Currency, "withdraw", req.IsTestMode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Provider not available for withdrawals",
			"details": err.Error(),
		})
		return
	}

	// Load provider from this instance (withdraw uses same provider type)
	_, err = h.providerLoader.LoadProviderFromInstance(ctx, instance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to load provider",
			"details": err.Error(),
		})
		return
	}

	// Calculate fees based on provider settings
	var fee float64
	var feeType string
	var estimatedMinutes int

	switch req.PayoutMethod {
	case "mobile_money":
		fee = req.Amount * 0.015 // 1.5% for mobile money
		feeType = "percentage"
		estimatedMinutes = 5
	case "bank_transfer":
		fee = req.Amount * 0.02 // 2% for bank transfer
		if fee < 5 {
			fee = 5 // Minimum fee
		}
		feeType = "hybrid"
		estimatedMinutes = 1440 // 24 hours
	case "paypal":
		fee = req.Amount*0.029 + 0.30 // PayPal standard fees
		feeType = "hybrid"
		estimatedMinutes = 60
	default:
		fee = req.Amount * 0.02
		feeType = "percentage"
		estimatedMinutes = 30
	}

	// Get min/max from instance settings
	minAmount := 100.0
	if instance.MinBalance != nil {
		minAmount = *instance.MinBalance
	}

	maxAmount := 5000000.0
	if instance.MaxBalance != nil {
		maxAmount = *instance.MaxBalance
	}

	c.JSON(http.StatusOK, PayoutQuoteResponse{
		Amount:           req.Amount,
		Currency:         req.Currency,
		Fee:              fee,
		FeeType:          feeType,
		AmountReceived:   req.Amount - fee,
		ExchangeRate:     1.0, // Same currency
		ReceivedCurrency: req.Currency,
		EstimatedMinutes: estimatedMinutes,
		Provider:         req.Provider,
		PayoutMethod:     req.PayoutMethod,
		MinAmount:        minAmount,
		MaxAmount:        maxAmount,
	})
}

// InitiatePayout initiates an external payout/withdrawal
// POST /api/v1/payouts/initiate
func (h *PayoutHandler) InitiatePayout(c *gin.Context) {
	var req InitiatePayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	transactionID := fmt.Sprintf("pay_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// 1. Verify user PIN first
	if err := h.verifyUserPIN(ctx, req.UserID, req.PIN, c.GetHeader("Authorization")); err != nil {
		log.Printf("[PayoutHandler] PIN verification failed for user %s: %v", req.UserID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
		return
	}

	// 2. Get best instance for this provider (withdraw)
	instance, err := h.providerLoader.CredentialsClient().GetBestInstanceWithCredentials(req.Provider, req.Country, req.Amount, req.Currency, "withdraw", req.IsTestMode)
	if err != nil {
		log.Printf("[PayoutHandler] Provider not available: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Provider not available",
			"details": err.Error(),
		})
		return
	}
	provider, err := h.providerLoader.LoadProviderFromInstance(ctx, instance)
	if err != nil {
		log.Printf("[PayoutHandler] Failed to load provider: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to load provider",
			"details": err.Error(),
		})
		return
	}

	// Check if withdrawal is enabled for this provider
	if !instance.WithdrawEnabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Withdrawals are not enabled for this provider",
		})
		return
	}

	// 3. Select best hot wallet for payout
	walletID, err := h.fundMovement.SelectBestWalletForWithdrawal(ctx, instance.ID, req.Currency, req.Amount)
	if err != nil {
		log.Printf("[PayoutHandler] No available wallet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No available wallet for this currency",
			"details": err.Error(),
		})
		return
	}

	// 4. Calculate fees
	fee := h.calculatePayoutFee(req.Amount, req.PayoutMethod, instance)
	amountReceived := req.Amount - fee

	// 5. Check user balance first
	userBalance, err := h.getUserWalletBalance(ctx, req.UserID, req.WalletID)
	if err != nil {
		log.Printf("[PayoutHandler] Failed to get user balance: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to verify wallet balance"})
		return
	}

	if userBalance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Insufficient balance. Available: %.2f %s, Required: %.2f %s",
				userBalance, req.Currency, req.Amount, req.Currency),
		})
		return
	}

	// 6. Determine recipient account for display
	recipientAccount := h.getRecipientAccountDisplay(req)

	// 7. Create payout transaction record
	payoutReq := &repository.CreatePayoutRequest{
		ID:                   transactionID,
		UserID:               req.UserID,
		WalletID:             req.WalletID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Fee:                  fee,
		AmountReceived:       amountReceived,
		ProviderCode:         req.Provider,
		PayoutMethod:         req.PayoutMethod,
		AggregatorInstanceID: instance.ID,
		HotWalletID:          walletID,
		RecipientName:        req.RecipientName,
		RecipientEmail:       req.RecipientEmail,
		RecipientPhone:       req.RecipientPhone,
		BankCode:             req.BankCode,
		BankName:             req.BankName,
		AccountNumber:        req.AccountNumber,
		IBAN:                 req.IBAN,
		MobileOperator:       req.MobileOperator,
		MobileNumber:         req.MobileNumber,
		PayPalEmail:          req.PayPalEmail,
		Narration:            req.Narration,
		Country:              req.Country,
		ExpiresAt:            time.Now().Add(DefaultPayoutTimeout),
		IPAddress:            c.ClientIP(),
		UserAgent:            c.Request.UserAgent(),
	}

	if _, err = h.payoutRepo.Create(ctx, payoutReq); err != nil {
		log.Printf("[PayoutHandler] Failed to create payout record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payout"})
		return
	}

	// 8. Process fund movement: User Wallet -> Hot Wallet
	err = h.fundMovement.ProcessWithdrawal(
		ctx,
		req.UserID,
		walletID,
		req.Amount,
		req.Currency,
		transactionID,
		req.Provider,
		recipientAccount,
	)
	if err != nil {
		log.Printf("[PayoutHandler] Fund movement failed: %v", err)
		_ = h.payoutRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Fund movement failed: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process withdrawal",
			"details": err.Error(),
		})
		return
	}

	// 9. Initiate payout with provider
	payoutRequest := h.buildProviderPayoutRequest(req, transactionID)
	response, err := h.initiateProviderPayout(ctx, provider, payoutRequest, instance)

	if err != nil {
		log.Printf("[PayoutHandler] Provider payout initiation failed: %v", err)
		// Reverse the fund movement
		_ = h.fundMovement.ConfirmWithdrawalPayout(ctx, transactionID, "failed", "")
		_ = h.payoutRepo.MarkFailed(ctx, transactionID, fmt.Sprintf("Provider error: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initiate payout with provider",
			"details": err.Error(),
		})
		return
	}

	// 10. Update transaction with provider reference
	_ = h.payoutRepo.UpdateProviderReference(ctx, transactionID, response.ProviderReference)
	_ = h.payoutRepo.UpdateStatus(ctx, transactionID, "processing", "Payout initiated with provider")

	// 11. Record usage in instance
	_ = h.providerLoader.RecordTransactionUsage(
		ctx,
		instance.ID,
		transactionID,
		req.Amount,
		req.Currency,
		"processing",
		response.ProviderReference,
	)

	// Get new balance
	newBalance, _ := h.getUserWalletBalance(ctx, req.UserID, req.WalletID)

	log.Printf("[PayoutHandler] ✅ Payout initiated: %s | Provider: %s | Amount: %.2f %s | Recipient: %s",
		transactionID, req.Provider, req.Amount, req.Currency, recipientAccount)

	c.JSON(http.StatusOK, InitiatePayoutResponse{
		TransactionID:        transactionID,
		Status:               "processing",
		Provider:             req.Provider,
		PayoutMethod:         req.PayoutMethod,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Fee:                  fee,
		AmountReceived:       amountReceived,
		RecipientName:        req.RecipientName,
		RecipientAccount:     recipientAccount,
		Message:              response.Message,
		AggregatorInstanceID: instance.ID,
		HotWalletID:          walletID,
		EstimatedDelivery:    response.EstimatedDelivery.Format(time.RFC3339),
		NewBalance:           &newBalance,
	})
}

// GetPayoutStatus returns the status of a payout
// GET /api/v1/payouts/:id/status
func (h *PayoutHandler) GetPayoutStatus(c *gin.Context) {
	transactionID := c.Param("id")

	ctx := c.Request.Context()
	payout, err := h.payoutRepo.GetByID(ctx, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
		return
	}

	recipientAccount := h.getRecipientAccountFromPayout(payout)

	c.JSON(http.StatusOK, PayoutStatusResponse{
		TransactionID:     payout.ID,
		Status:            payout.Status,
		ProviderReference: payout.ProviderReference,
		Amount:            payout.Amount,
		Currency:          payout.Currency,
		Fee:               payout.Fee,
		AmountReceived:    payout.AmountReceived,
		RecipientName:     payout.RecipientName,
		RecipientAccount:  recipientAccount,
		PayoutMethod:      payout.PayoutMethod,
		Provider:          payout.ProviderCode,
		StatusMessage:     payout.StatusMessage,
		CreatedAt:         payout.CreatedAt.Format(time.RFC3339),
		CompletedAt:       formatTime(payout.CompletedAt),
	})
}

// CancelPayout cancels a pending payout
// POST /api/v1/payouts/:id/cancel
func (h *PayoutHandler) CancelPayout(c *gin.Context) {
	transactionID := c.Param("id")

	ctx := c.Request.Context()
	payout, err := h.payoutRepo.GetByID(ctx, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payout not found"})
		return
	}

	// Can only cancel pending payouts
	if payout.Status != "pending" && payout.Status != "processing" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot cancel payout with status: %s", payout.Status),
		})
		return
	}

	// Reverse the fund movement
	err = h.fundMovement.ConfirmWithdrawalPayout(ctx, transactionID, "failed", "")
	if err != nil {
		log.Printf("[PayoutHandler] Failed to reverse fund movement: %v", err)
	}

	// Mark as cancelled
	_ = h.payoutRepo.UpdateStatus(ctx, transactionID, "cancelled", "Cancelled by user")

	log.Printf("[PayoutHandler] ❌ Payout cancelled: %s", transactionID)

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"transaction_id": transactionID,
		"status":         "cancelled",
		"message":        "Payout cancelled and funds returned to wallet",
	})
}

// GetUserPayouts returns payout history for a user
// GET /api/v1/payouts/user/:user_id
func (h *PayoutHandler) GetUserPayouts(c *gin.Context) {
	userID := c.Param("user_id")

	ctx := c.Request.Context()
	payouts, err := h.payoutRepo.GetByUserID(ctx, userID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payouts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payouts": payouts,
		"count":   len(payouts),
	})
}

// GetBanks returns available banks for a country
// GET /api/v1/payouts/banks?country=XX
func (h *PayoutHandler) GetBanks(c *gin.Context) {
	country := c.Query("country")
	if country == "" {
		country = "CI"
	}

	// TODO: Fetch from provider API or database
	banks := []map[string]string{
		{"code": "012", "name": "SGBCI (Société Générale)"},
		{"code": "030", "name": "BICICI (BNP Paribas)"},
		{"code": "062", "name": "BIAO-CI"},
		{"code": "064", "name": "SIB (Attijariwafa Bank)"},
		{"code": "065", "name": "NSIA Banque"},
		{"code": "066", "name": "ECOBANK"},
		{"code": "067", "name": "BOA (Bank of Africa)"},
		{"code": "068", "name": "BACI (Banque Atlantique)"},
		{"code": "069", "name": "CORIS Bank"},
		{"code": "074", "name": "UBA (United Bank for Africa)"},
		{"code": "077", "name": "Orabank"},
		{"code": "079", "name": "Diamond Bank"},
		{"code": "088", "name": "Bridge Bank"},
		{"code": "134", "name": "GTP Bank"},
	}

	c.JSON(http.StatusOK, gin.H{
		"banks":   banks,
		"country": country,
	})
}

// GetMobileOperators returns available mobile money operators for a country
// GET /api/v1/payouts/mobile-operators?country=XX
func (h *PayoutHandler) GetMobileOperators(c *gin.Context) {
	country := c.Query("country")
	if country == "" {
		country = "CI"
	}

	operators := map[string][]map[string]string{
		"CI": {
			{"code": "MTN", "name": "MTN Mobile Money", "prefix": "05,06"},
			{"code": "ORANGE", "name": "Orange Money", "prefix": "07,08"},
			{"code": "MOOV", "name": "Moov Money", "prefix": "01"},
			{"code": "WAVE", "name": "Wave", "prefix": ""},
		},
		"SN": {
			{"code": "ORANGE", "name": "Orange Money", "prefix": "77,78"},
			{"code": "FREE", "name": "Free Money", "prefix": "76"},
			{"code": "WAVE", "name": "Wave", "prefix": ""},
		},
		"NG": {
			{"code": "MTN", "name": "MTN MoMo", "prefix": "080,081"},
			{"code": "AIRTEL", "name": "Airtel Money", "prefix": "090,091"},
			{"code": "GLO", "name": "Glo Mobile Money", "prefix": "070,071"},
		},
		"GH": {
			{"code": "MTN", "name": "MTN Mobile Money", "prefix": "024,054,055"},
			{"code": "VODAFONE", "name": "Vodafone Cash", "prefix": "020,050"},
			{"code": "AIRTELTIGO", "name": "AirtelTigo Money", "prefix": "027,057,026,056"},
		},
		"KE": {
			{"code": "MPESA", "name": "M-Pesa", "prefix": "07"},
			{"code": "AIRTEL", "name": "Airtel Money", "prefix": "07"},
		},
	}

	ops := operators[country]
	if ops == nil {
		ops = []map[string]string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"operators": ops,
		"country":   country,
	})
}

// ==================== WEBHOOK HANDLERS ====================

// HandleWebhook handles payout webhooks from providers
// POST /api/v1/payouts/webhook/:provider
func (h *PayoutHandler) HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// Restore body for handler
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	switch provider {
	case "flutterwave":
		h.handleFlutterwavePayoutWebhook(c, body)
	case "stripe":
		h.handleStripePayoutWebhook(c, body)
	case "paystack":
		h.handlePaystackPayoutWebhook(c, body)
	case "paypal":
		h.handlePayPalPayoutWebhook(c, body)
	case "mtn_momo":
		h.handleMTNMoMoPayoutWebhook(c, body)
	case "orange_money":
		h.handleOrangeMoneyPayoutWebhook(c, body)
	case "wave":
		h.handleWavePayoutWebhook(c, body)
	case "thunes":
		h.handleThunesPayoutWebhook(c, body)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown provider"})
	}
}

func (h *PayoutHandler) handleFlutterwavePayoutWebhook(c *gin.Context, body []byte) {
	// Verify signature
	signature := c.GetHeader("verif-hash")
	if secret := h.webhookSecrets["flutterwave"]; secret != "" && signature != secret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var event struct {
		Event string `json:"event"`
		Data  struct {
			ID        int     `json:"id"`
			Reference string  `json:"reference"`
			Status    string  `json:"status"`
			Amount    float64 `json:"amount"`
			Currency  string  `json:"currency"`
			Message   string  `json:"complete_message"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()
	transactionRef := event.Data.Reference

	switch event.Event {
	case "transfer.completed":
		if event.Data.Status == "SUCCESSFUL" {
			h.completePayout(ctx, transactionRef, fmt.Sprintf("%d", event.Data.ID))
		} else {
			h.failPayout(ctx, transactionRef, event.Data.Message)
		}
	case "transfer.failed":
		h.failPayout(ctx, transactionRef, event.Data.Message)
	}

	log.Printf("[Flutterwave Payout Webhook] %s: %s -> %s", event.Event, transactionRef, event.Data.Status)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handleStripePayoutWebhook(c *gin.Context, body []byte) {
	signature := c.GetHeader("Stripe-Signature")
	if secret := h.webhookSecrets["stripe"]; secret != "" {
		if !verifyStripeWebhookSignature(body, signature, secret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID       string `json:"id"`
				Status   string `json:"status"`
				Metadata struct {
					TransactionID string `json:"transaction_id"`
				} `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()
	transactionRef := event.Data.Object.Metadata.TransactionID

	switch event.Type {
	case "payout.paid":
		h.completePayout(ctx, transactionRef, event.Data.Object.ID)
	case "payout.failed":
		h.failPayout(ctx, transactionRef, "Stripe payout failed")
	case "payout.canceled":
		h.failPayout(ctx, transactionRef, "Stripe payout canceled")
	}

	log.Printf("[Stripe Payout Webhook] %s: %s -> %s", event.Type, transactionRef, event.Data.Object.Status)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handlePaystackPayoutWebhook(c *gin.Context, body []byte) {
	signature := c.GetHeader("X-Paystack-Signature")
	if secret := h.webhookSecrets["paystack"]; secret != "" {
		if !verifyPaystackSignature(body, signature, secret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	var event struct {
		Event string `json:"event"`
		Data  struct {
			Reference string  `json:"reference"`
			Status    string  `json:"status"`
			Amount    float64 `json:"amount"`
			Currency  string  `json:"currency"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()
	transactionRef := event.Data.Reference

	switch event.Event {
	case "transfer.success":
		h.completePayout(ctx, transactionRef, event.Data.Reference)
	case "transfer.failed", "transfer.reversed":
		h.failPayout(ctx, transactionRef, "Paystack transfer failed")
	}

	log.Printf("[Paystack Payout Webhook] %s: %s -> %s", event.Event, transactionRef, event.Data.Status)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handlePayPalPayoutWebhook(c *gin.Context, body []byte) {
	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			BatchHeader struct {
				PayoutBatchID string `json:"payout_batch_id"`
				BatchStatus   string `json:"batch_status"`
			} `json:"batch_header"`
			Items []struct {
				PayoutItemID  string `json:"payout_item_id"`
				TransactionID string `json:"transaction_id"`
				ActivityID    string `json:"activity_id"`
				PayoutItemFee struct {
					Value string `json:"value"`
				} `json:"payout_item_fee"`
				TransactionStatus string `json:"transaction_status"`
				SenderBatchID     string `json:"sender_batch_id"`
			} `json:"items"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()

	switch event.EventType {
	case "PAYMENT.PAYOUTSBATCH.SUCCESS":
		for _, item := range event.Resource.Items {
			h.completePayout(ctx, item.SenderBatchID, item.PayoutItemID)
		}
	case "PAYMENT.PAYOUTSBATCH.DENIED":
		for _, item := range event.Resource.Items {
			h.failPayout(ctx, item.SenderBatchID, "PayPal payout denied")
		}
	}

	log.Printf("[PayPal Payout Webhook] %s", event.EventType)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handleMTNMoMoPayoutWebhook(c *gin.Context, body []byte) {
	var event struct {
		ExternalID    string  `json:"externalId"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		Status        string  `json:"status"`
		FinancialTxID string  `json:"financialTransactionId"`
		Reason        string  `json:"reason"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()

	switch event.Status {
	case "SUCCESSFUL":
		h.completePayout(ctx, event.ExternalID, event.FinancialTxID)
	case "FAILED":
		h.failPayout(ctx, event.ExternalID, event.Reason)
	}

	log.Printf("[MTN MoMo Payout Webhook] %s: %s -> %s", event.ExternalID, event.Status, event.Reason)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handleOrangeMoneyPayoutWebhook(c *gin.Context, body []byte) {
	var event struct {
		Status  string `json:"status"`
		TxnID   string `json:"txnid"`
		Message string `json:"message"`
		OrderID string `json:"order_id"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()

	switch event.Status {
	case "SUCCESS", "SUCCESSFUL":
		h.completePayout(ctx, event.OrderID, event.TxnID)
	case "FAILED", "FAILURE":
		h.failPayout(ctx, event.OrderID, event.Message)
	}

	log.Printf("[Orange Money Payout Webhook] %s: %s -> %s", event.OrderID, event.Status, event.Message)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handleWavePayoutWebhook(c *gin.Context, body []byte) {
	var event struct {
		ID        string  `json:"id"`
		Type      string  `json:"type"`
		ClientRef string  `json:"client_reference"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Status    string  `json:"status"`
		Error     string  `json:"error_message"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()

	switch event.Status {
	case "succeeded":
		h.completePayout(ctx, event.ClientRef, event.ID)
	case "failed":
		h.failPayout(ctx, event.ClientRef, event.Error)
	}

	log.Printf("[Wave Payout Webhook] %s: %s -> %s", event.ClientRef, event.Status, event.Error)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PayoutHandler) handleThunesPayoutWebhook(c *gin.Context, body []byte) {
	var event struct {
		Event string `json:"event"`
		Data  struct {
			ID         int    `json:"id"`
			ExternalID string `json:"external_id"`
			Status     string `json:"status"`
			StatusMsg  string `json:"status_message"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	ctx := c.Request.Context()

	switch event.Data.Status {
	case "50000", "COMPLETED":
		h.completePayout(ctx, event.Data.ExternalID, fmt.Sprintf("%d", event.Data.ID))
	case "40000", "FAILED":
		h.failPayout(ctx, event.Data.ExternalID, event.Data.StatusMsg)
	}

	log.Printf("[Thunes Payout Webhook] %s: %s -> %s", event.Data.ExternalID, event.Data.Status, event.Data.StatusMsg)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// ==================== HELPER FUNCTIONS ====================

func (h *PayoutHandler) completePayout(ctx context.Context, transactionRef, providerRef string) {
	// Update payout status
	err := h.payoutRepo.MarkCompleted(ctx, transactionRef, providerRef)
	if err != nil {
		log.Printf("[PayoutHandler] Failed to mark payout as completed: %v", err)
		return
	}

	// Confirm withdrawal in fund movement
	err = h.fundMovement.ConfirmWithdrawalPayout(ctx, transactionRef, "successful", providerRef)
	if err != nil {
		log.Printf("[PayoutHandler] Failed to confirm withdrawal: %v", err)
	}

	log.Printf("[PayoutHandler] ✅ Payout completed: %s | Provider Ref: %s", transactionRef, providerRef)
}

func (h *PayoutHandler) failPayout(ctx context.Context, transactionRef, reason string) {
	// Mark payout as failed
	err := h.payoutRepo.MarkFailed(ctx, transactionRef, reason)
	if err != nil {
		log.Printf("[PayoutHandler] Failed to mark payout as failed: %v", err)
		return
	}

	// Reverse the fund movement (refund user)
	err = h.fundMovement.ConfirmWithdrawalPayout(ctx, transactionRef, "failed", "")
	if err != nil {
		log.Printf("[PayoutHandler] Failed to reverse fund movement: %v", err)
	}

	log.Printf("[PayoutHandler] ❌ Payout failed: %s | Reason: %s (funds returned to user)", transactionRef, reason)
}

func (h *PayoutHandler) verifyUserPIN(ctx context.Context, userID, pin, authToken string) error {
	reqBody := map[string]string{
		"user_id": userID,
		"pin":     pin,
	}
	body, _ := json.Marshal(reqBody)

	// Call auth-service to verify PIN
	authServiceURL := "http://auth-service:8081/api/v1/verify-pin"
	req, _ := http.NewRequestWithContext(ctx, "POST", authServiceURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify PIN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PIN verification failed")
	}

	return nil
}

func (h *PayoutHandler) getUserWalletBalance(ctx context.Context, userID, walletID string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/wallets/%s/balance", h.walletServiceURL, walletID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Balance, nil
}

func (h *PayoutHandler) calculatePayoutFee(amount float64, payoutMethod string, instance *models.AggregatorInstanceWithDetails) float64 {
	var fee float64

	switch payoutMethod {
	case "mobile_money":
		fee = amount * 0.015 // 1.5%
		if fee < 50 {
			fee = 50 // Min fee 50 XOF
		}
	case "bank_transfer":
		fee = amount * 0.02 // 2%
		if fee < 500 {
			fee = 500 // Min fee 500 XOF
		}
	case "paypal":
		fee = amount*0.029 + 0.30 // PayPal fees
	default:
		fee = amount * 0.02
	}

	return fee
}

func (h *PayoutHandler) getRecipientAccountDisplay(req InitiatePayoutRequest) string {
	switch req.PayoutMethod {
	case "mobile_money":
		if len(req.MobileNumber) > 4 {
			return fmt.Sprintf("%s ***%s", req.MobileOperator, req.MobileNumber[len(req.MobileNumber)-4:])
		}
		return req.MobileNumber
	case "bank_transfer":
		if len(req.AccountNumber) > 4 {
			return fmt.Sprintf("%s ***%s", req.BankName, req.AccountNumber[len(req.AccountNumber)-4:])
		}
		return req.AccountNumber
	case "paypal":
		if len(req.PayPalEmail) > 5 {
			return fmt.Sprintf("PayPal: %s***", req.PayPalEmail[:5])
		}
		return req.PayPalEmail
	default:
		return "***"
	}
}

func (h *PayoutHandler) getRecipientAccountFromPayout(payout *repository.PayoutTransaction) string {
	switch payout.PayoutMethod {
	case "mobile_money":
		if payout.MobileNumber != "" && len(payout.MobileNumber) > 4 {
			return fmt.Sprintf("%s ***%s", payout.MobileOperator, payout.MobileNumber[len(payout.MobileNumber)-4:])
		}
	case "bank_transfer":
		if payout.AccountNumber != "" && len(payout.AccountNumber) > 4 {
			return fmt.Sprintf("%s ***%s", payout.BankName, payout.AccountNumber[len(payout.AccountNumber)-4:])
		}
	case "paypal":
		if payout.PayPalEmail != "" && len(payout.PayPalEmail) > 5 {
			return fmt.Sprintf("PayPal: %s***", payout.PayPalEmail[:5])
		}
	}
	return "***"
}

func (h *PayoutHandler) buildProviderPayoutRequest(req InitiatePayoutRequest, transactionID string) *providers.PayoutRequest {
	return &providers.PayoutRequest{
		ReferenceID:      transactionID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientEmail:   req.RecipientEmail,
		RecipientCountry: req.Country,
		PayoutMethod:     providers.PayoutMethod(req.PayoutMethod),
		BankCode:         req.BankCode,
		BankName:         req.BankName,
		AccountNumber:    req.AccountNumber,
		IBAN:             req.IBAN,
		SwiftCode:        req.SwiftCode,
		RoutingNumber:    req.RoutingNumber,
		MobileOperator:   req.MobileOperator,
		MobileNumber:     req.MobileNumber,
		Narration:        req.Narration,
		Metadata: map[string]string{
			"transaction_id": transactionID,
			"user_id":        req.UserID,
		},
	}
}

func (h *PayoutHandler) initiateProviderPayout(ctx context.Context, provider providers.CollectionProvider, req *providers.PayoutRequest, instance *models.AggregatorInstanceWithDetails) (*providers.PayoutResponse, error) {
	// Cast to PayoutProvider if supported
	payoutProvider, ok := provider.(providers.PayoutProvider)
	if !ok {
		// Provider doesn't support payouts, return mock success for demo
		if instance.ProviderCode == "demo" {
			return &providers.PayoutResponse{
				ProviderName:      "demo",
				ProviderReference: "DEMO_" + req.ReferenceID,
				ReferenceID:       req.ReferenceID,
				Status:            providers.PayoutStatusCompleted,
				Fee:               req.Amount * 0.015,
				AmountReceived:    req.Amount - (req.Amount * 0.015),
				EstimatedDelivery: time.Now().Add(5 * time.Minute),
				Message:           "Demo payout - instant success",
			}, nil
		}
		return nil, fmt.Errorf("provider does not support payouts")
	}

	// Use the real payout provider
	return payoutProvider.CreatePayout(ctx, req)
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func verifyPaystackSignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func verifyStripeWebhookSignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

// RegisterRoutes registers all payout routes
func (h *PayoutHandler) RegisterRoutes(router *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	// Public routes (authenticated users)
	payouts := router.Group("/payouts")
	{
		payouts.POST("/quote", h.GetPayoutQuote)
		payouts.POST("/initiate", h.InitiatePayout)
		payouts.GET("/:id/status", h.GetPayoutStatus)
		payouts.POST("/:id/cancel", h.CancelPayout)
		payouts.GET("/user/:user_id", h.GetUserPayouts)
		payouts.GET("/banks", h.GetBanks)
		payouts.GET("/mobile-operators", h.GetMobileOperators)

		// Webhooks (no auth - verified by signature)
		payouts.POST("/webhook/:provider", h.HandleWebhook)
	}
}
