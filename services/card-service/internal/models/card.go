package models

import (
	"time"
)

type Card struct {
	ID                string     `json:"id" db:"id"`
	UserID            string     `json:"user_id" db:"user_id"`
	CardNumber        string     `json:"card_number" db:"card_number"`           // Masqué: ****-****-****-1234
	CardNumberFull    string     `json:"-" db:"card_number_full"`                // Complet, chiffré
	CardType          string     `json:"card_type" db:"card_type"`               // prepaid, virtual, gift
	CardCategory      string     `json:"card_category" db:"card_category"`       // personal, business
	Currency          string     `json:"currency" db:"currency"`
	Balance           float64    `json:"balance" db:"balance"`
	AvailableBalance  float64    `json:"available_balance" db:"available_balance"`
	CardholderName    string     `json:"cardholder_name" db:"cardholder_name"`
	ExpiryMonth       int        `json:"expiry_month" db:"expiry_month"`
	ExpiryYear        int        `json:"expiry_year" db:"expiry_year"`
	CVV               string     `json:"-" db:"cvv"`                             // Chiffré
	PINHash           string     `json:"-" db:"pin_hash"`                        // Haché
	Status            string     `json:"status" db:"status"`                     // active, inactive, blocked, expired, frozen
	IsVirtual         bool       `json:"is_virtual" db:"is_virtual"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	ActivatedAt       *time.Time `json:"activated_at,omitempty" db:"activated_at"`
	ExpiresAt         time.Time  `json:"expires_at" db:"expires_at"`
	
	// Freeze/Block info
	FreezeReason      *string    `json:"freeze_reason,omitempty" db:"freeze_reason"`
	FrozenAt          *time.Time `json:"frozen_at,omitempty" db:"frozen_at"`
	BlockReason       *string    `json:"block_reason,omitempty" db:"block_reason"`
	BlockedAt         *time.Time `json:"blocked_at,omitempty" db:"blocked_at"`
	
	// Limites
	DailyLimit        float64 `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit      float64 `json:"monthly_limit" db:"monthly_limit"`
	SingleTxLimit     float64 `json:"single_tx_limit" db:"single_tx_limit"`
	ATMDailyLimit     float64 `json:"atm_daily_limit" db:"atm_daily_limit"`
	OnlineTxLimit     float64 `json:"online_tx_limit" db:"online_tx_limit"`
	
	// Utilisation actuelle
	DailySpent        float64 `json:"daily_spent" db:"daily_spent"`
	MonthlySpent      float64 `json:"monthly_spent" db:"monthly_spent"`
	ATMDailySpent     float64 `json:"atm_daily_spent" db:"atm_daily_spent"`
	
	// Paramètres
	AllowATM          bool `json:"allow_atm" db:"allow_atm"`
	AllowOnline       bool `json:"allow_online" db:"allow_online"`
	AllowInternational bool `json:"allow_international" db:"allow_international"`
	AllowContactless  bool `json:"allow_contactless" db:"allow_contactless"`
	
	// Auto-reload
	AutoReloadEnabled   bool     `json:"auto_reload_enabled" db:"auto_reload_enabled"`
	AutoReloadAmount    float64  `json:"auto_reload_amount" db:"auto_reload_amount"`
	AutoReloadThreshold float64  `json:"auto_reload_threshold" db:"auto_reload_threshold"`
	ReloadWalletID      *string  `json:"reload_wallet_id,omitempty" db:"reload_wallet_id"`
	AutoReloadWalletID  *string  `json:"auto_reload_wallet_id,omitempty" db:"auto_reload_wallet_id"`
	
	// Carte physique
	ShippingAddress   *string    `json:"shipping_address,omitempty" db:"shipping_address"`
	ShippingStatus    *string    `json:"shipping_status,omitempty" db:"shipping_status"`
	TrackingNumber    *string    `json:"tracking_number,omitempty" db:"tracking_number"`
	ShippedAt         *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	DeliveredAt       *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	
	// Externe
	ExternalCardID    *string `json:"external_card_id,omitempty" db:"external_card_id"` // ID chez le processeur
	IssuerID          string  `json:"issuer_id" db:"issuer_id"`                         // marqeta, etc.
	
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type CardTransaction struct {
	ID                    string     `json:"id" db:"id"`
	CardID                string     `json:"card_id" db:"card_id"`
	UserID                string     `json:"user_id" db:"user_id"`
	TransactionType       string     `json:"transaction_type" db:"transaction_type"` // purchase, withdrawal, load, refund
	Amount                float64    `json:"amount" db:"amount"`
	Currency              string     `json:"currency" db:"currency"`
	OriginalAmount        *float64   `json:"original_amount,omitempty" db:"original_amount"`
	OriginalCurrency      *string    `json:"original_currency,omitempty" db:"original_currency"`
	ExchangeRate          *float64   `json:"exchange_rate,omitempty" db:"exchange_rate"`
	Fee                   float64    `json:"fee" db:"fee"`
	MerchantName          *string    `json:"merchant_name,omitempty" db:"merchant_name"`
	MerchantCategory      *string    `json:"merchant_category,omitempty" db:"merchant_category"`
	MerchantCity          *string    `json:"merchant_city,omitempty" db:"merchant_city"`
	MerchantCountry       *string    `json:"merchant_country,omitempty" db:"merchant_country"`
	AuthorizationCode     *string    `json:"authorization_code,omitempty" db:"authorization_code"`
	ReferenceNumber       *string    `json:"reference_number,omitempty" db:"reference_number"`
	ExternalTransactionID *string    `json:"external_transaction_id,omitempty" db:"external_transaction_id"`
	Status                string     `json:"status" db:"status"` // pending, approved, declined, reversed
	DeclineReason         *string    `json:"decline_reason,omitempty" db:"decline_reason"`
	IsOnline              bool       `json:"is_online" db:"is_online"`
	IsInternational       bool       `json:"is_international" db:"is_international"`
	IsContactless         bool       `json:"is_contactless" db:"is_contactless"`
	ProcessedAt           *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

type GiftCard struct {
	ID              string     `json:"id" db:"id"`
	Code            string     `json:"code" db:"code"`           // Code de la gift card
	SenderID        string     `json:"sender_id" db:"sender_id"` // Qui l'envoie
	RecipientEmail  *string    `json:"recipient_email,omitempty" db:"recipient_email"`
	RecipientPhone  *string    `json:"recipient_phone,omitempty" db:"recipient_phone"`
	Amount          float64    `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	Message         *string    `json:"message,omitempty" db:"message"`
	Design          string     `json:"design" db:"design"`       // birthday, christmas, etc.
	Status          string     `json:"status" db:"status"`       // pending, sent, redeemed, expired
	RedeemedBy      *string    `json:"redeemed_by,omitempty" db:"redeemed_by"`
	RedeemedAt      *time.Time `json:"redeemed_at,omitempty" db:"redeemed_at"`
	ExpiresAt       time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// Request models
type CreateCardRequest struct {
	CardType         string  `json:"card_type" binding:"required,oneof=prepaid virtual gift"`
	CardCategory     string  `json:"card_category" binding:"required,oneof=personal business"`
	Currency         string  `json:"currency" binding:"required"`
	CardholderName   string  `json:"cardholder_name" binding:"required"`
	InitialAmount    float64 `json:"initial_amount" binding:"min=0"`
	SourceWalletID   *string `json:"source_wallet_id,omitempty"`
	
	// Limites personnalisées
	DailyLimit       *float64 `json:"daily_limit,omitempty"`
	MonthlyLimit     *float64 `json:"monthly_limit,omitempty"`
	SingleTxLimit    *float64 `json:"single_tx_limit,omitempty"`
	ATMDailyLimit    *float64 `json:"atm_daily_limit,omitempty"`
	
	// Paramètres
	AllowATM         *bool `json:"allow_atm,omitempty"`
	AllowOnline      *bool `json:"allow_online,omitempty"`
	AllowInternational *bool `json:"allow_international,omitempty"`
	AllowContactless *bool `json:"allow_contactless,omitempty"`
}

type CreateVirtualCardRequest struct {
	Currency         string  `json:"currency" binding:"required"`
	CardholderName   string  `json:"cardholder_name" binding:"required"`
	InitialAmount    float64 `json:"initial_amount" binding:"min=0"`
	SourceWalletID   string  `json:"source_wallet_id" binding:"required"`
	ValidityMonths   *int    `json:"validity_months,omitempty"` // 12, 24, 36, 48
	Purpose          string  `json:"purpose" binding:"required,oneof=online_shopping subscription single_use travel"`
}

type OrderPhysicalCardRequest struct {
	Currency         string `json:"currency" binding:"required"`
	CardholderName   string `json:"cardholder_name" binding:"required"`
	InitialAmount    float64 `json:"initial_amount" binding:"min=0"`
	SourceWalletID   string `json:"source_wallet_id" binding:"required"`
	
	// Adresse de livraison
	ShippingAddress  ShippingAddress `json:"shipping_address" binding:"required"`
	ExpressShipping  bool            `json:"express_shipping"`
	
	// Design de la carte
	CardDesign       string `json:"card_design" binding:"required"`
}

type ShippingAddress struct {
	FullName    string `json:"full_name" binding:"required"`
	AddressLine1 string `json:"address_line1" binding:"required"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required,len=3"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoadCardRequest struct {
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	SourceWalletID string  `json:"source_wallet_id" binding:"required"`
	Description    string  `json:"description,omitempty"`
}

type SetupAutoLoadRequest struct {
	ReloadAmount     float64 `json:"reload_amount" binding:"required,gt=0"`
	ReloadThreshold  float64 `json:"reload_threshold" binding:"required,gt=0"`
	SourceWalletID   string  `json:"source_wallet_id" binding:"required"`
	MaxReloadsPerDay *int    `json:"max_reloads_per_day,omitempty"`
}

type UpdateCardRequest struct {
	Name               *string `json:"name,omitempty"`
	AllowATM           *bool   `json:"allow_atm,omitempty"`
	AllowOnline        *bool   `json:"allow_online,omitempty"`
	AllowInternational *bool   `json:"allow_international,omitempty"`
	AllowContactless   *bool   `json:"allow_contactless,omitempty"`
}

type UpdateCardLimitsRequest struct {
	DailyLimit        *float64 `json:"daily_limit,omitempty"`
	MonthlyLimit      *float64 `json:"monthly_limit,omitempty"`
	SingleTxLimit     *float64 `json:"single_tx_limit,omitempty"`
	ATMDailyLimit     *float64 `json:"atm_daily_limit,omitempty"`
	OnlineTxLimit     *float64 `json:"online_tx_limit,omitempty"`
	
	AllowATM          *bool `json:"allow_atm,omitempty"`
	AllowOnline       *bool `json:"allow_online,omitempty"`
	AllowInternational *bool `json:"allow_international,omitempty"`
	AllowContactless  *bool `json:"allow_contactless,omitempty"`
}

type SetCardPINRequest struct {
	PIN string `json:"pin" binding:"required,len=4"`
}

type ChangeCardPINRequest struct {
	CurrentPIN string `json:"current_pin" binding:"required,len=4"`
	NewPIN     string `json:"new_pin" binding:"required,len=4"`
}

type ShippingStatus struct {
	Status         string     `json:"status"`
	TrackingNumber *string    `json:"tracking_number,omitempty"`
	Carrier        string     `json:"carrier,omitempty"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	EstimatedDelivery *time.Time `json:"estimated_delivery,omitempty"`
}

type CreateGiftCardRequest struct {
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required"`
	RecipientEmail *string `json:"recipient_email,omitempty"`
	RecipientPhone *string `json:"recipient_phone,omitempty"`
	Message        *string `json:"message,omitempty"`
	Design         string  `json:"design" binding:"required"`
	ValidityDays   *int    `json:"validity_days,omitempty"` // Par défaut 365 jours
	SourceWalletID string  `json:"source_wallet_id" binding:"required"`
}

type RedeemGiftCardRequest struct {
	Code           string `json:"code" binding:"required"`
	TargetWalletID string `json:"target_wallet_id" binding:"required"`
}

// Response models
type CardDetailsResponse struct {
	Card           Card              `json:"card"`
	CardNumber     string            `json:"card_number"`      // Démasqué pour le propriétaire
	CVV            string            `json:"cvv"`              // Démasqué pour le propriétaire
	QRCode         *string           `json:"qr_code,omitempty"`
	ApplePayToken  *string           `json:"apple_pay_token,omitempty"`
	GooglePayToken *string           `json:"google_pay_token,omitempty"`
}

type CardBalance struct {
	Available     float64 `json:"available"`
	Pending       float64 `json:"pending"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
	LastUpdated   time.Time `json:"last_updated"`
}

type CardLimits struct {
	DailyLimit        float64 `json:"daily_limit"`
	DailyRemaining    float64 `json:"daily_remaining"`
	MonthlyLimit      float64 `json:"monthly_limit"`
	MonthlyRemaining  float64 `json:"monthly_remaining"`
	SingleTxLimit     float64 `json:"single_tx_limit"`
	ATMDailyLimit     float64 `json:"atm_daily_limit"`
	ATMDailyRemaining float64 `json:"atm_daily_remaining"`
	OnlineTxLimit     float64 `json:"online_tx_limit"`
	Currency          string  `json:"currency"`
}

type CardFees struct {
	IssuanceFee       float64            `json:"issuance_fee"`
	MonthlyFee        float64            `json:"monthly_fee"`
	ATMWithdrawalFee  float64            `json:"atm_withdrawal_fee"`
	ForeignTxFee      float64            `json:"foreign_tx_fee"`
	ReloadFee         float64            `json:"reload_fee"`
	ReplacementFee    float64            `json:"replacement_fee"`
	CurrencyFees      map[string]float64 `json:"currency_fees"`
}

type CardUsageStatistics struct {
	TotalTransactions   int                    `json:"total_transactions"`
	TotalSpent          float64                `json:"total_spent"`
	AverageTransaction  float64                `json:"average_transaction"`
	TopMerchants        []MerchantSpending     `json:"top_merchants"`
	SpendingByCategory  map[string]float64     `json:"spending_by_category"`
	MonthlySpending     []MonthlySpending      `json:"monthly_spending"`
	TransactionsByType  map[string]int         `json:"transactions_by_type"`
}

type MerchantSpending struct {
	MerchantName string  `json:"merchant_name"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
}

type MonthlySpending struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}