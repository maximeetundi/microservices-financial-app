package models

import "time"

// AggregatorSetting represents the configuration for a payment aggregator
// Admin can enable/disable aggregators and configure their features
type AggregatorSetting struct {
	ID               string    `json:"id"`
	ProviderCode     string    `json:"provider_code"`     // stripe, paypal, flutterwave, etc.
	ProviderName     string    `json:"provider_name"`     // Display name
	LogoURL          string    `json:"logo_url"`          // Provider logo for UI
	IsEnabled        bool      `json:"is_enabled"`        // Master toggle
	DepositEnabled   bool      `json:"deposit_enabled"`   // Can receive deposits
	WithdrawEnabled  bool      `json:"withdraw_enabled"`  // Can send withdrawals
	SupportedRegions []string  `json:"supported_regions"` // Countries where available
	Priority         int       `json:"priority"`          // Higher = preferred (1-100)
	MinAmount        float64   `json:"min_amount"`        // Minimum transaction amount
	MaxAmount        float64   `json:"max_amount"`        // Maximum transaction amount (0 = unlimited)
	FeePercent       float64   `json:"fee_percent"`       // Fee percentage
	FeeFixed         float64   `json:"fee_fixed"`         // Fixed fee amount
	FeeCurrency      string    `json:"fee_currency"`      // Currency for fixed fee
	Description      string    `json:"description"`       // Admin notes
	MaintenanceMode  bool      `json:"maintenance_mode"`  // Temporary disable with message
	MaintenanceMsg   string    `json:"maintenance_msg"`   // Message to show users
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AggregatorForFrontend represents aggregator info sent to frontend
type AggregatorForFrontend struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	LogoURL         string  `json:"logo_url"`
	DepositEnabled  bool    `json:"deposit_enabled"`
	WithdrawEnabled bool    `json:"withdraw_enabled"`
	MinAmount       float64 `json:"min_amount"`
	MaxAmount       float64 `json:"max_amount"`
	FeePercent      float64 `json:"fee_percent"`
	FeeFixed        float64 `json:"fee_fixed"`
	FeeCurrency     string  `json:"fee_currency"`
	MaintenanceMode bool    `json:"maintenance_mode"`
	MaintenanceMsg  string  `json:"maintenance_msg,omitempty"`
}

// UpdateAggregatorRequest is the request body for updating aggregator settings
type UpdateAggregatorRequest struct {
	IsEnabled       *bool    `json:"is_enabled,omitempty"`
	DepositEnabled  *bool    `json:"deposit_enabled,omitempty"`
	WithdrawEnabled *bool    `json:"withdraw_enabled,omitempty"`
	Priority        *int     `json:"priority,omitempty"`
	MinAmount       *float64 `json:"min_amount,omitempty"`
	MaxAmount       *float64 `json:"max_amount,omitempty"`
	FeePercent      *float64 `json:"fee_percent,omitempty"`
	FeeFixed        *float64 `json:"fee_fixed,omitempty"`
	FeeCurrency     *string  `json:"fee_currency,omitempty"`
	Description     *string  `json:"description,omitempty"`
	MaintenanceMode *bool    `json:"maintenance_mode,omitempty"`
	MaintenanceMsg  *string  `json:"maintenance_msg,omitempty"`
}

// DefaultAggregators returns the list of default aggregators to seed
func DefaultAggregators() []AggregatorSetting {
	return []AggregatorSetting{
		{
			ProviderCode:     "stripe",
			ProviderName:     "Stripe",
			LogoURL:          "/icons/aggregators/stripe.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"US", "CA", "GB", "DE", "FR", "ES", "IT", "NL", "AU"},
			Priority:         90,
			MinAmount:        1.0,
			MaxAmount:        50000,
			FeePercent:       2.9,
			FeeFixed:         0.30,
			FeeCurrency:      "USD",
		},
		{
			ProviderCode:     "paypal",
			ProviderName:     "PayPal",
			LogoURL:          "/icons/aggregators/paypal.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"US", "CA", "GB", "DE", "FR", "ES", "IT", "AU", "BR", "MX"},
			Priority:         85,
			MinAmount:        1.0,
			MaxAmount:        10000,
			FeePercent:       3.49,
			FeeFixed:         0.49,
			FeeCurrency:      "USD",
		},
		{
			ProviderCode:     "flutterwave",
			ProviderName:     "Flutterwave",
			LogoURL:          "/icons/aggregators/flutterwave.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"NG", "GH", "KE", "UG", "TZ", "ZA", "RW", "CI", "CM", "SN"},
			Priority:         95,
			MinAmount:        100,
			MaxAmount:        5000000,
			FeePercent:       1.4,
			FeeFixed:         0,
			FeeCurrency:      "NGN",
		},
		{
			ProviderCode:     "mtn_momo",
			ProviderName:     "MTN Mobile Money",
			LogoURL:          "/icons/aggregators/mtn.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"GH", "UG", "RW", "CI", "CM", "BJ", "CG", "ZM"},
			Priority:         90,
			MinAmount:        500,
			MaxAmount:        2000000,
			FeePercent:       1.0,
			FeeFixed:         0,
			FeeCurrency:      "XOF",
		},
		{
			ProviderCode:     "orange_money",
			ProviderName:     "Orange Money",
			LogoURL:          "/icons/aggregators/orange.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"CI", "SN", "CM", "ML", "BF", "NE", "GN", "MG", "CD"},
			Priority:         88,
			MinAmount:        500,
			MaxAmount:        1500000,
			FeePercent:       1.5,
			FeeFixed:         0,
			FeeCurrency:      "XOF",
		},
		{
			ProviderCode:     "pesapal",
			ProviderName:     "Pesapal",
			LogoURL:          "/icons/aggregators/pesapal.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  false,
			SupportedRegions: []string{"KE", "UG", "TZ", "RW", "ZM", "MW"},
			Priority:         80,
			MinAmount:        10,
			MaxAmount:        500000,
			FeePercent:       3.5,
			FeeFixed:         0,
			FeeCurrency:      "KES",
		},
		{
			ProviderCode:     "chipper",
			ProviderName:     "Chipper Cash",
			LogoURL:          "/icons/aggregators/chipper.svg",
			IsEnabled:        true,
			DepositEnabled:   true,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"GH", "KE", "NG", "RW", "TZ", "UG", "ZA", "US", "GB"},
			Priority:         75,
			MinAmount:        1,
			MaxAmount:        10000,
			FeePercent:       0,
			FeeFixed:         0,
			FeeCurrency:      "USD",
		},
		{
			ProviderCode:     "thunes",
			ProviderName:     "Thunes",
			LogoURL:          "/icons/aggregators/thunes.svg",
			IsEnabled:        true,
			DepositEnabled:   false,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"*"}, // Global
			Priority:         70,
			MinAmount:        10,
			MaxAmount:        50000,
			FeePercent:       2.0,
			FeeFixed:         5,
			FeeCurrency:      "USD",
		},
		{
			ProviderCode:     "wise",
			ProviderName:     "Wise (TransferWise)",
			LogoURL:          "/icons/aggregators/wise.svg",
			IsEnabled:        false, // Disabled by default - requires business account
			DepositEnabled:   false,
			WithdrawEnabled:  true,
			SupportedRegions: []string{"*"}, // Global
			Priority:         65,
			MinAmount:        1,
			MaxAmount:        1000000,
			FeePercent:       0.5,
			FeeFixed:         0,
			FeeCurrency:      "USD",
		},
	}
}
