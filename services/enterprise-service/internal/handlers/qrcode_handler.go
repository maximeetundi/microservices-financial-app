package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	BaseURL string
}

func NewQRCodeHandler(baseURL string) *QRCodeHandler {
	return &QRCodeHandler{BaseURL: baseURL}
}

// GenerateEnterpriseQR returns a QR code for the Enterprise Profile
// GET /enterprises/:id/qrcode
func (h *QRCodeHandler) GenerateEnterpriseQR(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enterprise ID is required"})
		return
	}

	// URL to the Enterprise Public Profile (or deeply linked app url)
	// Example: https://app.maximeetundi.store/enterprises/:id
	url := fmt.Sprintf("%s/enterprises/%s", h.BaseURL, id)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

// GenerateServiceQR returns a QR code for a specific Service subscription
// GET /enterprises/:id/services/:serviceId/qrcode
func (h *QRCodeHandler) GenerateServiceQR(c *gin.Context) {
	entId := c.Param("id")
	svcId := c.Param("serviceId")

	if entId == "" || svcId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enterprise ID and Service ID are required"})
		return
	}

	// URL to the Service Subscription Form
	// Example: https://app.maximeetundi.store/enterprises/:id/subscribe?service_id=:svcId
	url := fmt.Sprintf("%s/enterprises/%s/subscribe?service_id=%s", h.BaseURL, entId, svcId)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}
