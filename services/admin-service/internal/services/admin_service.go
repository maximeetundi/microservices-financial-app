package services

import (
	"context"
	"errors"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo        *repository.AdminRepository
	mongoRepo   *repository.MongoRepository
	kafkaClient *messaging.KafkaClient
	config      *config.Config
}

func NewAdminService(repo *repository.AdminRepository, mongoRepo *repository.MongoRepository, kafkaClient *messaging.KafkaClient, cfg *config.Config) *AdminService {
	return &AdminService{
		repo:        repo,
		mongoRepo:   mongoRepo,
		kafkaClient: kafkaClient,
		config:      cfg,
	}
}

// ========== Authentication ==========

func (s *AdminService) Login(email, password string) (*models.AdminLoginResponse, error) {
	admin, err := s.repo.GetAdminByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	if !admin.IsActive {
		return nil, errors.New("account is disabled")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	
	// Get permissions
	permissions, _ := s.repo.GetAdminPermissions(admin.ID)
	
	// Generate JWT
	expiresAt := time.Now().Add(8 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id":    admin.ID,
		"email":       admin.Email,
		"role_id":     admin.RoleID,
		"permissions": permissions,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
		"type":        "admin",
	})
	
	tokenString, err := token.SignedString([]byte(s.config.AdminJWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	
	// Update last login
	s.repo.UpdateLastLogin(admin.ID)
	
	// Get role details
	role, _ := s.repo.GetRoleByID(admin.RoleID)
	admin.Role = role
	admin.PasswordHash = "" // Don't expose
	
	return &models.AdminLoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		Admin:     admin,
	}, nil
}

func (s *AdminService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.AdminJWTSecret), nil
	})
	
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	
	// Check if it's an admin token
	if claims["type"] != "admin" {
		return nil, errors.New("not an admin token")
	}
	
	return &claims, nil
}

// ========== Admin CRUD ==========

func (s *AdminService) CreateAdmin(req *models.CreateAdminRequest, createdBy string) (*models.AdminUser, error) {
	// Check if email exists
	existing, _ := s.repo.GetAdminByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	admin := models.NewAdminUser(req.Email, req.FirstName, req.LastName, req.RoleID, createdBy)
	admin.PasswordHash = string(hashedPassword)
	
	if err := s.repo.CreateAdminUser(admin); err != nil {
		return nil, err
	}
	
	admin.PasswordHash = ""
	return admin, nil
}

func (s *AdminService) GetAdmin(id string) (*models.AdminUser, error) {
	admin, err := s.repo.GetAdminByID(id)
	if err != nil {
		return nil, err
	}
	
	role, _ := s.repo.GetRoleByID(admin.RoleID)
	admin.Role = role
	admin.PasswordHash = ""
	
	return admin, nil
}

func (s *AdminService) GetAllAdmins(limit, offset int) ([]models.AdminUser, error) {
	return s.repo.GetAllAdmins(limit, offset)
}

func (s *AdminService) UpdateAdmin(id string, req *models.UpdateAdminRequest) error {
	updates := make(map[string]interface{})
	
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	
	if len(updates) == 0 {
		return nil
	}
	
	return s.repo.UpdateAdmin(id, updates)
}

func (s *AdminService) DeleteAdmin(id string) error {
	return s.repo.DeleteAdmin(id)
}

// ========== Roles ==========

func (s *AdminService) GetRoles() ([]models.AdminRole, error) {
	return s.repo.GetRoles()
}

func (s *AdminService) GetRole(id string) (*models.AdminRole, error) {
	return s.repo.GetRoleByID(id)
}

func (s *AdminService) GetAdminPermissions(adminID string) ([]string, error) {
	return s.repo.GetAdminPermissions(adminID)
}

// ========== Dashboard ==========

func (s *AdminService) GetDashboardStats() (map[string]interface{}, error) {
	return s.repo.GetDashboardStats()
}

// ========== Main DB Data ==========

