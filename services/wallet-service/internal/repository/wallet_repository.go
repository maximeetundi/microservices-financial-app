package repository

import (
	"database/sql"
	"fmt"
	"encoding/json"
	"time"

	"github.com/crypto-bank/wallet-service/internal/models"
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
							 wallet_address, private_key_encrypted, name, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, wallet.UserID, wallet.Currency, wallet.WalletType,
		wallet.Balance, wallet.FrozenBalance, wallet.WalletAddress,
		wallet.PrivateKeyEncrypted, wallet.Name, wallet.IsActive).Scan(
		&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *WalletRepository) GetByID(walletID string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, created_at, updated_at
		FROM wallets WHERE id = $1 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, walletID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive,
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

func (r *WalletRepository) GetByUserID(userID string) ([]*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, name, is_active, created_at, updated_at
		FROM wallets 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

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
			&wallet.Name, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt,
		)
		if err != nil {
			continue
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func (r *WalletRepository) GetByUserAndCurrency(userID, currency string) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, currency, wallet_type, balance, frozen_balance,
			   wallet_address, private_key_encrypted, name, is_active, created_at, updated_at
		FROM wallets 
		WHERE user_id = $1 AND currency = $2 AND is_active = true
	`

	var wallet models.Wallet
	err := r.db.QueryRow(query, userID, currency).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Currency, &wallet.WalletType,
		&wallet.Balance, &wallet.FrozenBalance, &wallet.WalletAddress,
		&wallet.PrivateKeyEncrypted, &wallet.Name, &wallet.IsActive,
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

func (r *WalletRepository) UpdateAddress(walletID, address, encryptedPrivateKey string) error {
	query := `
		UPDATE wallets 
		SET wallet_address = $1, private_key_encrypted = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(query, address, encryptedPrivateKey, walletID)
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
	query := `UPDATE wallets SET is_active = false, updated_at = NOW() WHERE id = $1`
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
	case "deposit", "receive":
		newBalance = currentBalance + amount
	case "withdrawal", "send":
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