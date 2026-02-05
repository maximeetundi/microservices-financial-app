package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PayoutRepository handles payout database operations
type PayoutRepository struct {
	db *sql.DB
}

// NewPayoutRepository creates a new payout repository
func NewPayoutRepository(db *sql.DB) *PayoutRepository {
	return &PayoutRepository{db: db}
}

// PayoutTransaction represents a payout transaction in the database
type PayoutTransaction struct {
	ID                   string     `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	WalletID             string     `json:"wallet_id" db:"wallet_id"`
	Amount               float64    `json:"amount" db:"amount"`
	Currency             string     `json:"currency" db:"currency"`
	Fee                  float64    `json:"fee" db:"fee"`
	AmountReceived       float64    `json:"amount_received" db:"amount_received"`
	ProviderCode         string     `json:"provider_code" db:"provider_code"`
	PayoutMethod         string     `json:"payout_method" db:"payout_method"`
	AggregatorInstanceID string     `json:"aggregator_instance_id" db:"aggregator_instance_id"`
	HotWalletID          string     `json:"hot_wallet_id" db:"hot_wallet_id"`
	ProviderReference    string     `json:"provider_reference" db:"provider_reference"`
	RecipientName        string     `json:"recipient_name" db:"recipient_name"`
	RecipientEmail       string     `json:"recipient_email" db:"recipient_email"`
	RecipientPhone       string     `json:"recipient_phone" db:"recipient_phone"`
	BankCode             string     `json:"bank_code" db:"bank_code"`
	BankName             string     `json:"bank_name" db:"bank_name"`
	AccountNumber        string     `json:"account_number" db:"account_number"`
	IBAN                 string     `json:"iban" db:"iban"`
	SwiftCode            string     `json:"swift_code" db:"swift_code"`
	RoutingNumber        string     `json:"routing_number" db:"routing_number"`
	MobileOperator       string     `json:"mobile_operator" db:"mobile_operator"`
	MobileNumber         string     `json:"mobile_number" db:"mobile_number"`
	PayPalEmail          string     `json:"paypal_email" db:"paypal_email"`
	Narration            string     `json:"narration" db:"narration"`
	Country              string     `json:"country" db:"country"`
	Status               string     `json:"status" db:"status"`
	StatusMessage        string     `json:"status_message" db:"status_message"`
	WebhookData          []byte     `json:"webhook_data" db:"webhook_data"`
	ExpiresAt            time.Time  `json:"expires_at" db:"expires_at"`
	CompletedAt          *time.Time `json:"completed_at" db:"completed_at"`
	FailedAt             *time.Time `json:"failed_at" db:"failed_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	IPAddress            string     `json:"ip_address" db:"ip_address"`
	UserAgent            string     `json:"user_agent" db:"user_agent"`
	Metadata             []byte     `json:"metadata" db:"metadata"`
}

// CreatePayoutRequest represents a request to create a payout
type CreatePayoutRequest struct {
	ID                   string
	UserID               string
	WalletID             string
	Amount               float64
	Currency             string
	Fee                  float64
	AmountReceived       float64
	ProviderCode         string
	PayoutMethod         string
	AggregatorInstanceID string
	HotWalletID          string
	RecipientName        string
	RecipientEmail       string
	RecipientPhone       string
	BankCode             string
	BankName             string
	AccountNumber        string
	IBAN                 string
	SwiftCode            string
	RoutingNumber        string
	MobileOperator       string
	MobileNumber         string
	PayPalEmail          string
	Narration            string
	Country              string
	ExpiresAt            time.Time
	IPAddress            string
	UserAgent            string
	Metadata             map[string]interface{}
}

