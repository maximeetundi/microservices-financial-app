package models

import "time"

// MemberRole defines the role of a member
type MemberRole string

const (
	MemberRolePresident  MemberRole = "president"
	MemberRoleTreasurer  MemberRole = "treasurer"
	MemberRoleSecretary  MemberRole = "secretary"
	MemberRoleMember     MemberRole = "member"
)

// MemberStatus defines the status of a member
type MemberStatus string

const (
	MemberStatusPending   MemberStatus = "pending"
	MemberStatusActive    MemberStatus = "active"
	MemberStatusSuspended MemberStatus = "suspended"
	MemberStatusLeft      MemberStatus = "left"
)

// Member represents a member of an association
type Member struct {
	ID                string       `json:"id"`
	AssociationID     string       `json:"association_id"`
	UserID            string       `json:"user_id"`
	Role              MemberRole   `json:"role"`
	JoinDate          time.Time    `json:"join_date"`
	Status            MemberStatus `json:"status"`
	ContributionsPaid float64      `json:"contributions_paid"`
	LoansReceived     float64      `json:"loans_received"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	
	// Populated from user service
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

// Note: JoinAssociationRequest and UpdateMemberRoleRequest are defined in association.go
