package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// SearchUsers searches for users by email or phone
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	
	if len(query) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be at least 3 characters"})
		return
	}

	query = strings.ToLower(strings.TrimSpace(query))

	var users []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	// Search by email or phone
	err := h.db.Table("users").
		Select("id, first_name || ' ' || last_name as name, email, phone").
		Where("LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(10).
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
