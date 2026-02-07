package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const MaxDepositNumbersPerUser = 8

type DepositNumber struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Country    string    `json:"country"`
	Phone      string    `json:"phone"`
	Label      string    `json:"label"`
	IsDefault  bool      `json:"is_default"`
	IsVerified bool      `json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DepositNumberRepository struct {
	db *sql.DB
}

func NewDepositNumberRepository(db *sql.DB) *DepositNumberRepository {
	return &DepositNumberRepository{db: db}
}

func (r *DepositNumberRepository) InitSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS user_deposit_numbers (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(100) NOT NULL,
			country VARCHAR(10) NOT NULL,
			phone VARCHAR(50) NOT NULL,
			label VARCHAR(100) DEFAULT '',
			is_default BOOLEAN DEFAULT FALSE,
			is_verified BOOLEAN DEFAULT TRUE,
			verified_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, country, phone)
		);

		CREATE INDEX IF NOT EXISTS idx_user_deposit_numbers_user ON user_deposit_numbers(user_id);
		CREATE INDEX IF NOT EXISTS idx_user_deposit_numbers_phone ON user_deposit_numbers(phone);
		CREATE INDEX IF NOT EXISTS idx_user_deposit_numbers_default ON user_deposit_numbers(user_id, is_default);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *DepositNumberRepository) CountByUser(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user_deposit_numbers WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *DepositNumberRepository) ListByUser(ctx context.Context, userID string) ([]*DepositNumber, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, country, phone, label, COALESCE(is_default, FALSE), COALESCE(is_verified, TRUE), verified_at, created_at, updated_at
		FROM user_deposit_numbers
		WHERE user_id = $1
		ORDER BY COALESCE(is_default, FALSE) DESC, updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*DepositNumber
	for rows.Next() {
		var n DepositNumber
		err := rows.Scan(&n.ID, &n.UserID, &n.Country, &n.Phone, &n.Label, &n.IsDefault, &n.IsVerified, &n.VerifiedAt, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			continue
		}
		res = append(res, &n)
	}
	return res, rows.Err()
}

func (r *DepositNumberRepository) GetByID(ctx context.Context, id string) (*DepositNumber, error) {
	var n DepositNumber
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, country, phone, label, COALESCE(is_default, FALSE), COALESCE(is_verified, TRUE), verified_at, created_at, updated_at
		FROM user_deposit_numbers
		WHERE id = $1
	`, id).Scan(&n.ID, &n.UserID, &n.Country, &n.Phone, &n.Label, &n.IsDefault, &n.IsVerified, &n.VerifiedAt, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *DepositNumberRepository) GetByUserCountryPhone(ctx context.Context, userID, country, phone string) (*DepositNumber, error) {
	var n DepositNumber
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, country, phone, label, COALESCE(is_default, FALSE), COALESCE(is_verified, TRUE), verified_at, created_at, updated_at
		FROM user_deposit_numbers
		WHERE user_id = $1 AND country = $2 AND phone = $3
	`, userID, country, phone).Scan(&n.ID, &n.UserID, &n.Country, &n.Phone, &n.Label, &n.IsDefault, &n.IsVerified, &n.VerifiedAt, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *DepositNumberRepository) Create(ctx context.Context, userID, country, phone, label string, isDefault bool) (*DepositNumber, error) {
	// Enforce max limit
	count, err := r.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= MaxDepositNumbersPerUser {
		return nil, fmt.Errorf("max deposit numbers reached (%d)", MaxDepositNumbersPerUser)
	}

	id := uuid.New().String()
	now := time.Now()

	// If new default, unset previous default
	if isDefault {
		_, _ = r.db.ExecContext(ctx, `UPDATE user_deposit_numbers SET is_default = FALSE, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1`, userID)
	}

	verifiedAt := now
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO user_deposit_numbers (id, user_id, country, phone, label, is_default, is_verified, verified_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,TRUE,$7,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)
	`, id, userID, country, phone, label, isDefault, verifiedAt)
	if err != nil {
		return nil, err
	}

	return &DepositNumber{
		ID:         id,
		UserID:     userID,
		Country:    country,
		Phone:      phone,
		Label:      label,
		IsDefault:  isDefault,
		IsVerified: true,
		VerifiedAt: &verifiedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (r *DepositNumberRepository) Delete(ctx context.Context, userID, id string) error {
	// Do not delete if it's the only number
	count, err := r.CountByUser(ctx, userID)
	if err != nil {
		return err
	}
	if count <= 1 {
		return fmt.Errorf("cannot delete last deposit number")
	}

	res, err := r.db.ExecContext(ctx, `DELETE FROM user_deposit_numbers WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
