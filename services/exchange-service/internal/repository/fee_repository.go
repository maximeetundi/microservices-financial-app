package repository

import (
	"database/sql"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
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
			Key:         "exchange_crypto_to_fiat",
			Description: "Fee for converting Crypto to Fiat",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "exchange_fiat_to_crypto",
			Description: "Fee for converting Fiat to Crypto",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "exchange_crypto_to_crypto",
			Description: "Fee for converting Crypto to Crypto",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "exchange_fiat_to_fiat",
			Description: "Fee for converting Fiat to Fiat",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "trading_buy",
			Description: "Fee for buying assets (Trading)",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		{
			Key:         "trading_sell",
			Description: "Fee for selling assets (Trading)",
			Type:        "percentage",
			Percentage:  0, // Free by default
		},
		// New fees
		{
			Key:         "donation_fee",
			Description: "Fee on donations received",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "ecommerce_purchase_fee",
			Description: "Fee on e-commerce purchases",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "event_ticket_fee",
			Description: "Fee on event ticket sales",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "bill_payment_fee",
			Description: "Fee for paying bills",
			Type:        "flat",
			FixedAmount: 0,
		},
		{
			Key:         "mobile_money_cashin_fee",
			Description: "Fee for Mobile Money deposits",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "mobile_money_cashout_fee",
			Description: "Fee for Mobile Money withdrawals",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "association_membership_fee",
			Description: "Fee on association membership payments",
			Type:        "percentage",
			Percentage:  0,
		},
		{
			Key:         "crowdfunding_contribution_fee",
			Description: "Fee on crowdfunding contributions",
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
