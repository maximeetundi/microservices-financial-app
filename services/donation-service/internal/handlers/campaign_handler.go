package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/services"
	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service *services.CampaignService
}

func NewCampaignHandler(service *services.CampaignService) *CampaignHandler {
	return &CampaignHandler{service: service}
}

func (h *CampaignHandler) Create(c *gin.Context) {
	var campaign models.Campaign
	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming Auth middleware sets UserID
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	campaign.CreatorID = userID

	if err := h.service.CreateCampaign(&campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "campaign created", "campaign": campaign})
}

func (h *CampaignHandler) Get(c *gin.Context) {
	id := c.Param("id")
	campaign, err := h.service.GetCampaign(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func (h *CampaignHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	// Filter by creator if requested
	creatorID := c.Query("creator_id")
	if creatorID != "" {
		campaigns, err := h.service.GetMyCampaigns(creatorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
		return
	}

	campaigns, err := h.service.ListCampaigns(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

func (h *CampaignHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetHeader("X-User-ID")
	
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize updates? Service checks ownership.
	
	if err := h.service.UpdateCampaign(id, userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "campaign updated"})
}

func (h *CampaignHandler) GetStats(c *gin.Context) {
	// Not implemented yet separately, available in GetCampaign
	c.Status(http.StatusNotImplemented)
}

// UploadImage Handler stub - implementation depends on MinIO setup via shared utility or direct
func (h *CampaignHandler) UploadImage(c *gin.Context) {
	// For MVP, just return dummy or implement if util available
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
