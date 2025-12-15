package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
)

// PaymentRequestRepository handles payment request database operations
type PaymentRequestRepository struct {
	db *sql.DB
}

// NewPaymentRequestRepository creates a new payment request repository
func NewPaymentRequestRepository(db *sql.DB) *PaymentRequestRepository {
	return &PaymentRequestRepository{db: db}
}

// InitTable creates the payment_requests table if it doesn't exist
func (r *PaymentRequestRepository) InitTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS payment_requests (
			id VARCHAR(64) PRIMARY KEY,
			merchant_id VARCHAR(64) NOT NULL,
			wallet_id VARCHAR(64) NOT NULL,
			type VARCHAR(20) NOT NULL,
			amount DECIMAL(20, 8),
			min_amount DECIMAL(20, 8),
			max_amount DECIMAL(20, 8),
			currency VARCHAR(10) NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			items_json TEXT,
			qr_code_data TEXT NOT NULL,
			payment_link VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP,
			never_expires BOOLEAN DEFAULT FALSE,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			paid_amount DECIMAL(20, 8),
			paid_by VARCHAR(64),
			paid_at TIMESTAMP,
			transaction_id VARCHAR(64),
			reusable BOOLEAN DEFAULT FALSE,
			times_used INTEGER DEFAULT 0,
			metadata_json TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_payment_requests_merchant ON payment_requests(merchant_id);
		CREATE INDEX IF NOT EXISTS idx_payment_requests_status ON payment_requests(status);

		CREATE TABLE IF NOT EXISTS payment_history (
			id VARCHAR(64) PRIMARY KEY,
			payment_request_id VARCHAR(64) NOT NULL,
			merchant_id VARCHAR(64) NOT NULL,
			customer_id VARCHAR(64) NOT NULL,
			amount DECIMAL(20, 8) NOT NULL,
			fee DECIMAL(20, 8) NOT NULL,
			net_amount DECIMAL(20, 8) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			transaction_id VARCHAR(64) NOT NULL,
			paid_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_payment_history_merchant ON payment_history(merchant_id);
		CREATE INDEX IF NOT EXISTS idx_payment_history_customer ON payment_history(customer_id);
	`
	
	_, err := r.db.Exec(query)
	return err
}

// Create creates a new payment request
func (r *PaymentRequestRepository) Create(payment *models.PaymentRequest) error {
	itemsJSON, _ := json.Marshal(payment.Items)
	metadataJSON, _ := json.Marshal(payment.Metadata)

	query := `
		INSERT INTO payment_requests (
			id, merchant_id, wallet_id, type, amount, min_amount, max_amount,
			currency, title, description, items_json, qr_code_data, payment_link,
			expires_at, never_expires, status, reusable, metadata_json, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.db.Exec(query,
		payment.ID, payment.MerchantID, payment.WalletID, payment.Type,
		payment.Amount, payment.MinAmount, payment.MaxAmount,
		payment.Currency, payment.Title, payment.Description,
		string(itemsJSON), payment.QRCodeData, payment.PaymentLink,
		payment.ExpiresAt, payment.NeverExpires, payment.Status,
		payment.Reusable, string(metadataJSON), payment.CreatedAt, payment.UpdatedAt,
	)

	return err
}

// GetByID gets a payment request by ID
func (r *PaymentRequestRepository) GetByID(id string) (*models.PaymentRequest, error) {
	query := `
		SELECT id, merchant_id, wallet_id, type, amount, min_amount, max_amount,
			currency, title, description, items_json, qr_code_data, payment_link,
			expires_at, never_expires, status, paid_amount, paid_by, paid_at,
			transaction_id, reusable, times_used, metadata_json, created_at, updated_at
		FROM payment_requests
		WHERE id = $1
	`

	var payment models.PaymentRequest
	var itemsJSON, metadataJSON sql.NullString
	var expiresAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID, &payment.MerchantID, &payment.WalletID, &payment.Type,
		&payment.Amount, &payment.MinAmount, &payment.MaxAmount,
		&payment.Currency, &payment.Title, &payment.Description,
		&itemsJSON, &payment.QRCodeData, &payment.PaymentLink,
		&expiresAt, &payment.NeverExpires, &payment.Status,
		&payment.PaidAmount, &payment.PaidBy, &payment.PaidAt,
		&payment.TransactionID, &payment.Reusable, &payment.TimesUsed,
		&metadataJSON, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment request not found")
		}
		return nil, err
	}

	if expiresAt.Valid {
		payment.ExpiresAt = &expiresAt.Time
	}

	if itemsJSON.Valid {
		json.Unmarshal([]byte(itemsJSON.String), &payment.Items)
	}

	if metadataJSON.Valid {
		json.Unmarshal([]byte(metadataJSON.String), &payment.Metadata)
	}

	return &payment, nil
}

