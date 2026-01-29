package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/google/uuid"
)

// AggregatorRepository manages aggregator settings in the database
type AggregatorRepository struct {
	db *sql.DB
}

// NewAggregatorRepository creates a new aggregator repository
func NewAggregatorRepository(db *sql.DB) *AggregatorRepository {
	return &AggregatorRepository{db: db}
}

// InitSchema creates the aggregator_settings table and seeds default data
func (r *AggregatorRepository) InitSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS aggregator_settings (
			id VARCHAR(36) PRIMARY KEY,
			provider_code VARCHAR(50) UNIQUE NOT NULL,
			provider_name VARCHAR(100) NOT NULL,
			logo_url VARCHAR(255),
			is_enabled BOOLEAN DEFAULT true,
			deposit_enabled BOOLEAN DEFAULT true,
			withdraw_enabled BOOLEAN DEFAULT true,
			supported_regions TEXT,
			priority INT DEFAULT 50,
			min_amount DECIMAL(20, 8) DEFAULT 0,
			max_amount DECIMAL(20, 8) DEFAULT 0,
			fee_percent DECIMAL(10, 4) DEFAULT 0,
			fee_fixed DECIMAL(20, 8) DEFAULT 0,
			fee_currency VARCHAR(10) DEFAULT 'USD',
			description TEXT,
			maintenance_mode BOOLEAN DEFAULT false,
			maintenance_msg TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create aggregator_settings table: %w", err)
	}

	log.Println("[AggregatorRepo] Table created/verified, seeding defaults...")
	return r.seedDefaults()
}

