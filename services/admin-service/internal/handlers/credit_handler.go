package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreditHandler handles all credit operations (individual, mass, promotions)
type CreditHandler struct {
	db     *sql.DB
	mainDB *sql.DB // Connection to main wallet-service DB for platform accounts
}

// NewCreditHandler creates a new credit handler
func NewCreditHandler(db *sql.DB, mainDB *sql.DB) *CreditHandler {
	return &CreditHandler{db: db, mainDB: mainDB}
}

// ========== Request/Response Models ==========

// SingleCreditRequest for individual user credit
type SingleCreditRequest struct {
	UserID     string  `json:"user_id" binding:"required"`
	Currency   string  `json:"currency" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Reason     string  `json:"reason" binding:"required"`
	ReasonType string  `json:"reason_type"` // compensation, bonus, promotion, refund, other
}

// MassCreditFilters for filtering users in mass credit
type MassCreditFilters struct {
	UserTypes            []string   `json:"user_types"` // individual, association, enterprise
	Countries            []string   `json:"countries"`  // CI, SN, FR, etc.
	TransactionDateFrom  *time.Time `json:"tx_date_from"`
	TransactionDateTo    *time.Time `json:"tx_date_to"`
	MinTransactionAmount *float64   `json:"min_tx_amount"`
	MaxTransactionAmount *float64   `json:"max_tx_amount"`
	WalletCurrency       string     `json:"wallet_currency"`
	KYCStatus            string     `json:"kyc_status"` // pending, approved, rejected
	HasActivity          *bool      `json:"has_activity"`
}

// MassCreditRequest for crediting multiple users
type MassCreditRequest struct {
	Filters      MassCreditFilters `json:"filters"`
	Amount       float64           `json:"amount" binding:"required,gt=0"`
	AmountType   string            `json:"amount_type"` // fixed, percentage
	Currency     string            `json:"currency" binding:"required"`
	Reason       string            `json:"reason" binding:"required"`
	ReasonType   string            `json:"reason_type"` // compensation, bonus, promotion, contest, other
	CampaignName string            `json:"campaign_name" binding:"required"`
	NotifyUsers  bool              `json:"notify_users"`
}

// PromotionCreditRequest for contest/promotion credits
type PromotionCreditRequest struct {
	CampaignName  string    `json:"campaign_name" binding:"required"`
	UserIDs       []string  `json:"user_ids" binding:"required"`
	Amounts       []float64 `json:"amounts"`        // Optional: different amounts per user
	UniformAmount float64   `json:"uniform_amount"` // Used if amounts is empty
	Currency      string    `json:"currency" binding:"required"`
	Reason        string    `json:"reason" binding:"required"`
	ReasonType    string    `json:"reason_type"`
	NotifyUsers   bool      `json:"notify_users"`
}

// CreditCampaign represents a credit campaign in the database
type CreditCampaign struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`   // single, mass, promotion
	Status       string     `json:"status"` // pending, processing, completed, failed
	Reason       string     `json:"reason"`
	ReasonType   string     `json:"reason_type"`
	Currency     string     `json:"currency"`
	TotalAmount  float64    `json:"total_amount"`
	UserCount    int        `json:"user_count"`
	SuccessCount int        `json:"success_count"`
	FailedCount  int        `json:"failed_count"`
	Filters      string     `json:"filters"` // JSON string of filters used
	AdminID      string     `json:"admin_id"`
	AdminName    string     `json:"admin_name"`
	HotWalletID  string     `json:"hot_wallet_id"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}

// CreditOperation represents a single credit operation
type CreditOperation struct {
	ID           string    `json:"id"`
	CampaignID   string    `json:"campaign_id"`
	UserID       string    `json:"user_id"`
	WalletID     string    `json:"wallet_id"`
	Currency     string    `json:"currency"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"` // pending, success, failed
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

// ========== Predefined Reason Types ==========
var ReasonTypes = []map[string]string{
	{"id": "compensation", "label": "Compensation (panne, erreur)", "icon": "🔧"},
	{"id": "bonus", "label": "Bonus de bienvenue", "icon": "🎁"},
	{"id": "promotion", "label": "Promotion spéciale", "icon": "🎉"},
	{"id": "contest", "label": "Prix concours/jeu", "icon": "🏆"},
	{"id": "refund", "label": "Remboursement", "icon": "💸"},
	{"id": "loyalty", "label": "Programme fidélité", "icon": "⭐"},
	{"id": "referral", "label": "Parrainage", "icon": "👥"},
	{"id": "other", "label": "Autre", "icon": "📝"},
}

// ========== Endpoints ==========

// GetReasonTypes returns predefined reason types
func (h *CreditHandler) GetReasonTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"reason_types": ReasonTypes})
}

