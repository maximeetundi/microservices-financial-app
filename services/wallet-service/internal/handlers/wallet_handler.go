package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService       *services.WalletService
	balanceService      *services.BalanceService
	systemConfigService *services.SystemConfigService
	cryptoService       *services.CryptoService // Injected
}

func NewWalletHandler(walletService *services.WalletService, balanceService *services.BalanceService, systemConfigService *services.SystemConfigService, cryptoService *services.CryptoService) *WalletHandler {
	return &WalletHandler{
		walletService:       walletService,
		balanceService:      balanceService,
		systemConfigService: systemConfigService,
		cryptoService:       cryptoService,
	}
}

func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	includeHiddenStr := c.DefaultQuery("include_hidden", "false")
	includeHidden, _ := strconv.ParseBool(includeHiddenStr)

	wallets, err := h.walletService.GetUserWallets(userID.(string), includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallets"})
		return
	}

	// Auto-create a default wallet if user has none (only if not getting hidden)
	if len(wallets) == 0 && !includeHidden {
		// Get user's country to determine currency
		userCountry := h.walletService.GetUserCountry(userID.(string))
		defaultCurrency := getCurrencyForCountry(userCountry)

		name := "Wallet Principal"
		desc := "Wallet créé automatiquement"
		defaultReq := &models.CreateWalletRequest{
			Currency:    defaultCurrency,
			WalletType:  "fiat",
			Name:        &name,
			Description: &desc,
		}

		newWallet, err := h.walletService.CreateWallet(userID.(string), defaultReq)
		if err == nil {
			wallets = append(wallets, newWallet)
		}
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (h *WalletHandler) GetPublicConfig(c *gin.Context) {
	isTestnet := h.systemConfigService.IsTestnetEnabled()
	c.JSON(http.StatusOK, gin.H{
		"testnet_enabled": isTestnet,
	})
}

func (h *WalletHandler) GetNetworks(c *gin.Context) {
	currency := c.Query("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency query parameter is required"})
		return
	}

	networks := h.cryptoService.GetAvailableNetworks(currency)
	c.JSON(http.StatusOK, gin.H{"networks": networks})
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")
	if walletID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet ID is required"})
		return
	}

	// Verify PIN
	pin := c.GetHeader("X-Transaction-Pin")
	if pin != "" {
		authClient := services.NewAuthClient()
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		valid, err := authClient.VerifyPin(userID.(string), pin, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN verification failed: " + err.Error()})
			return
		}
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
			return
		}
	} else {
		// PIN is required for deletion
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN is required for wallet deletion"})
		return
	}

	err := h.walletService.HideWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet hidden successfully"})
}

