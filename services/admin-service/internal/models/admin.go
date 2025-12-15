package models

import (
	"time"

	"github.com/google/uuid"
)

// AdminUser represents an administrator
type AdminUser struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	RoleID       string     `json:"role_id"`
	Role         *AdminRole `json:"role,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
}

// AdminRole defines a role with permissions
type AdminRole struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions []AdminPermission `json:"permissions"`
	IsSystem    bool              `json:"is_system"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AdminPermission defines a specific permission
type AdminPermission struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// Predefined permission codes
const (
	// User Management
	PermViewUsers       = "users.view"
	PermCreateUsers     = "users.create"
	PermUpdateUsers     = "users.update"
	PermBlockUsers      = "users.block"
	PermDeleteUsers     = "users.delete"
	
	// KYC Management
	PermViewKYC         = "kyc.view"
	PermApproveKYC      = "kyc.approve"
	PermRejectKYC       = "kyc.reject"
	
	// Transaction Management
	PermViewTransactions    = "transactions.view"
	PermBlockTransactions   = "transactions.block"
	PermRefundTransactions  = "transactions.refund"
	PermApproveTransactions = "transactions.approve"
	
	// Card Management
	PermViewCards       = "cards.view"
	PermFreezeCards     = "cards.freeze"
	PermBlockCards      = "cards.block"
	PermReplaceCards    = "cards.replace"
	
	// Wallet Management
	PermViewWallets     = "wallets.view"
	PermFreezeWallets   = "wallets.freeze"
	PermAdjustBalances  = "wallets.adjust"
	
	// Transfer Management
	PermViewTransfers   = "transfers.view"
	PermBlockTransfers  = "transfers.block"
	PermApproveTransfers = "transfers.approve"
	
	// Exchange Management
	PermViewExchanges   = "exchanges.view"
	PermSetRates        = "exchanges.rates"
	
	// System Management
	PermViewSystem      = "system.view"
	PermViewLogs        = "system.logs"
	PermManageSettings  = "system.settings"
	
	// Admin Management
	PermViewAdmins      = "admins.view"
	PermCreateAdmins    = "admins.create"
	PermUpdateAdmins    = "admins.update"
	PermDeleteAdmins    = "admins.delete"
	PermManageRoles     = "admins.roles"
	
	// Analytics
	PermViewAnalytics   = "analytics.view"
	PermExportReports   = "analytics.export"
)

// Predefined roles
var (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleSupport    = "support"
	RoleCompliance = "compliance"
	RoleAnalyst    = "analyst"
	RoleViewer     = "viewer"
)

// AdminAuditLog tracks all admin actions
type AdminAuditLog struct {
	ID          string    `json:"id"`
	AdminID     string    `json:"admin_id"`
	AdminEmail  string    `json:"admin_email"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id"`
	Details     string    `json:"details"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

// AdminSession for tracking active sessions
type AdminSession struct {
	ID        string    `json:"id"`
	AdminID   string    `json:"admin_id"`
	Token     string    `json:"token"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAdminUser creates a new admin user
func NewAdminUser(email, firstName, lastName, roleID, createdBy string) *AdminUser {
	now := time.Now()
	return &AdminUser{
		ID:        uuid.New().String(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RoleID:    roleID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: &createdBy,
	}
}

// AdminLoginRequest for admin login
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// AdminLoginResponse after successful login
type AdminLoginResponse struct {
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	Admin     *AdminUser `json:"admin"`
}

// CreateAdminRequest for creating new admin
type CreateAdminRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleID    string `json:"role_id" binding:"required"`
}

// UpdateAdminRequest for updating admin
type UpdateAdminRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	RoleID    *string `json:"role_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
