package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/exchange-service/internal/models"
)

type ExchangeRepository struct {
	db *sql.DB
}

func NewExchangeRepository(db *sql.DB) *ExchangeRepository {
	return &ExchangeRepository{db: db}
}

func (r *ExchangeRepository) Create(exchange *models.Exchange) error {
	query := `
		INSERT INTO exchanges (user_id, from_wallet_id, to_wallet_id, from_currency, to_currency, 
		                      from_amount, to_amount, exchange_rate, fee, fee_percentage, status,
		                      destination_amount, destination_currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	err := r.db.QueryRow(query, 
		exchange.UserID, exchange.FromWalletID, exchange.ToWalletID,
		exchange.FromCurrency, exchange.ToCurrency, exchange.FromAmount,
		exchange.ToAmount, exchange.ExchangeRate, exchange.Fee,
		exchange.FeePercentage, exchange.Status, exchange.DestinationAmount,
		exchange.DestinationCurrency).Scan(&exchange.ID, &exchange.CreatedAt)

	return err
}

func (r *ExchangeRepository) GetByID(id string) (*models.Exchange, error) {
	exchange := &models.Exchange{}
	query := `
		SELECT id, user_id, from_wallet_id, to_wallet_id, from_currency, to_currency,
		       from_amount, to_amount, exchange_rate, fee, fee_percentage, status,
		       destination_amount, destination_currency, created_at, updated_at, completed_at
		FROM exchanges WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&exchange.ID, &exchange.UserID, &exchange.FromWalletID, &exchange.ToWalletID,
		&exchange.FromCurrency, &exchange.ToCurrency, &exchange.FromAmount,
		&exchange.ToAmount, &exchange.ExchangeRate, &exchange.Fee,
		&exchange.FeePercentage, &exchange.Status, &exchange.DestinationAmount,
		&exchange.DestinationCurrency, &exchange.CreatedAt, &exchange.UpdatedAt,
		&exchange.CompletedAt)

	return exchange, err
}

func (r *ExchangeRepository) GetByUserID(userID string, limit int) ([]*models.Exchange, error) {
	query := `
		SELECT id, user_id, from_wallet_id, to_wallet_id, from_currency, to_currency,
		       from_amount, to_amount, exchange_rate, fee, fee_percentage, status,
		       destination_amount, destination_currency, created_at, updated_at, completed_at
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
			&exchange.FeePercentage, &exchange.Status, &exchange.DestinationAmount,
			&exchange.DestinationCurrency, &exchange.CreatedAt, &exchange.UpdatedAt,
			&exchange.CompletedAt)
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