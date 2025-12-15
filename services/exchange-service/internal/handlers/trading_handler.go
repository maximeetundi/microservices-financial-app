package handlers

import (
	"net/http"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
	"github.com/gin-gonic/gin"
)

type TradingHandler struct {
	tradingService *services.TradingService
}

func NewTradingHandler(tradingService *services.TradingService) *TradingHandler {
	return &TradingHandler{
		tradingService: tradingService,
	}
}

type MarketOrderRequest struct {
	Pair   string  `json:"pair" binding:"required"`
	Side   string  `json:"side" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type LimitOrderRequest struct {
	Pair   string  `json:"pair" binding:"required"`
	Side   string  `json:"side" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Price  float64 `json:"price" binding:"required,gt=0"`
}

type StopOrderRequest struct {
	Pair      string  `json:"pair" binding:"required"`
	Side      string  `json:"side" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	StopPrice float64 `json:"stop_price" binding:"required,gt=0"`
}

func (h *TradingHandler) PlaceMarketOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req MarketOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.tradingService.PlaceMarketOrder(userID.(string), req.Pair, req.Side, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *TradingHandler) PlaceLimitOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req LimitOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.tradingService.PlaceLimitOrder(userID.(string), req.Pair, req.Side, req.Amount, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *TradingHandler) PlaceStopOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req StopOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.tradingService.PlaceStopLossOrder(userID.(string), req.Pair, req.Side, req.Amount, req.StopPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *TradingHandler) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	orders, err := h.tradingService.GetUserOrders(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  len(orders),
	})
}

func (h *TradingHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	err := h.tradingService.CancelOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func (h *TradingHandler) GetPortfolio(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	portfolio, err := h.tradingService.GetPortfolio(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func (h *TradingHandler) GetOrderBook(c *gin.Context) {
	pair := c.Param("pair")
	if pair == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trading pair is required"})
		return
	}

	// Parse pair like "BTC/USD" to fromCurrency and toCurrency
	parts := strings.Split(pair, "/")
	var fromCurrency, toCurrency string
	if len(parts) == 2 {
		fromCurrency = parts[0]
		toCurrency = parts[1]
	} else {
		fromCurrency = pair
		toCurrency = "USD"
	}

	orders, err := h.tradingService.GetActiveOrders(fromCurrency, toCurrency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Separate buy and sell orders
	var buyOrders, sellOrders []interface{}
	for _, order := range orders {
		if order.Side == "buy" {
			buyOrders = append(buyOrders, order)
		} else {
			sellOrders = append(sellOrders, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"pair":        pair,
		"buy_orders":  buyOrders,
		"sell_orders": sellOrders,
		"timestamp":   "2024-01-01T00:00:00Z",
	})
}

func (h *TradingHandler) GetTickers(c *gin.Context) {
	// Simulate ticker data for common trading pairs
	tickers := []map[string]interface{}{
		{
			"symbol":        "BTC/USD",
			"price":         43500.0,
			"change_24h":    2.3,
			"volume_24h":    1234567.89,
			"high_24h":      44000.0,
			"low_24h":       42800.0,
			"bid":           43485.0,
			"ask":           43515.0,
			"last_updated":  "2024-01-01T00:00:00Z",
		},
		{
			"symbol":        "ETH/USD",
			"price":         2450.0,
			"change_24h":    -1.2,
			"volume_24h":    987654.32,
			"high_24h":      2480.0,
			"low_24h":       2420.0,
			"bid":           2448.0,
			"ask":           2452.0,
			"last_updated":  "2024-01-01T00:00:00Z",
		},
		{
			"symbol":        "EUR/USD",
			"price":         1.0856,
			"change_24h":    0.1,
			"volume_24h":    5000000.0,
			"high_24h":      1.0870,
			"low_24h":       1.0840,
			"bid":           1.0855,
			"ask":           1.0857,
			"last_updated":  "2024-01-01T00:00:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"tickers": tickers,
		"total":   len(tickers),
	})
}