// InitSchema creates the payout_transactions table
func (r *PayoutRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS payout_transactions (
		id VARCHAR(100) PRIMARY KEY,
		user_id UUID NOT NULL,
		wallet_id UUID NOT NULL,
		amount DECIMAL(20,8) NOT NULL,
		currency VARCHAR(10) NOT NULL,
		fee DECIMAL(20,8) DEFAULT 0,
		amount_received DECIMAL(20,8) DEFAULT 0,
		provider_code VARCHAR(50) NOT NULL,
		payout_method VARCHAR(50) NOT NULL,
		aggregator_instance_id UUID,
		hot_wallet_id VARCHAR(36),
		provider_reference VARCHAR(255),
		recipient_name VARCHAR(255),
		recipient_email VARCHAR(255),
		recipient_phone VARCHAR(50),
		bank_code VARCHAR(50),
		bank_name VARCHAR(255),
		account_number VARCHAR(100),
		iban VARCHAR(50),
		swift_code VARCHAR(20),
		routing_number VARCHAR(50),
		mobile_operator VARCHAR(50),
		mobile_number VARCHAR(50),
		paypal_email VARCHAR(255),
		narration TEXT,
		country VARCHAR(5),
		status VARCHAR(20) DEFAULT 'pending',
		status_message TEXT,
		webhook_data JSONB,
		expires_at TIMESTAMP,
		completed_at TIMESTAMP,
		failed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		ip_address VARCHAR(50),
		user_agent TEXT,
		metadata JSONB
	);

	CREATE INDEX IF NOT EXISTS idx_payout_user_id ON payout_transactions(user_id);
	CREATE INDEX IF NOT EXISTS idx_payout_status ON payout_transactions(status);
	CREATE INDEX IF NOT EXISTS idx_payout_provider ON payout_transactions(provider_code);
	CREATE INDEX IF NOT EXISTS idx_payout_created_at ON payout_transactions(created_at);
	CREATE INDEX IF NOT EXISTS idx_payout_provider_ref ON payout_transactions(provider_reference);
	`

	_, err := r.db.Exec(query)
	return err
}

// Create creates a new payout transaction
func (r *PayoutRepository) Create(ctx context.Context, req *CreatePayoutRequest) (*PayoutTransaction, error) {
	if req.ID == "" {
		req.ID = fmt.Sprintf("pay_%s_%d", uuid.New().String()[:8], time.Now().Unix())
	}

	metadataJSON, _ := json.Marshal(req.Metadata)

	query := `
		INSERT INTO payout_transactions (
			id, user_id, wallet_id, amount, currency, fee, amount_received,
			provider_code, payout_method, aggregator_instance_id, hot_wallet_id,
			recipient_name, recipient_email, recipient_phone,
			bank_code, bank_name, account_number, iban, swift_code, routing_number,
			mobile_operator, mobile_number, paypal_email, narration, country,
			status, expires_at, ip_address, user_agent, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
			'pending', $26, $27, $28, $29, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		RETURNING id, status, created_at
	`

	var payout PayoutTransaction
	err := r.db.QueryRowContext(ctx, query,
		req.ID, req.UserID, req.WalletID, req.Amount, req.Currency, req.Fee, req.AmountReceived,
		req.ProviderCode, req.PayoutMethod, req.AggregatorInstanceID, req.HotWalletID,
		req.RecipientName, req.RecipientEmail, req.RecipientPhone,
		req.BankCode, req.BankName, req.AccountNumber, req.IBAN, req.SwiftCode, req.RoutingNumber,
		req.MobileOperator, req.MobileNumber, req.PayPalEmail, req.Narration, req.Country,
		req.ExpiresAt, req.IPAddress, req.UserAgent, metadataJSON,
	).Scan(&payout.ID, &payout.Status, &payout.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	payout.Amount = req.Amount
	payout.Currency = req.Currency
	payout.Fee = req.Fee
	payout.AmountReceived = req.AmountReceived
	payout.ProviderCode = req.ProviderCode
	payout.PayoutMethod = req.PayoutMethod
	payout.RecipientName = req.RecipientName

	return &payout, nil
}

// GetByID retrieves a payout by ID
func (r *PayoutRepository) GetByID(ctx context.Context, id string) (*PayoutTransaction, error) {
	query := `
		SELECT id, user_id, wallet_id, amount, currency, fee, amount_received,
		       provider_code, payout_method, aggregator_instance_id, hot_wallet_id,
		       COALESCE(provider_reference, ''), recipient_name, COALESCE(recipient_email, ''),
		       COALESCE(recipient_phone, ''), COALESCE(bank_code, ''), COALESCE(bank_name, ''),
		       COALESCE(account_number, ''), COALESCE(iban, ''), COALESCE(swift_code, ''),
		       COALESCE(routing_number, ''), COALESCE(mobile_operator, ''), COALESCE(mobile_number, ''),
		       COALESCE(paypal_email, ''), COALESCE(narration, ''), COALESCE(country, ''),
		       status, COALESCE(status_message, ''), completed_at, failed_at, created_at, updated_at
		FROM payout_transactions
		WHERE id = $1
	`

	var payout PayoutTransaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payout.ID, &payout.UserID, &payout.WalletID, &payout.Amount, &payout.Currency,
		&payout.Fee, &payout.AmountReceived, &payout.ProviderCode, &payout.PayoutMethod,
		&payout.AggregatorInstanceID, &payout.HotWalletID, &payout.ProviderReference,
		&payout.RecipientName, &payout.RecipientEmail, &payout.RecipientPhone,
		&payout.BankCode, &payout.BankName, &payout.AccountNumber, &payout.IBAN,
		&payout.SwiftCode, &payout.RoutingNumber, &payout.MobileOperator, &payout.MobileNumber,
		&payout.PayPalEmail, &payout.Narration, &payout.Country, &payout.Status,
		&payout.StatusMessage, &payout.CompletedAt, &payout.FailedAt, &payout.CreatedAt, &payout.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("payout not found: %w", err)
	}

	return &payout, nil
}

// GetByProviderRef retrieves a payout by provider reference
func (r *PayoutRepository) GetByProviderRef(ctx context.Context, providerRef string) (*PayoutTransaction, error) {
	query := `
		SELECT id, user_id, wallet_id, amount, currency, fee, amount_received,
		       provider_code, payout_method, status
		FROM payout_transactions
		WHERE provider_reference = $1
	`

	var payout PayoutTransaction
	err := r.db.QueryRowContext(ctx, query, providerRef).Scan(
		&payout.ID, &payout.UserID, &payout.WalletID, &payout.Amount, &payout.Currency,
		&payout.Fee, &payout.AmountReceived, &payout.ProviderCode, &payout.PayoutMethod, &payout.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("payout not found: %w", err)
	}

	return &payout, nil
}

// GetByUserID retrieves payouts for a user
func (r *PayoutRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*PayoutTransaction, error) {
	query := `
		SELECT id, user_id, wallet_id, amount, currency, fee, amount_received,
		       provider_code, payout_method, COALESCE(provider_reference, ''),
		       recipient_name, COALESCE(recipient_email, ''), COALESCE(recipient_phone, ''),
		       COALESCE(bank_name, ''), COALESCE(account_number, ''),
		       COALESCE(mobile_operator, ''), COALESCE(mobile_number, ''),
		       COALESCE(paypal_email, ''), COALESCE(country, ''),
		       status, COALESCE(status_message, ''), completed_at, failed_at, created_at
		FROM payout_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query payouts: %w", err)
	}
	defer rows.Close()

	var payouts []*PayoutTransaction
	for rows.Next() {
		var p PayoutTransaction
		err := rows.Scan(
			&p.ID, &p.UserID, &p.WalletID, &p.Amount, &p.Currency, &p.Fee, &p.AmountReceived,
			&p.ProviderCode, &p.PayoutMethod, &p.ProviderReference,
			&p.RecipientName, &p.RecipientEmail, &p.RecipientPhone,
			&p.BankName, &p.AccountNumber, &p.MobileOperator, &p.MobileNumber,
			&p.PayPalEmail, &p.Country, &p.Status, &p.StatusMessage,
			&p.CompletedAt, &p.FailedAt, &p.CreatedAt,
		)
		if err != nil {
			continue
		}
		payouts = append(payouts, &p)
	}

	return payouts, nil
}

