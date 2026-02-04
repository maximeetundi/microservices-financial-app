package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DepositFundMovementService handles fund movements for deposits
// Moves money from hot wallets to user accounts
type DepositFundMovementService struct {
	db *sql.DB
}

func NewDepositFundMovementService(db *sql.DB) *DepositFundMovementService {
	return &DepositFundMovementService{db: db}
}

// ProcessDepositFromWallet credits user from hot wallet when external deposit is confirmed
// This implements double-entry bookkeeping: Debit Hot Wallet -> Credit User Account
func (s *DepositFundMovementService) ProcessDepositFromWallet(
	ctx context.Context,
	userID string,
	hotWalletID string,
	amount float64,
	currency string,
	referenceID string,
	providerName string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get hot wallet balance
	var hotWalletBalance float64
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM platform_accounts
		WHERE id = $1 AND account_type = 'hot'
		FOR UPDATE
	`, hotWalletID).Scan(&hotWalletBalance)

	if err != nil {
		return fmt.Errorf("get hot wallet balance: %w", err)
	}

	// 2. Check if hot wallet has enough balance
	if hotWalletBalance < amount {
		return fmt.Errorf("insufficient hot wallet balance: has %f, needs %f", hotWalletBalance, amount)
	}

	// 3. Debit hot wallet
	_, err = tx.ExecContext(ctx, `
		UPDATE platform_accounts
		SET balance = balance - $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, amount, hotWalletID)

	if err != nil {
		return fmt.Errorf("debit hot wallet: %w", err)
	}

	// 4. Get user wallet ID
	var userWalletID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM wallets
		WHERE user_id = $1 AND currency = $2
	`, userID, currency).Scan(&userWalletID)

	if err != nil {
		return fmt.Errorf("get user wallet: %w", err)
	}

	// 5. Credit user wallet
	_, err = tx.ExecContext(ctx, `
		UPDATE wallets
		SET balance = balance + $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, amount, userWalletID)

	if err != nil {
		return fmt.Errorf("credit user wallet: %w", err)
	}

	// 6. Create transaction record
	transactionID := uuid.New().String()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (
			id, user_id, wallet_id, type, amount, currency, status,
			description, reference_id, created_at
		) VALUES ($1, $2, $3, 'deposit', $4, $5, 'completed', $6, $7, CURRENT_TIMESTAMP)
	`, transactionID, userID, userWalletID, amount, currency,
		fmt.Sprintf("Deposit via %s", providerName), referenceID)

	if err != nil {
		return fmt.Errorf("create transaction record: %w", err)
	}

	// 7. Log the movement in platform_account_movements
	_, err = tx.ExecContext(ctx, `
		INSERT INTO platform_account_movements (
			account_id, movement_type, amount, currency, balance_after,
			related_user_id, related_transaction_id, description, created_at
		) VALUES ($1, 'debit', $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
	`, hotWalletID, amount, currency, hotWalletBalance-amount, userID, transactionID,
	`, hotWalletID, amount, currency, hotWalletBalanceAfter, userID, transactionID,
		fmt.Sprintf("Deposit to user via %s", providerName))

	if err != nil {
		return fmt.Errorf("log hot wallet movement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Success logging with details
	fmt.Printf("════════════════════════════════════════════════════════════\n")
	fmt.Printf("[DEPOSIT COMPLETED] Transaction: %s\n", transactionID)
	fmt.Printf("  Provider: %s\n", providerName)
	fmt.Printf("  Amount: %.2f %s\n", amount, currency)
	fmt.Printf("  User: %s\n", userID)
	fmt.Printf("  Hot Wallet: %s\n", hotWalletID)
	fmt.Printf("  Flow: Provider → Hot Wallet(%.2f) → User Account(+%.2f)\n", -amount, amount)
	fmt.Printf("  Hot Wallet Balance After: %.2f %s\n", hotWalletBalanceAfter, currency)
	fmt.Printf("════════════════════════════════════════════════════════════\n")

	return nil
}

// SelectBestWalletForDeposit selects the best hot wallet for a deposit
// Uses the DB function select_best_wallet_for_instance
func (s *DepositFundMovementService) SelectBestWalletForDeposit(
	ctx context.Context,
	instanceID string,
	currency string,
	amount float64,
) (string, error) {
	var walletID sql.NullString

	err := s.db.QueryRowContext(ctx, `
		SELECT select_best_wallet_for_instance($1, $2, $3)
	`, instanceID, currency, amount).Scan(&walletID)

	if err != nil {
		return "", fmt.Errorf("select best wallet: %w", err)
	}

	if !walletID.Valid || walletID.String == "" {
		return "", fmt.Errorf("no available wallet for instance %s with currency %s", instanceID, currency)
	}

	return walletID.String, nil
}

// CheckAndTriggerRecharge checks if a wallet needs recharge and triggers it
func (s *DepositFundMovementService) CheckAndTriggerRecharge(
	ctx context.Context,
	instanceWalletID string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get wallet info
	var (
		hotWalletID            string
		autoRechargeEnabled    bool
		rechargeThreshold      sql.NullFloat64
		rechargeTarget         sql.NullFloat64
		rechargeSourceWalletID sql.NullString
		walletBalance          float64
		currency               string
	)

	err = tx.QueryRowContext(ctx, `
		SELECT 
			aiw.hot_wallet_id,
			aiw.auto_recharge_enabled,
			aiw.recharge_threshold,
			aiw.recharge_target,
			aiw.recharge_source_wallet_id,
			pa.balance,
			pa.currency
		FROM aggregator_instance_wallets aiw
		JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id
		WHERE aiw.id = $1
		FOR UPDATE
	`, instanceWalletID).Scan(
		&hotWalletID, &autoRechargeEnabled, &rechargeThreshold,
		&rechargeTarget, &rechargeSourceWalletID, &walletBalance, &currency,
	)

	if err != nil {
		return fmt.Errorf("get wallet info: %w", err)
	}

	// Check if recharge is needed
	if !autoRechargeEnabled || !rechargeThreshold.Valid || !rechargeTarget.Valid {
		return nil // Recharge not configured
	}

	if walletBalance >= rechargeThreshold.Float64 {
		return nil // Balance is sufficient
	}

	if !rechargeSourceWalletID.Valid {
		return fmt.Errorf("recharge needed but no source wallet configured")
	}

	// Calculate recharge amount
	rechargeAmount := rechargeTarget.Float64 - walletBalance

	// Debit source wallet (usually cold wallet)
	_, err = tx.ExecContext(ctx, `
		UPDATE platform_accounts
		SET balance = balance - $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND balance >= $1
	`, rechargeAmount, rechargeSourceWalletID.String)

	if err != nil {
		return fmt.Errorf("debit source wallet: %w", err)
	}

	// Credit hot wallet
	_, err = tx.ExecContext(ctx, `
		UPDATE platform_accounts
		SET balance = balance + $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, rechargeAmount, hotWalletID)

	if err != nil {
		return fmt.Errorf("credit hot wallet: %w", err)
	}

	// Log the recharge movement
	_, err = tx.ExecContext(ctx, `
		INSERT INTO platform_account_movements (
			account_id, movement_type, amount, currency, balance_after,
			description, created_at
		) VALUES 
		($1, 'credit', $2, $3, $4, 'Auto-recharge from cold wallet', CURRENT_TIMESTAMP),
		($5, 'debit', $2, $3, (SELECT balance FROM platform_accounts WHERE id = $5), 'Auto-recharge to hot wallet', CURRENT_TIMESTAMP)
	`, hotWalletID, rechargeAmount, currency, walletBalance+rechargeAmount, rechargeSourceWalletID.String)

	if err != nil {
		return fmt.Errorf("log recharge movements: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit recharge: %w", err)
	}

	return nil
}

// UpdateWalletUsageStats updates statistics for a wallet after a transaction
func (s *DepositFundMovementService) UpdateWalletUsageStats(
	ctx context.Context,
	instanceWalletID string,
	amount float64,
	transactionType string,
) error {
	query := `
		UPDATE aggregator_instance_wallets
		SET transaction_count = transaction_count + 1,
		    last_used_at = $1
	`

	if transactionType == "deposit" {
		query += `, total_deposits = total_deposits + $2`
	} else if transactionType == "withdrawal" {
		query += `, total_withdrawals = total_withdrawals + $2`
	}

	query += ` WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, time.Now(), amount, instanceWalletID)
	return err
}
