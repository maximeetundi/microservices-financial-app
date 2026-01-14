package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionRepository(db *mongo.Database) *SubscriptionRepository {
	return &SubscriptionRepository{
		collection: db.Collection("subscriptions"),
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	res, err := r.collection.InsertOne(ctx, sub)
	if err != nil {
		return err
	}
	sub.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindDueSubscriptions finds subscriptions that need billing (NextBillingAt <= Now)
func (r *SubscriptionRepository) FindDueSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"status":          "ACTIVE",
		"next_billing_at": bson.M{"$lte": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	var subs []models.Subscription
	if err = cursor.All(ctx, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) FindByEnterprise(ctx context.Context, enterpriseID string) ([]models.Subscription, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"enterprise_id": oid})
	if err != nil {
		return nil, err
	}
	var subs []models.Subscription
	if err = cursor.All(ctx, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) FindByID(ctx context.Context, id string) (*models.Subscription, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var sub models.Subscription
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": sub.ID}, bson.M{"$set": sub})
	return err
}
