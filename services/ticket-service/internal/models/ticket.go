package models

import "time"

// Ticket represents a purchased ticket
type Ticket struct {
	ID            string            `json:"id" db:"id"`
	EventID       string            `json:"event_id" db:"event_id"`
	BuyerID       string            `json:"buyer_id" db:"buyer_id"`
	TierID        string            `json:"tier_id" db:"tier_id"`
	TierName      string            `json:"tier_name" db:"tier_name"`
	TierIcon      string            `json:"tier_icon" db:"tier_icon"`
	Price         float64           `json:"price" db:"price"`
	Currency      string            `json:"currency" db:"currency"`
	FormData      map[string]string `json:"form_data"`
	QRCode        string            `json:"qr_code" db:"qr_code"`
	TicketCode    string            `json:"ticket_code" db:"ticket_code"`
	Status        string            `json:"status" db:"status"` // pending, paid, used, cancelled, refunded
	TransactionID string            `json:"transaction_id" db:"transaction_id"`
	UsedAt        *time.Time        `json:"used_at,omitempty" db:"used_at"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	
	// Join fields
	EventTitle    string            `json:"event_title,omitempty"`
	EventDate     *time.Time        `json:"event_date,omitempty"`
	EventLocation string            `json:"event_location,omitempty"`
	BuyerName     string            `json:"buyer_name,omitempty"`
	BuyerEmail    string            `json:"buyer_email,omitempty"`
}

// Ticket status constants
const (
	TicketStatusPending   = "pending"
	TicketStatusPaid      = "paid"
	TicketStatusUsed      = "used"
	TicketStatusCancelled = "cancelled"
	TicketStatusRefunded  = "refunded"
)

// PurchaseTicketRequest represents the request to purchase a ticket
type PurchaseTicketRequest struct {
	EventID      string            `json:"event_id" binding:"required"`
	TierID       string            `json:"tier_id" binding:"required"`
	FormData     map[string]string `json:"form_data" binding:"required"`
	WalletID     string            `json:"wallet_id" binding:"required"`
	PIN          string            `json:"pin" binding:"required"`
}

// VerifyTicketRequest represents the request to verify a ticket
type VerifyTicketRequest struct {
	TicketCode string `json:"ticket_code"`
	QRData     string `json:"qr_data"` // Can be code or full QR data
}

// VerifyTicketResponse represents the response for ticket verification
type VerifyTicketResponse struct {
	Valid       bool              `json:"valid"`
	Ticket      *Ticket           `json:"ticket,omitempty"`
	Event       *Event            `json:"event,omitempty"`
	Message     string            `json:"message"`
	CanUse      bool              `json:"can_use"`
	AlreadyUsed bool              `json:"already_used"`
}

// TicketStats represents statistics for an event
type TicketStats struct {
	TotalTickets   int                `json:"total_tickets"`
	SoldTickets    int                `json:"sold_tickets"`
	UsedTickets    int                `json:"used_tickets"`
	TotalRevenue   float64            `json:"total_revenue"`
	Currency       string             `json:"currency"`
	TierBreakdown  []TierStats        `json:"tier_breakdown"`
}

// TierStats represents statistics for a single tier
type TierStats struct {
	TierID      string  `json:"tier_id"`
	TierName    string  `json:"tier_name"`
	TierIcon    string  `json:"tier_icon"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Sold        int     `json:"sold"`
	Used        int     `json:"used"`
	Revenue     float64 `json:"revenue"`
}
