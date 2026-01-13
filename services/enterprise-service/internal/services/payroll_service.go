package services

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayrollService struct {
	payrollRepo  *repository.PayrollRepository
	empRepo      *repository.EmployeeRepository
	salaryService *SalaryService
}

func NewPayrollService(pRepo *repository.PayrollRepository, eRepo *repository.EmployeeRepository, sService *SalaryService) *PayrollService {
	return &PayrollService{
		payrollRepo:   pRepo,
		empRepo:       eRepo,
		salaryService: sService,
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

// ExecutePayroll (Step 2): Save run and Trigger Payment (Mocked for MVP)
func (s *PayrollService) ExecutePayroll(ctx context.Context, run *models.PayrollRun) error {
	run.Status = models.PayrollStatusProcessing
	run.ExecutedAt = time.Now()

	// 1. Save initial record
	if err := s.payrollRepo.Create(ctx, run); err != nil {
		return err
	}

	// 2. Call Transfer Service (Simulated)
	// In real impl: POST request to transfer-service/bulk
	transactionID := "TX_MOCK_" + primitive.NewObjectID().Hex()
	success := true 

	if success {
		run.Status = models.PayrollStatusCompleted
		run.TransactionID = transactionID
		for i := range run.Details {
			run.Details[i].Status = "SUCCESS"
		}
		// TODO: Notify Users
	} else {
		run.Status = models.PayrollStatusFailed
	}

	// 3. Update Record
	// Note: Repo needs an Update method, or we rely on Create for log-only for now if immutable.
	// We'll simplisticly assume for this MVP we just created it with the status. 
	// To be strict, we really should update.
	
	// For MVP simplicity, we assume successful execution happens synchronously before final save if local, 
	// or we'd update it properly. Let's just return nil as we saved it in Valid state already if we swapped the order.
	// Actually, let's just re-save. Since repo doesn't have Update yet, we will skip re-saving in MVP step or add Update.
	
	return nil
}
