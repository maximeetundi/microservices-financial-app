package main

import (
	"context"
	"log"
	"os"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/handlers"
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

	// 3. Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "enterprise-service"})
	})

	// 4. Initialize Components
	entRepo := repository.NewEnterpriseRepository(mongoClient.Database(cfg.DBName))
	empRepo := repository.NewEmployeeRepository(mongoClient.Database(cfg.DBName))
	payRepo := repository.NewPayrollRepository(mongoClient.Database(cfg.DBName))
	subRepo := repository.NewSubscriptionRepository(mongoClient.Database(cfg.DBName))
	invRepo := repository.NewInvoiceRepository(mongoClient.Database(cfg.DBName))
	batchRepo := repository.NewBatchRepository(mongoClient.Database(cfg.DBName))
	
	salaryService := services.NewSalaryService()
	authClient := services.NewAuthClient()
	empService := services.NewEmployeeService(empRepo, salaryService, authClient)
	payService := services.NewPayrollService(payRepo, empRepo, salaryService)
	billService := services.NewBillingService(subRepo, invRepo, batchRepo)

	entHandler := handlers.NewEnterpriseHandler(entRepo)
	empHandler := handlers.NewEmployeeHandler(empService)
	payHandler := handlers.NewPayrollHandler(payService)
	billHandler := handlers.NewBillingHandler(billService)

	api := router.Group("/api/v1")
	{
		// Enterprise Routes
		api.POST("/enterprises", entHandler.CreateEnterprise)
		api.GET("/enterprises/:id", entHandler.GetEnterprise)
		api.GET("/enterprises/:id/employees", empHandler.ListEmployees)

		// Employee Routes
		api.POST("/employees/invite", empHandler.InviteEmployee)
		api.POST("/employees/accept", empHandler.AcceptInvitation)
		api.PUT("/employees/:id/promote", empHandler.PromoteEmployee)

		// Payroll Routes
		api.POST("/enterprises/:id/payroll/preview", payHandler.PreviewPayroll)
		api.POST("/enterprises/:id/payroll/execute", payHandler.ExecutePayroll)

		// Billing Routes
		api.POST("/invoices", billHandler.CreateInvoice)
		api.POST("/enterprises/:id/invoices/import", billHandler.ImportInvoices)
		api.POST("/enterprises/:id/invoices/batches/:batch_id/validate", billHandler.ValidateBatch)
		api.POST("/enterprises/:id/invoices/batches/:batch_id/schedule", billHandler.ScheduleBatch)
	}

	// 4. Start Server
	port := cfg.Port
	if port == "" {
		port = "8097"
	}
	log.Printf("Enterprise Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
