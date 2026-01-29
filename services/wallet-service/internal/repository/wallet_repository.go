package repository

import (
	"database/sql"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	query := `
		INSERT INTO wallets (user_id, currency, wallet_type, balance, frozen_balance, 
							 wallet_address, private_key_encrypted, name, is_active, is_hidden, external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, wallet.UserID, wallet.Currency, wallet.WalletType,
		wallet.Balance, wallet.FrozenBalance, wallet.WalletAddress,
		wallet.PrivateKeyEncrypted, wallet.Name, wallet.IsActive, wallet.IsHidden, wallet.ExternalID).Scan(
		&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *WalletRepository) InitSchema() error {
	// Create table if not exists (PostgreSQL)
	// SECURITY: private_key_encrypted stores AES-256-GCM encrypted private keys
	// The plaintext key is NEVER stored in the database
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS wallets (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			wallet_type VARCHAR(20) NOT NULL,
			balance DECIMAL(20, 8) DEFAULT 0,
			frozen_balance DECIMAL(20, 8) DEFAULT 0,
			wallet_address VARCHAR(255),
			private_key_encrypted TEXT,
			key_hash VARCHAR(64),
			name VARCHAR(255),
			is_active BOOLEAN DEFAULT TRUE,
			is_hidden BOOLEAN DEFAULT FALSE,
			external_id VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := r.db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create wallets table: %w", err)
	}

	// Migrations: Add missing columns if they don't exist
	// In Postgres, we can do ADD COLUMN IF NOT EXISTS
	migrations := []string{
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS is_hidden BOOLEAN DEFAULT FALSE",
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS external_id VARCHAR(255)",
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS key_hash VARCHAR(64)",
		// Hybrid deposit system fields
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS deposit_memo VARCHAR(64)",
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS sweep_status VARCHAR(20) DEFAULT 'none'",
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS last_swept_at TIMESTAMP",
		"ALTER TABLE wallets ADD COLUMN IF NOT EXISTS pending_sweep_amount DECIMAL(20, 8) DEFAULT 0",
	}

	for _, query := range migrations {
		_, err := r.db.Exec(query)
		if err != nil {
			// Log error but continue? Or fail?
			// Check if error is "column already exists" if on old postgres
			// But for now return error
			return fmt.Errorf("failed to run migration %s: %w", query, err)
		}
	}

	return nil
}

func (r *WalletRepository) GetByExternalID(externalID string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, is_hidden, external_id, created_at, updated_at
		FROM wallets 
		WHERE external_id = $1 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, externalID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive, &wallet.IsHidden, &wallet.ExternalID,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// GetByAddress finds a wallet by its blockchain address
func (r *WalletRepository) GetByAddress(address string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, is_hidden, external_id, created_at, updated_at
		FROM wallets 
		WHERE wallet_address = $1 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, address).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive, &wallet.IsHidden, &wallet.ExternalID,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil to indicate not found (for webhook handling)
		}
		return nil, fmt.Errorf("failed to get wallet by address: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) GetByID(walletID string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, is_hidden, created_at, updated_at
		FROM wallets WHERE id = $1 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, walletID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive, &wallet.IsHidden,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(userID string, includeHidden bool) ([]*models.Wallet, error) {
	var query string
	if includeHidden {
		query = `
			SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
				   wallet_address, name, is_active, is_hidden, created_at, updated_at
			FROM wallets 
			WHERE user_id = $1 AND is_active = true
			ORDER BY created_at DESC
		`
	} else {
		query = `
			SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
				   wallet_address, name, is_active, is_hidden, created_at, updated_at
			FROM wallets 
			WHERE user_id = $1 AND is_active = true AND is_hidden = false
			ORDER BY created_at DESC
		`
	}

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}
	defer rows.Close()

	var wallets []*models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		err := rows.Scan(
			&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
			&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
			&wallet.Name, &wallet.IsActive, &wallet.IsHidden, &wallet.CreatedAt, &wallet.UpdatedAt,
		)
		if err != nil {
			continue
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func (r *WalletRepository) GetByUserAndCurrency(userID, currency string) (*models.Wallet, error) {
	// Query handles both hidden and visible wallets to enforce uniqueness check properly
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, is_hidden, created_at, updated_at
		FROM wallets 
		WHERE user_id = $1 AND currency = $2 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, userID, currency).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive, &wallet.IsHidden,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(walletID string, balance, frozenBalance float64) error {
	query := `
		UPDATE wallets 
		SET balance = $1, frozen_balance = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(query, balance, frozenBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return nil
}

// UpdateAddress updates the wallet address and encrypted private key
// SECURITY: The private key should already be encrypted before calling this method
// keyHash is the SHA-256 hash of the plaintext key for verification purposes
func (r *WalletRepository) UpdateAddress(walletID, address, encryptedPrivateKey, keyHash string) error {
	query := `
		UPDATE wallets 
		SET wallet_address = $1, private_key_encrypted = $2, key_hash = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(query, address, encryptedPrivateKey, keyHash, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet address: %w", err)
	}

	return nil
}

func (r *WalletRepository) Freeze(walletID string) error {
	query := `
		UPDATE wallets 
		SET frozen_balance = balance, balance = 0, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(query, walletID)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet: %w", err)
	}

	return nil
}

