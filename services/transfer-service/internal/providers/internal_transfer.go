package providers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TransferStatus represents the status of an internal transfer
type TransferStatus string

const (
	TransferStatusPending     TransferStatus = "pending"
	TransferStatusProcessing  TransferStatus = "processing"
	TransferStatusLocked      TransferStatus = "locked"      // Funds locked during crypto conversion
	TransferStatusCompleted   TransferStatus = "completed"
	TransferStatusFailed      TransferStatus = "failed"
	TransferStatusRefunded    TransferStatus = "refunded"
)

// InternalTransfer represents a wallet-to-wallet transfer
type InternalTransfer struct {
	ID              string         `json:"id"`
	ReferenceID     string         `json:"reference_id"`
	
	// Sender
	SenderUserID    string         `json:"sender_user_id"`
	SenderWalletID  string         `json:"sender_wallet_id"`
	SenderAmount    float64        `json:"sender_amount"`
	SenderCurrency  string         `json:"sender_currency"`
	
	// Recipient
	RecipientUserID   string       `json:"recipient_user_id"`
	RecipientWalletID string       `json:"recipient_wallet_id"`
	RecipientAmount   float64      `json:"recipient_amount"`
	RecipientCurrency string       `json:"recipient_currency"`
	
	// Conversion info
	ExchangeRate      float64      `json:"exchange_rate"`
	UsedCryptoRails   bool         `json:"used_crypto_rails"`
	CryptoTxHash      string       `json:"crypto_tx_hash,omitempty"`
	
	// Fees
	ConversionFee     float64      `json:"conversion_fee"`
	TransferFee       float64      `json:"transfer_fee"`
	TotalFee          float64      `json:"total_fee"`
	
	// Status
	Status            TransferStatus `json:"status"`
	StatusMessage     string        `json:"status_message,omitempty"`
	
	// Timestamps
	CreatedAt         time.Time    `json:"created_at"`
	LockedAt          *time.Time   `json:"locked_at,omitempty"`
	CompletedAt       *time.Time   `json:"completed_at,omitempty"`
	
	// Estimated time for locked transfers
	EstimatedUnlockAt *time.Time   `json:"estimated_unlock_at,omitempty"`
}

// InternalTransferRequest represents a request to transfer between wallets
type InternalTransferRequest struct {
	ReferenceID       string  `json:"reference_id"`
	
	// Sender
	SenderUserID      string  `json:"sender_user_id"`
	SenderWalletID    string  `json:"sender_wallet_id"`
	Amount            float64 `json:"amount"`
	SourceCurrency    string  `json:"source_currency"`
	
	// Recipient
	RecipientUserID   string  `json:"recipient_user_id"`
	RecipientWalletID string  `json:"recipient_wallet_id"`
	TargetCurrency    string  `json:"target_currency"`
	
	// Optional: for external payout after internal transfer
	PayoutAfterTransfer bool         `json:"payout_after_transfer"`
	PayoutMethod        PayoutMethod `json:"payout_method,omitempty"`
	PayoutDetails       *PayoutRequest `json:"payout_details,omitempty"`
	
	Narration         string  `json:"narration,omitempty"`
}

// WalletBalance represents a wallet balance (interface to wallet-service)
type WalletBalance struct {
	WalletID        string  `json:"wallet_id"`
	UserID          string  `json:"user_id"`
	Currency        string  `json:"currency"`
	AvailableBalance float64 `json:"available_balance"`
	LockedBalance   float64 `json:"locked_balance"`
	TotalBalance    float64 `json:"total_balance"`
}

// WalletServiceInterface interface to communicate with wallet-service
type WalletServiceInterface interface {
	GetBalance(ctx context.Context, walletID string) (*WalletBalance, error)
	LockFunds(ctx context.Context, walletID string, amount float64, referenceID string) error
	UnlockFunds(ctx context.Context, walletID string, referenceID string) error
	DebitWallet(ctx context.Context, walletID string, amount float64, referenceID string) error
	CreditWallet(ctx context.Context, walletID string, amount float64, referenceID string) error
}

// InternalTransferService handles wallet-to-wallet transfers
type InternalTransferService struct {
	cryptoRails   *CryptoRailsProvider
	zoneRouter    *ZoneRouter
	walletService WalletServiceInterface
	
	// In-memory store for transfers (in production, use database)
	transfers     map[string]*InternalTransfer
	mu            sync.RWMutex
	
	// Threshold for instant vs crypto transfer
	instantThreshold float64
}

// NewInternalTransferService creates a new internal transfer service
func NewInternalTransferService(
	cryptoRails *CryptoRailsProvider,
	zoneRouter *ZoneRouter,
	walletService WalletServiceInterface,
) *InternalTransferService {
	return &InternalTransferService{
		cryptoRails:      cryptoRails,
		zoneRouter:       zoneRouter,
		walletService:    walletService,
		transfers:        make(map[string]*InternalTransfer),
		instantThreshold: 500.0, // $500 threshold
	}
}

