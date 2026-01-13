package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayrollStatus string

const (
	PayrollStatusDraft      PayrollStatus = "DRAFT"
	PayrollStatusProcessing PayrollStatus = "PROCESSING"
	PayrollStatusCompleted  PayrollStatus = "COMPLETED"
	PayrollStatusFailed     PayrollStatus = "FAILED"
)

// PayrollRun represents a monthly salary execution for an Enterprise
// Section 5: "Paiement des salaires" & Section 6 "Gestion financière globale"
type PayrollRun struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EnterpriseID primitive.ObjectID `bson:"enterprise_id" json:"enterprise_id"`
	
	PeriodMonth int `bson:"period_month" json:"period_month"` // 1-12
	PeriodYear  int `bson:"period_year" json:"period_year"`   // 2024
	
	TotalAmount    float64 `bson:"total_amount" json:"total_amount"`
	TotalEmployees int     `bson:"total_employees" json:"total_employees"`
	
	Status        PayrollStatus `bson:"status" json:"status"`
	TransactionID string        `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"` // Link to Transfer Service Batch ID
	
	Details []PayrollDetail `bson:"details" json:"details"`

	ExecutedAt time.Time `bson:"executed_at" json:"executed_at"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type PayrollDetail struct {
	EmployeeID   primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	EmployeeName string             `bson:"employee_name" json:"employee_name"`
	BaseSalary   float64            `bson:"base_salary" json:"base_salary"`
	Bonuses      float64            `bson:"bonuses" json:"bonuses"`
	Deductions   float64            `bson:"deductions" json:"deductions"`
	NetPay       float64            `bson:"net_pay" json:"net_pay"`
	Status       string             `bson:"status" json:"status"` // SUCCESS, FAILED
}
