package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/notification-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

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

	// Initialize notification service
	notificationService := services.NewNotificationService(rabbitChannel, cfg)

	// Start consumers
	if err := notificationService.Start(); err != nil {
		log.Fatal("Failed to start notification consumers:", err)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Notification Service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
