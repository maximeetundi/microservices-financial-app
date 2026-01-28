package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/google/uuid"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// InitSchema creates the table if it doesn't exist and seeds default data
func (r *SettingsRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS system_settings (
		id VARCHAR(36) PRIMARY KEY,
		key VARCHAR(100) UNIQUE NOT NULL,
		value TEXT NOT NULL,
		description TEXT,
		type VARCHAR(20) NOT NULL DEFAULT 'string',
		category VARCHAR(50) NOT NULL DEFAULT 'general',
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create system_settings table: %w", err)
	}

	return r.seedDefaults()
}

func (r *SettingsRepository) seedDefaults() error {
	defaults := []models.SystemSetting{
		{
			Key:         "rate_update_interval",
			Value:       "60",
			Description: "Intervalle de mise à jour des taux de change (en secondes)",
			Type:        "integer",
			Category:    "exchange",
		},
		{
			Key:         "crypto_network",
			Value:       "mainnet",
			Description: "Réseau crypto à utiliser (mainnet / testnet)",
			Type:        "string",
			Category:    "crypto",
		},
		{
			Key:         "binance_test_mode",
			Value:       "true",
			Description: "Utiliser le mode test de Binance",
			Type:        "boolean",
			Category:    "exchange",
		},
		{
			Key:         "rate_update_enabled",
			Value:       "true",
			Description: "Activer la mise à jour automatique des taux",
			Type:        "boolean",
			Category:    "exchange",
		},
	}

	for _, setting := range defaults {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM system_settings WHERE key = $1)", setting.Key).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			id := uuid.New().String()
			_, err := r.db.Exec(`
				INSERT INTO system_settings (id, key, value, description, type, category, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, id, setting.Key, setting.Value, setting.Description, setting.Type, setting.Category, time.Now())
			if err != nil {
				return fmt.Errorf("failed to seed setting %s: %w", setting.Key, err)
			}
		}
	}
	return nil
}

func (r *SettingsRepository) GetAll() ([]models.SystemSetting, error) {
	query := `SELECT id, key, value, description, type, category, updated_at FROM system_settings ORDER BY category, key`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var s models.SystemSetting
		err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Description, &s.Type, &s.Category, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SettingsRepository) GetByKey(key string) (*models.SystemSetting, error) {
	query := `SELECT id, key, value, description, type, category, updated_at FROM system_settings WHERE key = $1`
	var s models.SystemSetting
	err := r.db.QueryRow(query, key).Scan(&s.ID, &s.Key, &s.Value, &s.Description, &s.Type, &s.Category, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepository) Update(setting *models.SystemSetting) error {
	query := `
		UPDATE system_settings 
		SET value = $1, description = $2, updated_at = NOW()
		WHERE key = $3
	`
	result, err := r.db.Exec(query, setting.Value, setting.Description, setting.Key)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("setting with key '%s' not found", setting.Key)
	}
	return nil
}

func (r *SettingsRepository) GetByCategory(category string) ([]models.SystemSetting, error) {
	query := `SELECT id, key, value, description, type, category, updated_at FROM system_settings WHERE category = $1 ORDER BY key`
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var s models.SystemSetting
		err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Description, &s.Type, &s.Category, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}