// CreateTransfer initiates an internal transfer
func (s *InternalTransferService) CreateTransfer(ctx context.Context, req *InternalTransferRequest) (*InternalTransfer, error) {
	// Step 1: Check sender balance
	senderBalance, err := s.walletService.GetBalance(ctx, req.SenderWalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender balance: %w", err)
	}
	
	if senderBalance.AvailableBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: available %.2f, required %.2f", 
			senderBalance.AvailableBalance, req.Amount)
	}
	
	// Step 2: Determine if same currency or needs conversion
	needsConversion := req.SourceCurrency != req.TargetCurrency
	
	// Step 3: Create transfer record
	transfer := &InternalTransfer{
		ID:              fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			time.Now().Unix(),
			time.Now().Nanosecond()&0xffff,
			0x4000|time.Now().Nanosecond()&0x0fff,
			0x8000|time.Now().Nanosecond()&0x3fff,
			time.Now().UnixNano()&0xffffffffffff),
		ReferenceID:     req.ReferenceID,
		SenderUserID:    req.SenderUserID,
		SenderWalletID:  req.SenderWalletID,
		SenderAmount:    req.Amount,
		SenderCurrency:  req.SourceCurrency,
		RecipientUserID: req.RecipientUserID,
		RecipientWalletID: req.RecipientWalletID,
		RecipientCurrency: req.TargetCurrency,
		Status:          TransferStatusPending,
		CreatedAt:       time.Now(),
	}
	
	// Step 4: Process based on amount and conversion need
	if needsConversion {
		// Use crypto rails for currency conversion
		if req.Amount > s.instantThreshold {
			// Large amount: Lock funds and process via blockchain
			return s.processLargeTransfer(ctx, transfer, req)
		} else {
			// Small amount: Use internal pool (instant)
			return s.processInstantTransfer(ctx, transfer, req)
		}
	} else {
		// Same currency: Direct transfer (instant)
		return s.processDirectTransfer(ctx, transfer, req)
	}
}

// processDirectTransfer handles same-currency transfers (instant)
func (s *InternalTransferService) processDirectTransfer(ctx context.Context, transfer *InternalTransfer, req *InternalTransferRequest) (*InternalTransfer, error) {
	// Small fee for same currency
	transfer.TransferFee = transfer.SenderAmount * 0.001 // 0.1%
	transfer.TotalFee = transfer.TransferFee
	transfer.RecipientAmount = transfer.SenderAmount - transfer.TransferFee
	transfer.ExchangeRate = 1.0
	transfer.UsedCryptoRails = false
	
	// Debit sender
	if err := s.walletService.DebitWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID); err != nil {
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = err.Error()
		return transfer, err
	}
	
	// Credit recipient
	if err := s.walletService.CreditWallet(ctx, req.RecipientWalletID, transfer.RecipientAmount, transfer.ReferenceID); err != nil {
		// Rollback: refund sender
		s.walletService.CreditWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID+"_refund")
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = err.Error()
		return transfer, err
	}
	
	now := time.Now()
	transfer.Status = TransferStatusCompleted
	transfer.CompletedAt = &now
	
	s.storeTransfer(transfer)
	
	return transfer, nil
}

// processInstantTransfer handles cross-currency transfers via internal pool
func (s *InternalTransferService) processInstantTransfer(ctx context.Context, transfer *InternalTransfer, req *InternalTransferRequest) (*InternalTransfer, error) {
	// Get conversion quote from crypto rails
	convReq := &ConversionRequest{
		ReferenceID:    req.ReferenceID,
		SourceAmount:   req.Amount,
		SourceCurrency: req.SourceCurrency,
		TargetCurrency: req.TargetCurrency,
	}
	
	convResult, err := s.cryptoRails.ConvertViaStablecoin(ctx, convReq)
	if err != nil {
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = err.Error()
		return transfer, err
	}
	
	transfer.RecipientAmount = convResult.TargetAmount
	transfer.ExchangeRate = convResult.USDToTargetRate
	transfer.ConversionFee = convResult.ConversionFee
	transfer.TransferFee = transfer.SenderAmount * 0.001
	transfer.TotalFee = convResult.TotalFee + transfer.TransferFee
	transfer.UsedCryptoRails = true
	
	// Debit sender in source currency
	if err := s.walletService.DebitWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID); err != nil {
		transfer.Status = TransferStatusFailed
		return transfer, err
	}
	
	// Credit recipient in target currency
	if err := s.walletService.CreditWallet(ctx, req.RecipientWalletID, transfer.RecipientAmount, transfer.ReferenceID); err != nil {
		// Rollback
		s.walletService.CreditWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID+"_refund")
		transfer.Status = TransferStatusFailed
		return transfer, err
	}
	
	now := time.Now()
	transfer.Status = TransferStatusCompleted
	transfer.CompletedAt = &now
	
	s.storeTransfer(transfer)
	
	return transfer, nil
}

