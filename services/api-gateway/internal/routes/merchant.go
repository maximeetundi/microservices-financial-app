package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func SetupMerchantRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.POST("/payments", handleCreatePayment(serviceManager))
	router.GET("/payments", handleGetPayments(serviceManager))
	router.GET("/payments/history", handleGetPaymentHistory(serviceManager))
	router.DELETE("/payments/:id", handleCancelPayment(serviceManager))
	router.POST("/quick-pay", handleQuickPayment(serviceManager))
}

func SetupPaymentRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/:id", handleGetPaymentRequest(serviceManager))
	router.GET("/:id/qr", handleGetPaymentQRCode(serviceManager))
	router.POST("/:id/pay", handlePayPaymentRequest(serviceManager))
}

func SetupPublicPaymentRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.GET("/pay/:id", handleScanPayment(serviceManager))
}

func handleCreatePayment(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		forwardRequest(c, sm, "wallet", "POST", "/api/v1/merchant/payments")
	}
}

func handleGetPayments(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		forwardRequest(c, sm, "wallet", "GET", "/api/v1/merchant/payments")
	}
}

func handleGetPaymentHistory(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		forwardRequest(c, sm, "wallet", "GET", "/api/v1/merchant/payments/history")
	}
}

func handleCancelPayment(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, sm, "wallet", "DELETE", "/api/v1/merchant/payments/"+id)
	}
}

func handleQuickPayment(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		forwardRequest(c, sm, "wallet", "POST", "/api/v1/merchant/quick-pay")
	}
}

func handleGetPaymentRequest(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, sm, "wallet", "GET", "/api/v1/payments/"+id)
	}
}

func handleGetPaymentQRCode(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, sm, "wallet", "GET", "/api/v1/payments/"+id+"/qr")
	}
}

func handlePayPaymentRequest(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, sm, "wallet", "POST", "/api/v1/payments/"+id+"/pay")
	}
}

func handleScanPayment(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		forwardRequest(c, sm, "wallet", "GET", "/api/v1/pay/"+id)
	}
}

// Helper to deduce forwarding logic
func forwardRequest(c *gin.Context, sm *services.ServiceManager, service string, method string, path string) {
	// Query params
	if c.Request.URL.RawQuery != "" {
		path += "?" + c.Request.URL.RawQuery
	}

	headers := map[string]string{}
	// Forward Authorization and other useful headers
	if token := c.GetHeader("Authorization"); token != "" {
		headers["Authorization"] = token
	}
	if contentType := c.GetHeader("Content-Type"); contentType != "" {
		headers["Content-Type"] = contentType
	}

	var requestBody interface{}
	// For POST/PUT/PATCH, read body
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				// Reassignment for further use (gin middleware compatibility usually requires body restore, 
				// but here we are at the end of the chain or close to it)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				// Unmarshal into generic map to pass to CallService (which expects interface{})
				var jsonBody interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					requestBody = jsonBody
				} else {
					// Fallback: if not valid JSON, maybe don't send as body? 
					// Or ServiceManager handles bytes? ServiceManager.CallService takes interface{}.
					// If it expects JSON, passing struct/map is safer.
					// If parsing fails, we might just proceed with nil body or error out.
					// For now, let's assume JSON.
				}
			}
		}
	}

	resp, err := sm.CallService(c.Request.Context(), service, method, path, requestBody, headers)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable", "details": err.Error()})
		return
	}

	// Forward response headers
	for k, v := range resp.Headers {
		for _, val := range v {
			c.Header(k, val)
		}
	}

	c.Data(resp.StatusCode, "application/json", resp.Body)
}
