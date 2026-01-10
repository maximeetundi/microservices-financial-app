package models

import "time"

// Ticket represents a purchased ticket
type Ticket struct {
	ID            string            `json:"id" bson:"_id,omitempty"`
	EventID       string            `json:"event_id" bson:"event_id"`
	BuyerID       string            `json:"buyer_id" bson:"buyer_id"`
	TierID        string            `json:"tier_id" bson:"tier_id"`
	TierName      string            `json:"tier_name" bson:"tier_name"`
	TierIcon      string            `json:"tier_icon" bson:"tier_icon"`
	TierColor     string            `json:"tier_color" bson:"tier_color"`
	Price         float64           `json:"price" bson:"price"`
	Currency      string            `json:"currency" bson:"currency"`
	FormData      map[string]string `json:"form_data" bson:"form_data"`
	QRCode        string            `json:"qr_code" bson:"qr_code"`
	TicketCode    string            `json:"ticket_code" bson:"ticket_code"`
	Status        string            `json:"status" bson:"status"` // pending, paid, used, cancelled, refunded
	TransactionID string            `json:"transaction_id" bson:"transaction_id"`
	UsedAt        *time.Time        `json:"used_at,omitempty" bson:"used_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`

	// Join fields (not stored in ticket document)
	EventTitle    string     `json:"event_title,omitempty" bson:"-"`
	EventDate     *time.Time `json:"event_date,omitempty" bson:"-"`
	EventLocation string     `json:"event_location,omitempty" bson:"-"`
	EventCreatorID string    `json:"event_creator_id,omitempty" bson:"-"`
	BuyerName     string     `json:"buyer_name,omitempty" bson:"-"`
	BuyerEmail    string     `json:"buyer_email,omitempty" bson:"-"`
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
	EventID  string            `json:"event_id" binding:"required"`
	TierID   string            `json:"tier_id" binding:"required"`
	FormData map[string]string `json:"form_data" binding:"required"`
	WalletID string            `json:"wallet_id" binding:"required"`
	PIN      string            `json:"pin" binding:"required"`
}

// VerifyTicketRequest represents the request to verify a ticket
type VerifyTicketRequest struct {
	TicketCode string `json:"ticket_code"`
	QRData     string `json:"qr_data"` // Can be code or full QR data
}

// VerifyTicketResponse represents the response for ticket verification
type VerifyTicketResponse struct {
	Valid       bool    `json:"valid"`
	Ticket      *Ticket `json:"ticket,omitempty"`
	Event       *Event  `json:"event,omitempty"`
	Message     string  `json:"message"`
	CanUse      bool    `json:"can_use"`
	AlreadyUsed bool    `json:"already_used"`
}

// TicketStats represents statistics for an event
type TicketStats struct {
	TotalTickets  int         `json:"total_tickets"`
	SoldTickets   int         `json:"sold_tickets"`
	UsedTickets   int         `json:"used_tickets"`
	TotalRevenue  float64     `json:"total_revenue"`
	Currency      string      `json:"currency"`
	TierBreakdown []TierStats `json:"tier_breakdown"`
}

// TierStats represents statistics for a single tier
type TierStats struct {
	TierID   string  `json:"tier_id"`
	TierName string  `json:"tier_name"`
	TierIcon string  `json:"tier_icon"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Sold     int     `json:"sold"`
	Used     int     `json:"used"`
	Revenue  float64 `json:"revenue"`
}
