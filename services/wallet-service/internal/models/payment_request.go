package models

import "time"

// PaymentRequestType defines the type of payment
type PaymentRequestType string

const (
	PaymentTypeFixed    PaymentRequestType = "fixed"    // Fixed price (e.g., product)
	PaymentTypeVariable PaymentRequestType = "variable" // Customer defines amount
	PaymentTypeInvoice  PaymentRequestType = "invoice"  // Multiple items
)

// PaymentRequestStatus defines the status of a payment request
type PaymentRequestStatus string

const (
	PaymentStatusPending   PaymentRequestStatus = "pending"
	PaymentStatusPaid      PaymentRequestStatus = "paid"
	PaymentStatusExpired   PaymentRequestStatus = "expired"
	PaymentStatusCancelled PaymentRequestStatus = "cancelled"
)

// PaymentRequest represents a merchant payment request (QR code payment)
type PaymentRequest struct {
	ID          string             `json:"id" db:"id"`
	MerchantID  string             `json:"merchant_id" db:"merchant_id"`   // User ID of merchant
	WalletID    string             `json:"wallet_id" db:"wallet_id"`       // Merchant's wallet to receive
	Type        PaymentRequestType `json:"type" db:"type"`
	
	// Amount - nil for variable type (customer defines)
	Amount      *float64 `json:"amount,omitempty" db:"amount"`
	MinAmount   *float64 `json:"min_amount,omitempty" db:"min_amount"` // For variable type
	MaxAmount   *float64 `json:"max_amount,omitempty" db:"max_amount"` // For variable type
	Currency    string   `json:"currency" db:"currency"`
	
	// Description
	Title       string `json:"title" db:"title"`
	Description string `json:"description,omitempty" db:"description"`
	
	// For invoice type
	Items []PaymentItem `json:"items,omitempty" db:"-"`
	
	// QR and Link
	QRCodeData  string `json:"qr_code_data" db:"qr_code_data"`
	PaymentLink string `json:"payment_link" db:"payment_link"`
	
	// Expiration - nil means never expires
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	NeverExpires bool      `json:"never_expires" db:"never_expires"`
	
	// Status
	Status      PaymentRequestStatus `json:"status" db:"status"`
	
	// Payment info (when paid)
	PaidAmount   *float64   `json:"paid_amount,omitempty" db:"paid_amount"`
	PaidBy       *string    `json:"paid_by,omitempty" db:"paid_by"`           // Customer user ID
	PaidAt       *time.Time `json:"paid_at,omitempty" db:"paid_at"`
	TransactionID *string   `json:"transaction_id,omitempty" db:"transaction_id"`
	
	// Reusable - can be paid multiple times
	Reusable     bool `json:"reusable" db:"reusable"`
	TimesUsed    int  `json:"times_used" db:"times_used"`
	
	// Metadata
	Metadata    map[string]string `json:"metadata,omitempty" db:"-"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

// PaymentItem represents an item in an invoice payment
type PaymentItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"` // Quantity * UnitPrice
}

// PaymentHistory represents a completed payment
type PaymentHistory struct {
	ID              string    `json:"id" db:"id"`
	PaymentRequestID string   `json:"payment_request_id" db:"payment_request_id"`
	MerchantID      string    `json:"merchant_id" db:"merchant_id"`
	CustomerID      string    `json:"customer_id" db:"customer_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Fee             float64   `json:"fee" db:"fee"`
	NetAmount       float64   `json:"net_amount" db:"net_amount"`
	Currency        string    `json:"currency" db:"currency"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	PaidAt          time.Time `json:"paid_at" db:"paid_at"`
}

// ============================
// Request/Response DTOs
// ============================

// CreatePaymentRequestDTO is the request to create a payment request
type CreatePaymentRequestDTO struct {
	Type         PaymentRequestType `json:"type" binding:"required"`
	WalletID     string             `json:"wallet_id" binding:"required"`
	
	// Amount - required for fixed, optional for variable
	Amount       *float64 `json:"amount"`
	MinAmount    *float64 `json:"min_amount"`
	MaxAmount    *float64 `json:"max_amount"`
	Currency     string   `json:"currency" binding:"required"`
	
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	
	// For invoice type
	Items        []PaymentItem `json:"items"`
	
	// Expiration in minutes (0 or nil = use default, -1 = never expires)
	ExpiresInMinutes *int `json:"expires_in_minutes"`
	
	// Reusable - can be paid multiple times
	Reusable     bool `json:"reusable"`
	
	Metadata     map[string]string `json:"metadata"`
}

// PayPaymentRequestDTO is the request to pay a payment request
type PayPaymentRequestDTO struct {
	PaymentRequestID string  `json:"payment_request_id" binding:"required"`
	FromWalletID     string  `json:"from_wallet_id" binding:"required"`
	Amount           float64 `json:"amount"` // Required for variable type
	PIN              string  `json:"pin"`    // Optional security PIN
}

// PaymentRequestResponse is the response after creating a payment request
type PaymentRequestResponse struct {
	PaymentRequest *PaymentRequest `json:"payment_request"`
	QRCodeBase64   string          `json:"qr_code_base64"` // Base64 encoded PNG
	PaymentURL     string          `json:"payment_url"`
}

// QRCodeData represents the data encoded in the QR code
type QRCodeData struct {
	Type        string   `json:"type"`    // "cryptobank_payment"
	Version     int      `json:"version"` // 1
	PaymentID   string   `json:"payment_id"`
	MerchantName string  `json:"merchant_name"`
	Amount      *float64 `json:"amount,omitempty"`
	Currency    string   `json:"currency"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
}
