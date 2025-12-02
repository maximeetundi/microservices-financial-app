package main

import (
	"log"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connections
	db, err := database.Initialize(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	redisClient, err := database.InitializeRedis(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	defer redisClient.Close()

	mqChannel, err := database.InitializeRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer mqChannel.Close()

	// Initialize repositories
	exchangeRepo := repository.NewExchangeRepository(db)
	rateRepo := repository.NewRateRepository(db, redisClient)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	rateService := services.NewRateService(rateRepo, cfg)
	exchangeService := services.NewExchangeService(exchangeRepo, rateService, mqChannel, cfg)
	fiatExchangeService := services.NewFiatExchangeService(exchangeRepo, rateService, mqChannel, cfg)
	tradingService := services.NewTradingService(orderRepo, exchangeService, cfg)

	// Start rate updater in background
	rateService.StartRateUpdater()
	log.Println("Rate updater started")

	// Initialize handlers
	exchangeHandler := handlers.NewExchangeHandler(exchangeService, rateService)
	fiatHandler := handlers.NewFiatHandler(fiatExchangeService, rateService)
	tradingHandler := handlers.NewTradingHandler(tradingService)

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Security middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiter())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "exchange-service",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := r.Group("/api/v1")

	// Public routes (no authentication required)
	public := api.Group("/")
	{
		// Exchange rates (public)
		public.GET("/rates", exchangeHandler.GetRates)
		public.GET("/rates/:from/:to", exchangeHandler.GetSpecificRate)
		public.GET("/rates/:from/:to/history", exchangeHandler.GetRateHistory)
		public.GET("/convert", exchangeHandler.ConvertCurrency)
		public.GET("/supported-currencies", exchangeHandler.GetSupportedCurrencies)

		// Fiat exchange rates (public)
		public.GET("/fiat/rates", fiatHandler.GetFiatRates)
		public.GET("/fiat/rates/:from/:to", fiatHandler.GetSpecificFiatRate)
		public.GET("/fiat/rates/:from/:to/history", fiatHandler.GetFiatRateHistory)
		public.GET("/fiat/convert", fiatHandler.FiatConverter)
		public.GET("/fiat/currencies", fiatHandler.GetSupportedFiatCurrencies)
		public.GET("/fiat/fees/compare", fiatHandler.CompareBankingFees)

		// Trading data (public)
		public.GET("/trading/tickers", tradingHandler.GetTickers)
		public.GET("/trading/orderbook/:pair", tradingHandler.GetOrderBook)
	}

	// Protected routes (authentication required)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		// Exchange operations
		protected.POST("/exchange/quote", exchangeHandler.GetQuote)
		protected.POST("/exchange/execute", exchangeHandler.ExecuteExchange)
		protected.GET("/exchange/history", exchangeHandler.GetExchangeHistory)
		protected.GET("/exchange/:id", exchangeHandler.GetExchange)

		// Fiat exchange operations
		protected.POST("/fiat/quote", fiatHandler.GetFiatQuote)
		protected.POST("/fiat/execute", fiatHandler.ExecuteFiatExchange)

		// Trading operations
		protected.POST("/trading/market-order", tradingHandler.PlaceMarketOrder)
		protected.POST("/trading/limit-order", tradingHandler.PlaceLimitOrder)
		protected.POST("/trading/stop-order", tradingHandler.PlaceStopOrder)
		protected.GET("/trading/orders", tradingHandler.GetUserOrders)
		protected.POST("/trading/orders/:id/cancel", tradingHandler.CancelOrder)
		protected.GET("/trading/portfolio", tradingHandler.GetPortfolio)
	}

	// Admin routes (admin authentication required)
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(cfg.JWTSecret))
	admin.Use(middleware.AdminOnly())
	{
		admin.POST("/rates/update", exchangeHandler.ForceUpdateRates)
		admin.GET("/exchange/all", exchangeHandler.GetAllExchanges)
		admin.GET("/orders/all", tradingHandler.GetUserOrders)
	}

	// Start server
	port := ":8083"
	log.Printf("Exchange service starting on port %s", port)
	log.Printf("Environment: %s", cfg.Environment)
	
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}