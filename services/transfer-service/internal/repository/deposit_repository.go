package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// DepositTransaction represents a deposit transaction record
type DepositTransaction struct {
	ID                   string          `json:"id" db:"id"`
	UserID               string          `json:"user_id" db:"user_id"`
	Amount               float64         `json:"amount" db:"amount"`
	Currency             string          `json:"currency" db:"currency"`
	Fee                  float64         `json:"fee" db:"fee"`
	NetAmount            *float64        `json:"net_amount,omitempty" db:"net_amount"`
	ProviderCode         string          `json:"provider_code" db:"provider_code"`
	AggregatorInstanceID *string         `json:"aggregator_instance_id,omitempty" db:"aggregator_instance_id"`
	HotWalletID          *string         `json:"hot_wallet_id,omitempty" db:"hot_wallet_id"`
	PaymentURL           *string         `json:"payment_url,omitempty" db:"payment_url"`
	ProviderReference    *string         `json:"provider_reference,omitempty" db:"provider_reference"`
	UserEmail            *string         `json:"user_email,omitempty" db:"user_email"`
	UserPhone            *string         `json:"user_phone,omitempty" db:"user_phone"`
	UserCountry          *string         `json:"user_country,omitempty" db:"user_country"`
	Status               string          `json:"status" db:"status"`
	StatusMessage        *string         `json:"status_message,omitempty" db:"status_message"`
	WebhookReceivedAt    *time.Time      `json:"webhook_received_at,omitempty" db:"webhook_received_at"`
	WebhookData          json.RawMessage `json:"webhook_data,omitempty" db:"webhook_data"`
	WebhookVerified      bool            `json:"webhook_verified" db:"webhook_verified"`
	RetryCount           int             `json:"retry_count" db:"retry_count"`
	MaxRetries           int             `json:"max_retries" db:"max_retries"`
	ExpiresAt            *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	CompletedAt          *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CancelledAt          *time.Time      `json:"cancelled_at,omitempty" db:"cancelled_at"`
	FailedAt             *time.Time      `json:"failed_at,omitempty" db:"failed_at"`
	FailureReason        *string         `json:"failure_reason,omitempty" db:"failure_reason"`
	UserWalletID         *string         `json:"user_wallet_id,omitempty" db:"user_wallet_id"`
	WalletCredited       bool            `json:"wallet_credited" db:"wallet_credited"`
	WalletCreditedAt     *time.Time      `json:"wallet_credited_at,omitempty" db:"wallet_credited_at"`
	ReturnURL            *string         `json:"return_url,omitempty" db:"return_url"`
	CancelURL            *string         `json:"cancel_url,omitempty" db:"cancel_url"`
	Metadata             json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	IPAddress            *string         `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent            *string         `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// Deposit statuses
const (
	DepositStatusPending    = "pending"
	DepositStatusProcessing = "processing"
	DepositStatusCompleted  = "completed"
	DepositStatusFailed     = "failed"
	DepositStatusCancelled  = "cancelled"
	DepositStatusExpired    = "expired"
)

// CreateDepositRequest represents the data needed to create a deposit transaction
type CreateDepositRequest struct {
	ID                   string
	UserID               string
	Amount               float64
	Currency             string
	Fee                  float64
	ProviderCode         string
	AggregatorInstanceID string
	HotWalletID          string
	PaymentURL           string
	ProviderReference    string
	UserEmail            string
	UserPhone            string
	UserCountry          string
	ReturnURL            string
	CancelURL            string
	ExpiresAt            time.Time
	Metadata             map[string]interface{}
	IPAddress            string
	UserAgent            string
}

// DepositRepository handles deposit transaction database operations
type DepositRepository struct {
	db *sql.DB
}

// NewDepositRepository creates a new deposit repository
func NewDepositRepository(db *sql.DB) *DepositRepository {
	return &DepositRepository{db: db}
}

// Create creates a new deposit transaction
func (r *DepositRepository) Create(ctx context.Context, req *CreateDepositRequest) (*DepositTransaction, error) {
	metadataJSON, err := json.Marshal(req.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO deposit_transactions (
			id, user_id, amount, currency, fee, net_amount,
			provider_code, aggregator_instance_id, hot_wallet_id,
			payment_url, provider_reference,
			user_email, user_phone, user_country,
			status, return_url, cancel_url, expires_at,
			metadata, ip_address, user_agent
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11,
			$12, $13, $14,
			$15, $16, $17, $18,
			$19, $20, $21
		)
		RETURNING created_at, updated_at
	`

	netAmount := req.Amount - req.Fee

	var createdAt, updatedAt time.Time
	err = r.db.QueryRowContext(ctx, query,
		req.ID, req.UserID, req.Amount, req.Currency, req.Fee, netAmount,
		req.ProviderCode, nullString(req.AggregatorInstanceID), nullString(req.HotWalletID),
		nullString(req.PaymentURL), nullString(req.ProviderReference),
		nullString(req.UserEmail), nullString(req.UserPhone), nullString(req.UserCountry),
		DepositStatusPending, nullString(req.ReturnURL), nullString(req.CancelURL), req.ExpiresAt,
		metadataJSON, nullString(req.IPAddress), nullString(req.UserAgent),
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create deposit transaction: %w", err)
	}

	return &DepositTransaction{
		ID:                   req.ID,
		UserID:               req.UserID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Fee:                  req.Fee,
		NetAmount:            &netAmount,
		ProviderCode:         req.ProviderCode,
		AggregatorInstanceID: strPtr(req.AggregatorInstanceID),
		HotWalletID:          strPtr(req.HotWalletID),
		PaymentURL:           strPtr(req.PaymentURL),
		ProviderReference:    strPtr(req.ProviderReference),
		UserEmail:            strPtr(req.UserEmail),
		UserPhone:            strPtr(req.UserPhone),
		UserCountry:          strPtr(req.UserCountry),
		Status:               DepositStatusPending,
		ExpiresAt:            &req.ExpiresAt,
		ReturnURL:            strPtr(req.ReturnURL),
		CancelURL:            strPtr(req.CancelURL),
		Metadata:             metadataJSON,
		IPAddress:            strPtr(req.IPAddress),
		UserAgent:            strPtr(req.UserAgent),
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}, nil
}

// GetByID retrieves a deposit transaction by ID
func (r *DepositRepository) GetByID(ctx context.Context, id string) (*DepositTransaction, error) {
	query := `
		SELECT id, user_id, amount, currency, fee, net_amount,
			   provider_code, aggregator_instance_id, hot_wallet_id,
			   payment_url, provider_reference,
			   user_email, user_phone, user_country,
			   status, status_message,
			   webhook_received_at, webhook_data, webhook_verified,
			   retry_count, max_retries, expires_at,
			   completed_at, cancelled_at, failed_at, failure_reason,
			   user_wallet_id, wallet_credited, wallet_credited_at,
			   return_url, cancel_url, metadata, ip_address, user_agent,
			   created_at, updated_at
		FROM deposit_transactions
		WHERE id = $1
	`

	var tx DepositTransaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.UserID, &tx.Amount, &tx.Currency, &tx.Fee, &tx.NetAmount,
		&tx.ProviderCode, &tx.AggregatorInstanceID, &tx.HotWalletID,
		&tx.PaymentURL, &tx.ProviderReference,
		&tx.UserEmail, &tx.UserPhone, &tx.UserCountry,
		&tx.Status, &tx.StatusMessage,
		&tx.WebhookReceivedAt, &tx.WebhookData, &tx.WebhookVerified,
		&tx.RetryCount, &tx.MaxRetries, &tx.ExpiresAt,
		&tx.CompletedAt, &tx.CancelledAt, &tx.FailedAt, &tx.FailureReason,
		&tx.UserWalletID, &tx.WalletCredited, &tx.WalletCreditedAt,
		&tx.ReturnURL, &tx.CancelURL, &tx.Metadata, &tx.IPAddress, &tx.UserAgent,
		&tx.CreatedAt, &tx.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit transaction: %w", err)
	}

	return &tx, nil
}

// GetByProviderReference retrieves a deposit transaction by provider reference
func (r *DepositRepository) GetByProviderReference(ctx context.Context, providerCode, providerRef string) (*DepositTransaction, error) {
	query := `
		SELECT id, user_id, amount, currency, fee, net_amount,
			   provider_code, aggregator_instance_id, hot_wallet_id,
			   payment_url, provider_reference,
			   user_email, user_phone, user_country,
			   status, status_message,
			   webhook_received_at, webhook_data, webhook_verified,
			   retry_count, max_retries, expires_at,
			   completed_at, cancelled_at, failed_at, failure_reason,
			   user_wallet_id, wallet_credited, wallet_credited_at,
			   return_url, cancel_url, metadata, ip_address, user_agent,
			   created_at, updated_at
		FROM deposit_transactions
		WHERE provider_code = $1 AND provider_reference = $2
	`

	var tx DepositTransaction
	err := r.db.QueryRowContext(ctx, query, providerCode, providerRef).Scan(
		&tx.ID, &tx.UserID, &tx.Amount, &tx.Currency, &tx.Fee, &tx.NetAmount,
		&tx.ProviderCode, &tx.AggregatorInstanceID, &tx.HotWalletID,
		&tx.PaymentURL, &tx.ProviderReference,
		&tx.UserEmail, &tx.UserPhone, &tx.UserCountry,
		&tx.Status, &tx.StatusMessage,
		&tx.WebhookReceivedAt, &tx.WebhookData, &tx.WebhookVerified,
		&tx.RetryCount, &tx.MaxRetries, &tx.ExpiresAt,
		&tx.CompletedAt, &tx.CancelledAt, &tx.FailedAt, &tx.FailureReason,
		&tx.UserWalletID, &tx.WalletCredited, &tx.WalletCreditedAt,
		&tx.ReturnURL, &tx.CancelURL, &tx.Metadata, &tx.IPAddress, &tx.UserAgent,
		&tx.CreatedAt, &tx.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit by provider ref: %w", err)
	}

	return &tx, nil
}

// UpdateStatus updates the status of a deposit transaction
func (r *DepositRepository) UpdateStatus(ctx context.Context, id, status, message string) error {
	var completedAt, failedAt, cancelledAt *time.Time
	now := time.Now()

	switch status {
	case DepositStatusCompleted:
		completedAt = &now
	case DepositStatusFailed, DepositStatusExpired:
		failedAt = &now
	case DepositStatusCancelled:
		cancelledAt = &now
	}

	query := `
		UPDATE deposit_transactions
		SET status = $2,
			status_message = $3,
			completed_at = COALESCE($4, completed_at),
			failed_at = COALESCE($5, failed_at),
			cancelled_at = COALESCE($6, cancelled_at)
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status, message, completedAt, failedAt, cancelledAt)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("deposit transaction not found: %s", id)
	}

	return nil
}

