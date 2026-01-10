package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TierRepository struct {
	collection *mongo.Collection
}

func NewTierRepository(db *mongo.Database) *TierRepository {
	return &TierRepository{
		collection: db.Collection("ticket_tiers"),
	}
}

func (r *TierRepository) Create(tier *models.TicketTier) error {
	if tier.ID == "" {
		tier.ID = uuid.New().String()
	}
	tier.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), tier)
	return err
}

func (r *TierRepository) GetByEventID(eventID string) ([]*models.TicketTier, error) {
	var tiers []*models.TicketTier
	filter := bson.M{"event_id": eventID}
	
	opts := options.Find().SetSort(bson.D{
		{Key: "sort_order", Value: 1},
		{Key: "price", Value: 1},
	})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var tier models.TicketTier
		if err := cursor.Decode(&tier); err != nil {
			return nil, err
		}
		tiers = append(tiers, &tier)
	}

	return tiers, nil
}

func (r *TierRepository) GetByID(id string) (*models.TicketTier, error) {
	var tier models.TicketTier
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&tier)
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *TierRepository) IncrementSold(tierID string) error {
	filter := bson.M{"_id": tierID}
	update := bson.M{"$inc": bson.M{"sold": 1}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TierRepository) DecrementSold(tierID string) error {
	filter := bson.M{"_id": tierID, "sold": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"sold": -1}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *TierRepository) DeleteByEventID(eventID string) error {
	filter := bson.M{"event_id": eventID}
	_, err := r.collection.DeleteMany(context.Background(), filter)
	return err
}

func (r *TierRepository) CheckAvailability(tierID string) (bool, error) {
	var tier models.TicketTier
	filter := bson.M{"_id": tierID}
	
	// Only fetch needed fields
	opts := options.FindOne().SetProjection(bson.M{"quantity": 1, "sold": 1})
	
	err := r.collection.FindOne(context.Background(), filter, opts).Decode(&tier)
	if err != nil {
		return false, err
	}

	if tier.Quantity == -1 {
		return true, nil
	}

	return tier.Sold < tier.Quantity, nil
}
