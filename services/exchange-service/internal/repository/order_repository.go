package repository

import (
	"database/sql"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.TradingOrder) error {
	query := `
		INSERT INTO trading_orders (user_id, wallet_id, order_type, from_currency, to_currency, side, 
		                           amount, price, stop_price, filled_amount, remaining_amount, status, fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	return r.db.QueryRow(query,
		order.UserID, order.WalletID, order.OrderType, order.FromCurrency, order.ToCurrency,
		order.Side, order.Amount, order.Price, order.StopPrice, 
		order.FilledAmount, order.RemainingAmount, order.Status, order.Fee).Scan(&order.ID, &order.CreatedAt)
}

func (r *OrderRepository) GetOrderByID(id string) (*models.TradingOrder, error) {
	order := &models.TradingOrder{}
	query := `
		SELECT id, user_id, wallet_id, order_type, from_currency, to_currency, side, 
		       amount, price, stop_price, filled_amount, remaining_amount, status, fee,
		       executed_at, expires_at, created_at, updated_at
		FROM trading_orders WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.UserID, &order.WalletID, &order.OrderType, 
		&order.FromCurrency, &order.ToCurrency, &order.Side,
		&order.Amount, &order.Price, &order.StopPrice, 
		&order.FilledAmount, &order.RemainingAmount, &order.Status, &order.Fee,
		&order.ExecutedAt, &order.ExpiresAt, &order.CreatedAt, &order.UpdatedAt)

	return order, err
}

func (r *OrderRepository) GetOrdersByUser(userID string) ([]*models.TradingOrder, error) {
	query := `
		SELECT id, user_id, wallet_id, order_type, from_currency, to_currency, side,
		       amount, price, stop_price, filled_amount, remaining_amount, status, fee,
		       executed_at, expires_at, created_at, updated_at
		FROM trading_orders 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.TradingOrder
	for rows.Next() {
		order := &models.TradingOrder{}
		err := rows.Scan(
			&order.ID, &order.UserID, &order.WalletID, &order.OrderType,
			&order.FromCurrency, &order.ToCurrency, &order.Side,
			&order.Amount, &order.Price, &order.StopPrice,
			&order.FilledAmount, &order.RemainingAmount, &order.Status, &order.Fee,
			&order.ExecutedAt, &order.ExpiresAt, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetActiveOrders(fromCurrency, toCurrency string) ([]*models.TradingOrder, error) {
	query := `
		SELECT id, user_id, wallet_id, order_type, from_currency, to_currency, side,
		       amount, price, stop_price, filled_amount, remaining_amount, status, fee,
		       executed_at, expires_at, created_at, updated_at
		FROM trading_orders 
		WHERE from_currency = $1 AND to_currency = $2 AND status IN ('open', 'partial')
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.TradingOrder
	for rows.Next() {
		order := &models.TradingOrder{}
		err := rows.Scan(
			&order.ID, &order.UserID, &order.WalletID, &order.OrderType,
			&order.FromCurrency, &order.ToCurrency, &order.Side,
			&order.Amount, &order.Price, &order.StopPrice,
			&order.FilledAmount, &order.RemainingAmount, &order.Status, &order.Fee,
			&order.ExecutedAt, &order.ExpiresAt, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatus(id, status string) error {
	query := `UPDATE trading_orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *OrderRepository) UpdateFilledAmount(id string, filledAmount, remainingAmount float64) error {
	query := `UPDATE trading_orders SET filled_amount = $1, remaining_amount = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, filledAmount, remainingAmount, id)
	return err
}