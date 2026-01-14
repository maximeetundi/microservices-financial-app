package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

// ApprovalService handles multi-admin approval workflows
type ApprovalService struct {
	approvalRepo *repository.ApprovalRepository
	empRepo      *repository.EmployeeRepository
	entRepo      *repository.EnterpriseRepository
	notifClient  *NotificationClient
}

func NewApprovalService(
	approvalRepo *repository.ApprovalRepository,
	empRepo *repository.EmployeeRepository,
	entRepo *repository.EnterpriseRepository,
	notifClient *NotificationClient,
) *ApprovalService {
	return &ApprovalService{
		approvalRepo: approvalRepo,
		empRepo:      empRepo,
		entRepo:      entRepo,
		notifClient:  notifClient,
	}
}

// InitiateAction creates a new approval request for a sensitive action
func (s *ApprovalService) InitiateAction(
	ctx context.Context,
	enterpriseID string,
	initiatorEmployee *models.Employee,
	actionType string,
	actionName string,
	description string,
	payload map[string]interface{},
	amount float64,
	currency string,
) (*models.ActionApproval, error) {
	// Get enterprise
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return nil, errors.New("enterprise not found")
	}
	
	// Check if initiator is admin
	if !initiatorEmployee.IsAdmin() {
		return nil, errors.New("only admins can initiate actions requiring approval")
	}
	
	// Find applicable policy
	policy := s.findApplicablePolicy(ent, actionType, amount)
	if policy == nil || !policy.Enabled {
		// No policy or disabled - action can proceed without approval
		return nil, nil
	}
	
	// Count total admins
	admins, err := s.getEnterpriseAdmins(ctx, enterpriseID)
	if err != nil {
		return nil, err
	}
	totalAdmins := len(admins)
	
	// Calculate required approvals
	requiredApprovals := policy.MinApprovals
	if policy.RequireMajority {
		requiredApprovals = totalAdmins/2 + 1
	}
	if policy.RequireAllAdmins {
		requiredApprovals = totalAdmins
	}
	
	// Set expiration
	expirationHours := policy.ExpirationHours
	if expirationHours <= 0 {
		expirationHours = 24
	}
	
	approval := &models.ActionApproval{
		EnterpriseID:      ent.ID,
		ActionType:        actionType,
		ActionName:        actionName,
		Description:       description,
		Payload:           payload,
		Amount:            amount,
		Currency:          currency,
		RequiredApprovals: requiredApprovals,
		RequireMajority:   policy.RequireMajority,
		RequireAllAdmins:  policy.RequireAllAdmins,
		TotalAdmins:       totalAdmins,
		Status:            models.ApprovalStatusPending,
		Approvals: []models.AdminApproval{
			// Initiator automatically approves
			{
				AdminEmployeeID: initiatorEmployee.ID.Hex(),
				AdminUserID:     initiatorEmployee.UserID,
				AdminName:       initiatorEmployee.FirstName + " " + initiatorEmployee.LastName,
				Decision:        "APPROVED",
				DecidedAt:       time.Now(),
			},
		},
		InitiatedBy:     initiatorEmployee.ID.Hex(),
		InitiatorUserID: initiatorEmployee.UserID,
		InitiatorName:   initiatorEmployee.FirstName + " " + initiatorEmployee.LastName,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(time.Duration(expirationHours) * time.Hour),
	}
	
	// Check if already approved (single admin case)
	if approval.IsApproved() {
		approval.Status = models.ApprovalStatusApproved
	}
	
	// Save
	if err := s.approvalRepo.Create(ctx, approval); err != nil {
		return nil, err
	}
	
	// Notify other admins if still pending
	if approval.Status == models.ApprovalStatusPending {
		s.notifyAdminsOfPendingApproval(ctx, ent, admins, approval, initiatorEmployee.UserID)
	}
	
	return approval, nil
}

// ApproveAction adds an admin's approval to an action
func (s *ApprovalService) ApproveAction(ctx context.Context, approvalID string, admin *models.Employee) error {
	approval, err := s.approvalRepo.FindByID(ctx, approvalID)
	if err != nil {
		return errors.New("approval not found")
	}
	
	if approval.Status != models.ApprovalStatusPending {
		return errors.New("action is not pending approval")
	}
	
	if approval.IsExpired() {
		s.approvalRepo.UpdateStatus(ctx, approvalID, models.ApprovalStatusExpired)
		return errors.New("approval has expired")
	}
	
	if !admin.IsAdmin() {
		return errors.New("only admins can approve actions")
	}
	
	if approval.HasAdminVoted(admin.UserID) {
		return errors.New("you have already voted on this action")
	}
	
	// Add approval
	adminApproval := models.AdminApproval{
		AdminEmployeeID: admin.ID.Hex(),
		AdminUserID:     admin.UserID,
		AdminName:       admin.FirstName + " " + admin.LastName,
		Decision:        "APPROVED",
		DecidedAt:       time.Now(),
	}
	
	if err := s.approvalRepo.AddApproval(ctx, approvalID, adminApproval); err != nil {
		return err
	}
	
	// Refresh and check if now approved
	approval, _ = s.approvalRepo.FindByID(ctx, approvalID)
	if approval.IsApproved() {
		s.approvalRepo.UpdateStatus(ctx, approvalID, models.ApprovalStatusApproved)
	}
	
	return nil
}

