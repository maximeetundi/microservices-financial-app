package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DRAFT"
	InvoiceStatusSent      InvoiceStatus = "SENT"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusOverdue   InvoiceStatus = "OVERDUE"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

// Invoice represents a request for payment sent to a Client (Point 9)
type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID  primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	
	// Client Info (Can be linked to User or ad-hoc from Excel import)
	ClientID      string             `bson:"client_id,omitempty" json:"client_id,omitempty"` // Optional Link
	ClientName    string             `bson:"client_name" json:"client_name"`
	ClientContact string             `bson:"client_contact" json:"client_contact"` // Email or Phone

	// Service Info
	ServiceID     string             `bson:"service_id" json:"service_id"` // "utility_bill", "school_fee"
	ServiceName   string             `bson:"service_name" json:"service_name"`
	SubscriptionID primitive.ObjectID `bson:"subscription_id,omitempty" json:"subscription_id,omitempty"`
	
	Amount        float64            `bson:"amount" json:"amount"`
	Consumption   *float64           `bson:"consumption,omitempty" json:"consumption,omitempty"` // For Usage Billing
	Currency      string             `bson:"currency" json:"currency"`
	Description   string             `bson:"description" json:"description"`
	Type          string             `bson:"type" json:"type"` // e.g., "TUITION", "TRANSPORT"
	Metadata      map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	
	// Dates
	DueDate       time.Time          `bson:"due_date" json:"due_date"`
	SentAt        time.Time          `bson:"sent_at,omitempty" json:"sent_at,omitempty"`
	PaidAt        time.Time          `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	
	Status        InvoiceStatus      `bson:"status" json:"status"`
	PaymentRef    string             `bson:"payment_ref,omitempty" json:"payment_ref,omitempty"`

	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// Subscription represents a recurring contract (Point 7)
type Subscription struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID  primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	ClientID      string             `bson:"client_id" json:"client_id"` // Must be a registered User for Auto-Debit
	ClientName    string             `bson:"client_name" json:"client_name"` // Snapshot for display
	ExternalID    string             `bson:"external_id" json:"external_id"` // Matricule, Compteur, Ref
	
	ServiceID     string             `bson:"service_id" json:"service_id"` // e.g. "transport_monthly", "school_tranche"
	
	// Specialized Config (Polymorphic-ish)
	TransportDetails *TransportDetails `bson:"transport_details,omitempty" json:"transport_details,omitempty"`
	SchoolDetails    *SchoolDetails    `bson:"school_details,omitempty" json:"school_details,omitempty"`
	
	// Dynamic Form Data (Section 8: "Formulaire personnalisé")
	FormData     map[string]interface{} `bson:"form_data,omitempty" json:"form_data,omitempty"`
	
	Metadata     map[string]string   `bson:"metadata,omitempty" json:"metadata,omitempty"`
	
	Status        string             `bson:"status" json:"status"` // ACTIVE, CANCELLED
	
	// Billing Logic
	BillingFrequency string          `bson:"billing_frequency" json:"billing_frequency"` // Copied from ServiceDef or Overridden
	Amount        float64            `bson:"amount" json:"amount"`
	CustomInterval int               `bson:"custom_interval,omitempty" json:"custom_interval,omitempty"` // For CUSTOM frequency
	
	LastBillingAt *time.Time         `bson:"last_billing_at,omitempty" json:"last_billing_at,omitempty"`
	NextBillingAt time.Time          `bson:"next_billing_at" json:"next_billing_at"`
}

// InvoiceBatch (Point 11: Validation avant envoi)
type InvoiceBatch struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	Name         string             `bson:"name" json:"name"` // "Import Janvier"
	
	TotalInvoices int     `bson:"total_invoices" json:"total_invoices"`
	TotalAmount   float64 `bson:"total_amount" json:"total_amount"`
	
	Status        string    `bson:"status" json:"status"` // DRAFT, VALIDATED, SCHEDULED, PROCESSED
	ScheduledAt   time.Time `bson:"scheduled_at,omitempty" json:"scheduled_at,omitempty"`
	
	Invoices []Invoice `bson:"invoices" json:"invoices"` // Embedded for MVP simplicity
	
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// TransportDetails (Point 7a)
type TransportDetails struct {
	RouteID     string `bson:"route_id" json:"route_id"`
	ZoneID      string `bson:"zone_id" json:"zone_id"`
	PassType    string `bson:"pass_type" json:"pass_type"` // MONTHLY, QUARTERLY
}

// SchoolDetails (Point 7b)
type SchoolDetails struct {
	ClassID     string `bson:"class_id" json:"class_id"`
	StudentName string `bson:"student_name" json:"student_name"`
	TrancheID   string `bson:"tranche_id" json:"tranche_id"` // If paying by tranche
}
