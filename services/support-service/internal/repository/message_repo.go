package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *models.Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO messages (id, conversation_id, sender_id, sender_name, sender_type, content, content_type, attachments, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		msg.ID,
		msg.ConversationID,
		msg.SenderID,
		msg.SenderName,
		msg.SenderType,
		msg.Content,
		msg.ContentType,
		pq.Array(msg.Attachments),
		msg.IsRead,
		msg.CreatedAt,
	)

	return err
}

func (r *MessageRepository) GetByConversationID(conversationID string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, sender_name, sender_type, content, content_type, attachments, is_read, read_at, created_at
		FROM messages 
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.SenderName,
			&msg.SenderType,
			&msg.Content,
			&msg.ContentType,
			pq.Array(&msg.Attachments),
			&msg.IsRead,
			&msg.ReadAt,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Ensure attachments is never nil (for proper JSON serialization)
		if msg.Attachments == nil {
			msg.Attachments = []string{}
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepository) MarkAsRead(conversationID, userID string) error {
	now := time.Now()
	query := `UPDATE messages SET is_read = true, read_at = $1 WHERE conversation_id = $2 AND sender_id != $3 AND is_read = false`
	_, err := r.db.Exec(query, now, conversationID, userID)
	return err
}

func (r *MessageRepository) GetUnreadCount(conversationID, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE conversation_id = $1 AND sender_id != $2 AND is_read = false`
	err := r.db.QueryRow(query, conversationID, userID).Scan(&count)
	return count, err
}
