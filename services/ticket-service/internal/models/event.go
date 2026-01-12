package models

import "time"

// Event represents an event created by an organizer
type Event struct {
	ID            string       `json:"id" bson:"_id,omitempty"`
	CreatorID     string       `json:"creator_id" bson:"creator_id"`
	Title         string       `json:"title" bson:"title"`
	Description   string       `json:"description" bson:"description"`
	Location      string       `json:"location" bson:"location"`
	CoverImage    string       `json:"cover_image" bson:"cover_image"`
	StartDate     time.Time    `json:"start_date" bson:"start_date"`
	EndDate       time.Time    `json:"end_date" bson:"end_date"`
	SaleStartDate time.Time    `json:"sale_start_date" bson:"sale_start_date"`
	SaleEndDate   time.Time    `json:"sale_end_date" bson:"sale_end_date"`
	FormFields    []FormField  `json:"form_fields" bson:"form_fields"`
	TicketTiers   []TicketTier `json:"ticket_tiers" bson:"-"` // Loaded separately
	QRCode        string       `json:"qr_code" bson:"qr_code"`
	EventCode     string       `json:"event_code" bson:"event_code"`
	Status        string       `json:"status" bson:"status"` // draft, active, ended, cancelled
	Currency      string       `json:"currency" bson:"currency"`
	CreatedAt     time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" bson:"updated_at"`

	// Computed fields (not stored)
	TotalSold    int     `json:"total_sold,omitempty" bson:"-"`
	TotalRevenue float64 `json:"total_revenue,omitempty" bson:"-"`
	CreatorName  string  `json:"creator_name,omitempty" bson:"-"`
}

// FormField represents a custom form field for ticket purchase
type FormField struct {
	Name     string   `json:"name" bson:"name"`
	Label    string   `json:"label" bson:"label"`
	Type     string   `json:"type" bson:"type"` // text, email, phone, select, checkbox, number, date
	Required bool     `json:"required" bson:"required"`
	Options  []string `json:"options,omitempty" bson:"options,omitempty"` // for select type
}

// TicketTier represents a pricing tier for an event
type TicketTier struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	EventID     string    `json:"event_id" bson:"event_id"`
	Name        string    `json:"name" bson:"name"`
	Icon        string    `json:"icon" bson:"icon"` // emoji or icon name from predefined list
	Price       float64   `json:"price" bson:"price"`
	Quantity         int       `json:"quantity" bson:"quantity"` // -1 for unlimited
	DisplayRemaining bool      `json:"display_remaining" bson:"display_remaining"`
	Sold             int       `json:"sold" bson:"sold"`
	Description      string    `json:"description" bson:"description"`
	Benefits         []string  `json:"benefits" bson:"benefits"`
	Color            string    `json:"color" bson:"color"` // hex color for UI
	SortOrder        int       `json:"sort_order" bson:"sort_order"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
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
