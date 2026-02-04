package models

import "time"

// AggregatorInstance represents a specific instance of an aggregator with its own API keys
// Can be linked to MULTIPLE hot wallets
type AggregatorInstance struct {
	ID           string `json:"id"`
	AggregatorID string `json:"aggregator_id"`
	InstanceName string `json:"instance_name"`
	HotWalletID  string `json:"hot_wallet_id"`

	// API Credentials (encrypted in DB)
	APICredentials map[string]string `json:"api_credentials"`

	// Configuration
	Enabled    bool     `json:"enabled"`
	Priority   int      `json:"priority"`
	MinBalance *float64 `json:"min_balance,omitempty"`
	MaxBalance *float64 `json:"max_balance,omitempty"`

	// Transaction limits
	DailyLimit     *float64 `json:"daily_limit,omitempty"`
	MonthlyLimit   *float64 `json:"monthly_limit,omitempty"`
	DailyUsage     float64  `json:"daily_usage"`
	MonthlyUsage   float64  `json:"monthly_usage"`
	UsageResetDate string   `json:"usage_reset_date"`

	// Geographic restrictions
	RestrictedCountries []string `json:"restricted_countries,omitempty"`

	// Test/Production mode
	IsTestMode bool `json:"is_test_mode"`

	// Statistics
	TotalTransactions int        `json:"total_transactions"`
	TotalVolume       float64    `json:"total_volume"`
	LastUsedAt        *time.Time `json:"last_used_at,omitempty"`

	// Hot wallets linked to this instance
	Wallets []InstanceWallet `json:"wallets,omitempty"`

	// Metadata
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy *string   `json:"created_by,omitempty"`
}

