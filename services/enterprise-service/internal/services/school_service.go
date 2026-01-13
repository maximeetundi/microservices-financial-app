package services

import (
	"context"
	"errors"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
)

type SchoolService struct {
	entRepo *repository.EnterpriseRepository
	subRepo *repository.SubscriptionRepository
	invRepo *repository.InvoiceRepository
}

func NewSchoolService(entRepo *repository.EnterpriseRepository, subRepo *repository.SubscriptionRepository, invRepo *repository.InvoiceRepository) *SchoolService {
	return &SchoolService{
		entRepo: entRepo,
		subRepo: subRepo,
		invRepo: invRepo,
	}
}

// GetStudentStatus (Reporting for School)
func (s *SchoolService) GetStudentStatus(ctx context.Context, enterpriseID string, classID string) ([]map[string]interface{}, error) {
	// 1. Get Enterprise Config to know Tranches
	ent, err := s.entRepo.FindByID(ctx, enterpriseID)
	if err != nil { return nil, err }
	if ent.SchoolConfig == nil { return nil, errors.New("not a school") }
	
	// Find correct class config
	var classConfig models.Class
	found := false
	for _, c := range ent.SchoolConfig.Classes {
		if c.ID == classID {
			classConfig = c
			found = true
			break
		}
	}
	if !found { return nil, errors.New("class not found") }

	// 2. Get Subscriptions (Students) for this class
	// Assuming Subscription has metadata for "class_id"
	subs, err := s.subRepo.FindByEnterprise(ctx, enterpriseID)
	if err != nil { return nil, err }

	var report []map[string]interface{}

	for _, sub := range subs {
		// Filter by class (simplified, ideally repo does this)
		if sub.Metadata["class_id"] != classID {
			continue
		}

		// Calculate Financial Status
		invoices, _ := s.invRepo.FindBySubscriptionID(ctx, sub.ID.Hex())
		totalPaid := 0.0
		for _, inv := range invoices {
			if inv.Status == "PAID" {
				totalPaid += inv.Amount
			}
		}

		// Calculate Total Due based on active tranches (simplified: sum all tranches)
		totalDue := 0.0
		for _, t := range classConfig.Tranches {
			totalDue += t.Amount
		}

		status := "OK"
		if totalPaid < totalDue {
			status = "LATE" 
		}

		studentReport := map[string]interface{}{
			"student_id": sub.ClientID, // Using UserID as Student Identifier
			"name": sub.Metadata["student_name"],
			"total_due": totalDue,
			"total_paid": totalPaid,
			"balance": totalDue - totalPaid,
			"status": status,
		}
		report = append(report, studentReport)
	}
	
	return report, nil
}

// CreateTrancheInvoice Logic
func (s *SchoolService) CreateTrancheInvoice(ctx context.Context, subscriptionID string, trancheID string) (*models.Invoice, error) {
	// 1. Find Subscription
	sub, err := s.subRepo.FindByID(ctx, subscriptionID)
	if err != nil { return nil, err }

	// 2. Get Enterprise Config
	ent, err := s.entRepo.FindByID(ctx, sub.EnterpriseID.Hex())
	if err != nil { return nil, err }
	if ent.SchoolConfig == nil { return nil, errors.New("enterprise is not a school") }

	// 3. Find Tranche Details
	var trancheAmount float64
	var trancheName string
	classID := sub.Metadata["class_id"]
	
	found := false
	for _, c := range ent.SchoolConfig.Classes {
		if c.ID == classID {
			for _, t := range c.Tranches {
				if t.ID == trancheID {
					trancheAmount = t.Amount
					trancheName = t.Name
					found = true
					break
				}
			}
		}
	}
	if !found { return nil, errors.New("tranche not found") }

	// 4. Create Invoice
	inv := &models.Invoice{
		EnterpriseID: sub.EnterpriseID,
		SubscriptionID: sub.ID,
		Amount: trancheAmount,
		Status: "PENDING",
		Type: "TUITION", 
		Metadata: map[string]string{
			"tranche_id": trancheID,
			"tranche_name": trancheName,
			"student_name": sub.Metadata["student_name"],
		},
		// DueDate: tranche.Deadline...
	}

	if err := s.invRepo.Create(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}
