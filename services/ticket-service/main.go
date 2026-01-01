package main

import (
	"log"
	"os"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	// Initialize service
	ticketService := services.NewTicketService(eventRepo, tierRepo, ticketRepo)

	// Initialize handler
	handler := handlers.NewTicketHandler(ticketService)

	// Setup Gin
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
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
