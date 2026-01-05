package models

import "time"

// ChatMessage represents a message in the association chat
type ChatMessage struct {
	ID            string    `json:"id"`
	AssociationID string    `json:"association_id"`
	SenderID      string    `json:"sender_id"`
	SenderName    string    `json:"sender_name,omitempty"`
	Content       string    `json:"content"`
	IsAdminOnly   bool      `json:"is_admin_only"`
	CreatedAt     time.Time `json:"created_at"`
}

// SendMessageRequest is the request to send a chat message
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	IsAdminOnly bool   `json:"is_admin_only"`
}

// SolidarityEventType defines the type of solidarity event
type SolidarityEventType string

const (
	SolidarityDeceased  SolidarityEventType = "deceased"  // Deuil
	SolidarityMarriage  SolidarityEventType = "marriage"  // Mariage
	SolidarityBirth     SolidarityEventType = "birth"     // Naissance
	SolidarityIllness   SolidarityEventType = "illness"   // Maladie
	SolidarityOther     SolidarityEventType = "other"     // Autre
)

// SolidarityEventStatus defines the status
type SolidarityEventStatus string

const (
	SolidarityStatusActive SolidarityEventStatus = "active"
	SolidarityStatusClosed SolidarityEventStatus = "closed"
)

// SolidarityEvent represents a solidarity/aid event
type SolidarityEvent struct {
	ID              string                `json:"id"`
	AssociationID   string                `json:"association_id"`
	EventType       SolidarityEventType   `json:"event_type"`
	BeneficiaryID   string                `json:"beneficiary_id"`
	BeneficiaryName string                `json:"beneficiary_name,omitempty"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	TargetAmount    float64               `json:"target_amount"`
	CollectedAmount float64               `json:"collected_amount"`
	Status          SolidarityEventStatus `json:"status"`
	CreatedBy       string                `json:"created_by"`
	CreatedAt       time.Time             `json:"created_at"`
	ClosedAt        *time.Time            `json:"closed_at,omitempty"`
	Contributions   []SolidarityContribution `json:"contributions,omitempty"`
}

// CreateSolidarityEventRequest is the request to create a solidarity event
type CreateSolidarityEventRequest struct {
	EventType     SolidarityEventType `json:"event_type" binding:"required"`
	BeneficiaryID string              `json:"beneficiary_id" binding:"required"`
	Title         string              `json:"title" binding:"required"`
	Description   string              `json:"description"`
	TargetAmount  float64             `json:"target_amount"`
}

// SolidarityContribution represents a contribution to a solidarity event
type SolidarityContribution struct {
	ID              string     `json:"id"`
	EventID         string     `json:"event_id"`
	ContributorID   string     `json:"contributor_id"`
	ContributorName string     `json:"contributor_name,omitempty"`
	Amount          float64    `json:"amount"`
	Paid            bool       `json:"paid"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ContributeSolidarityRequest is the request to contribute to a solidarity event
type ContributeSolidarityRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	WalletID string  `json:"wallet_id" binding:"required"`
	Pin      string  `json:"pin" binding:"required"`
}

// CalledRound represents a round in a "called" tontine
type CalledRound struct {
	ID             string     `json:"id"`
	AssociationID  string     `json:"association_id"`
	BeneficiaryID  string     `json:"beneficiary_id"`
	BeneficiaryName string    `json:"beneficiary_name,omitempty"`
	RoundNumber    int        `json:"round_number"`
	TotalCollected float64    `json:"total_collected"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	Pledges        []CalledPledge `json:"pledges,omitempty"`
}

// CreateCalledRoundRequest is the request to start a new called round
type CreateCalledRoundRequest struct {
	BeneficiaryID string `json:"beneficiary_id" binding:"required"`
}

// CalledPledge represents a pledge in a called round
type CalledPledge struct {
	ID              string     `json:"id"`
	RoundID         string     `json:"round_id"`
	ContributorID   string     `json:"contributor_id"`
	ContributorName string     `json:"contributor_name,omitempty"`
	Amount          float64    `json:"amount"`
	Paid            bool       `json:"paid"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
}

// MakePledgeRequest is the request to make a pledge
type MakePledgeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// PayPledgeRequest is the request to pay a pledge
type PayPledgeRequest struct {
	WalletID string `json:"wallet_id" binding:"required"`
	Pin      string `json:"pin" binding:"required"`
}

// AssociationSettings represents configurable settings
type AssociationSettings struct {
	LateFeeAmount   float64 `json:"late_fee_amount"`
	LateFeePercent  float64 `json:"late_fee_percent"`
	FoodFee         float64 `json:"food_fee"`
	DrinkFee        float64 `json:"drink_fee"`
	RequiredApprovals int   `json:"required_approvals"`
}
