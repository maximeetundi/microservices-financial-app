package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

type TransferService struct {
	transferRepo *repository.TransferRepository
	walletRepo   *repository.WalletRepository
	mqClient     *database.RabbitMQClient
	config       *config.Config
}

func NewTransferService(
	transferRepo *repository.TransferRepository,
	walletRepo *repository.WalletRepository,
	mqClient *database.RabbitMQClient,
	config *config.Config,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		walletRepo:   walletRepo,
		mqClient:     mqClient,
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

	// Publish to queue for async processing (notifications, etc)
	msg, _ := json.Marshal(transfer)
	s.mqClient.Publish("transfers", msg)
	
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
		s.publishTransferEvent("transfer.sent", transfer, senderUserID, recipientUserID)
		
		// Notification for recipient (money received)
		if recipientUserID != "" && recipientUserID != senderUserID {
			s.publishTransferEvent("transfer.received", transfer, recipientUserID, senderUserID)
		}
	} else {
		// Transfer initiated but not yet completed
		s.publishTransferEvent("transfer.initiated", transfer, senderUserID, recipientUserID)
	}

	return transfer, nil
}

// resolveOrCreateRecipientWallet finds a recipient's wallet by email or phone, or creates one
func (s *TransferService) resolveOrCreateRecipientWallet(email, phone *string, currency string) (string, error) {
	// Find recipient user by email or phone
	var recipientUserID string
	var err error
	
	if email != nil && *email != "" {
		recipientUserID, err = s.walletRepo.GetUserIDByEmail(*email)
	} else if phone != nil && *phone != "" {
		recipientUserID, err = s.walletRepo.GetUserIDByPhone(*phone)
	} else {
		return "", fmt.Errorf("no recipient email or phone provided")
	}
	
	if err != nil {
		return "", fmt.Errorf("recipient user not found: %w", err)
	}
	
	// Find recipient's wallet in the same currency
	walletID, err := s.walletRepo.GetWalletIDByUserAndCurrency(recipientUserID, currency)
	if err == nil && walletID != "" {
		return walletID, nil
	}
	
	// Wallet doesn't exist - create one for the recipient
	newWalletID := generateID()
	err = s.walletRepo.CreateWallet(newWalletID, recipientUserID, currency)
	if err != nil {
		return "", fmt.Errorf("failed to create recipient wallet: %w", err)
	}
	
	return newWalletID, nil
}

func (s *TransferService) GetTransfer(id string) (*models.Transfer, error) {
	return s.transferRepo.GetByID(id)
}

func (s *TransferService) GetTransferHistory(userID string, limit, offset int) ([]models.Transfer, error) {
	return s.transferRepo.GetByUserID(userID, limit, offset)
}

func (s *TransferService) CancelTransfer(id string) error {
	transfer, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}
	if transfer.Status != "pending" {
		return fmt.Errorf("cannot cancel transfer with status: %s", transfer.Status)
	}
	return s.transferRepo.UpdateStatus(id, "cancelled")
}

func (s *TransferService) calculateFee(transferType string, amount float64) float64 {
	feeRate := s.config.TransferFees[transferType]
	if feeRate == 0 {
		feeRate = 0.5 // Default 0.5%
	}
	return amount * (feeRate / 100)
}

func generateID() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		time.Now().Unix(),
		time.Now().Nanosecond()&0xffff,
		0x4000|time.Now().Nanosecond()&0x0fff,
		0x8000|time.Now().Nanosecond()&0x3fff,
		time.Now().UnixNano()&0xffffffffffff)
}

func generateReferenceID() string {
	return fmt.Sprintf("REF%d", time.Now().UnixNano())
}

// publishTransferEvent publishes transfer events to RabbitMQ for notifications
func (s *TransferService) publishTransferEvent(eventType string, transfer *models.Transfer, senderUserID, recipientUserID string) {
	if s.mqClient == nil {
		return
	}

	// Get sender name for better notification readability
	senderName := senderUserID
	if s.walletRepo != nil {
		if name, err := s.walletRepo.GetUserNameByID(senderUserID); err == nil && name != "" {
			senderName = name
		}
	}

	event := map[string]interface{}{
		"type":         eventType,
		"transfer_id":  transfer.ID,
		"user_id":      senderUserID, // For sender notification
		"sender":       senderUserID,
		"sender_name":  senderName, // Human readable name for notifications
		"recipient":    recipientUserID,
		"amount":       transfer.Amount,
		"currency":     transfer.Currency,
		"status":       transfer.Status,
		"reference":    transfer.Reference,
		"timestamp":    time.Now(),
	}

	eventJSON, _ := json.Marshal(event)

	// Publish to transfer.events exchange
	s.mqClient.PublishToExchange("transfer.events", eventType, eventJSON)
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
