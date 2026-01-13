package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CampaignRepository struct {
	collection *mongo.Collection
}

func NewCampaignRepository(db *mongo.Database) *CampaignRepository {
	return &CampaignRepository{
		collection: db.Collection("campaigns"),
	}
}

func (r *CampaignRepository) Create(campaign *models.Campaign) error {
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	campaign.ID = primitive.NewObjectID()
	
	if campaign.Status == "" {
		campaign.Status = models.CampaignStatusDraft
	}

	_, err := r.collection.InsertOne(context.Background(), campaign)
	return err
}

func (r *CampaignRepository) GetByID(id string) (*models.Campaign, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var campaign models.Campaign
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&campaign)
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignRepository) GetByCode(code string) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.collection.FindOne(context.Background(), bson.M{"campaign_code": code}).Decode(&campaign)
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignRepository) Update(id string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": updates},
	)
	return err
}

func (r *CampaignRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	return err
}

func (r *CampaignRepository) List(limit, offset int64) ([]*models.Campaign, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	
	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campaigns []*models.Campaign
	if err = cursor.All(context.Background(), &campaigns); err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *CampaignRepository) ListByCreator(creatorID string) ([]*models.Campaign, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(context.Background(), bson.M{"creator_id": creatorID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campaigns []*models.Campaign // Corrected type
	if err = cursor.All(context.Background(), &campaigns); err != nil {
		return nil, err
	}
	return campaigns, nil
}

// IncrementCollectedAmounts updates the collected amount of a campaign
func (r *CampaignRepository) IncrementCollectedAmount(id string, amount float64) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{
			"$inc": bson.M{"collected_amount": amount},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}
