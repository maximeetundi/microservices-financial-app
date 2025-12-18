package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

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
		// Protected routes - apply JWT auth middleware
		api.POST("/transfers", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CreateTransfer)
		api.GET("/transfers", middleware.JWTAuth(cfg.JWTSecret), transferHandler.GetTransferHistory)
		api.GET("/transfers/:transfer_id", middleware.JWTAuth(cfg.JWTSecret), transferHandler.GetTransfer)
		api.POST("/transfers/:transfer_id/cancel", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CancelTransfer)

		// International transfers
		api.POST("/international", middleware.JWTAuth(cfg.JWTSecret), transferHandler.CreateInternationalTransfer)

		// Mobile money
		mobile := api.Group("/mobile")
		mobile.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			mobile.POST("/send", transferHandler.SendMobileMoney)
			mobile.POST("/receive", transferHandler.ReceiveMobileMoney)
			mobile.GET("/providers", transferHandler.GetMobileProviders)
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