package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	repo    *repository.EnterpriseRepository
	storage *services.StorageService
}

func NewEnterpriseHandler(repo *repository.EnterpriseRepository, storage *services.StorageService) *EnterpriseHandler {
	return &EnterpriseHandler{repo: repo, storage: storage}
}

// UploadLogo handles logo upload
func (h *EnterpriseHandler) UploadLogo(c *gin.Context) {
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

	// Upload to MinIO (folder "enterprises")
	url, err := h.storage.UploadFile(c.Request.Context(), file, header.Filename, header.Size, contentType, "enterprises")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// CreateEnterprise (Point 1)
func (h *EnterpriseHandler) CreateEnterprise(c *gin.Context) {
	var req models.Enterprise
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set OwnerID from JWT (assuming middleware sets "user_id")
	userID, _ := c.Get("user_id")
	req.OwnerID, _ = userID.(string)

	// Set Default Type if missing
	if req.Type == "" {
		req.Type = models.EnterpriseTypeService
	}

	if err := h.repo.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enterprise", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *EnterpriseHandler) GetEnterprise(c *gin.Context) {
	id := c.Param("id")
	ent, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}
	c.JSON(http.StatusOK, ent)
}

func (h *EnterpriseHandler) ListEnterprises(c *gin.Context) {
	enterprises, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list enterprises", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enterprises)
}

// UpdateEnterprise (Required for Settings)
func (h *EnterpriseHandler) UpdateEnterprise(c *gin.Context) {
	id := c.Param("id")
	
	// Fetch existing to ensure ownership & existence
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}

	// Verify Owner (simple check)
	userID, exists := c.Get("user_id")
	if !exists || existing.OwnerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this enterprise"})
		return
	}

	// Bind updates
	var req models.Enterprise
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve ID and Owner
	req.ID = existing.ID
	req.OwnerID = existing.OwnerID
	req.CreatedAt = existing.CreatedAt // Don't overwrite creation time

	if err := h.repo.Update(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update enterprise", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}
