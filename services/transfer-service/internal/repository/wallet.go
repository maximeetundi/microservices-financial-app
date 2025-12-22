package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
)

type WalletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetByID(id string) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `SELECT id, user_id, currency, wallet_type, balance, frozen_balance, is_active 
			  FROM wallets WHERE id = $1`
	
	err := r.db.Get(&wallet, query, id)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(userID string) ([]models.Wallet, error) {
	var wallets []models.Wallet
	query := `SELECT id, user_id, currency, wallet_type, balance, frozen_balance, is_active 
			  FROM wallets WHERE user_id = $1`
	
	err := r.db.Select(&wallets, query, userID)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *WalletRepository) UpdateBalance(walletID string, newBalance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, newBalance, walletID)
	return err
}
