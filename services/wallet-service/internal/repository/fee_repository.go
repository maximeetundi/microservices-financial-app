package repository

import (
	"database/sql"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/google/uuid"
)

type FeeRepository struct {
	db *sql.DB
}

func NewFeeRepository(db *sql.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

// InitSchema creates the table if it doesn't exist and seeds default data
func (r *FeeRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS fee_configs (
		id VARCHAR(36) PRIMARY KEY,
		key VARCHAR(50) UNIQUE NOT NULL,
		description TEXT,
		type VARCHAR(20) NOT NULL,
		percentage DECIMAL(10, 4) DEFAULT 0,
		fixed_amount DECIMAL(20, 8) DEFAULT 0,
		currency VARCHAR(10) DEFAULT 'USD',
		min_fee DECIMAL(20, 8) DEFAULT 0,
		max_fee DECIMAL(20, 8) DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create fee_configs table: %w", err)
	}

	return r.seedDefaults()
}

func (r *FeeRepository) seedDefaults() error {
	defaults := []models.FeeConfig{
		{
			Key:         "transfer_internal",
			Description: "Internal wallet to wallet transfer",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "transfer_external",
			Description: "External transfer to bank or mobile money",
			Type:        "hybrid",
			Percentage:  0, // Free by default
			FixedAmount: 0,
		},
		{
			Key:         "merchant_payment",
			Description: "Merchant payment fee",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "deposit",
			Description: "Deposit fee",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "withdrawal",
			Description: "Withdrawal fee",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "crypto_send",
			Description: "Crypto blockchain withdrawal/send fee (Platform fee)",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "transfer_bank",
			Description: "Bank transfer fee",
			Type:        "percentage", // or hybrid
			Percentage:  0,
		},
		{
			Key:         "transfer_mobile",
			Description: "Mobile money transfer fee",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "transfer_international",
			Description: "International transfer fee",
			Type:        "percentage",
			Percentage:  0,
		},
	}

	for _, config := range defaults {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM fee_configs WHERE key = $1)", config.Key).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			id := uuid.New().String()
			_, err := r.db.Exec(`
				INSERT INTO fee_configs (id, key, description, type, percentage, fixed_amount, currency, min_fee, max_fee, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`, id, config.Key, config.Description, config.Type, config.Percentage, config.FixedAmount, "USD", 0, 0, true)
			if err != nil {
				return fmt.Errorf("failed to seed fee config %s: %w", config.Key, err)
			}
		}
	}
	return nil
}

func (r *FeeRepository) GetAll() ([]models.FeeConfig, error) {
	query := `SELECT id, key, description, type, percentage, fixed_amount, currency, min_fee, max_fee, is_active, created_at, updated_at FROM fee_configs ORDER BY key`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.FeeConfig
	for rows.Next() {
		var c models.FeeConfig
		err := rows.Scan(&c.ID, &c.Key, &c.Description, &c.Type, &c.Percentage, &c.FixedAmount, &c.Currency, &c.MinFee, &c.MaxFee, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
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

func (r *FeeRepository) Update(config *models.FeeConfig) error {
	query := `
		UPDATE fee_configs 
		SET description = $1, type = $2, percentage = $3, fixed_amount = $4, currency = $5, min_fee = $6, max_fee = $7, is_active = $8, updated_at = NOW()
		WHERE id = $9
	`
	_, err := r.db.Exec(query, config.Description, config.Type, config.Percentage, config.FixedAmount, config.Currency, config.MinFee, config.MaxFee, config.IsActive, config.ID)
	return err
}
