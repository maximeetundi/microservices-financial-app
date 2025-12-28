package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/services"
	"github.com/gin-gonic/gin"
)

// KYCHandler handles KYC document operations
type KYCHandler struct {
	storage *services.StorageService
}

// NewKYCHandler creates a new KYC handler
func NewKYCHandler(storage *services.StorageService) *KYCHandler {
	return &KYCHandler{storage: storage}
}

// GetDocumentURL generates a presigned URL for secure document access
// POST /api/v1/admin/kyc/document-url
func (h *KYCHandler) GetDocumentURL(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_path is required"})
		return
	}

	if h.storage == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Generate presigned URL valid for 15 minutes
	url, err := h.storage.GetPresignedURL(context.Background(), req.FilePath, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secure URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":        url,
		"expires_in": "15m",
	})
}

// GetDocumentDownloadURL generates a presigned URL for downloading with filename
// POST /api/v1/admin/kyc/download-url
func (h *KYCHandler) GetDocumentDownloadURL(c *gin.Context) {
	var req struct {
		FilePath     string `json:"file_path" binding:"required"`
		DownloadName string `json:"download_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_path is required"})
		return
	}

	if h.storage == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage service not available"})
		return
	}

	// Generate presigned URL valid for 5 minutes (download)
	url, err := h.storage.GetPresignedURLWithDownload(context.Background(), req.FilePath, 5*time.Minute, req.DownloadName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":        url,
		"expires_in": "5m",
	})
}
