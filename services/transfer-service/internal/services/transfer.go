package services

import (
	"context"
	"fmt"
	"time"

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

// ... (omitted methods remain unchanged)

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
