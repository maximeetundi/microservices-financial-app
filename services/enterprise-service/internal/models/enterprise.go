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

// Enterprise corresponds to Section 1: "Création et gestion d’une entreprise"
type Enterprise struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      EnterpriseType     `bson:"type" json:"type"`
	OwnerID   string             `bson:"owner_id" json:"owner_id"` // Link to User ID
	Settings  EnterpriseSettings `bson:"settings" json:"settings"`
	
	// Dynamic Service Definitions (Section 7)
	CustomServices []ServiceDefinition `bson:"custom_services" json:"custom_services"`

	// Specialized Modules Config (Section 7a, 7b)
	TransportConfig *TransportConfig `bson:"transport_config,omitempty" json:"transport_config,omitempty"`
	SchoolConfig    *SchoolConfig    `bson:"school_config,omitempty" json:"school_config,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type EnterpriseSettings struct {
	Currency        string `bson:"currency" json:"currency"`
	PayrollDate     int    `bson:"payroll_date" json:"payroll_date"` // e.g., 25th of month
	AutoPaySalaries bool   `bson:"auto_pay_salaries" json:"auto_pay_salaries"`
}

// ServiceDefinition allows defining "Eau", "Electricité", "Gardiennage" dynamically.
type ServiceDefinition struct {
	ID        string  `bson:"id" json:"id"` // slug, e.g., "gym_membership"
	Name      string  `bson:"name" json:"name"`
	Unit      string  `bson:"unit" json:"unit"` // e.g., "Month", "Session"
	BasePrice float64 `bson:"base_price" json:"base_price"`
	
	// Dynamic Frequency Config
	BillingFrequency string `bson:"billing_frequency" json:"billing_frequency"` // DAILY, WEEKLY, MONTHLY, QUARTERLY, ANNUALLY, CUSTOM
	CustomInterval   int    `bson:"custom_interval,omitempty" json:"custom_interval,omitempty"` // If CUSTOM, e.g., every 15 days
}

// TransportConfig (Point 7a: Lignes, Zones, Tarifs)
type TransportConfig struct {
	Routes []Route `bson:"routes" json:"routes"`
	Zones  []Zone  `bson:"zones" json:"zones"`
}

type Route struct {
	ID          string  `bson:"id" json:"id"`
	Name        string  `bson:"name" json:"name"` // "Ligne 14"
	BasePrice   float64 `bson:"base_price" json:"base_price"`
}

type Zone struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"` // "Zone A - Centre"
}

// SchoolConfig (Point 7b: Classes, Tranches)
type SchoolConfig struct {
	Classes []Class `bson:"classes" json:"classes"`
}

type Class struct {
	ID          string    `bson:"id" json:"id"` // "CP", "CE1"
	Name        string    `bson:"name" json:"name"`
	TotalFees   float64   `bson:"total_fees" json:"total_fees"`
	Tranches    []Tranche `bson:"tranches" json:"tranches"`
}

type Tranche struct {
	ID        string    `bson:"id" json:"id"` // "T1", "T2"
	Name      string    `bson:"name" json:"name"` // "1ère Tranche"
	Amount    float64   `bson:"amount" json:"amount"`
	DueDate   time.Time `bson:"due_date" json:"due_date"`
}
