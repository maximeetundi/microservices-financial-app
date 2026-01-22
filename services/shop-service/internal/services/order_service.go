package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo      *repository.OrderRepository
	productRepo    *repository.ProductRepository
	shopRepo       *repository.ShopRepository
	walletClient   *WalletClient
	exchangeClient *ExchangeClient
	kafkaClient    *messaging.KafkaClient
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	shopRepo *repository.ShopRepository,
	walletClient *WalletClient,
	exchangeClient *ExchangeClient,
	kafkaClient *messaging.KafkaClient,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		shopRepo:       shopRepo,
		walletClient:   walletClient,
		exchangeClient: exchangeClient,
		kafkaClient:    kafkaClient,
	}
}

func (s *OrderService) Create(ctx context.Context, req *models.CreateOrderRequest, userID, userName, userEmail string) (*models.Order, error) {
	shopID, err := parseObjectID(req.ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid shop ID: %w", err)
	}
	
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found: %w", err)
	}
	
	// Validate and get buyer's wallet
	buyerWallet, err := s.walletClient.GetWallet(req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet: %w", err)
	}
	if buyerWallet.UserID != userID {
		return nil, fmt.Errorf("wallet does not belong to user")
	}
	
	// Build order items and calculate totals
	var items []models.OrderItem
	var subTotal float64
	
	for _, cartItem := range req.Items {
		productID, err := parseObjectID(cartItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}
		
		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", cartItem.ProductID)
		}
		
		if product.ShopID != shopID {
			return nil, fmt.Errorf("product does not belong to shop")
		}
		
		if product.Status != "active" {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}
		
		if product.Stock >= 0 && product.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s", product.Name)
		}
		
		itemTotal := product.Price * float64(cartItem.Quantity)
		
		imageURL := ""
		if len(product.Images) > 0 {
			imageURL = product.Images[0]
		}
		
		items = append(items, models.OrderItem{
			ProductID:    productID,
			ProductName:  product.Name,
			ProductImage: imageURL,
			VariantID:    cartItem.VariantID,
			Quantity:     cartItem.Quantity,
			UnitPrice:    product.Price,
			TotalPrice:   itemTotal,
			CustomValues: cartItem.CustomValues,
		})
		
		subTotal += itemTotal
	}
	
	// Calculate delivery fee
	deliveryFee := 0.0
	if req.DeliveryType == "delivery" && shop.Settings.AllowDelivery {
		deliveryFee = shop.Settings.DeliveryFee
	}
	
	totalAmount := subTotal + deliveryFee
	
	// Check min/max order amounts
	if shop.Settings.MinOrderAmount > 0 && totalAmount < shop.Settings.MinOrderAmount {
		return nil, fmt.Errorf("minimum order amount is %.2f %s", shop.Settings.MinOrderAmount, shop.Currency)
	}
	if shop.Settings.MaxOrderAmount > 0 && totalAmount > shop.Settings.MaxOrderAmount {
		return nil, fmt.Errorf("maximum order amount is %.2f %s", shop.Settings.MaxOrderAmount, shop.Currency)
	}
	
	// Handle currency conversion if needed
	var convertedAmount float64
	var exchangeRate float64
	buyerCurrency := buyerWallet.Currency
	
	if buyerCurrency != shop.Currency {
		converted, rate, err := s.exchangeClient.ConvertAmount(totalAmount, shop.Currency, buyerCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert currency: %w", err)
		}
		convertedAmount = converted
		exchangeRate = rate
	} else {
		convertedAmount = totalAmount
		exchangeRate = 1.0
	}
	
	// Check buyer has sufficient balance
	if buyerWallet.Balance < convertedAmount {
		return nil, fmt.Errorf("insufficient balance")
	}
	
	// Create order
	transactionID := primitive.NewObjectID().Hex()
	order := &models.Order{
		ShopID:            shopID,
		ShopName:          shop.Name,
		BuyerID:           userID,
		BuyerName:         userName,
		BuyerEmail:        userEmail,
		Items:             items,
		SubTotal:          subTotal,
		DeliveryFee:       deliveryFee,
		TotalAmount:       totalAmount,
		Currency:          shop.Currency,
		BuyerWalletID:     req.WalletID,
		BuyerCurrency:     buyerCurrency,
		ConvertedAmount:   convertedAmount,
		ExchangeRate:      exchangeRate,
		TransactionID:     transactionID,
		PaymentStatus:     "pending",
		OrderStatus:       "pending",
		DeliveryType:      req.DeliveryType,
		ShippingAddress:   req.ShippingAddress,
		Notes:             req.Notes,
	}
	
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	// Send payment request to Kafka
	paymentEvent := models.PaymentRequestEvent{
		TransactionID:       transactionID,
		SourceWalletID:      req.WalletID,
		UserID:              userID,
		Amount:              convertedAmount,
		Currency:            buyerCurrency,
		DestinationWalletID: shop.WalletID,
		DestinationUserID:   shop.OwnerID,
		Reference:           fmt.Sprintf("Order %s from %s", order.OrderNumber, shop.Name),
		OriginService:       "shop-service",
		MetaData: map[string]interface{}{
			"order_id":       order.ID.Hex(),
			"order_number":   order.OrderNumber,
			"shop_id":        shop.ID.Hex(),
			"shop_name":      shop.Name,
			"original_amount": totalAmount,
			"original_currency": shop.Currency,
		},
	}
	
	if s.kafkaClient != nil {
		if err := s.kafkaClient.PublishJSON("wallet.payment.request", paymentEvent); err != nil {
			log.Printf("Failed to publish payment request: %v", err)
		}
	}
	
	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, orderID, userID string) (*models.Order, error) {
	oid, err := parseObjectID(orderID)
	if err != nil {
		return nil, err
	}
	
	order, err := s.orderRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	
	// Check access (buyer or shop manager)
	if order.BuyerID != userID {
		shop, err := s.shopRepo.GetByID(ctx, order.ShopID)
		if err != nil || !hasShopPermission(shop, userID, "orders") {
			return nil, fmt.Errorf("access denied")
		}
	}
	
	return order, nil
}

