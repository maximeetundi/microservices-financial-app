package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientService struct {
	clientRepo  *repository.ClientRepository
	shopRepo    *repository.ShopRepository
	kafkaClient *messaging.KafkaClient
	authURL     string
}

func NewClientService(
	clientRepo *repository.ClientRepository,
	shopRepo *repository.ShopRepository,
	kafkaClient *messaging.KafkaClient,
) *ClientService {
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://auth-service:8081"
	}
	return &ClientService{
		clientRepo:  clientRepo,
		shopRepo:    shopRepo,
		kafkaClient: kafkaClient,
		authURL:     authURL,
	}
}

// InviteClient invites a client to a private shop
func (s *ClientService) InviteClient(ctx context.Context, shopID string, req *models.InviteClientRequest, inviterUserID string) (*models.ShopClient, error) {
	shopOID, err := parseObjectID(shopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID")
	}

	// Get shop and verify it's private
	shop, err := s.shopRepo.GetByID(ctx, shopOID)
	if err != nil {
		return nil, fmt.Errorf("shop not found")
	}

	// Only private shops need client invitations
	if shop.IsPublic {
		return nil, fmt.Errorf("public shops don't require client invitations")
	}

	// Check inviter has permission (owner or admin)
	if !s.hasManagePermission(shop, inviterUserID) {
		return nil, fmt.Errorf("permission denied")
	}

	// Check if client already invited
	existing, _ := s.clientRepo.GetByShopAndEmail(ctx, shopOID, req.Email)
	if existing != nil {
		if existing.Status == models.ClientStatusActive {
			return nil, fmt.Errorf("client already has access")
		}
		if existing.Status == models.ClientStatusPending {
			return nil, fmt.Errorf("invitation already pending")
		}
		// If revoked or declined, allow re-invite
	}

	client := &models.ShopClient{
		ShopID:    shopOID,
		Email:     req.Email,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    models.ClientStatusPending,
		InvitedBy: inviterUserID,
		Notes:     req.Notes,
		Tags:      req.Tags,
		Discount:  req.Discount,
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// Send notification to client (via email/push)
	s.sendInvitationNotification(client, shop)

	return client, nil
}

// GetPendingInvitations returns pending invitations for a user by their email
func (s *ClientService) GetPendingInvitations(ctx context.Context, email string) ([]models.ClientInvitationDetail, error) {
	clients, err := s.clientRepo.GetPendingInvitationsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var invitations []models.ClientInvitationDetail
	for _, c := range clients {
		shop, err := s.shopRepo.GetByID(ctx, c.ShopID)
		if err != nil {
			continue
		}

		invitations = append(invitations, models.ClientInvitationDetail{
			ID:        c.ID.Hex(),
			ShopID:    c.ShopID.Hex(),
			ShopName:  shop.Name,
			ShopLogo:  shop.LogoURL,
			InvitedBy: c.InvitedBy,
			Notes:     c.Notes,
			Status:    string(c.Status),
			InvitedAt: c.InvitedAt,
		})
	}

	return invitations, nil
}

// AcceptInvitation accepts a shop invitation with PIN verification
func (s *ClientService) AcceptInvitation(ctx context.Context, invitationID, userID, email, pin string) error {
	// Verify PIN with auth-service
	if err := s.verifyPIN(userID, pin); err != nil {
		return fmt.Errorf("invalid PIN: %w", err)
	}

	invOID, err := parseObjectID(invitationID)
	if err != nil {
		return fmt.Errorf("invalid invitation ID")
	}

	invitation, err := s.clientRepo.GetByID(ctx, invOID)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	// Verify the email matches
	if invitation.Email != email {
		return fmt.Errorf("invitation not for this user")
	}

	if invitation.Status != models.ClientStatusPending {
		return fmt.Errorf("invitation already processed")
	}

	// Accept invitation
	if err := s.clientRepo.AcceptInvitation(ctx, invOID, userID); err != nil {
		return fmt.Errorf("failed to accept invitation: %w", err)
	}

	// Update shop client count (could be in stats)
	// Notify shop owner
	s.sendAcceptedNotification(invitation, userID)

	return nil
}

// DeclineInvitation declines a shop invitation
func (s *ClientService) DeclineInvitation(ctx context.Context, invitationID, email string) error {
	invOID, err := parseObjectID(invitationID)
	if err != nil {
		return fmt.Errorf("invalid invitation ID")
	}

	invitation, err := s.clientRepo.GetByID(ctx, invOID)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	if invitation.Email != email {
		return fmt.Errorf("invitation not for this user")
	}

	if invitation.Status != models.ClientStatusPending {
		return fmt.Errorf("invitation already processed")
	}

	return s.clientRepo.DeclineInvitation(ctx, invOID)
}

// RevokeClientAccess revokes a client's access to a private shop
func (s *ClientService) RevokeClientAccess(ctx context.Context, shopID, clientID, userID string) error {
	shopOID, err := parseObjectID(shopID)
	if err != nil {
		return fmt.Errorf("invalid shop ID")
	}

	shop, err := s.shopRepo.GetByID(ctx, shopOID)
	if err != nil {
		return fmt.Errorf("shop not found")
	}

	if !s.hasManagePermission(shop, userID) {
		return fmt.Errorf("permission denied")
	}

	clientOID, err := parseObjectID(clientID)
	if err != nil {
		return fmt.Errorf("invalid client ID")
	}

	return s.clientRepo.RevokeAccess(ctx, clientOID)
}

// ListShopClients returns all clients for a shop
func (s *ClientService) ListShopClients(ctx context.Context, shopID, userID string, page, pageSize int, status string) (*models.ShopClientsResponse, error) {
	shopOID, err := parseObjectID(shopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID")
	}

	shop, err := s.shopRepo.GetByID(ctx, shopOID)
	if err != nil {
		return nil, fmt.Errorf("shop not found")
	}

	if !s.hasManagePermission(shop, userID) {
		return nil, fmt.Errorf("permission denied")
	}

	clients, total, err := s.clientRepo.ListByShop(ctx, shopOID, page, pageSize, status)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.ShopClientsResponse{
		Clients:    clients,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetMyPrivateShops returns private shops where user has access
func (s *ClientService) GetMyPrivateShops(ctx context.Context, userID string) ([]models.Shop, error) {
	clients, err := s.clientRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var shops []models.Shop
	for _, c := range clients {
		shop, err := s.shopRepo.GetByID(ctx, c.ShopID)
		if err == nil && shop.Status == "active" {
			shops = append(shops, *shop)
		}
	}

	return shops, nil
}

// HasShopAccess checks if a user can access a private shop
func (s *ClientService) HasShopAccess(ctx context.Context, shopID primitive.ObjectID, userID string) bool {
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return false
	}

	// Public shops are always accessible
	if shop.IsPublic {
		return true
	}

	// Owner and managers always have access
	if shop.OwnerID == userID {
		return true
	}
	for _, m := range shop.Managers {
		if m.UserID == userID && m.Status == "active" {
			return true
		}
	}

	// Check if user is an authorized client
	return s.clientRepo.HasAccess(ctx, shopID, userID)
}

// verifyPIN calls auth-service to verify user PIN
func (s *ClientService) verifyPIN(userID, pin string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
		"pin":     pin,
	})

	resp, err := http.Post(
		s.authURL+"/api/v1/security/verify-pin",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return fmt.Errorf("failed to verify PIN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PIN verification failed: %s", string(body))
	}

	return nil
}