// AggregatorInstanceWithDetails includes aggregator and hot wallet information
type AggregatorInstanceWithDetails struct {
	ID                  string     `json:"id"`
	InstanceName        string     `json:"instance_name"`
	Enabled             bool       `json:"enabled"`
	Priority            int        `json:"priority"`
	IsTestMode          bool       `json:"is_test_mode"`
	RestrictedCountries []string   `json:"restricted_countries,omitempty"`
	DailyLimit          *float64   `json:"daily_limit,omitempty"`
	MonthlyLimit        *float64   `json:"monthly_limit,omitempty"`
	DailyUsage          float64    `json:"daily_usage"`
	MonthlyUsage        float64    `json:"monthly_usage"`
	TotalTransactions   int        `json:"total_transactions"`
	TotalVolume         float64    `json:"total_volume"`
	LastUsedAt          *time.Time `json:"last_used_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Aggregator info
	AggregatorID      string `json:"aggregator_id"`
	ProviderCode      string `json:"provider_code"`
	ProviderName      string `json:"provider_name"`
	ProviderLogo      string `json:"provider_logo"`
	AggregatorEnabled bool   `json:"aggregator_enabled"`

	// Hot wallet info
	HotWalletID       string   `json:"hot_wallet_id"`
	HotWalletCurrency string   `json:"hot_wallet_currency"`
	HotWalletBalance  float64  `json:"hot_wallet_balance"`
	MinBalance        *float64 `json:"min_balance,omitempty"`
	MaxBalance        *float64 `json:"max_balance,omitempty"`

	// Availability
	AvailabilityStatus string `json:"availability_status"`
}

// CreateAggregatorInstanceRequest is the request body for creating an aggregator instance
type CreateAggregatorInstanceRequest struct {
	AggregatorID        string            `json:"aggregator_id" binding:"required"`
	InstanceName        string            `json:"instance_name" binding:"required"`
	HotWalletID         string            `json:"hot_wallet_id" binding:"required"`
	APICredentials      map[string]string `json:"api_credentials"`
	Priority            *int              `json:"priority,omitempty"`
	MinBalance          *float64          `json:"min_balance,omitempty"`
	MaxBalance          *float64          `json:"max_balance,omitempty"`
	DailyLimit          *float64          `json:"daily_limit,omitempty"`
	MonthlyLimit        *float64          `json:"monthly_limit,omitempty"`
	RestrictedCountries []string          `json:"restricted_countries,omitempty"`
	IsTestMode          *bool             `json:"is_test_mode,omitempty"`
	Notes               string            `json:"notes,omitempty"`
}

// InstanceWallet represents a hot wallet linked to an instance
type InstanceWallet struct {
	ID          string   `json:"id"`
	InstanceID  string   `json:"instance_id"`
	HotWalletID string   `json:"hot_wallet_id"`
	IsPrimary   bool     `json:"is_primary"`
	Priority    int      `json:"priority"`
	MinBalance  *float64 `json:"min_balance,omitempty"`
	MaxBalance  *float64 `json:"max_balance,omitempty"`

	// Auto-recharge configuration
	AutoRechargeEnabled    bool     `json:"auto_recharge_enabled"`
	RechargeThreshold      *float64 `json:"recharge_threshold,omitempty"`
	RechargeTarget         *float64 `json:"recharge_target,omitempty"`
	RechargeSourceWalletID *string  `json:"recharge_source_wallet_id,omitempty"`

	// Statistics
	TotalDeposits    float64    `json:"total_deposits"`
	TotalWithdrawals float64    `json:"total_withdrawals"`
	TransactionCount int        `json:"transaction_count"`
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`

	// Current wallet state (loaded from join)
	WalletCurrency string  `json:"wallet_currency,omitempty"`
	WalletBalance  float64 `json:"wallet_balance,omitempty"`

	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InstanceWalletAvailability represents availability info for a wallet
type InstanceWalletAvailability struct {
	InstanceWallet
	ProviderCode       string `json:"provider_code"`
	ProviderName       string `json:"provider_name"`
	InstanceEnabled    bool   `json:"instance_enabled"`
	AggregatorEnabled  bool   `json:"aggregator_enabled"`
	AvailabilityStatus string `json:"availability_status"`
}

// Availability status constants
const (
	WalletAvailable           = "available"
	WalletDisabled            = "wallet_disabled"
	WalletInstanceDisabled    = "instance_disabled"
	WalletAggregatorDisabled  = "aggregator_disabled"
	WalletInsufficientBalance = "insufficient_balance"
	WalletBalanceTooHigh      = "balance_too_high"
)

// CreateInstanceWalletRequest is the request to link a wallet to an instance
type CreateInstanceWalletRequest struct {
	InstanceID             string   `json:"instance_id" binding:"required"`
	HotWalletID            string   `json:"hot_wallet_id" binding:"required"`
	IsPrimary              *bool    `json:"is_primary,omitempty"`
	Priority               *int     `json:"priority,omitempty"`
	MinBalance             *float64 `json:"min_balance,omitempty"`
	MaxBalance             *float64 `json:"max_balance,omitempty"`
	AutoRechargeEnabled    *bool    `json:"auto_recharge_enabled,omitempty"`
	RechargeThreshold      *float64 `json:"recharge_threshold,omitempty"`
	RechargeTarget         *float64 `json:"recharge_target,omitempty"`
	RechargeSourceWalletID *string  `json:"recharge_source_wallet_id,omitempty"`
}

// IsAvailable checks if wallet is available for use
func (w *InstanceWalletAvailability) IsAvailable() bool {
	return w.AvailabilityStatus == WalletAvailable
}

// NeedsRecharge checks if wallet needs auto-recharge
func (w *InstanceWallet) NeedsRecharge() bool {
	if !w.AutoRechargeEnabled || w.RechargeThreshold == nil {
		return false
	}

	return w.WalletBalance < *w.RechargeThreshold
}

// GetRechargeAmount calculates how much to recharge
func (w *InstanceWallet) GetRechargeAmount() float64 {
	if !w.NeedsRecharge() || w.RechargeTarget == nil {
		return 0
	}

	return *w.RechargeTarget - w.WalletBalance
}
