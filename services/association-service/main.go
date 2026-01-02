package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/services"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	assocRepo := repository.NewAssociationRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	treasuryRepo := repository.NewTreasuryRepository(db)
	meetingRepo := repository.NewMeetingRepository(db)
	loanRepo := repository.NewLoanRepository(db)

	// Initialize service
	assocService := services.NewAssociationService(
		assocRepo,
		memberRepo,
		treasuryRepo,
		meetingRepo,
		loanRepo,
	)

	// Initialize handlers
	handler := handlers.NewHandler(assocService)

	// Setup Gin router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "association-service"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public/Admin routes (for admin dashboard without user JWT)
		api.GET("/associations", handler.GetAllAssociations)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Associations CRUD
			protected.POST("/associations", handler.CreateAssociation)
			protected.GET("/associations/me", handler.GetAssociations)
			protected.GET("/associations/:id", handler.GetAssociation)
			protected.PUT("/associations/:id", handler.UpdateAssociation)
			protected.DELETE("/associations/:id", handler.DeleteAssociation)

			// Members
			protected.POST("/associations/:id/join", handler.JoinAssociation)
			protected.POST("/associations/:id/leave", handler.LeaveAssociation)
			protected.GET("/associations/:id/members", handler.GetMembers)
			protected.PUT("/associations/:id/members/:uid/role", handler.UpdateMemberRole)

			// Contributions & Treasury
			protected.POST("/associations/:id/contributions", handler.PayContribution)
			protected.GET("/associations/:id/treasury", handler.GetTreasury)
			protected.POST("/associations/:id/distribute", handler.DistributeFunds)

			// Loans
			protected.POST("/associations/:id/loans", handler.RequestLoan)
			protected.GET("/associations/:id/loans", handler.GetLoans)
			protected.PUT("/loans/:loanId/approve", handler.ApproveLoan)
			protected.POST("/loans/:loanId/repay", handler.RepayLoan)

			// Meetings
			protected.POST("/associations/:id/meetings", handler.CreateMeeting)
			protected.GET("/associations/:id/meetings", handler.GetMeetings)
			protected.POST("/meetings/:meetingId/attendance", handler.RecordAttendance)
			protected.PUT("/meetings/:meetingId/minutes", handler.UpdateMinutes)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}

	log.Printf("Association service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