// GetByMerchantID gets all payment requests for a merchant
func (r *PaymentRequestRepository) GetByMerchantID(merchantID string, limit, offset int) ([]models.PaymentRequest, error) {
	query := `
		SELECT id, merchant_id, wallet_id, type, amount, currency, title,
			payment_link, expires_at, never_expires, status, reusable,
			times_used, created_at
		FROM payment_requests
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, merchantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentRequest
	for rows.Next() {
		var p models.PaymentRequest
		var expiresAt sql.NullTime

		err := rows.Scan(
			&p.ID, &p.MerchantID, &p.WalletID, &p.Type, &p.Amount,
			&p.Currency, &p.Title, &p.PaymentLink, &expiresAt,
			&p.NeverExpires, &p.Status, &p.Reusable, &p.TimesUsed, &p.CreatedAt,
		)
		if err != nil {
			continue
		}

		if expiresAt.Valid {
			p.ExpiresAt = &expiresAt.Time
		}

		payments = append(payments, p)
	}

	return payments, nil
}

// Update updates a payment request
func (r *PaymentRequestRepository) Update(payment *models.PaymentRequest) error {
	query := `
		UPDATE payment_requests
		SET status = $2, paid_amount = $3, paid_by = $4, paid_at = $5,
			transaction_id = $6, times_used = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		payment.ID, payment.Status, payment.PaidAmount, payment.PaidBy,
		payment.PaidAt, payment.TransactionID, payment.TimesUsed, time.Now(),
	)

	return err
}

// UpdateStatus updates only the status of a payment request
func (r *PaymentRequestRepository) UpdateStatus(id, status string) error {
	query := `UPDATE payment_requests SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, time.Now())
	return err
}

// CreateHistory creates a payment history record
func (r *PaymentRequestRepository) CreateHistory(history *models.PaymentHistory) error {
	query := `
		INSERT INTO payment_history (
			id, payment_request_id, merchant_id, customer_id,
			amount, fee, net_amount, currency, transaction_id, paid_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		history.ID, history.PaymentRequestID, history.MerchantID,
		history.CustomerID, history.Amount, history.Fee, history.NetAmount,
		history.Currency, history.TransactionID, history.PaidAt,
	)

	return err
}

// GetHistoryByMerchantID gets payment history for a merchant
func (r *PaymentRequestRepository) GetHistoryByMerchantID(merchantID string, limit, offset int) ([]models.PaymentHistory, error) {
	query := `
		SELECT id, payment_request_id, merchant_id, customer_id,
			amount, fee, net_amount, currency, transaction_id, paid_at
		FROM payment_history
		WHERE merchant_id = $1
		ORDER BY paid_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, merchantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PaymentHistory
	for rows.Next() {
		var h models.PaymentHistory
		err := rows.Scan(
			&h.ID, &h.PaymentRequestID, &h.MerchantID, &h.CustomerID,
			&h.Amount, &h.Fee, &h.NetAmount, &h.Currency,
			&h.TransactionID, &h.PaidAt,
		)
		if err != nil {
			continue
		}
		history = append(history, h)
	}

	return history, nil
}

// GetMerchantStats gets statistics for a merchant
func (r *PaymentRequestRepository) GetMerchantStats(merchantID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_payments,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(fee), 0) as total_fees,
			COALESCE(SUM(net_amount), 0) as total_net
		FROM payment_history
		WHERE merchant_id = $1
	`

	var totalPayments int
	var totalAmount, totalFees, totalNet float64

	err := r.db.QueryRow(query, merchantID).Scan(
		&totalPayments, &totalAmount, &totalFees, &totalNet,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_payments": totalPayments,
		"total_amount":   totalAmount,
		"total_fees":     totalFees,
		"total_net":      totalNet,
	}, nil
}
