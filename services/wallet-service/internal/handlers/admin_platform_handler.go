package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminPlatformHandler struct {
	platformService *services.PlatformAccountService
	cryptoService   *services.CryptoService
}

func NewAdminPlatformHandler(platformService *services.PlatformAccountService, cryptoService *services.CryptoService) *AdminPlatformHandler {
	return &AdminPlatformHandler{
		platformService: platformService,
		cryptoService:   cryptoService,
	}
}

// ==================== Platform Fiat Accounts ====================

// GetAccounts returns all platform fiat accounts
func (h *AdminPlatformHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.platformService.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get accounts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// GetAccount returns a single platform account by ID
func (h *AdminPlatformHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	account, err := h.platformService.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}
	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, account)
}

// CreateAccount creates a new platform account
func (h *AdminPlatformHandler) CreateAccount(c *gin.Context) {
	var req models.CreatePlatformAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.platformService.CreateAccount(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

// CreditAccount allows admin to manually credit a platform account
func (h *AdminPlatformHandler) CreditAccount(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("user_id") // From JWT

	var req models.AdminCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AccountID = id

	err := h.platformService.AdminCreditAccount(id, req.Amount, req.Description, adminID, req.Reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Account credited successfully",
		"account_id": id,
		"amount":     req.Amount,
	})
}

// DebitAccount allows admin to manually debit a platform account
func (h *AdminPlatformHandler) DebitAccount(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("user_id") // From JWT

	var req models.AdminCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.platformService.AdminDebitAccount(id, req.Amount, req.Description, adminID, req.Reference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Account debited successfully",
		"account_id": id,
		"amount":     req.Amount,
	})
}

// ==================== Crypto Wallets ====================

// GetCryptoWallets returns all platform crypto wallets
func (h *AdminPlatformHandler) GetCryptoWallets(c *gin.Context) {
	wallets, err := h.platformService.GetAllCryptoWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get crypto wallets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

// GetCryptoWallet returns a single crypto wallet by ID
func (h *AdminPlatformHandler) GetCryptoWallet(c *gin.Context) {
	id := c.Param("id")
	wallet, err := h.platformService.GetCryptoWalletByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet"})
		return
	}
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

// CreateCryptoWallet adds a new crypto wallet address
func (h *AdminPlatformHandler) CreateCryptoWallet(c *gin.Context) {
	var req models.CreatePlatformCryptoWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.platformService.CreateCryptoWallet(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wallet)
}

// SyncCryptoWalletBalance synchronizes balance from blockchain
func (h *AdminPlatformHandler) SyncCryptoWalletBalance(c *gin.Context) {
	id := c.Param("id")

	wallet, err := h.platformService.GetCryptoWalletByID(id)
	if err != nil || wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	if wallet.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet has no address configured"})
		return
	}

	// Use CryptoService/FailoverProvider to fetch blockchain balance
	blockchainBalance, err := h.cryptoService.GetBlockchainBalance(wallet.Currency, wallet.Address)
	if err != nil {
		log.Printf("[AdminPlatform] Failed to sync balance for wallet %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blockchain balance"})
		return
	}

	previousBalance := wallet.Balance

	// Update balance in database
	err = h.platformService.UpdateCryptoWalletBalance(id, blockchainBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance in database"})
		return
	}

	log.Printf("[AdminPlatform] Synced balance for wallet %s: %.8f -> %.8f %s",
		wallet.Label, previousBalance, blockchainBalance, wallet.Currency)

	c.JSON(http.StatusOK, gin.H{
		"message":          "Balance synced successfully",
		"wallet_id":        id,
		"previous_balance": previousBalance,
		"new_balance":      blockchainBalance,
		"currency":         wallet.Currency,
	})
}

// ==================== Admin Consolidation ====================

// ConsolidateFunds moves user wallet funds to platform hot/cold wallet (DB only)
func (h *AdminPlatformHandler) ConsolidateFunds(c *gin.Context) {
	adminID := c.GetString("user_id") // From JWT

	var req services.ConsolidateUserFundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AdminID = adminID

	// Validate target type
	if req.TargetType != "hot" && req.TargetType != "cold" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_type must be 'hot' or 'cold'"})
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	err := h.platformService.ConsolidateUserFunds(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Funds consolidated successfully",
		"user_wallet_id": req.UserWalletID,
		"target_type":    req.TargetType,
		"amount":         req.Amount,
		"currency":       req.Currency,
	})
}

// ==================== Transaction Ledger ====================

// GetTransactions returns platform transaction history
func (h *AdminPlatformHandler) GetTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.platformService.GetTransactions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// GetReconciliation returns balance reconciliation report
func (h *AdminPlatformHandler) GetReconciliation(c *gin.Context) {
	balances, err := h.platformService.GetReconciliationReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reconciliation report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balances": balances})
}

// TransferPlatformFunds handles manual transfers between platform wallets
func (h *AdminPlatformHandler) TransferPlatformFunds(c *gin.Context) {
	adminID := c.GetString("user_id")

	var req struct {
		SourceID    string  `json:"source_id" binding:"required"`
		TargetID    string  `json:"target_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.platformService.TransferPlatformFunds(req.SourceID, req.TargetID, req.Amount, req.Description, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Transfer successful",
		"amount":    req.Amount,
		"source_id": req.SourceID,
		"target_id": req.TargetID,
	})
}
