package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/wallet-service/internal/models"
	"github.com/crypto-bank/wallet-service/internal/services"
)

type WalletHandler struct {
	walletService  *services.WalletService
	balanceService *services.BalanceService
}

func NewWalletHandler(walletService *services.WalletService, balanceService *services.BalanceService) *WalletHandler {
	return &WalletHandler{
		walletService:  walletService,
		balanceService: balanceService,
	}
}

func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	wallets, err := h.walletService.GetUserWallets(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.CreateWallet(userID.(string), &req)
	if err != nil {
		if err.Error() == "wallet already exists for currency "+req.Currency {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"wallet": wallet})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	wallet, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

func (h *WalletHandler) UpdateWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	
	// First verify ownership
	_, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	var updateReq map[string]interface{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, only allow updating the name
	if name, ok := updateReq["name"]; ok {
		// TODO: Implement wallet update in service
		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated", "name": name})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	
	// Verify ownership
	_, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	balance, err := h.balanceService.GetBalance(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (h *WalletHandler) GetWalletTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	txType := c.Query("type")

	transactions, err := h.walletService.GetTransactionHistory(walletID, userID.(string), limit, offset, status, txType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *WalletHandler) FreezeWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.walletService.FreezeWallet(walletID, userID.(string), req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to freeze wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet frozen successfully"})
}

func (h *WalletHandler) UnfreezeWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	err := h.walletService.UnfreezeWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfreeze wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet unfrozen successfully"})
}

func (h *WalletHandler) GenerateCryptoWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Force wallet type to crypto
	req.WalletType = "crypto"

	wallet, err := h.walletService.CreateWallet(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate crypto wallet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"wallet": wallet})
}

func (h *WalletHandler) GetCryptoAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	address, err := h.walletService.GetCryptoAddress(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get crypto address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

func (h *WalletHandler) SendCrypto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	var req models.SendCryptoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.walletService.SendCrypto(walletID, userID.(string), &req)
	if err != nil {
		if err.Error() == "insufficient balance" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send crypto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

func (h *WalletHandler) GetPendingTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	transactions, err := h.walletService.GetPendingTransactions(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *WalletHandler) EstimateTransactionFee(c *gin.Context) {
	walletID := c.Param("wallet_id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify wallet ownership
	wallet, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	var req struct {
		Amount   float64 `json:"amount" binding:"required,gt=0"`
		Priority string  `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = "normal"
	}

	estimate, err := h.walletService.EstimateTransactionFee(wallet.Currency, req.Amount, req.Priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to estimate fee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"estimate": estimate})
}

func (h *WalletHandler) HandleCryptoConfirmation(c *gin.Context) {
	var req models.BlockchainConfirmation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the confirmation
	// This would typically update transaction status based on confirmations
	// For now, just acknowledge receipt

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation processed"})
}

func (h *WalletHandler) HandleCryptoDeposit(c *gin.Context) {
	var req struct {
		Address  string  `json:"address" binding:"required"`
		TxHash   string  `json:"tx_hash" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.walletService.ProcessCryptoDeposit(req.Address, req.TxHash, req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit processed"})
}