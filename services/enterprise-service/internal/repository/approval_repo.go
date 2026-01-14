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

type ApprovalRepository struct {
	collection *mongo.Collection
}

func NewApprovalRepository(db *mongo.Database) *ApprovalRepository {
	return &ApprovalRepository{
		collection: db.Collection("action_approvals"),
	}
}

func (r *ApprovalRepository) Create(ctx context.Context, approval *models.ActionApproval) error {
	approval.CreatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, approval)
	if err != nil {
		return err
	}
	approval.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ApprovalRepository) FindByID(ctx context.Context, id string) (*models.ActionApproval, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var approval models.ActionApproval
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&approval)
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

func (r *ApprovalRepository) FindPendingByEnterprise(ctx context.Context, enterpriseID string) ([]models.ActionApproval, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	
	filter := bson.M{
		"enterprise_id": oid,
		"status":        models.ApprovalStatusPending,
		"expires_at":    bson.M{"$gt": time.Now()},
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	
	var approvals []models.ActionApproval
	if err = cursor.All(ctx, &approvals); err != nil {
		return nil, err
	}
	return approvals, nil
}

func (r *ApprovalRepository) Update(ctx context.Context, approval *models.ActionApproval) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": approval.ID}, approval)
	return err
}

func (r *ApprovalRepository) AddApproval(ctx context.Context, approvalID string, adminApproval models.AdminApproval) error {
	oid, err := primitive.ObjectIDFromHex(approvalID)
	if err != nil {
		return err
	}
	
	update := bson.M{
		"$push": bson.M{"approvals": adminApproval},
	}
	
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

func (r *ApprovalRepository) UpdateStatus(ctx context.Context, approvalID string, status models.ActionApprovalStatus) error {
	oid, err := primitive.ObjectIDFromHex(approvalID)
	if err != nil {
		return err
	}
	
	update := bson.M{
		"$set": bson.M{"status": status},
	}
	
	if status == models.ApprovalStatusExecuted {
		now := time.Now()
		update["$set"].(bson.M)["executed_at"] = &now
	}
	
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

// ExpireOldApprovals marks expired pending approvals
func (r *ApprovalRepository) ExpireOldApprovals(ctx context.Context) (int64, error) {
	filter := bson.M{
		"status":     models.ApprovalStatusPending,
		"expires_at": bson.M{"$lt": time.Now()},
	}
	
	update := bson.M{
		"$set": bson.M{"status": models.ApprovalStatusExpired},
	}
	
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
