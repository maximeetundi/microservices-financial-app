package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

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

// UpdatePresence updates the user's last_seen timestamp (call this when user is active)
func (h *UserHandler) UpdatePresence(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE users SET last_seen = $1 WHERE id = $2
	`, time.Now(), userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update presence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetUserPresence gets the last_seen status of a user
func (h *UserHandler) GetUserPresence(c *gin.Context) {
	userID := c.Param("id")
	
	var lastSeen sql.NullTime
	err := h.db.QueryRow(`
		SELECT last_seen FROM users WHERE id = $1
	`, userID).Scan(&lastSeen)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Calculate status
	status := "offline"
	lastSeenStr := ""
	
	if lastSeen.Valid {
		lastSeenStr = lastSeen.Time.Format(time.RFC3339)
		diff := time.Since(lastSeen.Time)
		
		if diff.Minutes() < 5 {
			status = "online"
		} else if diff.Minutes() < 60 {
			status = "away"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"status":    status,
		"last_seen": lastSeenStr,
	})
}

// GetMultiplePresence gets presence for multiple users at once
func (h *UserHandler) GetMultiplePresence(c *gin.Context) {
	var req struct {
		UserIDs []string `json:"user_ids"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"presences": []interface{}{}})
		return
	}

	// Build query for multiple users
	presences := make([]map[string]interface{}, 0)
	
	for _, uid := range req.UserIDs {
		var lastSeen sql.NullTime
		err := h.db.QueryRow(`SELECT last_seen FROM users WHERE id = $1`, uid).Scan(&lastSeen)
		if err != nil {
			continue
		}

		status := "offline"
		lastSeenStr := ""
		
		if lastSeen.Valid {
			lastSeenStr = lastSeen.Time.Format(time.RFC3339)
			diff := time.Since(lastSeen.Time)
			
			if diff.Minutes() < 5 {
				status = "online"
			} else if diff.Minutes() < 60 {
				status = "away"
			}
		}

		presences = append(presences, map[string]interface{}{
			"user_id":   uid,
			"status":    status,
			"last_seen": lastSeenStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{"presences": presences})
}
