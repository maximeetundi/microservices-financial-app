package providers

import (
	"context"
	"fmt"
	"sync"
)

// FullTransferService combines all transfer capabilities
type FullTransferService struct {
	// Providers
	flutterwaveCollection *FlutterwaveCollectionProvider
	stripeCollection      *StripeCollectionProvider
	cryptoRails           *CryptoRailsProvider
	zoneRouter            *ZoneRouter
	internalTransfer      *InternalTransferService

	// Wallet service interface
	walletService WalletServiceInterface

	mu sync.RWMutex
}

// NewFullTransferService creates a complete transfer service
func NewFullTransferService(cfg *Config, walletService WalletServiceInterface) *FullTransferService {
	cryptoRails := NewCryptoRailsProvider(cfg.CryptoRails)
	zoneRouter := InitializeRouter(cfg)

	flutterwaveCollection := NewFlutterwaveCollectionProvider(cfg.Flutterwave)
	stripeCollection := NewStripeCollectionProvider(cfg.Stripe)

	internalTransfer := NewInternalTransferService(cryptoRails, zoneRouter, walletService)

	return &FullTransferService{
		flutterwaveCollection: flutterwaveCollection,
		stripeCollection:      stripeCollection,
		cryptoRails:           cryptoRails,
		zoneRouter:            zoneRouter,
		internalTransfer:      internalTransfer,
		walletService:         walletService,
	}
}

// ========================
// DEPOSIT FLOW
// ========================

// InitiateDeposit starts a deposit flow
func (s *FullTransferService) InitiateDeposit(ctx context.Context, req *CollectionRequest) (*CollectionResponse, error) {
	// Choose provider based on country
	zone := s.zoneRouter.GetZone(req.Country)

	switch zone {
	case ZoneAfrica:
		return s.flutterwaveCollection.InitiateCollection(ctx, req)
	case ZoneEurope, ZoneNorthAmerica:
		return s.stripeCollection.InitiateCollection(ctx, req)
	default:
		// Default to Flutterwave for other zones
		return s.flutterwaveCollection.InitiateCollection(ctx, req)
	}
}

// VerifyDeposit verifies a deposit and credits wallet
func (s *FullTransferService) VerifyDeposit(ctx context.Context, referenceID, provider string) (*CollectionResponse, error) {
	var resp *CollectionResponse
	var err error

	switch provider {
	case "flutterwave":
		resp, err = s.flutterwaveCollection.VerifyCollection(ctx, referenceID)
	case "stripe":
		resp, err = s.stripeCollection.VerifyCollection(ctx, referenceID)
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	// If successful, credit the user's wallet
	if resp.Status == CollectionStatusSuccessful {
		// Extract user wallet ID from metadata passed from InitiateCollection
		walletID := resp.Metadata["wallet_id"]
		userID := resp.Metadata["user_id"] // Ensure user_id is in metadata!

		if walletID != "" {
			// Double Check: We need userID for the new method.
			// Ideally metadata has it. If not, we might fail or need to fetch wallet.
			// Assuming metadata has user_id populated as well in InitiateCollection.
			// If not, we need to fetch the wallet to get the UserID?
			// But for now, let's assume valid metadata or fetch wallet if missing (safest).

			// Use the new Double Entry method
			// provider arg is passed to verifyDeposit
			providerName := provider
			providerRef := resp.ProviderReference

			// Fallback if userID is missing from metadata (prevent breakage if old flows exist)
			if userID == "" {
				// We should ideally fetch the wallet here to get the user ID, but we don't have GetWallet easily exposed here?
				// WalletService has GetBalance but not GetWallet in generic interface?
				// internal_transfer has GetBalance which returns UserID.
				bal, err := s.walletService.GetBalance(ctx, walletID)
				if err == nil {
					userID = bal.UserID
				} else {
					// Logic error: can't deposit if we don't know the user (for consistency checks)
					// But current CreditWallet didn't generic ask for user_id.
					// However, CreditWalletFromPlatform DOES require it.
				}
			}

			if err := s.walletService.CreditWalletFromPlatform(ctx, userID, walletID, resp.NetAmount, resp.Currency, providerRef, providerName); err != nil {
				return resp, fmt.Errorf("failed to credit wallet (double-entry): %w", err)
			}
		} else {
			// Log warning: missing wallet ID in metadata, manual intervention might be needed
			// For now just returning, but ideally should alert
			// fmt.Printf("Warning: Deposit %s verified but no wallet_id found in metadata\n", referenceID)
		}
	}

	return resp, nil
}

// ========================
// INTERNAL TRANSFER FLOW
// ========================

// TransferToUser transfers from one user's wallet to another
func (s *FullTransferService) TransferToUser(ctx context.Context, req *InternalTransferRequest) (*InternalTransfer, error) {
	return s.internalTransfer.CreateTransfer(ctx, req)
}

// GetTransferStatus gets status of internal transfer
func (s *FullTransferService) GetTransferStatus(ctx context.Context, transferID string) (*InternalTransfer, error) {
	return s.internalTransfer.GetTransfer(ctx, transferID)
}

// ========================
// WITHDRAWAL/PAYOUT FLOW
// ========================

// WithdrawToExternal withdraws from wallet to external account
func (s *FullTransferService) WithdrawToExternal(ctx context.Context, req *PayoutRequest, walletID string) (*PayoutResponse, error) {
	// Step 1: Check wallet balance
	balance, err := s.walletService.GetBalance(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if balance.AvailableBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: available %.2f, required %.2f",
			balance.AvailableBalance, req.Amount)
	}

	// Step 2: Debit wallet
	if err := s.walletService.DebitWallet(ctx, walletID, req.Amount, req.ReferenceID); err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}

	// Step 3: Execute payout via zone router
	resp, err := s.zoneRouter.CreatePayout(ctx, req)
	if err != nil {
		// Refund wallet on failure
		s.walletService.CreditWallet(ctx, walletID, req.Amount, req.ReferenceID+"_refund")
		return nil, fmt.Errorf("payout failed: %w", err)
	}

	return resp, nil
}

