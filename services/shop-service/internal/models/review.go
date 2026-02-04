package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a product review
type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	UserName  string             `json:"user_name" bson:"user_name"` // Snapshot of user name
	Rating    int                `json:"rating" bson:"rating"`       // 1-5
	Comment   string             `json:"comment" bson:"comment"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// CreateReviewRequest
type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required,max=1000"`
}

// ReviewListResponse
type ReviewListResponse struct {
	Reviews []Review `json:"reviews"`
	Total   int64    `json:"total"`
}
