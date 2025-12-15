package providers

import (
	"context"
	"fmt"
	"time"
)

// TransferOrchestrator orchestrates the complete transfer flow:
// 1. Source Fiat → Stablecoin (via CryptoRails)
// 2. Stablecoin → Target Fiat (via Zone Router)
type TransferOrchestrator struct {
	cryptoRails *CryptoRailsProvider
	zoneRouter  *ZoneRouter
}

// NewTransferOrchestrator creates a new orchestrator
func NewTransferOrchestrator(cryptoRails *CryptoRailsProvider, zoneRouter *ZoneRouter) *TransferOrchestrator {
	return &TransferOrchestrator{
		cryptoRails: cryptoRails,
		zoneRouter:  zoneRouter,
	}
}

// FullTransferRequest represents a complete international transfer
type FullTransferRequest struct {
	// Reference
	ReferenceID string `json:"reference_id"`
	
	// Source (sender)
	SenderUserID     string  `json:"sender_user_id"`
	SourceAmount     float64 `json:"source_amount"`
	SourceCurrency   string  `json:"source_currency"` // EUR, USD, GBP
	
	// Destination (recipient)
	RecipientName    string       `json:"recipient_name"`
	RecipientPhone   string       `json:"recipient_phone"`
	RecipientCountry string       `json:"recipient_country"` // CI, SN, KE, etc.
	TargetCurrency   string       `json:"target_currency"`   // XOF, KES, etc.
	PayoutMethod     PayoutMethod `json:"payout_method"`
	
	// Payout details
	MobileOperator string `json:"mobile_operator,omitempty"`
	MobileNumber   string `json:"mobile_number,omitempty"`
	BankCode       string `json:"bank_code,omitempty"`
	AccountNumber  string `json:"account_number,omitempty"`
	
	// Extra
	Narration string `json:"narration,omitempty"`
}

// FullTransferResponse represents the complete transfer result
type FullTransferResponse struct {
	ReferenceID string `json:"reference_id"`
	
	// Amounts
	SourceAmount   float64 `json:"source_amount"`
	SourceCurrency string  `json:"source_currency"`
	TargetAmount   float64 `json:"target_amount"`
	TargetCurrency string  `json:"target_currency"`
	
	// Fees breakdown
	CryptoConversionFee float64 `json:"crypto_conversion_fee"`
	NetworkFee          float64 `json:"network_fee"`
	PayoutFee           float64 `json:"payout_fee"`
	TotalFee            float64 `json:"total_fee"`
	
	// Exchange rates
	ExchangeRate float64 `json:"exchange_rate"`
	
	// Processing details
	UsedCryptoRails    bool   `json:"used_crypto_rails"`
	UsedInternalPool   bool   `json:"used_internal_pool"`
	PayoutProvider     string `json:"payout_provider"`
	PayoutMethod       string `json:"payout_method"`
	
	// Status
	Status           string    `json:"status"`
	EstimatedArrival time.Time `json:"estimated_arrival"`
	
	// References
	CryptoTxHash      string `json:"crypto_tx_hash,omitempty"`
	PayoutReference   string `json:"payout_reference,omitempty"`
}

// GetQuote gets a complete quote for the transfer
func (o *TransferOrchestrator) GetQuote(ctx context.Context, req *FullTransferRequest) (*FullTransferResponse, error) {
	// Step 1: Get crypto rails quote
	cryptoReq := &ConversionRequest{
		ReferenceID:      req.ReferenceID,
		SourceAmount:     req.SourceAmount,
		SourceCurrency:   req.SourceCurrency,
		TargetCurrency:   req.TargetCurrency,
		RecipientCountry: req.RecipientCountry,
		PayoutMethod:     req.PayoutMethod,
	}
	
	cryptoQuote, err := o.cryptoRails.GetQuote(ctx, cryptoReq)
	if err != nil {
		return nil, fmt.Errorf("crypto rails quote failed: %w", err)
	}
	
	// Step 2: Get payout quote from zone router
	payoutReq := &PayoutRequest{
		ReferenceID:      req.ReferenceID,
		Amount:           cryptoQuote.TargetAmount,
		Currency:         req.TargetCurrency,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientCountry: req.RecipientCountry,
		PayoutMethod:     req.PayoutMethod,
		MobileOperator:   req.MobileOperator,
		MobileNumber:     req.MobileNumber,
		BankCode:         req.BankCode,
		AccountNumber:    req.AccountNumber,
	}
	
	provider, err := o.zoneRouter.GetProvider(req.RecipientCountry, req.PayoutMethod)
	if err != nil {
		return nil, fmt.Errorf("no payout provider available: %w", err)
	}
	
	payoutQuote, err := provider.GetQuote(ctx, payoutReq)
	if err != nil {
		return nil, fmt.Errorf("payout quote failed: %w", err)
	}
	
	// Step 3: Combine into full response
	totalFee := cryptoQuote.TotalFee + payoutQuote.Fee
	finalAmount := cryptoQuote.TargetAmount - payoutQuote.Fee
	
	estimatedMinutes := cryptoQuote.EstimatedTime/60 + 5 // Add 5 min for payout
	estimatedArrival := time.Now().Add(time.Duration(estimatedMinutes) * time.Minute)
	
	return &FullTransferResponse{
		ReferenceID:         req.ReferenceID,
		SourceAmount:        req.SourceAmount,
		SourceCurrency:      req.SourceCurrency,
		TargetAmount:        finalAmount,
		TargetCurrency:      req.TargetCurrency,
		CryptoConversionFee: cryptoQuote.ConversionFee,
		NetworkFee:          cryptoQuote.NetworkFee,
		PayoutFee:           payoutQuote.Fee,
		TotalFee:            totalFee,
		ExchangeRate:        cryptoQuote.USDToTargetRate,
		UsedCryptoRails:     true,
		UsedInternalPool:    cryptoQuote.UsedInternalPool,
		PayoutProvider:      provider.GetName(),
		PayoutMethod:        string(req.PayoutMethod),
		Status:              "quote",
		EstimatedArrival:    estimatedArrival,
	}, nil
}

