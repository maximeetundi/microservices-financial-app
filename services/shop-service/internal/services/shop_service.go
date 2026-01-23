package services

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
)

type ShopService struct {
	shopRepo     *repository.ShopRepository
	productRepo  *repository.ProductRepository
	walletClient *WalletClient
	qrService    *QRService
	storage      *StorageService
}

func NewShopService(
	shopRepo *repository.ShopRepository,
	productRepo *repository.ProductRepository,
	walletClient *WalletClient,
	qrService *QRService,
	storage *StorageService,
) *ShopService {
	return &ShopService{
		shopRepo:     shopRepo,
		productRepo:  productRepo,
		walletClient: walletClient,
		qrService:    qrService,
		storage:      storage,
	}
}

func (s *ShopService) Create(ctx context.Context, req *models.CreateShopRequest, userID, token string) (*models.Shop, error) {
	// Validate wallet ownership
	valid, err := s.walletClient.ValidateWalletOwnership(req.WalletID, userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate wallet: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("wallet does not belong to user")
	}
	
	ownerID := userID
	ownerType := "user"

	if req.EnterpriseID != "" && req.OwnerType == "enterprise" {
		ownerID = req.EnterpriseID
		ownerType = "enterprise"
	}

	shop := &models.Shop{
		OwnerID:     ownerID,
		OwnerType:   ownerType,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		WalletID:    req.WalletID,
		Currency:    req.Currency,
		Tags:        req.Tags,
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
		Status:      "active",
		Settings: models.ShopSettings{
			AllowPickup:      true,
			AllowDelivery:    true,
			AutoAcceptOrders: true,
		},
		Stats: models.ShopStats{},
		Managers: []models.ShopManager{
			{
				UserID:      userID,
				Role:        "owner",
				Permissions: []string{"all"},
				Status:      "active",
			},
		},
	}
	
	if req.Address != nil {
		shop.Address = *req.Address
	}
	if req.ContactInfo != nil {
		shop.ContactInfo = *req.ContactInfo
	}
	
	if err := s.shopRepo.Create(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}
	
	// Generate QR code
	if s.qrService != nil && s.storage != nil {
		qrData, err := s.qrService.GenerateShopQR(shop.Slug)
		if err == nil {
			qrURL, err := s.storage.UploadQRCode(ctx, qrData, "shop", shop.ID.Hex())
			if err == nil {
				s.shopRepo.SetQRCode(ctx, shop.ID, qrURL)
				shop.QRCode = qrURL
			}
		}
	}
	
	return shop, nil
}

func (s *ShopService) GetBySlug(ctx context.Context, slug string) (*models.Shop, error) {
	return s.shopRepo.GetBySlug(ctx, slug)
}

func (s *ShopService) GetByWalletID(ctx context.Context, walletID string) (*models.Shop, error) {
	return s.shopRepo.GetByWalletID(ctx, walletID)
}

func (s *ShopService) GetByID(ctx context.Context, id string) (*models.Shop, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return nil, err
	}
	return s.shopRepo.GetByID(ctx, oid)
}

func (s *ShopService) ListPublic(ctx context.Context, page, pageSize int, search string) (*models.ShopListResponse, error) {
	shops, total, err := s.shopRepo.ListPublic(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &models.ShopListResponse{
		Shops:      shops,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *ShopService) GetMyShops(ctx context.Context, userID string) ([]models.Shop, error) {
	return s.shopRepo.GetByManagerID(ctx, userID)
}

func (s *ShopService) Update(ctx context.Context, shopID string, req *models.UpdateShopRequest, userID, token string) (*models.Shop, error) {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return nil, err
	}
	
	shop, err := s.shopRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	// Check permission
	if !s.hasPermission(shop, userID, "settings") {
		return nil, fmt.Errorf("permission denied")
	}
	
	// Apply updates
	if req.Name != nil {
		shop.Name = *req.Name
	}
	if req.Description != nil {
		shop.Description = *req.Description
	}
	if req.IsPublic != nil {
		shop.IsPublic = *req.IsPublic
	}
	if req.WalletID != nil {
		// Validate new wallet
		valid, err := s.walletClient.ValidateWalletOwnership(*req.WalletID, shop.OwnerID, token)
		if err != nil || !valid {
			return nil, fmt.Errorf("invalid wallet")
		}
		shop.WalletID = *req.WalletID
	}
	if req.Tags != nil {
		shop.Tags = req.Tags
	}
	if req.Settings != nil {
		shop.Settings = *req.Settings
	}
	if req.Address != nil {
		shop.Address = *req.Address
	}
	if req.ContactInfo != nil {
		shop.ContactInfo = *req.ContactInfo
	}
	if req.LogoURL != nil {
		shop.LogoURL = *req.LogoURL
	}
	if req.BannerURL != nil {
		shop.BannerURL = *req.BannerURL
	}
	
	if err := s.shopRepo.Update(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}
	
	return shop, nil
}

func (s *ShopService) Delete(ctx context.Context, shopID, userID string) error {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return err
	}
	
	shop, err := s.shopRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("shop not found: %w", err)
	}
	
	// Only owner can delete
	if shop.OwnerID != userID {
		return fmt.Errorf("only owner can delete shop")
	}
	
	return s.shopRepo.Delete(ctx, oid)
}

func (s *ShopService) InviteManager(ctx context.Context, shopID string, req *models.InviteManagerRequest, userID string) error {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return err
	}
	
	shop, err := s.shopRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("shop not found: %w", err)
	}
	
	// Only owner can invite
	if shop.OwnerID != userID {
		return fmt.Errorf("only owner can invite managers")
	}
	
	manager := models.ShopManager{
		Email:       req.Email,
		Role:        req.Role,
		Permissions: req.Permissions,
		Status:      "pending",
	}
	
	return s.shopRepo.AddManager(ctx, oid, manager)
}

func (s *ShopService) RemoveManager(ctx context.Context, shopID, targetUserID, userID string) error {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return err
	}
	
	shop, err := s.shopRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("shop not found: %w", err)
	}
	
	// Only owner can remove
	if shop.OwnerID != userID {
		return fmt.Errorf("only owner can remove managers")
	}
	
	// Cannot remove owner
	if targetUserID == shop.OwnerID {
		return fmt.Errorf("cannot remove owner")
	}
	
	return s.shopRepo.RemoveManager(ctx, oid, targetUserID)
}

func (s *ShopService) UpdateStats(ctx context.Context, shopID string) error {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return err
	}
	
	productCount, _ := s.productRepo.CountByShop(ctx, oid)
	
	stats := models.ShopStats{
		TotalProducts: int(productCount),
	}
	
	return s.shopRepo.UpdateStats(ctx, oid, stats)
}

func (s *ShopService) hasPermission(shop *models.Shop, userID string, permission string) bool {
	if shop.OwnerID == userID {
		return true
	}
	
	for _, m := range shop.Managers {
		if m.UserID == userID && m.Status == "active" {
			if m.Role == "admin" {
				return true
			}
			for _, p := range m.Permissions {
				if p == permission || p == "all" {
					return true
				}
			}
		}
	}
	return false
}

func (s *ShopService) HasAccess(ctx context.Context, shopID, userID string) bool {
	oid, err := parseObjectID(shopID)
	if err != nil {
		return false
	}
	
	shop, err := s.shopRepo.GetByID(ctx, oid)
	if err != nil {
		return false
	}
	
	return s.hasPermission(shop, userID, "any")
}
