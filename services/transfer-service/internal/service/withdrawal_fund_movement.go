package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// WithdrawalFundMovementService handles fund movements for withdrawals
type WithdrawalFundMovementService struct {
	db *sql.DB
}

// NewWithdrawalFundMovementService creates a new withdrawal fund movement service
func NewWithdrawalFundMovementService(db *sql.DB) *WithdrawalFundMovementService {
	return &WithdrawalFundMovementService{db: db}
}

// ProcessWithdrawal moves funds from user wallet to hot wallet before external payout
func (s *WithdrawalFundMovementService) ProcessWithdrawal(
	ctx context.Context,
	userID string,
	hotWalletID string,
	amount float64,
	currency string,
	transactionRef string,
	providerCode string,
	destinationAccount string, // Bank account or mobile money number
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get user wallet ID
	var userWalletID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM wallets 
		WHERE user_id = $1 AND currency = $2 AND is_active = true
		LIMIT 1
	`, userID, currency).Scan(&userWalletID)

	if err != nil {
		return fmt.Errorf("failed to find user wallet: %w", err)
	}

	// 2. Check user balance
	var userBalance float64
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE id = $1
	`, userWalletID).Scan(&userBalance)

	if err != nil {
		return fmt.Errorf("failed to get user balance: %w", err)
	}

	if userBalance < amount {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", userBalance, amount)
	}

	// 3. DEBIT User Wallet (user pays)
	_, err = tx.ExecContext(ctx, `
		UPDATE wallets 
		SET balance = balance - $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, amount, userWalletID)

	if err != nil {
		return fmt.Errorf("failed to debit user wallet: %w", err)
	}

	// 4. CREDIT Hot Wallet (platform receives for payout)
	var hotWalletBalanceAfter float64
	err = tx.QueryRowContext(ctx, `
		UPDATE platform_accounts 
		SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING balance
	`, amount, hotWalletID).Scan(&hotWalletBalanceAfter)

	if err != nil {
		return fmt.Errorf("failed to credit hot wallet: %w", err)
	}

	// 5. Create transaction record
	transactionID := uuid.New().String()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (
			id, user_id, wallet_id, type, amount, currency,
			status, reference, provider, destination_account,
			created_at, updated_at
		) VALUES ($1, $2, $3, 'withdrawal', $4, $5, 'pending', $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, transactionID, userID, userWalletID, amount, currency, transactionRef, providerCode, destinationAccount)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// 6. Log user wallet movement (debit)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallet_movements (
			id, wallet_id, transaction_id, movement_type, amount, currency,
			balance_after, description, created_at
		)
		SELECT $1, id, $2, 'debit', $3, $4, balance, $5, CURRENT_TIMESTAMP
		FROM wallets WHERE id = $6
	`, uuid.New().String(), transactionID, amount, currency,
		fmt.Sprintf("Withdrawal to %s via %s", destinationAccount, providerCode), userWalletID)

	if err != nil {
		return fmt.Errorf("failed to log user movement: %w", err)
	}

	// 7. Log platform account movement (credit)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO platform_account_movements (
			id, account_id, movement_type, amount, currency, balance_after,
			related_user_id, related_transaction_id, description, created_at
		) VALUES ($1, $2, 'credit', $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
	`, uuid.New().String(), hotWalletID, amount, currency, hotWalletBalanceAfter,
		userID, transactionID, fmt.Sprintf("Withdrawal staging for %s", destinationAccount))

	if err != nil {
		return fmt.Errorf("failed to log platform movement: %w", err)
	}

	// 8. Update instance wallet statistics
	_, err = tx.ExecContext(ctx, `
		UPDATE aggregator_instance_wallets
		SET total_withdrawals = total_withdrawals + $1,
		    transaction_count = transaction_count + 1,
		    last_used_at = CURRENT_TIMESTAMP
		WHERE hot_wallet_id = $2
	`, amount, hotWalletID)

	if err != nil {
		// Non-critical, just log
		fmt.Printf("Warning: failed to update instance wallet stats: %v\n", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("[WITHDRAWAL] User %s → Hot Wallet: %.2f %s (Ref: %s)\n",
		userID, amount, currency, transactionRef)

	return nil
}

// ConfirmWithdrawalPayout updates transaction status after provider confirms payout
func (s *WithdrawalFundMovementService) ConfirmWithdrawalPayout(
	ctx context.Context,
	transactionRef string,
	status string, // "successful", "failed"
	providerReference string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if status == "successful" {
		// Update transaction status
		_, err = tx.ExecContext(ctx, `
			UPDATE transactions
			SET status = 'completed',
			    provider_reference = $1,
			    completed_at = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP
			WHERE reference = $2
		`, providerReference, transactionRef)

		if err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		fmt.Printf("[WITHDRAWAL CONFIRMED] Ref: %s - Provider Ref: %s\n",
			transactionRef, providerReference)

	} else if status == "failed" {
		// Withdrawal failed - REVERSE the movement
		var userID, walletID, hotWalletID, currency string
		var amount float64

		err = tx.QueryRowContext(ctx, `
			SELECT t.user_id, t.wallet_id, t.amount, t.currency
			FROM transactions t
			WHERE t.reference = $1
		`, transactionRef).Scan(&userID, &walletID, &amount, &currency)

		if err != nil {
			return fmt.Errorf("failed to find transaction: %w", err)
		}

		// Find hot wallet used
		err = tx.QueryRowContext(ctx, `
			SELECT account_id FROM platform_account_movements
			WHERE related_transaction_id = (SELECT id FROM transactions WHERE reference = $1)
			AND movement_type = 'credit'
			LIMIT 1
		`, transactionRef).Scan(&hotWalletID)

		if err != nil {
			return fmt.Errorf("failed to find hot wallet: %w", err)
		}

		// REVERSE: Credit user back
		_, err = tx.ExecContext(ctx, `
			UPDATE wallets 
			SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, amount, walletID)

		if err != nil {
			return fmt.Errorf("failed to refund user: %w", err)
		}

		// REVERSE: Debit hot wallet back
		_, err = tx.ExecContext(ctx, `
			UPDATE platform_accounts 
			SET balance = balance - $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, amount, hotWalletID)

		if err != nil {
			return fmt.Errorf("failed to reverse hot wallet: %w", err)
		}

		// Update transaction status
		_, err = tx.ExecContext(ctx, `
			UPDATE transactions
			SET status = 'failed',
			    updated_at = CURRENT_TIMESTAMP
			WHERE reference = $1
		`, transactionRef)

		// Log reversal movements
		transactionID := uuid.New().String()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO wallet_movements (
				id, wallet_id, transaction_id, movement_type, amount, currency,
				balance_after, description, created_at
			)
			SELECT $1, id, $2, 'credit', $3, $4, balance, 'Withdrawal failed - Refund', CURRENT_TIMESTAMP
			FROM wallets WHERE id = $5
		`, uuid.New().String(), transactionID, amount, currency, walletID)

		fmt.Printf("[WITHDRAWAL FAILED - REVERSED] Ref: %s - Amount: %.2f %s\n",
			transactionRef, amount, currency)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SelectBestWalletForWithdrawal selects the best hot wallet for withdrawal payout
func (s *WithdrawalFundMovementService) SelectBestWalletForWithdrawal(
	ctx context.Context,
	instanceID string,
	currency string,
	amount float64,
) (string, error) {
	var walletID string

	query := `
		SELECT aiw.hot_wallet_id
		FROM select_best_wallet_for_instance($1, $2, $3) aiw
		LIMIT 1
	`

	err := s.db.QueryRowContext(ctx, query, instanceID, currency, amount).Scan(&walletID)
	if err != nil {
		return "", fmt.Errorf("no available withdrawal wallet: %w", err)
	}

	return walletID, nil
}
