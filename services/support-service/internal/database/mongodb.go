package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client        *mongo.Client
	Database      *mongo.Database
	Conversations *mongo.Collection
	Messages      *mongo.Collection
	Agents        *mongo.Collection
)

// Connect establishes connection to MongoDB
func Connect() error {
	uri := getEnv("MONGODB_URI", "mongodb://mongodb:27017")
	dbName := getEnv("MONGODB_DATABASE", "support_service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	Database = client.Database(dbName)
	Conversations = Database.Collection("conversations")
	Messages = Database.Collection("messages")
	Agents = Database.Collection("agents")

	log.Printf("Connected to MongoDB: %s/%s", uri, dbName)
	return nil
}

// InitSchema creates indexes for collections
func InitSchema() error {
	ctx := context.Background()

	// Conversations indexes
	convIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "agent_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}
	if _, err := Conversations.Indexes().CreateMany(ctx, convIndexes); err != nil {
		log.Printf("Warning: Failed to create conversations indexes: %v", err)
	}

	// Messages indexes
	msgIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "conversation_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},
	}
	if _, err := Messages.Indexes().CreateMany(ctx, msgIndexes); err != nil {
		log.Printf("Warning: Failed to create messages indexes: %v", err)
	}

	// Agents indexes
	agentIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "is_available", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
	}
	if _, err := Agents.Indexes().CreateMany(ctx, agentIndexes); err != nil {
		log.Printf("Warning: Failed to create agents indexes: %v", err)
	}

	log.Println("MongoDB indexes initialized")
	return nil
}

// Close closes the MongoDB connection
func Close() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		Client.Disconnect(ctx)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
