package models

import "time"

// TransactionType defines the type of treasury transaction
type TransactionType string

const (
	TransactionTypeContribution TransactionType = "contribution"
	TransactionTypeLoan         TransactionType = "loan"
	TransactionTypeRepayment    TransactionType = "repayment"
	TransactionTypeDistribution TransactionType = "distribution"
	TransactionTypeExpense      TransactionType = "expense"
)

// TransactionStatus defines the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// TreasuryTransaction represents a financial transaction in the association
type TreasuryTransaction struct {
	ID            string            `json:"id"`
	AssociationID string            `json:"association_id"`
	Type          TransactionType   `json:"type"`
	Amount        float64           `json:"amount"`
	FromMemberID  string            `json:"from_member_id,omitempty"`
	ToMemberID    string            `json:"to_member_id,omitempty"`
	Description   string            `json:"description"`
	Status        TransactionStatus `json:"status"`
	CreatedBy     string            `json:"created_by"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// RecordTransactionRequest is the request to record a transaction
type RecordTransactionRequest struct {
	Type         TransactionType `json:"type" binding:"required"`
	Amount       float64         `json:"amount" binding:"required,gt=0"`
	FromMemberID string          `json:"from_member_id"`
	ToMemberID   string          `json:"to_member_id"`
	Description  string          `json:"description"`
}

// TreasuryReport represents a financial report
type TreasuryReport struct {
	TotalBalance       float64                `json:"total_balance"`
	TotalContributions float64                `json:"total_contributions"`
	TotalLoans         float64                `json:"total_loans"`
	TotalRepayments    float64                `json:"total_repayments"`
	TotalDistributions float64                `json:"total_distributions"`
	TotalExpenses      float64                `json:"total_expenses"`
	Transactions       []*TreasuryTransaction `json:"transactions"`
}

