package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminSettingsHandler struct {
	settingsService *services.SettingsService
}

func NewAdminSettingsHandler(settingsService *services.SettingsService) *AdminSettingsHandler {
	return &AdminSettingsHandler{settingsService: settingsService}
}

// GetSettings returns all system settings
func (h *AdminSettingsHandler) GetSettings(c *gin.Context) {
	category := c.Query("category")

	var settings []models.SystemSetting
	var err error

	if category != "" {
		settings, err = h.settingsService.GetByCategory(category)
	} else {
		settings, err = h.settingsService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSetting updates a single system setting
func (h *AdminSettingsHandler) UpdateSetting(c *gin.Context) {
	var req models.SystemSetting
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Setting key is required"})
		return
	}

	if req.Value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Setting value is required"})
		return
	}

	err := h.settingsService.UpdateSetting(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully", "setting": req})
}

// GetSettingByKey returns a single setting by key
func (h *AdminSettingsHandler) GetSettingByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Setting key is required"})
		return
	}

	value := h.settingsService.GetSetting(key)
	if value == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}
