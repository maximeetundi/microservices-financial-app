package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/gin-gonic/gin"
)

type PlatformProxyHandler struct {
	config *config.Config
	client *http.Client
}

func NewPlatformProxyHandler(cfg *config.Config) *PlatformProxyHandler {
	return &PlatformProxyHandler{
		config: cfg,
		client: &http.Client{
			// Timeout to prevent hanging
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest forwards the request to the wallet service
func (h *PlatformProxyHandler) ProxyRequest(c *gin.Context) {
	// target path should appear after /fees/ or just match the suffix
	// For /platform/accounts -> /api/v1/admin/platform/accounts
	// Wallet service expects: /api/v1/admin/accounts (maybe?)
	// Let's check wallet-service router.
	// Assuming wallet-service has /api/v1/admin/platform/... or similar.
	// Users reported: /platform/accounts handled by AdminPlatformHandler

	// We will forward to: {WalletServiceURL}/api/v1/admin/platform/{wildcard}
	// The incoming path in admin-service main.go will be /api/v1/admin/platform/*path

	proxyPath := c.Param("path")
	targetURL := h.config.WalletServiceURL + "/api/v1/admin/platform" + proxyPath

	// Prepare request
	var body io.Reader
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore for other middlewares if any
		body = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Copy headers
	for k, v := range c.Request.Header {
		// Skip hop-by-hop headers
		if k == "Host" || k == "Content-Length" || k == "Connection" {
			continue
		}
		req.Header[k] = v
	}

	// Ensure Auth token is present (it should be in headers already)
	// If not, we might need to add it explicitly if using service-to-service auth

	// Execute
	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to wallet service"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Copy body
	io.Copy(c.Writer, resp.Body)
}
