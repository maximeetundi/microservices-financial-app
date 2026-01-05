package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// AssociationType defines the type of association
type AssociationType string

const (
	AssociationTypeTontine AssociationType = "tontine"
	AssociationTypeSavings AssociationType = "savings"
	AssociationTypeCredit  AssociationType = "credit"
	AssociationTypeGeneral AssociationType = "general"
)

// AssociationStatus defines the status of an association
type AssociationStatus string

const (
	AssociationStatusActive    AssociationStatus = "active"
	AssociationStatusSuspended AssociationStatus = "suspended"
	AssociationStatusClosed    AssociationStatus = "closed"
)

// Association represents a group/association
type Association struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            AssociationType   `json:"type"`
	Description     string            `json:"description"`
	Rules           JSONB             `json:"rules"`
	TotalMembers    int               `json:"total_members"`
	TreasuryBalance float64           `json:"treasury_balance"`
	Currency        string            `json:"currency"`
	Status          AssociationStatus `json:"status"`
	CreatedBy       string            `json:"created_by"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// CreateAssociationRequest is the request to create an association
type CreateAssociationRequest struct {
	Name        string          `json:"name" binding:"required"`
	Type        AssociationType `json:"type" binding:"required"`
	Description string          `json:"description"`
	Rules       JSONB           `json:"rules"`
	Currency    string          `json:"currency"`
	CreatorRole MemberRole      `json:"creator_role"` // Optional: role of the creator (default: president)
}

// UpdateAssociationRequest is the request to update an association
type UpdateAssociationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Rules       *JSONB  `json:"rules"`
	Status      *string `json:"status"`
}

// JoinAssociationRequest is the request to join an association
type JoinAssociationRequest struct {
	Message string `json:"message"`
}

// UpdateMemberRoleRequest is the request to update a member's role
type UpdateMemberRoleRequest struct {
	Role MemberRole `json:"role" binding:"required"`
}

// AssociationStats represents statistics for an association
type AssociationStats struct {
	TotalMembers       int     `json:"total_members"`
	ActiveMembers      int     `json:"active_members"`
	TreasuryBalance    float64 `json:"treasury_balance"`
	TotalContributions float64 `json:"total_contributions"`
	TotalLoans         float64 `json:"total_loans"`
	TotalRepayments    float64 `json:"total_repayments"`
	UpcomingMeetings   int     `json:"upcoming_meetings"`
}

// JSONB type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB value")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

