package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageService *services.StorageService
}

func NewUploadHandler(storageService *services.StorageService) *UploadHandler {
	return &UploadHandler{storageService: storageService}
}

// UploadMedia uploads a file to storage
func (h *UploadHandler) UploadMedia(c *gin.Context) {
	if h.storageService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage service not available"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Only allow images
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed."})
		return
	}

	// Read file
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Upload
	filename := filepath.Base(header.Filename)
	url, err := h.storageService.UploadFile(c.Request.Context(), data, filename, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
