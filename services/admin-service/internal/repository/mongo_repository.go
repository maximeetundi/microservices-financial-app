package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: client,
	}
}

// ========== Donation Service Collections ==========

func (r *MongoRepository) GetCampaigns(limit, offset int) ([]map[string]interface{}, error) {
	collection := r.client.Database("donation_service").Collection("campaigns")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var campaigns []map[string]interface{}
	if err = cursor.All(ctx, &campaigns); err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *MongoRepository) GetDonations(limit, offset int) ([]map[string]interface{}, error) {
	collection := r.client.Database("donation_service").Collection("donations")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var donations []map[string]interface{}
	if err = cursor.All(ctx, &donations); err != nil {
		return nil, err
	}
	return donations, nil
}

// ========== Ticket Service Collections ==========

func (r *MongoRepository) GetEvents(limit, offset int) ([]map[string]interface{}, error) {
	collection := r.client.Database("ticket_service").Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []map[string]interface{}
	if err = cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *MongoRepository) GetSoldTickets(limit, offset int) ([]map[string]interface{}, error) {
	collection := r.client.Database("ticket_service").Collection("tickets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []map[string]interface{}
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}
