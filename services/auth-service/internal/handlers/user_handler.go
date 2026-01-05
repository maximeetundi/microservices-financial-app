package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
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

	// Search by email or phone using raw SQL
	rows, err := h.db.Query(`
		SELECT id, first_name || ' ' || last_name as name, email, phone
		FROM users
		WHERE LOWER(email) LIKE $1 OR LOWER(phone) LIKE $1
		LIMIT 10
	`, "%"+query+"%")
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id,  name, email, phone string
		if err := rows.Scan(&id, &name, &email, &phone); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
			"phone": phone,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
