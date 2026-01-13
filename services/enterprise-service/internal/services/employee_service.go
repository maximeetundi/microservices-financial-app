package services

import (
	"context"
	"errors"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

type EmployeeService struct {
	repo          *repository.EmployeeRepository
	salaryService *SalaryService
	authClient    *AuthClient
}

func NewEmployeeService(repo *repository.EmployeeRepository, salaryService *SalaryService, authClient *AuthClient) *EmployeeService {
	return &EmployeeService{
		repo:          repo,
		salaryService: salaryService,
		authClient:    authClient,
	}
}

// ListByEnterprise (Point 2): Get all employees for an enterprise
func (s *EmployeeService) ListByEnterprise(ctx context.Context, enterpriseID string) ([]models.Employee, error) {
	return s.repo.FindByEnterprise(ctx, enterpriseID)
}

// InviteEmployee (Point 2): Adds employee & Sends notification (mocked here)
func (s *EmployeeService) InviteEmployee(ctx context.Context, emp *models.Employee) error {
	// 1. Calculate initial salary
	s.salaryService.CalculateNetSalary(&emp.SalaryConfig)
	
	// 2. Set as Pending
	emp.Status = models.EmployeeStatusPending
	emp.InvitedAt = time.Now()
	
	// TODO: Send Notification (Email/SMS)
	
	return s.repo.Create(ctx, emp)
}

// ConfirmEmployee (Point 2): Validates PIN and activates
func (s *EmployeeService) ConfirmEmployee(ctx context.Context, employeeID string, pin string) error {
	emp, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}

	if emp.Status == models.EmployeeStatusActive {
		return errors.New("employee already active")
	}

	// Verify PIN with Auth Service
	// We need UserID. Assuming the employee has already been linked to a UserID via invite or claims.
	// If UserID is empty, we cannot verify global PIN.
	if emp.UserID == "" {
		return errors.New("employee record has no associated user account")
	}

	valid, err := s.authClient.VerifyPin(emp.UserID, pin, "") // Token passed as empty for now, or need to propagate
	if err != nil {
		return err // "PIN verification failed" or service error
	}
	if !valid {
		return errors.New("invalid security PIN")
	}

	emp.Status = models.EmployeeStatusActive
	emp.AcceptedAt = time.Now()
	
	return s.repo.Update(ctx, emp)
}

// PromoteEmployee (Point 4): Updates position/salary and Logs History
func (s *EmployeeService) PromoteEmployee(ctx context.Context, employeeID string, newRole string, newSalary *models.SalaryConfig) error {
	emp, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}

	// 1. Log previous state to History
	historyEvent := models.CareerEvent{
		Date:        time.Now(),
		Type:        "PROMOTION",
		Description: "Promoted from " + emp.Profession + " to " + newRole,
		Previous:    map[string]interface{}{"role": emp.Profession, "salary": emp.SalaryConfig},
		New:         map[string]interface{}{"role": newRole, "salary": newSalary},
	}
	emp.History = append(emp.History, historyEvent)

	// 2. Apply new state
	emp.Profession = newRole
	if newSalary != nil {
		s.salaryService.CalculateNetSalary(newSalary)
		emp.SalaryConfig = *newSalary
	}

	return s.repo.Update(ctx, emp)
}
