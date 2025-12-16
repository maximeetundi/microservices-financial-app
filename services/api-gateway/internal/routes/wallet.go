package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
)

func SetupWalletRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("", handleGetWallets(serviceManager))
	router.POST("", handleCreateWallet(serviceManager))
	router.GET("/:wallet_id", handleGetWallet(serviceManager))
	router.PUT("/:wallet_id", handleUpdateWallet(serviceManager))
	router.GET("/:wallet_id/balance", handleGetBalance(serviceManager))
	router.GET("/:wallet_id/transactions", handleGetWalletTransactions(serviceManager))
	
	// High-value operations require KYC level 2
	router.POST("/:wallet_id/freeze", middleware.KYCRequired(2), handleFreezeWallet(serviceManager))
	router.POST("/:wallet_id/unfreeze", middleware.KYCRequired(2), handleUnfreezeWallet(serviceManager))
	
	// Crypto wallet specific
	crypto := router.Group("/crypto")
	{
		crypto.POST("/generate", middleware.KYCRequired(1), handleGenerateCryptoWallet(serviceManager))
		crypto.GET("/:wallet_id/address", handleGetCryptoAddress(serviceManager))
		crypto.POST("/:wallet_id/send", middleware.KYCRequired(2), handleSendCrypto(serviceManager))
		crypto.GET("/:wallet_id/pending", handleGetPendingTransactions(serviceManager))
	}
}

func handleGetWallets(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.GetWallets(c.Request.Context(), userID.(string), extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleCreateWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Currency    string `json:"currency" binding:"required"`
			WalletType  string `json:"wallet_type" binding:"required,oneof=fiat crypto"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		walletData := map[string]interface{}{
			"user_id":     userID.(string),
			"currency":    req.Currency,
			"wallet_type": req.WalletType,
			"name":        req.Name,
			"description": req.Description,
		}

		resp, err := sm.CreateWallet(c.Request.Context(), walletData, extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/wallets/"+walletID, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleUpdateWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		var req map[string]interface{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "PUT", "/wallets/"+walletID, req, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetBalance(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/wallets/"+walletID+"/balance", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetWalletTransactions(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		// Query parameters
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")
		status := c.Query("status")
		txType := c.Query("type")

		endpoint := "/wallets/" + walletID + "/transactions?limit=" + limit + "&offset=" + offset
		if status != "" {
			endpoint += "&status=" + status
		}
		if txType != "" {
			endpoint += "&type=" + txType
		}

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", endpoint, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleFreezeWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		var req struct {
			Reason string `json:"reason" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "POST", "/wallets/"+walletID+"/freeze", req, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleUnfreezeWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "POST", "/wallets/"+walletID+"/unfreeze", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGenerateCryptoWallet(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Currency string `json:"currency" binding:"required,oneof=BTC ETH LTC BCH"`
			Name     string `json:"name"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		walletData := map[string]interface{}{
			"user_id":     userID.(string),
			"currency":    req.Currency,
			"wallet_type": "crypto",
			"name":        req.Name,
		}

		resp, err := sm.CallService(c.Request.Context(), "wallet", "POST", "/crypto/generate", walletData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetCryptoAddress(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/crypto/"+walletID+"/address", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleSendCrypto(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		var req struct {
			ToAddress string  `json:"to_address" binding:"required"`
			Amount    float64 `json:"amount" binding:"required,gt=0"`
			GasPrice  *int64  `json:"gas_price,omitempty"`
			Note      string  `json:"note,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")
		userID, _ := c.Get("user_id")

		sendData := map[string]interface{}{
			"wallet_id":  walletID,
			"user_id":    userID.(string),
			"to_address": req.ToAddress,
			"amount":     req.Amount,
			"note":       req.Note,
		}

		if req.GasPrice != nil {
			sendData["gas_price"] = *req.GasPrice
		}

		resp, err := sm.CallService(c.Request.Context(), "wallet", "POST", "/crypto/"+walletID+"/send", sendData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetPendingTransactions(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("wallet_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/crypto/"+walletID+"/pending", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func extractBearerToken(authorization string) string {
	if len(authorization) > 7 && authorization[:7] == "Bearer " {
		return authorization[7:]
	}
	return ""
}