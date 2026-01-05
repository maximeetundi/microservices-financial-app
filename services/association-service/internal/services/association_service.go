package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/repository"
)

type AssociationService struct {
	assocRepo   *repository.AssociationRepository
	memberRepo  *repository.MemberRepository
	treasuryRepo *repository.TreasuryRepository
	meetingRepo *repository.MeetingRepository
	loanRepo    *repository.LoanRepository
	walletServiceURL string
}

func NewAssociationService(
	assocRepo *repository.AssociationRepository,
	memberRepo *repository.MemberRepository,
	treasuryRepo *repository.TreasuryRepository,
	meetingRepo *repository.MeetingRepository,
	loanRepo *repository.LoanRepository,
) *AssociationService {
	walletURL := os.Getenv("WALLET_SERVICE_URL")
	if walletURL == "" {
		walletURL = "http://wallet-service:8082"
	}
	return &AssociationService{
		assocRepo:   assocRepo,
		memberRepo:  memberRepo,
		treasuryRepo: treasuryRepo,
		meetingRepo: meetingRepo,
		loanRepo:    loanRepo,
		walletServiceURL: walletURL,
	}
}

// === Association CRUD ===

func (s *AssociationService) CreateAssociation(userID string, req *models.CreateAssociationRequest) (*models.Association, error) {
	currency := req.Currency
	if currency == "" {
		currency = "XOF"
	}

	assoc := &models.Association{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Rules:       req.Rules,
		Currency:    currency,
		Status:      models.AssociationStatusActive,
		CreatedBy:   userID,
	}

	if err := s.assocRepo.Create(assoc); err != nil {
		return nil, fmt.Errorf("failed to create association: %w", err)
	}

	// Determine creator's role (default: president)
	creatorRole := req.CreatorRole
	if creatorRole == "" {
		creatorRole = models.MemberRolePresident
	}

	// Add creator as member with chosen role
	member := &models.Member{
		AssociationID: assoc.ID,
		UserID:        userID,
		Role:          creatorRole,
		Status:        models.MemberStatusActive,
	}
	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	s.assocRepo.UpdateMemberCount(assoc.ID, 1)
	assoc.TotalMembers = 1

	return assoc, nil
}

func (s *AssociationService) GetAssociation(id string) (*models.Association, error) {
	return s.assocRepo.GetByID(id)
}

func (s *AssociationService) GetMyAssociations(userID string, limit, offset int) ([]*models.Association, error) {
	return s.assocRepo.GetByUser(userID, limit, offset)
}

func (s *AssociationService) GetAllAssociations(limit, offset int) ([]*models.Association, error) {
	return s.assocRepo.GetAll(limit, offset)
}

func (s *AssociationService) UpdateAssociation(id, userID string, req *models.UpdateAssociationRequest) (*models.Association, error) {
	assoc, err := s.assocRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("association not found")
	}

	// Check permission (creator or president)
	if !s.hasAdminPermission(assoc, userID) {
		return nil, fmt.Errorf("not authorized")
	}

	if req.Name != nil {
		assoc.Name = *req.Name
	}
	if req.Description != nil {
		assoc.Description = *req.Description
	}
	if req.Rules != nil {
		assoc.Rules = *req.Rules
	}
	if req.Status != nil {
		assoc.Status = models.AssociationStatus(*req.Status)
	}

	if err := s.assocRepo.Update(assoc); err != nil {
		return nil, err
	}

	return assoc, nil
}

func (s *AssociationService) DeleteAssociation(id, userID string) error {
	assoc, err := s.assocRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	if assoc.CreatedBy != userID {
		return fmt.Errorf("only creator can delete association")
	}

	return s.assocRepo.Delete(id)
}

// === Members ===

func (s *AssociationService) JoinAssociation(associationID, userID, message string) (*models.Member, error) {
	// Check if already a member
	existing, _ := s.memberRepo.GetByAssociationAndUser(associationID, userID)
	if existing != nil {
		if existing.Status == models.MemberStatusActive {
			return nil, fmt.Errorf("already a member")
		}
		// Reactivate
		s.memberRepo.UpdateStatus(existing.ID, models.MemberStatusActive)
		s.assocRepo.UpdateMemberCount(associationID, 1)
		existing.Status = models.MemberStatusActive
		return existing, nil
	}

	member := &models.Member{
		AssociationID: associationID,
		UserID:        userID,
		Role:          models.MemberRoleMember,
		Status:        models.MemberStatusActive, // Auto-approve for now
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, fmt.Errorf("failed to join: %w", err)
	}

	s.assocRepo.UpdateMemberCount(associationID, 1)

	return member, nil
}

