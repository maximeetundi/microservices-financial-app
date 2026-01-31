package providers

import (
	"context"
	"time"
)

// PayoutMethod represents the payout delivery method
type PayoutMethod string

const (
	PayoutMethodMobileMoney  PayoutMethod = "mobile_money"
	PayoutMethodBankTransfer PayoutMethod = "bank_transfer"
	PayoutMethodCashPickup   PayoutMethod = "cash_pickup"
	PayoutMethodWallet       PayoutMethod = "wallet"
	PayoutMethodCard         PayoutMethod = "card"
	PayoutMethodSEPA         PayoutMethod = "sepa"
	PayoutMethodSEPAInstant  PayoutMethod = "sepa_instant"
	PayoutMethodACH          PayoutMethod = "ach"
	PayoutMethodWire         PayoutMethod = "wire"
)

// PayoutStatus represents the status of a payout
type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
	PayoutStatusCancelled  PayoutStatus = "cancelled"
)

// PayoutStatusAccordingTo converts a string status to PayoutStatus
func PayoutStatusAccordingTo(status string) PayoutStatus {
	switch status {
	case "completed", "success", "paid", "COMPLETED", "SUCCESS", "PAID":
		return PayoutStatusCompleted
	case "pending", "processing", "initiated", "PENDING", "PROCESSING", "INITIATED":
		return PayoutStatusPending
	case "failed", "error", "rejected", "FAILED", "ERROR", "REJECTED":
		return PayoutStatusFailed
	case "cancelled", "CANCELLED":
		return PayoutStatusCancelled
	default:
		return PayoutStatusPending
	}
}

// PayoutRequest represents a payout request to a provider
type PayoutRequest struct {
	// Transaction details
	ReferenceID string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`

	// Recipient details
	RecipientName    string `json:"recipient_name"`
	RecipientPhone   string `json:"recipient_phone,omitempty"`
	RecipientEmail   string `json:"recipient_email,omitempty"`
	RecipientCountry string `json:"recipient_country"`

	// Payout method details
	PayoutMethod PayoutMethod `json:"payout_method"`

	// Bank details (for bank transfers)
	BankCode      string `json:"bank_code,omitempty"`
	BankName      string `json:"bank_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	IBAN          string `json:"iban,omitempty"`
	SwiftCode     string `json:"swift_code,omitempty"`
	RoutingNumber string `json:"routing_number,omitempty"`

	// Mobile Money details
	MobileOperator string `json:"mobile_operator,omitempty"`
	MobileNumber   string `json:"mobile_number,omitempty"`

	// Additional info
	Narration string            `json:"narration,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// PayoutResponse represents the response from a payout provider
type PayoutResponse struct {
	ProviderName      string       `json:"provider_name"`
	ProviderReference string       `json:"provider_reference"`
	ReferenceID       string       `json:"reference_id"`
	Status            PayoutStatus `json:"status"`
	Message           string       `json:"message,omitempty"`
	Fee               float64      `json:"fee"`
	ExchangeRate      float64      `json:"exchange_rate,omitempty"`
	AmountReceived    float64      `json:"amount_received,omitempty"`
	ReceivedCurrency  string       `json:"received_currency,omitempty"`
	EstimatedDelivery time.Time    `json:"estimated_delivery,omitempty"`
}

// PayoutStatusResponse represents status check response
type PayoutStatusResponse struct {
	ReferenceID       string       `json:"reference_id"`
	ProviderReference string       `json:"provider_reference"`
	Status            PayoutStatus `json:"status"`
	Message           string       `json:"message,omitempty"`
	CompletedAt       *time.Time   `json:"completed_at,omitempty"`
}

// AvailableMethod represents an available payout method for a country
type AvailableMethod struct {
	Code             string       `json:"code"`
	Name             string       `json:"name"`
	Type             string       `json:"type"` // "mobile", "bank", "wallet", etc.
	Method           PayoutMethod `json:"method"`
	EstimatedMinutes int          `json:"estimated_minutes"`
	Fee              float64      `json:"fee"`
	FeeType          string       `json:"fee_type"` // "flat" or "percentage"
	MinAmount        float64      `json:"min_amount"`
	MaxAmount        float64      `json:"max_amount"`
	Countries        []string     `json:"countries,omitempty"`
	Currencies       []string     `json:"currencies,omitempty"`
	RequiredFields   []string     `json:"required_fields,omitempty"`
}

// MobileOperator represents a Mobile Money operator
type MobileOperator struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	Countries    []string `json:"countries"`
	NumberPrefix []string `json:"number_prefix,omitempty"`
}

// Bank represents a bank in a country
type Bank struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// PayoutProvider is the interface that all payout providers must implement
type PayoutProvider interface {
	// GetName returns the provider name
	GetName() string

	// GetSupportedCountries returns list of supported countries (ISO 2-letter codes)
	GetSupportedCountries() []string

	// GetAvailableMethods returns available payout methods for a country
	GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error)

	// GetBanks returns list of banks for a country (for bank transfers)
	GetBanks(ctx context.Context, country string) ([]Bank, error)

	// GetMobileOperators returns list of mobile operators for a country
	GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error)

	// ValidateRecipient validates recipient details before payout
	ValidateRecipient(ctx context.Context, req *PayoutRequest) error

	// GetQuote gets a quote for the payout (fees, exchange rate, etc.)
	GetQuote(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error)

	// CreatePayout creates a payout transaction
	CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error)

	// GetPayoutStatus checks the status of a payout
	GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error)

	// CancelPayout cancels a pending payout (if supported)
	CancelPayout(ctx context.Context, referenceID string) error
}
