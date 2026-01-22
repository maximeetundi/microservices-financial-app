package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ClientStatus defines the status of a private shop client
type ClientStatus string

const (
	ClientStatusPending  ClientStatus = "pending"  // Invitation sent, awaiting acceptance
	ClientStatusActive   ClientStatus = "active"   // Client accepted and can access shop
	ClientStatusRevoked  ClientStatus = "revoked"  // Access revoked by owner
	ClientStatusDeclined ClientStatus = "declined" // Client declined invitation
)

// ShopClient represents an authorized client for a private shop
type ShopClient struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID    primitive.ObjectID `json:"shop_id" bson:"shop_id"`
	UserID    string             `json:"user_id,omitempty" bson:"user_id,omitempty"` // Set when accepted
	
	// Invitation info (before acceptance)
	Email       string `json:"email" bson:"email"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
	FirstName   string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	
	// Status
	Status    ClientStatus `json:"status" bson:"status"`
	InvitedBy string       `json:"invited_by" bson:"invited_by"` // UserID who sent invitation
	
	// Timestamps
	InvitedAt  time.Time  `json:"invited_at" bson:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty" bson:"accepted_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	
	// Optional: Client-specific settings
	Notes       string   `json:"notes,omitempty" bson:"notes,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"` // VIP, Wholesale, etc.
	Discount    float64  `json:"discount,omitempty" bson:"discount,omitempty"` // Client-specific discount %
	
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// InviteClientRequest - request to invite a client to a private shop
type InviteClientRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	Phone     string   `json:"phone,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Discount  float64  `json:"discount,omitempty"`
}

// AcceptClientInvitationRequest - request to accept an invitation (requires PIN)
type AcceptClientInvitationRequest struct {
	InvitationID string `json:"invitation_id" binding:"required"`
	PIN          string `json:"pin" binding:"required,len=6"`
}

// ClientInvitationResponse - response after invitation is sent
type ClientInvitationResponse struct {
	ID        string    `json:"id"`
	ShopID    string    `json:"shop_id"`
	ShopName  string    `json:"shop_name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	InvitedAt time.Time `json:"invited_at"`
}

// MyInvitationsResponse - list of pending invitations for a user
type MyInvitationsResponse struct {
	Invitations []ClientInvitationDetail `json:"invitations"`
}

// ClientInvitationDetail - detailed view of an invitation
type ClientInvitationDetail struct {
	ID          string    `json:"id"`
	ShopID      string    `json:"shop_id"`
	ShopName    string    `json:"shop_name"`
	ShopLogo    string    `json:"shop_logo,omitempty"`
	InvitedBy   string    `json:"invited_by"`
	InviterName string    `json:"inviter_name,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	Status      string    `json:"status"`
	InvitedAt   time.Time `json:"invited_at"`
}

// ShopClientsResponse - list of shop clients
type ShopClientsResponse struct {
	Clients    []ShopClient `json:"clients"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
