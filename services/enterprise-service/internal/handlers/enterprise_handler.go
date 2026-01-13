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
