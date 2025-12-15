package services

import (
	"fmt"
	"strings"
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

// parsePair splits a trading pair like "BTC/USD" into fromCurrency and toCurrency
func parsePair(pair string) (string, string) {
	parts := strings.Split(pair, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return pair, "USD"
}

func (s *TradingService) PlaceMarketOrder(userID, pair, side string, amount float64) (*models.TradingOrder, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	if side != "buy" && side != "sell" {
		return nil, fmt.Errorf("side must be 'buy' or 'sell'")
	}

	fromCurrency, toCurrency := parsePair(pair)

	order := &models.TradingOrder{
		UserID:          userID,
		OrderType:       "market",
		FromCurrency:    fromCurrency,
		ToCurrency:      toCurrency,
		Side:            side,
		Amount:          amount,
		RemainingAmount: amount,
		Status:          "open",
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

	fromCurrency, toCurrency := parsePair(pair)

	order := &models.TradingOrder{
		UserID:          userID,
		OrderType:       "limit",
		FromCurrency:    fromCurrency,
		ToCurrency:      toCurrency,
		Side:            side,
		Amount:          amount,
		Price:           &price,
		RemainingAmount: amount,
		Status:          "open",
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

	fromCurrency, toCurrency := parsePair(pair)

	order := &models.TradingOrder{
		UserID:          userID,
		OrderType:       "stop_loss",
		FromCurrency:    fromCurrency,
		ToCurrency:      toCurrency,
		Side:            side,
		Amount:          amount,
		StopPrice:       &stopPrice,
		RemainingAmount: amount,
		Status:          "open",
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

func (s *TradingService) GetActiveOrders(fromCurrency, toCurrency string) ([]*models.TradingOrder, error) {
	return s.orderRepo.GetActiveOrders(fromCurrency, toCurrency)
}

func (s *TradingService) GetOrderBook(pair string) (*models.OrderBook, error) {
	fromCurrency, toCurrency := parsePair(pair)
	
	// Simulated order book
	orderBook := &models.OrderBook{
		Symbol: pair,
		Bids: []models.OrderLevel{
			{Price: 43480.0, Amount: 1.5, Count: 3},
			{Price: 43470.0, Amount: 2.3, Count: 5},
			{Price: 43460.0, Amount: 0.8, Count: 2},
			{Price: 43450.0, Amount: 4.2, Count: 8},
			{Price: 43440.0, Amount: 1.1, Count: 2},
		},
		Asks: []models.OrderLevel{
			{Price: 43520.0, Amount: 1.2, Count: 2},
			{Price: 43530.0, Amount: 2.0, Count: 4},
			{Price: 43540.0, Amount: 0.6, Count: 1},
			{Price: 43550.0, Amount: 3.5, Count: 6},
			{Price: 43560.0, Amount: 1.8, Count: 3},
		},
		UpdatedAt: time.Now(),
	}
	
	// Use currencies to avoid unused variable warning
	_ = fromCurrency
	_ = toCurrency
	
	return orderBook, nil
}

func (s *TradingService) GetRecentTrades(pair string) ([]*models.Trade, error) {
	fromCurrency, toCurrency := parsePair(pair)
	
	// Simulated recent trades
	trades := []*models.Trade{
		{ID: "t1", FromCurrency: fromCurrency, ToCurrency: toCurrency, Amount: 0.5, Price: 43500.0, Fee: 10.88, ExecutedAt: time.Now().Add(-1 * time.Minute)},
		{ID: "t2", FromCurrency: fromCurrency, ToCurrency: toCurrency, Amount: 1.2, Price: 43510.0, Fee: 26.11, ExecutedAt: time.Now().Add(-2 * time.Minute)},
		{ID: "t3", FromCurrency: fromCurrency, ToCurrency: toCurrency, Amount: 0.3, Price: 43495.0, Fee: 6.52, ExecutedAt: time.Now().Add(-3 * time.Minute)},
		{ID: "t4", FromCurrency: fromCurrency, ToCurrency: toCurrency, Amount: 2.0, Price: 43520.0, Fee: 43.52, ExecutedAt: time.Now().Add(-5 * time.Minute)},
		{ID: "t5", FromCurrency: fromCurrency, ToCurrency: toCurrency, Amount: 0.8, Price: 43480.0, Fee: 17.39, ExecutedAt: time.Now().Add(-7 * time.Minute)},
	}
	
	return trades, nil
}

func (s *TradingService) CancelOrder(orderID string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, "cancelled")
}

func (s *TradingService) GetPortfolio(userID string) (*models.Portfolio, error) {
	// This is a simplified portfolio calculation
	// In a real system, this would aggregate from wallet service
	
	portfolio := &models.Portfolio{
		TotalValue:   0,
		TotalPnL:     0,
		TotalPnLPerc: 0,
		Holdings:     []models.Holding{},
		Performance: models.PerformanceMetrics{
			TotalTrades:   10,
			WinningTrades: 7,
			LosingTrades:  3,
			WinRate:       70.0,
			TotalVolume:   50000.0,
			TotalFees:     75.0,
			BestTrade:     1250.0,
			WorstTrade:    -320.0,
			ProfitFactor:  2.5,
		},
	}

	// Simulate some holdings
	holdings := []models.Holding{
		{
			Currency:      "BTC",
			Amount:        0.5,
			Value:         21750.0,
			AvgBuyPrice:   40000.0,
			CurrentPrice:  43500.0,
			PnL:           1750.0,
			PnLPercentage: 8.75,
		},
		{
			Currency:      "ETH",
			Amount:        10.0,
			Value:         24500.0,
			AvgBuyPrice:   2200.0,
			CurrentPrice:  2450.0,
			PnL:           2500.0,
			PnLPercentage: 11.36,
		},
		{
			Currency:      "USD",
			Amount:        3750.0,
			Value:         3750.0,
			AvgBuyPrice:   1.0,
			CurrentPrice:  1.0,
			PnL:           0.0,
			PnLPercentage: 0.0,
		},
	}

	portfolio.Holdings = holdings
	
	// Calculate total value and PnL
	for _, holding := range holdings {
		portfolio.TotalValue += holding.Value
		portfolio.TotalPnL += holding.PnL
	}
	
	if portfolio.TotalValue > 0 {
		portfolio.TotalPnLPerc = (portfolio.TotalPnL / (portfolio.TotalValue - portfolio.TotalPnL)) * 100
	}

	return portfolio, nil
}

func (s *TradingService) processMarketOrder(order *models.TradingOrder) {
	// Simulate order processing time
	time.Sleep(2 * time.Second)

	// Update order as filled
	s.orderRepo.UpdateOrderStatus(order.ID, "filled")
	s.orderRepo.UpdateFilledAmount(order.ID, order.Amount, 0)

	order.Status = "filled"
	order.FilledAmount = order.Amount
	order.RemainingAmount = 0
	now := time.Now()
	order.ExecutedAt = &now
}