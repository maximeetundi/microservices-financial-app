package models

import "time"

// Event represents an event created by an organizer
type Event struct {
	ID             string       `json:"id" db:"id"`
	CreatorID      string       `json:"creator_id" db:"creator_id"`
	Title          string       `json:"title" db:"title"`
	Description    string       `json:"description" db:"description"`
	Location       string       `json:"location" db:"location"`
	CoverImage     string       `json:"cover_image" db:"cover_image"`
	StartDate      time.Time    `json:"start_date" db:"start_date"`
	EndDate        time.Time    `json:"end_date" db:"end_date"`
	SaleStartDate  time.Time    `json:"sale_start_date" db:"sale_start_date"`
	SaleEndDate    time.Time    `json:"sale_end_date" db:"sale_end_date"`
	FormFields     []FormField  `json:"form_fields"`
	TicketTiers    []TicketTier `json:"ticket_tiers"`
	QRCode         string       `json:"qr_code" db:"qr_code"`
	EventCode      string       `json:"event_code" db:"event_code"`
	Status         string       `json:"status" db:"status"` // draft, active, ended, cancelled
	Currency       string       `json:"currency" db:"currency"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	
	// Computed fields
	TotalSold      int          `json:"total_sold,omitempty"`
	TotalRevenue   float64      `json:"total_revenue,omitempty"`
	CreatorName    string       `json:"creator_name,omitempty"`
}

// FormField represents a custom form field for ticket purchase
type FormField struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"` // text, email, phone, select, checkbox, number, date
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"` // for select type
}

// TicketTier represents a pricing tier for an event
type TicketTier struct {
	ID          string   `json:"id" db:"id"`
	EventID     string   `json:"event_id" db:"event_id"`
	Name        string   `json:"name" db:"name"`
	Icon        string   `json:"icon" db:"icon"` // emoji or icon name from predefined list
	Price       float64  `json:"price" db:"price"`
	Quantity    int      `json:"quantity" db:"quantity"` // -1 for unlimited
	Sold        int      `json:"sold" db:"sold"`
	Description string   `json:"description" db:"description"`
	Benefits    []string `json:"benefits"`
	Color       string   `json:"color" db:"color"` // hex color for UI
	SortOrder   int      `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Predefined icons for ticket tiers
var AvailableTierIcons = []string{
	"⭐", "🌟", "✨", "💎", "👑", "🏆", "🎖️", "🥇", "🥈", "🥉",
	"🎫", "🎟️", "🎪", "🎭", "🎬", "🎵", "🎸", "🎤", "🎧", "🎹",
	"🔥", "💫", "⚡", "🌈", "🎯", "🚀", "💥", "🎉", "🎊", "🎁",
	"VIP", "PRO", "GOLD", "SILVER", "BRONZE", "PREMIUM", "ELITE",
}

// Event status constants
const (
	EventStatusDraft     = "draft"
	EventStatusActive    = "active"
	EventStatusEnded     = "ended"
	EventStatusCancelled = "cancelled"
)

// CreateEventRequest represents the request to create an event
type CreateEventRequest struct {
	Title         string       `json:"title" binding:"required"`
	Description   string       `json:"description"`
	Location      string       `json:"location"`
	CoverImage    string       `json:"cover_image"`
	StartDate     time.Time    `json:"start_date" binding:"required"`
	EndDate       time.Time    `json:"end_date" binding:"required"`
	SaleStartDate time.Time    `json:"sale_start_date" binding:"required"`
	SaleEndDate   time.Time    `json:"sale_end_date" binding:"required"`
	FormFields    []FormField  `json:"form_fields"`
	TicketTiers   []TicketTier `json:"ticket_tiers" binding:"required"`
	Currency      string       `json:"currency"`
}

// UpdateEventRequest represents the request to update an event
type UpdateEventRequest struct {
	Title         *string      `json:"title,omitempty"`
	Description   *string      `json:"description,omitempty"`
	Location      *string      `json:"location,omitempty"`
	CoverImage    *string      `json:"cover_image,omitempty"`
	StartDate     *time.Time   `json:"start_date,omitempty"`
	EndDate       *time.Time   `json:"end_date,omitempty"`
	SaleStartDate *time.Time   `json:"sale_start_date,omitempty"`
	SaleEndDate   *time.Time   `json:"sale_end_date,omitempty"`
	FormFields    []FormField  `json:"form_fields,omitempty"`
	TicketTiers   []TicketTier `json:"ticket_tiers,omitempty"`
	Status        *string      `json:"status,omitempty"`
}
