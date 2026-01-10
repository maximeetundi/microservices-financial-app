package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *MessageRepository) Create(msg *models.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	msg.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), msg)
	return err
}

func (r *MessageRepository) GetByConversationID(conversationID string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	filter := bson.M{"conversation_id": conversationID}
	
	opts := options.Find().
		SetSort(bson.M{"created_at": 1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		// Ensure attachments slice is initialized
		if msg.Attachments == nil {
			msg.Attachments = []string{}
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *MessageRepository) MarkAsRead(conversationID, userID string) error {
	filter := bson.M{
		"conversation_id": conversationID,
		"sender_id":       bson.M{"$ne": userID},
		"is_read":         false,
	}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateMany(context.Background(), filter, update)
	return err
}

func (r *MessageRepository) GetUnreadCount(conversationID, userID string) (int, error) {
	filter := bson.M{
		"conversation_id": conversationID,
		"sender_id":       bson.M{"$ne": userID},
		"is_read":         false,
	}
	count, err := r.collection.CountDocuments(context.Background(), filter)
	return int(count), err
}
