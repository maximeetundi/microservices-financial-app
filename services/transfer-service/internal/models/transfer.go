package models

import (
	"time"
)

type Transfer struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	FromWalletID    *string    `json:"from_wallet_id,omitempty" db:"from_wallet_id"`
	ToWalletID      *string    `json:"to_wallet_id,omitempty" db:"to_wallet_id"`
	RecipientEmail  *string    `json:"recipient_email,omitempty" db:"recipient_email"`
	RecipientPhone  *string    `json:"recipient_phone,omitempty" db:"recipient_phone"`
	Amount          float64    `json:"amount" db:"amount"`
	Fee             float64    `json:"fee" db:"fee"`
	Currency        string     `json:"currency" db:"currency"`
	ExchangeRate    *float64   `json:"exchange_rate,omitempty" db:"exchange_rate"`
	DestinationAmount *float64 `json:"destination_amount,omitempty" db:"destination_amount"`
	DestinationCurrency *string `json:"destination_currency,omitempty" db:"destination_currency"`
	TransferType    string     `json:"transfer_type" db:"transfer_type"` // domestic, international, mobile_money
	Status          string     `json:"status" db:"status"` // pending, processing, completed, failed, cancelled
	Description     *string    `json:"description,omitempty" db:"description"`
	Reference       *string    `json:"reference,omitempty" db:"reference"`
	ExternalRef     *string    `json:"external_ref,omitempty" db:"external_ref"`
	ProviderTxID    *string    `json:"provider_tx_id,omitempty" db:"provider_tx_id"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	
	// Additional fields for different transfer types
	InternationalDetails *InternationalTransferDetails `json:"international_details,omitempty" db:"-"`
	MobileMoneyDetails   *MobileMoneyDetails          `json:"mobile_money_details,omitempty" db:"-"`

	// Enriched Data (not stored in DB)
	SenderDetails    *UserDetails `json:"sender_details,omitempty" db:"-"`
	RecipientDetails *UserDetails `json:"recipient_details,omitempty" db:"-"`
}

type UserDetails struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type InternationalTransferDetails struct {
	RecipientName     string `json:"recipient_name" binding:"required"`
	RecipientAddress  string `json:"recipient_address" binding:"required"`
	RecipientCountry  string `json:"recipient_country" binding:"required,len=3"`
	BankName          string `json:"bank_name" binding:"required"`
	BankAccount       string `json:"bank_account" binding:"required"`
	SwiftCode         string `json:"swift_code,omitempty"`
	RoutingNumber     string `json:"routing_number,omitempty"`
	Purpose           string `json:"purpose" binding:"required"`
}

type MobileMoneyDetails struct {
	Provider      string `json:"provider" binding:"required"`
	RecipientPhone string `json:"recipient_phone" binding:"required"`
	Country       string `json:"country" binding:"required,len=3"`
}

type BulkTransfer struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	FromWalletID    *string   `json:"from_wallet_id,omitempty" db:"from_wallet_id"`
	Currency        string    `json:"currency" db:"currency"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	TotalFee        float64   `json:"total_fee" db:"total_fee"`
	TotalCount      int       `json:"total_count" db:"total_count"`
	ProcessedCount  int       `json:"processed_count" db:"processed_count"`
	SuccessCount    int       `json:"success_count" db:"success_count"`
	FailedCount     int       `json:"failed_count" db:"failed_count"`
	Status          string    `json:"status" db:"status"` // pending, processing, completed, failed, cancelled
	Description     *string   `json:"description,omitempty" db:"description"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	
	// Individual transfers
	Transfers []Transfer `json:"transfers,omitempty" db:"-"`
}

type MobileMoneyProvider struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Countries   []string `json:"countries"`
	Currencies  []string `json:"currencies"`
	MinAmount   float64  `json:"min_amount"`
	MaxAmount   float64  `json:"max_amount"`
	Fee         float64  `json:"fee"`
	IsAvailable bool     `json:"is_available"`
}

type TransferRequest struct {
	FromWalletID   string  `json:"from_wallet_id" binding:"required"`
	ToWalletID     *string `json:"to_wallet_id,omitempty"`
	ToEmail        *string `json:"to_email,omitempty"`
	ToPhone        *string `json:"to_phone,omitempty"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required"`
	Description    *string `json:"description,omitempty"`
	Reference      *string `json:"reference,omitempty"`
	Type           string  `json:"type,omitempty"` // transfer_domestic, transfer_p2p, etc.
}

