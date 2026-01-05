package models

import "time"

// Permission types for custom roles
type Permission string

const (
	PermissionManageMembers     Permission = "manage_members"
	PermissionManageMeetings    Permission = "manage_meetings"
	PermissionApproveLoan       Permission = "approve_loans"
	PermissionApproveDistribution Permission = "approve_distributions"
	PermissionViewTreasury      Permission = "view_treasury"
	PermissionRecordContributions Permission = "record_contributions"
	PermissionManageRoles       Permission = "manage_roles"
	PermissionManageChat        Permission = "manage_chat"
	PermissionManageSolidarity  Permission = "manage_solidarity"
)

// AssociationRole represents a custom role with permissions
type AssociationRole struct {
	ID            string       `json:"id"`
	AssociationID string       `json:"association_id"`
	Name          string       `json:"name"`
	Permissions   []Permission `json:"permissions"`
	IsDefault     bool         `json:"is_default"`
	CreatedAt     time.Time    `json:"created_at"`
}

// CreateRoleRequest is the request to create a custom role
type CreateRoleRequest struct {
	Name        string       `json:"name" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
}

// Approver represents a designated approver for multi-signature
type Approver struct {
	ID            string    `json:"id"`
	AssociationID string    `json:"association_id"`
	MemberID      string    `json:"member_id"`
	MemberName    string    `json:"member_name,omitempty"` // Populated on read
	Position      int       `json:"position"`
	CreatedAt     time.Time `json:"created_at"`
}

// SetApproversRequest is the request to set approvers
type SetApproversRequest struct {
	MemberIDs []string `json:"member_ids" binding:"required,min=5,max=5"` // Exactly 5 members
}

// ApprovalRequestStatus defines the status of an approval request
type ApprovalRequestStatus string

const (
	ApprovalStatusPending  ApprovalRequestStatus = "pending"
	ApprovalStatusApproved ApprovalRequestStatus = "approved"
	ApprovalStatusRejected ApprovalRequestStatus = "rejected"
	ApprovalStatusExpired  ApprovalRequestStatus = "expired"
)

// ApprovalRequestType defines what kind of request needs approval
type ApprovalRequestType string

const (
	ApprovalTypeLoan         ApprovalRequestType = "loan"
	ApprovalTypeDistribution ApprovalRequestType = "distribution"
)

// ApprovalRequest represents a request that needs multi-signature approval
type ApprovalRequest struct {
	ID                string                `json:"id"`
	AssociationID     string                `json:"association_id"`
	RequestType       ApprovalRequestType   `json:"request_type"`
	ReferenceID       string                `json:"reference_id"`
	Amount            float64               `json:"amount"`
	RequesterID       string                `json:"requester_id"`
	RequesterName     string                `json:"requester_name,omitempty"`
	Status            ApprovalRequestStatus `json:"status"`
	RequiredApprovals int                   `json:"required_approvals"`
	CurrentApprovals  int                   `json:"current_approvals,omitempty"`
	CurrentRejections int                   `json:"current_rejections,omitempty"`
	Description       string                `json:"description"`
	Votes             []ApprovalVote        `json:"votes,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

// VoteType defines the type of vote
type VoteType string

const (
	VoteApprove VoteType = "approve"
	VoteReject  VoteType = "reject"
)

// ApprovalVote represents a vote on an approval request
type ApprovalVote struct {
	ID          string    `json:"id"`
	RequestID   string    `json:"request_id"`
	ApproverID  string    `json:"approver_id"`
	ApproverName string   `json:"approver_name,omitempty"`
	Vote        VoteType  `json:"vote"`
	Comment     string    `json:"comment"`
	VotedAt     time.Time `json:"voted_at"`
}

// VoteRequest is the request to vote on an approval
type VoteRequest struct {
	Vote    VoteType `json:"vote" binding:"required"`
	Comment string   `json:"comment"`
}
