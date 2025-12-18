package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
	"net/http"
)

func SetupNotificationRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/", handleGetNotifications(serviceManager))
	router.PUT("/:notification_id/read", handleMarkAsRead(serviceManager))
	router.DELETE("/:notification_id", handleDeleteNotification(serviceManager))
	router.POST("/settings", handleUpdateNotificationSettings(serviceManager))
}

func handleGetNotifications(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		token := c.GetHeader("Authorization")

		resp, err := sm.CallService(c.Request.Context(), "notification", "GET", "/api/v1/notifications?user_id="+userID.(string), nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
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
