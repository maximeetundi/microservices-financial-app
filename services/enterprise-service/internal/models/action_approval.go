package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActionApprovalStatus represents the state of an approval request
type ActionApprovalStatus string

const (
	ApprovalStatusPending  ActionApprovalStatus = "PENDING"
	ApprovalStatusApproved ActionApprovalStatus = "APPROVED"
	ApprovalStatusRejected ActionApprovalStatus = "REJECTED"
	ApprovalStatusExpired  ActionApprovalStatus = "EXPIRED"
	ApprovalStatusExecuted ActionApprovalStatus = "EXECUTED"
)

// ActionApproval represents a pending action that requires multi-admin approval
type ActionApproval struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	
	// Action Details
	ActionType   string                 `bson:"action_type" json:"action_type"` // PAYROLL, TRANSACTION, etc.
	ActionName   string                 `bson:"action_name" json:"action_name"` // Human readable
	Description  string                 `bson:"description" json:"description"`
	Payload      map[string]interface{} `bson:"payload" json:"payload"`         // Action data to execute
	
	// Amount (if applicable)
	Amount   float64 `bson:"amount,omitempty" json:"amount,omitempty"`
	Currency string  `bson:"currency,omitempty" json:"currency,omitempty"`
	
	// Approval Requirements
	RequiredApprovals int  `bson:"required_approvals" json:"required_approvals"`
	RequireMajority   bool `bson:"require_majority" json:"require_majority"`
	RequireAllAdmins  bool `bson:"require_all_admins" json:"require_all_admins"`
	TotalAdmins       int  `bson:"total_admins" json:"total_admins"` // For majority calculation
	
	// Current State
	Status   ActionApprovalStatus `bson:"status" json:"status"`
	Approvals []AdminApproval     `bson:"approvals" json:"approvals"`
	
	// Initiator
	InitiatedBy     string    `bson:"initiated_by" json:"initiated_by"`         // Employee ID
	InitiatorUserID string    `bson:"initiator_user_id" json:"initiator_user_id"` // User ID
	InitiatorName   string    `bson:"initiator_name" json:"initiator_name"`
	
	// Timestamps
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	ExpiresAt  time.Time  `bson:"expires_at" json:"expires_at"`
	ExecutedAt *time.Time `bson:"executed_at,omitempty" json:"executed_at,omitempty"`
}

// AdminApproval represents a single admin's approval or rejection
type AdminApproval struct {
	AdminEmployeeID string    `bson:"admin_employee_id" json:"admin_employee_id"`
	AdminUserID     string    `bson:"admin_user_id" json:"admin_user_id"`
	AdminName       string    `bson:"admin_name" json:"admin_name"`
	Decision        string    `bson:"decision" json:"decision"` // APPROVED, REJECTED
	Reason          string    `bson:"reason,omitempty" json:"reason,omitempty"` // For rejections
	DecidedAt       time.Time `bson:"decided_at" json:"decided_at"`
}

// GetApprovalCount returns the number of approvals
func (a *ActionApproval) GetApprovalCount() int {
	count := 0
	for _, approval := range a.Approvals {
		if approval.Decision == "APPROVED" {
			count++
		}
	}
	return count
}

// GetRejectionCount returns the number of rejections
func (a *ActionApproval) GetRejectionCount() int {
	count := 0
	for _, approval := range a.Approvals {
		if approval.Decision == "REJECTED" {
			count++
		}
	}
	return count
}

// IsApproved checks if the action has received enough approvals
func (a *ActionApproval) IsApproved() bool {
	approvals := a.GetApprovalCount()
	
	if a.RequireAllAdmins {
		return approvals >= a.TotalAdmins
	}
	
	if a.RequireMajority {
		return approvals > a.TotalAdmins/2
	}
	
	return approvals >= a.RequiredApprovals
}

// IsRejected checks if the action has been rejected by majority
func (a *ActionApproval) IsRejected() bool {
	rejections := a.GetRejectionCount()
	
	if a.RequireAllAdmins {
		return rejections >= 1 // Any rejection blocks if all required
	}
	
	// Rejected if majority rejected or cannot reach approval threshold
	remaining := a.TotalAdmins - len(a.Approvals)
	possibleApprovals := a.GetApprovalCount() + remaining
	
	return possibleApprovals < a.RequiredApprovals
}

// HasAdminVoted checks if a specific admin has already voted
func (a *ActionApproval) HasAdminVoted(adminUserID string) bool {
	for _, approval := range a.Approvals {
		if approval.AdminUserID == adminUserID {
			return true
		}
	}
	return false
}

// IsExpired checks if the approval request has expired
func (a *ActionApproval) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}
