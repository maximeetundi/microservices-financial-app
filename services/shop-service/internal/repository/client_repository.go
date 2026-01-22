package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) *ClientRepository {
	coll := db.Collection("shop_clients")
	
	// Create indexes
	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "shop_id", Value: 1}, {Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "shop_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	
	return &ClientRepository{collection: coll}
}

func (r *ClientRepository) Create(ctx context.Context, client *models.ShopClient) error {
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	client.InvitedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, client)
	if err != nil {
		return err
	}
	client.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.ShopClient, error) {
	var client models.ShopClient
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetByShopAndEmail(ctx context.Context, shopID primitive.ObjectID, email string) (*models.ShopClient, error) {
	var client models.ShopClient
	err := r.collection.FindOne(ctx, bson.M{"shop_id": shopID, "email": email}).Decode(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetByShopAndUserID(ctx context.Context, shopID primitive.ObjectID, userID string) (*models.ShopClient, error) {
	var client models.ShopClient
	err := r.collection.FindOne(ctx, bson.M{"shop_id": shopID, "user_id": userID, "status": models.ClientStatusActive}).Decode(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetPendingInvitationsByEmail returns all pending invitations for an email
func (r *ClientRepository) GetPendingInvitationsByEmail(ctx context.Context, email string) ([]models.ShopClient, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"email":  email,
		"status": models.ClientStatusPending,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var clients []models.ShopClient
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

// ListByShop returns all clients for a shop
func (r *ClientRepository) ListByShop(ctx context.Context, shopID primitive.ObjectID, page, pageSize int, status string) ([]models.ShopClient, int64, error) {
	filter := bson.M{"shop_id": shopID}
	if status != "" {
		filter["status"] = status
	}
	
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var clients []models.ShopClient
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, 0, err
	}
	return clients, total, nil
}

// ListActiveByUserID returns all active shop accesses for a user
func (r *ClientRepository) ListActiveByUserID(ctx context.Context, userID string) ([]models.ShopClient, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id": userID,
		"status":  models.ClientStatusActive,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var clients []models.ShopClient
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) Update(ctx context.Context, client *models.ShopClient) error {
	client.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": client.ID}, client)
	return err
}

// AcceptInvitation updates the client status to active
func (r *ClientRepository) AcceptInvitation(ctx context.Context, id primitive.ObjectID, userID string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"user_id":     userID,
			"status":      models.ClientStatusActive,
			"accepted_at": now,
			"updated_at":  now,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// RevokeAccess revokes a client's access to store
func (r *ClientRepository) RevokeAccess(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":     models.ClientStatusRevoked,
			"revoked_at": now,
			"updated_at": now,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// DeclineInvitation marks invitation as declined
func (r *ClientRepository) DeclineInvitation(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":     models.ClientStatusDeclined,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// HasAccess checks if a user has access to a private shop
func (r *ClientRepository) HasAccess(ctx context.Context, shopID primitive.ObjectID, userID string) bool {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"shop_id": shopID,
		"user_id": userID,
		"status":  models.ClientStatusActive,
	})
	return err == nil && count > 0
}

func (r *ClientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
