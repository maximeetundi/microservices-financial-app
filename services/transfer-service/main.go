package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/services"
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

	// Initialize message queue
	mqClient, err := database.InitializeRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer mqClient.Close()

	// Initialize repositories
	transferRepo := repository.NewTransferRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// Initialize services
	transferService := services.NewTransferService(transferRepo, walletRepo, mqClient, cfg)
	mobilemoneyService := services.NewMobileMoneyService(cfg)
	internationalService := services.NewInternationalTransferService(cfg)
	complianceService := services.NewComplianceService(cfg)
	
	// Initialize handlers
	transferHandler := handlers.NewTransferHandler(
		transferService, 
		mobilemoneyService, 
		internationalService,
		complianceService,
	)

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
			"service": "transfer-service",
		})
	})

	// Transfer routes
	api := router.Group("/api/v1")
	{
		// Protected routes
		protected := api.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.POST("/transfers", transferHandler.CreateTransfer)
			protected.GET("/transfers", transferHandler.GetTransferHistory)
			protected.GET("/transfers/:transfer_id", transferHandler.GetTransfer)
			protected.POST("/transfers/:transfer_id/cancel", transferHandler.CancelTransfer)

			// International transfers
			protected.POST("/international", transferHandler.CreateInternationalTransfer)

			// Mobile money
			mobile := protected.Group("/mobile")
			{
				mobile.POST("/send", transferHandler.SendMobileMoney)
				mobile.POST("/receive", transferHandler.ReceiveMobileMoney)
				mobile.GET("/providers", transferHandler.GetMobileProviders)
			}

			// Bulk transfers
			bulk := protected.Group("/bulk")
			{
				bulk.POST("/", transferHandler.CreateBulkTransfer)
				bulk.GET("/:batch_id", transferHandler.GetBulkTransferStatus)
				bulk.POST("/:batch_id/approve", transferHandler.ApproveBulkTransfer)
			}
		}

		// Webhook endpoints
		webhooks := api.Group("/webhooks")
		{
			webhooks.POST("/mobile/callback", transferHandler.HandleMobileMoneyCallback)
			webhooks.POST("/bank/callback", transferHandler.HandleBankCallback)
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