func (s *AssociationService) LeaveAssociation(associationID, userID string) error {
	member, err := s.memberRepo.GetByAssociationAndUser(associationID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if member.Role == models.MemberRolePresident {
		return fmt.Errorf("president cannot leave, transfer role first")
	}

	s.memberRepo.UpdateStatus(member.ID, models.MemberStatusLeft)
	s.assocRepo.UpdateMemberCount(associationID, -1)

	return nil
}

func (s *AssociationService) GetMembers(associationID string) ([]*models.Member, error) {
	return s.memberRepo.GetByAssociation(associationID)
}

func (s *AssociationService) UpdateMemberRole(associationID, targetUserID, requestingUserID string, role models.MemberRole) error {
	assoc, err := s.assocRepo.GetByID(associationID)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, requestingUserID) {
		return fmt.Errorf("not authorized")
	}

	member, err := s.memberRepo.GetByAssociationAndUser(associationID, targetUserID)
	if err != nil {
		return fmt.Errorf("member not found")
	}

	return s.memberRepo.UpdateRole(member.ID, role)
}

// === Treasury & Contributions ===

// PayContribution - Deducts from user's wallet and adds to treasury
func (s *AssociationService) PayContribution(associationID, userID, walletID, pin string, amount float64, period, description string) (*models.TreasuryTransaction, error) {
	// Get association
	assoc, err := s.assocRepo.GetByID(associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found")
	}

	// Get member
	member, err := s.memberRepo.GetByAssociationAndUser(associationID, userID)
	if err != nil {
		return nil, fmt.Errorf("not a member of this association")
	}

	if member.Status != models.MemberStatusActive {
		return nil, fmt.Errorf("membership not active")
	}

	// Deduct from wallet using wallet-service
	err = s.deductFromWallet(userID, walletID, pin, amount, assoc.Currency, 
		fmt.Sprintf("Cotisation %s - %s", assoc.Name, period))
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	// Create treasury transaction
	tx := &models.TreasuryTransaction{
		AssociationID: associationID,
		Type:          models.TransactionTypeContribution,
		Amount:        amount,
		FromMemberID:  member.ID,
		Description:   description,
		Status:        models.TransactionStatusCompleted,
		CreatedBy:     userID,
	}

	if err := s.treasuryRepo.CreateTransaction(tx); err != nil {
		// TODO: Refund wallet on failure
		return nil, fmt.Errorf("failed to record transaction: %w", err)
	}

	// Update balances
	s.assocRepo.UpdateTreasuryBalance(associationID, amount)
	s.memberRepo.UpdateContributions(member.ID, amount)

	return tx, nil
}

func (s *AssociationService) GetTreasuryReport(associationID string) (*models.TreasuryReport, error) {
	return s.treasuryRepo.GetReport(associationID)
}

func (s *AssociationService) GetTransactions(associationID string, limit, offset int) ([]*models.TreasuryTransaction, error) {
	return s.treasuryRepo.GetByAssociation(associationID, limit, offset)
}

// DistributeFunds - For tontine: distribute treasury to a member
func (s *AssociationService) DistributeFunds(associationID, userID, recipientMemberID, walletID string, amount float64) (*models.TreasuryTransaction, error) {
	assoc, err := s.assocRepo.GetByID(associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, userID) {
		return nil, fmt.Errorf("not authorized")
	}

	if assoc.TreasuryBalance < amount {
		return nil, fmt.Errorf("insufficient treasury balance")
	}

	recipient, err := s.memberRepo.GetByID(recipientMemberID)
	if err != nil {
		return nil, fmt.Errorf("recipient member not found")
	}

	// Credit recipient's wallet
	err = s.creditToWallet(recipient.UserID, walletID, amount, assoc.Currency,
		fmt.Sprintf("Distribution %s", assoc.Name))
	if err != nil {
		return nil, fmt.Errorf("distribution payment failed: %w", err)
	}

	// Create transaction
	tx := &models.TreasuryTransaction{
		AssociationID: associationID,
		Type:          models.TransactionTypeDistribution,
		Amount:        amount,
		ToMemberID:    recipientMemberID,
		Description:   fmt.Sprintf("Distribution à membre"),
		Status:        models.TransactionStatusCompleted,
		CreatedBy:     userID,
	}

	if err := s.treasuryRepo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	s.assocRepo.UpdateTreasuryBalance(associationID, -amount)

	return tx, nil
}

