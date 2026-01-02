package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/services"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	cardOperationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "card_operations_total", Help: "Total card operations"}, []string{"type", "status"})
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" { c.Next(); return }
		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath(); if path == "" { path = c.Request.URL.Path }
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

	// Initialize message queue
	mqClient, err := database.InitializeRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer mqClient.Close()

	// Initialize repositories
	cardRepo := repository.NewCardRepository(db)
	transactionRepo := repository.NewCardTransactionRepository(db)

	// Initialize services
	cardIssuer := services.NewCardIssuerService(cfg)
	cardService := services.NewCardService(cardRepo, transactionRepo, cardIssuer, mqClient.GetChannel(), cfg)
	walletClient := services.NewWalletClient(cfg.WalletServiceURL)
	
	// Initialize handlers
	cardHandler := handlers.NewCardHandler(cardService, walletClient)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	// CORS configuration - Must be first!
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())
	router.Use(middleware.SecurityHeaders(), middleware.RateLimiter())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "card-service"}) })

	// Card routes
	api := router.Group("/api/v1")
	
	// Protected routes - use Group with middleware
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		// Card management
		protected.GET("/cards", cardHandler.GetUserCards)
		protected.POST("/cards", cardHandler.CreateCard)
		protected.GET("/cards/:card_id", cardHandler.GetCard)
		protected.PUT("/cards/:card_id", cardHandler.UpdateCard)
		protected.DELETE("/cards/:card_id", cardHandler.DeleteCard)

		// Card operations
		protected.POST("/cards/:card_id/activate", cardHandler.ActivateCard)
		protected.POST("/cards/:card_id/deactivate", cardHandler.DeactivateCard)
		protected.POST("/cards/:card_id/freeze", cardHandler.FreezeCard)
		protected.POST("/cards/:card_id/unfreeze", cardHandler.UnfreezeCard)
		protected.POST("/cards/:card_id/block", cardHandler.BlockCard)

		// Card loading (recharge)
		protected.POST("/cards/:card_id/load", cardHandler.LoadCard)
		protected.POST("/cards/:card_id/topup", cardHandler.LoadCard) // Alias for frontend compatibility
		protected.POST("/cards/:card_id/auto-load", cardHandler.SetupAutoLoad)
		protected.DELETE("/cards/:card_id/auto-load", cardHandler.CancelAutoLoad)

		// Card limits
		protected.GET("/cards/:card_id/limits", cardHandler.GetCardLimits)
		protected.PUT("/cards/:card_id/limits", cardHandler.UpdateCardLimits)

		// Card transactions
		protected.GET("/cards/:card_id/transactions", cardHandler.GetCardTransactions)
		protected.GET("/cards/:card_id/balance", cardHandler.GetCardBalance)

		// Card details and security
		protected.GET("/cards/:card_id/details", cardHandler.GetCardDetails)
		protected.POST("/cards/:card_id/pin", cardHandler.SetCardPIN)
		protected.PUT("/cards/:card_id/pin", cardHandler.ChangeCardPIN)
		protected.POST("/cards/:card_id/reset-pin", cardHandler.ResetCardPIN)

		// Virtual card specific
		virtual := protected.Group("/cards/virtual")
		{
			virtual.POST("/", cardHandler.CreateVirtualCard)
			virtual.POST("/:card_id/regenerate", cardHandler.RegenerateVirtualCard)
			virtual.GET("/:card_id/qr", cardHandler.GetCardQR)
		}

		// Physical card specific
		physical := protected.Group("/cards/physical")
		{
			physical.POST("/", cardHandler.OrderPhysicalCard)
			physical.GET("/:card_id/shipping", cardHandler.GetShippingStatus)
			physical.POST("/:card_id/activate-physical", cardHandler.ActivatePhysicalCard)
			physical.POST("/:card_id/report-lost", cardHandler.ReportLostCard)
			physical.POST("/:card_id/replacement", cardHandler.RequestReplacement)
		}

		// Gift cards
		gift := protected.Group("/cards/gift")
		{
			gift.POST("/", cardHandler.CreateGiftCard)
			gift.POST("/send", cardHandler.SendGiftCard)
			gift.POST("/redeem", cardHandler.RedeemGiftCard)
			gift.GET("/", cardHandler.GetGiftCards)
		}
	}

	// Webhook endpoints for card processors
	webhooks := api.Group("/webhooks")
	{
		webhooks.POST("/marqeta/transaction", cardHandler.HandleMarqetaTransaction)
		webhooks.POST("/marqeta/auth", cardHandler.HandleMarqetaAuth)
		webhooks.POST("/issuer/callback", cardHandler.HandleIssuerCallback)
	}

	// Public endpoints
	public := api.Group("/public")
	{
		public.GET("/supported-currencies", cardHandler.GetSupportedCurrencies)
		public.GET("/fees", cardHandler.GetCardFees)
		public.GET("/limits", cardHandler.GetCardLimits)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Card Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}