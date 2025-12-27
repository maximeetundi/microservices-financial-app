package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
)

type PreferencesHandler struct {
	repo *repository.PreferencesRepository
}

func NewPreferencesHandler(repo *repository.PreferencesRepository) *PreferencesHandler {
	return &PreferencesHandler{repo: repo}
}

// ============ User Preferences ============

// GetPreferences returns user preferences
func (h *PreferencesHandler) GetPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	prefs, err := h.repo.GetUserPreferences(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences updates user preferences
func (h *PreferencesHandler) UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefs, err := h.repo.UpdateUserPreferences(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// ============ Notification Preferences ============

// GetNotificationPrefs returns notification preferences
func (h *PreferencesHandler) GetNotificationPrefs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	prefs, err := h.repo.GetNotificationPrefs(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdateNotificationPrefs updates notification preferences
func (h *PreferencesHandler) UpdateNotificationPrefs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.UpdateNotificationPrefsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefs, err := h.repo.UpdateNotificationPrefs(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// ============ KYC ============

// GetKYCStatus returns the KYC status and documents
func (h *PreferencesHandler) GetKYCStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docs, err := h.repo.GetKYCDocuments(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get KYC documents"})
		return
	}

	// Determine overall status based on documents
	status := "pending"
	allApproved := true
	hasRejected := false
	hasSubmitted := false

	for _, doc := range docs {
		if doc.Status == "rejected" {
			hasRejected = true
			allApproved = false
		} else if doc.Status == "pending" {
			hasSubmitted = true
			allApproved = false
		} else if doc.Status != "approved" {
			allApproved = false
		}
	}

	if len(docs) > 0 {
		if allApproved && len(docs) >= 3 {
			status = "verified"
		} else if hasRejected {
			status = "rejected"
		} else if hasSubmitted {
			status = "submitted"
		}
	}

	level := 0
	if status == "verified" {
		level = 2
	} else if len(docs) > 0 {
		level = 1
	}

	c.JSON(http.StatusOK, models.KYCStatusResponse{
		Status:    status,
		Level:     level,
		Documents: docs,
	})
}

// GetKYCDocuments returns all KYC documents for the user
func (h *PreferencesHandler) GetKYCDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docs, err := h.repo.GetKYCDocuments(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

// UploadKYCDocument handles KYC document upload
func (h *PreferencesHandler) UploadKYCDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docType := c.PostForm("type")
	if docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document type is required"})
		return
	}

	if docType != "identity" && docType != "selfie" && docType != "address" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document file is required"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg":               true,
		"image/png":                true,
		"image/gif":                true,
		"image/heic":               true,
		"image/heif":               true,
		"application/pdf":          true,
		"application/octet-stream": true, // Fallback for unknown types
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: JPEG, PNG, GIF, HEIC, PDF"})
		return
	}

	// Create upload directory
	uploadDir := fmt.Sprintf("uploads/kyc/%s", userID.(string))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s_%s_%s%s", docType, time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Save file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create document record
	doc := &models.KYCDocument{
		UserID:   userID.(string),
		Type:     docType,
		FileName: header.Filename,
		FilePath: filePath,
		FileSize: header.Size,
		MimeType: contentType,
	}

	if err := h.repo.CreateKYCDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": doc,
	})
}
