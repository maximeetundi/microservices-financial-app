package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminNotification represents an admin notification
type AdminNotification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      string    `json:"data,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationHandler handles admin notification endpoints
type NotificationHandler struct {
	db *sql.DB
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

// GetNotifications returns all admin notifications
// GET /api/v1/admin/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rows, err := h.db.Query(`
		SELECT id, type, title, message, COALESCE(data::text, '{}'), is_read, created_at
		FROM admin_notifications
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	defer rows.Close()

	var notifications []AdminNotification
	for rows.Next() {
		var n AdminNotification
		if err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.Message, &n.Data, &n.IsRead, &n.CreatedAt); err != nil {
			continue
		}
		notifications = append(notifications, n)
	}

	if notifications == nil {
		notifications = []AdminNotification{}
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// GetUnreadCount returns the count of unread notifications
// GET /api/v1/admin/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	var count int
	err := h.db.QueryRow(`SELECT COUNT(*) FROM admin_notifications WHERE is_read = FALSE`).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// MarkAsRead marks a notification as read
// POST /api/v1/admin/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec(`UPDATE admin_notifications SET is_read = TRUE WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllAsRead marks all notifications as read
// POST /api/v1/admin/notifications/mark-all-read
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	_, err := h.db.Exec(`UPDATE admin_notifications SET is_read = TRUE WHERE is_read = FALSE`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CreateNotification creates a new admin notification (for internal use/testing)
// POST /api/v1/admin/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req struct {
		Type    string `json:"type" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
		Data    string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	_, err := h.db.Exec(`
		INSERT INTO admin_notifications (id, type, title, message, data)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`, id, req.Type, req.Title, req.Message, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "success": true})
}

// DeleteOldNotifications deletes notifications older than 30 days
// DELETE /api/v1/admin/notifications/cleanup
func (h *NotificationHandler) DeleteOldNotifications(c *gin.Context) {
	result, err := h.db.Exec(`DELETE FROM admin_notifications WHERE created_at < NOW() - INTERVAL '30 days'`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old notifications"})
		return
	}

	count, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"deleted": count})
}
