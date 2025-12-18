package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
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

// SetupUserRoutes removed (moved to user.go)
// SetupCardRoutes removed (moved to card.go - if exists, or should be cleaned up)
// SetupNotificationRoutes removed
// SetupAdminRoutes removed

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
		
		resp, err := sm.CallService(c.Request.Context(), "exchange", "GET", "/api/v1/rates/"+from+"/"+to, nil, nil)
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
		resp, err := sm.CallService(c.Request.Context(), "exchange", "POST", "/api/v1/quote", req, map[string]string{
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

		endpoint := "/api/v1/exchange/history?user_id=" + userID.(string)
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

		resp, err := sm.CallService(c.Request.Context(), "exchange", "GET", "/api/v1/exchange/"+exchangeID, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// User and Card handlers removed (moved to respective files)

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

		resp, err := sm.CallService(c.Request.Context(), "exchange", "POST", "/api/v1/trading/limit-order", orderData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

// Notification handlers (moved to notification.go)

// Admin handlers (moved to admin.go)

// Placeholder handlers for other routes
// User specific placeholders moved to user.go, card to card.go etc.
// Keeping only Exchange specific placeholders if any.

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
