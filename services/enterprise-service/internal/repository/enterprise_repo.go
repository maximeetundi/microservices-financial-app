package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnterpriseRepository struct {
	collection *mongo.Collection
}

func NewEnterpriseRepository(db *mongo.Database) *EnterpriseRepository {
	return &EnterpriseRepository{
		collection: db.Collection("enterprises"),
	}
}

func (r *EnterpriseRepository) Create(ctx context.Context, ent *models.Enterprise) error {
	ent.CreatedAt = time.Now()
	ent.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, ent)
	if err != nil {
		return err
	}
	ent.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *EnterpriseRepository) FindByID(ctx context.Context, id string) (*models.Enterprise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var ent models.Enterprise
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&ent)
	if err != nil {
		return nil, err
	}
	return &ent, nil
}

func (r *EnterpriseRepository) Update(ctx context.Context, ent *models.Enterprise) error {
	ent.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": ent.ID}, ent)
	return err
}
