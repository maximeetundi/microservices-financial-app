package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadFile handles file upload to MinIO for chat attachments
func (h *MessageHandler) UploadFile(c *gin.Context) {
	// Get file from multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Get file type from form
	fileType := c.PostForm("type") // image, audio, document, video
	if fileType == "" {
		fileType = "document"
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	objectName := fmt.Sprintf("%s/%s-%s%s", fileType, time.Now().Format("2006-01-02"), uuid.New().String(), ext)

	// Upload to MinIO
	ctx := context.Background()
	_, err = h.storage.UploadFileStream(ctx, objectName, file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Generate presigned URL (valid for 7 days)
	url, err := h.storage.GetPresignedURL(ctx, objectName, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":       url,
		"filename":  header.Filename,
		"size":      header.Size,
		"mime_type": header.Header.Get("Content-Type"),
	})
}
