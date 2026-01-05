package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/handlers"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/services"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration", Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}}, []string{"method", "path", "status"})
	contributionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{Name: "contributions_total", Help: "Total contributions"}, []string{"association", "status"})
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" { c.Next(); return }
		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath(); if path == "" { path = c.Request.URL.Path }
		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(time.Since(start).Seconds())
	}
}

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

	// New repositories for extended features
	roleRepo := repository.NewRoleRepository(db)
	approvalRepo := repository.NewApprovalRepository(db)
	chatRepo := repository.NewChatRepository(db)
	solidarityRepo := repository.NewSolidarityRepository(db)
	calledRepo := repository.NewCalledRoundRepository(db)

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
	extHandler := handlers.NewExtendedHandler(db.DB, roleRepo, approvalRepo, chatRepo, solidarityRepo, calledRepo, memberRepo)

	// Setup Gin router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), prometheusMiddleware())

	// CORS - Note: Cannot use wildcard "*" with AllowCredentials: true
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://app.maximeetundi.store", "https://admin.maximeetundi.store", "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Authorization", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "association-service"}) })

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

			// === NEW FEATURES ===

			// Custom Roles
			protected.POST("/associations/:id/roles", extHandler.CreateRole)
			protected.GET("/associations/:id/roles", extHandler.GetRoles)
			protected.DELETE("/associations/:id/roles/:roleId", extHandler.DeleteRole)

			// Multi-Signature Approvers
			protected.POST("/associations/:id/approvers", extHandler.SetApprovers)
			protected.GET("/associations/:id/approvers", extHandler.GetApprovers)
			protected.GET("/associations/:id/approvals", extHandler.GetPendingApprovals)
			protected.POST("/approvals/:requestId/vote", extHandler.VoteOnApproval)

			// Chat
			protected.POST("/associations/:id/chat", extHandler.SendMessage)
			protected.GET("/associations/:id/chat", extHandler.GetMessages)

			// Solidarity Events
			protected.POST("/associations/:id/solidarity", extHandler.CreateSolidarityEvent)
			protected.GET("/associations/:id/solidarity", extHandler.GetSolidarityEvents)
			protected.POST("/solidarity/:eventId/contribute", extHandler.ContributeToSolidarity)

			// Called Tontine
			protected.POST("/associations/:id/rounds", extHandler.CreateCalledRound)
			protected.GET("/associations/:id/rounds", extHandler.GetCalledRounds)
			protected.POST("/rounds/:roundId/pledge", extHandler.MakePledge)
			protected.POST("/rounds/:roundId/pledges/:pledgeId/pay", extHandler.PayPledge)
			protected.GET("/rounds/:roundId/pledges", extHandler.GetPledges)

			// Emergency Fund (Caisse de Secours)
			protected.POST("/associations/:id/emergency-fund", extHandler.CreateEmergencyFund)
			protected.GET("/associations/:id/emergency-fund", extHandler.GetEmergencyFund)
			protected.POST("/associations/:id/emergency-fund/contribute", extHandler.ContributeToEmergencyFund)
			protected.POST("/associations/:id/emergency-fund/withdraw", extHandler.RequestEmergencyWithdrawal)
			protected.GET("/associations/:id/emergency-fund/withdrawals", extHandler.GetEmergencyWithdrawals)
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