// getCurrencyForCountry returns the appropriate currency for a country code
func getCurrencyForCountry(countryCode string) string {
	// UEMOA - West African CFA Franc (XOF)
	xofCountries := map[string]bool{
		"BJ": true, "BF": true, "CI": true, "GW": true,
		"ML": true, "NE": true, "SN": true, "TG": true,
	}
	// CEMAC - Central African CFA Franc (XAF)
	xafCountries := map[string]bool{
		"CM": true, "CF": true, "TD": true, "CG": true, "GQ": true, "GA": true,
	}
	// Eurozone
	eurCountries := map[string]bool{
		"FR": true, "DE": true, "IT": true, "ES": true, "PT": true,
		"NL": true, "BE": true, "AT": true, "IE": true, "FI": true,
		"GR": true, "SK": true, "SI": true, "EE": true, "LV": true,
		"LT": true, "CY": true, "MT": true, "LU": true,
	}

	if countryCode == "GB" || countryCode == "UK" {
		return "GBP"
	}
	if countryCode == "US" || countryCode == "PR" {
		return "USD"
	}
	if countryCode == "CA" {
		return "CAD"
	}
	if xofCountries[countryCode] {
		return "XOF"
	}
	if xafCountries[countryCode] || countryCode == "Cameroun" || countryCode == "Cameroon" {
		return "XAF"
	}
	if eurCountries[countryCode] {
		return "EUR"
	}

	// Other African countries
	switch countryCode {
	case "MA":
		return "MAD"
	case "DZ":
		return "DZD"
	case "TN":
		return "TND"
	case "EG":
		return "EGP"
	case "NG":
		return "NGN"
	case "GH":
		return "GHS"
	case "KE":
		return "KES"
	case "ZA":
		return "ZAR"
	}

	// Default to USD for unknown countries
	return "USD"
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
		log.Printf("Error creating wallet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
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

// GetDashboardSummary returns summary statistics for the user's dashboard
func (h *WalletHandler) GetDashboardSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get all wallets for user
	wallets, err := h.walletService.GetUserWallets(userID.(string), false)
	if err != nil {
		// Return empty summary on error
		c.JSON(http.StatusOK, gin.H{
			"totalBalance":     0,
			"cryptoBalance":    0,
			"cardsBalance":     0,
			"activeCards":      0,
			"monthlyTransfers": 0,
			"monthlyVolume":    0,
		})
		return
	}

	// Calculate totals
	var totalBalance, cryptoBalance, fiatBalance float64
	for _, wallet := range wallets {
		if wallet.WalletType == "crypto" {
			cryptoBalance += wallet.Balance
		} else {
			fiatBalance += wallet.Balance
		}
		totalBalance += wallet.Balance
	}

	c.JSON(http.StatusOK, gin.H{
		"totalBalance":     totalBalance,
		"cryptoBalance":    cryptoBalance,
		"cardsBalance":     0, // TODO: Integrate with card service
		"activeCards":      0, // TODO: Integrate with card service
		"monthlyTransfers": 0, // TODO: Calculate from transactions
		"monthlyVolume":    0, // TODO: Calculate from transactions
	})
}

// GetRecentActivity returns recent transactions/activity for the user
func (h *WalletHandler) GetRecentActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}

	// Get user wallets first
	wallets, err := h.walletService.GetUserWallets(userID.(string), false)
	if err != nil || len(wallets) == 0 {
		c.JSON(http.StatusOK, gin.H{"activities": []interface{}{}})
		return
	}

	// Get transactions from all wallets
	var activities []map[string]interface{}
	for _, wallet := range wallets {
		transactions, err := h.walletService.GetTransactionHistory(wallet.ID, userID.(string), 5, 0, "", "")
		if err == nil {
			for _, tx := range transactions {
				description := ""
				if tx.Description != nil {
					description = *tx.Description
				}
				activities = append(activities, map[string]interface{}{
					"id":          tx.ID,
					"icon":        getTransactionIcon(tx.TransactionType),
					"title":       getTransactionTitle(tx.TransactionType),
					"description": description,
					"amount":      tx.Amount,
					"currency":    tx.Currency,
					"time":        tx.CreatedAt,
					"bgColor":     getTransactionBgColor(tx.TransactionType),
				})
			}
		}
		if len(activities) >= limit {
			break
		}
	}

	// Limit results
	if len(activities) > limit {
		activities = activities[:limit]
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

func getTransactionIcon(txType string) string {
	icons := map[string]string{
		"deposit":  "↓",
		"withdraw": "↑",
		"transfer": "💸",
		"exchange": "💱",
		"payment":  "💳",
	}
	if icon, ok := icons[txType]; ok {
		return icon
	}
	return "💰"
}

func getTransactionTitle(txType string) string {
	titles := map[string]string{
		"deposit":  "Dépôt",
		"withdraw": "Retrait",
		"transfer": "Transfert",
		"exchange": "Échange",
		"payment":  "Paiement",
	}
	if title, ok := titles[txType]; ok {
		return title
	}
	return "Transaction"
}

func getTransactionBgColor(txType string) string {
	colors := map[string]string{
		"deposit":  "bg-green-500",
		"withdraw": "bg-red-500",
		"transfer": "bg-purple-500",
		"exchange": "bg-blue-500",
		"payment":  "bg-orange-500",
	}
	if color, ok := colors[txType]; ok {
		return color
	}
	return "bg-gray-500"
}

// ProcessInterServiceTransaction handles transaction requests from other services
func (h *WalletHandler) ProcessInterServiceTransaction(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.walletService.ProcessTransaction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// CreditWalletFromDeposit handles credit requests from transfer-service after successful deposit
// POST /api/v1/wallets/credit
func (h *WalletHandler) CreditWalletFromDeposit(c *gin.Context) {
	var req struct {
		UserID         string  `json:"user_id" binding:"required"`
		WalletID       string  `json:"wallet_id"`
		Amount         float64 `json:"amount" binding:"required,gt=0"`
		Currency       string  `json:"currency" binding:"required"`
		TransactionRef string  `json:"transaction_ref" binding:"required"`
		Type           string  `json:"type"` // deposit, refund, etc.
		Description    string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	// If no wallet_id provided, find default wallet for currency
	walletID := req.WalletID
	if walletID == "" {
		wallet, err := h.walletService.GetWalletByCurrency(req.UserID, req.Currency)
		if err != nil || wallet == nil {
			// Create wallet if doesn't exist
			newWallet, err := h.walletService.CreateWallet(req.UserID, &models.CreateWalletRequest{
				Currency:   req.Currency,
				WalletType: "fiat",
			})
			if err != nil {
				log.Printf("[CreditWallet] Failed to create wallet for user %s: %v", req.UserID, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to find or create wallet",
					"success": false,
				})
				return
			}
			walletID = newWallet.ID
		} else {
			walletID = wallet.ID
		}
	}

	// Process the deposit
	description := req.Description
	if description == "" {
		description = "Dépôt via agrégateur de paiement"
	}

	err := h.walletService.ProcessDepositFromPlatform(
		req.UserID,
		walletID,
		req.Amount,
		req.Currency,
		req.TransactionRef,
		"aggregator",
	)

	if err != nil {
		log.Printf("[CreditWallet] Failed to credit wallet %s: %v", walletID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	// Get updated balance
	wallet, _ := h.walletService.GetWallet(walletID, req.UserID)
	var newBalance float64
	if wallet != nil {
		newBalance = wallet.Balance
	}

	log.Printf("[CreditWallet] ✅ Credited %.2f %s to wallet %s (user: %s, ref: %s)",
		req.Amount, req.Currency, walletID, req.UserID, req.TransactionRef)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"wallet_id":   walletID,
		"new_balance": newBalance,
		"amount":      req.Amount,
		"currency":    req.Currency,
	})
}

// ProcessPlatformDeposit handles deposit requests coming from platform (aggregators)
// It performs double-entry bookkeeping (Debit Platform Reserve -> Credit User Wallet)
func (h *WalletHandler) ProcessPlatformDeposit(c *gin.Context) {
	var req struct {
		UserID       string  `json:"user_id" binding:"required"`
		WalletID     string  `json:"wallet_id" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
		Currency     string  `json:"currency" binding:"required"`
		ProviderRef  string  `json:"provider_ref" binding:"required"`
		ProviderName string  `json:"provider_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.walletService.ProcessDepositFromPlatform(req.UserID, req.WalletID, req.Amount, req.Currency, req.ProviderRef, req.ProviderName)
	if err != nil {
		// Log the error
		log.Printf("Failed to process platform deposit: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	var req struct {
		Amount   float64 `json:"amount" binding:"required,gt=0"`
		Method   string  `json:"method" binding:"required"` // Provider name: demo, flutterwave, stripe, etc.
		Provider string  `json:"provider"`                  // Alias for method
		Country  string  `json:"country"`                   // User's country for routing
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Support both 'method' and 'provider' for backwards compatibility
	provider := req.Method
	if req.Provider != "" {
		provider = req.Provider
	}

	// Verify wallet exists and belongs to user
	wallet, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Check if this is a demo provider (instant credit)
	if provider == "demo" || provider == "Mode Démo" {
		// DEMO MODE: Instant credit using ProcessDepositFromPlatform
		transactionID := "dep_demo_" + strconv.FormatInt(time.Now().Unix(), 10)

		err = h.walletService.ProcessDepositFromPlatform(
			userID.(string),
			walletID,
			req.Amount,
			wallet.Currency,
			transactionID,
			"demo",
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process deposit: " + err.Error()})
			return
		}

		// Get updated balance
		updatedWallet, _ := h.walletService.GetWallet(walletID, userID.(string))
		newBalance := wallet.Balance + req.Amount
		if updatedWallet != nil {
			newBalance = updatedWallet.Balance
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Deposit successful (demo mode)",
			"transaction_id": transactionID,
			"status":         "completed",
			"amount":         req.Amount,
			"new_balance":    newBalance,
			"currency":       wallet.Currency,
		})
		return
	}

	// REAL PROVIDER MODE: Call transfer-service to initiate payment
	transferServiceURL := os.Getenv("TRANSFER_SERVICE_URL")
	if transferServiceURL == "" {
		transferServiceURL = "http://transfer-service:8084"
	}

	// Prepare request to transfer-service
	depositReq := map[string]interface{}{
		"user_id":    userID.(string),
		"amount":     req.Amount,
		"currency":   wallet.Currency,
		"provider":   provider,
		"country":    req.Country,
		"return_url": c.Request.Header.Get("Origin") + "/wallet?deposit=success",
	}

	reqBody, _ := json.Marshal(depositReq)

	// Create HTTP request with Authorization header forwarded
	httpReq, err := http.NewRequest(
		"POST",
		transferServiceURL+"/api/v1/deposits/initiate",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		log.Printf("Failed to create request to transfer-service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate payment"})
		return
	}

	// Set headers - forward Authorization from original request
	httpReq.Header.Set("Content-Type", "application/json")
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		httpReq.Header.Set("Authorization", authHeader)
	}

	// Execute the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)

	if err != nil {
		log.Printf("Failed to call transfer-service: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Payment service temporarily unavailable",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var transferResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&transferResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, transferResp)
		return
	}

	// Return payment URL to frontend
	c.JSON(http.StatusOK, gin.H{
		"status":         transferResp["status"],
		"transaction_id": transferResp["transaction_id"],
		"payment_url":    transferResp["payment_url"],
		"provider":       provider,
		"amount":         req.Amount,
		"currency":       wallet.Currency,
		"message":        "Please complete payment on the provider's page",
	})
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	walletID := c.Param("wallet_id")

	// Support both standard withdraw and ticket purchase withdraw
	var req struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Destination string  `json:"destination"` // Optional for ticket purchase
		Description string  `json:"description"` // For ticket purchase
		PIN         string  `json:"pin"`         // PIN for verification
		Currency    string  `json:"currency"`    // Target currency for conversion
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify wallet exists and belongs to user
	wallet, err := h.walletService.GetWallet(walletID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Check if wallet is active
	if !wallet.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wallet is not active"})
		return
	}

	// Verify PIN if provided
	if req.PIN != "" {
		authClient := services.NewAuthClient()
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		valid, err := authClient.VerifyPin(userID.(string), req.PIN, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN verification failed: " + err.Error()})
			return
		}
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
			return
		}
	}

	// Calculate amount to debit with currency conversion if needed
	amountToDebit := req.Amount
	if req.Currency != "" && req.Currency != wallet.Currency {
		exchangeClient := services.NewExchangeClient()
		convertedAmount, err := exchangeClient.ConvertAmount(req.Amount, req.Currency, wallet.Currency)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Currency conversion failed: " + err.Error()})
			return
		}
		amountToDebit = convertedAmount
	}

	// Check if sufficient balance
	if wallet.Balance < amountToDebit {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Insufficient balance",
			"available": wallet.Balance,
			"required":  amountToDebit,
			"currency":  wallet.Currency,
		})
		return
	}

	// Perform the actual withdrawal/debit
	err = h.balanceService.UpdateBalance(walletID, amountToDebit, "withdraw")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process withdrawal: " + err.Error()})
		return
	}

	// Generate transaction ID
	transactionID := "txn_" + strconv.FormatInt(time.Now().UnixNano(), 36)

	// Calculate new balance
	newBalance := wallet.Balance - amountToDebit

	c.JSON(http.StatusOK, gin.H{
		"message":        "Withdrawal successful",
		"transaction_id": transactionID,
		"new_balance":    newBalance,
	})
}

