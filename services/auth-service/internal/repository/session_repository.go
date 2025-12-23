package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewSessionRepository(db *sql.DB, redisClient *redis.Client) *SessionRepository {
	return &SessionRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *SessionRepository) Create(session *models.Session) error {
	// Store in PostgreSQL
	query := `
		INSERT INTO user_sessions (user_id, session_token, refresh_token, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, session.UserID, session.SessionToken, session.RefreshToken,
		session.IPAddress, session.UserAgent, session.ExpiresAt).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Store session in Redis for quick access
	ctx := context.Background()
	sessionData, _ := json.Marshal(session)
	
	// Store session token mapping
	r.redis.Set(ctx, "session:"+session.SessionToken, sessionData, time.Until(session.ExpiresAt))
	
	// Store refresh token mapping
	r.redis.Set(ctx, "refresh:"+session.RefreshToken, session.UserID, time.Until(session.ExpiresAt))

	return nil
}

func (r *SessionRepository) GetByToken(sessionToken string) (*models.Session, error) {
	ctx := context.Background()
	
	// Try Redis first
	sessionData, err := r.redis.Get(ctx, "session:"+sessionToken).Result()
	if err == nil {
		var session models.Session
		if json.Unmarshal([]byte(sessionData), &session) == nil {
			return &session, nil
		}
	}

	// Fallback to PostgreSQL
	query := `
		SELECT id, user_id, session_token, refresh_token, ip_address, user_agent,
			   expires_at, is_active, created_at
		FROM user_sessions
		WHERE session_token = $1 AND is_active = true AND expires_at > NOW()
	`

	var session models.Session
	err = r.db.QueryRow(query, sessionToken).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.RefreshToken,
		&session.IPAddress, &session.UserAgent, &session.ExpiresAt, &session.IsActive,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Update Redis cache
	sessionBytes, _ := json.Marshal(session)
	r.redis.Set(ctx, "session:"+sessionToken, sessionBytes, time.Until(session.ExpiresAt))

	return &session, nil
}

func (r *SessionRepository) GetByRefreshToken(refreshToken string) (*models.Session, error) {
	ctx := context.Background()
	
	log.Printf("GetByRefreshToken: Looking for token (first 20 chars): %s...", refreshToken[:min(20, len(refreshToken))])
	
	// Try Redis first
	userID, err := r.redis.Get(ctx, "refresh:"+refreshToken).Result()
	if err == nil {
		log.Printf("GetByRefreshToken: Found in Redis for user ID=%s", userID)
		// Get full session from database
		return r.getByRefreshTokenFromDB(refreshToken, userID)
	}
	log.Printf("GetByRefreshToken: Not found in Redis (error: %v), falling back to DB", err)

	// Fallback to PostgreSQL only
	return r.getByRefreshTokenFromDB(refreshToken, "")
}

func (r *SessionRepository) getByRefreshTokenFromDB(refreshToken, userID string) (*models.Session, error) {
	log.Printf("getByRefreshTokenFromDB: Querying DB for refresh token (first 20 chars): %s...", refreshToken[:min(20, len(refreshToken))])
	
	query := `
		SELECT id, user_id, session_token, refresh_token, ip_address, user_agent,
			   expires_at, is_active, created_at
		FROM user_sessions
		WHERE refresh_token = $1 AND is_active = true AND expires_at > NOW()
	`

	var session models.Session
	err := r.db.QueryRow(query, refreshToken).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.RefreshToken,
		&session.IPAddress, &session.UserAgent, &session.ExpiresAt, &session.IsActive,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Check if it exists but is expired or inactive
			checkQuery := `
				SELECT id, is_active, expires_at, expires_at < NOW() as is_expired
				FROM user_sessions 
				WHERE refresh_token = $1
			`
			var id string
			var isActive bool
			var expiresAt time.Time
			var isExpired bool
			checkErr := r.db.QueryRow(checkQuery, refreshToken).Scan(&id, &isActive, &expiresAt, &isExpired)
			if checkErr == nil {
				log.Printf("getByRefreshTokenFromDB: Session found but not valid - ID=%s, is_active=%v, expires_at=%v, is_expired=%v", 
					id, isActive, expiresAt, isExpired)
			} else {
				log.Printf("getByRefreshTokenFromDB: No session found with this refresh token at all")
			}
			return nil, fmt.Errorf("refresh token not found")
		}
		log.Printf("getByRefreshTokenFromDB: DB error: %v", err)
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	log.Printf("getByRefreshTokenFromDB: Found valid session ID=%s for user ID=%s", session.ID, session.UserID)
	return &session, nil
}

func (r *SessionRepository) Update(session *models.Session) error {
	query := `
		UPDATE user_sessions 
		SET session_token = $1, refresh_token = $2, expires_at = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(query, session.SessionToken, session.RefreshToken, session.ExpiresAt, session.ID)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	// Update Redis cache
	ctx := context.Background()
	sessionData, _ := json.Marshal(session)
	r.redis.Set(ctx, "session:"+session.SessionToken, sessionData, time.Until(session.ExpiresAt))
	r.redis.Set(ctx, "refresh:"+session.RefreshToken, session.UserID, time.Until(session.ExpiresAt))

	return nil
}

