package main

import (
	"log"
	"os"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8091"
	}

	// Initialize database
	db, err := database.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Gin
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "association-service"})
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Associations
		v1.POST("/associations", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented yet - coming soon"})
		})
		v1.GET("/associations", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented yet - coming soon"})
		})
		v1.GET("/associations/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented yet - coming soon"})
		})

		// Placeholder routes to show structure
		v1.POST("/associations/:id/join", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Join association - coming soon"})
		})
		v1.GET("/associations/:id/members", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "List members - coming soon"})
		})
		v1.POST("/associations/:id/meetings", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Create meeting - coming soon"})
		})
		v1.POST("/associations/:id/contributions", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Record contribution - coming soon"})
		})
		v1.POST("/associations/:id/loans", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Request loan - coming soon"})
		})
	}

	log.Printf("Association service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
