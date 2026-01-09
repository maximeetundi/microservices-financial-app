package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	ticketsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "tickets_sold_total", Help: "Total tickets sold"}, []string{"event", "tier"})
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
	log.Println("----------------------------------------------------------------")
	log.Println("--- STARTING TICKET SERVICE VERSION 2.0 (DB FIX APPLIED) ---")
	log.Println("----------------------------------------------------------------")

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize schema
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize repositories
	eventRepo := repository.NewEventRepository(database.DB)
	tierRepo := repository.NewTierRepository(database.DB)
	ticketRepo := repository.NewTicketRepository(database.DB)

	// Initialize RabbitMQ
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	mqClient, err := database.InitializeRabbitMQ(rabbitmqURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize RabbitMQ: %v", err)
	} else {
		defer mqClient.Close()
		
		// Start Payment Status Consumer
		paymentConsumer := services.NewPaymentStatusConsumer(mqClient, ticketRepo)
		if err := paymentConsumer.Start(); err != nil {
			log.Printf("Warning: Failed to start PaymentStatusConsumer: %v", err)
		}
	}

	// Initialize service
	ticketService := services.NewTicketService(eventRepo, tierRepo, ticketRepo, mqClient)

	// Initialize MinIO storage
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "minio:9000"
	}
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin123"
	}
	minioPublicURL := os.Getenv("MINIO_PUBLIC_URL")
	if minioPublicURL == "" {
		minioPublicURL = "https://minio.maximeetundi.store"
	}
	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true"

	storageService, err := services.NewStorageService(minioEndpoint, minioAccessKey, minioSecretKey, "event-images", minioUseSSL, minioPublicURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v", err)
	}

	// Initialize handlers
	handler := handlers.NewTicketHandler(ticketService)
	uploadHandler := handlers.NewUploadHandler(storageService)

	// Setup Gin
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		api.GET("/events/active", handler.GetActiveEvents)
		api.GET("/events/code/:code", handler.GetEventByCode)
		api.GET("/events/public/:id", handler.GetEvent)
		api.GET("/icons", handler.GetAvailableIcons)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Event management (organizer)
			protected.POST("/events", handler.CreateEvent)
			protected.GET("/events", handler.GetMyEvents)
			protected.GET("/events/:id", handler.GetEvent)
			protected.PUT("/events/:id", handler.UpdateEvent)
			protected.DELETE("/events/:id", handler.DeleteEvent)
			protected.POST("/events/:id/publish", handler.PublishEvent)
			protected.GET("/events/:id/tickets", handler.GetEventTickets)
			protected.GET("/events/:id/stats", handler.GetEventStats)

			// Upload image (for events)
			protected.POST("/upload", uploadHandler.UploadImage)

			// Ticket purchase (buyer)
			protected.POST("/tickets/purchase", handler.PurchaseTicket)
			protected.GET("/tickets", handler.GetMyTickets)
			protected.GET("/tickets/:id", handler.GetTicket)

			// Ticket verification (organizer)
			protected.POST("/tickets/verify", handler.VerifyTicket)
			protected.POST("/tickets/:id/use", handler.UseTicket)
		}
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Ticket service starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
