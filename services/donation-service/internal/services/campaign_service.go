package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"

	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/repository"
)

type CampaignService struct {
	repo *repository.CampaignRepository
}

func NewCampaignService(repo *repository.CampaignRepository) *CampaignService {
	return &CampaignService{repo: repo}
}

func (s *CampaignService) CreateCampaign(campaign *models.Campaign) error {
	if campaign.Title == "" || campaign.Description == "" {
		return errors.New("title and description are required")
	}
	if campaign.CreatorID == "" {
		return errors.New("creator_id is required")
	}
	
	// Generate Campaign Code
	campaign.CampaignCode = s.generateCampaignCode()

	// Generate QR Code
	qrData := fmt.Sprintf("ZEKORA_CAMPAIGN:%s", campaign.CampaignCode)
	campaign.QRCode = s.generateQRCodeBase64(qrData)

	campaign.Status = models.CampaignStatusActive // Default to active
	
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	
	return s.repo.Create(campaign)
}

func (s *CampaignService) GetCampaign(id string) (*models.Campaign, error) {
	return s.repo.GetByID(id)
}

func (s *CampaignService) ListCampaigns(limit, offset int64) ([]*models.Campaign, error) {
	return s.repo.List(limit, offset)
}

func (s *CampaignService) GetMyCampaigns(creatorID string) ([]*models.Campaign, error) {
	return s.repo.ListByCreator(creatorID)
}

func (s *CampaignService) UpdateCampaign(id, creatorID string, updates map[string]interface{}) error {
	// Verify ownership
	campaign, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if campaign.CreatorID != creatorID {
		return errors.New("unauthorized")
	}
	
	return s.repo.Update(id, updates)
}

// Helpers

func (s *CampaignService) generateCampaignCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	code := strings.ToUpper(base64.RawURLEncoding.EncodeToString(b))
	return "CPN-" + code[:8]
}

func (s *CampaignService) generateQRCodeBase64(data string) string {
	qr, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qr)
}
