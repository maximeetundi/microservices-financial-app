package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomField represents a customizable field for a product
type CustomField struct {
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"` // text, number, select, color
	Required    bool     `json:"required" bson:"required"`
	Options     []string `json:"options,omitempty" bson:"options,omitempty"` // For select type
	Placeholder string   `json:"placeholder,omitempty" bson:"placeholder,omitempty"`
	MinValue    float64  `json:"min_value,omitempty" bson:"min_value,omitempty"`
	MaxValue    float64  `json:"max_value,omitempty" bson:"max_value,omitempty"`
}

// ProductVariant represents a variant of a product (size, color, etc.)
type ProductVariant struct {
	ID         string            `json:"id" bson:"id"`
	Name       string            `json:"name" bson:"name"`
	SKU        string            `json:"sku" bson:"sku"`
	Price      float64           `json:"price" bson:"price"`
	Stock      int               `json:"stock" bson:"stock"`
	ImageURL   string            `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Attributes map[string]string `json:"attributes" bson:"attributes"`
}

// Product represents an item for sale in a shop
type Product struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID         primitive.ObjectID `json:"shop_id" bson:"shop_id"`
	CategoryID     primitive.ObjectID `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Slug           string             `json:"slug" bson:"slug"`
	Description    string             `json:"description" bson:"description"`
	ShortDesc      string             `json:"short_desc" bson:"short_desc"`
	Price          float64            `json:"price" bson:"price"`
	CompareAtPrice float64            `json:"compare_at_price,omitempty" bson:"compare_at_price,omitempty"`
	Currency       string             `json:"currency" bson:"currency"`
	Images         []string           `json:"images" bson:"images"`
	IsDigital      bool               `json:"is_digital" bson:"is_digital"`
	DigitalFileURL string             `json:"digital_file_url,omitempty" bson:"digital_file_url,omitempty"`
	LicenseText    string             `json:"license_text,omitempty" bson:"license_text,omitempty"`
	Stock          int                `json:"stock" bson:"stock"` // -1 for unlimited
	SKU            string             `json:"sku,omitempty" bson:"sku,omitempty"`
	Barcode        string             `json:"barcode,omitempty" bson:"barcode,omitempty"`
	Weight         float64            `json:"weight,omitempty" bson:"weight,omitempty"` // in grams
	IsCustomizable bool               `json:"is_customizable" bson:"is_customizable"`
	CustomFields   []CustomField      `json:"custom_fields,omitempty" bson:"custom_fields,omitempty"`
	Variants       []ProductVariant   `json:"variants,omitempty" bson:"variants,omitempty"`
	Tags           []string           `json:"tags" bson:"tags"`
	QRCode         string             `json:"qr_code" bson:"qr_code"`
	Status         string             `json:"status" bson:"status"` // active, draft, archived
	IsFeatured     bool               `json:"is_featured" bson:"is_featured"`
	SoldCount      int                `json:"sold_count" bson:"sold_count"`
	ViewCount      int                `json:"view_count" bson:"view_count"`
	AverageRating  float64            `json:"average_rating" bson:"average_rating"`
	ReviewCount    int                `json:"review_count" bson:"review_count"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateProductRequest is the request to create a new product
type CreateProductRequest struct {
	ShopID         string        `json:"shop_id" binding:"required"`
	CategoryID     string        `json:"category_id"`
	Name           string        `json:"name" binding:"required"`
	Description    string        `json:"description"`
	ShortDesc      string        `json:"short_desc"`
	Price          float64       `json:"price" binding:"required,gt=0"`
	CompareAtPrice float64       `json:"compare_at_price"`
	Images         []string      `json:"images"`
	IsDigital      bool          `json:"is_digital"`
	DigitalFileURL string        `json:"digital_file_url"`
	LicenseText    string        `json:"license_text"`
	Stock          int           `json:"stock"`
	SKU            string        `json:"sku"`
	Weight         float64       `json:"weight"`
	IsCustomizable bool          `json:"is_customizable"`
	CustomFields   []CustomField `json:"custom_fields"`
	Tags           []string      `json:"tags"`
	Status         string        `json:"status"`
	IsFeatured     bool          `json:"is_featured"`
}

// UpdateProductRequest is the request to update a product
type UpdateProductRequest struct {
	CategoryID     *string       `json:"category_id"`
	Name           *string       `json:"name"`
	Description    *string       `json:"description"`
	ShortDesc      *string       `json:"short_desc"`
	Price          *float64      `json:"price"`
	CompareAtPrice *float64      `json:"compare_at_price"`
	Images         []string      `json:"images"`
	IsDigital      *bool         `json:"is_digital"`
	DigitalFileURL *string       `json:"digital_file_url"`
	LicenseText    *string       `json:"license_text"`
	Stock          *int          `json:"stock"`
	SKU            *string       `json:"sku"`
	Weight         *float64      `json:"weight"`
	IsCustomizable *bool         `json:"is_customizable"`
	CustomFields   []CustomField `json:"custom_fields"`
	Tags           []string      `json:"tags"`
	Status         *string       `json:"status"`
	IsFeatured     *bool         `json:"is_featured"`
}

// ProductListResponse for paginated product listing
type ProductListResponse struct {
	Products   []Product `json:"products"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}
