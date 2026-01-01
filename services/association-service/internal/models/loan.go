package models

import "time"

// LoanStatus defines the status of a loan
type LoanStatus string

const (
	LoanStatusPending   LoanStatus = "pending"
	LoanStatusApproved  LoanStatus = "approved"
	LoanStatusActive    LoanStatus = "active"
	LoanStatusPaid      LoanStatus = "paid"
	LoanStatusDefaulted LoanStatus = "defaulted"
	LoanStatusRejected  LoanStatus = "rejected"
)

// Loan represents a loan from the association to a member
type Loan struct {
	ID            string     `json:"id"`
	AssociationID string     `json:"association_id"`
	BorrowerID    string     `json:"borrower_id"`
	Amount        float64    `json:"amount"`
	InterestRate  float64    `json:"interest_rate"`
	Duration      int        `json:"duration"` // in months
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	Repayments    JSONB      `json:"repayments"` // array of repayment records
	Status        LoanStatus `json:"status"`
	ApprovedBy    string     `json:"approved_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Loan Request is the request to request a loan
type LoanRequest struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" binding:"gte=0"`
	Duration     int     `json:"duration" binding:"required,gt=0"` // months
	Reason       string  `json:"reason"`
}

// ApproveLoanRequest is the request to approve/reject a loan
type ApproveLoanRequest struct {
	Approve bool   `json:"approve"`
	Reason  string `json:"reason"`
}

// RepaymentRequest is the request to record a loan repayment
type RepaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
