package services

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	shopRepo     *repository.ShopRepository
	qrService    *QRService
	storage      *StorageService
}

func NewProductService(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	shopRepo *repository.ShopRepository,
	qrService *QRService,
	storage *StorageService,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		shopRepo:     shopRepo,
		qrService:    qrService,
		storage:      storage,
	}
}

func (s *ProductService) Create(ctx context.Context, req *models.CreateProductRequest, userID string) (*models.Product, error) {
	shopID, err := parseObjectID(req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}
	
	// Get shop and verify permission
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return nil, fmt.Errorf("permission denied")
	}
	
	product := &models.Product{
		ShopID:         shopID,
		Name:           req.Name,
		Description:    req.Description,
		ShortDesc:      req.ShortDesc,
		Price:          req.Price,
		CompareAtPrice: req.CompareAtPrice,
		Currency:       shop.Currency,
		Images:         req.Images,
		IsDigital:      req.IsDigital,
		DigitalFileURL: req.DigitalFileURL,
		LicenseText:    req.LicenseText,
		Stock:          req.Stock,
		SKU:            req.SKU,
		Weight:         req.Weight,
		IsCustomizable: req.IsCustomizable,
		CustomFields:   req.CustomFields,
		Tags:           req.Tags,
		Status:         "active",
		IsFeatured:     req.IsFeatured,
	}
	
	if req.Status != "" {
		product.Status = req.Status
	}
	
	if req.CategoryID != "" {
		catID, err := parseObjectID(req.CategoryID)
		if err == nil {
			product.CategoryID = catID
		}
	}
	
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	// Generate QR code
	if s.qrService != nil && s.storage != nil {
		qrData, err := s.qrService.GenerateProductQR(shop.Slug, product.Slug)
		if err == nil {
			qrURL, err := s.storage.UploadQRCode(ctx, qrData, "product", product.ID.Hex())
			if err == nil {
				s.productRepo.SetQRCode(ctx, product.ID, qrURL)
				product.QRCode = qrURL
			}
		}
	}
	
	// Update category product count
	if !product.CategoryID.IsZero() {
		count, _ := s.productRepo.CountByCategory(ctx, product.CategoryID)
		s.categoryRepo.UpdateProductCount(ctx, product.CategoryID, int(count))
	}
	
	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return nil, err
	}
	
	product, err := s.productRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	
	// Increment view count
	s.productRepo.IncrementViewCount(ctx, oid)
	
	return product, nil
}

func (s *ProductService) GetBySlug(ctx context.Context, shopSlug, productSlug string) (*models.Product, error) {
	shop, err := s.shopRepo.GetBySlug(ctx, shopSlug)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	product, err := s.productRepo.GetBySlug(ctx, shop.ID, productSlug)
	if err != nil {
		return nil, err
	}
	
	// Increment view count
	s.productRepo.IncrementViewCount(ctx, product.ID)
	
	return product, nil
}

func (s *ProductService) ListByShop(ctx context.Context, shopSlug string, page, pageSize int, categorySlug, status, search string) (*models.ProductListResponse, error) {
	shop, err := s.shopRepo.GetBySlug(ctx, shopSlug)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	var categoryID *primitive.ObjectID
	if categorySlug != "" {
		cat, err := s.categoryRepo.GetBySlug(ctx, shop.ID, categorySlug)
		if err == nil {
			categoryID = &cat.ID
		}
	}
	
	products, total, err := s.productRepo.ListByShop(ctx, shop.ID, page, pageSize, categoryID, status, search)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &models.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *ProductService) Update(ctx context.Context, productID string, req *models.UpdateProductRequest, userID string) (*models.Product, error) {
	oid, err := parseObjectID(productID)
	if err != nil {
		return nil, err
	}
	
	product, err := s.productRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, product.ShopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return nil, fmt.Errorf("permission denied")
	}
	
	oldCategoryID := product.CategoryID
	
	// Apply updates
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.ShortDesc != nil {
		product.ShortDesc = *req.ShortDesc
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CompareAtPrice != nil {
		product.CompareAtPrice = *req.CompareAtPrice
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.IsCustomizable != nil {
		product.IsCustomizable = *req.IsCustomizable
	}
	if req.CustomFields != nil {
		product.CustomFields = req.CustomFields
	}
	if req.Tags != nil {
		product.Tags = req.Tags
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
	}
	if req.DigitalFileURL != nil {
		product.DigitalFileURL = *req.DigitalFileURL
	}
	if req.LicenseText != nil {
		product.LicenseText = *req.LicenseText
	}
	if req.CategoryID != nil {
		if *req.CategoryID == "" {
			product.CategoryID = primitive.NilObjectID
		} else {
			catID, err := parseObjectID(*req.CategoryID)
			if err == nil {
				product.CategoryID = catID
			}
		}
	}
	
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	// Update category counts if category changed
	if oldCategoryID != product.CategoryID {
		if !oldCategoryID.IsZero() {
			count, _ := s.productRepo.CountByCategory(ctx, oldCategoryID)
			s.categoryRepo.UpdateProductCount(ctx, oldCategoryID, int(count))
		}
		if !product.CategoryID.IsZero() {
			count, _ := s.productRepo.CountByCategory(ctx, product.CategoryID)
			s.categoryRepo.UpdateProductCount(ctx, product.CategoryID, int(count))
		}
	}
	
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, productID, userID string) error {
	oid, err := parseObjectID(productID)
	if err != nil {
		return err
	}
	
	product, err := s.productRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, product.ShopID)
	if err != nil {
		return fmt.Errorf("shop not found: %w", err)
	}
	
	if !hasShopPermission(shop, userID, "products") {
		return fmt.Errorf("permission denied")
	}
	
	if err := s.productRepo.Delete(ctx, oid); err != nil {
		return err
	}
	
	// Update category count
	if !product.CategoryID.IsZero() {
		count, _ := s.productRepo.CountByCategory(ctx, product.CategoryID)
		s.categoryRepo.UpdateProductCount(ctx, product.CategoryID, int(count))
	}
	
	return nil
}

func hasShopPermission(shop *models.Shop, userID string, permission string) bool {
	if shop.OwnerID == userID {
		return true
	}
	
	for _, m := range shop.Managers {
		if m.UserID == userID && m.Status == "active" {
			if m.Role == "admin" || m.Role == "owner" {
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
