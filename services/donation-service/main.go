package main

import (
	"context"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/handlers"
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
	kafkaClient, err := messaging.NewKafkaClient([]string{cfg.KafkaBrokers}, cfg.KafkaGroupID)
	if err != nil {
		log.Printf("Warning: Failed to connect to Kafka: %v (running without messaging)", err)
		// Assuming we can run without kafka for dev, but really should fail for prod features.
		// For now, allow it but functionality will be limited.
	} else {
		log.Println("Connected to Kafka successfully")
		defer kafkaClient.Close()
	}

	// 4. Initialize Repositories
	campaignRepo := repository.NewCampaignRepository(db)
	donationRepo := repository.NewDonationRepository(db)

	// 5. Initialize Services
	walletClient := services.NewWalletClient()
	
	campaignService := services.NewCampaignService(campaignRepo)
	donationService := services.NewDonationService(donationRepo, campaignRepo, walletClient, kafkaClient)
	
	paymentConsumer := services.NewPaymentConsumer(kafkaClient, donationRepo, campaignRepo)

	// 6. Start Payment Consumer
	if kafkaClient != nil {
		go paymentConsumer.Start()
	}

	// 7. Initialize Handlers
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	donationHandler := handlers.NewDonationHandler(donationService)

	// 8. Setup Router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID"},
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
		
		// Campaigns
		api.POST("/campaigns", campaignHandler.Create)
		api.GET("/campaigns", campaignHandler.List)
		api.GET("/campaigns/:id", campaignHandler.Get)
		api.PUT("/campaigns/:id", campaignHandler.Update)
		
		// Donations
		api.POST("/donations", donationHandler.Initiate)
		api.GET("/donations", donationHandler.List)
		
		// Utility
		api.POST("/upload", campaignHandler.UploadImage) // Stub
	}

	log.Printf("Donation service starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
