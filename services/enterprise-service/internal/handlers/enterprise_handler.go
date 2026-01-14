package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	repo *repository.EnterpriseRepository
}

func NewEnterpriseHandler(repo *repository.EnterpriseRepository) *EnterpriseHandler {
	return &EnterpriseHandler{repo: repo}
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
