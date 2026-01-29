package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

type PreferencesHandler struct {
	repo      *repository.PreferencesRepository
	storage   *services.StorageService
	publisher *services.EventPublisher
}

func NewPreferencesHandler(repo *repository.PreferencesRepository, storage *services.StorageService, publisher *services.EventPublisher) *PreferencesHandler {
	return &PreferencesHandler{repo: repo, storage: storage, publisher: publisher}
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

// UploadKYCDocument handles KYC document upload to Minio
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

	// Upload to Minio
	folder := fmt.Sprintf("kyc/%s/%s", userID.(string), docType)
	fileURL, err := h.storage.UploadFile(context.Background(), file, header.Filename, header.Size, contentType, folder)
	if err != nil {
		fmt.Printf("Error uploading to Minio: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Get optional document metadata (only for identity documents)
	var documentNumber, expiryDate, identitySubType *string
	if docType == "identity" {
		if num := c.PostForm("document_number"); num != "" {
			documentNumber = &num
		}
		if date := c.PostForm("expiry_date"); date != "" {
			expiryDate = &date
		}
		if subType := c.PostForm("identity_sub_type"); subType != "" {
			identitySubType = &subType
		}
	}

	// Create document record
	doc := &models.KYCDocument{
		UserID:          userID.(string),
		Type:            docType,
		IdentitySubType: identitySubType,
		FileName:        header.Filename,
		FilePath:        fileURL, // Store the full Minio URL as file path
		FileSize:        header.Size,
		MimeType:        contentType,
		DocumentNumber:  documentNumber,
		ExpiryDate:      expiryDate,
	}

	if err := h.repo.CreateKYCDocument(doc, fileURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	// Update user's KYC status to 'pending' so they appear in admin panel
	if err := h.repo.UpdateUserKYCStatus(userID.(string), "pending"); err != nil {
		// Log but don't fail - document was saved successfully
		fmt.Printf("Warning: Failed to update KYC status for user %s: %v\n", userID.(string), err)
	}

	// Publish KYC submitted event for admin notification
	if h.publisher != nil {
		// Get user info for the notification
		userEmail, _ := c.Get("user_email")
		userName, _ := c.Get("user_name")
		emailStr, _ := userEmail.(string)
		nameStr, _ := userName.(string)
		if nameStr == "" {
			nameStr = "Un utilisateur"
		}
		h.publisher.PublishKYCSubmitted(userID.(string), emailStr, nameStr, docType)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": doc,
		"file_url": fileURL,
	})
}
