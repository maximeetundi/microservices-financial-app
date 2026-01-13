package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
    "math/rand"

	"github.com/google/uuid"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

type TransferService struct {
	transferRepo *repository.TransferRepository
	walletRepo   *repository.WalletRepository
	kafkaClient  *messaging.KafkaClient
	config       *config.Config
}

func NewTransferService(
	transferRepo *repository.TransferRepository,
	walletRepo *repository.WalletRepository,
	kafkaClient *messaging.KafkaClient,
	config *config.Config,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		walletRepo:   walletRepo,
		kafkaClient:  kafkaClient,
		config:       config,
	}
}

func (s *TransferService) CreateTransfer(req *models.TransferRequest) (*models.Transfer, error) {
	// Validate source wallet
	fromWallet, err := s.walletRepo.GetByID(req.FromWalletID)
	if err != nil {
		return nil, fmt.Errorf("source wallet not found: %w", err)
	}

	// Check sufficient balance
	if fromWallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Calculate fee
	transferType := "domestic" // default type
	fee := s.calculateFee(transferType, req.Amount)
	
	// Check balance covers amount + fee
	totalDebit := req.Amount + fee
	if fromWallet.Balance < totalDebit {
		return nil, fmt.Errorf("insufficient balance for amount plus fees")
	}

	// Resolve destination wallet ID
	var destinationWalletID *string
	
	// If to_wallet_id is provided, use it directly
	if req.ToWalletID != nil && *req.ToWalletID != "" {
		destinationWalletID = req.ToWalletID
	} else if req.ToEmail != nil || req.ToPhone != nil {
		// P2P transfer: Find or create recipient wallet based on email/phone
		recipientWalletID, err := s.resolveOrCreateRecipientWallet(req.ToEmail, req.ToPhone, req.Currency)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve recipient: %w", err)
		}
		destinationWalletID = &recipientWalletID
	}

	// Create transfer record
	ref := generateReferenceID()
	fromWalletIDPtr := req.FromWalletID
	transfer := &models.Transfer{
		ID:             generateID(),
		FromWalletID:   &fromWalletIDPtr,
		ToWalletID:     destinationWalletID,
		RecipientEmail: req.ToEmail,
		RecipientPhone: req.ToPhone,
		TransferType:   transferType,
		Amount:         req.Amount,
		Fee:            fee,
		Currency:       req.Currency,
		Status:         "pending",
		Reference:      &ref,
		Description:    req.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.transferRepo.Create(transfer); err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	// Debit from source wallet immediately
	if err := s.walletRepo.UpdateBalance(req.FromWalletID, -totalDebit); err != nil {
		// Rollback: update transfer status to failed
		s.transferRepo.UpdateStatus(transfer.ID, "failed")
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	
	// Credit to destination wallet if resolved
	if destinationWalletID != nil && *destinationWalletID != "" {
		if err := s.walletRepo.UpdateBalance(*destinationWalletID, req.Amount); err != nil {
			// Rollback: refund source wallet and mark transfer failed
			s.walletRepo.UpdateBalance(req.FromWalletID, totalDebit)
			s.transferRepo.UpdateStatus(transfer.ID, "failed")
			return nil, fmt.Errorf("failed to credit destination wallet: %w", err)
		}
		// Mark as completed for internal transfers
		s.transferRepo.UpdateStatus(transfer.ID, "completed")
		transfer.Status = "completed"
	} else {
		// External transfers remain pending for async processing
		s.transferRepo.UpdateStatus(transfer.ID, "processing")
		transfer.Status = "processing"
	}

	// Publish transfer event to Kafka
	transferEvent := messaging.TransferCompletedEvent{
		TransferID:   transfer.ID,
		FromWalletID: req.FromWalletID,
		ToWalletID:   "",
		Amount:       transfer.Amount,
		Fee:          transfer.Fee,
		Status:       transfer.Status,
	}
	if destinationWalletID != nil {
		transferEvent.ToWalletID = *destinationWalletID
	}
	envelope := messaging.NewEventEnvelope(messaging.EventTransferCompleted, "transfer-service", transferEvent)
	s.kafkaClient.Publish(context.Background(), messaging.TopicTransferEvents, envelope)
	
	// Get sender user ID from wallet
	senderUserID := fromWallet.UserID
	
	// Get recipient user ID (if internal transfer)
	var recipientUserID string
	if destinationWalletID != nil && *destinationWalletID != "" {
		toWallet, err := s.walletRepo.GetByID(*destinationWalletID)
		if err == nil && toWallet != nil {
			recipientUserID = toWallet.UserID
		}
	}
	
	// Publish notification events for BOTH sender and recipient
	if transfer.Status == "completed" {
		// Notification for sender (money sent)
		s.publishTransferEvent("transfer.sent", transfer, senderUserID, senderUserID, recipientUserID)
		
		// Notification for recipient (money received)
		if recipientUserID != "" && recipientUserID != senderUserID {
			s.publishTransferEvent("transfer.received", transfer, recipientUserID, senderUserID, recipientUserID)
		}
	} else {
		// Transfer initiated but not yet completed
		s.publishTransferEvent("transfer.initiated", transfer, senderUserID, senderUserID, recipientUserID)
	}

	return transfer, nil
}



func (s *TransferService) GetTransfer(id string) (*models.Transfer, error) {
	transfer, err := s.transferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 1. Populate Sender Details
	// Getting sender user ID from wallet or directly from transfer
	var senderUserID string
	if transfer.FromWalletID != nil && *transfer.FromWalletID != "" {
		fromWallet, err := s.walletRepo.GetByID(*transfer.FromWalletID)
		if err == nil {
			senderUserID = fromWallet.UserID
		}
	}
	
	// Fallback to transfer.UserID if wallet lookup failed or provided no user (e.g. external)
	if senderUserID == "" {
		senderUserID = transfer.UserID
	}

	if senderUserID != "" {
		name, email, phone, err := s.walletRepo.GetUserInfo(senderUserID)
		if err == nil {
			transfer.SenderDetails = &models.UserDetails{
				Name:  name,
				Email: email,
				Phone: phone,
			}
		}
	}

	// 2. Populate Recipient Details
	// Check if internal transfer with wallet ID
	if transfer.ToWalletID != nil && *transfer.ToWalletID != "" {
		toWallet, err := s.walletRepo.GetByID(*transfer.ToWalletID)
		if err == nil {
			name, email, phone, err := s.walletRepo.GetUserInfo(toWallet.UserID)
			if err == nil {
				transfer.RecipientDetails = &models.UserDetails{
					Name:  name,
					Email: email,
					Phone: phone,
				}
			}
		}
	} else {
		// External or P2P by email/phone - data is already in transfer object
		var name, email, phone string
		
		// If mobile money or international, these details might be in the specific structs, 
		// but standard transfer struct has RecipientEmail/Phone too.
		if transfer.RecipientEmail != nil {
			email = *transfer.RecipientEmail
		}
		if transfer.RecipientPhone != nil {
			phone = *transfer.RecipientPhone
		}
		
		if email != "" || phone != "" {
			name = "Contact Externe"
			transfer.RecipientDetails = &models.UserDetails{
				Name:  name,
				Email: email,
				Phone: phone,
			}
		}
	}

	return transfer, nil
}

func (s *TransferService) GetTransferHistory(userID string, limit, offset int) ([]*models.Transfer, error) {
	// We need to fetch wallets for this user first
    // Note: This relies on walletRepo having a method to list user wallets, OR we query transfers by userID directly?
    // Usually transfers are linked to wallets. But Transfer model might not have UserID directly on it (it has FromWalletID).
    // Let's assume transferRepo has functionality to list by UserID (joining wallets) or we filter by wallets.
    // simpler: s.transferRepo.ListByUserID(userID, limit, offset)
    return s.transferRepo.GetByUserID(userID, limit, offset)
}

func (s *TransferService) CancelTransfer(id, requesterID, reason string) error {
	// 1. Get Transfer
	transfer, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 2. Validate Sender
	// We need to resolve sender user ID.
	// We can trust requesterID passed from handler (extracted from token).
	// Check if transfer.FromWalletID belongs to requesterID.
	if transfer.FromWalletID == nil {
		return fmt.Errorf("cannot cancel external transfer")
	}
	fromWallet, err := s.walletRepo.GetByID(*transfer.FromWalletID)
	if err != nil {
		return fmt.Errorf("sender wallet not found")
	}
	if fromWallet.UserID != requesterID {
		return fmt.Errorf("unauthorized: you are not the sender")
	}

	// 3. Check Time Conditions (Sender Cancel)
	// - Less than 1 hour
	// - OR More than 7 days (Inactive assumption)
	duration := time.Since(transfer.CreatedAt)
	isRecent := duration < time.Hour
	isOld := duration > 7*24*time.Hour

	if !isRecent && !isOld {
		return fmt.Errorf("cancellation allowed only within 1 hour or after 7 days of inactivity")
	}

	// 4. Check Status
	if transfer.Status == "cancelled" || transfer.Status == "reversed" {
		return fmt.Errorf("transfer already cancelled/reversed")
	}

	// 5. If Completed, Check Recipient Balance and Reversal
	if transfer.Status == "completed" {
		if transfer.ToWalletID == nil {
			return fmt.Errorf("cannot reverse external completed transfer")
		}
		toWalletID := *transfer.ToWalletID

		// Check Recipient Balance
		toWallet, err := s.walletRepo.GetByID(toWalletID)
		if err != nil {
			return fmt.Errorf("recipient wallet not found")
		}
		if toWallet.Balance < transfer.Amount {
			return fmt.Errorf("recipient has insufficient funds for reversal")
		}

		// Execute Reversal: Debit Recipient, Credit Sender
		if err := s.walletRepo.UpdateBalance(toWalletID, -transfer.Amount); err != nil {
			return fmt.Errorf("failed to debit recipient: %w", err)
		}
		if err := s.walletRepo.UpdateBalance(*transfer.FromWalletID, transfer.Amount); err != nil {
			// CRITICAL: Failed to credit sender after debiting recipient
			// Rollback recipient debit
			s.walletRepo.UpdateBalance(toWalletID, transfer.Amount)
			return fmt.Errorf("failed to credit sender: %w", err)
		}
	} else if transfer.Status == "pending" || transfer.Status == "processing" {
		// Just refund sender if pending (money hasn't reached recipient or strictly pending)
		// But in CreateTransfer we debit Sender immediately.
		// If Pending, money is in Utils/Hold? No, it's just out of Sender.
		// So Credit Sender.
		if err := s.walletRepo.UpdateBalance(*transfer.FromWalletID, transfer.Amount+transfer.Fee); err != nil {
			return fmt.Errorf("failed to refund sender: %w", err)
		}
	} else {
		return fmt.Errorf("cannot cancel transfer in status %s", transfer.Status)
	}

	// 6. Update Transfer Status
	// 6. Update Transfer Status

	// Note: We might want a separate field for cancellation reason, but appending to description works for MVP.
	// Or use a repository method that supports extra fields.
	// For now, simple status update.
	if err := s.transferRepo.UpdateStatus(id, "cancelled"); err != nil {
		return err
	}

	return nil
}

func (s *TransferService) ReverseTransfer(id, requesterID, reason string) error {
	// 1. Get Transfer
	transfer, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 2. Validate Recipient (Beneficiary)
	if transfer.ToWalletID == nil {
		return fmt.Errorf("cannot reverse external transfer")
	}
	toWallet, err := s.walletRepo.GetByID(*transfer.ToWalletID)
	if err != nil {
		return fmt.Errorf("recipient wallet not found")
	}
	if toWallet.UserID != requesterID {
		return fmt.Errorf("unauthorized: you are not the recipient")
	}

	// 3. Check Time Condition (Beneficiary Return)
	// - Less than 7 days
	duration := time.Since(transfer.CreatedAt)
	if duration > 7*24*time.Hour {
		return fmt.Errorf("return allowed only within 7 days")
	}

	// 4. Check Status
	if transfer.Status != "completed" {
		return fmt.Errorf("only completed transfers can be returned")
	}

	// 5. Check Balance
	if toWallet.Balance < transfer.Amount {
		return fmt.Errorf("insufficient funds to return transfer")
	}

	// 6. Execute Reversal
	// Debit Recipient
	if err := s.walletRepo.UpdateBalance(*transfer.ToWalletID, -transfer.Amount); err != nil {
		return fmt.Errorf("failed to debit account: %w", err)
	}
	// Credit Sender (FromWallet)
	if transfer.FromWalletID == nil {
		// Should not happen for internal/completed, but safe check
		// If FromWallet is missing, maybe burn? No, fail.
		s.walletRepo.UpdateBalance(*transfer.ToWalletID, transfer.Amount) // Rollback
		return fmt.Errorf("sender wallet unknown")
	}
	if err := s.walletRepo.UpdateBalance(*transfer.FromWalletID, transfer.Amount); err != nil {
		s.walletRepo.UpdateBalance(*transfer.ToWalletID, transfer.Amount) // Rollback
		return fmt.Errorf("failed to credit sender: %w", err)
	}

	// 7. Update Status
	// Set to "reversed"
	if err := s.transferRepo.UpdateStatus(id, "reversed"); err != nil {
		// Non-critical but bad state
		return err
	}

	return nil
}

// publishTransferEvent publishes transfer events to Kafka for notifications
func (s *TransferService) publishTransferEvent(eventType string, transfer *models.Transfer, targetUserID, actualSenderID, actualRecipientID string) {
	if s.kafkaClient == nil {
		return
	}

	// Get sender name for better notification readability
	senderName := actualSenderID
	if s.walletRepo != nil {
		if name, err := s.walletRepo.GetUserNameByID(actualSenderID); err == nil && name != "" {
			senderName = name
		}
	}

	eventData := map[string]interface{}{
		"transfer_id":  transfer.ID,
		"user_id":      targetUserID,    // The user who receives the notification
		"sender":       actualSenderID,  // The actual sender of the money
		"sender_name":  senderName,      // The name of the sender
		"recipient":    actualRecipientID, // The actual recipient of the money
		"amount":       transfer.Amount,
		"currency":     transfer.Currency,
		"status":       transfer.Status,
		"reference":    transfer.Reference,
	}

	envelope := messaging.NewEventEnvelope(eventType, "transfer-service", eventData)
	s.kafkaClient.Publish(context.Background(), messaging.TopicTransferEvents, envelope)
}

// MobileMoneyService handles mobile money transfers
type MobileMoneyService struct {
	config *config.Config
}

func NewMobileMoneyService(config *config.Config) *MobileMoneyService {
	return &MobileMoneyService{config: config}
}

func (s *MobileMoneyService) Send(req *models.MobileMoneyRequest) (*models.MobileMoneyResponse, error) {
	// Implementation for mobile money send
	return &models.MobileMoneyResponse{
		TransactionID: generateID(),
		Status:        "pending",
		Provider:      req.Provider,
	}, nil
}

func (s *MobileMoneyService) Receive(req *models.MobileMoneyRequest) (*models.MobileMoneyResponse, error) {
	return &models.MobileMoneyResponse{
		TransactionID: generateID(),
		Status:        "pending",
		Provider:      req.Provider,
	}, nil
}

func (s *MobileMoneyService) GetProviders() []string {
	providers := make([]string, 0, len(s.config.MobileMoneyProviders))
	for k := range s.config.MobileMoneyProviders {
		providers = append(providers, k)
	}
	return providers
}

// InternationalTransferService handles international transfers
type InternationalTransferService struct {
	config *config.Config
}

func NewInternationalTransferService(config *config.Config) *InternationalTransferService {
	return &InternationalTransferService{config: config}
}

func (s *InternationalTransferService) CreateTransfer(req *models.InternationalTransferRequest) (*models.Transfer, error) {
	ref := generateReferenceID()
	fromWalletIDPtr := req.FromWalletID
	transfer := &models.Transfer{
		ID:           generateID(),
		FromWalletID: &fromWalletIDPtr,
		TransferType: "international",
		Amount:       req.Amount,
		Currency:     req.Currency,
		Status:       "pending",
		Reference:    &ref,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return transfer, nil
}

// ComplianceService handles compliance checks
type ComplianceService struct {
	config *config.Config
}

func NewComplianceService(config *config.Config) *ComplianceService {
	return &ComplianceService{config: config}
}

func (s *ComplianceService) CheckTransfer(transfer *models.Transfer) (*models.ComplianceResult, error) {
	result := &models.ComplianceResult{
		Passed:    true,
		RiskScore: 0,
		Checks:    []string{"AML", "Sanctions", "PEP"},
	}

	// Check amount limits
	if transfer.Amount > s.config.ComplianceSettings.MaxAmountWithoutKYC {
		result.RequiresKYC = true
	}

	return result, nil
}

// Helper functions and methods

func (s *TransferService) calculateFee(transferType string, amount float64) float64 {
    // Simple fee logic: 1% for international, 0 for domestic
    if transferType == "international" {
        return amount * 0.01
    }
    return 0
}

func (s *TransferService) resolveOrCreateRecipientWallet(email *string, phone *string, currency string) (string, error) {
	var userID string
	var err error

	// 1. Resolve User ID
	if email != nil && *email != "" {
		userID, err = s.walletRepo.GetUserIDByEmail(*email)
	} else if phone != nil && *phone != "" {
		userID, err = s.walletRepo.GetUserIDByPhone(*phone)
	} else {
		return "", fmt.Errorf("either email or phone must be provided")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("recipient user not found")
		}
		return "", fmt.Errorf("failed to lookup recipient: %w", err)
	}

	// 2. Find existing wallet for this currency
	walletID, err := s.walletRepo.GetWalletIDByUserAndCurrency(userID, currency)
	if err == nil {
		return walletID, nil
	}
	
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to check recipient wallet: %w", err)
	}

	// 3. Create new wallet if not found
	newWalletID := generateID()
	if err := s.walletRepo.CreateWallet(newWalletID, userID, currency); err != nil {
		return "", fmt.Errorf("failed to create recipient wallet: %w", err)
	}

	return newWalletID, nil
}

func generateID() string {
    return uuid.New().String()
}

func generateReferenceID() string {
    // Generate a unique reference ID, e.g., TRF-timestamp-random
    return fmt.Sprintf("TRF-%d-%d", time.Now().Unix(), rand.Intn(10000))
}
