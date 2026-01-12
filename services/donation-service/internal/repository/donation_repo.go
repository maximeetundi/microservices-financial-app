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

type DonationRepository struct {
	collection *mongo.Collection
}

func NewDonationRepository(db *mongo.Database) *DonationRepository {
	return &DonationRepository{
		collection: db.Collection("donations"),
	}
}

func (r *DonationRepository) Create(donation *models.Donation) error {
	donation.CreatedAt = time.Now()
	donation.UpdatedAt = time.Now()
	donation.ID = primitive.NewObjectID()
	
	if donation.Status == "" {
		donation.Status = models.DonationStatusPending
	}

	_, err := r.collection.InsertOne(context.Background(), donation)
	return err
}

func (r *DonationRepository) GetByID(id string) (*models.Donation, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var donation models.Donation
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&donation)
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *DonationRepository) UpdateStatus(id string, status models.DonationStatus, txID string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}
	if txID != "" {
		update["transaction_id"] = txID
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	return err
}

func (r *DonationRepository) ListByCampaign(campaignID string, limit, offset int64) ([]*models.Donation, error) {
	campObjID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	// Only show paid donations in public list usually? 
	// Or maybe pending too for the user?
	// Assuming this is used for "Wall of Donors", so only Paid.
	
	filter := bson.M{
		"campaign_id": campObjID,
		"status":      models.DonationStatusPaid,
	}

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var donations []*models.Donation
	if err = cursor.All(context.Background(), &donations); err != nil {
		return nil, err
	}
	return donations, nil
}

func (r *DonationRepository) ListByDonor(donorID string) ([]*models.Donation, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(context.Background(), bson.M{"donor_id": donorID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var donations []*models.Donation
	if err = cursor.All(context.Background(), &donations); err != nil {
		return nil, err
	}
	return donations, nil
}

// GetByTransactionID retrieves a donation by its transaction ID (used in Payment Consumer)
func (r *DonationRepository) GetByTransactionID(txID string) (*models.Donation, error) {
	var donation models.Donation
	// Note: We might store txID in "transaction_id" or if it matches the PaymentRequestID.
	// For "request_id" we usually store it? 
	// In ticket-service we used "ReferenceID".
	// Let's assume we search by transaction_id field if populated, or we need to look differently.
	// Actually, when initiating payment, we set RequestID = "DON-<ID>".
	// PaymentStatusEvent returns RequestID.
	// So we should find by RequestID/Reference?
	// But our model doesn't have Request ID explicitly other than ID.
	// So we might need to parse ID from RequestID string, or store RequestID.
	
	// Implementation decision: We will put DonationID in the RequestID of payment.
	// So we can just fetch by ID if we extract it. 
	// BUT, adding GetByTransactionID is useful if we store the resulting TransactionID.
	
	err := r.collection.FindOne(context.Background(), bson.M{"transaction_id": txID}).Decode(&donation)
	if err != nil {
		return nil, err
	}
	return &donation, nil
}
