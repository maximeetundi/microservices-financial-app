package repository

import (
	"database/sql"
	"fmt"

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
		INSERT INTO trading_orders (user_id, order_type, pair, side, amount, price, stop_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	return r.db.QueryRow(query,
		order.UserID, order.OrderType, order.Pair, order.Side,
		order.Amount, order.Price, order.StopPrice, order.Status).Scan(&order.ID, &order.CreatedAt)
}

func (r *OrderRepository) GetOrderByID(id string) (*models.TradingOrder, error) {
	order := &models.TradingOrder{}
	query := `
		SELECT id, user_id, order_type, pair, side, amount, price, stop_price, 
		       status, filled_amount, created_at, updated_at
		FROM trading_orders WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.UserID, &order.OrderType, &order.Pair,
		&order.Side, &order.Amount, &order.Price, &order.StopPrice,
		&order.Status, &order.FilledAmount, &order.CreatedAt, &order.UpdatedAt)

	return order, err
}

func (r *OrderRepository) GetOrdersByUser(userID string) ([]*models.TradingOrder, error) {
	query := `
		SELECT id, user_id, order_type, pair, side, amount, price, stop_price,
		       status, filled_amount, created_at, updated_at
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
			&order.ID, &order.UserID, &order.OrderType, &order.Pair,
			&order.Side, &order.Amount, &order.Price, &order.StopPrice,
			&order.Status, &order.FilledAmount, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetActiveOrders(pair string) ([]*models.TradingOrder, error) {
	query := `
		SELECT id, user_id, order_type, pair, side, amount, price, stop_price,
		       status, filled_amount, created_at, updated_at
		FROM trading_orders 
		WHERE pair = $1 AND status = 'pending'
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, pair)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.TradingOrder
	for rows.Next() {
		order := &models.TradingOrder{}
		err := rows.Scan(
			&order.ID, &order.UserID, &order.OrderType, &order.Pair,
			&order.Side, &order.Amount, &order.Price, &order.StopPrice,
			&order.Status, &order.FilledAmount, &order.CreatedAt, &order.UpdatedAt)
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

func (r *OrderRepository) UpdateFilledAmount(id string, filledAmount float64) error {
	query := `UPDATE trading_orders SET filled_amount = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, filledAmount, id)
	return err
}