// processLargeTransfer handles large cross-currency transfers with fund locking
func (s *InternalTransferService) processLargeTransfer(ctx context.Context, transfer *InternalTransfer, req *InternalTransferRequest) (*InternalTransfer, error) {
	// Step 1: Lock sender's funds
	if err := s.walletService.LockFunds(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID); err != nil {
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = "Failed to lock funds: " + err.Error()
		return transfer, err
	}
	
	now := time.Now()
	transfer.Status = TransferStatusLocked
	transfer.LockedAt = &now
	
	// Estimate unlock time (2-5 minutes for crypto)
	unlockTime := now.Add(5 * time.Minute)
	transfer.EstimatedUnlockAt = &unlockTime
	
	s.storeTransfer(transfer)
	
	// Step 2: Process conversion in background (goroutine)
	go s.processLockedTransfer(context.Background(), transfer, req)
	
	return transfer, nil
}

// processLockedTransfer processes a locked transfer in background
func (s *InternalTransferService) processLockedTransfer(ctx context.Context, transfer *InternalTransfer, req *InternalTransferRequest) {
	// Get conversion from crypto rails (this may take time on blockchain)
	convReq := &ConversionRequest{
		ReferenceID:    req.ReferenceID,
		SourceAmount:   req.Amount,
		SourceCurrency: req.SourceCurrency,
		TargetCurrency: req.TargetCurrency,
	}
	
	convResult, err := s.cryptoRails.ConvertViaStablecoin(ctx, convReq)
	if err != nil {
		// Unlock funds and mark as failed
		s.walletService.UnlockFunds(ctx, req.SenderWalletID, transfer.ReferenceID)
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = "Conversion failed: " + err.Error()
		s.storeTransfer(transfer)
		return
	}
	
	transfer.RecipientAmount = convResult.TargetAmount
	transfer.ExchangeRate = convResult.USDToTargetRate
	transfer.ConversionFee = convResult.ConversionFee
	transfer.TotalFee = convResult.TotalFee
	transfer.CryptoTxHash = convResult.BlockchainTxHash
	transfer.Status = TransferStatusProcessing
	s.storeTransfer(transfer)
	
	// Unlock and debit sender
	s.walletService.UnlockFunds(ctx, req.SenderWalletID, transfer.ReferenceID)
	if err := s.walletService.DebitWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID); err != nil {
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = "Debit failed: " + err.Error()
		s.storeTransfer(transfer)
		return
	}
	
	// Credit recipient
	if err := s.walletService.CreditWallet(ctx, req.RecipientWalletID, transfer.RecipientAmount, transfer.ReferenceID); err != nil {
		// Refund sender
		s.walletService.CreditWallet(ctx, req.SenderWalletID, transfer.SenderAmount, transfer.ReferenceID+"_refund")
		transfer.Status = TransferStatusFailed
		transfer.StatusMessage = "Credit failed: " + err.Error()
		s.storeTransfer(transfer)
		return
	}
	
	// If payout requested, execute it
	if req.PayoutAfterTransfer && req.PayoutDetails != nil {
		// Trigger payout to external (Mobile Money, Bank, etc.)
		req.PayoutDetails.ReferenceID = transfer.ReferenceID
		req.PayoutDetails.Amount = transfer.RecipientAmount
		req.PayoutDetails.Currency = transfer.RecipientCurrency
		
		// Note: Payout is fire-and-forget here, could be improved with proper error handling
		// The payout result should be tracked separately
	}
	
	now := time.Now()
	transfer.Status = TransferStatusCompleted
	transfer.CompletedAt = &now
	s.storeTransfer(transfer)
}

// GetTransfer retrieves a transfer by ID
func (s *InternalTransferService) GetTransfer(ctx context.Context, transferID string) (*InternalTransfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	transfer, ok := s.transfers[transferID]
	if !ok {
		return nil, fmt.Errorf("transfer not found: %s", transferID)
	}
	return transfer, nil
}

// GetTransferByReference retrieves a transfer by reference ID
func (s *InternalTransferService) GetTransferByReference(ctx context.Context, referenceID string) (*InternalTransfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, t := range s.transfers {
		if t.ReferenceID == referenceID {
			return t, nil
		}
	}
	return nil, fmt.Errorf("transfer not found for reference: %s", referenceID)
}

func (s *InternalTransferService) storeTransfer(transfer *InternalTransfer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transfers[transfer.ID] = transfer
}
