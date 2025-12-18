package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func SetupAdminRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/users", handleAdminGetUsers(serviceManager))
	router.GET("/users/:user_id", handleAdminGetUser(serviceManager))
	router.PUT("/users/:user_id/kyc", handleAdminUpdateKYC(serviceManager))
	router.GET("/transactions", handleAdminGetTransactions(serviceManager))
	router.POST("/transactions/:tx_id/investigate", handleInvestigateTransaction(serviceManager))
	router.GET("/analytics", handleGetAnalytics(serviceManager))
}

func handleAdminGetUsers(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		limit := c.DefaultQuery("limit", "50")
		offset := c.DefaultQuery("offset", "0")

		resp, err := sm.CallService(c.Request.Context(), "user", "GET", "/api/v1/admin/users?limit="+limit+"&offset="+offset, nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
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
