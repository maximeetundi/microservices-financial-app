package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents a product category in a shop
type Category struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID      primitive.ObjectID `json:"shop_id" bson:"shop_id"`
	ParentID    primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	QRCode      string             `json:"qr_code" bson:"qr_code"`
	Order       int                `json:"order" bson:"order"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	ProductCount int               `json:"product_count" bson:"product_count"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateCategoryRequest is the request to create a new category
type CreateCategoryRequest struct {
	ShopID      string `json:"shop_id" binding:"required"`
	ParentID    string `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Order       int    `json:"order"`
}

// UpdateCategoryRequest is the request to update a category
type UpdateCategoryRequest struct {
	ParentID    *string `json:"parent_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	Order       *int    `json:"order"`
	IsActive    *bool   `json:"is_active"`
}

// CategoryWithChildren includes nested subcategories
type CategoryWithChildren struct {
	Category
	Children []CategoryWithChildren `json:"children,omitempty"`
}
