package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/gin-gonic/gin"
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
	// Incoming: /api/v1/admin/settings[/...]

	// We might have path parameters handled by Gin if we register specific routes,
	// or we can use wildcard.
	// Since main.go will likely register GET/PUT /settings separately, we should make this generic or handle path construction.

	targetURL := h.config.ExchangeServiceURL + "/api/v1/admin/settings"

	// If path param key is present in context (e.g. /settings/:key)
	key := c.Param("key")
	if key != "" {
		targetURL += "/" + key
	}

	// Prepare request
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

	// Copy headers
	for k, v := range c.Request.Header {
		if k == "Host" || k == "Content-Length" || k == "Connection" {
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
