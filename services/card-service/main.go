package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/services"
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
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "card-service",
		})
	})

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