func (h *WalletHandler) HandleTatumWebhook(c *gin.Context) {
	// Tatum Webhook Payload Structure (Simplified)
	var req struct {
		AccountId        string `json:"accountId"`
		Amount           string `json:"amount"` // String to avoid float precision issues initially
		SubscriptionType string `json:"subscriptionType"`
		TxId             string `json:"txId"`
		Currency         string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// We only care about ACCOUNT_INCOMING_BLOCKCHAIN_TRANSACTION
	if req.SubscriptionType != "ACCOUNT_INCOMING_BLOCKCHAIN_TRANSACTION" {
		c.JSON(http.StatusOK, gin.H{"message": "Ignored subscription type"})
		return
	}

	amountFloat, _ := strconv.ParseFloat(req.Amount, 64)

	err := h.walletService.ProcessTatumDeposit(req.AccountId, amountFloat, req.Currency, req.TxId)
	if err != nil {
		// Log error but maybe return 200 to prevent Tatum retries if it's a logic error?
		// Or 500 to retry?
		// Return 400/500 to alert.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit processed"})
}

// GetPortfolio returns the user's portfolio with asset allocation
func (h *WalletHandler) GetPortfolio(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get all wallets for user
	wallets, err := h.walletService.GetUserWallets(userID.(string), false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"assets":      []interface{}{},
			"total_value": 0,
			"change_24h":  0,
			"change_pct":  0,
		})
		return
	}

	// Build portfolio data
	var assets []map[string]interface{}
	var totalValue float64

	for _, wallet := range wallets {
		asset := map[string]interface{}{
			"id":         wallet.ID,
			"name":       wallet.Currency,
			"symbol":     wallet.Currency,
			"type":       wallet.WalletType,
			"balance":    wallet.Balance,
			"value":      wallet.Balance, // TODO: Convert to base currency
			"change_24h": 0.0,            // TODO: Get from exchange service
			"change_pct": 0.0,            // TODO: Get from exchange service
			"allocation": 0.0,            // Will be calculated below
		}
		assets = append(assets, asset)
		totalValue += wallet.Balance
	}

	// Calculate allocation percentages
	for i := range assets {
		if totalValue > 0 {
			assets[i]["allocation"] = (assets[i]["value"].(float64) / totalValue) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"assets":      assets,
		"total_value": totalValue,
		"change_24h":  0,   // TODO: Calculate from exchange rates
		"change_pct":  0.0, // TODO: Calculate from exchange rates
	})
}

