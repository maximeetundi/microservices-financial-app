package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

type EmployeeService struct {
	repo          *repository.EmployeeRepository
	salaryService *SalaryService
	authClient    *AuthClient
	entRepo       *repository.EnterpriseRepository
	notifClient   *NotificationClient
}

func NewEmployeeService(
	repo *repository.EmployeeRepository, 
	salaryService *SalaryService, 
	authClient *AuthClient,
	entRepo *repository.EnterpriseRepository,
	notifClient *NotificationClient,
) *EmployeeService {
	return &EmployeeService{
		repo:          repo,
		salaryService: salaryService,
		authClient:    authClient,
		entRepo:       entRepo,
		notifClient:   notifClient,
	}
}

// ListByEnterprise (Point 2): Get all employees for an enterprise
func (s *EmployeeService) ListByEnterprise(ctx context.Context, enterpriseID string) ([]models.Employee, error) {
	return s.repo.FindByEnterprise(ctx, enterpriseID)
}

// GetEmployeeByUserAndEnterprise finds an employee by user ID and enterprise ID
func (s *EmployeeService) GetEmployeeByUserAndEnterprise(ctx context.Context, userID, enterpriseID string) (*models.Employee, error) {
	employees, err := s.repo.FindByEnterprise(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}
	
	for _, emp := range employees {
		if emp.UserID == userID && emp.Status == models.EmployeeStatusActive {
			return &emp, nil
		}
	}
	
	return nil, errors.New("employee not found for this user in this enterprise")
}

// InviteEmployee (Point 2): Adds employee & Sends notification
func (s *EmployeeService) InviteEmployee(ctx context.Context, emp *models.Employee) error {
	// 1. Calculate initial salary
	s.salaryService.CalculateNetSalary(&emp.SalaryConfig)
	
	// 2. Try to find the user by email or phone to link accounts
	if emp.UserID == "" && s.authClient != nil {
		userID, err := s.authClient.FindUserByContact(emp.Email, emp.PhoneNumber)
		if err == nil && userID != "" {
			emp.UserID = userID
			log.Printf("Found existing user %s for employee %s %s", userID, emp.FirstName, emp.LastName)
		}
	}
	
	// 3. Set as Pending
	emp.Status = models.EmployeeStatusPending
	emp.InvitedAt = time.Now()
	
	// 4. Create employee record
	if err := s.repo.Create(ctx, emp); err != nil {
		return err
	}
	
	// 5. Send Notification to employee (if they have a user account)
	if emp.UserID != "" && s.notifClient != nil && s.entRepo != nil {
		// Fetch enterprise name for notification
		ent, err := s.entRepo.FindByID(ctx, emp.EnterpriseID.Hex())
		if err == nil {
			if err := s.notifClient.NotifyEmployeeInvitation(ctx, emp.UserID, ent.Name); err != nil {
				log.Printf("Failed to send invitation notification: %v", err)
			} else {
				log.Printf("Sent invitation notification to user %s for enterprise %s", emp.UserID, ent.Name)
			}
		}
	} else {
		log.Printf("No user account found for employee %s - notification not sent", emp.Email)
	}
	
	return nil
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
	if emp.UserID == "" {
		return errors.New("employee record has no associated user account")
	}

	valid, err := s.authClient.VerifyPin(emp.UserID, pin, "")
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid security PIN")
	}

	emp.Status = models.EmployeeStatusActive
	emp.AcceptedAt = time.Now()
	
	if err := s.repo.Update(ctx, emp); err != nil {
		return err
	}
	
	// Send notification to enterprise owner
	if s.notifClient != nil && s.entRepo != nil {
		ent, err := s.entRepo.FindByID(ctx, emp.EnterpriseID.Hex())
		if err == nil {
			employeeName := emp.FirstName + " " + emp.LastName
			if err := s.notifClient.NotifyEmployeeAcceptance(ctx, ent.OwnerID, employeeName, ent.Name); err != nil {
				log.Printf("Failed to send acceptance notification: %v", err)
			}
		}
	}
	
	return nil
}

// PromoteEmployee (Point 4): Updates position/salary and Logs History
func (s *EmployeeService) PromoteEmployee(ctx context.Context, employeeID string, newRole string, newSalary *models.SalaryConfig, permissions *models.AdminPermission) error {
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
		New:         map[string]interface{}{"role": newRole, "salary": newSalary, "permissions": permissions},
	}
	emp.History = append(emp.History, historyEvent)

	// 2. Apply new state
    if newRole != "" {
	    emp.Profession = newRole
        if newRole == "ADMIN" {
            emp.Role = models.EmployeeRoleAdmin
        }
    }
    
    // Apply permissions if provided
    if permissions != nil {
        emp.Permissions = *permissions
    }
    
	if newSalary != nil {
		s.salaryService.CalculateNetSalary(newSalary)
		emp.SalaryConfig = *newSalary
	}

	return s.repo.Update(ctx, emp)
}

// TerminateEmployee (Point 4): Terminates employee and stops payments
func (s *EmployeeService) TerminateEmployee(ctx context.Context, enterpriseID, employeeID string) error {
	emp, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}

	if emp.EnterpriseID.Hex() != enterpriseID {
		return errors.New("employee does not belong to this enterprise")
	}

	// Log termination event
	historyEvent := models.CareerEvent{
		Date:        time.Now(),
		Type:        "TERMINATION",
		Description: "Employee terminated",
		Previous:    map[string]interface{}{"status": emp.Status},
		New:         map[string]interface{}{"status": models.EmployeeStatusTerminated},
	}
	emp.History = append(emp.History, historyEvent)
	emp.Status = models.EmployeeStatusTerminated

	return s.repo.Update(ctx, emp)
}

