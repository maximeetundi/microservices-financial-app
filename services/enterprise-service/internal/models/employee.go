package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmployeeStatus string

const (
	EmployeeStatusPending EmployeeStatus = "PENDING_INVITE" // 2. "Tant que l’invitation n’est pas acceptée..."
	EmployeeStatusActive  EmployeeStatus = "ACTIVE"
	EmployeeStatusTerminated EmployeeStatus = "TERMINATED"
)

// Employee corresponds to Section 2: "Gestion avancée des employés"
type Employee struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	UserID       string             `bson:"user_id,omitempty" json:"user_id,omitempty"` // Link to existing User account
	
	// Profile
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Profession  string `bson:"profession" json:"profession"`
	Email       string `bson:"email" json:"email"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`

	// Section 2: Consentement & Sécurité
	Status     EmployeeStatus `bson:"status" json:"status"`
	PinHash    string         `bson:"-" json:"-"` // Not stored here, verified against Auth service or local hash depending on arch
	InvitedAt  time.Time      `bson:"invited_at" json:"invited_at"`
	AcceptedAt time.Time      `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`

	// Section 3: Configuration complète du salaire
	SalaryConfig SalaryConfig `bson:"salary_config" json:"salary_config"`
	
	// Section 4: Historique & Évolution
	History []CareerEvent `bson:"history" json:"history"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type SalaryConfig struct {
	BaseAmount  float64       `bson:"base_amount" json:"base_amount"`
	Frequency   string        `bson:"frequency" json:"frequency"` // MONTHLY, WEEKLY
	Deductions  []FinancialItem `bson:"deductions" json:"deductions"` // Impôts, Cotisations
	Bonuses     []FinancialItem `bson:"bonuses" json:"bonuses"`       // Primes, Avantages
	NetPayable  float64       `bson:"net_payable" json:"net_payable"` // Calcul automatique
}

type FinancialItem struct {
	Name   string  `bson:"name" json:"name"`
	Type   string  `bson:"type" json:"type"` // PERCENTAGE or FIXED
	Value  float64 `bson:"value" json:"value"`
	Amount float64 `bson:"amount" json:"amount"` // Calculated absolute amount
}

// CareerEvent tracks "Promouvoir", "Rétrograder", "Licencier"
type CareerEvent struct {
	Date        time.Time `bson:"date" json:"date"`
	Type        string    `bson:"type" json:"type"` // PROMOTION, DEMOTION, SALARY_CHANGE, TERMINATION
	Description string    `bson:"description" json:"description"`
	Previous    interface{} `bson:"previous,omitempty" json:"previous,omitempty"`
	New         interface{} `bson:"new,omitempty" json:"new,omitempty"`
}
