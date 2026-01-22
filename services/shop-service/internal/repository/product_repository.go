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

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	coll := db.Collection("products")
	
	// Create indexes
	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "shop_id", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "shop_id", Value: 1}}},
		{Keys: bson.D{{Key: "category_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
		{Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}}},
	})
	
	return &ProductRepository{collection: coll}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Slug = generateProductSlug(product.Name)
	
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, shopID primitive.ObjectID, slug string) (*models.Product, error) {
	var product models.Product
	filter := bson.M{"shop_id": shopID, "slug": slug}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) ListByShop(ctx context.Context, shopID primitive.ObjectID, page, pageSize int, categoryID *primitive.ObjectID, status string, search string) ([]models.Product, int64, error) {
	filter := bson.M{"shop_id": shopID}
	
	if categoryID != nil && !categoryID.IsZero() {
		filter["category_id"] = *categoryID
	}
	if status != "" {
		filter["status"] = status
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
		SetSort(bson.D{{Key: "is_featured", Value: -1}, {Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProductRepository) IncrementViewCount(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{"$inc": bson.M{"view_count": 1}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *ProductRepository) IncrementSoldCount(ctx context.Context, id primitive.ObjectID, quantity int) error {
	update := bson.M{
		"$inc": bson.M{
			"sold_count": quantity,
			"stock":      -quantity,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *ProductRepository) CountByShop(ctx context.Context, shopID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"shop_id": shopID, "status": "active"})
}

func (r *ProductRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"category_id": categoryID})
}

func (r *ProductRepository) SetQRCode(ctx context.Context, productID primitive.ObjectID, qrCodeURL string) error {
	update := bson.M{
		"$set": bson.M{
			"qr_code":    qrCodeURL,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": productID}, update)
	return err
}

func generateProductSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug + "-" + primitive.NewObjectID().Hex()[:6]
}
