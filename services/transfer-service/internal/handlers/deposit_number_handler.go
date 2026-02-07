package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type DepositNumberHandler struct {
	repo *repository.DepositNumberRepository
}

func NewDepositNumberHandler(repo *repository.DepositNumberRepository) *DepositNumberHandler {
	return &DepositNumberHandler{repo: repo}
}

type CreateDepositNumberRequest struct {
	Country string `json:"country" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Label   string `json:"label"`
	IsDefault bool `json:"is_default"`
}

func (h *DepositNumberHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)
	if strings.TrimSpace(userIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	numbers, err := h.repo.ListByUser(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load deposit numbers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"numbers": numbers, "max": repository.MaxDepositNumbersPerUser})
}

func (h *DepositNumberHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)
	if strings.TrimSpace(userIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateDepositNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Country = strings.ToUpper(strings.TrimSpace(req.Country))
	req.Phone = strings.TrimSpace(req.Phone)
	req.Label = strings.TrimSpace(req.Label)

	// Deduplicate: if exists, return it
	existing, err := h.repo.GetByUserCountryPhone(c.Request.Context(), userIDStr, req.Country, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing number"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusOK, gin.H{"number": existing})
		return
	}

	n, err := h.repo.Create(c.Request.Context(), userIDStr, req.Country, req.Phone, req.Label, req.IsDefault)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"number": n})
}

func (h *DepositNumberHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr := fmt.Sprintf("%v", userID)
	if strings.TrimSpace(userIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	id = strings.TrimSpace(id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
		return
	}

	err := h.repo.Delete(c.Request.Context(), userIDStr, id)
	if err != nil {
		if err.Error() == "cannot delete last deposit number" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vous devez garder au moins un numéro"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
