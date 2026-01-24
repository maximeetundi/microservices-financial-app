package models

import "time"

type FeeConfig struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"` // unique identifier, e.g., 'transfer_internal'
	Description string    `json:"description"`
	Type        string    `json:"type"` // 'percentage', 'fixed', 'hybrid'
	Percentage  float64   `json:"percentage"`
	FixedAmount float64   `json:"fixed_amount"`
	Currency    string    `json:"currency"`
	MinFee      float64   `json:"min_fee"`
	MaxFee      float64   `json:"max_fee"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
