package repository

import (
	"database/sql"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
)

type FeeRepository struct {
	db *sql.DB
}

func NewFeeRepository(db *sql.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) GetByKey(key string) (*models.FeeConfig, error) {
	query := `SELECT id, key, description, type, percentage, fixed_amount, currency, min_fee, max_fee, is_active, created_at, updated_at FROM fee_configs WHERE key = $1`
	var c models.FeeConfig
	err := r.db.QueryRow(query, key).Scan(&c.ID, &c.Key, &c.Description, &c.Type, &c.Percentage, &c.FixedAmount, &c.Currency, &c.MinFee, &c.MaxFee, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