func (r *SessionRepository) Revoke(sessionToken string) error {
	// Get session first to clean up Redis
	session, err := r.GetByToken(sessionToken)
	if err == nil {
		ctx := context.Background()
		r.redis.Del(ctx, "session:"+sessionToken)
		r.redis.Del(ctx, "refresh:"+session.RefreshToken)
	}

	// Deactivate in PostgreSQL
	query := `UPDATE user_sessions SET is_active = false WHERE session_token = $1`
	_, err = r.db.Exec(query, sessionToken)
	return err
}

func (r *SessionRepository) RevokeByID(sessionID string) error {
	// Get session first to clean up Redis
	query := `
		SELECT session_token, refresh_token
		FROM user_sessions
		WHERE id = $1 AND is_active = true
	`

	var sessionToken, refreshToken string
	err := r.db.QueryRow(query, sessionID).Scan(&sessionToken, &refreshToken)
	if err == nil {
		ctx := context.Background()
		r.redis.Del(ctx, "session:"+sessionToken)
		r.redis.Del(ctx, "refresh:"+refreshToken)
	}

	// Deactivate in PostgreSQL
	updateQuery := `UPDATE user_sessions SET is_active = false WHERE id = $1`
	_, err = r.db.Exec(updateQuery, sessionID)
	return err
}

func (r *SessionRepository) RevokeAllUserSessions(userID string) error {
	// Get all active sessions for user
	query := `
		SELECT session_token, refresh_token
		FROM user_sessions
		WHERE user_id = $1 AND is_active = true
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	ctx := context.Background()
	var tokensToDelete []string

	for rows.Next() {
		var sessionToken, refreshToken string
		if err := rows.Scan(&sessionToken, &refreshToken); err == nil {
			tokensToDelete = append(tokensToDelete, "session:"+sessionToken)
			tokensToDelete = append(tokensToDelete, "refresh:"+refreshToken)
		}
	}

	// Clean up Redis
	if len(tokensToDelete) > 0 {
		r.redis.Del(ctx, tokensToDelete...)
	}

	// Deactivate all sessions in PostgreSQL
	updateQuery := `UPDATE user_sessions SET is_active = false WHERE user_id = $1`
	_, err = r.db.Exec(updateQuery, userID)
	return err
}

func (r *SessionRepository) GetUserSessions(userID string) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, session_token, ip_address, user_agent, expires_at, is_active, created_at
		FROM user_sessions
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(
			&session.ID, &session.UserID, &session.SessionToken,
			&session.IPAddress, &session.UserAgent, &session.ExpiresAt,
			&session.IsActive, &session.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Don't return the actual tokens in the list
		session.SessionToken = ""
		session.RefreshToken = ""
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *SessionRepository) CleanupExpired() error {
	// Clean up PostgreSQL
	query := `UPDATE user_sessions SET is_active = false WHERE expires_at <= NOW() AND is_active = true`
	_, err := r.db.Exec(query)
	
	// Redis keys will expire automatically
	return err
}