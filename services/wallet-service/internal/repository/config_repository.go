package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// InitSchema creates the necessary tables
func (r *ConfigRepository) InitSchema() error {
	// Table: system_configs
	configQuery := `
		CREATE TABLE IF NOT EXISTS system_configs (
			key VARCHAR(100) PRIMARY KEY,
			value TEXT,
			fixed_amount DECIMAL(30, 18),
			percentage_amount DECIMAL(10, 4),
			is_enabled BOOLEAN DEFAULT true,
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`
	if _, err := r.db.Exec(configQuery); err != nil {
		return fmt.Errorf("failed to create system_configs table: %w", err)
	}

	// Table: user_usage_stats
	usageQuery := `
		CREATE TABLE IF NOT EXISTS user_usage_stats (
			user_id VARCHAR(255) NOT NULL,
			period_key VARCHAR(50) NOT NULL, -- e.g., "DAILY_2024-01-29"
			currency VARCHAR(10) NOT NULL,
			total_amount DECIMAL(30, 18) DEFAULT 0,
			transaction_count INT DEFAULT 0,
			updated_at TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (user_id, period_key, currency)
		);
	`
	if _, err := r.db.Exec(usageQuery); err != nil {
		return fmt.Errorf("failed to create user_usage_stats table: %w", err)
	}

	return nil
}

// UpsertSystemConfig inserts or updates a system configuration
func (r *ConfigRepository) UpsertSystemConfig(cfg *models.SystemConfig) error {
	query := `
		INSERT INTO system_configs (key, value, fixed_amount, percentage_amount, is_enabled, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			fixed_amount = EXCLUDED.fixed_amount,
			percentage_amount = EXCLUDED.percentage_amount,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.Exec(query, cfg.Key, cfg.Value, cfg.FixedAmount, cfg.PercentageAmount, cfg.IsEnabled, time.Now())
	if err != nil {
		return fmt.Errorf("failed to upsert config %s: %w", cfg.Key, err)
	}
	return nil
}

// GetSystemConfigs retrieves all system configurations
func (r *ConfigRepository) GetSystemConfigs() ([]models.SystemConfig, error) {
	query := `SELECT key, value, fixed_amount, percentage_amount, is_enabled, updated_at FROM system_configs`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.SystemConfig
	for rows.Next() {
		var c models.SystemConfig
		var val sql.NullString
		// Handle potential NULLs if needed, though schema has some non-nulls.
		// Value can be empty string.
		if err := rows.Scan(&c.Key, &val, &c.FixedAmount, &c.PercentageAmount, &c.IsEnabled, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if val.Valid {
			c.Value = val.String
		}
		configs = append(configs, c)
	}
	return configs, nil
}

// GetUserUsage retrieves usage stats for a specific user and period
func (r *ConfigRepository) GetUserUsage(userID, periodKey, currency string) (*models.UserUsageStats, error) {
	query := `
		SELECT user_id, period_key, currency, total_amount, transaction_count, updated_at
		FROM user_usage_stats
		WHERE user_id = $1 AND period_key = $2 AND currency = $3
	`
	var stats models.UserUsageStats
	err := r.db.QueryRow(query, userID, periodKey, currency).Scan(
		&stats.UserID, &stats.PeriodKey, &stats.Currency, &stats.TotalAmount, &stats.TransactionCount, &stats.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return &models.UserUsageStats{
			UserID: userID, PeriodKey: periodKey, Currency: currency, TotalAmount: 0, TransactionCount: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// IncrementUserUsage updates the user's usage stats
func (r *ConfigRepository) IncrementUserUsage(userID, periodKey, currency string, amount float64) error {
	query := `
		INSERT INTO user_usage_stats (user_id, period_key, currency, total_amount, transaction_count, updated_at)
		VALUES ($1, $2, $3, $4, 1, $5)
		ON CONFLICT (user_id, period_key, currency) DO UPDATE SET
			total_amount = user_usage_stats.total_amount + EXCLUDED.total_amount,
			transaction_count = user_usage_stats.transaction_count + 1,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.Exec(query, userID, periodKey, currency, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}
	return nil
}
