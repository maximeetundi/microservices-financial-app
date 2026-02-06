package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/lib/pq"
)

type AggregatorInstanceRepository struct {
	db *sql.DB
}

func NewAggregatorInstanceRepository(db *sql.DB) *AggregatorInstanceRepository {
	return &AggregatorInstanceRepository{db: db}
}

// GetByID retrieves an instance by ID - tries provider_instances first, then aggregator_instances
func (r *AggregatorInstanceRepository) GetByID(ctx context.Context, id string) (*models.AggregatorInstance, error) {
	// First try provider_instances (admin-service table with vault_secret_path)
	providerQuery := `
		SELECT pi.id, pp.id as aggregator_id, pi.name, pi.hot_wallet_id, pi.vault_secret_path,
		       pi.is_active, pi.priority, pi.is_global, pi.is_paused, pi.pause_reason,
		       pi.request_count, pi.last_used_at, pi.created_at, pi.updated_at,
		       pp.name as provider_code
		FROM provider_instances pi
		JOIN payment_providers pp ON pi.provider_id = pp.id
		WHERE pi.id = $1
	`

	var instance models.AggregatorInstance
	var vaultPath sql.NullString
	var pauseReason sql.NullString
	var lastUsedAt sql.NullTime
	var providerCode string

	err := r.db.QueryRowContext(ctx, providerQuery, id).Scan(
		&instance.ID, &instance.AggregatorID, &instance.InstanceName, &instance.HotWalletID,
		&vaultPath, &instance.Enabled, &instance.Priority, &instance.IsGlobal, &instance.IsPaused,
		&pauseReason, &instance.TotalTransactions, &lastUsedAt,
		&instance.CreatedAt, &instance.UpdatedAt, &providerCode,
	)

	if err == nil {
		// Found in provider_instances
		if vaultPath.Valid {
			instance.VaultSecretPath = vaultPath.String
		}
		if pauseReason.Valid {
			instance.PauseReason = &pauseReason.String
		}
		if lastUsedAt.Valid {
			instance.LastUsedAt = &lastUsedAt.Time
		}
		instance.APICredentials = make(map[string]string)
		log.Printf("[AggregatorInstanceRepo] Found instance %s in provider_instances, vault_path: %s", id, instance.VaultSecretPath)
		return &instance, nil
	}

	// Fallback to aggregator_instances table
	query := `
		SELECT id, aggregator_id, instance_name, hot_wallet_id, api_credentials,
		       enabled, priority, min_balance, max_balance, daily_limit, monthly_limit,
		       daily_usage, monthly_usage, usage_reset_date, restricted_countries,
		       is_test_mode, total_transactions, total_volume, last_used_at,
		       notes, created_at, updated_at, created_by
		FROM aggregator_instances
		WHERE id = $1
	`

	var apiCredsBytes []byte
	var restrictedCountries pq.StringArray

	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&instance.ID, &instance.AggregatorID, &instance.InstanceName, &instance.HotWalletID,
		&apiCredsBytes, &instance.Enabled, &instance.Priority, &instance.MinBalance,
		&instance.MaxBalance, &instance.DailyLimit, &instance.MonthlyLimit,
		&instance.DailyUsage, &instance.MonthlyUsage, &instance.UsageResetDate,
		&restrictedCountries, &instance.IsTestMode, &instance.TotalTransactions,
		&instance.TotalVolume, &instance.LastUsedAt, &instance.Notes,
		&instance.CreatedAt, &instance.UpdatedAt, &instance.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	// Parse JSON credentials
	if len(apiCredsBytes) > 0 {
		if err := json.Unmarshal(apiCredsBytes, &instance.APICredentials); err != nil {
			return nil, fmt.Errorf("parse api credentials: %w", err)
		}
	} else {
		instance.APICredentials = make(map[string]string)
	}

	instance.RestrictedCountries = restrictedCountries

	return &instance, nil
}

