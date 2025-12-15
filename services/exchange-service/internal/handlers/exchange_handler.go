package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
)

type ExchangeHandler struct {
	exchangeService *services.ExchangeService
	rateService     *services.RateService
	tradingService  *services.TradingService
	walletClient    *services.WalletClient
}

func NewExchangeHandler(
	exchangeService *services.ExchangeService,
	rateService *services.RateService,
	tradingService *services.TradingService,
	walletClient *services.WalletClient,
) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
		rateService:     rateService,
		tradingService:  tradingService,
		walletClient:    walletClient,
	}
}

func (h *ExchangeHandler) GetExchangeRates(c *gin.Context) {
	rates, err := h.rateService.GetAllRates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exchange rates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

func (h *ExchangeHandler) GetSpecificRate(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")

	rate, err := h.rateService.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange rate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rate": rate})
}

func (h *ExchangeHandler) GetMarkets(c *gin.Context) {
	markets, err := h.rateService.GetMarkets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get markets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"markets": markets})
}

func (h *ExchangeHandler) GetQuote(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote, err := h.exchangeService.GetQuote(
		userID.(string),
		req.FromCurrency,
		req.ToCurrency,
		req.FromAmount,
		req.ToAmount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quote": quote})
}

func (h *ExchangeHandler) ExecuteExchange(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchange, err := h.exchangeService.ExecuteExchange(
		userID.(string),
		req.QuoteID,
		req.FromWalletID,
		req.ToWalletID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Exchange initiated successfully",
		"exchange": exchange,
	})
}

func (h *ExchangeHandler) GetExchangeHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	exchanges, err := h.exchangeService.GetUserExchanges(userID.(string), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exchange history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exchanges": exchanges})
}

func (h *ExchangeHandler) GetExchange(c *gin.Context) {
	exchangeID := c.Param("exchange_id")

	exchange, err := h.exchangeService.GetExchange(exchangeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exchange": exchange})
}

func (h *ExchangeHandler) PlaceBuyOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.BuyOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.exchangeService.BuyCrypto(
		userID.(string),
		req.Currency,
		req.PayCurrency,
		req.Amount,
		req.OrderType,
		req.LimitPrice,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Buy order placed successfully",
		"order":   order,
	})
}

func (h *ExchangeHandler) PlaceSellOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.SellOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.exchangeService.SellCrypto(
		userID.(string),
		req.Currency,
		req.ReceiveCurrency,
		req.Amount,
		req.OrderType,
		req.LimitPrice,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sell order placed successfully",
		"order":   order,
	})
}

func (h *ExchangeHandler) CreateLimitOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.LimitOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order *models.TradingOrder
	var err error

	if req.OrderType == "buy" {
		order, err = h.exchangeService.BuyCrypto(
			userID.(string),
			req.ToCurrency,
			req.FromCurrency,
			req.Amount,
			"limit",
			&req.LimitPrice,
		)
	} else {
		order, err = h.exchangeService.SellCrypto(
			userID.(string),
			req.FromCurrency,
			req.ToCurrency,
			req.Amount,
			"limit",
			&req.LimitPrice,
		)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Limit order created successfully",
		"order":   order,
	})
}

func (h *ExchangeHandler) CreateStopLoss(c *gin.Context) {
	var req models.StopLossRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implémenter les ordres stop-loss
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Stop-loss orders coming soon",
		"request": req,
	})
}

func (h *ExchangeHandler) GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	orders, err := h.exchangeService.GetUserOrders(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *ExchangeHandler) GetActiveOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	orders, err := h.exchangeService.GetUserOrders(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *ExchangeHandler) CancelOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("order_id")

	err := h.exchangeService.CancelOrder(userID.(string), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func (h *ExchangeHandler) GetPortfolio(c *gin.Context) {
	userID, _ := c.Get("user_id")

	portfolio, err := h.exchangeService.GetPortfolio(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get portfolio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *ExchangeHandler) GetPerformance(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// TODO: Implémenter les métriques de performance
	c.JSON(http.StatusOK, gin.H{
		"message": "Performance metrics coming soon",
		"user_id": userID,
	})
}

func (h *ExchangeHandler) GetOrderBook(c *gin.Context) {
	pair := c.Param("pair")

	orderBook, err := h.tradingService.GetOrderBook(pair)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orderbook": orderBook})
}

func (h *ExchangeHandler) GetRecentTrades(c *gin.Context) {
	pair := c.Param("pair")

	trades, err := h.tradingService.GetRecentTrades(pair)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trades": trades})
}

// P2P Trading handlers
func (h *ExchangeHandler) CreateP2POffer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.P2POfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implémenter le trading P2P
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "P2P trading coming soon",
		"user_id": userID,
		"request": req,
	})
}

func (h *ExchangeHandler) GetP2POffers(c *gin.Context) {
	// TODO: Implémenter la liste des offres P2P
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "P2P offers listing coming soon",
	})
}

func (h *ExchangeHandler) AcceptP2POffer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	offerID := c.Param("offer_id")

	// TODO: Implémenter l'acceptation d'offre P2P
	c.JSON(http.StatusNotImplemented, gin.H{
		"message":  "P2P offer acceptance coming soon",
		"user_id":  userID,
		"offer_id": offerID,
	})
}

func (h *ExchangeHandler) GetP2PTrades(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// TODO: Implémenter l'historique P2P
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "P2P trade history coming soon",
		"user_id": userID,
	})
}

// Admin handlers
func (h *ExchangeHandler) UpdateRates(c *gin.Context) {
	// TODO: Implémenter la mise à jour manuelle des taux
	c.JSON(http.StatusOK, gin.H{"message": "Rates update initiated"})
}

func (h *ExchangeHandler) GetAnalytics(c *gin.Context) {
	// TODO: Implémenter les analytics admin
	c.JSON(http.StatusOK, gin.H{
		"message": "Analytics data",
		"data": map[string]interface{}{
			"total_exchanges": 1250,
			"total_volume":    450000.50,
			"active_users":    89,
		},
	})
}

func (h *ExchangeHandler) GetTradingVolume(c *gin.Context) {
	period := c.DefaultQuery("period", "24h")
	
	// TODO: Implémenter les stats de volume
	c.JSON(http.StatusOK, gin.H{
		"period": period,
		"volume": 125000.75,
		"trades": 450,
	})
}

// Webhook handlers
func (h *ExchangeHandler) HandlePriceUpdate(c *gin.Context) {
	var priceUpdate map[string]interface{}
	if err := c.ShouldBindJSON(&priceUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Traiter la mise à jour de prix externe
	c.JSON(http.StatusOK, gin.H{"message": "Price update processed"})
}