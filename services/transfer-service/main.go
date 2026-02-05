package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/providers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/service"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal   = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	transfersTotal      = promauto.NewCounterVec(prometheus.CounterOpts{Name: "transfers_total", Help: "Total transfers"}, []string{"type", "status"})
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

	// Initialize Kafka client
	kafkaClient, err := database.InitializeKafka(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatal("Failed to initialize Kafka:", err)
	}
	defer kafkaClient.Close()

	// Initialize repositories
	transferRepo := repository.NewTransferRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	feeRepo := repository.NewFeeRepository(db)
	depositRepo := repository.NewDepositRepository(db)
	payoutRepo := repository.NewPayoutRepository(db)

	// Initialize payout schema
	if err := payoutRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to init payout schema: %v", err)
	}

	// Load provider configuration (Vault/Env)
	providerCfg := providers.LoadConfig()

	// Initialize services
	walletClient := services.NewWalletClient(cfg)
	exchangeClient := services.NewExchangeClient()
	enterpriseClient := services.NewEnterpriseClient(cfg)
	shopClient := services.NewShopClient(cfg)
	feeService := services.NewFeeService(feeRepo)
	transferService := services.NewTransferService(transferRepo, walletRepo, kafkaClient, enterpriseClient, shopClient, feeService, cfg)

	// Initialize Mobile Money Service with Router
	mobilemoneyService := services.NewMobileMoneyService(cfg, providerCfg)

	internationalService := services.NewInternationalTransferService(cfg)
	complianceService := services.NewComplianceService(cfg)

	// Start Kafka Consumers
	paymentConsumer := services.NewPaymentRequestConsumer(kafkaClient, walletClient, exchangeClient, walletRepo, transferRepo)
	if err := paymentConsumer.Start(); err != nil {
		log.Printf("Warning: Failed to start PaymentRequestConsumer: %v", err)
	}

	// Initialize handlers
	transferHandler := handlers.NewTransferHandler(
		transferService,
		mobilemoneyService,
		internationalService,
		complianceService,
	)

	// Initialize Admin Client (Proxy for Payment Methods)
	adminClient := services.NewAdminClient(cfg)

	// Initialize Aggregator (Payment Providers)
	aggregatorRepo := repository.NewAggregatorRepository(db)
	if err := aggregatorRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to init aggregator schema: %v", err)
	}
	aggregatorHandler := handlers.NewAggregatorHandler(aggregatorRepo, adminClient)

	// Initialize Instance-Based Provider System for Deposits
	instanceRepo := repository.NewAggregatorInstanceRepository(db)
	providerLoader := providers.NewInstanceBasedProviderLoader(instanceRepo)
	fundMovementService := service.NewDepositFundMovementService(db)
	walletServiceURL := cfg.WalletServiceURL
	if walletServiceURL == "" {
		walletServiceURL = "http://wallet-service:8083"
	}

	// Webhook secrets for signature verification (from env/vault)
	webhookSecrets := map[string]string{
		"flutterwave": os.Getenv("FLUTTERWAVE_WEBHOOK_SECRET"),
		"stripe":      os.Getenv("STRIPE_WEBHOOK_SECRET"),
		"paystack":    os.Getenv("PAYSTACK_WEBHOOK_SECRET"),
		"cinetpay":    os.Getenv("CINETPAY_WEBHOOK_SECRET"),
		"paypal":      os.Getenv("PAYPAL_WEBHOOK_SECRET"),
	}

	// Initialize Full Deposit Handler (production-ready)
	depositHandler := handlers.NewDepositHandler(
		depositRepo,
		instanceRepo,
		providerLoader,
		fundMovementService,
		walletServiceURL,
		webhookSecrets,
	)

	// Initialize Withdrawal/Payout Handler
	withdrawalFundMovement := service.NewWithdrawalFundMovementService(db)
	payoutHandler := handlers.NewPayoutHandler(
		payoutRepo,
		instanceRepo,
		providerLoader,
		withdrawalFundMovement,
		walletServiceURL,
		webhookSecrets,
	)

	// Initialize and start deposit expiry service (background job)
	depositExpiryService := service.NewDepositExpiryService(depositRepo, 5*time.Minute)
	depositExpiryService.Start()
	defer depositExpiryService.Stop()

	// Initialize Admin Instance Handler for pause/update functionality
	adminInstanceHandler := handlers.NewAdminInstanceHandler(instanceRepo, aggregatorRepo)

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
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "transfer-service"}) })

	// Transfer routes
	api := router.Group("/api/v1")
	{
		// Register Aggregator Routes (Proxy Pattern)
		// Public routes: /api/v1/aggregators
		// Admin routes: /api/v1/admin/aggregators (Protected)
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.JWTAuth(cfg.JWTSecret)) // Apply Auth to admin routes
		aggregatorHandler.RegisterRoutes(api, adminGroup)

		// Admin Instance Management Routes
		instances := adminGroup.Group("/instances")
		{
			instances.GET("", adminInstanceHandler.ListInstances)
			instances.GET("/:id", adminInstanceHandler.GetInstance)
			instances.POST("", adminInstanceHandler.CreateInstance)
			instances.PUT("/:id", adminInstanceHandler.UpdateInstance)
			instances.POST("/:id/pause", adminInstanceHandler.TogglePause)
			instances.GET("/:id/stats", adminInstanceHandler.GetInstanceStats)
			instances.POST("/:id/wallets", adminInstanceHandler.LinkWallet)
			instances.DELETE("/:id/wallets/:wallet_id", adminInstanceHandler.UnlinkWallet)
			instances.PUT("/:id/wallets/:wallet_id", adminInstanceHandler.UpdateInstanceWallet)
		}

		// Protected routes - apply JWT auth middleware
		api.POST("/transfers", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CreateTransfer)
		api.GET("/transfers", middleware.JWTAuth(cfg.JWTSecret), transferHandler.GetTransferHistory)

		// Additional lookup routes (consistent with frontend useApi.ts)
		api.GET("/transfers/banks", transferHandler.GetBanks)
		api.GET("/transfers/mobile-operators", transferHandler.GetMobileProviders)
		api.POST("/transfers/validate-recipient", transferHandler.ValidateRecipient)
		api.GET("/transfers/fees", transferHandler.GetFees)

		api.GET("/transfers/:transfer_id", middleware.JWTAuth(cfg.JWTSecret), transferHandler.GetTransfer)
		api.POST("/transfers/:transfer_id/cancel", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CancelTransfer)
		api.POST("/transfers/:transfer_id/reverse", middleware.JWTAuth(cfg.JWTSecret), transferHandler.ReverseTransfer)

		// International transfers
		api.POST("/international", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CreateInternationalTransfer)

		// Mobile money
		mobile := api.Group("/mobile")
		mobile.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			mobile.POST("/send", transferHandler.SendMobileMoney)
			mobile.POST("/receive", transferHandler.ReceiveMobileMoney)
			mobile.GET("/providers", transferHandler.GetMobileProviders) // Keep for backward compatibility if needed
		}

		// Bulk transfers
		bulk := api.Group("/bulk")
		bulk.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			bulk.POST("/", transferHandler.CreateBulkTransfer)
			bulk.GET("/:batch_id", transferHandler.GetBulkTransferStatus)
			bulk.POST("/:batch_id/approve", transferHandler.ApproveBulkTransfer)
		}

		// Webhook endpoints
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/mobile/callback", transferHandler.HandleMobileMoneyCallback)
			webhooks.POST("/bank/callback", transferHandler.HandleBankCallback)
		}

		// Deposit/Collection routes
		deposits := api.Group("/deposits")
		{
			deposits.POST("/initiate", middleware.JWTAuth(cfg.JWTSecret), depositHandler.InitiateDeposit)
			deposits.GET("/:id/status", middleware.JWTAuth(cfg.JWTSecret), depositHandler.GetDepositStatus)
			deposits.POST("/:id/cancel", middleware.JWTAuth(cfg.JWTSecret), depositHandler.CancelDeposit)
			deposits.GET("/user/:user_id", middleware.JWTAuth(cfg.JWTSecret), depositHandler.GetUserDeposits)

			// Webhooks (no auth - verified by signature)
			deposits.POST("/webhook/:provider", depositHandler.HandleWebhook)
		}

		// Payout/Withdrawal routes (external transfers: Mobile Money, Bank, PayPal)
		payouts := api.Group("/payouts")
		{
			payouts.POST("/quote", middleware.JWTAuth(cfg.JWTSecret), payoutHandler.GetPayoutQuote)
			payouts.POST("/initiate", middleware.JWTAuth(cfg.JWTSecret), payoutHandler.InitiatePayout)
			payouts.GET("/:id/status", middleware.JWTAuth(cfg.JWTSecret), payoutHandler.GetPayoutStatus)
			payouts.POST("/:id/cancel", middleware.JWTAuth(cfg.JWTSecret), payoutHandler.CancelPayout)
			payouts.GET("/user/:user_id", middleware.JWTAuth(cfg.JWTSecret), payoutHandler.GetUserPayouts)
			payouts.GET("/banks", payoutHandler.GetBanks)
			payouts.GET("/mobile-operators", payoutHandler.GetMobileOperators)

			// Payout webhooks (no auth - verified by signature)
			payouts.POST("/webhook/:provider", payoutHandler.HandleWebhook)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Transfer Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
