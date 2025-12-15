package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
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
	// Note: mqClient is the channel for wallet-service since InitializeRabbitMQ returns *amqp.Channel

	// Initialize repositories
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	cryptoService := services.NewCryptoService(cfg)
	balanceService := services.NewBalanceService(walletRepo, redisClient)
	walletService := services.NewWalletService(walletRepo, transactionRepo, cryptoService, balanceService, mqClient)
	
	// Start RabbitMQ consumer for inter-service communication
	consumer := services.NewConsumer(mqClient, walletService)
	if err := consumer.Start(); err != nil {
		log.Printf("Warning: Failed to start RabbitMQ consumer: %v", err)
	}
	
	// Initialize handlers
	walletHandler := handlers.NewWalletHandler(walletService, balanceService)

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
			"service": "wallet-service",
		})
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

			// Crypto wallet specific routes
			crypto := protected.Group("/crypto")
			{
				crypto.POST("/generate", walletHandler.GenerateCryptoWallet)
				crypto.GET("/:wallet_id/address", walletHandler.GetCryptoAddress)
				crypto.POST("/:wallet_id/send", walletHandler.SendCrypto)
				crypto.GET("/:wallet_id/pending", walletHandler.GetPendingTransactions)
				crypto.POST("/:wallet_id/estimate-fee", walletHandler.EstimateTransactionFee)
			}
		}

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