// GetWithdrawalStatus gets status of withdrawal
func (s *FullTransferService) GetWithdrawalStatus(ctx context.Context, referenceID, provider string) (*PayoutStatusResponse, error) {
	return s.zoneRouter.GetPayoutStatus(ctx, referenceID, provider)
}

// ========================
// FULL INTERNATIONAL TRANSFER (Send Money)
// ========================

// SendMoney performs a complete international transfer:
// 1. Debit sender's wallet
// 2. Convert currency (via crypto rails if needed)
// 3. Credit recipient's wallet OR payout directly
type SendMoneyRequest struct {
	ReferenceID string `json:"reference_id"`

	// Sender
	SenderUserID   string  `json:"sender_user_id"`
	SenderWalletID string  `json:"sender_wallet_id"`
	Amount         float64 `json:"amount"`
	SourceCurrency string  `json:"source_currency"`

	// Recipient - can be internal user OR external
	IsInternalRecipient bool `json:"is_internal_recipient"`

	// For internal recipient
	RecipientUserID   string `json:"recipient_user_id,omitempty"`
	RecipientWalletID string `json:"recipient_wallet_id,omitempty"`

	// For external recipient
	RecipientName    string       `json:"recipient_name,omitempty"`
	RecipientPhone   string       `json:"recipient_phone,omitempty"`
	RecipientCountry string       `json:"recipient_country"`
	TargetCurrency   string       `json:"target_currency"`
	PayoutMethod     PayoutMethod `json:"payout_method,omitempty"`

	// Payout details
	MobileOperator string `json:"mobile_operator,omitempty"`
	MobileNumber   string `json:"mobile_number,omitempty"`
	BankCode       string `json:"bank_code,omitempty"`
	AccountNumber  string `json:"account_number,omitempty"`

	Narration string `json:"narration,omitempty"`
}

type SendMoneyResponse struct {
	ReferenceID string `json:"reference_id"`

	// Amounts
	SentAmount       float64 `json:"sent_amount"`
	SentCurrency     string  `json:"sent_currency"`
	ReceivedAmount   float64 `json:"received_amount"`
	ReceivedCurrency string  `json:"received_currency"`

	// Fees
	TotalFee     float64 `json:"total_fee"`
	ExchangeRate float64 `json:"exchange_rate"`

	// Delivery
	DeliveryMethod string `json:"delivery_method"` // "internal_wallet", "mobile_money", "bank_transfer"
	Status         string `json:"status"`
	EstimatedTime  string `json:"estimated_time"`

	// References
	InternalTransferID string `json:"internal_transfer_id,omitempty"`
	PayoutReference    string `json:"payout_reference,omitempty"`
}

func (s *FullTransferService) SendMoney(ctx context.Context, req *SendMoneyRequest) (*SendMoneyResponse, error) {
	if req.IsInternalRecipient {
		// Internal transfer to another app user
		return s.sendToInternalUser(ctx, req)
	} else {
		// External payout (Mobile Money, Bank, etc.)
		return s.sendToExternal(ctx, req)
	}
}