func (s *OrderService) ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) (*models.OrderListResponse, error) {
	orders, total, err := s.orderRepo.ListByBuyer(ctx, buyerID, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *OrderService) ListByShop(ctx context.Context, shopID, userID string, page, pageSize int, status string) (*models.OrderListResponse, error) {
	shopOID, err := parseObjectID(shopID)
	if err != nil {
		return nil, err
	}
	
	shop, err := s.shopRepo.GetByID(ctx, shopOID)
	if err != nil {
		return nil, fmt.Errorf("shop not found")
	}
	
	if !hasShopPermission(shop, userID, "orders") {
		return nil, fmt.Errorf("access denied")
	}
	
	orders, total, err := s.orderRepo.ListByShop(ctx, shopOID, page, pageSize, status)
	if err != nil {
		return nil, err
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID string, req *models.UpdateOrderStatusRequest, userID string) (*models.Order, error) {
	oid, err := parseObjectID(orderID)
	if err != nil {
		return nil, err
	}
	
	order, err := s.orderRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}
	
	shop, err := s.shopRepo.GetByID(ctx, order.ShopID)
	if err != nil {
		return nil, fmt.Errorf("shop not found")
	}
	
	if !hasShopPermission(shop, userID, "orders") {
		return nil, fmt.Errorf("access denied")
	}
	
	order.OrderStatus = req.Status
	if req.TrackingNumber != "" {
		order.TrackingNumber = req.TrackingNumber
	}
	if req.SellerNotes != "" {
		order.SellerNotes = req.SellerNotes
	}
	if req.Status == "delivered" {
		now := time.Now()
		order.CompletedAt = &now
	}
	
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}
	
	// Send notification to buyer
	if s.kafkaClient != nil {
		notification := models.NotificationEvent{
			UserID:    order.BuyerID,
			Type:      "order_status",
			Title:     "Mise à jour de commande",
			Message:   fmt.Sprintf("Votre commande %s est maintenant: %s", order.OrderNumber, req.Status),
			Data:      map[string]interface{}{"order_id": order.ID.Hex()},
			Timestamp: time.Now(),
		}
		s.kafkaClient.PublishJSON("notification.send", notification)
	}
	
	return order, nil
}

