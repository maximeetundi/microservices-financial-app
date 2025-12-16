package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/google/uuid"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	conv.ID = uuid.New().String()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
	conv.Status = models.ConversationStatusOpen

	query := `
		INSERT INTO conversations (id, user_id, user_name, user_email, agent_type, subject, category, status, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		conv.ID,
		conv.UserID,
		conv.UserName,
		conv.UserEmail,
		conv.AgentType,
		conv.Subject,
		conv.Category,
		conv.Status,
		conv.Priority,
		conv.CreatedAt,
		conv.UpdatedAt,
	)

	return err
}

func (r *ConversationRepository) GetByID(id string) (*models.Conversation, error) {
	query := `
		SELECT id, user_id, user_name, user_email, agent_id, agent_type, subject, category, 
		       status, priority, last_message, last_message_at, unread_count, message_count,
		       rating, feedback, resolved_at, created_at, updated_at
		FROM conversations WHERE id = $1
	`

	conv := &models.Conversation{}
	err := r.db.QueryRow(query, id).Scan(
		&conv.ID,
		&conv.UserID,
		&conv.UserName,
		&conv.UserEmail,
		&conv.AgentID,
		&conv.AgentType,
		&conv.Subject,
		&conv.Category,
		&conv.Status,
		&conv.Priority,
		&conv.LastMessage,
		&conv.LastMessageAt,
		&conv.UnreadCount,
		&conv.MessageCount,
		&conv.Rating,
		&conv.Feedback,
		&conv.ResolvedAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *ConversationRepository) GetByUserID(userID string, limit, offset int) ([]*models.Conversation, error) {
	query := `
		SELECT id, user_id, user_name, user_email, agent_id, agent_type, subject, category, 
		       status, priority, last_message, last_message_at, unread_count, message_count,
		       rating, feedback, resolved_at, created_at, updated_at
		FROM conversations 
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.UserName,
			&conv.UserEmail,
			&conv.AgentID,
			&conv.AgentType,
			&conv.Subject,
			&conv.Category,
			&conv.Status,
			&conv.Priority,
			&conv.LastMessage,
			&conv.LastMessageAt,
			&conv.UnreadCount,
			&conv.MessageCount,
			&conv.Rating,
			&conv.Feedback,
			&conv.ResolvedAt,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *ConversationRepository) GetAll(status string, limit, offset int) ([]*models.Conversation, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, user_id, user_name, user_email, agent_id, agent_type, subject, category, 
			       status, priority, last_message, last_message_at, unread_count, message_count,
			       rating, feedback, resolved_at, created_at, updated_at
			FROM conversations 
			WHERE status = $1
			ORDER BY 
				CASE priority 
					WHEN 'urgent' THEN 1 
					WHEN 'high' THEN 2 
					WHEN 'medium' THEN 3 
					ELSE 4 
				END,
				updated_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{status, limit, offset}
	} else {
		query = `
			SELECT id, user_id, user_name, user_email, agent_id, agent_type, subject, category, 
			       status, priority, last_message, last_message_at, unread_count, message_count,
			       rating, feedback, resolved_at, created_at, updated_at
			FROM conversations 
			ORDER BY 
				CASE priority 
					WHEN 'urgent' THEN 1 
					WHEN 'high' THEN 2 
					WHEN 'medium' THEN 3 
					ELSE 4 
				END,
				updated_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.UserID,
			&conv.UserName,
			&conv.UserEmail,
			&conv.AgentID,
			&conv.AgentType,
			&conv.Subject,
			&conv.Category,
			&conv.Status,
			&conv.Priority,
			&conv.LastMessage,
			&conv.LastMessageAt,
			&conv.UnreadCount,
			&conv.MessageCount,
			&conv.Rating,
			&conv.Feedback,
			&conv.ResolvedAt,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *ConversationRepository) UpdateStatus(id string, status models.ConversationStatus) error {
	query := `UPDATE conversations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *ConversationRepository) AssignAgent(id, agentID string) error {
	query := `UPDATE conversations SET agent_id = $1, status = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, agentID, models.ConversationStatusActive, time.Now(), id)
	return err
}

func (r *ConversationRepository) UpdateLastMessage(id, message string) error {
	now := time.Now()
	query := `UPDATE conversations SET last_message = $1, last_message_at = $2, message_count = message_count + 1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, message, now, id)
	return err
}

func (r *ConversationRepository) Close(id string, rating int, feedback string) error {
	now := time.Now()
	query := `UPDATE conversations SET status = $1, rating = $2, feedback = $3, resolved_at = $4, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(query, models.ConversationStatusClosed, rating, feedback, now, id)
	return err
}

func (r *ConversationRepository) GetStats() (*models.SupportStats, error) {
	stats := &models.SupportStats{}

	// Total conversations
	r.db.QueryRow(`SELECT COUNT(*) FROM conversations`).Scan(&stats.TotalConversations)

	// Open conversations
	r.db.QueryRow(`SELECT COUNT(*) FROM conversations WHERE status IN ('open', 'pending', 'active')`).Scan(&stats.OpenConversations)

	// Resolved today
	r.db.QueryRow(`SELECT COUNT(*) FROM conversations WHERE status = 'resolved' AND resolved_at >= CURRENT_DATE`).Scan(&stats.ResolvedToday)

	// Pending conversations
	r.db.QueryRow(`SELECT COUNT(*) FROM conversations WHERE status = 'pending'`).Scan(&stats.PendingConversations)

	// Average satisfaction (from ratings)
	r.db.QueryRow(`SELECT COALESCE(AVG(rating), 0) FROM conversations WHERE rating IS NOT NULL`).Scan(&stats.CustomerSatisfaction)

	return stats, nil
}