// === Loans ===

func (s *AssociationService) RequestLoan(associationID, userID string, req *models.LoanRequest) (*models.Loan, error) {
	member, err := s.memberRepo.GetByAssociationAndUser(associationID, userID)
	if err != nil {
		return nil, fmt.Errorf("not a member")
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, req.Duration, 0)

	loan := &models.Loan{
		AssociationID: associationID,
		BorrowerID:    member.ID,
		Amount:        req.Amount,
		InterestRate:  req.InterestRate,
		Duration:      req.Duration,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        models.LoanStatusPending,
	}

	if err := s.loanRepo.Create(loan); err != nil {
		return nil, fmt.Errorf("failed to create loan request: %w", err)
	}

	return loan, nil
}

func (s *AssociationService) ApproveLoan(loanID, approverID, walletID string, approve bool, reason string) error {
	loan, err := s.loanRepo.GetByID(loanID)
	if err != nil {
		return fmt.Errorf("loan not found")
	}

	assoc, err := s.assocRepo.GetByID(loan.AssociationID)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, approverID) {
		return fmt.Errorf("not authorized")
	}

	if !approve {
		return s.loanRepo.Approve(loanID, approverID, models.LoanStatusRejected)
	}

	// Check treasury balance
	if assoc.TreasuryBalance < loan.Amount {
		return fmt.Errorf("insufficient treasury balance")
	}

	// Get borrower
	borrower, err := s.memberRepo.GetByID(loan.BorrowerID)
	if err != nil {
		return fmt.Errorf("borrower not found")
	}

	// Credit borrower's wallet
	err = s.creditToWallet(borrower.UserID, walletID, loan.Amount, assoc.Currency,
		fmt.Sprintf("Prêt %s", assoc.Name))
	if err != nil {
		return fmt.Errorf("loan payment failed: %w", err)
	}

	// Update loan status
	if err := s.loanRepo.Approve(loanID, approverID, models.LoanStatusActive); err != nil {
		return err
	}

	// Create treasury transaction
	tx := &models.TreasuryTransaction{
		AssociationID: loan.AssociationID,
		Type:          models.TransactionTypeLoan,
		Amount:        loan.Amount,
		ToMemberID:    loan.BorrowerID,
		Description:   fmt.Sprintf("Prêt approuvé"),
		Status:        models.TransactionStatusCompleted,
		CreatedBy:     approverID,
	}
	s.treasuryRepo.CreateTransaction(tx)

	// Update balances
	s.assocRepo.UpdateTreasuryBalance(loan.AssociationID, -loan.Amount)
	s.memberRepo.UpdateLoansReceived(loan.BorrowerID, loan.Amount)

	return nil
}

func (s *AssociationService) RepayLoan(loanID, userID, walletID, pin string, amount float64) error {
	loan, err := s.loanRepo.GetByID(loanID)
	if err != nil {
		return fmt.Errorf("loan not found")
	}

	borrower, err := s.memberRepo.GetByID(loan.BorrowerID)
	if err != nil {
		return fmt.Errorf("borrower not found")
	}

	if borrower.UserID != userID {
		return fmt.Errorf("only borrower can repay")
	}

	assoc, err := s.assocRepo.GetByID(loan.AssociationID)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	// Deduct from wallet
	err = s.deductFromWallet(userID, walletID, pin, amount, assoc.Currency,
		fmt.Sprintf("Remboursement prêt %s", assoc.Name))
	if err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	// Record repayment
	repayments := loan.Repayments
	if repayments == nil {
		repayments = make(models.JSONB)
	}
	repaymentList, _ := repayments["list"].([]interface{})
	repaymentList = append(repaymentList, map[string]interface{}{
		"amount":  amount,
		"date":    time.Now().Format(time.RFC3339),
	})
	repayments["list"] = repaymentList
	repayments["total_paid"] = s.getTotalRepaid(repaymentList)

	s.loanRepo.UpdateRepayments(loanID, repayments)

	// Check if fully repaid
	totalWithInterest := loan.Amount * (1 + loan.InterestRate/100)
	if repayments["total_paid"].(float64) >= totalWithInterest {
		s.loanRepo.UpdateStatus(loanID, models.LoanStatusPaid)
	}

	// Treasury transaction
	tx := &models.TreasuryTransaction{
		AssociationID: loan.AssociationID,
		Type:          models.TransactionTypeRepayment,
		Amount:        amount,
		FromMemberID:  loan.BorrowerID,
		Description:   fmt.Sprintf("Remboursement prêt"),
		Status:        models.TransactionStatusCompleted,
		CreatedBy:     userID,
	}
	s.treasuryRepo.CreateTransaction(tx)

	s.assocRepo.UpdateTreasuryBalance(loan.AssociationID, amount)

	return nil
}

