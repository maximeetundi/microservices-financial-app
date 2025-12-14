package models

import (
	"time"
)

type User struct {
	ID            string     `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Phone         string     `json:"phone" db:"phone"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	FirstName     string     `json:"first_name" db:"first_name"`
	LastName      string     `json:"last_name" db:"last_name"`
	DateOfBirth   *time.Time `json:"date_of_birth" db:"date_of_birth"`
	Country       string     `json:"country" db:"country"`
	KYCStatus     string     `json:"kyc_status" db:"kyc_status"`
	KYCLevel      int        `json:"kyc_level" db:"kyc_level"`
	Role          string     `json:"role" db:"role"` // user, admin, support
	IsActive      bool       `json:"is_active" db:"is_active"`
	TwoFAEnabled  bool       `json:"two_fa_enabled" db:"two_fa_enabled"`
	TwoFASecret   string     `json:"-" db:"two_fa_secret"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	PhoneVerified bool       `json:"phone_verified" db:"phone_verified"`
	LastLoginAt   *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	
	// Security fields
	FailedAttempts int        `json:"-" db:"failed_attempts"`
	LockedUntil    *time.Time `json:"-" db:"locked_until"`
}

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	SessionToken string    `json:"-" db:"session_token"`
	RefreshToken string    `json:"-" db:"refresh_token"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type VerificationToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Type      string    `json:"type" db:"type"` // email_verification, phone_verification, password_reset
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	TwoFACode string `json:"two_fa_code,omitempty"`
}

type RegisterRequest struct {
	Email       string    `json:"email" binding:"required,email"`
	Phone       string    `json:"phone" binding:"required"`
	Password    string    `json:"password" binding:"required,min=8"`
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Country     string    `json:"country" binding:"required,len=3"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         *User  `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type Enable2FAResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// BackupCode represents a one-time backup code for 2FA recovery
type BackupCode struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Code      string    `json:"-" db:"code"`       // Hashed code
	Used      bool      `json:"used" db:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents a security audit log entry
type AuditLog struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	EventType string    `json:"event_type" db:"event_type"` // login, logout, password_change, 2fa_enable, etc.
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Success   bool      `json:"success" db:"success"`
	Details   *string   `json:"details,omitempty" db:"details"` // JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ChangePasswordRequest for password changes
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}