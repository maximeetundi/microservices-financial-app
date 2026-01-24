package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	log.Println("================================================================")
	log.Println("=== STARTING SHOP SERVICE ===")
	log.Println("================================================================")

	// 1. Load Config
	cfg := config.Load()
	middleware.SetJWTSecret(cfg.JWTSecret)

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
	var kafkaClient *messaging.KafkaClient
	if cfg.KafkaBrokers != "" {
		kafkaClient = messaging.NewKafkaClient([]string{cfg.KafkaBrokers}, cfg.KafkaGroupID)
		log.Println("Kafka client initialized")
		defer kafkaClient.Close()
	}

	// 4. Initialize Storage
	var storageService *services.StorageService
	if cfg.MinioEndpoint != "" {
		var err error
		storageService, err = services.NewStorageService(
			cfg.MinioEndpoint,
			cfg.MinioAccessKey,
			cfg.MinioSecretKey,
			cfg.MinioBucket,
			cfg.MinioUseSSL,
			cfg.MinioPublicURL,
		)
		if err != nil {
			log.Printf("Warning: Failed to initialize storage service: %v", err)
		} else {
			log.Println("Storage service initialized")
		}
	}

	// 5. Initialize Repositories
	shopRepo := repository.NewShopRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	clientRepo := repository.NewClientRepository(db)
	log.Println("Repositories initialized")

	// 6. Initialize External Clients
	walletClient := services.NewWalletClient()
	exchangeClient := services.NewExchangeClient()
	
	// QR Service (using app URL)
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://app.maximeetundi.store"
	}
	qrService := services.NewQRService(appURL)

	// 7. Initialize Services
	shopService := services.NewShopService(shopRepo, productRepo, walletClient, qrService, storageService)
	productService := services.NewProductService(productRepo, categoryRepo, shopRepo, qrService, storageService)
	categoryService := services.NewCategoryService(categoryRepo, shopRepo, qrService, storageService)
	orderService := services.NewOrderService(orderRepo, productRepo, shopRepo, walletClient, exchangeClient, kafkaClient)
	clientService := services.NewClientService(clientRepo, shopRepo, kafkaClient)
	log.Println("Services initialized")

	// 8. Initialize Payment Consumer
	if kafkaClient != nil {
		paymentConsumer := services.NewPaymentConsumer(kafkaClient, orderService)
		go paymentConsumer.Start()
		log.Println("Payment consumer started")
	}

	// 9. Initialize Handlers
	shopHandler := handlers.NewShopHandler(shopService)
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	orderHandler := handlers.NewOrderHandler(orderService)
	uploadHandler := handlers.NewUploadHandler(storageService)
	clientHandler := handlers.NewClientHandler(clientService)

	// 10. Setup Router
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

	// Health & Metrics
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "shop-service"}) })
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API Routes
	api := router.Group("/api/v1")
	{
		// Public routes
		api.GET("/shops", shopHandler.List)
		api.GET("/shops/by-wallet/:wallet_id", shopHandler.GetByWalletID)
		api.GET("/shops/:id", shopHandler.Get)
		api.GET("/shops/:id/products", productHandler.ListByShop)
		api.GET("/shops/:id/products/:productSlug", productHandler.Get)
		api.GET("/shops/:id/categories", categoryHandler.ListByShop)
		api.GET("/shops/:id/categories/tree", categoryHandler.ListWithHierarchy)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Shops
			protected.POST("/shops", shopHandler.Create)
			protected.GET("/my-shops", shopHandler.GetMyShops)
			protected.PUT("/shops/:id", shopHandler.Update)
			protected.DELETE("/shops/:id", shopHandler.Delete)
			protected.POST("/shops/:id/managers", shopHandler.InviteManager)
			protected.DELETE("/shops/:id/managers/:userId", shopHandler.RemoveManager)
			protected.GET("/shops/:id/products", productHandler.ListByShop)

			protected.GET("/shops/:id/categories", categoryHandler.ListByShop)
			protected.GET("/shops/:id/categories/tree", categoryHandler.ListWithHierarchy)

			// Client Invitations (for private shops)
			protected.POST("/shops/:id/clients", clientHandler.InviteClient)
			protected.GET("/shops/:id/clients", clientHandler.ListShopClients)
			protected.DELETE("/shops/:id/clients/:clientId", clientHandler.RevokeClientAccess)
			protected.GET("/my-invitations", clientHandler.GetMyInvitations)
			protected.POST("/invitations/accept", clientHandler.AcceptInvitation)
			protected.DELETE("/invitations/:id", clientHandler.DeclineInvitation)
			protected.GET("/my-private-shops", clientHandler.GetMyPrivateShops)

			// Products
			protected.POST("/products", productHandler.Create)
			protected.GET("/products/:id", productHandler.GetByID)
			protected.PUT("/products/:id", productHandler.Update)
			protected.DELETE("/products/:id", productHandler.Delete)

			// Categories
			protected.POST("/categories", categoryHandler.Create)
			protected.GET("/categories/:id", categoryHandler.Get)
			protected.PUT("/categories/:id", categoryHandler.Update)
			protected.DELETE("/categories/:id", categoryHandler.Delete)

			// Orders
			protected.POST("/orders", orderHandler.Create)
			protected.GET("/orders", orderHandler.ListMyOrders)
			protected.GET("/orders/:id", orderHandler.Get)
			protected.PUT("/orders/:id/status", orderHandler.UpdateStatus)
			protected.POST("/orders/:id/refund", orderHandler.Refund)
			protected.GET("/shop-orders/:shopId", orderHandler.ListShopOrders)

			// Upload
			protected.POST("/upload", uploadHandler.UploadMedia)
		}
	}

	log.Printf("Shop service starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

