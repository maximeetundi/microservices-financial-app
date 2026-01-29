package models

import (
	"time"
)

// SystemConfig represents a cached system setting
type SystemConfig struct {
	Key              string    `json:"key" db:"key"`
	Value            string    `json:"value" db:"value"`
	FixedAmount      float64   `json:"fixed_amount" db:"fixed_amount"`
	PercentageAmount float64   `json:"percentage_amount" db:"percentage_amount"`
	IsEnabled        bool      `json:"is_enabled" db:"is_enabled"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserUsageStats tracks a user's transaction usage for limit enforcement
type UserUsageStats struct {
	UserID           string    `json:"user_id" db:"user_id"`
	PeriodKey        string    `json:"period_key" db:"period_key"` // e.g., "DAILY_2024-01-29", "MONTHLY_2024-01"
	Currency         string    `json:"currency" db:"currency"`
	TotalAmount      float64   `json:"total_amount" db:"total_amount"`
	TransactionCount int       `json:"transaction_count" db:"transaction_count"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
