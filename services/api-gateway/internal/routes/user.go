package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func SetupUserRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/profile", handleGetProfile(serviceManager))
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
	}
}
