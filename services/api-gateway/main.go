package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/routes"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	serviceManager := services.NewServiceManager(cfg)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration - Allow all origins to fix 301 redirect CORS issues
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false, // Must be false when AllowAllOrigins is true
	}))

	// Security middleware
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "api-gateway",
		})
	})

	// API routes - with service prefixes
	
	// Auth Service routes: /auth-service/api/v1/...
	authService := router.Group("/auth-service/api/v1")
	{
		auth := authService.Group("/auth")
		routes.SetupAuthRoutes(auth, serviceManager)
	}

	// Wallet Service routes: /wallet-service/api/v1/...
	walletService := router.Group("/wallet-service/api/v1")
	walletService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupWalletRoutes(walletService.Group("/wallets"), serviceManager)
		routes.SetupDashboardRoutes(walletService.Group("/dashboard"), serviceManager)
	}

	// Transfer Service routes: /transfer-service/api/v1/...
	transferService := router.Group("/transfer-service/api/v1")
	transferService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupTransferRoutes(transferService.Group("/transfers"), serviceManager)
	}

	// Exchange Service routes: /exchange-service/api/v1/...
	exchangeService := router.Group("/exchange-service/api/v1")
	exchangeService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupExchangeRoutes(exchangeService.Group("/exchange"), serviceManager)
	}

	// Card Service routes: /card-service/api/v1/...
	cardService := router.Group("/card-service/api/v1")
	cardService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupCardRoutes(cardService.Group("/cards"), serviceManager)
	}

	// Notification Service routes: /notification-service/api/v1/...
	notificationService := router.Group("/notification-service/api/v1")
	notificationService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupNotificationRoutes(notificationService.Group("/notifications"), serviceManager)
	}

	// User routes: /auth-service/api/v1/users/...
	userService := authService.Group("/users")
	userService.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		routes.SetupUserRoutes(userService, serviceManager)
	}

	// Legacy API routes (for backward compatibility) - keeping /api/v1 prefix
	api := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		routes.SetupAuthRoutes(auth, serviceManager)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			routes.SetupUserRoutes(protected.Group("/users"), serviceManager)
			routes.SetupWalletRoutes(protected.Group("/wallets"), serviceManager)
			routes.SetupTransferRoutes(protected.Group("/transfers"), serviceManager)
			routes.SetupExchangeRoutes(protected.Group("/exchange"), serviceManager)
			routes.SetupCardRoutes(protected.Group("/cards"), serviceManager)
			routes.SetupNotificationRoutes(protected.Group("/notifications"), serviceManager)
		}

		// Admin routes (through gateway)
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			routes.SetupAdminRoutes(admin, serviceManager)
		}
	}


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}