func (s *AdminService) GetUsers(limit, offset int) ([]map[string]interface{}, error) {
	return s.repo.GetUsersFromMainDB(limit, offset)
}

func (s *AdminService) GetTransactions(limit, offset int) ([]map[string]interface{}, error) {
	return s.repo.GetTransactionsFromMainDB(limit, offset)
}

func (s *AdminService) GetCards(limit, offset int) ([]map[string]interface{}, error) {
	return s.repo.GetCardsFromMainDB(limit, offset)
}

func (s *AdminService) GetWallets(limit, offset int) ([]map[string]interface{}, error) {
	return s.repo.GetWalletsFromMainDB(limit, offset)
}

func (s *AdminService) GetUserKYCDocuments(userID string) ([]map[string]interface{}, error) {
	return s.repo.GetUserKYCDocuments(userID)
}

// ========== Mongo Data (Read-Only) ==========

func (s *AdminService) GetCampaigns(limit, offset int) ([]map[string]interface{}, error) {
	if s.mongoRepo == nil {
		return nil, errors.New("mongo repository not initialized")
	}
	return s.mongoRepo.GetCampaigns(limit, offset)
}

func (s *AdminService) GetDonations(limit, offset int) ([]map[string]interface{}, error) {
	if s.mongoRepo == nil {
		return nil, errors.New("mongo repository not initialized")
	}
	return s.mongoRepo.GetDonations(limit, offset)
}

func (s *AdminService) GetEvents(limit, offset int) ([]map[string]interface{}, error) {
	if s.mongoRepo == nil {
		return nil, errors.New("mongo repository not initialized")
	}
	return s.mongoRepo.GetEvents(limit, offset)
}

func (s *AdminService) GetSoldTickets(limit, offset int) ([]map[string]interface{}, error) {
	if s.mongoRepo == nil {
		return nil, errors.New("mongo repository not initialized")
	}
	return s.mongoRepo.GetSoldTickets(limit, offset)
}

// ========== Admin Actions via Kafka ==========

func (s *AdminService) BlockUser(userID, reason, adminID string) error {
	// Get user email for notification
	user, _ := s.repo.GetUserByID(userID)
	
	// Direct database update for immediate effect
	if err := s.repo.BlockUser(userID); err != nil {
		return err
	}
	
	// Publish notification to user
	if user != nil {
		userNotif := map[string]interface{}{
			"user_id": userID,
			"email":   user["email"],
			"phone":   user["phone"],
			"reason":  reason,
		}
		s.publishEvent("user.blocked", userNotif)
		
		// Also publish admin notification
		adminNotif := map[string]interface{}{
			"user_id":    userID,
			"user_email": user["email"],
			"admin_id":   adminID,
			"reason":     reason,
		}
		s.publishEvent("admin.user_blocked", adminNotif)
	}
	
	// Also publish to Kafka for other services
	cmd := map[string]interface{}{
		"action":    "block_user",
		"user_id":   userID,
		"reason":    reason,
		"admin_id":  adminID,
	}
	return s.publishCommand("user.block", cmd)
}

func (s *AdminService) UnblockUser(userID, adminID string) error {
	// Get user email for notification
	user, _ := s.repo.GetUserByID(userID)
	
	// Direct database update for immediate effect
	if err := s.repo.UnblockUser(userID); err != nil {
		return err
	}
	
	// Publish notification to user
	if user != nil {
		userNotif := map[string]interface{}{
			"user_id": userID,
			"email":   user["email"],
			"phone":   user["phone"],
		}
		s.publishEvent("user.unblocked", userNotif)
		
		// Also publish admin notification
		adminNotif := map[string]interface{}{
			"user_id":    userID,
			"user_email": user["email"],
			"admin_id":   adminID,
		}
		s.publishEvent("admin.user_unblocked", adminNotif)
	}
	
	// Also publish to Kafka for other services
	cmd := map[string]interface{}{
		"action":   "unblock_user",
		"user_id":  userID,
		"admin_id": adminID,
	}
	return s.publishCommand("user.unblock", cmd)
}

