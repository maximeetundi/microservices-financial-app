package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	BaseURL string
	EntRepo *repository.EnterpriseRepository
}

func NewQRCodeHandler(baseURL string, entRepo *repository.EnterpriseRepository) *QRCodeHandler {
	return &QRCodeHandler{
		BaseURL: baseURL,
		EntRepo: entRepo,
	}
}

// GenerateEnterpriseQR returns a QR code for the Enterprise Profile
// GET /enterprises/:id/qrcode
func (h *QRCodeHandler) GenerateEnterpriseQR(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enterprise ID is required"})
		return
	}

	// URL to the Enterprise Public Profile
	url := fmt.Sprintf("%s/enterprises/%s/subscribe", h.BaseURL, id)

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

	// Verify Service Privacy
	ent, err := h.EntRepo.FindByID(context.Background(), entId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}

	isPrivate := false
	found := false

	for _, group := range ent.ServiceGroups {
		for _, svc := range group.Services {
			if svc.ID == svcId {
				found = true
				if group.IsPrivate {
					isPrivate = true
				}
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if isPrivate {
		// Strict Requirement: Private services cannot have public QR codes generated via API
		c.JSON(http.StatusForbidden, gin.H{"error": "QR Code generation is disabled for private services"})
		return
	}

	// URL to the Service Subscription Form
	url := fmt.Sprintf("%s/enterprises/%s/subscribe?service_id=%s", h.BaseURL, entId, svcId)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}
