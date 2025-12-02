package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
)

func SetupTransferRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.POST("/", middleware.KYCRequired(1), handleCreateTransfer(serviceManager))
	router.GET("/", handleGetTransferHistory(serviceManager))
	router.GET("/:transfer_id", handleGetTransfer(serviceManager))
	router.POST("/:transfer_id/cancel", handleCancelTransfer(serviceManager))
	
	// International transfers require KYC level 2
	router.POST("/international", middleware.KYCRequired(2), handleInternationalTransfer(serviceManager))
	
	// Mobile money transfers
	mobile := router.Group("/mobile")
	{
		mobile.POST("/send", middleware.KYCRequired(1), handleMobileMoneySend(serviceManager))
		mobile.POST("/receive", handleMobileMoneyReceive(serviceManager))
		mobile.GET("/providers", handleGetMobileProviders(serviceManager))
	}
	
	// Bulk transfers (for businesses)
	bulk := router.Group("/bulk")
	bulk.Use(middleware.KYCRequired(3))
	{
		bulk.POST("/", handleBulkTransfer(serviceManager))
		bulk.GET("/:batch_id", handleGetBulkTransferStatus(serviceManager))
		bulk.POST("/:batch_id/approve", handleApproveBulkTransfer(serviceManager))
	}
}

func handleCreateTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID string  `json:"from_wallet_id" binding:"required"`
			ToWalletID   string  `json:"to_wallet_id,omitempty"`
			ToEmail      string  `json:"to_email,omitempty"`
			ToPhone      string  `json:"to_phone,omitempty"`
			Amount       float64 `json:"amount" binding:"required,gt=0"`
			Currency     string  `json:"currency" binding:"required"`
			Description  string  `json:"description,omitempty"`
			Reference    string  `json:"reference,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate recipient
		if req.ToWalletID == "" && req.ToEmail == "" && req.ToPhone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Recipient required (wallet_id, email, or phone)"})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		transferData := map[string]interface{}{
			"user_id":        userID.(string),
			"from_wallet_id": req.FromWalletID,
			"amount":         req.Amount,
			"currency":       req.Currency,
			"description":    req.Description,
			"reference":      req.Reference,
		}

		if req.ToWalletID != "" {
			transferData["to_wallet_id"] = req.ToWalletID
		}
		if req.ToEmail != "" {
			transferData["to_email"] = req.ToEmail
		}
		if req.ToPhone != "" {
			transferData["to_phone"] = req.ToPhone
		}

		resp, err := sm.CreateTransfer(c.Request.Context(), transferData, extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetTransferHistory(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		// Query parameters
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")
		status := c.Query("status")
		currency := c.Query("currency")
		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")

		endpoint := "/transfers?user_id=" + userID.(string) + "&limit=" + limit + "&offset=" + offset
		if status != "" {
			endpoint += "&status=" + status
		}
		if currency != "" {
			endpoint += "&currency=" + currency
		}
		if dateFrom != "" {
			endpoint += "&date_from=" + dateFrom
		}
		if dateTo != "" {
			endpoint += "&date_to=" + dateTo
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "GET", endpoint, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		transferID := c.Param("transfer_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "transfer", "GET", "/transfers/"+transferID, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleCancelTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		transferID := c.Param("transfer_id")
		var req struct {
			Reason string `json:"reason" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/transfers/"+transferID+"/cancel", req, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleInternationalTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID      string  `json:"from_wallet_id" binding:"required"`
			RecipientName     string  `json:"recipient_name" binding:"required"`
			RecipientAddress  string  `json:"recipient_address" binding:"required"`
			RecipientCountry  string  `json:"recipient_country" binding:"required,len=3"`
			BankName          string  `json:"bank_name" binding:"required"`
			BankAccount       string  `json:"bank_account" binding:"required"`
			SwiftCode         string  `json:"swift_code,omitempty"`
			RoutingNumber     string  `json:"routing_number,omitempty"`
			Amount            float64 `json:"amount" binding:"required,gt=0"`
			Currency          string  `json:"currency" binding:"required"`
			Purpose           string  `json:"purpose" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		transferData := map[string]interface{}{
			"user_id":            userID.(string),
			"from_wallet_id":     req.FromWalletID,
			"recipient_name":     req.RecipientName,
			"recipient_address":  req.RecipientAddress,
			"recipient_country":  req.RecipientCountry,
			"bank_name":          req.BankName,
			"bank_account":       req.BankAccount,
			"swift_code":         req.SwiftCode,
			"routing_number":     req.RoutingNumber,
			"amount":             req.Amount,
			"currency":           req.Currency,
			"purpose":            req.Purpose,
			"transfer_type":      "international",
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/international", transferData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleMobileMoneySend(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID string  `json:"from_wallet_id" binding:"required"`
			ToPhone      string  `json:"to_phone" binding:"required"`
			Provider     string  `json:"provider" binding:"required"`
			Amount       float64 `json:"amount" binding:"required,gt=0"`
			Currency     string  `json:"currency" binding:"required"`
			Country      string  `json:"country" binding:"required,len=3"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		transferData := map[string]interface{}{
			"user_id":        userID.(string),
			"from_wallet_id": req.FromWalletID,
			"to_phone":       req.ToPhone,
			"provider":       req.Provider,
			"amount":         req.Amount,
			"currency":       req.Currency,
			"country":        req.Country,
			"transfer_type":  "mobile_money",
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/mobile/send", transferData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleMobileMoneyReceive(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ToWalletID     string  `json:"to_wallet_id" binding:"required"`
			FromPhone      string  `json:"from_phone" binding:"required"`
			Provider       string  `json:"provider" binding:"required"`
			Amount         float64 `json:"amount" binding:"required,gt=0"`
			Currency       string  `json:"currency" binding:"required"`
			TransactionRef string  `json:"transaction_ref" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		receiveData := map[string]interface{}{
			"user_id":         userID.(string),
			"to_wallet_id":    req.ToWalletID,
			"from_phone":      req.FromPhone,
			"provider":        req.Provider,
			"amount":          req.Amount,
			"currency":        req.Currency,
			"transaction_ref": req.TransactionRef,
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/mobile/receive", receiveData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetMobileProviders(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		country := c.Query("country")
		
		endpoint := "/mobile/providers"
		if country != "" {
			endpoint += "?country=" + country
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "GET", endpoint, nil, nil)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleBulkTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID string `json:"from_wallet_id" binding:"required"`
			Transfers    []struct {
				ToEmail     string  `json:"to_email,omitempty"`
				ToPhone     string  `json:"to_phone,omitempty"`
				ToWalletID  string  `json:"to_wallet_id,omitempty"`
				Amount      float64 `json:"amount" binding:"required,gt=0"`
				Reference   string  `json:"reference,omitempty"`
				Description string  `json:"description,omitempty"`
			} `json:"transfers" binding:"required,min=1,max=1000"`
			Currency    string `json:"currency" binding:"required"`
			Description string `json:"description,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		bulkData := map[string]interface{}{
			"user_id":        userID.(string),
			"from_wallet_id": req.FromWalletID,
			"transfers":      req.Transfers,
			"currency":       req.Currency,
			"description":    req.Description,
		}

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/bulk", bulkData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetBulkTransferStatus(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		batchID := c.Param("batch_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "transfer", "GET", "/bulk/"+batchID, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleApproveBulkTransfer(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		batchID := c.Param("batch_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "transfer", "POST", "/bulk/"+batchID+"/approve", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}