package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize admin database
	adminDB, err := database.InitializeAdminDB(cfg.AdminDBURL)
	if err != nil {
		log.Fatal("Failed to initialize admin database:", err)
	}
	defer adminDB.Close()

	// Initialize main database (read-only access)
	mainDB, err := database.InitializeMainDB(cfg.MainDBURL)
	if err != nil {
		log.Fatal("Failed to initialize main database:", err)
	}
	defer mainDB.Close()

	// Initialize RabbitMQ
	mqClient, err := database.InitializeRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ:", err)
	}
	defer mqClient.Close()

	// Initialize repository
	repo := repository.NewAdminRepository(adminDB, mainDB)

	// Initialize service
	adminService := services.NewAdminService(repo, mqClient, cfg)

	// Initialize storage service for presigned URLs
	storageService, err := services.NewStorageService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize Minio storage: %v (document viewing will fail)", err)
	}

	// Initialize handlers
	handler := handlers.NewAdminHandler(adminService)
	kycHandler := handlers.NewKYCHandler(storageService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	
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
			"service": "admin-service",
		})
	})

	// Public routes (no auth)
	public := router.Group("/api/v1/admin")
	{
		public.POST("/login", handler.Login)
	}

	// Protected routes (require admin auth)
	protected := router.Group("/api/v1/admin")
	protected.Use(middleware.AdminAuth(cfg.AdminJWTSecret))
	{
		// Current admin info
		protected.GET("/me", handler.GetCurrentAdmin)

		// Dashboard
		protected.GET("/dashboard", handler.GetDashboard)

		// Admin management
		admins := protected.Group("/admins")
		admins.Use(middleware.RequirePermission(models.PermViewAdmins))
		{
			admins.GET("", handler.GetAllAdmins)
			admins.GET("/:id", handler.GetAdmin)
			admins.POST("", middleware.RequirePermission(models.PermCreateAdmins), handler.CreateAdmin)
			admins.PUT("/:id", middleware.RequirePermission(models.PermUpdateAdmins), handler.UpdateAdmin)
			admins.DELETE("/:id", middleware.RequirePermission(models.PermDeleteAdmins), handler.DeleteAdmin)
		}

		// Roles
		roles := protected.Group("/roles")
		roles.Use(middleware.RequirePermission(models.PermViewAdmins))
		{
			roles.GET("", handler.GetRoles)
			roles.GET("/:id", handler.GetRole)
		}

		// Users management
		users := protected.Group("/users")
		users.Use(middleware.RequirePermission(models.PermViewUsers))
		{
			users.GET("", handler.GetUsers)
			users.GET("/:id/kyc/documents", middleware.RequirePermission(models.PermViewKYC), handler.GetUserKYCDocuments)
			users.POST("/:id/block", middleware.RequirePermission(models.PermBlockUsers), handler.BlockUser)
			users.POST("/:id/unblock", middleware.RequirePermission(models.PermBlockUsers), handler.UnblockUser)
		}

		// KYC management
		kyc := protected.Group("/kyc")
		kyc.Use(middleware.RequirePermission(models.PermViewKYC))
		{
			kyc.POST("/:id/approve", middleware.RequirePermission(models.PermApproveKYC), handler.ApproveKYC)
			kyc.POST("/:id/reject", middleware.RequirePermission(models.PermRejectKYC), handler.RejectKYC)
			// Secure document access - generates presigned URLs
			kyc.POST("/document-url", kycHandler.GetDocumentURL)
			kyc.POST("/download-url", kycHandler.GetDocumentDownloadURL)
		}

		// Transactions
		transactions := protected.Group("/transactions")
		transactions.Use(middleware.RequirePermission(models.PermViewTransactions))
		{
			transactions.GET("", handler.GetTransactions)
			transactions.POST("/:id/block", middleware.RequirePermission(models.PermBlockTransactions), handler.BlockTransaction)
			transactions.POST("/:id/refund", middleware.RequirePermission(models.PermRefundTransactions), handler.RefundTransaction)
		}

		// Cards
		cards := protected.Group("/cards")
		cards.Use(middleware.RequirePermission(models.PermViewCards))
		{
			cards.GET("", handler.GetCards)
			cards.POST("/:id/freeze", middleware.RequirePermission(models.PermFreezeCards), handler.FreezeCard)
			cards.POST("/:id/block", middleware.RequirePermission(models.PermBlockCards), handler.BlockCard)
		}

		// Wallets
		wallets := protected.Group("/wallets")
		wallets.Use(middleware.RequirePermission(models.PermViewWallets))
		{
			wallets.GET("", handler.GetWallets)
			wallets.POST("/:id/freeze", middleware.RequirePermission(models.PermFreezeWallets), handler.FreezeWallet)
		}

		// Audit logs
		logs := protected.Group("/logs")
		logs.Use(middleware.RequirePermission(models.PermViewLogs))
		{
			logs.GET("", handler.GetAuditLogs)
		}

		// Payment Providers management
		paymentHandler := handlers.NewPaymentHandler(mainDB)
		payments := protected.Group("/payment-providers")
		{
			payments.GET("", paymentHandler.GetPaymentProviders)
			payments.GET("/:id", paymentHandler.GetPaymentProvider)
			payments.POST("", paymentHandler.CreatePaymentProvider)
			payments.PUT("/:id", paymentHandler.UpdatePaymentProvider)
			payments.DELETE("/:id", paymentHandler.DeletePaymentProvider)
			payments.POST("/:id/toggle-status", paymentHandler.ToggleProviderStatus)
			payments.POST("/:id/toggle-demo", paymentHandler.ToggleDemoMode)
			payments.POST("/:id/test", paymentHandler.TestProviderConnection)
			payments.POST("/:id/countries", paymentHandler.AddProviderCountry)
			payments.DELETE("/:id/countries/:country", paymentHandler.RemoveProviderCountry)
		}
	}

	// Public endpoint for payment methods (no auth required for wallet)
	publicPayments := router.Group("/api/v1")
	{
		paymentHandler := handlers.NewPaymentHandler(mainDB)
		publicPayments.GET("/payment-methods", paymentHandler.GetPaymentMethodsForCountry)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	log.Printf("Admin Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
