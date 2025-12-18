package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func SetupUserRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/profile", handleGetProfile(serviceManager))
	router.PUT("/profile", handleUpdateProfile(serviceManager))
	router.GET("/kyc", handleGetKYCStatus(serviceManager))
	router.POST("/kyc/upload", handleUploadKYCDocument(serviceManager))
	router.GET("/settings", handleGetSettings(serviceManager))
	router.PUT("/settings", handleUpdateSettings(serviceManager))
	router.GET("/lookup", handleLookupUser(serviceManager))
}

func handleGetProfile(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proxy to auth service (acting as user service)
		resp, err := sm.CallService(c.Request.Context(), "auth", "GET", "/api/v1/users/profile", nil, map[string]string{
			"Authorization": c.GetHeader("Authorization"),
		})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleLookupUser(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		phone := c.Query("phone")

		params := ""
		if email != "" {
			params = "?email=" + email
		} else if phone != "" {
			params = "?phone=" + phone
		}

		// Proxy to auth service
		resp, err := sm.CallService(c.Request.Context(), "auth", "GET", "/api/v1/users/lookup"+params, nil, map[string]string{
			"Authorization": c.GetHeader("Authorization"),
		})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
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

func extractBearerToken(authorization string) string {
	if len(authorization) > 7 && authorization[:7] == "Bearer " {
		return authorization[7:]
	}
	return ""
}
