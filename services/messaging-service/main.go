package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:secure_password@mongodb:27017/messaging?authSource=admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Ping database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db = client.Database("messaging")
	log.Println("Connected to MongoDB")

	// Setup Gin router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://app.maximeetundi.store", "https://admin.maximeetundi.store", "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { 
		c.JSON(200, gin.H{"status": "ok", "service": "messaging-service"}) 
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// User-to-user messaging
		api.GET("/conversations", getConversations)
		api.GET("/conversations/:id", getConversation)
		api.POST("/conversations", createConversation)
		api.POST("/conversations/:id/messages", sendMessage)
		api.GET("/conversations/:id/messages", getMessages)
		api.POST("/messages/:id/read", markAsRead)
		api.DELETE("/conversations/:id", deleteConversation)

		// Association chat
		api.GET("/associations/:id/chat", getAssociationChat)
		api.POST("/associations/:id/chat", sendAssociationMessage)
		api.DELETE("/associations/:id/chat/:messageId", deleteAssociationMessage)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8095"
	}

	log.Printf("Messaging service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Placeholder handlers
func getConversations(c *gin.Context) {
	c.JSON(200, gin.H{"conversations": []interface{}{}})
}

func getConversation(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func createConversation(c *gin.Context) {
	c.JSON(201, gin.H{})
}

func sendMessage(c *gin.Context) {
	c.JSON(201, gin.H{})
}

func getMessages(c *gin.Context) {
	c.JSON(200, gin.H{"messages": []interface{}{}})
}

func markAsRead(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func deleteConversation(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func getAssociationChat(c *gin.Context) {
	c.JSON(200, gin.H{"messages": []interface{}{}})
}

func sendAssociationMessage(c *gin.Context) {
	c.JSON(201, gin.H{})
}

func deleteAssociationMessage(c *gin.Context) {
	c.JSON(200, gin.H{})
}
