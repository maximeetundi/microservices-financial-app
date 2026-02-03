package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TestModeHandler handles admin test mode operations
// Allows manual crediting of user accounts without waiting for webhooks
type TestModeHandler struct {
	db             *sql.DB
	walletClient   WalletServiceClient
	platformClient PlatformAccountService
}

// WalletServiceClient interface for wallet operations
type WalletServiceClient interface {
	CreditWallet(ctx context.Context, walletID string, amount float64, reference string) error
	GetWalletByUserAndCurrency(ctx context.Context, userID, currency string) (*WalletInfo, error)
}

// PlatformAccountService interface for platform operations
type PlatformAccountService interface {
	DebitPlatformHotWallet(ctx context.Context, currency string, amount float64, reference, description string) error
}

// WalletInfo represents wallet information
type WalletInfo struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

// NewTestModeHandler creates a new test mode handler
func NewTestModeHandler(db *sql.DB, walletClient WalletServiceClient, platformClient PlatformAccountService) *TestModeHandler {
	return &TestModeHandler{
		db:             db,
		walletClient:   walletClient,
		platformClient: platformClient,
	}
}

// ManualCreditRequest represents a manual credit request
type ManualCreditRequest struct {
	UserID     string  `json:"user_id" binding:"required"`
	Currency   string  `json:"currency" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Reason     string  `json:"reason" binding:"required"`
	Aggregator string  `json:"aggregator"` // Optional: simulates which aggregator
	TestMode   bool    `json:"test_mode"`  // If true, doesn't affect real balances
}

// ManualCreditResponse represents the response for a manual credit
type ManualCreditResponse struct {
	Success       bool    `json:"success"`
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	WalletID      string  `json:"wallet_id"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance,omitempty"`
	Reason        string  `json:"reason"`
	IsTestMode    bool    `json:"is_test_mode"`
	Timestamp     string  `json:"timestamp"`
}

