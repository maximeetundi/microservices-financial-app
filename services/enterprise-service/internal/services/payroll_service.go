package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayrollService struct {
	payrollRepo   *repository.PayrollRepository
	empRepo       *repository.EmployeeRepository
	salaryService *SalaryService
	entRepo       *repository.EnterpriseRepository
	notifClient   *NotificationClient
}

func NewPayrollService(
	pRepo *repository.PayrollRepository, 
	eRepo *repository.EmployeeRepository, 
	sService *SalaryService,
	entRepo *repository.EnterpriseRepository,
	notifClient *NotificationClient,
) *PayrollService {
	return &PayrollService{
		payrollRepo:   pRepo,
		empRepo:       eRepo,
		salaryService: sService,
		entRepo:       entRepo,
		notifClient:   notifClient,
	}
}

// PreparePayroll (Step 1): Calculate totals but don't execute
func (s *PayrollService) PreparePayroll(ctx context.Context, enterpriseID string) (*models.PayrollRun, error) {
	employees, err := s.empRepo.FindByEnterprise(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}

	entOID, _ := primitive.ObjectIDFromHex(enterpriseID)
	now := time.Now()

	run := &models.PayrollRun{
		EnterpriseID: entOID,
		PeriodMonth:  int(now.Month()),
		PeriodYear:   now.Year(),
		Status:       models.PayrollStatusDraft,
		Details:      []models.PayrollDetail{},
	}

	for _, emp := range employees {
		// Only pay active employees
		if emp.Status != models.EmployeeStatusActive {
			continue
		}

		// Recalculate to be sure
		s.salaryService.CalculateNetSalary(&emp.SalaryConfig)

		detail := models.PayrollDetail{
			EmployeeID:   emp.ID,
			EmployeeName: emp.FirstName + " " + emp.LastName,
			BaseSalary:   emp.SalaryConfig.BaseAmount,
			NetPay:       emp.SalaryConfig.NetPayable,
			Status:       "PENDING",
		}
		
		// Sum bonuses/deductions for reporting
		for _, b := range emp.SalaryConfig.Bonuses { detail.Bonuses += b.Amount }
		for _, d := range emp.SalaryConfig.Deductions { detail.Deductions += d.Amount }

		run.Details = append(run.Details, detail)
		run.TotalAmount += detail.NetPay
		run.TotalEmployees++
	}

	return run, nil
}

// ExecutePayroll (Step 2): Save run and Trigger Payment with Notifications
func (s *PayrollService) ExecutePayroll(ctx context.Context, run *models.PayrollRun) error {
	run.Status = models.PayrollStatusProcessing
	run.ExecutedAt = time.Now()

	// 1. Save initial record
	if err := s.payrollRepo.Create(ctx, run); err != nil {
		return err
	}

	// 2. Call Transfer Service (Simulated - TODO: connect to transfer-service)
	transactionID := "TX_MOCK_" + primitive.NewObjectID().Hex()
	success := true 

	if success {
		run.Status = models.PayrollStatusCompleted
		run.TransactionID = transactionID
		for i := range run.Details {
			run.Details[i].Status = "SUCCESS"
		}
		
		// 3. Send Notifications
		s.sendPayrollNotifications(ctx, run)
	} else {
		run.Status = models.PayrollStatusFailed
	}
	
	return nil
}

// sendPayrollNotifications sends notifications to owner and employees
func (s *PayrollService) sendPayrollNotifications(ctx context.Context, run *models.PayrollRun) {
	if s.notifClient == nil || s.entRepo == nil {
		return
	}
	
	// Get enterprise info
	ent, err := s.entRepo.FindByID(ctx, run.EnterpriseID.Hex())
	if err != nil {
		log.Printf("Failed to get enterprise for notifications: %v", err)
		return
	}
	
	period := fmt.Sprintf("%02d/%d", run.PeriodMonth, run.PeriodYear)
	
	// Notify owner
	if err := s.notifClient.NotifyPayrollExecution(ctx, ent.OwnerID, run.TotalAmount, run.TotalEmployees, period); err != nil {
		log.Printf("Failed to send payroll notification to owner: %v", err)
	}
	
	// Notify each employee
	employees, err := s.empRepo.FindByEnterprise(ctx, run.EnterpriseID.Hex())
	if err != nil {
		log.Printf("Failed to get employees for notifications: %v", err)
		return
	}
	
	// Create map of employee payments
	paymentMap := make(map[string]float64)
	for _, detail := range run.Details {
		paymentMap[detail.EmployeeID.Hex()] = detail.NetPay
	}
	
	for _, emp := range employees {
		if netPay, ok := paymentMap[emp.ID.Hex()]; ok && emp.UserID != "" {
			if err := s.notifClient.NotifySalaryPayment(ctx, emp.UserID, ent.Name, netPay, period); err != nil {
				log.Printf("Failed to send salary notification to %s: %v", emp.Email, err)
			}
		}
	}
}

// ListPayrollRuns returns payroll history for an enterprise
func (s *PayrollService) ListPayrollRuns(ctx context.Context, enterpriseID string, year int) ([]models.PayrollRun, error) {
	return s.payrollRepo.FindByEnterpriseAndYear(ctx, enterpriseID, year)
}

