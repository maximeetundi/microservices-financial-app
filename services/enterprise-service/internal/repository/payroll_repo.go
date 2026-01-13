package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PayrollRepository struct {
	collection *mongo.Collection
}

func NewPayrollRepository(db *mongo.Database) *PayrollRepository {
	return &PayrollRepository{
		collection: db.Collection("payrolls"),
	}
}

func (r *PayrollRepository) Create(ctx context.Context, p *models.PayrollRun) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *PayrollRepository) FindByEnterprise(ctx context.Context, enterpriseID string) ([]models.PayrollRun, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"enterprise_id": oid})
	if err != nil {
		return nil, err
	}
	var runs []models.PayrollRun
	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}
	return runs, nil
}
