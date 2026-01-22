package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents a single item in an order
type OrderItem struct {
	ProductID    primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName  string             `json:"product_name" bson:"product_name"`
	ProductImage string             `json:"product_image" bson:"product_image"`
	VariantID    string             `json:"variant_id,omitempty" bson:"variant_id,omitempty"`
	VariantName  string             `json:"variant_name,omitempty" bson:"variant_name,omitempty"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	UnitPrice    float64            `json:"unit_price" bson:"unit_price"`
	TotalPrice   float64            `json:"total_price" bson:"total_price"`
	CustomValues map[string]string  `json:"custom_values,omitempty" bson:"custom_values,omitempty"`
}

// Order represents a purchase transaction
type Order struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderNumber       string             `json:"order_number" bson:"order_number"`
	ShopID            primitive.ObjectID `json:"shop_id" bson:"shop_id"`
	ShopName          string             `json:"shop_name" bson:"shop_name"`
	BuyerID           string             `json:"buyer_id" bson:"buyer_id"`
	BuyerName         string             `json:"buyer_name" bson:"buyer_name"`
	BuyerEmail        string             `json:"buyer_email" bson:"buyer_email"`
	Items             []OrderItem        `json:"items" bson:"items"`
	SubTotal          float64            `json:"sub_total" bson:"sub_total"`
	DeliveryFee       float64            `json:"delivery_fee" bson:"delivery_fee"`
	TotalAmount       float64            `json:"total_amount" bson:"total_amount"`
	Currency          string             `json:"currency" bson:"currency"`
	BuyerWalletID     string             `json:"buyer_wallet_id" bson:"buyer_wallet_id"`
	BuyerCurrency     string             `json:"buyer_currency" bson:"buyer_currency"`
	ConvertedAmount   float64            `json:"converted_amount,omitempty" bson:"converted_amount,omitempty"`
	ExchangeRate      float64            `json:"exchange_rate,omitempty" bson:"exchange_rate,omitempty"`
	TransactionID     string             `json:"transaction_id" bson:"transaction_id"`
	PaymentStatus     string             `json:"payment_status" bson:"payment_status"` // pending, completed, failed, refunded
	OrderStatus       string             `json:"order_status" bson:"order_status"`     // pending, confirmed, processing, shipped, delivered, cancelled
	DeliveryType      string             `json:"delivery_type" bson:"delivery_type"`   // pickup, delivery
	ShippingAddress   *Address           `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	Notes             string             `json:"notes" bson:"notes"`
	SellerNotes       string             `json:"seller_notes" bson:"seller_notes"`
	TrackingNumber    string             `json:"tracking_number,omitempty" bson:"tracking_number,omitempty"`
	RefundReason      string             `json:"refund_reason,omitempty" bson:"refund_reason,omitempty"`
	RefundedAt        *time.Time         `json:"refunded_at,omitempty" bson:"refunded_at,omitempty"`
	CompletedAt       *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// CartItem for checkout request
type CartItem struct {
	ProductID    string            `json:"product_id" binding:"required"`
	VariantID    string            `json:"variant_id"`
	Quantity     int               `json:"quantity" binding:"required,gt=0"`
	CustomValues map[string]string `json:"custom_values"`
}

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	ShopID          string     `json:"shop_id" binding:"required"`
	Items           []CartItem `json:"items" binding:"required,min=1"`
	WalletID        string     `json:"wallet_id" binding:"required"`
	DeliveryType    string     `json:"delivery_type" binding:"required"` // pickup, delivery
	ShippingAddress *Address   `json:"shipping_address"`
	Notes           string     `json:"notes"`
}

// UpdateOrderStatusRequest for updating order status
type UpdateOrderStatusRequest struct {
	Status         string `json:"status" binding:"required"`
	TrackingNumber string `json:"tracking_number"`
	SellerNotes    string `json:"seller_notes"`
}

// RefundOrderRequest for refunding an order
type RefundOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// OrderListResponse for paginated order listing
type OrderListResponse struct {
	Orders     []Order `json:"orders"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}
