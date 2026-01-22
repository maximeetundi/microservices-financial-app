package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettings contains configuration for the shop
type ShopSettings struct {
	AllowPickup      bool   `json:"allow_pickup" bson:"allow_pickup"`
	AllowDelivery    bool   `json:"allow_delivery" bson:"allow_delivery"`
	DeliveryFee      float64 `json:"delivery_fee" bson:"delivery_fee"`
	MinOrderAmount   float64 `json:"min_order_amount" bson:"min_order_amount"`
	MaxOrderAmount   float64 `json:"max_order_amount" bson:"max_order_amount"`
	AutoAcceptOrders bool   `json:"auto_accept_orders" bson:"auto_accept_orders"`
}

// ShopStats contains statistics for the shop
type ShopStats struct {
	TotalProducts   int     `json:"total_products" bson:"total_products"`
	TotalOrders     int     `json:"total_orders" bson:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue" bson:"total_revenue"`
	AverageRating   float64 `json:"average_rating" bson:"average_rating"`
	TotalReviews    int     `json:"total_reviews" bson:"total_reviews"`
}

// ShopManager represents a user who can manage the shop
type ShopManager struct {
	UserID      string    `json:"user_id" bson:"user_id"`
	Email       string    `json:"email" bson:"email"`
	FirstName   string    `json:"first_name" bson:"first_name"`
	LastName    string    `json:"last_name" bson:"last_name"`
	Role        string    `json:"role" bson:"role"` // owner, admin, editor
	Permissions []string  `json:"permissions" bson:"permissions"`
	InvitedAt   time.Time `json:"invited_at" bson:"invited_at"`
	JoinedAt    time.Time `json:"joined_at,omitempty" bson:"joined_at,omitempty"`
	Status      string    `json:"status" bson:"status"` // pending, active, revoked
}

// Shop represents a store in the marketplace
type Shop struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerID     string             `json:"owner_id" bson:"owner_id"`
	OwnerType   string             `json:"owner_type" bson:"owner_type"` // user, enterprise
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	LogoURL     string             `json:"logo_url" bson:"logo_url"`
	BannerURL   string             `json:"banner_url" bson:"banner_url"`
	IsPublic    bool               `json:"is_public" bson:"is_public"`
	WalletID    string             `json:"wallet_id" bson:"wallet_id"`
	Currency    string             `json:"currency" bson:"currency"`
	Managers    []ShopManager      `json:"managers" bson:"managers"`
	Tags        []string           `json:"tags" bson:"tags"`
	QRCode      string             `json:"qr_code" bson:"qr_code"`
	Status      string             `json:"status" bson:"status"` // active, suspended, closed
	Settings    ShopSettings       `json:"settings" bson:"settings"`
	Stats       ShopStats          `json:"stats" bson:"stats"`
	Address     Address            `json:"address,omitempty" bson:"address,omitempty"`
	ContactInfo ContactInfo        `json:"contact_info,omitempty" bson:"contact_info,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Address for shop location or shipping
type Address struct {
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	Country    string `json:"country" bson:"country"`
	PostalCode string `json:"postal_code" bson:"postal_code"`
}

// ContactInfo for shop contact details
type ContactInfo struct {
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Website  string `json:"website" bson:"website"`
}

// CreateShopRequest is the request to create a new shop
type CreateShopRequest struct {
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description"`
	IsPublic    bool         `json:"is_public"`
	WalletID    string       `json:"wallet_id" binding:"required"`
	Currency    string       `json:"currency" binding:"required"`
	Tags        []string     `json:"tags"`
	Address     *Address     `json:"address"`
	ContactInfo *ContactInfo `json:"contact_info"`
	LogoURL     string       `json:"logo_url"`
	BannerURL   string       `json:"banner_url"`
}

// UpdateShopRequest is the request to update a shop
type UpdateShopRequest struct {
	Name        *string       `json:"name"`
	Description *string       `json:"description"`
	IsPublic    *bool         `json:"is_public"`
	WalletID    *string       `json:"wallet_id"`
	Tags        []string      `json:"tags"`
	Settings    *ShopSettings `json:"settings"`
	Address     *Address      `json:"address"`
	ContactInfo *ContactInfo  `json:"contact_info"`
	LogoURL     *string       `json:"logo_url"`
	BannerURL   *string       `json:"banner_url"`
}

// InviteManagerRequest is the request to invite a manager
type InviteManagerRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Role        string   `json:"role" binding:"required"` // admin, editor
	Permissions []string `json:"permissions"`
}

// ShopListResponse for paginated shop listing
type ShopListResponse struct {
	Shops      []Shop `json:"shops"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}
