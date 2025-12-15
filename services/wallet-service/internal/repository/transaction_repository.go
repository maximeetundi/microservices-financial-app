package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (from_wallet_id, to_wallet_id, transaction_type, amount, fee, 
								 currency, status, blockchain_tx_hash, reference_id, description, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	var metadataJSON *string
	if transaction.Metadata != nil {
		metadataJSON = transaction.Metadata
	}

	err := r.db.QueryRow(query,
		transaction.FromWalletID, transaction.ToWalletID, transaction.TransactionType,
		transaction.Amount, transaction.Fee, transaction.Currency, transaction.Status,
		transaction.BlockchainTxHash, transaction.ReferenceID, transaction.Description,
		metadataJSON).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetByID(transactionID string) (*models.Transaction, error) {
	query := `
		SELECT id, from_wallet_id, to_wallet_id, transaction_type, amount, fee, currency,
			   status, blockchain_tx_hash, reference_id, description, metadata, created_at, updated_at
		FROM transactions WHERE id = $1
	`

	var transaction models.Transaction
	err := r.db.QueryRow(query, transactionID).Scan(
		&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID,
		&transaction.TransactionType, &transaction.Amount, &transaction.Fee,
		&transaction.Currency, &transaction.Status, &transaction.BlockchainTxHash,
		&transaction.ReferenceID, &transaction.Description, &transaction.Metadata,
		&transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetByWalletID(walletID string, limit, offset int, status, txType string) ([]*models.Transaction, error) {
	query := `
		SELECT t.id, t.from_wallet_id, t.to_wallet_id, t.transaction_type, t.amount, t.fee,
			   t.currency, t.status, t.blockchain_tx_hash, t.reference_id, t.description,
			   t.metadata, t.created_at, t.updated_at
		FROM transactions t
		WHERE (t.from_wallet_id = $1 OR t.to_wallet_id = $1)
	`

	args := []interface{}{walletID}
	argCount := 1

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, status)
	}

	if txType != "" {
		argCount++
		query += fmt.Sprintf(" AND t.transaction_type = $%d", argCount)
		args = append(args, txType)
	}

	query += " ORDER BY t.created_at DESC"

	if limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)

		if offset > 0 {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID,
			&transaction.TransactionType, &transaction.Amount, &transaction.Fee,
			&transaction.Currency, &transaction.Status, &transaction.BlockchainTxHash,
			&transaction.ReferenceID, &transaction.Description, &transaction.Metadata,
			&transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByUserID(userID string, limit, offset int, currency, status string) ([]*models.Transaction, error) {
	query := `
		SELECT t.id, t.from_wallet_id, t.to_wallet_id, t.transaction_type, t.amount, t.fee,
			   t.currency, t.status, t.blockchain_tx_hash, t.reference_id, t.description,
			   t.metadata, t.created_at, t.updated_at
		FROM transactions t
		LEFT JOIN wallets w1 ON t.from_wallet_id = w1.id
		LEFT JOIN wallets w2 ON t.to_wallet_id = w2.id
		WHERE (w1.user_id = $1 OR w2.user_id = $1)
	`

	args := []interface{}{userID}
	argCount := 1

	if currency != "" {
		argCount++
		query += fmt.Sprintf(" AND t.currency = $%d", argCount)
		args = append(args, currency)
	}

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, status)
	}

	query += " ORDER BY t.created_at DESC"

	if limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)

		if offset > 0 {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID,
			&transaction.TransactionType, &transaction.Amount, &transaction.Fee,
			&transaction.Currency, &transaction.Status, &transaction.BlockchainTxHash,
			&transaction.ReferenceID, &transaction.Description, &transaction.Metadata,
			&transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateStatus(transactionID, status string, blockchainTxHash *string) error {
	query := `
		UPDATE transactions 
		SET status = $1, blockchain_tx_hash = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, blockchainTxHash, transactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetPendingCryptoTransactions(currency string) ([]*models.Transaction, error) {
	query := `
		SELECT t.id, t.from_wallet_id, t.to_wallet_id, t.transaction_type, t.amount, t.fee,
			   t.currency, t.status, t.blockchain_tx_hash, t.reference_id, t.description,
			   t.metadata, t.created_at, t.updated_at
		FROM transactions t
		WHERE t.currency = $1 AND t.status = 'pending' AND t.blockchain_tx_hash IS NOT NULL
		ORDER BY t.created_at ASC
	`

	rows, err := r.db.Query(query, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending crypto transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID,
			&transaction.TransactionType, &transaction.Amount, &transaction.Fee,
			&transaction.Currency, &transaction.Status, &transaction.BlockchainTxHash,
			&transaction.ReferenceID, &transaction.Description, &transaction.Metadata,
			&transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByBlockchainHash(txHash string) (*models.Transaction, error) {
	query := `
		SELECT id, from_wallet_id, to_wallet_id, transaction_type, amount, fee, currency,
			   status, blockchain_tx_hash, reference_id, description, metadata, created_at, updated_at
		FROM transactions WHERE blockchain_tx_hash = $1
	`

	var transaction models.Transaction
	err := r.db.QueryRow(query, txHash).Scan(
		&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID,
		&transaction.TransactionType, &transaction.Amount, &transaction.Fee,
		&transaction.Currency, &transaction.Status, &transaction.BlockchainTxHash,
		&transaction.ReferenceID, &transaction.Description, &transaction.Metadata,
		&transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetTransactionsSummary(userID string, fromDate, toDate time.Time) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_transactions,
			COALESCE(SUM(CASE WHEN t.transaction_type IN ('deposit', 'receive') THEN t.amount ELSE 0 END), 0) as total_received,
			COALESCE(SUM(CASE WHEN t.transaction_type IN ('withdrawal', 'send') THEN t.amount ELSE 0 END), 0) as total_sent,
			COALESCE(SUM(t.fee), 0) as total_fees
		FROM transactions t
		LEFT JOIN wallets w1 ON t.from_wallet_id = w1.id
		LEFT JOIN wallets w2 ON t.to_wallet_id = w2.id
		WHERE (w1.user_id = $1 OR w2.user_id = $1)
		  AND t.status = 'completed'
		  AND t.created_at >= $2 
		  AND t.created_at <= $3
	`

	var summary map[string]interface{} = make(map[string]interface{})
	var totalTx int
	var totalReceived, totalSent, totalFees float64

	err := r.db.QueryRow(query, userID, fromDate, toDate).Scan(
		&totalTx, &totalReceived, &totalSent, &totalFees,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary: %w", err)
	}

	summary["total_transactions"] = totalTx
	summary["total_received"] = totalReceived
	summary["total_sent"] = totalSent
	summary["total_fees"] = totalFees
	summary["net_amount"] = totalReceived - totalSent - totalFees

	return summary, nil
}

func (r *TransactionRepository) CreateWithBalanceUpdate(transaction *models.Transaction, fromWalletUpdate, toWalletUpdate *models.Wallet) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create transaction
	insertQuery := `
		INSERT INTO transactions (from_wallet_id, to_wallet_id, transaction_type, amount, fee, 
								 currency, status, blockchain_tx_hash, reference_id, description, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(insertQuery,
		transaction.FromWalletID, transaction.ToWalletID, transaction.TransactionType,
		transaction.Amount, transaction.Fee, transaction.Currency, transaction.Status,
		transaction.BlockchainTxHash, transaction.ReferenceID, transaction.Description,
		transaction.Metadata).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update from wallet balance if provided
	if fromWalletUpdate != nil {
		_, err = tx.Exec("UPDATE wallets SET balance = $1, frozen_balance = $2, updated_at = NOW() WHERE id = $3",
			fromWalletUpdate.Balance, fromWalletUpdate.FrozenBalance, fromWalletUpdate.ID)
		if err != nil {
			return fmt.Errorf("failed to update from wallet balance: %w", err)
		}
	}

	// Update to wallet balance if provided
	if toWalletUpdate != nil {
		_, err = tx.Exec("UPDATE wallets SET balance = $1, frozen_balance = $2, updated_at = NOW() WHERE id = $3",
			toWalletUpdate.Balance, toWalletUpdate.FrozenBalance, toWalletUpdate.ID)
		if err != nil {
			return fmt.Errorf("failed to update to wallet balance: %w", err)
		}
	}

	return tx.Commit()
}