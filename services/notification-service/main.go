package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Database
	db, err := database.Initialize(cfg.DBUrl)
	if err != nil {
		log.Println("Warning: Database not available, notification storage disabled:", err)
		// Continue without DB - notifications won't be persisted
	}
	if db != nil {
		defer db.Close()
	}

	// Initialize RabbitMQ
	rabbitConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}
	defer rabbitChannel.Close()

	// Declare all exchanges (for receiving events)
	exchanges := []string{
		"wallet.events",
		"transaction.events",
		"transfer.events",
		"exchange.events",
		"card.events",
		"notification.events",
	}

	for _, exchange := range exchanges {
		err = rabbitChannel.ExchangeDeclare(
			exchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to declare exchange %s: %v", exchange, err)
		}
	}

	// Initialize repositories first (needed for both API and event consumers)
	var notificationRepo *repository.NotificationRepository
	if db != nil {
		notificationRepo = repository.NewNotificationRepository(db)
	}

	// Initialize notification service (for consuming events) - pass repository for persistence
	notificationService := services.NewNotificationService(rabbitChannel, cfg, notificationRepo)

	// Start consumers
	if err := notificationService.Start(); err != nil {
		log.Fatal("Failed to start notification consumers:", err)
	}

	// Initialize handlers
	var notificationHandler *handlers.NotificationHandler
	if notificationRepo != nil {
		notificationHandler = handlers.NewNotificationHandler(notificationRepo)
	}

	// Setup Gin for health checks and admin endpoints
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
			"service": "notification-service",
		})
	})

	// Admin endpoints
	router.GET("/api/v1/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "running",
			"consumers":  5,
			"message":    "Notification service is consuming events",
		})
	})

	// Notification API routes (protected by JWT)
	if notificationHandler != nil {
		api := router.Group("/api/v1")
		api.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			api.GET("/notifications", notificationHandler.GetNotifications)
			api.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
			api.POST("/notifications/:id/read", notificationHandler.MarkAsRead)
			api.POST("/notifications/read-all", notificationHandler.MarkAllAsRead)
			api.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Notification Service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
