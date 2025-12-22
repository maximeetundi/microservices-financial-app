package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/models"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()
	n.IsRead = false

	query := `INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, n.ID, n.UserID, n.Type, n.Title, n.Message, n.Data, n.IsRead, n.CreatedAt)
	return err
}

func (r *NotificationRepository) GetByUserID(userID string, limit, offset int) ([]models.Notification, error) {
	query := `SELECT id, user_id, type, title, message, data, is_read, read_at, created_at
			  FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Data, &n.IsRead, &n.ReadAt, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepository) GetUnreadCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) MarkAsRead(id, userID string) error {
	query := `UPDATE notifications SET is_read = true, read_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := r.db.Exec(query, time.Now(), id, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	query := `UPDATE notifications SET is_read = true, read_at = $1 WHERE user_id = $2 AND is_read = false`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *NotificationRepository) Delete(id, userID string) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, id, userID)
	return err
}
