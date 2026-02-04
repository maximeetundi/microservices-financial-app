package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{
		collection: db.Collection("reviews"),
	}
}

func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	review.CreatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	review.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ReviewRepository) GetByProduct(ctx context.Context, productID primitive.ObjectID, page, pageSize int) ([]models.Review, int64, error) {
	filter := bson.M{"product_id": productID}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