func (s *AssociationService) GetLoans(associationID string, limit, offset int) ([]*models.Loan, error) {
	return s.loanRepo.GetByAssociation(associationID, limit, offset)
}

// === Meetings ===

func (s *AssociationService) CreateMeeting(associationID, userID string, req *models.CreateMeetingRequest) (*models.Meeting, error) {
	assoc, err := s.assocRepo.GetByID(associationID)
	if err != nil {
		return nil, fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, userID) {
		return nil, fmt.Errorf("not authorized")
	}

	meeting := &models.Meeting{
		AssociationID: associationID,
		Title:         req.Title,
		Date:          req.Date,
		Location:      req.Location,
		Agenda:        req.Agenda,
		Status:        models.MeetingStatusScheduled,
		CreatedBy:     userID,
	}

	if err := s.meetingRepo.Create(meeting); err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}

	return meeting, nil
}

func (s *AssociationService) GetMeetings(associationID string, limit, offset int) ([]*models.Meeting, error) {
	return s.meetingRepo.GetByAssociation(associationID, limit, offset)
}

func (s *AssociationService) RecordAttendance(meetingID, userID string, attendance map[string]bool) error {
	meeting, err := s.meetingRepo.GetByID(meetingID)
	if err != nil {
		return fmt.Errorf("meeting not found")
	}

	assoc, err := s.assocRepo.GetByID(meeting.AssociationID)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, userID) {
		return fmt.Errorf("not authorized")
	}

	attendanceJSON := make(models.JSONB)
	for k, v := range attendance {
		attendanceJSON[k] = v
	}

	return s.meetingRepo.UpdateAttendance(meetingID, attendanceJSON)
}

func (s *AssociationService) UpdateMinutes(meetingID, userID, minutes string) error {
	meeting, err := s.meetingRepo.GetByID(meetingID)
	if err != nil {
		return fmt.Errorf("meeting not found")
	}

	assoc, err := s.assocRepo.GetByID(meeting.AssociationID)
	if err != nil {
		return fmt.Errorf("association not found")
	}

	if !s.hasAdminPermission(assoc, userID) {
		return fmt.Errorf("not authorized")
	}

	return s.meetingRepo.UpdateMinutes(meetingID, minutes)
}

// === Wallet Integration ===

func (s *AssociationService) deductFromWallet(userID, walletID, pin string, amount float64, currency, description string) error {
	// Call wallet-service to create internal transfer (deduction)
	payload := map[string]interface{}{
		"type":        "internal",
		"amount":      amount,
		"currency":    currency,
		"recipient":   "ASSOCIATION_TREASURY", // Special recipient for internal tracking
		"description": description,
		"pin":         pin,
		"wallet_id":   walletID,
	}

	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/wallets/%s/withdraw", s.walletServiceURL, walletID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("wallet service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if msg, ok := errResp["error"].(string); ok {
			return fmt.Errorf(msg)
		}
		return fmt.Errorf("wallet deduction failed with status %d", resp.StatusCode)
	}

	return nil
}

func (s *AssociationService) creditToWallet(userID, walletID string, amount float64, currency, description string) error {
	// Call wallet-service to credit wallet
	payload := map[string]interface{}{
		"amount":      amount,
		"method":      "association_distribution",
		"description": description,
	}

	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/wallets/%s/deposit", s.walletServiceURL, walletID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("wallet service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("wallet credit failed with status %d", resp.StatusCode)
	}

	return nil
}

// === Helpers ===

func (s *AssociationService) hasAdminPermission(assoc *models.Association, userID string) bool {
	if assoc.CreatedBy == userID {
		return true
	}

	member, err := s.memberRepo.GetByAssociationAndUser(assoc.ID, userID)
	if err != nil {
		return false
	}

	return member.Role == models.MemberRolePresident || member.Role == models.MemberRoleTreasurer
}

func (s *AssociationService) getTotalRepaid(repaymentList []interface{}) float64 {
	var total float64
	for _, r := range repaymentList {
		if rep, ok := r.(map[string]interface{}); ok {
			if amt, ok := rep["amount"].(float64); ok {
				total += amt
			}
		}
	}
	return total
}
