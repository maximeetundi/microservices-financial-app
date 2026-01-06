package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContactsHandler struct {
	db *sql.DB
}

func NewContactsHandler(db *sql.DB) *ContactsHandler {
	return &ContactsHandler{db: db}
}

// Contact represents a user's saved contact
type Contact struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ensureContactsTable creates the contacts table if it doesn't exist
func (h *ContactsHandler) ensureContactsTable() error {
	_, err := h.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_contacts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			phone VARCHAR(50),
			email VARCHAR(255),
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id, phone),
			UNIQUE(user_id, email)
		)
	`)
	return err
}

// GetContacts returns all contacts for the current user
func (h *ContactsHandler) GetContacts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Ensure table exists
	h.ensureContactsTable()

	rows, err := h.db.Query(`
		SELECT id, user_id, phone, email, name, created_at, updated_at
		FROM user_contacts
		WHERE user_id = $1
		ORDER BY name ASC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}
	defer rows.Close()

	contacts := []Contact{}
	for rows.Next() {
		var contact Contact
		var email sql.NullString
		if err := rows.Scan(&contact.ID, &contact.UserID, &contact.Phone, &email, &contact.Name, &contact.CreatedAt, &contact.UpdatedAt); err != nil {
			continue
		}
		if email.Valid {
			contact.Email = email.String
		}
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

// AddContact adds a single contact
func (h *ContactsHandler) AddContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if req.Phone == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone or email is required"})
		return
	}

	// Ensure table exists
	h.ensureContactsTable()

	id := uuid.New().String()
	now := time.Now()

	// Try to insert, on conflict update the name
	_, err := h.db.Exec(`
		INSERT INTO user_contacts (id, user_id, phone, email, name, created_at, updated_at)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6, $6)
		ON CONFLICT (user_id, phone) WHERE phone IS NOT NULL
		DO UPDATE SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`, id, userID, req.Phone, req.Email, req.Name, now)

	if err != nil {
		// Try email conflict
		_, err = h.db.Exec(`
			INSERT INTO user_contacts (id, user_id, phone, email, name, created_at, updated_at)
			VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6, $6)
			ON CONFLICT (user_id, email) WHERE email IS NOT NULL
			DO UPDATE SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
		`, id, userID, req.Phone, req.Email, req.Name, now)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add contact"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"contact": Contact{
			ID:        id,
			UserID:    userID.(string),
			Phone:     req.Phone,
			Email:     req.Email,
			Name:      req.Name,
			CreatedAt: now,
			UpdatedAt: now,
		},
	})
}

// SyncContacts bulk syncs contacts from mobile
// Only contacts whose phone/email matches a registered user are saved
func (h *ContactsHandler) SyncContacts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Contacts []struct {
			Phone string `json:"phone"`
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"contacts"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Ensure table exists
	h.ensureContactsTable()

	synced := 0
	skipped := 0
	for _, contact := range req.Contacts {
		if contact.Name == "" || (contact.Phone == "" && contact.Email == "") {
			continue
		}

		// Check if this phone/email belongs to a registered user
		var userExists bool
		if contact.Phone != "" {
			err := h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`, contact.Phone).Scan(&userExists)
			if err != nil || !userExists {
				// Try with normalized phone (remove spaces, dashes)
				err = h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE REPLACE(REPLACE(phone, ' ', ''), '-', '') = REPLACE(REPLACE($1, ' ', ''), '-', ''))`, contact.Phone).Scan(&userExists)
			}
		}
		if !userExists && contact.Email != "" {
			h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER($1))`, contact.Email).Scan(&userExists)
		}

		if !userExists {
			skipped++
			continue
		}

		id := uuid.New().String()
		now := time.Now()

		// Upsert based on phone or email
		_, err := h.db.Exec(`
			INSERT INTO user_contacts (id, user_id, phone, email, name, created_at, updated_at)
			VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6, $6)
			ON CONFLICT (user_id, phone) WHERE phone IS NOT NULL
			DO UPDATE SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
		`, id, userID, contact.Phone, contact.Email, contact.Name, now)

		if err == nil {
			synced++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"synced":  synced,
		"skipped": skipped,
		"total":   len(req.Contacts),
		"message": "Contacts synced successfully (only registered users)",
	})
}

// DeleteContact deletes a contact by ID
func (h *ContactsHandler) DeleteContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contactID := c.Param("id")

	result, err := h.db.Exec(`
		DELETE FROM user_contacts
		WHERE id = $1 AND user_id = $2
	`, contactID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LookupContactName looks up a contact name by phone or email
func (h *ContactsHandler) LookupContactName(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	phone := c.Query("phone")
	email := c.Query("email")

	if phone == "" && email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone or email required"})
		return
	}

	var name string
	var err error

	if phone != "" {
		err = h.db.QueryRow(`
			SELECT name FROM user_contacts
			WHERE user_id = $1 AND phone = $2
		`, userID, phone).Scan(&name)
	}

	if (err != nil || name == "") && email != "" {
		err = h.db.QueryRow(`
			SELECT name FROM user_contacts
			WHERE user_id = $1 AND email = $2
		`, userID, email).Scan(&name)
	}

	if err != nil || name == "" {
		c.JSON(http.StatusOK, gin.H{"found": false, "name": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"found": true, "name": name})
}
