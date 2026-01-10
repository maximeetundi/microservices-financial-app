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
	Client     *mongo.Client
	Database   *mongo.Database
	Events     *mongo.Collection
	Tiers      *mongo.Collection
	Tickets    *mongo.Collection
)

// Connect establishes connection to MongoDB
func Connect() error {
	uri := getEnv("MONGODB_URI", "mongodb://mongodb:27017")
	dbName := getEnv("MONGODB_DATABASE", "ticket_service")

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
	Events = Database.Collection("events")
	Tiers = Database.Collection("ticket_tiers")
	Tickets = Database.Collection("tickets")

	log.Printf("Connected to MongoDB: %s/%s", uri, dbName)
	return nil
}

// InitSchema creates indexes for collections
func InitSchema() error {
	ctx := context.Background()

	// Events indexes
	eventsIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "creator_id", Value: 1}}},
		{Keys: bson.D{{Key: "event_code", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "start_date", Value: 1}}},
	}
	if _, err := Events.Indexes().CreateMany(ctx, eventsIndexes); err != nil {
		log.Printf("Warning: Failed to create events indexes: %v", err)
	}

	// Tiers indexes
	tiersIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "event_id", Value: 1}}},
	}
	if _, err := Tiers.Indexes().CreateMany(ctx, tiersIndexes); err != nil {
		log.Printf("Warning: Failed to create tiers indexes: %v", err)
	}

	// Tickets indexes
	ticketsIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "event_id", Value: 1}}},
		{Keys: bson.D{{Key: "buyer_id", Value: 1}}},
		{Keys: bson.D{{Key: "ticket_code", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "transaction_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	}
	if _, err := Tickets.Indexes().CreateMany(ctx, ticketsIndexes); err != nil {
		log.Printf("Warning: Failed to create tickets indexes: %v", err)
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