// ExecuteTransfer executes the full transfer
func (o *TransferOrchestrator) ExecuteTransfer(ctx context.Context, req *FullTransferRequest) (*FullTransferResponse, error) {
	// Step 1: Execute crypto rails conversion
	cryptoReq := &ConversionRequest{
		ReferenceID:      req.ReferenceID,
		SourceAmount:     req.SourceAmount,
		SourceCurrency:   req.SourceCurrency,
		TargetCurrency:   req.TargetCurrency,
		RecipientCountry: req.RecipientCountry,
		PayoutMethod:     req.PayoutMethod,
	}
	
	cryptoResult, err := o.cryptoRails.ConvertViaStablecoin(ctx, cryptoReq)
	if err != nil {
		return nil, fmt.Errorf("crypto conversion failed: %w", err)
	}
	
	// Step 2: Execute payout via zone router
	payoutReq := &PayoutRequest{
		ReferenceID:      req.ReferenceID,
		Amount:           cryptoResult.TargetAmount,
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
	
	payoutResult, err := o.zoneRouter.CreatePayout(ctx, payoutReq)
	if err != nil {
		// TODO: Implement rollback/refund logic
		return nil, fmt.Errorf("payout failed: %w", err)
	}
	
	// Step 3: Update pool if internal was used
	if cryptoResult.UsedInternalPool {
		o.cryptoRails.UpdatePoolBalance(o.cryptoRails.config.PreferredStablecoin, -cryptoResult.StablecoinAmount)
	}
	
	// Step 4: Build response
	totalFee := cryptoResult.TotalFee + payoutResult.Fee
	estimatedArrival := time.Now().Add(time.Duration(cryptoResult.EstimatedTime) * time.Second)
	
	return &FullTransferResponse{
		ReferenceID:         req.ReferenceID,
		SourceAmount:        req.SourceAmount,
		SourceCurrency:      req.SourceCurrency,
		TargetAmount:        payoutResult.AmountReceived,
		TargetCurrency:      req.TargetCurrency,
		CryptoConversionFee: cryptoResult.ConversionFee,
		NetworkFee:          cryptoResult.NetworkFee,
		PayoutFee:           payoutResult.Fee,
		TotalFee:            totalFee,
		ExchangeRate:        cryptoResult.USDToTargetRate,
		UsedCryptoRails:     true,
		UsedInternalPool:    cryptoResult.UsedInternalPool,
		PayoutProvider:      payoutResult.ProviderName,
		PayoutMethod:        string(req.PayoutMethod),
		Status:              string(payoutResult.Status),
		EstimatedArrival:    estimatedArrival,
		CryptoTxHash:        cryptoResult.BlockchainTxHash,
		PayoutReference:     payoutResult.ProviderReference,
	}, nil
}

// GetTransferStatus gets the current status of a transfer
func (o *TransferOrchestrator) GetTransferStatus(ctx context.Context, referenceID, payoutProvider string) (*FullTransferResponse, error) {
	status, err := o.zoneRouter.GetPayoutStatus(ctx, referenceID, payoutProvider)
	if err != nil {
		return nil, err
	}
	
	return &FullTransferResponse{
		ReferenceID: referenceID,
		Status:      string(status.Status),
	}, nil
}
