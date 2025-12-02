package services

import (
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
)

type TradingService struct {
	orderRepo       *repository.OrderRepository
	exchangeService *ExchangeService
	config          *config.Config
}

func NewTradingService(orderRepo *repository.OrderRepository, exchangeService *ExchangeService, cfg *config.Config) *TradingService {
	return &TradingService{
		orderRepo:       orderRepo,
		exchangeService: exchangeService,
		config:          cfg,
	}
}

func (s *TradingService) PlaceMarketOrder(userID, pair, side string, amount float64) (*models.TradingOrder, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	if side != "buy" && side != "sell" {
		return nil, fmt.Errorf("side must be 'buy' or 'sell'")
	}

	order := &models.TradingOrder{
		UserID:    userID,
		OrderType: "market",
		Pair:      pair,
		Side:      side,
		Amount:    amount,
		Status:    "pending",
	}

	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Process market order immediately
	go s.processMarketOrder(order)

	return order, nil
}

func (s *TradingService) PlaceLimitOrder(userID, pair, side string, amount, price float64) (*models.TradingOrder, error) {
	if amount <= 0 || price <= 0 {
		return nil, fmt.Errorf("amount and price must be greater than 0")
	}

	order := &models.TradingOrder{
		UserID:    userID,
		OrderType: "limit",
		Pair:      pair,
		Side:      side,
		Amount:    amount,
		Price:     &price,
		Status:    "pending",
	}

	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Limit orders wait in the order book
	return order, nil
}

func (s *TradingService) PlaceStopLossOrder(userID, pair, side string, amount, stopPrice float64) (*models.TradingOrder, error) {
	if amount <= 0 || stopPrice <= 0 {
		return nil, fmt.Errorf("amount and stop price must be greater than 0")
	}

	order := &models.TradingOrder{
		UserID:    userID,
		OrderType: "stop_loss",
		Pair:      pair,
		Side:      side,
		Amount:    amount,
		StopPrice: &stopPrice,
		Status:    "pending",
	}

	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *TradingService) GetUserOrders(userID string) ([]*models.TradingOrder, error) {
	return s.orderRepo.GetOrdersByUser(userID)
}

func (s *TradingService) GetActiveOrders(pair string) ([]*models.TradingOrder, error) {
	return s.orderRepo.GetActiveOrders(pair)
}

func (s *TradingService) CancelOrder(orderID string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, "cancelled")
}

func (s *TradingService) GetPortfolio(userID string) (*models.Portfolio, error) {
	// This is a simplified portfolio calculation
	// In a real system, this would aggregate from wallet service
	
	portfolio := &models.Portfolio{
		UserID:       userID,
		TotalValue:   0,
		BaseCurrency: "USD",
		Holdings:     []models.PortfolioHolding{},
		Performance: models.PortfolioPerformance{
			TotalReturn:      5.2,  // 5.2%
			TotalReturnValue: 520.0,
			DayReturn:        0.8,   // 0.8%
			DayReturnValue:   80.0,
			WeekReturn:       2.1,   // 2.1%
			MonthReturn:      8.5,   // 8.5%
			YearReturn:       45.2,  // 45.2%
		},
		LastUpdated: time.Now(),
	}

	// Simulate some holdings
	holdings := []models.PortfolioHolding{
		{
			Currency:       "BTC",
			Amount:         0.5,
			Value:          21750.0,
			Percentage:     43.5,
			Change24h:      2.3,
			ChangeValue24h: 500.25,
		},
		{
			Currency:       "ETH",
			Amount:         10.0,
			Value:          24500.0,
			Percentage:     49.0,
			Change24h:      -1.2,
			ChangeValue24h: -294.0,
		},
		{
			Currency:       "USD",
			Amount:         3750.0,
			Value:          3750.0,
			Percentage:     7.5,
			Change24h:      0.0,
			ChangeValue24h: 0.0,
		},
	}

	portfolio.Holdings = holdings
	
	// Calculate total value
	for _, holding := range holdings {
		portfolio.TotalValue += holding.Value
	}

	return portfolio, nil
}

func (s *TradingService) processMarketOrder(order *models.TradingOrder) {
	// Simulate order processing time
	time.Sleep(2 * time.Second)

	// Update order as filled
	s.orderRepo.UpdateOrderStatus(order.ID, "filled")
	s.orderRepo.UpdateFilledAmount(order.ID, order.Amount)

	order.Status = "filled"
	order.FilledAmount = order.Amount
	now := time.Now()
	order.FilledAt = &now
}