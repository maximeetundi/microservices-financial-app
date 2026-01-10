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
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/services"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" { c.Next(); return }
		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath(); if path == "" { path = c.Request.URL.Path }
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(time.Since(start).Seconds())
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := database.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize schema: %v", err)
	}

	// Initialize repositories
	convRepo := repository.NewConversationRepository(database.Database)
	msgRepo := repository.NewMessageRepository(database.Database)
	agentRepo := repository.NewAgentRepository(database.Database)

	// Initialize event publisher for admin notifications
	var eventPublisher *services.EventPublisher
	publisher, err := services.NewEventPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize event publisher: %v (admin notifications disabled)", err)
	} else {
		eventPublisher = publisher
	}

	// Initialize storage service
	storageService, err := services.NewStorageService(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOBucket,
		cfg.MinIOUseSSL,
		cfg.MinIOPublicURL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v", err)
	}

	// Initialize services
	supportService := services.NewSupportService(convRepo, msgRepo, agentRepo, cfg, eventPublisher, storageService)

	// Initialize handlers
	supportHandler := handlers.NewSupportHandler(supportService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())
	router.Use(middleware.SecurityHeaders(), middleware.RateLimiter())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "support-service", "database": "mongodb"}) })

	// API routes
	api := router.Group("/api/v1")
	{
		// Protected client routes
		support := api.Group("/support")
		support.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// Conversation management
			support.POST("/conversations", supportHandler.CreateConversation)
			support.GET("/conversations", supportHandler.GetConversations)
			support.GET("/conversations/:id", supportHandler.GetConversation)
			support.POST("/conversations/:id/messages", supportHandler.SendMessage)
			support.GET("/conversations/:id/messages", supportHandler.GetMessages)
			support.POST("/conversations/:id/escalate", supportHandler.EscalateConversation)
			support.PUT("/conversations/:id/close", supportHandler.CloseConversation)
			support.POST("/upload", supportHandler.UploadFile)
		}

		// Admin routes - use AdminJWTSecret if available, otherwise fall back to JWTSecret
		adminJWTSecret := cfg.AdminJWTSecret
		if adminJWTSecret == "" {
			adminJWTSecret = cfg.JWTSecret
		}
		admin := api.Group("/admin/support")
		admin.Use(middleware.JWTAuth(adminJWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			// Conversation management
			admin.GET("/conversations", supportHandler.AdminGetConversations)
			admin.GET("/conversations/:id", supportHandler.GetConversation)
			admin.GET("/conversations/:id/messages", supportHandler.GetMessages)
			admin.POST("/conversations/:id/messages", supportHandler.AdminSendMessage)
			admin.POST("/upload", supportHandler.UploadFile)
			admin.PUT("/conversations/:id/assign", supportHandler.AdminAssignAgent)
			admin.PUT("/conversations/:id/close", supportHandler.CloseConversation)


			// Agent management
			admin.GET("/agents", supportHandler.AdminGetAgents)
			admin.PUT("/agents/:id/availability", supportHandler.AdminUpdateAgentAvailability)

			// Statistics
			admin.GET("/stats", supportHandler.AdminGetStats)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	log.Printf("Support Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
