package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
)

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

func handleGetCards(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "card", "GET", "/api/v1/cards?user_id="+userID.(string), nil, map[string]string{
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

		resp, err := sm.CallService(c.Request.Context(), "card", "POST", "/api/v1/cards", cardData, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
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