// UpdateProviderReference updates the provider reference after initiation
func (r *DepositRepository) UpdateProviderReference(ctx context.Context, id, providerRef, paymentURL string) error {
	query := `
		UPDATE deposit_transactions
		SET provider_reference = $2, payment_url = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, providerRef, paymentURL)
	return err
}

// SaveWebhookData saves the webhook payload and marks it as received
func (r *DepositRepository) SaveWebhookData(ctx context.Context, id string, webhookData interface{}, verified bool) error {
	dataJSON, err := json.Marshal(webhookData)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook data: %w", err)
	}

	query := `
		UPDATE deposit_transactions
		SET webhook_received_at = CURRENT_TIMESTAMP,
			webhook_data = $2,
			webhook_verified = $3
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, query, id, dataJSON, verified)
	return err
}

// MarkWalletCredited marks the user's wallet as credited
func (r *DepositRepository) MarkWalletCredited(ctx context.Context, id, userWalletID string) error {
	query := `
		UPDATE deposit_transactions
		SET user_wallet_id = $2,
			wallet_credited = true,
			wallet_credited_at = CURRENT_TIMESTAMP,
			status = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, userWalletID, DepositStatusCompleted)
	return err
}

// MarkFailed marks a transaction as failed
func (r *DepositRepository) MarkFailed(ctx context.Context, id, reason string) error {
	query := `
		UPDATE deposit_transactions
		SET status = $2,
			failure_reason = $3,
			failed_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, DepositStatusFailed, reason)
	return err
}

