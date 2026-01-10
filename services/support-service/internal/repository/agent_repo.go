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

type AgentRepository struct {
	collection *mongo.Collection
}

func NewAgentRepository(db *mongo.Database) *AgentRepository {
	return &AgentRepository{
		collection: db.Collection("agents"),
	}
}

func (r *AgentRepository) Create(agent *models.Agent) error {
	if agent.ID == "" {
		agent.ID = uuid.New().String()
	}
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), agent)
	return err
}

func (r *AgentRepository) GetByID(id string) (*models.Agent, error) {
	var agent models.Agent
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.Background(), filter).Decode(&agent)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepository) GetAvailable(agentType models.AgentType) ([]*models.Agent, error) {
	var agents []*models.Agent
	filter := bson.M{
		"type":         agentType,
		"is_available": true,
		"$expr":        bson.M{"$lt": bson.A{"$active_chats", "$max_chats"}},
	}
	
	opts := options.Find().SetSort(bson.D{
		{Key: "active_chats", Value: 1},
		{Key: "rating", Value: -1},
	})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var agent models.Agent
		if err := cursor.Decode(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

func (r *AgentRepository) GetAll() ([]*models.Agent, error) {
	var agents []*models.Agent
	opts := options.Find().SetSort(bson.D{
		{Key: "type", Value: 1},
		{Key: "name", Value: 1},
	})

	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var agent models.Agent
		if err := cursor.Decode(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

func (r *AgentRepository) UpdateAvailability(id string, isAvailable bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_available": isAvailable,
			"updated_at":   time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *AgentRepository) IncrementActiveChats(id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"active_chats": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *AgentRepository) DecrementActiveChats(id string) error {
	// Standard decrement but prevent below 0 handled by logic or simply just dec
	// Mongo $inc -1 is fine, logic ensures it shouldn't go below 0 usually.
	// We can use $max to ensure it doesn't go below 0 if concurrent updates are risky, 
	// but simpler here:
	
	filter := bson.M{"_id": id, "active_chats": bson.M{"$gt": 0}}
	update := bson.M{
		"$inc": bson.M{"active_chats": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *AgentRepository) GetAIAgent() (*models.Agent, error) {
	filter := bson.M{"type": models.AgentTypeAI}
	var agent models.Agent
	
	err := r.collection.FindOne(context.Background(), filter).Decode(&agent)
	if err == mongo.ErrNoDocuments {
		// Create default AI agent
		agent = models.Agent{
			Name:        "Assistant IA",
			Email:       "ai@zekora.com",
			Type:        models.AgentTypeAI,
			Avatar:      "🤖",
			IsAvailable: true,
			MaxChats:    1000,
			ActiveChats: 0,
			Rating:      5.0,
		}
		if err := r.Create(&agent); err != nil {
			return nil, err
		}
		return &agent, nil
	}
	if err != nil {
		return nil, err
	}

	return &agent, nil
}
