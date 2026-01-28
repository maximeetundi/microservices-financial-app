package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/exchange-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal   = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	exchangesTotal      = promauto.NewCounterVec(prometheus.CounterOpts{Name: "exchanges_total", Help: "Total exchanges"}, []string{"pair", "status"})
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(time.Since(start).Seconds())
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.InitializeRedis(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize Kafka
	kafkaClient, err := database.InitializeKafka(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatal("Failed to initialize Kafka:", err)
	}
	defer kafkaClient.Close()

	// Initialize repositories
	exchangeRepo := repository.NewExchangeRepository(db)
	rateRepo := repository.NewRateRepository(db, redisClient)
	orderRepo := repository.NewOrderRepository(db)
	feeRepo := repository.NewFeeRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	// Initialize tables
	if err := feeRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize fee tables: %v", err)
	}
	if err := settingsRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize settings tables: %v", err)
	}
	if err := exchangeRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize exchange tables: %v", err)
	}

	// Initialize Kafka publisher wrapper
	kafkaPublisher := services.NewKafkaPublisher(kafkaClient)

	// Initialize services
	walletClient := services.NewWalletClient(cfg.WalletServiceURL)
	settingsService := services.NewSettingsService(settingsRepo)
	rateService := services.NewRateServiceWithSettings(rateRepo, cfg, settingsService)
	feeService := services.NewFeeService(feeRepo)

	// Start background rate updater (fetches prices from CoinGecko/Binance and stores in DB)
	rateService.StartRateUpdater()
	log.Printf("Background rate updater started (interval: %ds)", settingsService.GetRateUpdateInterval())
	exchangeService := services.NewExchangeService(exchangeRepo, orderRepo, rateService, feeService, kafkaPublisher, walletClient, cfg)
	tradingService := services.NewTradingService(orderRepo, exchangeService, cfg)
	fiatExchangeService := services.NewFiatExchangeService(exchangeRepo, rateService, feeService, kafkaPublisher, cfg)

	// Initialize handlers
	exchangeHandler := handlers.NewExchangeHandler(exchangeService, rateService, tradingService, walletClient)
	adminFeeHandler := handlers.NewAdminFeeHandler(feeService)
	adminSettingsHandler := handlers.NewAdminSettingsHandler(settingsService)
	fiatHandler := handlers.NewFiatHandler(fiatExchangeService, rateService)

	// Start Kafka Consumers
	paymentConsumer := services.NewPaymentStatusConsumer(kafkaClient, exchangeService)
	paymentConsumer.SetFiatExchangeService(fiatExchangeService)
	if err := paymentConsumer.Start(); err != nil {
		log.Printf("Warning: Failed to start PaymentStatusConsumer: %v", err)
	}

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())
	router.Use(middleware.SecurityHeaders(), middleware.RateLimiter())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "exchange-service"}) })

	// Exchange routes
	api := router.Group("/api/v1")
	{
		// Public routes (no auth required)
		api.GET("/rates", exchangeHandler.GetExchangeRates)
		api.GET("/rates/:from/:to", exchangeHandler.GetSpecificRate)
		api.GET("/markets", exchangeHandler.GetMarkets)
		api.GET("/orderbook/:pair", exchangeHandler.GetOrderBook)
		api.GET("/trades/:pair", exchangeHandler.GetRecentTrades)

		// Protected routes - using Group instead of Use to get RouterGroup
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// Exchange operations
			protected.POST("/quote", exchangeHandler.GetQuote)
			protected.POST("/execute", exchangeHandler.ExecuteExchange)
			protected.GET("/history", exchangeHandler.GetExchangeHistory)
			protected.GET("/exchange/:exchange_id", exchangeHandler.GetExchange)

			// Trading operations
			trading := protected.Group("/trading")
			{
				trading.POST("/buy", exchangeHandler.PlaceBuyOrder)
				trading.POST("/sell", exchangeHandler.PlaceSellOrder)
				trading.POST("/limit-order", exchangeHandler.CreateLimitOrder)
				trading.POST("/stop-loss", exchangeHandler.CreateStopLoss)
				trading.GET("/orders", exchangeHandler.GetUserOrders)
				trading.GET("/orders/active", exchangeHandler.GetActiveOrders)
				trading.DELETE("/orders/:order_id", exchangeHandler.CancelOrder)
				trading.GET("/portfolio", exchangeHandler.GetPortfolio)
				trading.GET("/performance", exchangeHandler.GetPerformance)
			}

			// P2P Trading
			p2p := protected.Group("/p2p")
			{
				p2p.POST("/offers", exchangeHandler.CreateP2POffer)
				p2p.GET("/offers", exchangeHandler.GetP2POffers)
				p2p.POST("/offers/:offer_id/accept", exchangeHandler.AcceptP2POffer)
				p2p.GET("/trades", exchangeHandler.GetP2PTrades)
			}

			// Fiat Exchange (Fiat-to-Fiat conversions)
			fiat := protected.Group("/fiat")
			{
				fiat.POST("/quote", fiatHandler.GetFiatQuote)
				fiat.POST("/exchange", fiatHandler.ExecuteFiatExchange)
				fiat.GET("/rates", fiatHandler.GetFiatRates)
				fiat.GET("/rates/:from/:to", fiatHandler.GetSpecificFiatRate)
				fiat.GET("/currencies", fiatHandler.GetSupportedFiatCurrencies)
				fiat.GET("/convert", fiatHandler.FiatConverter)
				fiat.GET("/compare-fees", fiatHandler.CompareBankingFees)
			}
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/rates/update", exchangeHandler.UpdateRates)
			admin.GET("/analytics", exchangeHandler.GetAnalytics)
			admin.GET("/volume", exchangeHandler.GetTradingVolume)

			// Fee Configuration
			admin.GET("/fees", adminFeeHandler.GetFees)
			admin.PUT("/fees", adminFeeHandler.UpdateFee)

			// System Settings
			admin.GET("/settings", adminSettingsHandler.GetSettings)
			admin.PUT("/settings", adminSettingsHandler.UpdateSetting)
			admin.GET("/settings/:key", adminSettingsHandler.GetSettingByKey)
		}

		// Webhook endpoints for price feeds
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/price-update", exchangeHandler.HandlePriceUpdate)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Exchange Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
