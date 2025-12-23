package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/google/uuid"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Create(agent *models.Agent) error {
	agent.ID = uuid.New().String()
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()

	query := `
		INSERT INTO agents (id, name, email, type, avatar, is_available, max_chats, active_chats, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		agent.ID,
		agent.Name,
		agent.Email,
		agent.Type,
		agent.Avatar,
		agent.IsAvailable,
		agent.MaxChats,
		agent.ActiveChats,
		agent.Rating,
		agent.CreatedAt,
		agent.UpdatedAt,
	)

	return err
}

func (r *AgentRepository) GetByID(id string) (*models.Agent, error) {
	query := `
		SELECT id, name, email, type, avatar, is_available, max_chats, active_chats, rating, created_at, updated_at
		FROM agents WHERE id = $1
	`

	agent := &models.Agent{}
	err := r.db.QueryRow(query, id).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Email,
		&agent.Type,
		&agent.Avatar,
		&agent.IsAvailable,
		&agent.MaxChats,
		&agent.ActiveChats,
		&agent.Rating,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (r *AgentRepository) GetAvailable(agentType models.AgentType) ([]*models.Agent, error) {
	query := `
		SELECT id, name, email, type, avatar, is_available, max_chats, active_chats, rating, created_at, updated_at
		FROM agents 
		WHERE type = $1 AND is_available = true AND active_chats < max_chats
		ORDER BY active_chats ASC, rating DESC
	`

	rows, err := r.db.Query(query, agentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		agent := &models.Agent{}
		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Email,
			&agent.Type,
			&agent.Avatar,
			&agent.IsAvailable,
			&agent.MaxChats,
			&agent.ActiveChats,
			&agent.Rating,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

func (r *AgentRepository) GetAll() ([]*models.Agent, error) {
	query := `
		SELECT id, name, email, type, avatar, is_available, max_chats, active_chats, rating, created_at, updated_at
		FROM agents 
		ORDER BY type, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		agent := &models.Agent{}
		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Email,
			&agent.Type,
			&agent.Avatar,
			&agent.IsAvailable,
			&agent.MaxChats,
			&agent.ActiveChats,
			&agent.Rating,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

func (r *AgentRepository) UpdateAvailability(id string, isAvailable bool) error {
	query := `UPDATE agents SET is_available = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, isAvailable, time.Now(), id)
	return err
}

func (r *AgentRepository) IncrementActiveChats(id string) error {
	query := `UPDATE agents SET active_chats = active_chats + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *AgentRepository) DecrementActiveChats(id string) error {
	query := `UPDATE agents SET active_chats = GREATEST(active_chats - 1, 0), updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *AgentRepository) GetAIAgent() (*models.Agent, error) {
	query := `SELECT id, name, email, type, avatar, is_available, max_chats, active_chats, rating, created_at, updated_at FROM agents WHERE type = 'ai' LIMIT 1`
	
	agent := &models.Agent{}
	err := r.db.QueryRow(query).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Email,
		&agent.Type,
		&agent.Avatar,
		&agent.IsAvailable,
		&agent.MaxChats,
		&agent.ActiveChats,
		&agent.Rating,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create default AI agent if not exists
		agent = &models.Agent{
			Name:        "Assistant IA",
			Email:       "ai@zekora.com",
			Type:        models.AgentTypeAI,
			Avatar:      "🤖",
			IsAvailable: true,
			MaxChats:    1000,
			ActiveChats: 0,
			Rating:      5.0,
		}
		r.Create(agent)
		return agent, nil
	}

	return agent, err
}
