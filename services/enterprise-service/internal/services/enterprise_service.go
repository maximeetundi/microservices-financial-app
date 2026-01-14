package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

// EnterpriseService handles business logic for enterprise management
type EnterpriseService struct {
	entRepo      *repository.EnterpriseRepository
	empRepo      *repository.EmployeeRepository
	walletClient *WalletClient
	notifClient  *NotificationClient
}

func NewEnterpriseService(
	entRepo *repository.EnterpriseRepository,
	empRepo *repository.EmployeeRepository,
	walletClient *WalletClient,
	notifClient *NotificationClient,
) *EnterpriseService {
	return &EnterpriseService{
		entRepo:      entRepo,
		empRepo:      empRepo,
		walletClient: walletClient,
		notifClient:  notifClient,
	}
}

// CreateEnterprise creates a new enterprise with:
// 1. Default wallet (auto-created)
// 2. Owner as first admin employee
// 3. Default security policies
func (s *EnterpriseService) CreateEnterprise(ctx context.Context, ent *models.Enterprise, ownerName, ownerEmail string) error {
	// Set default type if missing
	if ent.Type == "" {
		ent.Type = models.EnterpriseTypeService
	}
	
	// Set default currency
	if ent.Settings.DefaultCurrency == "" {
		ent.Settings.DefaultCurrency = "XOF"
	}
	
	// Set default security policies
	ent.SecurityPolicies = models.GetDefaultSecurityPolicies()
	
	// Create enterprise
	if err := s.entRepo.Create(ctx, ent); err != nil {
		return err
	}
	
	// Create default wallet for enterprise
	if s.walletClient != nil {
		wallet, err := s.walletClient.CreateWallet(ctx, ent.OwnerID, ent.Settings.DefaultCurrency)
		if err != nil {
			log.Printf("Warning: Failed to create default wallet for enterprise %s: %v", ent.ID.Hex(), err)
		} else {
			ent.DefaultWalletID = wallet.ID
			ent.WalletIDs = []string{wallet.ID}
			// Update enterprise with wallet info
			s.entRepo.Update(ctx, ent)
		}
	}
	
	// Create owner as first admin employee
	owner := &models.Employee{
		EnterpriseID: ent.ID,
		UserID:       ent.OwnerID,
		FirstName:    ownerName,
		LastName:     "",
		Email:        ownerEmail,
		Profession:   "Directeur",
		Role:         models.EmployeeRoleOwner,
		Status:       models.EmployeeStatusActive, // Owner is automatically active
		AcceptedAt:   time.Now(),
		InvitedAt:    time.Now(),
	}
	
	// Parse name if it contains space
	if spaceIdx := indexOf(ownerName, " "); spaceIdx > 0 {
		owner.FirstName = ownerName[:spaceIdx]
		owner.LastName = ownerName[spaceIdx+1:]
	}
	
	if err := s.empRepo.Create(ctx, owner); err != nil {
		log.Printf("Warning: Failed to create owner employee record for enterprise %s: %v", ent.ID.Hex(), err)
	}
	
	return nil
}

// AddWallet creates an additional wallet for the enterprise
func (s *EnterpriseService) AddWallet(ctx context.Context, enterpriseID, currency string) (*WalletInfo, error) {
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}
	
	if s.walletClient == nil {
		return nil, nil
	}
	
	wallet, err := s.walletClient.CreateWallet(ctx, ent.OwnerID, currency)
	if err != nil {
		return nil, err
	}
	
	// Add to enterprise wallet list
	ent.WalletIDs = append(ent.WalletIDs, wallet.ID)
	if err := s.entRepo.Update(ctx, ent); err != nil {
		return nil, err
	}
	
	return wallet, nil
}

// GetWallets returns all wallets for an enterprise
func (s *EnterpriseService) GetWallets(ctx context.Context, enterpriseID string) ([]WalletInfo, error) {
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}
	
	if s.walletClient == nil {
		return nil, nil
	}
	
	// Fetch wallet details
	var wallets []WalletInfo
	for _, walletID := range ent.WalletIDs {
		wallet, err := s.walletClient.GetWallet(ctx, walletID)
		if err != nil {
			log.Printf("Warning: Failed to get wallet %s: %v", walletID, err)
			continue
		}
		wallets = append(wallets, *wallet)
	}
	
	return wallets, nil
}

// SetEmployeeAsAdmin promotes an employee to admin role
func (s *EnterpriseService) SetEmployeeAsAdmin(ctx context.Context, enterpriseID, employeeID string) error {
	emp, err := s.empRepo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}
	
	if emp.EnterpriseID.Hex() != enterpriseID {
		return errors.New("employee does not belong to this enterprise")
	}
	
	// Add career event
	emp.History = append(emp.History, models.CareerEvent{
		Date:        time.Now(),
		Type:        "ROLE_CHANGE",
		Description: "Promu administrateur",
		Previous:    map[string]interface{}{"role": emp.Role},
		New:         map[string]interface{}{"role": models.EmployeeRoleAdmin},
	})
	
	emp.Role = models.EmployeeRoleAdmin
	return s.empRepo.Update(ctx, emp)
}

// RemoveAdminRole removes admin role from an employee (cannot remove owner)
func (s *EnterpriseService) RemoveAdminRole(ctx context.Context, enterpriseID, employeeID string) error {
	emp, err := s.empRepo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}
	
	if emp.EnterpriseID.Hex() != enterpriseID {
		return errors.New("employee does not belong to this enterprise")
	}
	
	if emp.Role == models.EmployeeRoleOwner {
		return errors.New("cannot remove owner role")
	}
	
	// Add career event
	emp.History = append(emp.History, models.CareerEvent{
		Date:        time.Now(),
		Type:        "ROLE_CHANGE",
		Description: "Retrait des droits administrateur",
		Previous:    map[string]interface{}{"role": emp.Role},
		New:         map[string]interface{}{"role": models.EmployeeRoleStandard},
	})
	
	emp.Role = models.EmployeeRoleStandard
	return s.empRepo.Update(ctx, emp)
}

// UpdateSecurityPolicies updates the security policies for an enterprise
func (s *EnterpriseService) UpdateSecurityPolicies(ctx context.Context, enterpriseID string, policies []models.SecurityPolicy) error {
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return err
	}
	
	ent.SecurityPolicies = policies
	return s.entRepo.Update(ctx, ent)
}

// GetAdmins returns all admin employees for an enterprise
func (s *EnterpriseService) GetAdmins(ctx context.Context, enterpriseID string) ([]models.Employee, error) {
	employees, err := s.empRepo.FindByEnterprise(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}
	
	var admins []models.Employee
	for _, emp := range employees {
		if emp.IsAdmin() && emp.Status == models.EmployeeStatusActive {
			admins = append(admins, emp)
		}
	}
	return admins, nil
}

// helper function
func indexOf(s string, substr string) int {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
