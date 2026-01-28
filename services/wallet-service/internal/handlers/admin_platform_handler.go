package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminPlatformHandler struct {
	platformService *services.PlatformAccountService
}

func NewAdminPlatformHandler(platformService *services.PlatformAccountService) *AdminPlatformHandler {
	return &AdminPlatformHandler{platformService: platformService}
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

// SyncCryptoWalletBalance updates the balance from blockchain (placeholder)
func (h *AdminPlatformHandler) SyncCryptoWalletBalance(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement actual blockchain balance sync via Tatum or other provider
	// For now, just return success with current balance
	wallet, err := h.platformService.GetCryptoWalletByID(id)
	if err != nil || wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Balance sync not implemented yet",
		"wallet":  wallet,
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
