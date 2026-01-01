package models

import "time"

// ContributionStatus defines the status of a contribution
type ContributionStatus string

const (
	ContributionStatusPending ContributionStatus = "pending"
	ContributionStatusPaid    ContributionStatus = "paid"
	ContributionStatusLate    ContributionStatus = "late"
)

// Contribution represents a member's contribution
type Contribution struct {
	ID            string             `json:"id"`
	AssociationID string             `json:"association_id"`
	MemberID      string             `json:"member_id"`
	Amount        float64            `json:"amount"`
	Period        string             `json:"period"` // e.g., "2024-01" or "week-1"
	DueDate       time.Time          `json:"due_date"`
	PaidDate      *time.Time         `json:"paid_date,omitempty"`
	Status        ContributionStatus `json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// RecordContributionRequest is the request to record a contribution
type RecordContributionRequest struct {
	MemberID string  `json:"member_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Period   string  `json:"period" binding:"required"`
}