// GetStats returns statistics for a given period
func (h *WalletHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	period := c.DefaultQuery("period", "month")

	// Get all wallets for user
	wallets, err := h.walletService.GetUserWallets(userID.(string), false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"period":            period,
			"total_deposits":    0,
			"total_withdrawals": 0,
			"total_transfers":   0,
			"total_exchanges":   0,
			"transaction_count": 0,
			"volume":            0,
		})
		return
	}

	// Calculate period stats
	var totalDeposits, totalWithdrawals, totalTransfers, totalExchanges float64
	var transactionCount int

	for _, wallet := range wallets {
		// Get transactions for this wallet in the period
		transactions, err := h.walletService.GetTransactionHistory(wallet.ID, userID.(string), 100, 0, "", "")
		if err == nil {
			for _, tx := range transactions {
				transactionCount++
				switch tx.TransactionType {
				case "deposit":
					totalDeposits += tx.Amount
				case "withdraw":
					totalWithdrawals += tx.Amount
				case "transfer":
					totalTransfers += tx.Amount
				case "exchange":
					totalExchanges += tx.Amount
				}
			}
		}
	}

	volume := totalDeposits + totalWithdrawals + totalTransfers + totalExchanges

	c.JSON(http.StatusOK, gin.H{
		"period":            period,
		"total_deposits":    totalDeposits,
		"total_withdrawals": totalWithdrawals,
		"total_transfers":   totalTransfers,
		"total_exchanges":   totalExchanges,
		"transaction_count": transactionCount,
		"volume":            volume,
	})
}

// === Internal Service-to-Service Handlers ===

func (h *WalletHandler) GetWalletsInternal(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	wallets, err := h.walletService.GetUserWallets(userID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (h *WalletHandler) CreateWalletInternal(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := "Portefeuille " + req.Currency
	desc := "Auto-generated wallet"
	createReq := &models.CreateWalletRequest{
		Currency:    req.Currency,
		WalletType:  "fiat", // Default to fiat for these auto-creations
		Name:        &name,
		Description: &desc,
	}

	wallet, err := h.walletService.CreateWallet(req.UserID, createReq)
	if err != nil {
		if err.Error() == "wallet already exists for currency "+req.Currency {
			// Try to find existing wallet
			wallets, _ := h.walletService.GetUserWallets(req.UserID, true)
			for _, w := range wallets {
				if w.Currency == req.Currency {
					c.JSON(http.StatusOK, gin.H{"wallet": w, "existing": true})
					return
				}
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"wallet": wallet})
}