// MarkCancelled marks a transaction as cancelled by user
func (r *DepositRepository) MarkCancelled(ctx context.Context, id string) error {
	query := `
		UPDATE deposit_transactions
		SET status = $2,
			cancelled_at = CURRENT_TIMESTAMP,
			status_message = 'Cancelled by user'
		WHERE id = $1 AND status = $3
	`

	result, err := r.db.ExecContext(ctx, query, id, DepositStatusCancelled, DepositStatusPending)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction not found or cannot be cancelled")
	}

	return nil
}

// GetPendingExpired returns all pending transactions that have expired
func (r *DepositRepository) GetPendingExpired(ctx context.Context) ([]*DepositTransaction, error) {
	query := `
		SELECT id, user_id, amount, currency, provider_code, status, expires_at, created_at
		FROM deposit_transactions
		WHERE status = $1
		  AND expires_at IS NOT NULL
		  AND expires_at < CURRENT_TIMESTAMP
		ORDER BY expires_at ASC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, DepositStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*DepositTransaction
	for rows.Next() {
		var tx DepositTransaction
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.Amount, &tx.Currency,
			&tx.ProviderCode, &tx.Status, &tx.ExpiresAt, &tx.CreatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// ExpirePendingTransactions marks all expired pending transactions as expired
func (r *DepositRepository) ExpirePendingTransactions(ctx context.Context) (int, error) {
	query := `
		UPDATE deposit_transactions
		SET status = $1,
			status_message = 'Transaction expired due to timeout',
			failed_at = CURRENT_TIMESTAMP
		WHERE status = $2
		  AND expires_at IS NOT NULL
		  AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := r.db.ExecContext(ctx, query, DepositStatusExpired, DepositStatusPending)
	if err != nil {
		return 0, err
	}

	rows, _ := result.RowsAffected()
	return int(rows), nil
}

