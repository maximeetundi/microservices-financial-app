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

func (r *EnterpriseRepository) FindAll(ctx context.Context) ([]*models.Enterprise, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var enterprises []*models.Enterprise
	if err := cursor.All(ctx, &enterprises); err != nil {
		return nil, err
	}
	
	if enterprises == nil {
		enterprises = []*models.Enterprise{}
	}
	
	return enterprises, nil
}

// FindByUserAccess returns enterprises where user is owner OR is an employee
func (r *EnterpriseRepository) FindByUserAccess(ctx context.Context, userID string) ([]*models.Enterprise, error) {
	// Get enterprises where user is owner
	ownerFilter := bson.M{"owner_id": userID}
	
	cursor, err := r.collection.Find(ctx, ownerFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ownedEnterprises []*models.Enterprise
	if err := cursor.All(ctx, &ownedEnterprises); err != nil {
		return nil, err
	}
	
	// Get enterprises where user is an employee
	// Query employees collection for this user
	employeesCol := r.collection.Database().Collection("employees")
	empCursor, err := employeesCol.Find(ctx, bson.M{
		"user_id": userID,
		"status": bson.M{"$in": []string{"ACTIVE", "PENDING", "PENDING_INVITE"}},
	})
	if err != nil {
		// If employees query fails, just return owned enterprises
		if ownedEnterprises == nil {
			ownedEnterprises = []*models.Enterprise{}
		}
		return ownedEnterprises, nil
	}
	defer empCursor.Close(ctx)
	
	// Collect enterprise IDs from employee records
	enterpriseIDSet := make(map[string]bool)
	for _, ent := range ownedEnterprises {
		enterpriseIDSet[ent.ID.Hex()] = true
	}
	
	for empCursor.Next(ctx) {
		var emp struct {
			EnterpriseID primitive.ObjectID `bson:"enterprise_id"`
		}
		if empCursor.Decode(&emp) == nil {
			enterpriseIDSet[emp.EnterpriseID.Hex()] = true
		}
	}
	
	// If user has employee records in other enterprises, fetch those too
	var additionalIDs []primitive.ObjectID
	for _, ent := range ownedEnterprises {
		delete(enterpriseIDSet, ent.ID.Hex()) // Already have these
	}
	for idHex := range enterpriseIDSet {
		if oid, err := primitive.ObjectIDFromHex(idHex); err == nil {
			additionalIDs = append(additionalIDs, oid)
		}
	}
	
	if len(additionalIDs) > 0 {
		addCursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": additionalIDs}})
		if err == nil {
			defer addCursor.Close(ctx)
			var additionalEnterprises []*models.Enterprise
			if addCursor.All(ctx, &additionalEnterprises) == nil {
				ownedEnterprises = append(ownedEnterprises, additionalEnterprises...)
			}
		}
	}
	
	if ownedEnterprises == nil {
		ownedEnterprises = []*models.Enterprise{}
	}
	
	return ownedEnterprises, nil
}

// Delete removes an enterprise by ID
func (r *EnterpriseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

