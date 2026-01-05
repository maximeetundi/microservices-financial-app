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
	
	"messaging-service/internal/handlers"
	"messaging-service/internal/middleware"
	"messaging-service/internal/services"
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

	storageService, err := services.NewStorageService(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		"chat-attachments",
		minioUseSSL,
		minioPublicURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	log.Println("Connected to MinIO")

	// Create handler with storage
	handler := handlers.NewMessageHandler(db, storageService)

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

	// API routes (protected with JWT)
	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))
	{
		// File upload
		api.POST("/upload", handler.UploadFile)
		
		// User-to-user messaging
		api.GET("/conversations", handler.GetConversations)
		api.POST("/conversations", handler.CreateConversation)
		api.GET("/conversations/:id/messages", handler.GetMessages)
		api.POST("/conversations/:id/messages", handler.SendMessage)
		api.POST("/messages/:id/read", handler.MarkAsRead)
		api.DELETE("/conversations/:id", handler.DeleteConversation)

		// Association chat (centralized in messaging-service)
		api.GET("/associations/:id/chat", handler.GetAssociationChat)
		api.POST("/associations/:id/chat", handler.SendAssociationMessage)
		api.DELETE("/associations/:id/chat/:messageId", handler.DeleteAssociationMessage)
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
