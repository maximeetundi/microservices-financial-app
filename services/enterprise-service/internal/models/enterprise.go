package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnterpriseType string

const (
	EnterpriseTypeSchool    EnterpriseType = "SCHOOL"
	EnterpriseTypeTransport EnterpriseType = "TRANSPORT"
	EnterpriseTypeService   EnterpriseType = "SERVICE" // Generic
	EnterpriseTypeUtility   EnterpriseType = "UTILITY" // Electricity, Water
)

// Enterprise corresponds to Section 1: "Création et gestion d'une entreprise"
type Enterprise struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Type        EnterpriseType     `bson:"type" json:"type"` // Default: SERVICE
	Status      string             `bson:"status" json:"status"` // ACTIVE, SUSPENDED, PENDING
	Logo        string             `bson:"logo,omitempty" json:"logo,omitempty"`
	EmployeeCountRange string    `bson:"employee_count_range" json:"employee_count_range"` // "1-10", "11-50", etc.
	RegistrationNumber string    `bson:"registration_number,omitempty" json:"registration_number,omitempty"`
	
	OwnerID   string             `bson:"owner_id" json:"owner_id"` // Link to User ID (auto-admin)
	Settings  EnterpriseSettings `bson:"settings" json:"settings"`
	
	// === Enterprise Wallets ===
	DefaultWalletID string   `bson:"default_wallet_id" json:"default_wallet_id"` // Auto-created
	WalletIDs       []string `bson:"wallet_ids" json:"wallet_ids"`               // Additional wallets
	
	// === Security Policies (Dynamic Multi-Approval) ===
	SecurityPolicies []SecurityPolicy `bson:"security_policies" json:"security_policies"`
	
	// === Job Positions (Postes avec salaires par défaut) ===
	JobPositions []JobPosition `bson:"job_positions" json:"job_positions"`
	
	// Dynamic Service Definitions (Section 7)
	ServiceGroups []ServiceGroup `bson:"service_groups" json:"service_groups"`

	// Specialized Modules Config (Section 7a, 7b)
	TransportConfig *TransportConfig `bson:"transport_config,omitempty" json:"transport_config,omitempty"`
	SchoolConfig    *SchoolConfig    `bson:"school_config,omitempty" json:"school_config,omitempty"`

	// Data Export Tracking (required before deletion)
	LastExportedAt time.Time `bson:"last_exported_at,omitempty" json:"last_exported_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// JobPosition defines a job/role with default salary in the enterprise
type JobPosition struct {
	ID            string  `bson:"id" json:"id"`
	Name          string  `bson:"name" json:"name"`                     // Ex: "Enseignant", "Comptable"
	Description   string  `bson:"description" json:"description"`
	Department    string  `bson:"department" json:"department"`         // "Finance", "Pédagogie", "IT"
	DefaultSalary float64 `bson:"default_salary" json:"default_salary"` // Salaire de base
	Currency      string  `bson:"currency" json:"currency"`             // XOF, EUR, etc.
}

// SecurityPolicy defines approval rules for specific actions
type SecurityPolicy struct {
	ID              string   `bson:"id" json:"id"`
	Name            string   `bson:"name" json:"name"`              // "Paiement de salaires", "Transfert important"
	ActionType      string   `bson:"action_type" json:"action_type"` // See ActionType constants
	Enabled         bool     `bson:"enabled" json:"enabled"`
	
	// Approval Rules
	MinApprovals    int      `bson:"min_approvals" json:"min_approvals"`       // Nombre minimum d'admins
	RequireMajority bool     `bson:"require_majority" json:"require_majority"` // Ou majorité des admins
	RequireAllAdmins bool    `bson:"require_all_admins" json:"require_all_admins"` // Tous les admins
	
	// Conditions (optionnelles)
	ThresholdAmount float64  `bson:"threshold_amount,omitempty" json:"threshold_amount,omitempty"` // Montant minimum pour appliquer
	
	// Timeout
	ExpirationHours int      `bson:"expiration_hours" json:"expiration_hours"` // Délai avant expiration (defaut: 24h)
}

// Action types that can require multi-approval
const (
	ActionTypeTransaction        = "TRANSACTION"         // Transferts, retraits
	ActionTypePayroll            = "PAYROLL"             // Paiement des salaires
	ActionTypeEmployeeTerminate  = "EMPLOYEE_TERMINATE"  // Licenciement
	ActionTypeEmployeePromote    = "EMPLOYEE_PROMOTE"    // Promotion
	ActionTypeSettingsChange     = "SETTINGS_CHANGE"     // Modification des paramètres
	ActionTypeWalletCreate       = "WALLET_CREATE"       // Création de wallet
	ActionTypeAdminChange        = "ADMIN_CHANGE"        // Ajout/retrait d'admin
	ActionTypeServiceCreate      = "SERVICE_CREATE"      // Création de service
	ActionTypeInvoiceBatch       = "INVOICE_BATCH"       // Envoi de factures en lot
	ActionTypeEnterpriseDelete   = "ENTERPRISE_DELETE"   // Suppression d'entreprise (export requis)
)

type ServiceGroup struct {
	ID        string              `bson:"id" json:"id"`
	Name      string              `bson:"name" json:"name"` // e.g. "Pension 6ème"
	IsPrivate bool                `bson:"is_private" json:"is_private"` // If true, hidden from QR scan
	Currency  string              `bson:"currency" json:"currency"`     // Specific currency for this group
	WalletID  string              `bson:"wallet_id,omitempty" json:"wallet_id,omitempty"` // Destination wallet for payments
	Services  []ServiceDefinition `bson:"services" json:"services"`
}

type EnterpriseSettings struct {
	PayrollDate     int    `bson:"payroll_date" json:"payroll_date"` // e.g., 25th of month
	AutoPaySalaries bool   `bson:"auto_pay_salaries" json:"auto_pay_salaries"`
	DefaultCurrency string `bson:"default_currency" json:"default_currency"` // XOF, EUR, etc.
}

// ServiceDefinition allows defining "Eau", "Electricité", "Gardiennage" dynamically.
type ServiceDefinition struct {
	ID        string  `bson:"id" json:"id"` // slug, e.g., "gym_membership"
	Name      string  `bson:"name" json:"name"`
	Unit      string  `bson:"unit" json:"unit"` // e.g., "kWh", "m³", "Session"
	BasePrice float64 `bson:"base_price" json:"base_price"`
	
	// Dynamic Frequency Config
	BillingType      string  `bson:"billing_type" json:"billing_type"` // FIXED, USAGE
	BillingFrequency string  `bson:"billing_frequency" json:"billing_frequency"` // DAILY, WEEKLY, MONTHLY, QUARTERLY, ANNUALLY, CUSTOM
	CustomInterval   int     `bson:"custom_interval,omitempty" json:"custom_interval,omitempty"` // If CUSTOM, e.g., every 15 days
    
	// Dynamic Pricing Config (for USAGE billing type)
	PricingMode    string         `bson:"pricing_mode,omitempty" json:"pricing_mode,omitempty"` // FIXED, TIERED, THRESHOLD
	PricingTiers   []PricingTier  `bson:"pricing_tiers,omitempty" json:"pricing_tiers,omitempty"` // For TIERED/THRESHOLD mode
    
    PaymentSchedule  []PaymentScheduleItem `bson:"payment_schedule,omitempty" json:"payment_schedule,omitempty"`
    FormSchema       []ReqFormField        `bson:"form_schema,omitempty" json:"form_schema,omitempty"`
    PenaltyConfig    *PenaltyConfig        `bson:"penalty_config,omitempty" json:"penalty_config,omitempty"`
}

// PricingTier defines a consumption tier with optional bonuses
// Used for tiered pricing (electricity, water, etc.)
type PricingTier struct {
	MinConsumption float64 `bson:"min_consumption" json:"min_consumption"` // e.g., 0, 100, 500
	MaxConsumption float64 `bson:"max_consumption" json:"max_consumption"` // e.g., 100, 500, -1 (unlimited)
	PricePerUnit   float64 `bson:"price_per_unit" json:"price_per_unit"`   // Base price for units in this tier
	FixedBonus     float64 `bson:"fixed_bonus,omitempty" json:"fixed_bonus,omitempty"`     // Fixed amount added when consumption enters this tier
	PercentBonus   float64 `bson:"percent_bonus,omitempty" json:"percent_bonus,omitempty"` // Percentage of tier total added as bonus
	Label          string  `bson:"label,omitempty" json:"label,omitempty"` // e.g., "Palier 1", "Tranche sociale"
}


type PenaltyConfig struct {
    Type        string  `bson:"type" json:"type"`                                   // FIXED, PERCENTAGE, HYBRID
    Value       float64 `bson:"value" json:"value"`                                 // Fixed amount (for FIXED and HYBRID)
    Percentage  float64 `bson:"percentage,omitempty" json:"percentage,omitempty"`   // Percentage (for PERCENTAGE and HYBRID)
    Frequency   string  `bson:"frequency" json:"frequency"`                         // DAILY, WEEKLY, MONTHLY, QUARTERLY, SEMIANNUAL, ANNUAL
    GracePeriod int     `bson:"grace_period" json:"grace_period"`                   // Grace period value
    GraceUnit   string  `bson:"grace_unit,omitempty" json:"grace_unit,omitempty"`   // DAYS, WEEKS, MONTHS (default: DAYS)
    MaxPenaltyMonths int `bson:"max_penalty_months,omitempty" json:"max_penalty_months,omitempty"` // Max months to apply penalty (default: 6, 0 = forever)
    MaxPenaltyAmount float64 `bson:"max_penalty_amount,omitempty" json:"max_penalty_amount,omitempty"` // Max penalty amount (0 = no limit)
}

type ReqFormField struct {
    Key      string   `bson:"key" json:"key"`
    Label    string   `bson:"label" json:"label"`
    Type     string   `bson:"type" json:"type"`
    Required bool     `bson:"required" json:"required"`
    Options  []string `bson:"options,omitempty" json:"options,omitempty"`
}

type PaymentScheduleItem struct {
    Name      string  `bson:"name" json:"name"`
    StartDate string  `bson:"start_date" json:"start_date"` // Format: YYYY-MM-DD
    EndDate   string  `bson:"end_date" json:"end_date"`     // Format: YYYY-MM-DD
    Amount    float64 `bson:"amount" json:"amount"`
}

// TransportConfig (Point 7a: Lignes, Zones, Tarifs)
type TransportConfig struct {
	Routes []Route `bson:"routes" json:"routes"`
	Zones  []Zone  `bson:"zones" json:"zones"`
}

type Route struct {
	ID          string  `bson:"id" json:"id"`
	Name        string  `bson:"name" json:"name"`
	BasePrice   float64 `bson:"base_price" json:"base_price"`
}

type Zone struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

// SchoolConfig (Point 7b: Classes, Tranches)
type SchoolConfig struct {
	Classes []Class `bson:"classes" json:"classes"`
}

type Class struct {
	ID          string    `bson:"id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	TotalFees   float64   `bson:"total_fees" json:"total_fees"`
	Tranches    []Tranche `bson:"tranches" json:"tranches"`
}

type Tranche struct {
	ID        string    `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Amount    float64   `bson:"amount" json:"amount"`
	DueDate   time.Time `bson:"due_date" json:"due_date"`
}

// GetDefaultSecurityPolicies returns sensible default policies for new enterprises
func GetDefaultSecurityPolicies() []SecurityPolicy {
	return []SecurityPolicy{
		{
			ID:              "policy_payroll",
			Name:            "Paiement des salaires",
			ActionType:      ActionTypePayroll,
			Enabled:         false, // Désactivé par défaut
			MinApprovals:    1,
			ExpirationHours: 24,
		},
		{
			ID:              "policy_large_transfer",
			Name:            "Transferts importants",
			ActionType:      ActionTypeTransaction,
			Enabled:         false,
			MinApprovals:    2,
			ThresholdAmount: 1000000, // 1M XOF
			ExpirationHours: 24,
		},
		{
			ID:              "policy_admin_change",
			Name:            "Modification des administrateurs",
			ActionType:      ActionTypeAdminChange,
			Enabled:         true, // Activé par défaut pour sécurité
			RequireAllAdmins: true,
			ExpirationHours: 48,
		},
	}
}

