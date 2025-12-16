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

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3002",
			"https://app.maximeetundi.store",
			"https://admin.maximeetundi.store",
			"https://api.app.maximeetundi.store",
			"https://api.admin.maximeetundi.store",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
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

	// API routes
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

		// Admin routes
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