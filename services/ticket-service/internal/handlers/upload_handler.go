package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/services"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage *services.StorageService
}

func NewUploadHandler(storage *services.StorageService) *UploadHandler {
	return &UploadHandler{storage: storage}
}

// UploadImage handles image upload to MinIO
func (h *UploadHandler) UploadImage(c *gin.Context) {
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
		contentType = "image/jpeg"
	}

	// Check if it's an image
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed (jpeg, png, gif, webp)"})
		return
	}

	// Limit file size to 10MB
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Upload to MinIO
	url, err := h.storage.UploadFile(c.Request.Context(), file, header.Filename, header.Size, contentType, "events")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"message": "File uploaded successfully",
	})
}
