package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DonationStatus string

const (
	DonationStatusPending DonationStatus = "pending"
	DonationStatusPaid    DonationStatus = "paid"
	DonationStatusFailed  DonationStatus = "failed"
	DonationStatusRefunded DonationStatus = "refunded"
)

type DonationFrequency string

const (
	FrequencyOneTime  DonationFrequency = "one_time"
	FrequencyMonthly  DonationFrequency = "monthly"
	FrequencyQuarterly DonationFrequency = "quarterly"
	FrequencyAnnually DonationFrequency = "annually"
)

type Donation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CampaignID      primitive.ObjectID `bson:"campaign_id" json:"campaign_id"`
	DonorID         string             `bson:"donor_id,omitempty" json:"donor_id,omitempty"` // Empty if guest (future support) or purely anonymous in DB? No, usually keep ID but flag anonymous.
	
	Amount          float64           `bson:"amount" json:"amount"`
	Currency        string            `bson:"currency" json:"currency"`
	Message         string            `bson:"message" json:"message"`
	IsAnonymous     bool              `bson:"is_anonymous" json:"is_anonymous"` // Hides identity from public lists
	
	Frequency       DonationFrequency `bson:"frequency" json:"frequency"`
	NextPaymentDate *time.Time        `bson:"next_payment_date,omitempty" json:"next_payment_date,omitempty"`

	Status          DonationStatus    `bson:"status" json:"status"`
	TransactionID   string            `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	PaymentWalletID string            `bson:"payment_wallet_id,omitempty" json:"payment_wallet_id,omitempty"`
	
	FormData        map[string]interface{} `bson:"form_data,omitempty" json:"form_data,omitempty"`
	
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}
