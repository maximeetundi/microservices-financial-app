package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/auth-service/internal/config"
	"github.com/crypto-bank/auth-service/internal/database"
	"github.com/crypto-bank/auth-service/internal/handlers"
	"github.com/crypto-bank/auth-service/internal/middleware"
	"github.com/crypto-bank/auth-service/internal/repository"
	"github.com/crypto-bank/auth-service/internal/services"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db, redisClient)

	// Initialize services
	authService := services.NewAuthService(userRepo, sessionRepo, cfg)
	emailService := services.NewEmailService(cfg)
	smsService := services.NewSMSService(cfg)
	totpService := services.NewTOTPService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, emailService, smsService, totpService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Security middleware
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// Auth routes
	api := router.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		api.POST("/logout", middleware.JWTAuth(cfg.JWTSecret), authHandler.Logout)
		api.POST("/forgot-password", authHandler.ForgotPassword)
		api.POST("/reset-password", authHandler.ResetPassword)
		api.POST("/verify-email", authHandler.VerifyEmail)
		api.POST("/verify-phone", authHandler.VerifyPhone)
		
		// 2FA routes
		api.POST("/enable-2fa", middleware.JWTAuth(cfg.JWTSecret), authHandler.Enable2FA)
		api.POST("/verify-2fa", middleware.JWTAuth(cfg.JWTSecret), authHandler.Verify2FA)
		api.POST("/disable-2fa", middleware.JWTAuth(cfg.JWTSecret), authHandler.Disable2FA)
		
		// Session management
		api.GET("/sessions", middleware.JWTAuth(cfg.JWTSecret), authHandler.GetSessions)
		api.DELETE("/sessions/:session_id", middleware.JWTAuth(cfg.JWTSecret), authHandler.RevokeSession)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}