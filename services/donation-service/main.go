package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	log.Println("----------------------------------------------------------------")
	log.Println("--- STARTING DONATION SERVICE ---")
	log.Println("----------------------------------------------------------------")

	// 1. Load Config
	cfg := config.Load()

	// 2. Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	db := client.Database(cfg.MongoDBName)

	// 3. Connect to Kafka
	// 3. Connect to Kafka
	kafkaClient := messaging.NewKafkaClient([]string{cfg.KafkaBrokers}, cfg.KafkaGroupID)
	log.Println("Kafka client initialized")
	defer kafkaClient.Close()

	// 4. Initialize Repositories
	campaignRepo := repository.NewCampaignRepository(db)
	donationRepo := repository.NewDonationRepository(db)

	// 5. Initialize Services
	walletClient := services.NewWalletClient()
	userClient := services.NewUserClient()
	
	campaignService := services.NewCampaignService(campaignRepo)
	donationService := services.NewDonationService(donationRepo, campaignRepo, walletClient, userClient, kafkaClient)
	
	paymentConsumer := services.NewPaymentConsumer(kafkaClient, donationRepo, campaignRepo)

	// 6. Start Payment Consumer
	if kafkaClient != nil {
		go paymentConsumer.Start()
	}

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

	storageService, err := services.NewStorageService(minioEndpoint, minioAccessKey, minioSecretKey, "campaign-media", minioUseSSL, minioPublicURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v", err)
	}

	// 7. Initialize Handlers
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	donationHandler := handlers.NewDonationHandler(donationService)
	uploadHandler := handlers.NewUploadHandler(storageService)

	// 8. Setup Router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	
	api := router.Group("/api/v1")
	{
		// Internal routes (if any)
		// Public routes
		
	// Public Routes
	api.GET("/campaigns", campaignHandler.List)
	api.GET("/campaigns/:id", campaignHandler.Get)
	api.GET("/donations", donationHandler.List)

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/campaigns", campaignHandler.Create)
		protected.PUT("/campaigns/:id", campaignHandler.Update)
		protected.POST("/donations", donationHandler.Initiate)
		protected.POST("/upload", uploadHandler.UploadMedia)
	}
	}

	log.Printf("Donation service starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