// GetAvailableInstancesForAggregator gets all available instances for an aggregator
// filtered by country and sorted by priority
func (r *AggregatorInstanceRepository) GetAvailableInstancesForAggregator(
	ctx context.Context,
	aggregatorID string,
	country string,
	amount float64,
) ([]*models.AggregatorInstanceWithDetails, error) {
	query := `
		SELECT * FROM aggregator_instances_with_details
		WHERE aggregator_id = $1
		  AND aggregator_enabled = true
		  AND availability_status = 'available'
		  AND (
		      restricted_countries IS NULL
		      OR array_length(restricted_countries, 1) IS NULL
		      OR $2 = ANY(restricted_countries)
		  )
		  AND (daily_limit IS NULL OR (daily_limit - daily_usage) >= $3)
		  AND (monthly_limit IS NULL OR (monthly_limit - monthly_usage) >= $3)
		ORDER BY priority DESC, hot_wallet_balance DESC
	`

	rows, err := r.db.QueryContext(ctx, query, aggregatorID, country, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.AggregatorInstanceWithDetails
	for rows.Next() {
		instance, err := r.scanInstanceWithDetails(rows)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, rows.Err()
}

// GetBestInstanceForProvider returns the best instance for a given provider code and country
// First tries provider_instances (admin-service), then aggregator_instances
func (r *AggregatorInstanceRepository) GetBestInstanceForProvider(
	ctx context.Context,
	providerCode string,
	country string,
	amount float64,
) (*models.AggregatorInstanceWithDetails, error) {
	// First, try to get from provider_instances (admin-service table)
	providerQuery := `
		SELECT pi.id, pi.name, pi.is_active, pi.priority, pi.is_paused, pi.is_global,
		       pi.pause_reason, pi.paused_at, pi.vault_secret_path,
		       pi.request_count, pi.last_used_at, pi.created_at, pi.updated_at,
		       pp.id as aggregator_id, pp.name as provider_code, pp.display_name as provider_name,
		       pp.logo_url as provider_logo, pp.is_active as aggregator_enabled,
		       COALESCE(pi.hot_wallet_id, '') as hot_wallet_id,
		       COALESCE(pa.currency, 'XOF') as hot_wallet_currency,
		       COALESCE(pa.balance, 0) as hot_wallet_balance
		FROM provider_instances pi
		JOIN payment_providers pp ON pi.provider_id = pp.id
		LEFT JOIN platform_accounts pa ON pi.hot_wallet_id = pa.id
		WHERE pp.name = $1
		  AND pp.is_active = true
		  AND pi.is_active = true
		  AND pi.is_paused = false
		ORDER BY pi.is_primary DESC, pi.priority DESC
		LIMIT 1
	`

	var instance models.AggregatorInstanceWithDetails
	var vaultPath sql.NullString
	var pauseReason sql.NullString
	var pausedAt sql.NullTime
	var lastUsedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, providerQuery, providerCode).Scan(
		&instance.ID, &instance.InstanceName, &instance.Enabled, &instance.Priority,
		&instance.IsPaused, &instance.IsGlobal, &pauseReason, &pausedAt, &vaultPath,
		&instance.TotalTransactions, &lastUsedAt, &instance.CreatedAt, &instance.UpdatedAt,
		&instance.AggregatorID, &instance.ProviderCode, &instance.ProviderName,
		&instance.ProviderLogo, &instance.AggregatorEnabled, &instance.HotWalletID,
		&instance.HotWalletCurrency, &instance.HotWalletBalance,
	)

	if err == nil {
		// Found in provider_instances
		if vaultPath.Valid {
			instance.VaultSecretPath = vaultPath.String
		}
		if pauseReason.Valid {
			instance.PauseReason = &pauseReason.String
		}
		if pausedAt.Valid {
			instance.PausedAt = &pausedAt.Time
		}
		if lastUsedAt.Valid {
			instance.LastUsedAt = &lastUsedAt.Time
		}
		instance.AvailabilityStatus = models.WalletAvailable
		instance.APICredentials = make(map[string]string)

		// IMPORTANT:
		// provider_instances (admin UI table) and aggregator_instances (transfer routing table)
		// are different schemas.
		// Wallet selection for deposits/withdrawals uses select_best_wallet_for_instance(),
		// which reads aggregator_instance_wallets (linked to aggregator_instances).
		// If we return a provider_instances ID here, wallet selection will fail.
		// So we map to the corresponding aggregator_instances ID.
		var aggSettingsID string
		if err := r.db.QueryRowContext(ctx, `SELECT id FROM aggregator_settings WHERE provider_code = $1`, providerCode).Scan(&aggSettingsID); err == nil {
			var aggInstanceID string
			if err := r.db.QueryRowContext(ctx, `
				SELECT id
				FROM aggregator_instances
				WHERE aggregator_id = $1 AND instance_name = $2
				ORDER BY priority DESC
				LIMIT 1
			`, aggSettingsID, instance.InstanceName).Scan(&aggInstanceID); err == nil {
				instance.ID = aggInstanceID
				instance.AggregatorID = aggSettingsID
			}
		}

		log.Printf("[AggregatorInstanceRepo] Found instance for %s in provider_instances: %s (vault: %s)",
			providerCode, instance.ID, instance.VaultSecretPath)
		return &instance, nil
	}

	log.Printf("[AggregatorInstanceRepo] provider_instances query failed for %s: %v, trying aggregator_instances", providerCode, err)

	// Fallback to aggregator_instances_with_details view
	query := `
		SELECT * FROM aggregator_instances_with_details
		WHERE provider_code = $1
		  AND aggregator_enabled = true
		  AND is_paused = false
		  AND availability_status = 'available'
		  AND (
		      is_global = true
		      OR restricted_countries IS NULL
		      OR array_length(restricted_countries, 1) IS NULL
		      OR $2 = ANY(restricted_countries)
		  )
		  AND (daily_limit IS NULL OR (daily_limit - daily_usage) >= $3)
		  AND (monthly_limit IS NULL OR (monthly_limit - monthly_usage) >= $3)
		ORDER BY priority DESC, hot_wallet_balance DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, providerCode, country, amount)
	return r.scanInstanceWithDetails(row)
}

// Create creates a new instance
func (r *AggregatorInstanceRepository) Create(ctx context.Context, req *models.CreateAggregatorInstanceRequest, createdBy string) (*models.AggregatorInstance, error) {
	query := `
		INSERT INTO aggregator_instances (
			aggregator_id, instance_name, hot_wallet_id, api_credentials,
			priority, min_balance, max_balance, daily_limit, monthly_limit,
			restricted_countries, is_test_mode, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	apiCredsBytes, err := json.Marshal(req.APICredentials)
	if err != nil {
		return nil, fmt.Errorf("marshal api credentials: %w", err)
	}

	priority := 50
	if req.Priority != nil {
		priority = *req.Priority
	}

	isTestMode := true
	if req.IsTestMode != nil {
		isTestMode = *req.IsTestMode
	}

	var instance models.AggregatorInstance
	err = r.db.QueryRowContext(ctx, query,
		req.AggregatorID, req.InstanceName, req.HotWalletID, apiCredsBytes,
		priority, req.MinBalance, req.MaxBalance, req.DailyLimit, req.MonthlyLimit,
		pq.Array(req.RestrictedCountries), isTestMode, req.Notes, createdBy,
	).Scan(&instance.ID, &instance.CreatedAt, &instance.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Populate the rest of the fields
	instance.AggregatorID = req.AggregatorID
	instance.InstanceName = req.InstanceName
	instance.HotWalletID = req.HotWalletID
	instance.APICredentials = req.APICredentials
	instance.Priority = priority
	instance.MinBalance = req.MinBalance
	instance.MaxBalance = req.MaxBalance
	instance.DailyLimit = req.DailyLimit
	instance.MonthlyLimit = req.MonthlyLimit
	instance.RestrictedCountries = req.RestrictedCountries
	instance.IsTestMode = isTestMode
	instance.Notes = req.Notes
	instance.Enabled = true

	return &instance, nil
}

// IncrementUsage increments the usage for an instance
func (r *AggregatorInstanceRepository) IncrementUsage(ctx context.Context, instanceID string, amount float64) error {
	query := `
		UPDATE aggregator_instances
		SET daily_usage = daily_usage + $1,
		    monthly_usage = monthly_usage + $1,
		    total_transactions = total_transactions + 1,
		    total_volume = total_volume + $1,
		    last_used_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, amount, instanceID)
	return err
}

// ResetDailyUsage resets daily usage for all instances
func (r *AggregatorInstanceRepository) ResetDailyUsage(ctx context.Context) error {
	query := `
		UPDATE aggregator_instances
		SET daily_usage = 0,
		    usage_reset_date = CURRENT_DATE
		WHERE usage_reset_date < CURRENT_DATE
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// LogTransaction logs a transaction for an instance
func (r *AggregatorInstanceRepository) LogTransaction(
	ctx context.Context,
	instanceID string,
	transactionID string,
	transactionType string,
	amount float64,
	currency string,
	status string,
	providerReference string,
) error {
	query := `
		INSERT INTO aggregator_instance_transactions (
			instance_id, transaction_id, transaction_type, amount, currency, status, provider_reference
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query, instanceID, transactionID, transactionType, amount, currency, status, providerReference)
	return err
}

// GetInstanceWithWallets retrieves instance with all linked wallets
func (r *AggregatorInstanceRepository) GetInstanceWithWallets(ctx context.Context, instanceID string) (*models.AggregatorInstance, error) {
	// Get base instance
	instance, err := r.GetByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Get linked wallets
	wallets, err := r.GetInstanceWallets(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	instance.Wallets = wallets
	return instance, nil
}

// GetInstanceWallets gets all wallets linked to an instance
func (r *AggregatorInstanceRepository) GetInstanceWallets(ctx context.Context, instanceID string) ([]models.InstanceWallet, error) {
	query := `
		SELECT
			aiw.id, aiw.instance_id, aiw.hot_wallet_id, aiw.is_primary,
			aiw.priority, aiw.min_balance, aiw.max_balance,
			aiw.auto_recharge_enabled, aiw.recharge_threshold, aiw.recharge_target,
			aiw.recharge_source_wallet_id, aiw.total_deposits, aiw.total_withdrawals,
			aiw.transaction_count, aiw.last_used_at, aiw.enabled,
			aiw.created_at, aiw.updated_at,
			pa.currency, pa.balance
		FROM aggregator_instance_wallets aiw
		JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id
		WHERE aiw.instance_id = $1
		ORDER BY aiw.is_primary DESC, aiw.priority DESC
	`

	rows, err := r.db.QueryContext(ctx, query, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.InstanceWallet
	for rows.Next() {
		var wallet models.InstanceWallet
		err := rows.Scan(
			&wallet.ID, &wallet.InstanceID, &wallet.HotWalletID, &wallet.IsPrimary,
			&wallet.Priority, &wallet.MinBalance, &wallet.MaxBalance,
			&wallet.AutoRechargeEnabled, &wallet.RechargeThreshold, &wallet.RechargeTarget,
			&wallet.RechargeSourceWalletID, &wallet.TotalDeposits, &wallet.TotalWithdrawals,
			&wallet.TransactionCount, &wallet.LastUsedAt, &wallet.Enabled,
			&wallet.CreatedAt, &wallet.UpdatedAt,
			&wallet.WalletCurrency, &wallet.WalletBalance,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, rows.Err()
}

// LinkWallet links a hot wallet to an instance
func (r *AggregatorInstanceRepository) LinkWallet(ctx context.Context, req *models.CreateInstanceWalletRequest) error {
	query := `
		INSERT INTO aggregator_instance_wallets (
			instance_id, hot_wallet_id, is_primary, priority, min_balance, max_balance,
			auto_recharge_enabled, recharge_threshold, recharge_target, recharge_source_wallet_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	isPrimary := false
	if req.IsPrimary != nil {
		isPrimary = *req.IsPrimary
	}

	priority := 50
	if req.Priority != nil {
		priority = *req.Priority
	}

	autoRecharge := false
	if req.AutoRechargeEnabled != nil {
		autoRecharge = *req.AutoRechargeEnabled
	}

	_, err := r.db.ExecContext(ctx, query,
		req.InstanceID, req.HotWalletID, isPrimary, priority, req.MinBalance, req.MaxBalance,
		autoRecharge, req.RechargeThreshold, req.RechargeTarget, req.RechargeSourceWalletID,
	)

	return err
}

// UnlinkWallet removes a wallet from an instance
func (r *AggregatorInstanceRepository) UnlinkWallet(ctx context.Context, instanceID, walletID string) error {
	query := `DELETE FROM aggregator_instance_wallets WHERE instance_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, instanceID, walletID)
	return err
}

// UpdateWallet updates wallet configuration
func (r *AggregatorInstanceRepository) UpdateWallet(ctx context.Context, instanceID, walletID string, updates map[string]interface{}) error {
	// Build dynamic UPDATE query
	query := `UPDATE aggregator_instance_wallets SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{instanceID, walletID}
	argIndex := 3

	for key, value := range updates {
		query += fmt.Sprintf(", %s = $%d", key, argIndex)
		args = append(args, value)
		argIndex++
	}

	query += ` WHERE instance_id = $1 AND id = $2`

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// SetPaused sets the pause status of an instance
func (r *AggregatorInstanceRepository) SetPaused(
	ctx context.Context,
	instanceID string,
	paused bool,
	reason string,
	adminID string,
) error {
	var query string
	var args []interface{}

	if paused {
		query = `
			UPDATE aggregator_instances
			SET is_paused = true,
			    pause_reason = $2,
			    paused_at = CURRENT_TIMESTAMP,
			    paused_by = $3,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		args = []interface{}{instanceID, reason, adminID}
	} else {
		query = `
			UPDATE aggregator_instances
			SET is_paused = false,
			    pause_reason = NULL,
			    paused_at = NULL,
			    paused_by = NULL,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		args = []interface{}{instanceID}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateInstance updates an instance's basic fields
func (r *AggregatorInstanceRepository) UpdateInstance(
	ctx context.Context,
	instanceID string,
	name string,
	priority int,
	isActive bool,
	isPrimary bool,
	isGlobal bool,
	vaultPath string,
) error {
	query := `
		UPDATE aggregator_instances
		SET instance_name = $2,
		    priority = $3,
		    enabled = $4,
		    is_global = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, instanceID, name, priority, isActive, isGlobal)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Helper function to scan instance with details
// Note: The view columns must match this order exactly
func (r *AggregatorInstanceRepository) scanInstanceWithDetails(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.AggregatorInstanceWithDetails, error) {
	var instance models.AggregatorInstanceWithDetails
	var restrictedCountries pq.StringArray

	err := scanner.Scan(
		&instance.ID, &instance.InstanceName, &instance.Enabled, &instance.Priority,
		&instance.IsTestMode, &instance.IsPaused, &instance.IsGlobal, &instance.PauseReason,
		&instance.PausedAt, &restrictedCountries, &instance.DailyLimit, &instance.MonthlyLimit,
		&instance.DailyUsage, &instance.MonthlyUsage, &instance.TotalTransactions,
		&instance.TotalVolume, &instance.LastUsedAt, &instance.CreatedAt, &instance.UpdatedAt,
		&instance.AggregatorID, &instance.ProviderCode, &instance.ProviderName,
		&instance.ProviderLogo, &instance.AggregatorEnabled, &instance.HotWalletID,
		&sql.NullString{}, // account_type
		&instance.HotWalletCurrency, &instance.HotWalletBalance,
		&instance.MinBalance, &instance.MaxBalance, &instance.AvailabilityStatus,
	)

	if err != nil {
		return nil, err
	}

	instance.RestrictedCountries = restrictedCountries
	return &instance, nil
}