// GetByUserID returns all deposit transactions for a user
func (r *DepositRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*DepositTransaction, error) {
	query := `
		SELECT id, user_id, amount, currency, fee, net_amount,
			   provider_code, status, status_message,
			   completed_at, cancelled_at, failed_at,
			   wallet_credited, created_at, updated_at
		FROM deposit_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*DepositTransaction
	for rows.Next() {
		var tx DepositTransaction
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.Amount, &tx.Currency, &tx.Fee, &tx.NetAmount,
			&tx.ProviderCode, &tx.Status, &tx.StatusMessage,
			&tx.CompletedAt, &tx.CancelledAt, &tx.FailedAt,
			&tx.WalletCredited, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// IncrementRetryCount increments the retry count for a transaction
func (r *DepositRepository) IncrementRetryCount(ctx context.Context, id string) error {
	query := `
		UPDATE deposit_transactions
		SET retry_count = retry_count + 1
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetStats returns statistics for deposit transactions
func (r *DepositRepository) GetStats(ctx context.Context, providerCode string, startDate, endDate time.Time) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'pending') AS pending,
			COUNT(*) FILTER (WHERE status = 'failed') AS failed,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
			COUNT(*) FILTER (WHERE status = 'expired') AS expired,
			COALESCE(SUM(amount) FILTER (WHERE status = 'completed'), 0) AS total_amount,
			COALESCE(SUM(fee) FILTER (WHERE status = 'completed'), 0) AS total_fees
		FROM deposit_transactions
		WHERE ($1 = '' OR provider_code = $1)
		  AND created_at BETWEEN $2 AND $3
	`

	var total, completed, pending, failed, cancelled, expired int
	var totalAmount, totalFees float64

	err := r.db.QueryRowContext(ctx, query, providerCode, startDate, endDate).Scan(
		&total, &completed, &pending, &failed, &cancelled, &expired,
		&totalAmount, &totalFees,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":        total,
		"completed":    completed,
		"pending":      pending,
		"failed":       failed,
		"cancelled":    cancelled,
		"expired":      expired,
		"total_amount": totalAmount,
		"total_fees":   totalFees,
	}, nil
}

// Helper functions
func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
