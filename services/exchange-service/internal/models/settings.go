package models

import (
	"time"
)

// SystemSetting represents a configurable system setting
type SystemSetting struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`         // unique identifier, e.g., 'rate_update_interval'
	Value       string    `json:"value"`       // stored as string, parsed by service
	Description string    `json:"description"` // human-readable description
	Type        string    `json:"type"`        // 'integer', 'string', 'boolean'
	Category    string    `json:"category"`    // 'exchange', 'crypto', 'general'
	UpdatedAt   time.Time `json:"updated_at"`
}
