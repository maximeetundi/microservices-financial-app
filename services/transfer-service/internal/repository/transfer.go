package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Create(transfer *models.Transfer) error {
	query := `
		INSERT INTO transactions (id, from_wallet_id, to_wallet_id, transaction_type, amount, fee, currency, status, reference_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query,
		transfer.ID,
		transfer.FromWalletID,
		transfer.ToWalletID,
		transfer.TransferType,
		transfer.Amount,
		transfer.Fee,
		transfer.Currency,
		transfer.Status,
		transfer.Reference,
		transfer.Description,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *TransferRepository) GetByID(id string) (*models.Transfer, error) {
	query := `
		SELECT id, from_wallet_id, to_wallet_id, transaction_type, amount, fee, currency, status, reference_id, description, created_at, updated_at
		FROM transactions WHERE id = $1
	`
	var transfer models.Transfer
	err := r.db.QueryRow(query, id).Scan(
		&transfer.ID,
		&transfer.FromWalletID,
		&transfer.ToWalletID,
		&transfer.TransferType,
		&transfer.Amount,
		&transfer.Fee,
		&transfer.Currency,
		&transfer.Status,
		&transfer.Reference,
		&transfer.Description,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *TransferRepository) GetByUserID(userID string, limit, offset int) ([]*models.Transfer, error) {
	query := `
		SELECT t.id, t.from_wallet_id, t.to_wallet_id, t.transaction_type, t.amount, t.fee, t.currency, t.status, t.reference_id, t.description, t.created_at, t.updated_at
		FROM transactions t
		JOIN wallets w ON (t.from_wallet_id = w.id OR t.to_wallet_id = w.id)
		WHERE w.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*models.Transfer
	for rows.Next() {
		var t models.Transfer
		err := rows.Scan(
			&t.ID, &t.FromWalletID, &t.ToWalletID, &t.TransferType, &t.Amount, &t.Fee,
			&t.Currency, &t.Status, &t.Reference, &t.Description,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, &t)
	}
	return transfers, nil
}

func (r *TransferRepository) UpdateStatus(id, status string) error {
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetByID(id string) (*models.Wallet, error) {
	query := `SELECT id, user_id, currency, wallet_type, balance, frozen_balance, is_active FROM wallets WHERE id = $1`
	var wallet models.Wallet
	err := r.db.QueryRow(query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Currency,
		&wallet.WalletType,
		&wallet.Balance,
		&wallet.FrozenBalance,
		&wallet.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateBalance(id string, amount float64) error {
	query := `UPDATE wallets SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *WalletRepository) GetUserIDByEmail(email string) (string, error) {
	query := `SELECT id FROM users WHERE email = $1`
	var userID string
	err := r.db.QueryRow(query, email).Scan(&userID)
	return userID, err
}

func (r *WalletRepository) GetUserIDByPhone(phone string) (string, error) {
	query := `SELECT id FROM users WHERE phone = $1`
	var userID string
	err := r.db.QueryRow(query, phone).Scan(&userID)
	return userID, err
}

func (r *WalletRepository) GetWalletIDByUserAndCurrency(userID, currency string) (string, error) {
	query := `SELECT id FROM wallets WHERE user_id = $1 AND currency = $2 AND is_active = true LIMIT 1`
	var walletID string
	err := r.db.QueryRow(query, userID, currency).Scan(&walletID)
	return walletID, err
}

func (r *WalletRepository) CreateWallet(id, userID, currency string) error {
	query := `INSERT INTO wallets (id, user_id, currency, wallet_type, balance, frozen_balance, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, 'fiat', 0, 0, true, $4, $4)`
	_, err := r.db.Exec(query, id, userID, currency, time.Now())
	return err
}

// GetUserNameByID retrieves the full name of a user by their ID
func (r *WalletRepository) GetUserNameByID(userID string) (string, error) {
	query := `SELECT COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') FROM users WHERE id = $1`
	var fullName string
	err := r.db.QueryRow(query, userID).Scan(&fullName)
	if err != nil {
		return userID, err // Return ID as fallback
	}
	if fullName == " " || fullName == "" {
		return userID, nil // Return ID if name is empty
	}
	return fullName, nil
}

// GetUserInfo retrieves detailed user info by ID
func (r *WalletRepository) GetUserInfo(userID string) (string, string, string, error) {
	query := `SELECT COALESCE(first_name, '') || ' ' || COALESCE(last_name, ''), email, phone FROM users WHERE id = $1`
	var fullName, email, phone string
	err := r.db.QueryRow(query, userID).Scan(&fullName, &email, &phone)
	if err != nil {
		return "", "", "", err
	}
	// Fallback if name is empty
	if fullName == " " || fullName == "" {
		fullName = "Utilisateur" // Generic name
	}
	return fullName, email, phone, nil
}
