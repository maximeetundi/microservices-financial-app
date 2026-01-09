package repository

import (
	"database/sql"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
)

type ExchangeRepository struct {
	db *sql.DB
}

func NewExchangeRepository(db *sql.DB) *ExchangeRepository {
	return &ExchangeRepository{db: db}
}

func (r *ExchangeRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS exchanges (
		id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id VARCHAR(36) NOT NULL,
		from_wallet_id VARCHAR(36) NOT NULL,
		to_wallet_id VARCHAR(36) NOT NULL,
		from_currency VARCHAR(10) NOT NULL,
		to_currency VARCHAR(10) NOT NULL,
		from_amount DECIMAL(20, 8) NOT NULL,
		to_amount DECIMAL(20, 8) NOT NULL,
		exchange_rate DECIMAL(20, 8) NOT NULL,
		fee DECIMAL(20, 8) DEFAULT 0,
		fee_percentage DECIMAL(10, 4) DEFAULT 0,
		status VARCHAR(20) NOT NULL,
		quote_id VARCHAR(36),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP
	);
	`
	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	// Migration: Add fee_percentage if it doesn't exist
	alterQuery := `ALTER TABLE exchanges ADD COLUMN IF NOT EXISTS fee_percentage DECIMAL(10, 4) DEFAULT 0;`
	if _, err := r.db.Exec(alterQuery); err != nil {
		return err
	}

	// Create quotes table if not exists
	quotesQuery := `
	CREATE TABLE IF NOT EXISTS quotes (
		id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id VARCHAR(36) NOT NULL,
		from_currency VARCHAR(10) NOT NULL,
		to_currency VARCHAR(10) NOT NULL,
		from_amount DECIMAL(20, 8) NOT NULL,
		to_amount DECIMAL(20, 8) NOT NULL,
		exchange_rate DECIMAL(20, 8) NOT NULL,
		fee DECIMAL(20, 8) DEFAULT 0,
		fee_percentage DECIMAL(10, 4) DEFAULT 0,
		valid_until TIMESTAMP NOT NULL,
		estimated_delivery VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := r.db.Exec(quotesQuery); err != nil {
		return err
	}

	// Migration: Add fee_percentage to quotes if it doesn't exist
	alterQuotesQuery := `ALTER TABLE quotes ADD COLUMN IF NOT EXISTS fee_percentage DECIMAL(10, 4) DEFAULT 0;`
	if _, err := r.db.Exec(alterQuotesQuery); err != nil {
		return err
	}

	return nil
}

func (r *ExchangeRepository) Create(exchange *models.Exchange) error {
	query := `
		INSERT INTO exchanges (user_id, from_wallet_id, to_wallet_id, from_currency, to_currency, 
		                      from_amount, to_amount, exchange_rate, fee, fee_percentage, status, quote_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	err := r.db.QueryRow(query, 
		exchange.UserID, exchange.FromWalletID, exchange.ToWalletID,
		exchange.FromCurrency, exchange.ToCurrency, exchange.FromAmount,
		exchange.ToAmount, exchange.ExchangeRate, exchange.Fee,
		exchange.FeePercentage, exchange.Status, exchange.QuoteID).Scan(&exchange.ID, &exchange.CreatedAt)

	return err
}

func (r *ExchangeRepository) GetByID(id string) (*models.Exchange, error) {
	exchange := &models.Exchange{}
	query := `
		SELECT id, user_id, from_wallet_id, to_wallet_id, from_currency, to_currency,
		       from_amount, to_amount, exchange_rate, fee, fee_percentage, status,
		       quote_id, created_at, updated_at, completed_at
		FROM exchanges WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&exchange.ID, &exchange.UserID, &exchange.FromWalletID, &exchange.ToWalletID,
		&exchange.FromCurrency, &exchange.ToCurrency, &exchange.FromAmount,
		&exchange.ToAmount, &exchange.ExchangeRate, &exchange.Fee,
		&exchange.FeePercentage, &exchange.Status, &exchange.QuoteID,
		&exchange.CreatedAt, &exchange.UpdatedAt, &exchange.CompletedAt)

	return exchange, err
}

func (r *ExchangeRepository) GetByUserID(userID string, limit int) ([]*models.Exchange, error) {
	query := `
		SELECT id, user_id, from_wallet_id, to_wallet_id, from_currency, to_currency,
		       from_amount, to_amount, exchange_rate, fee, fee_percentage, status,
		       quote_id, created_at, updated_at, completed_at
		FROM exchanges 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchanges []*models.Exchange
	for rows.Next() {
		exchange := &models.Exchange{}
		err := rows.Scan(
			&exchange.ID, &exchange.UserID, &exchange.FromWalletID, &exchange.ToWalletID,
			&exchange.FromCurrency, &exchange.ToCurrency, &exchange.FromAmount,
			&exchange.ToAmount, &exchange.ExchangeRate, &exchange.Fee,
			&exchange.FeePercentage, &exchange.Status, &exchange.QuoteID,
			&exchange.CreatedAt, &exchange.UpdatedAt, &exchange.CompletedAt)
		if err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

func (r *ExchangeRepository) UpdateStatus(id, status string) error {
	query := `UPDATE exchanges SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *ExchangeRepository) CreateQuote(quote *models.Quote) error {
	query := `
		INSERT INTO quotes (user_id, from_currency, to_currency, from_amount, to_amount,
		                   exchange_rate, fee, fee_percentage, valid_until, estimated_delivery)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	return r.db.QueryRow(query,
		quote.UserID, quote.FromCurrency, quote.ToCurrency,
		quote.FromAmount, quote.ToAmount, quote.ExchangeRate,
		quote.Fee, quote.FeePercentage, quote.ValidUntil,
		quote.EstimatedDelivery).Scan(&quote.ID, &quote.CreatedAt)
}

func (r *ExchangeRepository) GetQuote(id string) (*models.Quote, error) {
	quote := &models.Quote{}
	query := `
		SELECT id, user_id, from_currency, to_currency, from_amount, to_amount,
		       exchange_rate, fee, fee_percentage, valid_until, estimated_delivery, created_at
		FROM quotes WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&quote.ID, &quote.UserID, &quote.FromCurrency, &quote.ToCurrency,
		&quote.FromAmount, &quote.ToAmount, &quote.ExchangeRate,
		&quote.Fee, &quote.FeePercentage, &quote.ValidUntil,
		&quote.EstimatedDelivery, &quote.CreatedAt)

	return quote, err
}