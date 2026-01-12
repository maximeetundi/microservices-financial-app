package handlers

import (
	"net/http"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/services"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage *services.StorageService
}

func NewUploadHandler(storage *services.StorageService) *UploadHandler {
	return &UploadHandler{storage: storage}
}

// UploadMedia handles image and video uploads to MinIO
func (h *UploadHandler) UploadMedia(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback detection could be added here
		contentType = "application/octet-stream"
	}

	// Check if it's an image or video
	isImage := strings.HasPrefix(contentType, "image/")
	isVideo := strings.HasPrefix(contentType, "video/")

	if !isImage && !isVideo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image and video files are allowed"})
		return
	}

	// Limit file size (10MB for images, 50MB for videos)
	maxSize := int64(10 * 1024 * 1024) // 10MB default
	if isVideo {
		maxSize = 50 * 1024 * 1024 // 50MB for video
	}

	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit (10MB for images, 50MB for videos)"})
		return
	}

	folder := "campaigns"
	if isVideo {
		folder = "campaign-videos"
	}

	// Upload to MinIO
	url, err := h.storage.UploadFile(c.Request.Context(), file, header.Filename, header.Size, contentType, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"type":    contentType,
		"message": "File uploaded successfully",
	})
}