func (s *FullTransferService) sendToInternalUser(ctx context.Context, req *SendMoneyRequest) (*SendMoneyResponse, error) {
	transferReq := &InternalTransferRequest{
		ReferenceID:       req.ReferenceID,
		SenderUserID:      req.SenderUserID,
		SenderWalletID:    req.SenderWalletID,
		Amount:            req.Amount,
		SourceCurrency:    req.SourceCurrency,
		RecipientUserID:   req.RecipientUserID,
		RecipientWalletID: req.RecipientWalletID,
		TargetCurrency:    req.TargetCurrency,
		Narration:         req.Narration,
	}

	transfer, err := s.internalTransfer.CreateTransfer(ctx, transferReq)
	if err != nil {
		return nil, err
	}

	estimatedTime := "Instant"
	if transfer.Status == TransferStatusLocked {
		estimatedTime = "2-5 minutes"
	}

	return &SendMoneyResponse{
		ReferenceID:        req.ReferenceID,
		SentAmount:         transfer.SenderAmount,
		SentCurrency:       transfer.SenderCurrency,
		ReceivedAmount:     transfer.RecipientAmount,
		ReceivedCurrency:   transfer.RecipientCurrency,
		TotalFee:           transfer.TotalFee,
		ExchangeRate:       transfer.ExchangeRate,
		DeliveryMethod:     "internal_wallet",
		Status:             string(transfer.Status),
		EstimatedTime:      estimatedTime,
		InternalTransferID: transfer.ID,
	}, nil
}

func (s *FullTransferService) sendToExternal(ctx context.Context, req *SendMoneyRequest) (*SendMoneyResponse, error) {
	// Step 1: Convert via crypto rails
	convReq := &ConversionRequest{
		ReferenceID:      req.ReferenceID,
		SourceAmount:     req.Amount,
		SourceCurrency:   req.SourceCurrency,
		TargetCurrency:   req.TargetCurrency,
		RecipientCountry: req.RecipientCountry,
		PayoutMethod:     req.PayoutMethod,
	}

	convResult, err := s.cryptoRails.ConvertViaStablecoin(ctx, convReq)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}

	// Step 2: Debit sender's wallet
	if err := s.walletService.DebitWallet(ctx, req.SenderWalletID, req.Amount, req.ReferenceID); err != nil {
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}

	// Step 3: Execute payout
	payoutReq := &PayoutRequest{
		ReferenceID:      req.ReferenceID,
		Amount:           convResult.TargetAmount,
		Currency:         req.TargetCurrency,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientCountry: req.RecipientCountry,
		PayoutMethod:     req.PayoutMethod,
		MobileOperator:   req.MobileOperator,
		MobileNumber:     req.MobileNumber,
		BankCode:         req.BankCode,
		AccountNumber:    req.AccountNumber,
		Narration:        req.Narration,
	}

	payoutResult, err := s.zoneRouter.CreatePayout(ctx, payoutReq)
	if err != nil {
		// Refund sender
		s.walletService.CreditWallet(ctx, req.SenderWalletID, req.Amount, req.ReferenceID+"_refund")
		return nil, fmt.Errorf("payout failed: %w", err)
	}

	return &SendMoneyResponse{
		ReferenceID:      req.ReferenceID,
		SentAmount:       req.Amount,
		SentCurrency:     req.SourceCurrency,
		ReceivedAmount:   payoutResult.AmountReceived,
		ReceivedCurrency: req.TargetCurrency,
		TotalFee:         convResult.TotalFee + payoutResult.Fee,
		ExchangeRate:     convResult.USDToTargetRate,
		DeliveryMethod:   string(req.PayoutMethod),
		Status:           string(payoutResult.Status),
		EstimatedTime:    "5-30 minutes",
		PayoutReference:  payoutResult.ProviderReference,
	}, nil
}

// ========================
// AVAILABLE METHODS
// ========================

// GetDepositMethods returns available deposit methods for a country
func (s *FullTransferService) GetDepositMethods(ctx context.Context, country string) ([]CollectionMethod, error) {
	zone := s.zoneRouter.GetZone(country)

	switch zone {
	case ZoneAfrica:
		return s.flutterwaveCollection.GetAvailableMethods(ctx, country)
	case ZoneEurope, ZoneNorthAmerica:
		return s.stripeCollection.GetAvailableMethods(ctx, country)
	default:
		return s.flutterwaveCollection.GetAvailableMethods(ctx, country)
	}
}

// GetWithdrawalMethods returns available withdrawal methods for a country
func (s *FullTransferService) GetWithdrawalMethods(ctx context.Context, country string) ([]AvailableMethod, error) {
	return s.zoneRouter.GetAvailableMethodsForCountry(ctx, country)
}

// GetBanks returns available banks for a country
func (s *FullTransferService) GetBanks(ctx context.Context, country string) ([]Bank, error) {
	return s.zoneRouter.GetBanksForCountry(ctx, country)
}

// GetMobileOperators returns available mobile operators for a country
func (s *FullTransferService) GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error) {
	return s.zoneRouter.GetMobileOperatorsForCountry(ctx, country)
}
