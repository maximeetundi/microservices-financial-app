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

type ConversationRepository struct {
	collection *mongo.Collection
}

func NewConversationRepository(db *mongo.Database) *ConversationRepository {
	return &ConversationRepository{
		collection: db.Collection("conversations"),
	}
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	conv.ID = uuid.New().String()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
	conv.Status = models.ConversationStatusOpen

	_, err := r.collection.InsertOne(context.Background(), conv)
	return err
}

func (r *ConversationRepository) GetByID(id string) (*models.Conversation, error) {
	var conv models.Conversation
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&conv)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *ConversationRepository) GetByUserID(userID string, limit, offset int) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	filter := bson.M{"user_id": userID}
	
	opts := options.Find().
		SetSort(bson.M{"updated_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var conv models.Conversation
		if err := cursor.Decode(&conv); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (r *ConversationRepository) GetAll(status string, limit, offset int) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	// Custom sort: Priority then UpdatedAt
	// Note: In MongoDB, we can't easily do 'CASE WHEN' in sort without aggregation.
	// For simplicity, we'll sort by updated_at descending, effectively showing newest first.
	// If priority sorting is strictly required, we'd need an aggregation pipeline or store an integer priority.
	opts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var conv models.Conversation
		if err := cursor.Decode(&conv); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (r *ConversationRepository) UpdateStatus(id string, status models.ConversationStatus) error {
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

func (r *ConversationRepository) AssignAgent(id, agentID string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"agent_id":   agentID,
			"status":     models.ConversationStatusActive,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ConversationRepository) UpdateLastMessage(id, message string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"last_message":    message,
			"last_message_at": time.Now(),
			"updated_at":      time.Now(),
		},
		"$inc": bson.M{"message_count": 1},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ConversationRepository) Close(id string, rating int, feedback string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":      models.ConversationStatusClosed,
			"rating":      rating,
			"feedback":    feedback,
			"resolved_at": time.Now(),
			"updated_at":  time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ConversationRepository) GetStats() (*models.SupportStats, error) {
	stats := &models.SupportStats{}

	// Aggregation to get most stats in one go
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total_conversations": bson.M{"$sum": 1},
			"open_conversations": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$in": bson.A{"$status", bson.A{"open", "pending", "active"}}},
						1,
						0,
					},
				},
			},
			"resolved_today": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$and": bson.A{
							bson.M{"$eq": bson.A{"$status", "resolved"}},
							bson.M{"$gte": bson.A{"$resolved_at", time.Now().Truncate(24 * time.Hour)}},
						}},
						1,
						0,
					},
				},
			},
			"pending_conversations": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$status", "pending"}},
						1,
						0,
					},
				},
			},
			"rating_sum": bson.M{"$sum": "$rating"},
			"rating_count": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$rating", false}},
						1,
						0,
					},
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		var result struct {
			TotalConversations   int     `bson:"total_conversations"`
			OpenConversations    int     `bson:"open_conversations"`
			ResolvedToday        int     `bson:"resolved_today"`
			PendingConversations int     `bson:"pending_conversations"`
			RatingSum            float64 `bson:"rating_sum"`
			RatingCount          int     `bson:"rating_count"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.TotalConversations = result.TotalConversations
			stats.OpenConversations = result.OpenConversations
			stats.ResolvedToday = result.ResolvedToday
			stats.PendingConversations = result.PendingConversations
			if result.RatingCount > 0 {
				stats.CustomerSatisfaction = result.RatingSum / float64(result.RatingCount)
			}
		}
	}

	return stats, nil
}
