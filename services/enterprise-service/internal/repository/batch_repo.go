package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchRepository struct {
	collection *mongo.Collection
}

func NewBatchRepository(db *mongo.Database) *BatchRepository {
	return &BatchRepository{
		collection: db.Collection("invoice_batches"),
	}
}

func (r *BatchRepository) Create(ctx context.Context, batch *models.InvoiceBatch) error {
	batch.CreatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, batch)
	if err != nil {
		return err
	}
	batch.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *BatchRepository) FindByID(ctx context.Context, id string) (*models.InvoiceBatch, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var batch models.InvoiceBatch
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&batch)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepository) Update(ctx context.Context, batch *models.InvoiceBatch) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": batch.ID}, batch)
	return err
}

func (r *BatchRepository) FindPendingScheduled(ctx context.Context) ([]models.InvoiceBatch, error) {
	// Status = SCHEDULED and Time <= Now
	filter := bson.M{
		"status":       "SCHEDULED",
		"scheduled_at": bson.M{"$lte": time.Now()},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var batches []models.InvoiceBatch
	if err = cursor.All(ctx, &batches); err != nil {
		return nil, err
	}
	return batches, nil
}