// SingleCredit credits a single user
func (h *CreditHandler) SingleCredit(c *gin.Context) {
	var req SingleCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, _ := c.Get("admin_id")
	adminIDStr, _ := adminID.(string)

	// Get hot wallet for this currency
	hotWallet, err := h.getHotWallet(req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Hot wallet not found for %s: %v", req.Currency, err)})
		return
	}

	// Check hot wallet balance
	balance, err := h.getHotWalletBalance(hotWallet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check hot wallet balance"})
		return
	}

	if balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Insufficient hot wallet balance. Available: %.2f %s, Required: %.2f %s",
				balance, req.Currency, req.Amount, req.Currency),
		})
		return
	}

	// Create campaign for tracking
	campaignID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO credit_campaigns (id, name, type, status, reason, reason_type, currency, total_amount, user_count, admin_id, hot_wallet_id)
		VALUES ($1, $2, 'single', 'processing', $3, $4, $5, $6, 1, $7, $8)
	`, campaignID, fmt.Sprintf("Crédit individuel - %s", req.UserID[:8]), req.Reason, req.ReasonType, req.Currency, req.Amount, adminIDStr, hotWallet.ID)
	if err != nil {
		log.Printf("[CreditHandler] Failed to create campaign: %v", err)
	}

	// Get user's wallet for the currency
	walletID, err := h.getUserWallet(req.UserID, req.Currency)
	if err != nil {
		h.updateCampaignStatus(campaignID, "failed", 0, 1)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User wallet not found for %s: %v", req.Currency, err)})
		return
	}

	// Execute credit transfer from hot wallet to user wallet
	newBalance, err := h.executeCredit(hotWallet.ID, walletID, req.Amount, req.Currency, req.Reason, campaignID)
	if err != nil {
		h.updateCampaignStatus(campaignID, "failed", 0, 1)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Credit failed: %v", err)})
		return
	}

	h.updateCampaignStatus(campaignID, "completed", 1, 0)

	// Log the operation
	h.logCreditOperation(campaignID, req.UserID, walletID, req.Currency, req.Amount, "success", "")

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"campaign_id": campaignID,
		"new_balance": newBalance,
		"currency":    req.Currency,
		"message":     fmt.Sprintf("Crédit de %.2f %s effectué avec succès", req.Amount, req.Currency),
	})
}

// MassCreditPreview previews mass credit without executing
func (h *CreditHandler) MassCreditPreview(c *gin.Context) {
	var req MassCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get matching users
	users, err := h.getFilteredUsers(req.Filters, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to filter users: %v", err)})
		return
	}

	totalAmount := float64(len(users)) * req.Amount

	// Get hot wallet balance
	hotWallet, err := h.getHotWallet(req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Hot wallet not found for %s", req.Currency)})
		return
	}

	balance, _ := h.getHotWalletBalance(hotWallet.ID)

	c.JSON(http.StatusOK, gin.H{
		"user_count":         len(users),
		"amount_per_user":    req.Amount,
		"total_amount":       totalAmount,
		"currency":           req.Currency,
		"hot_wallet_balance": balance,
		"sufficient_funds":   balance >= totalAmount,
		"users_preview":      users[:min(10, len(users))], // First 10 users as preview
	})
}

// MassCredit executes mass credit to filtered users
func (h *CreditHandler) MassCredit(c *gin.Context) {
	var req MassCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, _ := c.Get("admin_id")
	adminIDStr, _ := adminID.(string)

	// Get matching users
	users, err := h.getFilteredUsers(req.Filters, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to filter users: %v", err)})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No users match the specified filters"})
		return
	}

	totalAmount := float64(len(users)) * req.Amount

	// Get hot wallet
	hotWallet, err := h.getHotWallet(req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Hot wallet not found for %s", req.Currency)})
		return
	}

	// Check balance
	balance, _ := h.getHotWalletBalance(hotWallet.ID)
	if balance < totalAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Insufficient hot wallet balance. Need %.2f %s, have %.2f %s",
				totalAmount, req.Currency, balance, req.Currency),
		})
		return
	}

	// Create campaign
	campaignID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO credit_campaigns (id, name, type, status, reason, reason_type, currency, total_amount, user_count, admin_id, hot_wallet_id)
		VALUES ($1, $2, 'mass', 'processing', $3, $4, $5, $6, $7, $8, $9)
	`, campaignID, req.CampaignName, req.Reason, req.ReasonType, req.Currency, totalAmount, len(users), adminIDStr, hotWallet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	// Process credits in background (for large batches)
	go h.processMassCredit(campaignID, hotWallet.ID, users, req.Amount, req.Currency, req.Reason, req.NotifyUsers)

	c.JSON(http.StatusAccepted, gin.H{
		"success":      true,
		"campaign_id":  campaignID,
		"user_count":   len(users),
		"total_amount": totalAmount,
		"status":       "processing",
		"message":      fmt.Sprintf("Crédit de masse lancé pour %d utilisateurs", len(users)),
	})
}