func (s *AdminService) ApproveKYC(userID, level, adminID string) error {
	// Direct database update for immediate effect
	kycLevel := 2 // Default to level 2 for verified
	if level == "basic" {
		kycLevel = 1
	} else if level == "full" {
		kycLevel = 3
	}
	
	if err := s.repo.UpdateUserKYCStatus(userID, "verified", kycLevel); err != nil {
		return err
	}
	
	// Also publish to Kafka for other services
	cmd := map[string]interface{}{
		"action":   "approve_kyc",
		"user_id":  userID,
		"level":    level,
		"admin_id": adminID,
	}
	return s.publishCommand("kyc.approve", cmd)
}

func (s *AdminService) RejectKYC(userID, reason, adminID string) error {
	// Direct database update for immediate effect
	if err := s.repo.UpdateUserKYCStatus(userID, "rejected", 0); err != nil {
		return err
	}
	
	// Also publish to Kafka for other services
	cmd := map[string]interface{}{
		"action":   "reject_kyc",
		"user_id":  userID,
		"reason":   reason,
		"admin_id": adminID,
	}
	return s.publishCommand("kyc.reject", cmd)
}

func (s *AdminService) BlockTransaction(txID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":         "block_transaction",
		"transaction_id": txID,
		"reason":         reason,
		"admin_id":       adminID,
	}
	return s.publishCommand("transaction.block", cmd)
}

func (s *AdminService) RefundTransaction(txID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":         "refund_transaction",
		"transaction_id": txID,
		"reason":         reason,
		"admin_id":       adminID,
	}
	return s.publishCommand("transaction.refund", cmd)
}

func (s *AdminService) FreezeCard(cardID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":   "freeze_card",
		"card_id":  cardID,
		"reason":   reason,
		"admin_id": adminID,
	}
	return s.publishCommand("card.freeze", cmd)
}

func (s *AdminService) BlockCard(cardID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":   "block_card",
		"card_id":  cardID,
		"reason":   reason,
		"admin_id": adminID,
	}
	return s.publishCommand("card.block", cmd)
}

func (s *AdminService) FreezeWallet(walletID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":    "freeze_wallet",
		"wallet_id": walletID,
		"reason":    reason,
		"admin_id":  adminID,
	}
	return s.publishCommand("wallet.freeze", cmd)
}

func (s *AdminService) BlockTransfer(transferID, reason, adminID string) error {
	cmd := map[string]interface{}{
		"action":      "block_transfer",
		"transfer_id": transferID,
		"reason":      reason,
		"admin_id":    adminID,
	}
	return s.publishCommand("transfer.block", cmd)
}

func (s *AdminService) publishCommand(eventType string, cmd map[string]interface{}) error {
	if s.kafkaClient == nil {
		return nil
	}
	
	envelope := messaging.NewEventEnvelope(eventType, "admin-service", cmd)
	return s.kafkaClient.Publish(context.Background(), messaging.TopicUserEvents, envelope)
}

func (s *AdminService) publishEvent(eventType string, data map[string]interface{}) {
	if s.kafkaClient == nil {
		return
	}
	
	envelope := messaging.NewEventEnvelope(eventType, "admin-service", data)
	s.kafkaClient.Publish(context.Background(), messaging.TopicUserEvents, envelope)
}

// ========== Audit Logs ==========

func (s *AdminService) CreateAuditLog(adminID, adminEmail, action, resource, resourceID, details, ipAddress, userAgent string) {
	log := &models.AdminAuditLog{
		AdminID:    adminID,
		AdminEmail: adminEmail,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	s.repo.CreateAuditLog(log)
}

func (s *AdminService) GetAuditLogs(limit, offset int) ([]models.AdminAuditLog, error) {
	return s.repo.GetAuditLogs(limit, offset, nil)
}
