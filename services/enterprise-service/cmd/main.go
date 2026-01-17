package main

import (
	"context"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Connect to MongoDB
	mongoClient, err := database.ConnectMongoDB(cfg.MongoDBURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()
	log.Println("Connected to MongoDB")

	// 3. Initialize Kafka Client for Notifications
	var kafkaClient *messaging.KafkaClient
	var notificationClient *services.NotificationClient
	
	if len(cfg.KafkaBrokers) > 0 && cfg.KafkaBrokers[0] != "" {
		kafkaClient = messaging.NewKafkaClient(cfg.KafkaBrokers, "enterprise-service-group")
		notificationClient = services.NewNotificationClient(kafkaClient)
		log.Println("Kafka client initialized for notifications")
		
		defer func() {
			if kafkaClient != nil {
				kafkaClient.Close()
			}
		}()
	} else {
		log.Println("Warning: Kafka not configured, notifications will be disabled")
	}

	// 4. Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://app.maximeetundi.store", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "enterprise-service"})
	})

	// 5. Initialize Components
	entRepo := repository.NewEnterpriseRepository(mongoClient.Database(cfg.DBName))
	empRepo := repository.NewEmployeeRepository(mongoClient.Database(cfg.DBName))
	payRepo := repository.NewPayrollRepository(mongoClient.Database(cfg.DBName))
	subRepo := repository.NewSubscriptionRepository(mongoClient.Database(cfg.DBName))
	invRepo := repository.NewInvoiceRepository(mongoClient.Database(cfg.DBName))
	batchRepo := repository.NewBatchRepository(mongoClient.Database(cfg.DBName))
	
	// Initialize Storage Service
	storageService, err := services.NewStorageService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
		cfg.PublicURL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage service: %v", err)
	}

	salaryService := services.NewSalaryService()
	authClient := services.NewAuthClient()
	
	// Initialize services with notification client
	empService := services.NewEmployeeService(empRepo, salaryService, authClient, entRepo, notificationClient)
	payService := services.NewPayrollService(payRepo, empRepo, salaryService, entRepo, notificationClient)
	billService := services.NewBillingService(subRepo, invRepo, batchRepo, entRepo, notificationClient)

	entHandler := handlers.NewEnterpriseHandler(entRepo, empRepo, storageService)
	empHandler := handlers.NewEmployeeHandler(empService)
	payHandler := handlers.NewPayrollHandler(payService)
	billHandler := handlers.NewBillingHandler(billService)

	api := router.Group("/api/v1")
	
	// Public routes (no auth required)
	api.GET("/invitations/:id", empHandler.GetInvitationDetails) // Get invitation details for accept page
	
	// Apply JWT Auth Middleware to protected routes
	api.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		// Enterprise Routes
		api.GET("/enterprises", entHandler.ListEnterprises)
		api.POST("/enterprises", entHandler.CreateEnterprise)
		api.POST("/enterprises/logo", entHandler.UploadLogo)
		api.PUT("/enterprises/:id", entHandler.UpdateEnterprise)
		api.GET("/enterprises/:id", entHandler.GetEnterprise)
		api.GET("/enterprises/:id/employees", empHandler.ListEmployees)
		api.DELETE("/enterprises/:id/employees/:employeeId", empHandler.TerminateEmployee)

		// Enterprise Data Export Routes (required before deletion)
		exportHandler := handlers.NewExportHandler(entRepo, empRepo, subRepo, invRepo, payRepo)
		api.GET("/enterprises/:id/export", exportHandler.ExportEnterpriseData)
		api.GET("/enterprises/:id/export/status", exportHandler.CheckExportStatus)

		// Employee Routes
		api.POST("/employees/invite", empHandler.InviteEmployee)
		api.POST("/employees/accept", empHandler.AcceptInvitation)
		api.PUT("/employees/:id/promote", empHandler.PromoteEmployee)
		api.GET("/employees/me", empHandler.GetMyEmployee) // Get current user's employee record for RBAC

		// Payroll Routes
		api.POST("/enterprises/:id/payroll/preview", payHandler.PreviewPayroll)
		api.POST("/enterprises/:id/payroll/execute", payHandler.ExecutePayroll)
		api.GET("/enterprises/:id/payroll/runs", payHandler.ListPayrollRuns) // New endpoint

		// Billing Routes
		api.POST("/invoices", billHandler.CreateInvoice)
		api.POST("/enterprises/:id/invoices/import", billHandler.ImportInvoices)
		api.POST("/enterprises/:id/invoices/batches", billHandler.CreateBatchInvoices)
		api.POST("/enterprises/:id/invoices/batches/:batch_id/validate", billHandler.ValidateBatch)
		api.POST("/enterprises/:id/invoices/batches/:batch_id/schedule", billHandler.ScheduleBatch)
		
		// QR Code Routes
		qrHandler := handlers.NewQRCodeHandler("https://app.maximeetundi.store", entRepo)
		api.GET("/enterprises/:id/qrcode", qrHandler.GenerateEnterpriseQR)
		api.GET("/enterprises/:id/services/:serviceId/qrcode", qrHandler.GenerateServiceQR)
		api.GET("/enterprises/:id/groups/:groupId/qrcode", qrHandler.GenerateServiceGroupQR)

		// Subscription Routes
		subHandler := handlers.NewSubscriptionHandler(subRepo)
		api.GET("/enterprises/:id/subscriptions", subHandler.ListSubscriptions)
		api.POST("/enterprises/:id/subscriptions", subHandler.CreateSubscription)

		// Multi-Approval System Routes
		approvalRepo := repository.NewApprovalRepository(mongoClient.Database(cfg.DBName))
		approvalService := services.NewApprovalService(approvalRepo, empRepo, entRepo, notificationClient)
		approvalHandler := handlers.NewApprovalHandler(approvalService, empService)
		api.GET("/enterprises/:id/approvals", approvalHandler.GetPendingApprovals)
		api.POST("/enterprises/:id/actions", approvalHandler.InitiateAction)
		api.POST("/approvals/:id/approve", approvalHandler.ApproveAction)
		api.POST("/approvals/:id/reject", approvalHandler.RejectAction)
	}

	// 6. Start Server
	port := cfg.Port
	if port == "" {
		port = "8097"
	}
	log.Printf("Enterprise Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

