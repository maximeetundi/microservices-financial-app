package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	coll := db.Collection("orders")
	
	// Create indexes
	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "order_number", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "shop_id", Value: 1}}},
		{Keys: bson.D{{Key: "buyer_id", Value: 1}}},
		{Keys: bson.D{{Key: "transaction_id", Value: 1}}},
		{Keys: bson.D{{Key: "payment_status", Value: 1}}},
		{Keys: bson.D{{Key: "order_status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	
	return &OrderRepository{collection: coll}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.OrderNumber = generateOrderNumber()
	
	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByTransactionID(ctx context.Context, transactionID string) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"transaction_id": transactionID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"order_number": orderNumber}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]models.Order, int64, error) {
	filter := bson.M{"buyer_id": buyerID}
	
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
	
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) ListByShop(ctx context.Context, shopID primitive.ObjectID, page, pageSize int, status string) ([]models.Order, int64, error) {
	filter := bson.M{"shop_id": shopID}
	if status != "" {
		filter["order_status"] = status
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
	
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	order.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	return err
}

func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, transactionID, status string) error {
	update := bson.M{
		"$set": bson.M{
			"payment_status": status,
			"updated_at":     time.Now(),
		},
	}
	if status == "completed" {
		update["$set"].(bson.M)["order_status"] = "confirmed"
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"transaction_id": transactionID}, update)
	return err
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID primitive.ObjectID, status string) error {
	update := bson.M{
		"$set": bson.M{
			"order_status": status,
			"updated_at":   time.Now(),
		},
	}
	if status == "delivered" {
		now := time.Now()
		update["$set"].(bson.M)["completed_at"] = &now
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": orderID}, update)
	return err
}

func (r *OrderRepository) MarkRefunded(ctx context.Context, orderID primitive.ObjectID, reason string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"payment_status": "refunded",
			"order_status":   "cancelled",
			"refund_reason":  reason,
			"refunded_at":    &now,
			"updated_at":     now,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": orderID}, update)
	return err
}

func (r *OrderRepository) CountByShop(ctx context.Context, shopID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"shop_id": shopID, "payment_status": "completed"})
}

func (r *OrderRepository) SumRevenueByShop(ctx context.Context, shopID primitive.ObjectID) (float64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"shop_id": shopID, "payment_status": "completed"}},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$total_amount"}}},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	
	var results []struct {
		Total float64 `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Total, nil
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d-%s", time.Now().Unix(), primitive.NewObjectID().Hex()[:6])
}
