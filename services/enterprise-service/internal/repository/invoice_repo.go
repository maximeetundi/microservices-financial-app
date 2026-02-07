package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepository struct {
	collection *mongo.Collection
}

func NewInvoiceRepository(db *mongo.Database) *InvoiceRepository {
	return &InvoiceRepository{
		collection: db.Collection("invoices"),
	}
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *models.Invoice) error {
	inv.CreatedAt = time.Now()
	inv.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, inv)
	if err != nil {
		return err
	}
	inv.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *InvoiceRepository) FindByEnterprise(ctx context.Context, enterpriseID string) ([]models.Invoice, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"enterprise_id": oid})
	if err != nil {
		return nil, err
	}
	var invs []models.Invoice
	if err = cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

func (r *InvoiceRepository) FindBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Invoice, error) {
	oid, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"subscription_id": oid})
	if err != nil {
		return nil, err
	}
	var invs []models.Invoice
	if err = cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

func (r *InvoiceRepository) FindByEnterpriseAndClientID(ctx context.Context, enterpriseID, clientID string) ([]models.Invoice, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"enterprise_id": oid, "client_id": clientID})
	if err != nil {
		return nil, err
	}
	var invs []models.Invoice
	if err = cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}