// TestCreditLog represents a log entry for test credits
type TestCreditLog struct {
	ID         string    `json:"id"`
	AdminID    string    `json:"admin_id"`
	UserID     string    `json:"user_id"`
	WalletID   string    `json:"wallet_id"`
	Currency   string    `json:"currency"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
	Aggregator string    `json:"aggregator"`
	IsTestMode bool      `json:"is_test_mode"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// InitSchema creates the test_credit_logs table
func (h *TestModeHandler) InitSchema() error {
	_, err := h.db.Exec(`
		CREATE TABLE IF NOT EXISTS test_credit_logs (
			id VARCHAR(36) PRIMARY KEY,
			admin_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			wallet_id VARCHAR(36),
			currency VARCHAR(10) NOT NULL,
			amount DECIMAL(20, 8) NOT NULL,
			reason TEXT NOT NULL,
			aggregator VARCHAR(50),
			is_test_mode BOOLEAN DEFAULT false,
			status VARCHAR(20) DEFAULT 'completed',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create test_credit_logs table: %w", err)
	}
	log.Println("[TestModeHandler] Table test_credit_logs created/verified")
	return nil
}

// ManualCredit credits a user's wallet manually (for testing or admin override)
// POST /admin/test-mode/credit
func (h *TestModeHandler) ManualCredit(c *gin.Context) {
	// Get admin ID from context (set by auth middleware)
	adminID := c.GetString("admin_id")
	if adminID == "" {
		adminID = c.GetString("user_id") // Fallback
	}

	var req ManualCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	transactionID := uuid.New().String()

	log.Printf("[TestMode] Admin %s initiating manual credit: %.2f %s to user %s",
		adminID, req.Amount, req.Currency, req.UserID)

	// Get user's wallet for the currency
	wallet, err := h.walletClient.GetWalletByUserAndCurrency(ctx, req.UserID, req.Currency)
	if err != nil {
		log.Printf("[TestMode] Error getting wallet: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User wallet not found for currency " + req.Currency})
		return
	}

	var newBalance float64

	if !req.TestMode {
		// REAL MODE: Actually credit the wallet
		// 1. Debit from platform hot wallet
		if h.platformClient != nil {
			if err := h.platformClient.DebitPlatformHotWallet(ctx, req.Currency, req.Amount,
				"ADMIN_MANUAL_CREDIT_"+transactionID, "Manual credit by admin"); err != nil {
				log.Printf("[TestMode] Warning: Failed to debit platform wallet: %v", err)
				// Continue anyway - admin override
			}
		}

		// 2. Credit user wallet
		reference := fmt.Sprintf("ADMIN_CREDIT_%s_%s", req.Aggregator, transactionID)
		if err := h.walletClient.CreditWallet(ctx, wallet.ID, req.Amount, reference); err != nil {
			log.Printf("[TestMode] Error crediting wallet: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to credit wallet: " + err.Error()})
			return
		}

		newBalance = wallet.Balance + req.Amount
		log.Printf("[TestMode] ✅ REAL credit completed: %.2f %s to wallet %s", req.Amount, req.Currency, wallet.ID)
	} else {
		// TEST MODE: Just simulate, don't actually credit
		newBalance = wallet.Balance // Balance unchanged
		log.Printf("[TestMode] 🧪 TEST credit simulated: %.2f %s (no actual transfer)", req.Amount, req.Currency)
	}

	// Log the credit operation
	_, err = h.db.Exec(`
		INSERT INTO test_credit_logs (id, admin_id, user_id, wallet_id, currency, amount, reason, aggregator, is_test_mode, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'completed')
	`, transactionID, adminID, req.UserID, wallet.ID, req.Currency, req.Amount, req.Reason, req.Aggregator, req.TestMode)
	if err != nil {
		log.Printf("[TestMode] Warning: Failed to log credit: %v", err)
	}

	response := ManualCreditResponse{
		Success:       true,
		TransactionID: transactionID,
		UserID:        req.UserID,
		WalletID:      wallet.ID,
		Currency:      req.Currency,
		Amount:        req.Amount,
		NewBalance:    newBalance,
		Reason:        req.Reason,
		IsTestMode:    req.TestMode,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// SimulateWebhook simulates a webhook callback for testing
// POST /admin/test-mode/simulate-webhook
func (h *TestModeHandler) SimulateWebhook(c *gin.Context) {
	var req struct {
		Aggregator    string  `json:"aggregator" binding:"required"` // stripe, flutterwave, etc.
		EventType     string  `json:"event_type" binding:"required"` // payment.success, etc.
		UserID        string  `json:"user_id" binding:"required"`
		Currency      string  `json:"currency" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		TransactionID string  `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if req.TransactionID == "" {
		req.TransactionID = uuid.New().String()
	}

	log.Printf("[TestMode] Simulating %s webhook: %s for %.2f %s", req.Aggregator, req.EventType, req.Amount, req.Currency)

	// Get user's wallet
	wallet, err := h.walletClient.GetWalletByUserAndCurrency(ctx, req.UserID, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User wallet not found for currency " + req.Currency})
		return
	}

	// Simulate successful payment - credit the wallet
	switch req.EventType {
	case "payment.success", "payment_intent.succeeded", "charge.completed":
		// 1. Debit from platform hot wallet first
		if h.platformClient != nil {
			if err := h.platformClient.DebitPlatformHotWallet(ctx, req.Currency, req.Amount,
				"SIMULATED_"+req.Aggregator+"_"+req.TransactionID, "Simulated webhook deposit"); err != nil {
				log.Printf("[TestMode] Warning: Failed to debit platform wallet: %v", err)
				// Continue anyway - admin override for testing
			}
		}

		// 2. Credit user wallet
		reference := fmt.Sprintf("SIMULATED_%s_%s", req.Aggregator, req.TransactionID)
		if err := h.walletClient.CreditWallet(ctx, wallet.ID, req.Amount, reference); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to credit wallet"})
			return
		}
		log.Printf("[TestMode] ✅ Simulated payment success - credited %.2f %s (debited from hot wallet)", req.Amount, req.Currency)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown event type: " + req.EventType})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        fmt.Sprintf("Simulated %s %s webhook processed", req.Aggregator, req.EventType),
		"transaction_id": req.TransactionID,
		"wallet_id":      wallet.ID,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})
}

// GetCreditLogs returns all manual credit logs
// GET /admin/test-mode/logs
func (h *TestModeHandler) GetCreditLogs(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, admin_id, user_id, wallet_id, currency, amount, reason, aggregator, is_test_mode, status, created_at
		FROM test_credit_logs
		ORDER BY created_at DESC
		LIMIT 100
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	defer rows.Close()

	var logs []TestCreditLog
	for rows.Next() {
		var l TestCreditLog
		var walletID, aggregator sql.NullString
		err := rows.Scan(&l.ID, &l.AdminID, &l.UserID, &walletID, &l.Currency, &l.Amount,
			&l.Reason, &aggregator, &l.IsTestMode, &l.Status, &l.CreatedAt)
		if err != nil {
			continue
		}
		if walletID.Valid {
			l.WalletID = walletID.String
		}
		if aggregator.Valid {
			l.Aggregator = aggregator.String
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

// QuickCredit provides quick credit buttons by currency
// POST /admin/test-mode/quick-credit
func (h *TestModeHandler) QuickCredit(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Currency string `json:"currency" binding:"required"`
		Preset   string `json:"preset" binding:"required"` // small, medium, large
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define preset amounts per currency
	presets := map[string]map[string]float64{
		"USD":  {"small": 10, "medium": 100, "large": 1000},
		"EUR":  {"small": 10, "medium": 100, "large": 1000},
		"GBP":  {"small": 10, "medium": 100, "large": 1000},
		"NGN":  {"small": 5000, "medium": 50000, "large": 500000},
		"KES":  {"small": 1000, "medium": 10000, "large": 100000},
		"GHS":  {"small": 100, "medium": 1000, "large": 10000},
		"XOF":  {"small": 5000, "medium": 50000, "large": 500000},
		"XAF":  {"small": 5000, "medium": 50000, "large": 500000},
		"ZAR":  {"small": 100, "medium": 1000, "large": 10000},
		"BTC":  {"small": 0.0001, "medium": 0.001, "large": 0.01},
		"ETH":  {"small": 0.01, "medium": 0.1, "large": 1},
		"USDT": {"small": 10, "medium": 100, "large": 1000},
	}

	currencyPresets, ok := presets[req.Currency]
	if !ok {
		// Default preset
		currencyPresets = map[string]float64{"small": 10, "medium": 100, "large": 1000}
	}

	amount, ok := currencyPresets[req.Preset]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preset. Use: small, medium, large"})
		return
	}

	// Redirect to manual credit with the amount
	creditReq := ManualCreditRequest{
		UserID:   req.UserID,
		Currency: req.Currency,
		Amount:   amount,
		Reason:   fmt.Sprintf("Quick credit (%s) for testing", req.Preset),
		TestMode: false,
	}

	// Create new context with the request
	c.Set("quick_credit", true)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "credit_request", creditReq))

	// Perform the credit
	adminID := c.GetString("admin_id")
	if adminID == "" {
		adminID = c.GetString("user_id")
	}

	ctx := c.Request.Context()
	transactionID := uuid.New().String()

	wallet, err := h.walletClient.GetWalletByUserAndCurrency(ctx, req.UserID, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User wallet not found"})
		return
	}

	reference := fmt.Sprintf("QUICK_CREDIT_%s_%s", req.Preset, transactionID)

	// 1. Debit from platform hot wallet first
	if h.platformClient != nil {
		if err := h.platformClient.DebitPlatformHotWallet(ctx, req.Currency, amount,
			reference, "Quick credit by admin"); err != nil {
			log.Printf("[TestMode] Warning: Failed to debit platform wallet: %v", err)
			// Continue anyway - admin override for testing
		}
	}

	// 2. Credit user wallet
	if err := h.walletClient.CreditWallet(ctx, wallet.ID, amount, reference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to credit wallet"})
		return
	}

	// Log it
	h.db.Exec(`
		INSERT INTO test_credit_logs (id, admin_id, user_id, wallet_id, currency, amount, reason, is_test_mode, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, false, 'completed')
	`, transactionID, adminID, req.UserID, wallet.ID, req.Currency, amount, creditReq.Reason)

	log.Printf("[TestMode] ⚡ Quick credit %s: %.2f %s to user %s (debited from hot wallet)", req.Preset, amount, req.Currency, req.UserID)

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"transaction_id": transactionID,
		"user_id":        req.UserID,
		"wallet_id":      wallet.ID,
		"currency":       req.Currency,
		"amount":         amount,
		"preset":         req.Preset,
		"new_balance":    wallet.Balance + amount,
	})
}

// GetPresets returns available quick credit presets
// GET /admin/test-mode/presets
func (h *TestModeHandler) GetPresets(c *gin.Context) {
	presets := map[string]map[string]float64{
		"USD":  {"small": 10, "medium": 100, "large": 1000},
		"EUR":  {"small": 10, "medium": 100, "large": 1000},
		"GBP":  {"small": 10, "medium": 100, "large": 1000},
		"NGN":  {"small": 5000, "medium": 50000, "large": 500000},
		"KES":  {"small": 1000, "medium": 10000, "large": 100000},
		"GHS":  {"small": 100, "medium": 1000, "large": 10000},
		"XOF":  {"small": 5000, "medium": 50000, "large": 500000},
		"XAF":  {"small": 5000, "medium": 50000, "large": 500000},
		"ZAR":  {"small": 100, "medium": 1000, "large": 10000},
		"BTC":  {"small": 0.0001, "medium": 0.001, "large": 0.01},
		"ETH":  {"small": 0.01, "medium": 0.1, "large": 1},
		"USDT": {"small": 10, "medium": 100, "large": 1000},
	}

	c.JSON(http.StatusOK, gin.H{"presets": presets})
}

// RegisterRoutes registers test mode routes
func (h *TestModeHandler) RegisterRoutes(adminRouter *gin.RouterGroup) {
	testMode := adminRouter.Group("/test-mode")
	{
		testMode.POST("/credit", h.ManualCredit)
		testMode.POST("/simulate-webhook", h.SimulateWebhook)
		testMode.POST("/quick-credit", h.QuickCredit)
		testMode.GET("/logs", h.GetCreditLogs)
		testMode.GET("/presets", h.GetPresets)
	}
}