// seedDefaults inserts default aggregator settings if not exist
func (r *AggregatorRepository) seedDefaults() error {
	defaults := models.DefaultAggregators()

	for _, agg := range defaults {
		// Check if exists
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM aggregator_settings WHERE provider_code = $1)", agg.ProviderCode).Scan(&exists)
		if err != nil {
			log.Printf("[AggregatorRepo] Error checking %s: %v", agg.ProviderCode, err)
			continue
		}
		if exists {
			continue
		}

		// Insert new
		id := uuid.New().String()
		regionsJSON, _ := json.Marshal(agg.SupportedRegions)

		_, err = r.db.Exec(`
			INSERT INTO aggregator_settings 
				(id, provider_code, provider_name, logo_url, is_enabled, deposit_enabled, withdraw_enabled, 
				 supported_regions, priority, min_amount, max_amount, fee_percent, fee_fixed, fee_currency)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, id, agg.ProviderCode, agg.ProviderName, agg.LogoURL, agg.IsEnabled, agg.DepositEnabled, agg.WithdrawEnabled,
			string(regionsJSON), agg.Priority, agg.MinAmount, agg.MaxAmount, agg.FeePercent, agg.FeeFixed, agg.FeeCurrency)

		if err != nil {
			log.Printf("[AggregatorRepo] Failed to seed %s: %v", agg.ProviderCode, err)
		} else {
			log.Printf("[AggregatorRepo] ✅ Seeded aggregator: %s", agg.ProviderName)
		}
	}

	return nil
}

// GetAll returns all aggregator settings
func (r *AggregatorRepository) GetAll() ([]models.AggregatorSetting, error) {
	rows, err := r.db.Query(`
		SELECT id, provider_code, provider_name, logo_url, is_enabled, deposit_enabled, withdraw_enabled,
			   supported_regions, priority, min_amount, max_amount, fee_percent, fee_fixed, fee_currency,
			   description, maintenance_mode, maintenance_msg, created_at, updated_at
		FROM aggregator_settings
		ORDER BY priority DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregators []models.AggregatorSetting
	for rows.Next() {
		var a models.AggregatorSetting
		var regionsJSON, desc, maintenanceMsg sql.NullString

		err := rows.Scan(
			&a.ID, &a.ProviderCode, &a.ProviderName, &a.LogoURL, &a.IsEnabled, &a.DepositEnabled, &a.WithdrawEnabled,
			&regionsJSON, &a.Priority, &a.MinAmount, &a.MaxAmount, &a.FeePercent, &a.FeeFixed, &a.FeeCurrency,
			&desc, &a.MaintenanceMode, &maintenanceMsg, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if regionsJSON.Valid {
			json.Unmarshal([]byte(regionsJSON.String), &a.SupportedRegions)
		}
		if desc.Valid {
			a.Description = desc.String
		}
		if maintenanceMsg.Valid {
			a.MaintenanceMsg = maintenanceMsg.String
		}

		aggregators = append(aggregators, a)
	}

	return aggregators, nil
}

// GetByCode returns an aggregator by its provider code
func (r *AggregatorRepository) GetByCode(code string) (*models.AggregatorSetting, error) {
	var a models.AggregatorSetting
	var regionsJSON, desc, maintenanceMsg sql.NullString

	err := r.db.QueryRow(`
		SELECT id, provider_code, provider_name, logo_url, is_enabled, deposit_enabled, withdraw_enabled,
			   supported_regions, priority, min_amount, max_amount, fee_percent, fee_fixed, fee_currency,
			   description, maintenance_mode, maintenance_msg, created_at, updated_at
		FROM aggregator_settings WHERE provider_code = $1
	`, code).Scan(
		&a.ID, &a.ProviderCode, &a.ProviderName, &a.LogoURL, &a.IsEnabled, &a.DepositEnabled, &a.WithdrawEnabled,
		&regionsJSON, &a.Priority, &a.MinAmount, &a.MaxAmount, &a.FeePercent, &a.FeeFixed, &a.FeeCurrency,
		&desc, &a.MaintenanceMode, &maintenanceMsg, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if regionsJSON.Valid {
		json.Unmarshal([]byte(regionsJSON.String), &a.SupportedRegions)
	}
	if desc.Valid {
		a.Description = desc.String
	}
	if maintenanceMsg.Valid {
		a.MaintenanceMsg = maintenanceMsg.String
	}

	return &a, nil
}

// GetEnabled returns all enabled aggregators for frontend
func (r *AggregatorRepository) GetEnabled() ([]models.AggregatorSetting, error) {
	rows, err := r.db.Query(`
		SELECT id, provider_code, provider_name, logo_url, is_enabled, deposit_enabled, withdraw_enabled,
			   supported_regions, priority, min_amount, max_amount, fee_percent, fee_fixed, fee_currency,
			   description, maintenance_mode, maintenance_msg, created_at, updated_at
		FROM aggregator_settings
		WHERE is_enabled = true
		ORDER BY priority DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aggregators []models.AggregatorSetting
	for rows.Next() {
		var a models.AggregatorSetting
		var regionsJSON, desc, maintenanceMsg sql.NullString

		err := rows.Scan(
			&a.ID, &a.ProviderCode, &a.ProviderName, &a.LogoURL, &a.IsEnabled, &a.DepositEnabled, &a.WithdrawEnabled,
			&regionsJSON, &a.Priority, &a.MinAmount, &a.MaxAmount, &a.FeePercent, &a.FeeFixed, &a.FeeCurrency,
			&desc, &a.MaintenanceMode, &maintenanceMsg, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if regionsJSON.Valid {
			json.Unmarshal([]byte(regionsJSON.String), &a.SupportedRegions)
		}
		if desc.Valid {
			a.Description = desc.String
		}
		if maintenanceMsg.Valid {
			a.MaintenanceMsg = maintenanceMsg.String
		}

		aggregators = append(aggregators, a)
	}

	return aggregators, nil
}

// GetEnabledForDeposit returns aggregators enabled for deposits
func (r *AggregatorRepository) GetEnabledForDeposit(country string) ([]models.AggregatorForFrontend, error) {
	aggregators, err := r.GetEnabled()
	if err != nil {
		return nil, err
	}

	var result []models.AggregatorForFrontend
	for _, a := range aggregators {
		if !a.DepositEnabled {
			continue
		}

		// Check if country is supported
		if !r.isCountrySupported(a.SupportedRegions, country) {
			continue
		}

		result = append(result, models.AggregatorForFrontend{
			Code:            a.ProviderCode,
			Name:            a.ProviderName,
			LogoURL:         a.LogoURL,
			DepositEnabled:  a.DepositEnabled,
			WithdrawEnabled: a.WithdrawEnabled,
			MinAmount:       a.MinAmount,
			MaxAmount:       a.MaxAmount,
			FeePercent:      a.FeePercent,
			FeeFixed:        a.FeeFixed,
			FeeCurrency:     a.FeeCurrency,
			MaintenanceMode: a.MaintenanceMode,
			MaintenanceMsg:  a.MaintenanceMsg,
		})
	}

	return result, nil
}

// GetEnabledForWithdraw returns aggregators enabled for withdrawals
func (r *AggregatorRepository) GetEnabledForWithdraw(country string) ([]models.AggregatorForFrontend, error) {
	aggregators, err := r.GetEnabled()
	if err != nil {
		return nil, err
	}

	var result []models.AggregatorForFrontend
	for _, a := range aggregators {
		if !a.WithdrawEnabled {
			continue
		}

		// Check if country is supported
		if !r.isCountrySupported(a.SupportedRegions, country) {
			continue
		}

		result = append(result, models.AggregatorForFrontend{
			Code:            a.ProviderCode,
			Name:            a.ProviderName,
			LogoURL:         a.LogoURL,
			DepositEnabled:  a.DepositEnabled,
			WithdrawEnabled: a.WithdrawEnabled,
			MinAmount:       a.MinAmount,
			MaxAmount:       a.MaxAmount,
			FeePercent:      a.FeePercent,
			FeeFixed:        a.FeeFixed,
			FeeCurrency:     a.FeeCurrency,
			MaintenanceMode: a.MaintenanceMode,
			MaintenanceMsg:  a.MaintenanceMsg,
		})
	}

	return result, nil
}

func (r *AggregatorRepository) isCountrySupported(regions []string, country string) bool {
	if len(regions) == 0 {
		return true
	}
	for _, region := range regions {
		if region == "*" || region == country {
			return true
		}
	}
	return false
}

// Update updates an aggregator's settings
func (r *AggregatorRepository) Update(code string, req *models.UpdateAggregatorRequest) error {
	// Build dynamic update query
	setClauses := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argCount := 2

	if req.IsEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_enabled = $%d", argCount))
		args = append(args, *req.IsEnabled)
		argCount++
	}
	if req.DepositEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("deposit_enabled = $%d", argCount))
		args = append(args, *req.DepositEnabled)
		argCount++
	}
	if req.WithdrawEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("withdraw_enabled = $%d", argCount))
		args = append(args, *req.WithdrawEnabled)
		argCount++
	}
	if req.Priority != nil {
		setClauses = append(setClauses, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *req.Priority)
		argCount++
	}
	if req.MinAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("min_amount = $%d", argCount))
		args = append(args, *req.MinAmount)
		argCount++
	}
	if req.MaxAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_amount = $%d", argCount))
		args = append(args, *req.MaxAmount)
		argCount++
	}
	if req.FeePercent != nil {
		setClauses = append(setClauses, fmt.Sprintf("fee_percent = $%d", argCount))
		args = append(args, *req.FeePercent)
		argCount++
	}
	if req.FeeFixed != nil {
		setClauses = append(setClauses, fmt.Sprintf("fee_fixed = $%d", argCount))
		args = append(args, *req.FeeFixed)
		argCount++
	}
	if req.FeeCurrency != nil {
		setClauses = append(setClauses, fmt.Sprintf("fee_currency = $%d", argCount))
		args = append(args, *req.FeeCurrency)
		argCount++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}
	if req.MaintenanceMode != nil {
		setClauses = append(setClauses, fmt.Sprintf("maintenance_mode = $%d", argCount))
		args = append(args, *req.MaintenanceMode)
		argCount++
	}
	if req.MaintenanceMsg != nil {
		setClauses = append(setClauses, fmt.Sprintf("maintenance_msg = $%d", argCount))
		args = append(args, *req.MaintenanceMsg)
		argCount++
	}

	args = append(args, code)

	query := fmt.Sprintf("UPDATE aggregator_settings SET %s WHERE provider_code = $%d",
		joinStrings(setClauses, ", "), argCount)

	_, err := r.db.Exec(query, args...)
	return err
}

// SetEnabled quickly enables or disables an aggregator
func (r *AggregatorRepository) SetEnabled(code string, enabled bool) error {
	_, err := r.db.Exec(`
		UPDATE aggregator_settings 
		SET is_enabled = $1, updated_at = $2 
		WHERE provider_code = $3
	`, enabled, time.Now(), code)
	return err
}

// SetMaintenanceMode sets maintenance mode on/off
func (r *AggregatorRepository) SetMaintenanceMode(code string, maintenance bool, msg string) error {
	_, err := r.db.Exec(`
		UPDATE aggregator_settings 
		SET maintenance_mode = $1, maintenance_msg = $2, updated_at = $3 
		WHERE provider_code = $4
	`, maintenance, msg, time.Now(), code)
	return err
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
