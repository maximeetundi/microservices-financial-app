package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/services"
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

	// Initialize repositories
	convRepo := repository.NewConversationRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	agentRepo := repository.NewAgentRepository(db)

	// Initialize event publisher for admin notifications
	var eventPublisher *services.EventPublisher
	publisher, err := services.NewEventPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize event publisher: %v (admin notifications disabled)", err)
	} else {
		eventPublisher = publisher
	}

	// Initialize services
	supportService := services.NewSupportService(convRepo, msgRepo, agentRepo, cfg, eventPublisher)

	// Initialize handlers
	supportHandler := handlers.NewSupportHandler(supportService)

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
			"service": "support-service",
		})
	})

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
		}

		// Admin routes
		admin := api.Group("/admin/support")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			// Conversation management
			admin.GET("/conversations", supportHandler.AdminGetConversations)
			admin.GET("/conversations/:id", supportHandler.GetConversation)
			admin.GET("/conversations/:id/messages", supportHandler.GetMessages)
			admin.POST("/conversations/:id/messages", supportHandler.AdminSendMessage)
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
