package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/exchange-service/internal/config"
	"github.com/crypto-bank/exchange-service/internal/database"
	"github.com/crypto-bank/exchange-service/internal/handlers"
	"github.com/crypto-bank/exchange-service/internal/middleware"
	"github.com/crypto-bank/exchange-service/internal/repository"
	"github.com/crypto-bank/exchange-service/internal/services"
)

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

	// Initialize message queue
	mqClient, err := database.InitializeRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer mqClient.Close()

	// Initialize repositories
	exchangeRepo := repository.NewExchangeRepository(db)
	rateRepo := repository.NewRateRepository(db, redisClient)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	rateService := services.NewRateService(rateRepo, cfg)
	exchangeService := services.NewExchangeService(exchangeRepo, rateService, mqClient, cfg)
	tradingService := services.NewTradingService(orderRepo, exchangeService, cfg)
	walletClient := services.NewWalletClient(cfg.WalletServiceURL)
	
	// Initialize handlers
	exchangeHandler := handlers.NewExchangeHandler(exchangeService, rateService, tradingService, walletClient)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "exchange-service",
		})
	})

	// Exchange routes
	api := router.Group("/api/v1")
	{
		// Public routes (no auth required)
		api.GET("/rates", exchangeHandler.GetExchangeRates)
		api.GET("/rates/:from/:to", exchangeHandler.GetSpecificRate)
		api.GET("/markets", exchangeHandler.GetMarkets)
		api.GET("/orderbook/:pair", exchangeHandler.GetOrderBook)
		api.GET("/trades/:pair", exchangeHandler.GetRecentTrades)

		// Protected routes
		protected := api.Use(middleware.JWTAuth(cfg.JWTSecret))
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
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/rates/update", exchangeHandler.UpdateRates)
			admin.GET("/analytics", exchangeHandler.GetAnalytics)
			admin.GET("/volume", exchangeHandler.GetTradingVolume)
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