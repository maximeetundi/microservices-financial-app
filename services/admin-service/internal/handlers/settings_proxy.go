package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SettingsProxyHandler struct {
	config *config.Config
	client *http.Client
}

func NewSettingsProxyHandler(cfg *config.Config) *SettingsProxyHandler {
	return &SettingsProxyHandler{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest forwards the request to the exchange service
func (h *SettingsProxyHandler) ProxyRequest(c *gin.Context) {
	// Target: {ExchangeServiceURL}/api/v1/admin/settings
	targetURL := h.config.ExchangeServiceURL + "/api/v1/admin/settings"

	key := c.Param("key")
	if key != "" {
		targetURL += "/" + key
	}

	// Prepare request body
	var body io.Reader
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		body = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Generate internal token for Exchange Service
	// Exchange Service expects a user token with role=admin signed with JWT_SECRET
	token, err := h.generateProxyToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate proxy auth token"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Copy other headers
	for k, v := range c.Request.Header {
		if k == "Host" || k == "Content-Length" || k == "Connection" || k == "Authorization" {
			continue
		}
		req.Header[k] = v
	}

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to exchange service"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (h *SettingsProxyHandler) generateProxyToken(c *gin.Context) (string, error) {
	adminID := c.GetString("admin_id")
	email := c.GetString("email") // Check if email is stored as 'email' or 'admin_email'
	if email == "" {
		email = c.GetString("admin_email")
	}

	// Claims expected by Exchange Service middleware
	claims := jwt.MapClaims{
		"user_id":   adminID, // Map admin_id to user_id
		"email":     email,
		"role":      "admin",                                // Force admin role
		"kyc_level": 3,                                      // Assume full access
		"exp":       time.Now().Add(1 * time.Minute).Unix(), // Short lived
		"iat":       time.Now().Unix(),
		"type":      "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.config.JWTSecret))
}