// RejectAction adds an admin's rejection to an action
func (s *ApprovalService) RejectAction(ctx context.Context, approvalID string, admin *models.Employee, reason string) error {
	approval, err := s.approvalRepo.FindByID(ctx, approvalID)
	if err != nil {
		return errors.New("approval not found")
	}
	
	if approval.Status != models.ApprovalStatusPending {
		return errors.New("action is not pending approval")
	}
	
	if !admin.IsAdmin() {
		return errors.New("only admins can reject actions")
	}
	
	if approval.HasAdminVoted(admin.UserID) {
		return errors.New("you have already voted on this action")
	}
	
	// Add rejection
	adminApproval := models.AdminApproval{
		AdminEmployeeID: admin.ID.Hex(),
		AdminUserID:     admin.UserID,
		AdminName:       admin.FirstName + " " + admin.LastName,
		Decision:        "REJECTED",
		Reason:          reason,
		DecidedAt:       time.Now(),
	}
	
	if err := s.approvalRepo.AddApproval(ctx, approvalID, adminApproval); err != nil {
		return err
	}
	
	// Refresh and check if now rejected
	approval, _ = s.approvalRepo.FindByID(ctx, approvalID)
	if approval.IsRejected() {
		s.approvalRepo.UpdateStatus(ctx, approvalID, models.ApprovalStatusRejected)
	}
	
	return nil
}

// GetPendingApprovals returns pending approvals for an enterprise
func (s *ApprovalService) GetPendingApprovals(ctx context.Context, enterpriseID string) ([]models.ActionApproval, error) {
	return s.approvalRepo.FindPendingByEnterprise(ctx, enterpriseID)
}

// MarkAsExecuted marks an approved action as executed
func (s *ApprovalService) MarkAsExecuted(ctx context.Context, approvalID string) error {
	return s.approvalRepo.UpdateStatus(ctx, approvalID, models.ApprovalStatusExecuted)
}

// RequiresApproval checks if an action requires multi-approval
func (s *ApprovalService) RequiresApproval(ctx context.Context, enterpriseID, actionType string, amount float64) (bool, error) {
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil {
		return false, err
	}
	
	policy := s.findApplicablePolicy(ent, actionType, amount)
	return policy != nil && policy.Enabled, nil
}

// findApplicablePolicy finds the security policy that applies to an action
func (s *ApprovalService) findApplicablePolicy(ent *models.Enterprise, actionType string, amount float64) *models.SecurityPolicy {
	for _, policy := range ent.SecurityPolicies {
		if policy.ActionType == actionType && policy.Enabled {
			// Check threshold if applicable
			if policy.ThresholdAmount > 0 && amount < policy.ThresholdAmount {
				continue
			}
			return &policy
		}
	}
	return nil
}

// getEnterpriseAdmins returns all admins for an enterprise
func (s *ApprovalService) getEnterpriseAdmins(ctx context.Context, enterpriseID string) ([]models.Employee, error) {
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

// notifyAdminsOfPendingApproval sends notifications to admins
func (s *ApprovalService) notifyAdminsOfPendingApproval(
	ctx context.Context,
	ent *models.Enterprise,
	admins []models.Employee,
	approval *models.ActionApproval,
	initiatorUserID string,
) {
	if s.notifClient == nil {
		return
	}
	
	for _, admin := range admins {
		if admin.UserID == initiatorUserID {
			continue // Don't notify initiator
		}
		
		if err := s.notifClient.NotifyUser(ctx, admin.UserID, "enterprise.approval.pending",
			"Approbation requise",
			approval.ActionName+" requiert votre approbation pour "+ent.Name,
			map[string]interface{}{
				"enterprise_id": ent.ID.Hex(),
				"approval_id":   approval.ID.Hex(),
				"action_type":   approval.ActionType,
				"amount":        approval.Amount,
				"action":        "approve_action",
			},
		); err != nil {
			log.Printf("Failed to send approval notification: %v", err)
		}
	}
}
