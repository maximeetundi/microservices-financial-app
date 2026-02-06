package repository

import (
	"context"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopRepository struct {
	collection *mongo.Collection
}

func NewShopRepository(db *mongo.Database) *ShopRepository {
	coll := db.Collection("shops")
	
	// Create indexes
	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "owner_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "is_public", Value: 1}}},
		{Keys: bson.D{{Key: "managers.user_id", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
	})
	
	return &ShopRepository{collection: coll}
}

func (r *ShopRepository) Create(ctx context.Context, shop *models.Shop) error {
	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()
	shop.Slug = generateSlug(shop.Name)
	
	result, err := r.collection.InsertOne(ctx, shop)
	if err != nil {
		return err
	}
	shop.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ShopRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Shop, error) {
	var shop models.Shop
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&shop)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *ShopRepository) GetBySlug(ctx context.Context, slug string) (*models.Shop, error) {
	var shop models.Shop
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&shop)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *ShopRepository) GetByWalletID(ctx context.Context, walletID string) (*models.Shop, error) {
	var shop models.Shop
	err := r.collection.FindOne(ctx, bson.M{"wallet_id": walletID}).Decode(&shop)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *ShopRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]models.Shop, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var shops []models.Shop
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopRepository) GetByManagerID(ctx context.Context, userID string) ([]models.Shop, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"managers.user_id": userID, "managers.status": "active"},
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var shops []models.Shop
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopRepository) ListPublic(ctx context.Context, page, pageSize int, search string) ([]models.Shop, int64, error) {
	filter := bson.M{
		"is_public": true,
		"status":    "active",
	}
	
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
			{"tags": bson.M{"$regex": search, "$options": "i"}},
		}
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
	
	var shops []models.Shop
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, 0, err
	}
	return shops, total, nil
}

func (r *ShopRepository) Update(ctx context.Context, shop *models.Shop) error {
	shop.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": shop.ID}, shop)
	return err
}

func (r *ShopRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ShopRepository) UpdateStats(ctx context.Context, shopID primitive.ObjectID, stats models.ShopStats) error {
	update := bson.M{
		"$set": bson.M{
			"stats":      stats,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": shopID}, update)
	return err
}

func (r *ShopRepository) AddManager(ctx context.Context, shopID primitive.ObjectID, manager models.ShopManager) error {
	update := bson.M{
		"$push": bson.M{"managers": manager},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": shopID}, update)
	return err
}

func (r *ShopRepository) RemoveManager(ctx context.Context, shopID primitive.ObjectID, userID string) error {
	update := bson.M{
		"$pull": bson.M{"managers": bson.M{"user_id": userID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": shopID}, update)
	return err
}

func (r *ShopRepository) SetQRCode(ctx context.Context, shopID primitive.ObjectID, qrCodeURL string) error {
	update := bson.M{
		"$set": bson.M{
			"qr_code":    qrCodeURL,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": shopID}, update)
	return err
}

func (r *ShopRepository) SetTrustBadges(ctx context.Context, shopID primitive.ObjectID, badges []models.ShopTrustBadge) error {
	update := bson.M{
		"$set": bson.M{
			"trust_badges": badges,
			"updated_at":   time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": shopID}, update)
	return err
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	// Add timestamp suffix to ensure uniqueness
	return slug + "-" + primitive.NewObjectID().Hex()[:6]
}