// UpdateStatus updates the status of a payout
func (r *PayoutRepository) UpdateStatus(ctx context.Context, id, status, message string) error {
	query := `
		UPDATE payout_transactions
		SET status = $1, status_message = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, message, id)
	return err
}

// UpdateProviderReference updates the provider reference
func (r *PayoutRepository) UpdateProviderReference(ctx context.Context, id, providerRef string) error {
	query := `
		UPDATE payout_transactions
		SET provider_reference = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, providerRef, id)
	return err
}

// MarkCompleted marks a payout as completed
func (r *PayoutRepository) MarkCompleted(ctx context.Context, id, providerRef string) error {
	query := `
		UPDATE payout_transactions
		SET status = 'completed',
		    provider_reference = COALESCE(NULLIF($1, ''), provider_reference),
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, providerRef, id)
	return err
}

// MarkFailed marks a payout as failed
func (r *PayoutRepository) MarkFailed(ctx context.Context, id, reason string) error {
	query := `
		UPDATE payout_transactions
		SET status = 'failed',
		    status_message = $1,
		    failed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, reason, id)
	return err
}

// MarkCancelled marks a payout as cancelled
func (r *PayoutRepository) MarkCancelled(ctx context.Context, id, reason string) error {
	query := `
		UPDATE payout_transactions
		SET status = 'cancelled',
		    status_message = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, reason, id)
	return err
}

// UpdateWebhookData stores webhook payload for debugging
func (r *PayoutRepository) UpdateWebhookData(ctx context.Context, id string, data []byte) error {
	query := `
		UPDATE payout_transactions
		SET webhook_data = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, data, id)
	return err
}

