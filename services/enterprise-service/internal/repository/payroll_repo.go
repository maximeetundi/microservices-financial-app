package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *PayrollRepository) FindByEnterpriseAndYear(ctx context.Context, enterpriseID string, year int) ([]models.PayrollRun, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"enterprise_id": oid,
		"period_year":   year,
	}

	opts := options.Find().SetSort(bson.D{{Key: "period_month", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var runs []models.PayrollRun
	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}
	return runs, nil
}

// Update updates an existing payroll run
func (r *PayrollRepository) Update(ctx context.Context, p *models.PayrollRun) error {
	p.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": p.ID},
		bson.M{"$set": p},
	)
	return err
}