func (s *OrderService) Refund(ctx context.Context, orderID string, req *models.RefundOrderRequest, userID string) error {
	oid, err := parseObjectID(orderID)
	if err != nil {
		return err
	}
	
	order, err := s.orderRepo.GetByID(ctx, oid)
	if err != nil {
		return fmt.Errorf("order not found")
	}
	
	shop, err := s.shopRepo.GetByID(ctx, order.ShopID)
	if err != nil {
		return fmt.Errorf("shop not found")
	}
	
	if !hasShopPermission(shop, userID, "orders") {
		return fmt.Errorf("access denied")
	}
	
	if order.PaymentStatus != "completed" {
		return fmt.Errorf("cannot refund unpaid order")
	}
	
	if order.PaymentStatus == "refunded" {
		return fmt.Errorf("order already refunded")
	}
	
	// Send refund request to Kafka
	refundEvent := models.PaymentRequestEvent{
		TransactionID:       primitive.NewObjectID().Hex(),
		SourceWalletID:      shop.WalletID,
		UserID:              shop.OwnerID,
		Amount:              order.ConvertedAmount,
		Currency:            order.BuyerCurrency,
		DestinationWalletID: order.BuyerWalletID,
		DestinationUserID:   order.BuyerID,
		Reference:           fmt.Sprintf("Refund for order %s", order.OrderNumber),
		OriginService:       "shop-service",
		MetaData: map[string]interface{}{
			"order_id":     order.ID.Hex(),
			"order_number": order.OrderNumber,
			"refund":       true,
			"reason":       req.Reason,
		},
	}
	
	if s.kafkaClient != nil {
		if err := s.kafkaClient.PublishJSON("wallet.payment.request", refundEvent); err != nil {
			log.Printf("Failed to publish refund request: %v", err)
			return fmt.Errorf("failed to initiate refund")
		}
	}
	
	// Mark order as refunded
	return s.orderRepo.MarkRefunded(ctx, oid, req.Reason)
}

// HandlePaymentStatus processes payment status updates from Kafka
func (s *OrderService) HandlePaymentStatus(ctx context.Context, event *models.PaymentStatusEvent) error {
	order, err := s.orderRepo.GetByTransactionID(ctx, event.TransactionID)
	if err != nil {
		// Not our transaction
		return nil
	}
	
	if err := s.orderRepo.UpdatePaymentStatus(ctx, event.TransactionID, event.Status); err != nil {
		return err
	}
	
	if event.Status == "completed" {
		// Update product stock
		for _, item := range order.Items {
			s.productRepo.IncrementSoldCount(ctx, item.ProductID, item.Quantity)
		}
		
		// Update shop stats
		orderCount, _ := s.orderRepo.CountByShop(ctx, order.ShopID)
		revenue, _ := s.orderRepo.SumRevenueByShop(ctx, order.ShopID)
		productCount, _ := s.productRepo.CountByShop(ctx, order.ShopID)
		
		shop, _ := s.shopRepo.GetByID(ctx, order.ShopID)
		if shop != nil {
			shop.Stats.TotalOrders = int(orderCount)
			shop.Stats.TotalRevenue = revenue
			shop.Stats.TotalProducts = int(productCount)
			s.shopRepo.Update(ctx, shop)
		}
		
		// Notify buyer
		if s.kafkaClient != nil {
			notification := models.NotificationEvent{
				UserID:    order.BuyerID,
				Type:      "order_paid",
				Title:     "Paiement confirmé",
				Message:   fmt.Sprintf("Votre paiement pour la commande %s a été confirmé", order.OrderNumber),
				Data:      map[string]interface{}{"order_id": order.ID.Hex()},
				Timestamp: time.Now(),
			}
			s.kafkaClient.PublishJSON("notification.send", notification)
		}
	}
	
	return nil
}