type InternationalTransferRequest struct {
	FromWalletID         string  `json:"from_wallet_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required,gt=0"`
	Currency             string  `json:"currency" binding:"required"`
	DestinationCurrency  string  `json:"destination_currency" binding:"required"`
	RecipientName        string  `json:"recipient_name" binding:"required"`
	RecipientAddress     string  `json:"recipient_address" binding:"required"`
	RecipientCountry     string  `json:"recipient_country" binding:"required,len=3"`
	BankName             string  `json:"bank_name" binding:"required"`
	BankAccount          string  `json:"bank_account" binding:"required"`
	SwiftCode            string  `json:"swift_code,omitempty"`
	RoutingNumber        string  `json:"routing_number,omitempty"`
	Purpose              string  `json:"purpose" binding:"required"`
}

type MobileMoneyTransferRequest struct {
	FromWalletID   string  `json:"from_wallet_id" binding:"required"`
	ToPhone        string  `json:"to_phone" binding:"required"`
	Provider       string  `json:"provider" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required"`
	Country        string  `json:"country" binding:"required,len=3"`
}

type BulkTransferRequest struct {
	FromWalletID string `json:"from_wallet_id" binding:"required"`
	Currency     string `json:"currency" binding:"required"`
	Description  *string `json:"description,omitempty"`
	Transfers    []struct {
		ToEmail     *string `json:"to_email,omitempty"`
		ToPhone     *string `json:"to_phone,omitempty"`
		ToWalletID  *string `json:"to_wallet_id,omitempty"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Reference   *string `json:"reference,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"transfers" binding:"required,min=1,max=1000"`
}

type TransferQuote struct {
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	DestinationAmount   float64 `json:"destination_amount"`
	DestinationCurrency string  `json:"destination_currency"`
	ExchangeRate        float64 `json:"exchange_rate"`
	Fee                 float64 `json:"fee"`
	TotalCost           float64 `json:"total_cost"`
	EstimatedDelivery   string  `json:"estimated_delivery"`
	ValidUntil          time.Time `json:"valid_until"`
}

type TransferLimits struct {
	Currency       string  `json:"currency"`
	DailyLimit     float64 `json:"daily_limit"`
	MonthlyLimit   float64 `json:"monthly_limit"`
	RemainingDaily float64 `json:"remaining_daily"`
	RemainingMonthly float64 `json:"remaining_monthly"`
	SingleLimit    float64 `json:"single_limit"`
}

type ComplianceCheck struct {
	ID           string    `json:"id" db:"id"`
	TransferID   string    `json:"transfer_id" db:"transfer_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	CheckType    string    `json:"check_type" db:"check_type"` // kyc, aml, sanctions, pep
	Status       string    `json:"status" db:"status"` // pending, passed, failed, manual_review
	RiskScore    *int      `json:"risk_score,omitempty" db:"risk_score"`
	Details      *string   `json:"details,omitempty" db:"details"` // JSON
	ReviewedBy   *string   `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Wallet represents a user wallet
type Wallet struct {
	ID            string  `json:"id" db:"id"`
	UserID        string  `json:"user_id" db:"user_id"`
	Currency      string  `json:"currency" db:"currency"`
	WalletType    string  `json:"wallet_type" db:"wallet_type"`
	Balance       float64 `json:"balance" db:"balance"`
	FrozenBalance float64 `json:"frozen_balance" db:"frozen_balance"`
	IsActive      bool    `json:"is_active" db:"is_active"`
}

// MobileMoneyRequest for mobile money operations
type MobileMoneyRequest struct {
	Provider   string  `json:"provider" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Currency   string  `json:"currency" binding:"required"`
	WalletID   string  `json:"wallet_id" binding:"required"`
}

// MobileMoneyResponse for mobile money operations
type MobileMoneyResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Provider      string `json:"provider"`
	Message       string `json:"message,omitempty"`
}

// ComplianceResult represents compliance check result
type ComplianceResult struct {
	Passed      bool     `json:"passed"`
	RiskScore   int      `json:"risk_score"`
	Checks      []string `json:"checks"`
	RequiresKYC bool     `json:"requires_kyc"`
	Message     string   `json:"message,omitempty"`
}

// Type alias for backward compatibility
func (t *Transfer) GetType() string {
	if t.TransferType != "" {
		return t.TransferType
	}
	return "domestic"
}

// Metadata field for JSON storage
type Metadata map[string]interface{}

// ReferenceID getter
func (t *Transfer) GetReferenceID() string {
	if t.Reference != nil {
		return *t.Reference
	}
	return ""
}