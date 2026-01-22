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

type CategoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	coll := db.Collection("categories")
	
	// Create indexes
	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "shop_id", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "shop_id", Value: 1}}},
		{Keys: bson.D{{Key: "parent_id", Value: 1}}},
	})
	
	return &CategoryRepository{collection: coll}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.Slug = generateCategorySlug(category.Name)
	category.IsActive = true
	
	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Category, error) {
	var category models.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, shopID primitive.ObjectID, slug string) (*models.Category, error) {
	var category models.Category
	filter := bson.M{"shop_id": shopID, "slug": slug}
	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) ListByShop(ctx context.Context, shopID primitive.ObjectID) ([]models.Category, error) {
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "name", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"shop_id": shopID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) ListRootCategories(ctx context.Context, shopID primitive.ObjectID) ([]models.Category, error) {
	filter := bson.M{
		"shop_id":   shopID,
		"parent_id": primitive.NilObjectID,
		"is_active": true,
	}
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) ListChildren(ctx context.Context, parentID primitive.ObjectID) ([]models.Category, error) {
	filter := bson.M{"parent_id": parentID, "is_active": true}
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	category.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *CategoryRepository) UpdateProductCount(ctx context.Context, categoryID primitive.ObjectID, count int) error {
	update := bson.M{
		"$set": bson.M{
			"product_count": count,
			"updated_at":    time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": categoryID}, update)
	return err
}

func (r *CategoryRepository) SetQRCode(ctx context.Context, categoryID primitive.ObjectID, qrCodeURL string) error {
	update := bson.M{
		"$set": bson.M{
			"qr_code":    qrCodeURL,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": categoryID}, update)
	return err
}

func generateCategorySlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	return slug + "-" + primitive.NewObjectID().Hex()[:6]
}
