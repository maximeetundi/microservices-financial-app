package repository

import (
	"context"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository struct {
	collection *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Database) *EmployeeRepository {
	return &EmployeeRepository{
		collection: db.Collection("employees"),
	}
}

func (r *EmployeeRepository) Create(ctx context.Context, emp *models.Employee) error {
	emp.CreatedAt = time.Now()
	emp.UpdatedAt = time.Now()
	// Only set default status if not already set (owner is created with ACTIVE)
	if emp.Status == "" {
		emp.Status = models.EmployeeStatusPending
	}
	res, err := r.collection.InsertOne(ctx, emp)
	if err != nil {
		return err
	}
	emp.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id string) (*models.Employee, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var emp models.Employee
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&emp)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) FindByEnterprise(ctx context.Context, enterpriseID string) ([]models.Employee, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"enterprise_id": oid})
	if err != nil {
		return nil, err
	}
	var employees []models.Employee
	if err = cursor.All(ctx, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, emp *models.Employee) error {
	emp.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": emp.ID}, emp)
	return err
}

// FindByEnterpriseAndContact checks if an employee with given email or phone already exists in the enterprise
func (r *EmployeeRepository) FindByEnterpriseAndContact(ctx context.Context, enterpriseID, email, phone string) (*models.Employee, error) {
	oid, err := primitive.ObjectIDFromHex(enterpriseID)
	if err != nil {
		return nil, err
	}
	
	// Build OR filter for email or phone
	var orFilters []bson.M
	if email != "" {
		orFilters = append(orFilters, bson.M{"email": email})
	}
	if phone != "" {
		orFilters = append(orFilters, bson.M{"phone_number": phone})
	}
	
	if len(orFilters) == 0 {
		return nil, nil // No email or phone to check
	}
	
	filter := bson.M{
		"enterprise_id": oid,
		"$or":           orFilters,
	}
	
	var emp models.Employee
	err = r.collection.FindOne(ctx, filter).Decode(&emp)
	if err != nil {
		return nil, err // Not found or error
	}
	return &emp, nil
}
