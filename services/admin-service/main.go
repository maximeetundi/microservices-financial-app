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

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/services"
)

var (
	httpRequestsTotal   = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	adminActionsTotal   = promauto.NewCounterVec(prometheus.CounterOpts{Name: "admin_actions_total", Help: "Total admin actions"}, []string{"action", "admin"})
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

	// Migrate Aggregator Schema (tables & views for transfer-service)
	if err := database.MigrateAggregatorSchema(mainDB); err != nil {
		log.Printf("Warning: failed to migrate aggregator schema: %v", err)
	}

	// Seed provider instances (needs access to both Admin and Main DBs)
	if err := database.SeedProviderInstances(adminDB, mainDB); err != nil {
		log.Printf("Warning: failed to seed provider instances: %v", err)
	}

	// Initialize MongoDB
	mongoClient, err := database.InitializeMongoDB(cfg.MongoDBURI)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	defer database.CloseMongoDB(mongoClient)

	// Initialize Kafka
	kafkaClient, err := database.InitializeKafka(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatal("Failed to initialize Kafka:", err)
	}
	defer kafkaClient.Close()

	// Initialize repository
	repo := repository.NewAdminRepository(adminDB, mainDB)
	mongoRepo := repository.NewMongoRepository(mongoClient)

	// Initialize service
	adminService := services.NewAdminService(repo, mongoRepo, kafkaClient, cfg)

	// Initialize system settings (fees, etc)
	if err := adminService.InitializeSettings(); err != nil {
		log.Printf("Warning: Failed to initialize settings: %v", err)
	}

	// Initialize storage service for presigned URLs
	storageService, err := services.NewStorageService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
		cfg.MinioPublicURL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize Minio storage: %v (document viewing will fail)", err)
	}

	// Initialize handlers
	handler := handlers.NewAdminHandler(adminService)
	kycHandler := handlers.NewKYCHandler(storageService)
	notifHandler := handlers.NewNotificationHandler(adminDB)

	// Start Kafka event consumer for admin notifications
	eventConsumer := services.NewEventConsumer(kafkaClient, adminDB)
	if err := eventConsumer.StartConsuming(); err != nil {
		log.Printf("Warning: Failed to start Kafka event consumer: %v", err)
	}

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())
	router.Use(middleware.SecurityHeaders())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "admin-service"}) })

	// Public routes (no auth)
	public := router.Group("/api/v1/admin")
	{
		public.POST("/login", handler.Login)
		// Public payment methods endpoint (moved here to ensure reachability via gateway)
		// Note: Use mainDB because payment_providers table is in the main database, not admin database
		publicPaymentHandler := handlers.NewPaymentHandler(mainDB)
		public.GET("/payment-methods", publicPaymentHandler.GetPaymentMethodsForCountry)
	}

	// Platform Proxy Handler
	proxyHandler := handlers.NewPlatformProxyHandler(cfg)

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
			kyc.GET("/requests", handler.GetAllKYCRequests)
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

		// Admin Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notifHandler.GetNotifications)
			notifications.GET("/unread-count", notifHandler.GetUnreadCount)
			notifications.POST("/:id/read", notifHandler.MarkAsRead)
			notifications.POST("/mark-all-read", notifHandler.MarkAllAsRead)
			notifications.POST("", notifHandler.CreateNotification)
			notifications.DELETE("/cleanup", notifHandler.DeleteOldNotifications)
		}

		// Fee Configuration
		fees := protected.Group("/fees")
		fees.Use(middleware.RequirePermission(models.PermManageSettings))
		{
			fees.GET("", handler.GetFeeConfigs)
			fees.POST("", handler.CreateFeeConfig)
			fees.PUT("/:key", handler.UpdateFeeConfig)
		}

		// System Settings (Proxied to Exchange Service)
		settingsProxy := handlers.NewSettingsProxyHandler(cfg)
		settings := protected.Group("/settings")
		settings.Use(middleware.RequirePermission(models.PermManageSettings))
		{
			settings.GET("", settingsProxy.ProxyRequest)
			settings.PUT("", settingsProxy.ProxyRequest)
			settings.GET("/:key", settingsProxy.ProxyRequest)
		}

		// Donations & Campaigns
		protected.GET("/campaigns", handler.GetCampaigns)
		protected.GET("/donations", handler.GetDonations)

		// Events & Tickets
		protected.GET("/events", handler.GetEvents)
		protected.GET("/sold-tickets", handler.GetSoldTickets)

		// Payment Providers management
		paymentHandler := handlers.NewPaymentHandler(adminDB)
		payments := protected.Group("/payment-providers")
		{
			payments.GET("", paymentHandler.GetPaymentProviders)
			payments.GET("/:id", paymentHandler.GetPaymentProvider)
			payments.POST("", paymentHandler.CreatePaymentProvider)
			payments.PUT("/:id", paymentHandler.UpdatePaymentProvider)
			payments.DELETE("/:id", paymentHandler.DeletePaymentProvider)
			payments.POST("/:id/toggle-status", paymentHandler.ToggleProviderStatus)
			payments.POST("/:id/toggle-demo", paymentHandler.ToggleDemoMode)
			payments.POST("/:id/toggle-deposit", paymentHandler.ToggleDepositEnabled)
			payments.POST("/:id/toggle-withdraw", paymentHandler.ToggleWithdrawEnabled)
			payments.PUT("/:id/limits", paymentHandler.UpdateProviderLimits)
			payments.PUT("/:id/fees", paymentHandler.UpdateProviderFees)
			payments.PUT("/:id/priority", paymentHandler.UpdateProviderPriority)
			payments.POST("/:id/test", paymentHandler.TestProviderConnection)
			payments.POST("/:id/countries", paymentHandler.AddProviderCountry)
			payments.POST("/:id/countries/:country/toggle", paymentHandler.ToggleCountryStatus)
			payments.DELETE("/:id/countries/:country", paymentHandler.RemoveProviderCountry)
		}

		// Provider Instances management (multi-key support)
		instanceHandler := handlers.NewInstanceHandler(adminDB)
		instances := protected.Group("/provider-instances")
		{
			instances.GET("", instanceHandler.GetAllInstances)
		}
		// Nested under providers
		payments.GET("/:id/instances", instanceHandler.GetProviderInstances)
		payments.POST("/:id/instances", instanceHandler.CreateProviderInstance)
		payments.GET("/:id/instances/best", instanceHandler.SelectBestInstance)
		payments.PUT("/:id/instances/:instanceId", instanceHandler.UpdateProviderInstance)
		payments.DELETE("/:id/instances/:instanceId", instanceHandler.DeleteProviderInstance)
		payments.POST("/:id/instances/:instanceId/link-wallet", instanceHandler.LinkHotWallet)
		payments.POST("/:id/instances/:instanceId/test", instanceHandler.TestInstance)
		payments.POST("/:id/instances/:instanceId/pause", instanceHandler.ToggleInstancePause)
		// Instance credentials management (Vault integration)
		payments.GET("/:id/instances/:instanceId/credentials", instanceHandler.GetInstanceCredentials)
		payments.PUT("/:id/instances/:instanceId/credentials", instanceHandler.UpdateInstanceCredentials)
		// Instance wallet management (multi-wallet support)
		payments.GET("/:id/instances/:instanceId/wallets", instanceHandler.GetInstanceWallets)
		payments.POST("/:id/instances/:instanceId/wallets", instanceHandler.AddInstanceWallet)
		payments.DELETE("/:id/instances/:instanceId/wallets/:walletId", instanceHandler.RemoveInstanceWallet)
		payments.PUT("/:id/instances/:instanceId/wallets/:walletId", instanceHandler.ToggleInstanceWallet)
		payments.GET("/:id/instances/:instanceId/best-wallet", instanceHandler.GetBestWalletForCurrency)

		// Platform Accounts Proxy
		// Forward /platform/* to wallet-service
		platform := protected.Group("/platform")
		{
			platform.Any("/*path", proxyHandler.ProxyRequest)
		}

		// Credit Management System (replaces test-mode)
		creditHandler := handlers.NewCreditHandler(adminDB, mainDB)
		credits := protected.Group("/credits")
		{
			// Predefined reason types
			credits.GET("/reason-types", creditHandler.GetReasonTypes)

			// Hot wallets list with balances
			credits.GET("/hot-wallets", creditHandler.GetHotWallets)

			// Individual credit
			credits.POST("/single", creditHandler.SingleCredit)

			// Mass credit
			credits.POST("/mass/preview", creditHandler.MassCreditPreview)
			credits.POST("/mass", creditHandler.MassCredit)

			// Promotion/Contest credits
			credits.POST("/promotion", creditHandler.PromotionCredit)

			// Campaigns history
			credits.GET("/campaigns", creditHandler.GetCampaigns)
			credits.GET("/campaigns/:id", creditHandler.GetCampaignDetails)

			// Audit logs
			credits.GET("/logs", creditHandler.GetCreditLogs)
		}

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