func (s *ClientService) hasManagePermission(shop *models.Shop, userID string) bool {
	if shop.OwnerID == userID {
		return true
	}
	for _, m := range shop.Managers {
		if m.UserID == userID && m.Status == "active" {
			if m.Role == "admin" || m.Role == "owner" {
				return true
			}
			for _, p := range m.Permissions {
				if p == "clients" || p == "all" {
					return true
				}
			}
		}
	}
	return false
}

func (s *ClientService) sendInvitationNotification(client *models.ShopClient, shop *models.Shop) {
	if s.kafkaClient == nil {
		return
	}

	event := messaging.NewEventEnvelope("shop_invitation", "shop-service", map[string]interface{}{
		"invitation_id": client.ID.Hex(),
		"shop_id":       shop.ID.Hex(),
		"shop_name":     shop.Name,
		"email":         client.Email,
		"title":         "Invitation à une boutique",
		"message":       fmt.Sprintf("Vous avez reçu une invitation pour accéder à la boutique %s", shop.Name),
	})

	if err := s.kafkaClient.Publish(context.Background(), "notification.send.email", event); err != nil {
		log.Printf("Failed to send invitation notification: %v", err)
	}
}

func (s *ClientService) sendAcceptedNotification(client *models.ShopClient, userID string) {
	if s.kafkaClient == nil {
		return
	}

	event := messaging.NewEventEnvelope("shop_invitation_accepted", "shop-service", map[string]interface{}{
		"user_id":   client.InvitedBy,
		"client_id": client.ID.Hex(),
		"shop_id":   client.ShopID.Hex(),
		"title":     "Invitation acceptée",
		"message":   fmt.Sprintf("%s a accepté votre invitation", client.Email),
	})

	if err := s.kafkaClient.Publish(context.Background(), "notification.send", event); err != nil {
		log.Printf("Failed to send accepted notification: %v", err)
	}
}
