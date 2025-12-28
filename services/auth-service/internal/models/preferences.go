package models

import (
	"time"
)

// KYCDocument represents a document uploaded for KYC verification
type KYCDocument struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	Type            string     `json:"type" db:"type"` // identity, selfie, address
	FileName        string     `json:"file_name" db:"file_name"`
	FilePath        string     `json:"-" db:"file_path"`
	FileSize        int64      `json:"file_size" db:"file_size"`
	MimeType        string     `json:"mime_type" db:"mime_type"`
	DocumentNumber  *string    `json:"document_number,omitempty" db:"document_number"` // ID/Passport/License number
	ExpiryDate      *string    `json:"expiry_date,omitempty" db:"expiry_date"`         // Document expiration date (YYYY-MM-DD)
	Status          string     `json:"status" db:"status"` // pending, approved, rejected
	RejectionReason *string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewedBy      *string    `json:"reviewed_by,omitempty" db:"reviewed_by"`
	UploadedAt      time.Time  `json:"uploaded_at" db:"uploaded_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// UserPreferences stores user application preferences
type UserPreferences struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Theme           string    `json:"theme" db:"theme"` // dark, light, system
	Language        string    `json:"language" db:"language"` // fr, en, es, ar
	Currency        string    `json:"currency" db:"currency"` // USD, EUR, XOF, etc.
	Timezone        string    `json:"timezone" db:"timezone"`
	NumberFormat    string    `json:"number_format" db:"number_format"` // fr, en
	DateFormat      string    `json:"date_format" db:"date_format"` // DD/MM/YYYY, MM/DD/YYYY
	HideBalances    bool      `json:"hide_balances" db:"hide_balances"`
	AnalyticsEnabled bool     `json:"analytics_enabled" db:"analytics_enabled"`
	AutoLockMinutes int       `json:"auto_lock_minutes" db:"auto_lock_minutes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NotificationPreferences stores notification settings
type NotificationPreferences struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	
	// Push notifications
	PushEnabled     bool      `json:"push_enabled" db:"push_enabled"`
	
	// Transaction alerts
	TransferReceived bool     `json:"transfer_received" db:"transfer_received"`
	TransferSent     bool     `json:"transfer_sent" db:"transfer_sent"`
	CardPayment      bool     `json:"card_payment" db:"card_payment"`
	LowBalance       bool     `json:"low_balance" db:"low_balance"`
	
	// Security alerts
	NewLogin         bool     `json:"new_login" db:"new_login"`
	PasswordChange   bool     `json:"password_change" db:"password_change"`
	OtpViaSMS        bool     `json:"otp_via_sms" db:"otp_via_sms"`
	
	// Email preferences
	WeeklyReport     bool     `json:"weekly_report" db:"weekly_report"`
	Newsletter       bool     `json:"newsletter" db:"newsletter"`
	Promotions       bool     `json:"promotions" db:"promotions"`
	
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UploadKYCDocumentRequest for uploading KYC documents
type UploadKYCDocumentRequest struct {
	Type string `form:"type" binding:"required,oneof=identity selfie address"`
}

// UpdatePreferencesRequest for updating user preferences
type UpdatePreferencesRequest struct {
	Theme           *string `json:"theme,omitempty"`
	Language        *string `json:"language,omitempty"`
	Currency        *string `json:"currency,omitempty"`
	Timezone        *string `json:"timezone,omitempty"`
	NumberFormat    *string `json:"number_format,omitempty"`
	DateFormat      *string `json:"date_format,omitempty"`
	HideBalances    *bool   `json:"hide_balances,omitempty"`
	AnalyticsEnabled *bool  `json:"analytics_enabled,omitempty"`
	AutoLockMinutes *int    `json:"auto_lock_minutes,omitempty"`
}

// UpdateNotificationPrefsRequest for updating notification preferences
type UpdateNotificationPrefsRequest struct {
	PushEnabled      *bool `json:"push_enabled,omitempty"`
	TransferReceived *bool `json:"transfer_received,omitempty"`
	TransferSent     *bool `json:"transfer_sent,omitempty"`
	CardPayment      *bool `json:"card_payment,omitempty"`
	LowBalance       *bool `json:"low_balance,omitempty"`
	NewLogin         *bool `json:"new_login,omitempty"`
	PasswordChange   *bool `json:"password_change,omitempty"`
	OtpViaSMS        *bool `json:"otp_via_sms,omitempty"`
	WeeklyReport     *bool `json:"weekly_report,omitempty"`
	Newsletter       *bool `json:"newsletter,omitempty"`
	Promotions       *bool `json:"promotions,omitempty"`
}

// KYCStatusResponse for KYC status endpoint
type KYCStatusResponse struct {
	Status     string         `json:"status"` // pending, submitted, verified, rejected
	Level      int            `json:"level"`
	Documents  []KYCDocument  `json:"documents"`
	Message    string         `json:"message,omitempty"`
}