func (r *WalletRepository) Unfreeze(walletID string) error {
	query := `
		UPDATE wallets 
		SET balance = frozen_balance, frozen_balance = 0, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(query, walletID)
	if err != nil {
		return fmt.Errorf("failed to unfreeze wallet: %w", err)
	}

	return nil
}

func (r *WalletRepository) Delete(walletID string) error {
	// Hard delete or deactivate
	query := `UPDATE wallets SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, walletID)
	return err
}

func (r *WalletRepository) Hide(walletID string) error {
	// Soft delete / Hide
	query := `UPDATE wallets SET is_hidden = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, walletID)
	return err
}

func (r *WalletRepository) Unhide(walletID string) error {
	// Unhide
	query := `UPDATE wallets SET is_hidden = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, walletID)
	return err
}

// Transaction methods with balance updates
func (r *WalletRepository) UpdateBalanceWithTransaction(walletID string, amount float64, transactionType string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current balance
	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	var newBalance float64
	switch transactionType {
	case "deposit", "receive", "credit":
		newBalance = currentBalance + amount
	case "withdrawal", "send", "debit":
		if currentBalance < amount {
			return fmt.Errorf("insufficient balance")
		}
		newBalance = currentBalance - amount
	default:
		return fmt.Errorf("invalid transaction type: %s", transactionType)
	}

	// Update balance
	_, err = tx.Exec("UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2", newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return tx.Commit()
}

func (r *WalletRepository) GetWalletStats(userID string) (*models.WalletStats, error) {
	// Get total balance across all wallets
	balanceQuery := `
		SELECT COALESCE(SUM(balance), 0) as total_balance, COUNT(*) as wallet_count
		FROM wallets 
		WHERE user_id = $1 AND is_active = true
	`

	var stats models.WalletStats
	var walletCount int
	err := r.db.QueryRow(balanceQuery, userID).Scan(&stats.TotalBalance, &walletCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet stats: %w", err)
	}

	// Get balance by wallet type
	balanceByTypeQuery := `
		SELECT wallet_type, COALESCE(SUM(balance), 0) as total_balance
		FROM wallets 
		WHERE user_id = $1 AND is_active = true
		GROUP BY wallet_type
	`

	rows, err := r.db.Query(balanceByTypeQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance by type: %w", err)
	}
	defer rows.Close()

	stats.BalanceByType = make(map[string]float64)
	for rows.Next() {
		var walletType string
		var balance float64
		if err := rows.Scan(&walletType, &balance); err == nil {
			stats.BalanceByType[walletType] = balance
		}
	}

	// Get total transactions count
	transactionCountQuery := `
		SELECT COUNT(*) 
		FROM transactions t
		JOIN wallets w ON (t.from_wallet_id = w.id OR t.to_wallet_id = w.id)
		WHERE w.user_id = $1
	`

	err = r.db.QueryRow(transactionCountQuery, userID).Scan(&stats.TotalTransactions)
	if err != nil {
		stats.TotalTransactions = 0
	}

	return &stats, nil
}

// GetUserCountry retrieves the user's country code from users table
func (r *WalletRepository) GetUserCountry(userID string) (string, error) {
	query := `SELECT country FROM users WHERE id = $1`
	var country string
	err := r.db.QueryRow(query, userID).Scan(&country)
	if err != nil {
		return "", err
	}
	return country, nil
}

// ==================== Sweep System Methods ====================

// GetWalletsNeedingSweep returns all wallets with pending sweep amounts
func (r *WalletRepository) GetWalletsNeedingSweep() ([]models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, is_hidden,
			   external_id, deposit_memo, sweep_status, last_swept_at, pending_sweep_amount,
			   created_at, updated_at
		FROM wallets 
		WHERE pending_sweep_amount > 0 
		  AND is_active = true 
		  AND wallet_type = 'crypto'
		  AND (sweep_status IS NULL OR sweep_status != 'pending')
		ORDER BY pending_sweep_amount DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets needing sweep: %w", err)
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		var sweepStatus sql.NullString
		var depositMemo sql.NullString
		var lastSweptAt sql.NullTime

		err := rows.Scan(
			&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
			&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
			&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive, &wallet.IsHidden,
			&wallet.ExternalID, &depositMemo, &sweepStatus, &lastSweptAt, &wallet.PendingSweepAmount,
			&wallet.CreatedAt, &wallet.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if sweepStatus.Valid {
			wallet.SweepStatus = sweepStatus.String
		}
		if depositMemo.Valid {
			wallet.DepositMemo = &depositMemo.String
		}
		if lastSweptAt.Valid {
			wallet.LastSweptAt = &lastSweptAt.Time
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// UpdateSweepStatus updates the sweep status and clears the pending amount
func (r *WalletRepository) UpdateSweepStatus(walletID, status string, sweptAmount float64) error {
	query := `
		UPDATE wallets 
		SET sweep_status = $1, 
			pending_sweep_amount = pending_sweep_amount - $2,
			last_swept_at = NOW(),
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, sweptAmount, walletID)
	if err != nil {
		return fmt.Errorf("failed to update sweep status: %w", err)
	}
	return nil
}

// AddPendingSweepAmount adds to the pending sweep amount (after a deposit)
func (r *WalletRepository) AddPendingSweepAmount(walletID string, amount float64) error {
	query := `
		UPDATE wallets 
		SET pending_sweep_amount = pending_sweep_amount + $1,
			sweep_status = 'pending',
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, amount, walletID)
	if err != nil {
		return fmt.Errorf("failed to add pending sweep amount: %w", err)
	}
	return nil
}

// UpdateDepositMemo sets the deposit memo for a wallet
func (r *WalletRepository) UpdateDepositMemo(walletID, memo string) error {
	query := `UPDATE wallets SET deposit_memo = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, memo, walletID)
	return err
}

// GetByDepositMemo finds a wallet by its deposit memo (for XRP/XLM/TON webhooks)
func (r *WalletRepository) GetByDepositMemo(memo, currency string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, name, is_active, is_hidden, deposit_memo, created_at, updated_at
		FROM wallets 
		WHERE deposit_memo = $1 AND currency = $2 AND is_active = true
	`

	var wallet models.Wallet
	var depositMemo sql.NullString

	err := r.db.QueryRow(query, memo, currency).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.Name, &wallet.IsActive, &wallet.IsHidden, &depositMemo,
		&wallet.CreatedAt, &wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet by deposit memo: %w", err)
	}

	if depositMemo.Valid {
		wallet.DepositMemo = &depositMemo.String
	}

	return &wallet, nil
}