// PromotionCredit credits specific users for promotions/contests
func (h *CreditHandler) PromotionCredit(c *gin.Context) {
	var req PromotionCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one user ID is required"})
		return
	}

	adminID, _ := c.Get("admin_id")
	adminIDStr, _ := adminID.(string)

	// Calculate total amount
	var totalAmount float64
	if len(req.Amounts) == len(req.UserIDs) {
		for _, amt := range req.Amounts {
			totalAmount += amt
		}
	} else {
		totalAmount = req.UniformAmount * float64(len(req.UserIDs))
	}

	// Get hot wallet
	hotWallet, err := h.getHotWallet(req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Hot wallet not found for %s", req.Currency)})
		return
	}

	// Check balance
	balance, _ := h.getHotWalletBalance(hotWallet.ID)
	if balance < totalAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Insufficient hot wallet balance. Need %.2f %s, have %.2f %s",
				totalAmount, req.Currency, balance, req.Currency),
		})
		return
	}

	// Create campaign
	campaignID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO credit_campaigns (id, name, type, status, reason, reason_type, currency, total_amount, user_count, admin_id, hot_wallet_id)
		VALUES ($1, $2, 'promotion', 'processing', $3, $4, $5, $6, $7, $8, $9)
	`, campaignID, req.CampaignName, req.Reason, req.ReasonType, req.Currency, totalAmount, len(req.UserIDs), adminIDStr, hotWallet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	// Process in background
	go h.processPromotionCredit(campaignID, hotWallet.ID, req.UserIDs, req.Amounts, req.UniformAmount, req.Currency, req.Reason, req.NotifyUsers)

	c.JSON(http.StatusAccepted, gin.H{
		"success":      true,
		"campaign_id":  campaignID,
		"user_count":   len(req.UserIDs),
		"total_amount": totalAmount,
		"status":       "processing",
		"message":      fmt.Sprintf("Promotion lancée pour %d utilisateurs", len(req.UserIDs)),
	})
}

// GetCampaigns returns list of credit campaigns
func (h *CreditHandler) GetCampaigns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	campaignType := c.Query("type")
	status := c.Query("status")

	offset := (page - 1) * limit

	query := `
		SELECT id, name, type, status, reason, reason_type, currency, total_amount, 
		       user_count, success_count, failed_count, admin_id, created_at, completed_at
		FROM credit_campaigns
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if campaignType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, campaignType)
		argIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
		return
	}
	defer rows.Close()

	var campaigns []CreditCampaign
	for rows.Next() {
		var camp CreditCampaign
		err := rows.Scan(&camp.ID, &camp.Name, &camp.Type, &camp.Status, &camp.Reason, &camp.ReasonType,
			&camp.Currency, &camp.TotalAmount, &camp.UserCount, &camp.SuccessCount, &camp.FailedCount,
			&camp.AdminID, &camp.CreatedAt, &camp.CompletedAt)
		if err != nil {
			continue
		}
		campaigns = append(campaigns, camp)
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

// GetCampaignDetails returns detailed info about a campaign
func (h *CreditHandler) GetCampaignDetails(c *gin.Context) {
	campaignID := c.Param("id")

	var camp CreditCampaign
	err := h.db.QueryRow(`
		SELECT id, name, type, status, reason, reason_type, currency, total_amount,
		       user_count, success_count, failed_count, admin_id, hot_wallet_id, created_at, completed_at
		FROM credit_campaigns WHERE id = $1
	`, campaignID).Scan(&camp.ID, &camp.Name, &camp.Type, &camp.Status, &camp.Reason, &camp.ReasonType,
		&camp.Currency, &camp.TotalAmount, &camp.UserCount, &camp.SuccessCount, &camp.FailedCount,
		&camp.AdminID, &camp.HotWalletID, &camp.CreatedAt, &camp.CompletedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Get operations for this campaign
	rows, err := h.db.Query(`
		SELECT id, user_id, wallet_id, currency, amount, status, error_message, created_at
		FROM credit_operations WHERE campaign_id = $1 ORDER BY created_at DESC LIMIT 100
	`, campaignID)
	if err == nil {
		defer rows.Close()
		var ops []CreditOperation
		for rows.Next() {
			var op CreditOperation
			rows.Scan(&op.ID, &op.UserID, &op.WalletID, &op.Currency, &op.Amount, &op.Status, &op.ErrorMessage, &op.CreatedAt)
			ops = append(ops, op)
		}
		c.JSON(http.StatusOK, gin.H{"campaign": camp, "operations": ops})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaign": camp})
}

// GetCreditLogs returns all credit operations for audit
func (h *CreditHandler) GetCreditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	rows, err := h.db.Query(`
		SELECT co.id, co.campaign_id, co.user_id, co.wallet_id, co.currency, co.amount, 
		       co.status, co.error_message, co.created_at, cc.name as campaign_name, cc.reason_type
		FROM credit_operations co
		LEFT JOIN credit_campaigns cc ON co.campaign_id = cc.id
		ORDER BY co.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, campaignID, userID, walletID, currency, status, errorMsg, campaignName, reasonType string
		var amount float64
		var createdAt time.Time
		err := rows.Scan(&id, &campaignID, &userID, &walletID, &currency, &amount, &status, &errorMsg, &createdAt, &campaignName, &reasonType)
		if err != nil {
			continue
		}
		logs = append(logs, map[string]interface{}{
			"id":            id,
			"campaign_id":   campaignID,
			"campaign_name": campaignName,
			"user_id":       userID,
			"wallet_id":     walletID,
			"currency":      currency,
			"amount":        amount,
			"status":        status,
			"error_message": errorMsg,
			"reason_type":   reasonType,
			"created_at":    createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetHotWallets returns all hot wallets with balances
func (h *CreditHandler) GetHotWallets(c *gin.Context) {
	rows, err := h.mainDB.Query(`
		SELECT id, name, currency, balance, account_type, is_active 
		FROM platform_accounts 
		WHERE account_type = 'operations' 
		ORDER BY currency
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hot wallets"})
		return
	}
	defer rows.Close()

	var wallets []map[string]interface{}
	for rows.Next() {
		var id, name, currency, accountType string
		var balance float64
		var isActive bool
		rows.Scan(&id, &name, &currency, &balance, &accountType, &isActive)
		wallets = append(wallets, map[string]interface{}{
			"id":           id,
			"name":         name,
			"currency":     currency,
			"balance":      balance,
			"account_type": accountType,
			"is_active":    isActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{"hot_wallets": wallets})
}

// ========== Helper Functions ==========

type HotWallet struct {
	ID       string
	Currency string
	Balance  float64
}

type FilteredUser struct {
	UserID   string  `json:"user_id"`
	Email    string  `json:"email"`
	WalletID string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}

func (h *CreditHandler) getHotWallet(currency string) (*HotWallet, error) {
	var hw HotWallet
	err := h.mainDB.QueryRow(`
		SELECT id, currency FROM platform_accounts 
		WHERE currency = $1 AND account_type = 'operations' AND is_active = true
	`, currency).Scan(&hw.ID, &hw.Currency)
	if err != nil {
		return nil, err
	}
	return &hw, nil
}

func (h *CreditHandler) getHotWalletBalance(walletID string) (float64, error) {
	var balance float64
	err := h.mainDB.QueryRow(`SELECT balance FROM platform_accounts WHERE id = $1`, walletID).Scan(&balance)
	return balance, err
}

func (h *CreditHandler) getUserWallet(userID, currency string) (string, error) {
	var walletID string
	err := h.mainDB.QueryRow(`
		SELECT id FROM wallets WHERE user_id = $1 AND currency = $2 AND status = 'active'
	`, userID, currency).Scan(&walletID)
	return walletID, err
}

func (h *CreditHandler) getFilteredUsers(filters MassCreditFilters, currency string) ([]FilteredUser, error) {
	query := `
		SELECT DISTINCT w.user_id, u.email, w.id as wallet_id, w.balance
		FROM wallets w
		JOIN users u ON w.user_id = u.id::text
		WHERE w.currency = $1 AND w.status = 'active'
	`
	args := []interface{}{currency}
	argIdx := 2

	// Apply filters
	if len(filters.UserTypes) > 0 {
		query += fmt.Sprintf(" AND u.account_type = ANY($%d)", argIdx)
		args = append(args, filters.UserTypes)
		argIdx++
	}

	if len(filters.Countries) > 0 {
		query += fmt.Sprintf(" AND u.country = ANY($%d)", argIdx)
		args = append(args, filters.Countries)
		argIdx++
	}

	if filters.TransactionDateFrom != nil && filters.TransactionDateTo != nil {
		query += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM transactions t WHERE t.user_id = w.user_id 
			AND t.created_at BETWEEN $%d AND $%d
		)`, argIdx, argIdx+1)
		args = append(args, filters.TransactionDateFrom, filters.TransactionDateTo)
		argIdx += 2
	}

	if filters.MinTransactionAmount != nil {
		query += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM transactions t WHERE t.user_id = w.user_id AND t.amount >= $%d
		)`, argIdx)
		args = append(args, *filters.MinTransactionAmount)
		argIdx++
	}

	if filters.KYCStatus != "" {
		query += fmt.Sprintf(" AND u.kyc_status = $%d", argIdx)
		args = append(args, filters.KYCStatus)
		argIdx++
	}

	query += " LIMIT 10000" // Safety limit

	rows, err := h.mainDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []FilteredUser
	for rows.Next() {
		var user FilteredUser
		if err := rows.Scan(&user.UserID, &user.Email, &user.WalletID, &user.Balance); err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}

func (h *CreditHandler) executeCredit(fromWalletID, toWalletID string, amount float64, currency, reason, campaignID string) (float64, error) {
	tx, err := h.mainDB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Debit from hot wallet
	_, err = tx.Exec(`UPDATE platform_accounts SET balance = balance - $1 WHERE id = $2`, amount, fromWalletID)
	if err != nil {
		return 0, fmt.Errorf("failed to debit hot wallet: %v", err)
	}

	// Credit to user wallet
	var newBalance float64
	err = tx.QueryRow(`
		UPDATE wallets SET balance = balance + $1 WHERE id = $2 RETURNING balance
	`, amount, toWalletID).Scan(&newBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to credit user wallet: %v", err)
	}

	// Record transaction with proper description
	_, err = tx.Exec(`
		INSERT INTO transactions (id, user_id, type, amount, currency, status, description, created_at)
		SELECT gen_random_uuid(), user_id, 'credit', $1, $2, 'completed', $3, NOW()
		FROM wallets WHERE id = $4
	`, amount, currency, reason, toWalletID)
	if err != nil {
		log.Printf("[CreditHandler] Warning: Failed to record transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (h *CreditHandler) processMassCredit(campaignID, hotWalletID string, users []FilteredUser, amount float64, currency, reason string, notify bool) {
	successCount := 0
	failedCount := 0

	for _, user := range users {
		_, err := h.executeCredit(hotWalletID, user.WalletID, amount, currency, reason, campaignID)
		if err != nil {
			h.logCreditOperation(campaignID, user.UserID, user.WalletID, currency, amount, "failed", err.Error())
			failedCount++
		} else {
			h.logCreditOperation(campaignID, user.UserID, user.WalletID, currency, amount, "success", "")
			successCount++
		}
	}

	h.updateCampaignStatus(campaignID, "completed", successCount, failedCount)
}

func (h *CreditHandler) processPromotionCredit(campaignID, hotWalletID string, userIDs []string, amounts []float64, uniformAmount float64, currency, reason string, notify bool) {
	successCount := 0
	failedCount := 0

	for i, userID := range userIDs {
		amount := uniformAmount
		if len(amounts) == len(userIDs) {
			amount = amounts[i]
		}

		walletID, err := h.getUserWallet(userID, currency)
		if err != nil {
			h.logCreditOperation(campaignID, userID, "", currency, amount, "failed", "Wallet not found")
			failedCount++
			continue
		}

		_, err = h.executeCredit(hotWalletID, walletID, amount, currency, reason, campaignID)
		if err != nil {
			h.logCreditOperation(campaignID, userID, walletID, currency, amount, "failed", err.Error())
			failedCount++
		} else {
			h.logCreditOperation(campaignID, userID, walletID, currency, amount, "success", "")
			successCount++
		}
	}

	h.updateCampaignStatus(campaignID, "completed", successCount, failedCount)
}

func (h *CreditHandler) updateCampaignStatus(campaignID, status string, successCount, failedCount int) {
	completedAt := "NULL"
	if status == "completed" || status == "failed" {
		completedAt = "NOW()"
	}
	h.db.Exec(fmt.Sprintf(`
		UPDATE credit_campaigns 
		SET status = $1, success_count = $2, failed_count = $3, completed_at = %s
		WHERE id = $4
	`, completedAt), status, successCount, failedCount, campaignID)
}

func (h *CreditHandler) logCreditOperation(campaignID, userID, walletID, currency string, amount float64, status, errorMsg string) {
	h.db.Exec(`
		INSERT INTO credit_operations (id, campaign_id, user_id, wallet_id, currency, amount, status, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, uuid.New().String(), campaignID, userID, walletID, currency, amount, status, errorMsg)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
