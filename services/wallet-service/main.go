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
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
)

// Prometheus metrics
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}},
		[]string{"method", "path", "status"},
	)
	transactionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "wallet_transactions_total", Help: "Total wallet transactions"},
		[]string{"type", "status"},
	)
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" { c.Next(); return }
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath(); if path == "" { path = c.Request.URL.Path }
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
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

	// Initialize Kafka client
	kafkaClient, err := database.InitializeKafka(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatal("Failed to initialize Kafka:", err)
	}
	defer kafkaClient.Close()

	// Initialize repositories
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	paymentRepo := repository.NewPaymentRequestRepository(db)
	feeRepo := repository.NewFeeRepository(db)
	
	// Initialize payment tables
	if err := paymentRepo.InitTable(); err != nil {
		log.Printf("Warning: Failed to initialize payment tables: %v", err)
	}

	// Initialize fee tables
	if err := feeRepo.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize fee tables: %v", err)
	}

	// Initialize services
	feeService := services.NewFeeService(feeRepo)
	cryptoService := services.NewCryptoService(cfg)
	exchangeClient := services.NewExchangeClient()
	balanceService := services.NewBalanceService(walletRepo, redisClient, exchangeClient, kafkaClient)
	walletService := services.NewWalletService(walletRepo, transactionRepo, cryptoService, balanceService, feeService, kafkaClient)
	merchantService := services.NewMerchantPaymentService(paymentRepo, walletService, feeService, cfg, kafkaClient)
	
	// Start Kafka consumer for inter-service communication
	consumer := services.NewConsumer(kafkaClient, walletService)
	if err := consumer.Start(); err != nil {
		log.Printf("Warning: Failed to start Kafka consumer: %v", err)
	}
	
	// Initialize handlers
	walletHandler := handlers.NewWalletHandler(walletService, balanceService)
	merchantHandler := handlers.NewMerchantPaymentHandler(merchantService)
	adminFeeHandler := handlers.NewAdminFeeHandler(feeService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(prometheusMiddleware())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "wallet-service"})
	})

	// Wallet routes
	api := router.Group("/api/v1")
	{
		// Protected routes - use Group to get RouterGroup type
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/wallets", walletHandler.GetWallets)
			protected.POST("/wallets", walletHandler.CreateWallet)
			protected.GET("/wallets/:wallet_id", walletHandler.GetWallet)
			protected.PUT("/wallets/:wallet_id", walletHandler.UpdateWallet)
			protected.GET("/wallets/:wallet_id/balance", walletHandler.GetBalance)
			protected.GET("/wallets/:wallet_id/transactions", walletHandler.GetWalletTransactions)
			protected.POST("/wallets/:wallet_id/freeze", walletHandler.FreezeWallet)
			protected.POST("/wallets/:wallet_id/unfreeze", walletHandler.UnfreezeWallet)
			
			// Deposit and Withdraw routes (added to match frontend)
			protected.POST("/wallets/:wallet_id/deposit", walletHandler.Deposit)
			protected.POST("/wallets/:wallet_id/withdraw", walletHandler.Withdraw)

			// Crypto wallet specific routes
			crypto := protected.Group("/crypto")
			{
				crypto.POST("/generate", walletHandler.GenerateCryptoWallet)
				crypto.GET("/:wallet_id/address", walletHandler.GetCryptoAddress)
				crypto.POST("/:wallet_id/send", walletHandler.SendCrypto)
				crypto.GET("/:wallet_id/pending", walletHandler.GetPendingTransactions)
				crypto.POST("/:wallet_id/estimate-fee", walletHandler.EstimateTransactionFee)
			}

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", walletHandler.GetDashboardSummary)
				dashboard.GET("/activity", walletHandler.GetRecentActivity)
				dashboard.GET("/portfolio", walletHandler.GetPortfolio)
				dashboard.GET("/stats", walletHandler.GetStats)
			}

			// Merchant payment routes
			merchant := protected.Group("/merchant")
			{
				merchant.POST("/payments", merchantHandler.CreatePaymentRequest)
				merchant.GET("/payments", merchantHandler.GetMerchantPayments)
				merchant.GET("/payments/history", merchantHandler.GetMerchantHistory)
				merchant.DELETE("/payments/:id", merchantHandler.CancelPaymentRequest)
				merchant.POST("/quick-pay", merchantHandler.QuickPaymentRequest)
			}

			// Payment routes (for paying)
			payments := protected.Group("/payments")
			{
				payments.GET("/:id", merchantHandler.GetPaymentRequest)
				payments.GET("/:id/qr", merchantHandler.GetPaymentQRCode)
				payments.POST("/:id/pay", merchantHandler.PayPaymentRequest)
			}

			// Admin Fee Management
			admin := protected.Group("/admin")
			{
				admin.GET("/fees", adminFeeHandler.GetFees)
				admin.PUT("/fees", adminFeeHandler.UpdateFee)
			}
		}

		// Service-to-Service routes (should be protected by internal network or secret)
		api.POST("/wallets/transaction", walletHandler.ProcessInterServiceTransaction)
		
		// Internal Wallet Management (for other services like transfer-service)
		internal := api.Group("/internal")
		{
			internal.GET("/wallets", walletHandler.GetWalletsInternal)
			internal.POST("/wallets", walletHandler.CreateWalletInternal)
		}

		// Public payment scan endpoint (no auth required)
		api.GET("/pay/:id", merchantHandler.ScanPayment)

		// Webhook endpoints for blockchain notifications
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/crypto/confirmation", walletHandler.HandleCryptoConfirmation)
			webhooks.POST("/crypto/deposit", walletHandler.HandleCryptoDeposit)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Wallet Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}