// GetPendingPayouts retrieves all pending payouts (for status checks)
func (r *PayoutRepository) GetPendingPayouts(ctx context.Context, olderThan time.Duration) ([]*PayoutTransaction, error) {
	cutoff := time.Now().Add(-olderThan)

	query := `
		SELECT id, user_id, wallet_id, amount, currency, provider_code, payout_method,
		       COALESCE(provider_reference, ''), status, created_at
		FROM payout_transactions
		WHERE status IN ('pending', 'processing')
		  AND created_at < $1
		ORDER BY created_at ASC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending payouts: %w", err)
	}
	defer rows.Close()

	var payouts []*PayoutTransaction
	for rows.Next() {
		var p PayoutTransaction
		err := rows.Scan(
			&p.ID, &p.UserID, &p.WalletID, &p.Amount, &p.Currency,
			&p.ProviderCode, &p.PayoutMethod, &p.ProviderReference, &p.Status, &p.CreatedAt,
		)
		if err != nil {
			continue
		}
		payouts = append(payouts, &p)
	}

	return payouts, nil
}

// ExpireOldPayouts marks old pending payouts as expired
func (r *PayoutRepository) ExpireOldPayouts(ctx context.Context) (int64, error) {
	query := `
		UPDATE payout_transactions
		SET status = 'expired',
		    status_message = 'Payout expired due to timeout',
		    failed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE status = 'pending'
		  AND expires_at IS NOT NULL
		  AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GetPayoutStats returns statistics for a time period
func (r *PayoutRepository) GetPayoutStats(ctx context.Context, since time.Time) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'completed') as completed,
			COUNT(*) FILTER (WHERE status = 'failed') as failed,
			COUNT(*) FILTER (WHERE status = 'pending' OR status = 'processing') as pending,
			COALESCE(SUM(amount) FILTER (WHERE status = 'completed'), 0) as total_volume,
			COALESCE(SUM(fee) FILTER (WHERE status = 'completed'), 0) as total_fees
		FROM payout_transactions
		WHERE created_at >= $1
	`

	var stats struct {
		Total       int     `json:"total"`
		Completed   int     `json:"completed"`
		Failed      int     `json:"failed"`
		Pending     int     `json:"pending"`
		TotalVolume float64 `json:"total_volume"`
		TotalFees   float64 `json:"total_fees"`
	}

	err := r.db.QueryRowContext(ctx, query, since).Scan(
		&stats.Total, &stats.Completed, &stats.Failed, &stats.Pending,
		&stats.TotalVolume, &stats.TotalFees,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":        stats.Total,
		"completed":    stats.Completed,
		"failed":       stats.Failed,
		"pending":      stats.Pending,
		"total_volume": stats.TotalVolume,
		"total_fees":   stats.TotalFees,
	}, nil
}

// GetPayoutsByProvider returns payouts grouped by provider
func (r *PayoutRepository) GetPayoutsByProvider(ctx context.Context, since time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT
			provider_code,
			COUNT(*) as count,
			COALESCE(SUM(amount), 0) as volume,
			COUNT(*) FILTER (WHERE status = 'completed') as success_count,
			COUNT(*) FILTER (WHERE status = 'failed') as failed_count
		FROM payout_transactions
		WHERE created_at >= $1
		GROUP BY provider_code
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var provider string
		var count, successCount, failedCount int
		var volume float64

		if err := rows.Scan(&provider, &count, &volume, &successCount, &failedCount); err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"provider":      provider,
			"count":         count,
			"volume":        volume,
			"success_count": successCount,
			"failed_count":  failedCount,
		})
	}

	return results, nil
}
