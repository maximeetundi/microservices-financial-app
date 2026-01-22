package services

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	shopRepo     *repository.ShopRepository
	qrService    *QRService
	storage      *StorageService
}

func NewCategoryService(
	categoryRepo *repository.CategoryRepository,
	shopRepo *repository.ShopRepository,
	qrService *QRService,
	storage *StorageService,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		shopRepo:     shopRepo,
		qrService:    qrService,
		storage:      storage,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *models.CreateCategoryRequest, userID string) (*models.Category, error) {
	shopID, err := parseObjectID(req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return nil, fmt.Errorf("permission denied")
	}
	
	category := &models.Category{
		ShopID:       shopID,
		Name:         req.Name,
		Description:  req.Description,
		ImageURL:     req.ImageURL,
		Order:        req.Order,
		IsActive:     true,
		ProductCount: 0,
	}
	
	if req.ParentID != "" {
		parentID, err := parseObjectID(req.ParentID)
		if err == nil {
			category.ParentID = parentID
		}
	}
	
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	// Generate QR code
	if s.qrService != nil && s.storage != nil {
		qrData, err := s.qrService.GenerateCategoryQR(shop.Slug, category.Slug)
		if err == nil {
			qrURL, err := s.storage.UploadQRCode(ctx, qrData, "category", category.ID.Hex())
			if err == nil {
				s.categoryRepo.SetQRCode(ctx, category.ID, qrURL)
				category.QRCode = qrURL
			}
		}
	}
	
	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return nil, err
	}
	return s.categoryRepo.GetByID(ctx, oid)
}

func (s *CategoryService) ListByShop(ctx context.Context, shopSlug string) ([]models.Category, error) {
	shop, err := s.shopRepo.GetBySlug(ctx, shopSlug)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	return s.categoryRepo.ListByShop(ctx, shop.ID)
}

func (s *CategoryService) ListWithHierarchy(ctx context.Context, shopSlug string) ([]models.CategoryWithChildren, error) {
	shop, err := s.shopRepo.GetBySlug(ctx, shopSlug)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	categories, err := s.categoryRepo.ListByShop(ctx, shop.ID)
	if err != nil {
		return nil, err
	}
	
	return buildCategoryTree(categories), nil
}

func (s *CategoryService) Update(ctx context.Context, categoryID string, req *models.UpdateCategoryRequest, userID string) (*models.Category, error) {
	oid, err := parseObjectID(categoryID)
	if err != nil {
		return nil, err
	}
	
	category, err := s.categoryRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, category.ShopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return nil, fmt.Errorf("permission denied")
	}
	
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ImageURL != nil {
		category.ImageURL = *req.ImageURL
	}
	if req.Order != nil {
		category.Order = *req.Order
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, categoryID, userID string) error {
	oid, err := parseObjectID(categoryID)
	if err != nil {
		return err
	}
	
	category, err := s.categoryRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, category.ShopID)
	if err != nil {
		return fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return fmt.Errorf("permission denied")
	}
	
	return s.categoryRepo.Delete(ctx, oid)
}

func buildCategoryTree(categories []models.Category) []models.CategoryWithChildren {
	// Build a map of categories by ID
	categoryMap := make(map[string]*models.CategoryWithChildren)
	for i := range categories {
		cwc := models.CategoryWithChildren{
			Category: categories[i],
			Children: []models.CategoryWithChildren{},
		}
		categoryMap[categories[i].ID.Hex()] = &cwc
	}
	
	// Build the tree
	var roots []models.CategoryWithChildren
	for i := range categories {
		cat := &categories[i]
		if cat.ParentID.IsZero() {
			roots = append(roots, *categoryMap[cat.ID.Hex()])
		} else {
			parentID := cat.ParentID.Hex()
			if parent, ok := categoryMap[parentID]; ok {
				parent.Children = append(parent.Children, *categoryMap[cat.ID.Hex()])
			}
		}
	}
	
	return roots
}
