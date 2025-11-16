package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/api-gateway/internal/services"
	"github.com/crypto-bank/api-gateway/internal/middleware"
)

func SetupExchangeRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/rates", handleGetExchangeRates(serviceManager))
	router.GET("/rates/:from/:to", handleGetSpecificRate(serviceManager))
	router.POST("/quote", middleware.KYCRequired(1), handleGetQuote(serviceManager))
	router.POST("/execute", middleware.KYCRequired(1), handleExecuteExchange(serviceManager))
	router.GET("/history", handleGetExchangeHistory(serviceManager))
	router.GET("/:exchange_id", handleGetExchange(serviceManager))
	
	// Advanced trading features
	trading := router.Group("/trading")
	trading.Use(middleware.KYCRequired(2))
	{
		trading.POST("/limit-order", handleCreateLimitOrder(serviceManager))
		trading.POST("/stop-loss", handleCreateStopLoss(serviceManager))
		trading.GET("/orders", handleGetOrders(serviceManager))
		trading.DELETE("/orders/:order_id", handleCancelOrder(serviceManager))
		trading.GET("/orderbook/:pair", handleGetOrderBook(serviceManager))
	}
}

func SetupUserRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/profile", handleGetProfile(serviceManager))
	router.PUT("/profile", handleUpdateProfile(serviceManager))
	router.GET("/kyc", handleGetKYCStatus(serviceManager))
	router.POST("/kyc/upload", handleUploadKYCDocument(serviceManager))
	router.GET("/settings", handleGetSettings(serviceManager))
	router.PUT("/settings", handleUpdateSettings(serviceManager))
}

func SetupCardRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/", handleGetCards(serviceManager))
	router.POST("/", middleware.KYCRequired(2), handleCreateCard(serviceManager))
	router.GET("/:card_id", handleGetCard(serviceManager))
	router.PUT("/:card_id", handleUpdateCard(serviceManager))
	router.POST("/:card_id/activate", handleActivateCard(serviceManager))
	router.POST("/:card_id/freeze", handleFreezeCard(serviceManager))
	router.POST("/:card_id/load", handleLoadCard(serviceManager))
	router.GET("/:card_id/transactions", handleGetCardTransactions(serviceManager))
}

func SetupNotificationRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/", handleGetNotifications(serviceManager))
	router.PUT("/:notification_id/read", handleMarkAsRead(serviceManager))
	router.DELETE("/:notification_id", handleDeleteNotification(serviceManager))
	router.POST("/settings", handleUpdateNotificationSettings(serviceManager))
}

func SetupAdminRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/users", handleAdminGetUsers(serviceManager))
	router.GET("/users/:user_id", handleAdminGetUser(serviceManager))
	router.PUT("/users/:user_id/kyc", handleAdminUpdateKYC(serviceManager))
	router.GET("/transactions", handleAdminGetTransactions(serviceManager))
	router.POST("/transactions/:tx_id/investigate", handleInvestigateTransaction(serviceManager))
	router.GET("/analytics", handleGetAnalytics(serviceManager))
}

// Exchange handlers
func handleGetExchangeRates(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := sm.GetExchangeRates(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetSpecificRate(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		
		resp, err := sm.CallService(c.Request.Context(), "exchange", "GET", "/rates/"+from+"/"+to, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetQuote(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromCurrency string  `json:"from_currency" binding:"required"`
			ToCurrency   string  `json:"to_currency" binding:"required"`
			Amount       float64 `json:"amount" binding:"required,gt=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")
		resp, err := sm.CallService(c.Request.Context(), "exchange", "POST", "/quote", req, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleExecuteExchange(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromWalletID string  `json:"from_wallet_id" binding:"required"`
			ToWalletID   string  `json:"to_wallet_id" binding:"required"`
			Amount       float64 `json:"amount" binding:"required,gt=0"`
			QuoteID      string  `json:"quote_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		exchangeData := map[string]interface{}{
			"user_id":        userID.(string),
			"from_wallet_id": req.FromWalletID,
			"to_wallet_id":   req.ToWalletID,
			"amount":         req.Amount,
			"quote_id":       req.QuoteID,
		}

		resp, err := sm.CreateExchange(c.Request.Context(), exchangeData, extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetExchangeHistory(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		endpoint := "/exchange/history?user_id=" + userID.(string)
		resp, err := sm.CallService(c.Request.Context(), "exchange", "GET", endpoint, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetExchange(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		exchangeID := c.Param("exchange_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "exchange", "GET", "/exchange/"+exchangeID, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// User handlers
func handleGetProfile(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.GetUser(c.Request.Context(), userID.(string), extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleUpdateProfile(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		var userData map[string]interface{}
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.UpdateUser(c.Request.Context(), userID.(string), userData, extractBearerToken(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Card handlers
func handleGetCards(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "card", "GET", "/cards?user_id="+userID.(string), nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleCreateCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CardType     string  `json:"card_type" binding:"required,oneof=prepaid virtual"`
			Currency     string  `json:"currency" binding:"required"`
			DailyLimit   float64 `json:"daily_limit,omitempty"`
			MonthlyLimit float64 `json:"monthly_limit,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		cardData := map[string]interface{}{
			"user_id":   userID.(string),
			"card_type": req.CardType,
			"currency":  req.Currency,
		}

		if req.DailyLimit > 0 {
			cardData["daily_limit"] = req.DailyLimit
		}
		if req.MonthlyLimit > 0 {
			cardData["monthly_limit"] = req.MonthlyLimit
		}

		resp, err := sm.CallService(c.Request.Context(), "card", "POST", "/cards", cardData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Trading handlers
func handleCreateLimitOrder(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			FromCurrency string  `json:"from_currency" binding:"required"`
			ToCurrency   string  `json:"to_currency" binding:"required"`
			Amount       float64 `json:"amount" binding:"required,gt=0"`
			LimitPrice   float64 `json:"limit_price" binding:"required,gt=0"`
			OrderType    string  `json:"order_type" binding:"required,oneof=buy sell"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		orderData := map[string]interface{}{
			"user_id":       userID.(string),
			"from_currency": req.FromCurrency,
			"to_currency":   req.ToCurrency,
			"amount":        req.Amount,
			"limit_price":   req.LimitPrice,
			"order_type":    req.OrderType,
			"execution_type": "limit",
		}

		resp, err := sm.CallService(c.Request.Context(), "exchange", "POST", "/trading/limit-order", orderData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Notification handlers
func handleGetNotifications(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "notification", "GET", "/notifications?user_id="+userID.(string), nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Admin handlers
func handleAdminGetUsers(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")

		resp, err := sm.CallService(c.Request.Context(), "user", "GET", "/admin/users?limit="+limit+"&offset="+offset, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Placeholder handlers for other routes
func handleGetKYCStatus(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleUploadKYCDocument(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetSettings(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleUpdateSettings(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleUpdateCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleActivateCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleFreezeCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleLoadCard(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetCardTransactions(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleMarkAsRead(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleDeleteNotification(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleUpdateNotificationSettings(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleAdminGetUser(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleAdminUpdateKYC(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleAdminGetTransactions(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleInvestigateTransaction(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetAnalytics(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleCreateStopLoss(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetOrders(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleCancelOrder(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}
func handleGetOrderBook(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
}