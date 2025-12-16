package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func SetupDashboardRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/summary", handleGetDashboardSummary(serviceManager))
	router.GET("/activity", handleGetRecentActivity(serviceManager))
}

func handleGetDashboardSummary(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/dashboard/summary?user_id="+userID.(string), nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			// Return empty dashboard data on error
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

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleGetRecentActivity(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")
		limit := c.DefaultQuery("limit", "10")

		resp, err := sm.CallService(c.Request.Context(), "wallet", "GET", "/dashboard/activity?user_id="+userID.(string)+"&limit="+limit, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			// Return empty activity on error
			c.JSON(http.StatusOK, gin.H{
				"activities": []interface{}{},
			})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}
