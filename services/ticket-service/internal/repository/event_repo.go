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

type EventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(db *mongo.Database) *EventRepository {
	return &EventRepository{
		collection: db.Collection("events"),
	}
}

func (r *EventRepository) Create(event *models.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), event)
	return err
}

func (r *EventRepository) GetByID(id string) (*models.Event, error) {
	var event models.Event
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetByCode(code string) (*models.Event, error) {
	var event models.Event
	filter := bson.M{"event_code": code}
	err := r.collection.FindOne(context.Background(), filter).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetByCreator(creatorID string, limit, offset int) ([]*models.Event, error) {
	var events []*models.Event
	filter := bson.M{"creator_id": creatorID}
	
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *EventRepository) GetActive(limit, offset int) ([]*models.Event, error) {
	var events []*models.Event
	filter := bson.M{
		"status": "active",
		"sale_end_date": bson.M{"$gt": time.Now()},
	}

	opts := options.Find().
		SetSort(bson.M{"start_date": 1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	event.UpdatedAt = time.Now()
	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *EventRepository) UpdateStatus(id, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *EventRepository) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}
