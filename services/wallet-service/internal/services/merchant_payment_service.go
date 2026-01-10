package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
)

// Default expiration time in minutes
const DefaultExpirationMinutes = 60 // 1 hour

// MerchantPaymentService handles merchant payment requests
type MerchantPaymentService struct {
	paymentRepo   *repository.PaymentRequestRepository
	walletService *WalletService
	feeService    *FeeService
	config        *config.Config
	kafkaClient   *messaging.KafkaClient
	baseURL       string
}

// NewMerchantPaymentService creates a new merchant payment service
func NewMerchantPaymentService(
	paymentRepo *repository.PaymentRequestRepository,
	walletService *WalletService,
	feeService *FeeService,
	cfg *config.Config,
	kafkaClient *messaging.KafkaClient,
) *MerchantPaymentService {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://app.zekora.com"
	}
	
	return &MerchantPaymentService{
		paymentRepo:   paymentRepo,
		walletService: walletService,
		feeService:    feeService,
		config:        cfg,
		kafkaClient:   kafkaClient,
		baseURL:       baseURL,
	}
}

// CreatePaymentRequest creates a new payment request with QR code
func (s *MerchantPaymentService) CreatePaymentRequest(merchantID string, req *models.CreatePaymentRequestDTO) (*models.PaymentRequestResponse, error) {
	// Validate wallet ownership
	wallet, err := s.walletService.GetWallet(req.WalletID, merchantID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found or not owned by merchant: %w", err)
	}

	// Validate based on type
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Generate unique ID
	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixNano())

	// Calculate expiration
	var expiresAt *time.Time
	neverExpires := false
	
	if req.ExpiresInMinutes != nil {
		if *req.ExpiresInMinutes == -1 {
			// Never expires
			neverExpires = true
		} else if *req.ExpiresInMinutes > 0 {
			exp := time.Now().Add(time.Duration(*req.ExpiresInMinutes) * time.Minute)
			expiresAt = &exp
		} else {
			// Use default
			exp := time.Now().Add(DefaultExpirationMinutes * time.Minute)
			expiresAt = &exp
		}
	} else {
		// Use default expiration
		exp := time.Now().Add(DefaultExpirationMinutes * time.Minute)
		expiresAt = &exp
	}

	// Calculate total for invoice type
	var amount *float64
	if req.Type == models.PaymentTypeInvoice && len(req.Items) > 0 {
		total := 0.0
		for i := range req.Items {
			req.Items[i].TotalPrice = float64(req.Items[i].Quantity) * req.Items[i].UnitPrice
			total += req.Items[i].TotalPrice
		}
		amount = &total
	} else {
		amount = req.Amount
	}

	// Generate payment link
	paymentLink := fmt.Sprintf("%s/pay/%s", s.baseURL, paymentID)

	// Create QR code data
	qrData := models.QRCodeData{
		Type:         "zekora_payment",
		Version:      1,
		PaymentID:    paymentID,
		MerchantName: merchantID, // TODO: Get merchant name from user service
		Amount:       amount,
		Currency:     req.Currency,
		Title:        req.Title,
		URL:          paymentLink,
	}
	
	qrDataJSON, _ := json.Marshal(qrData)

	// Create payment request
	paymentRequest := &models.PaymentRequest{
		ID:           paymentID,
		MerchantID:   merchantID,
		WalletID:     req.WalletID,
		Type:         req.Type,
		Amount:       amount,
		MinAmount:    req.MinAmount,
		MaxAmount:    req.MaxAmount,
		Currency:     wallet.Currency,
		Title:        req.Title,
		Description:  req.Description,
		Items:        req.Items,
		QRCodeData:   string(qrDataJSON),
		PaymentLink:  paymentLink,
		ExpiresAt:    expiresAt,
		NeverExpires: neverExpires,
		Status:       models.PaymentStatusPending,
		Reusable:     req.Reusable,
		TimesUsed:    0,
		Metadata:     req.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := s.paymentRepo.Create(paymentRequest); err != nil {
		return nil, fmt.Errorf("failed to create payment request: %w", err)
	}

	// Generate QR code image
	qrCodeBase64, err := s.generateQRCode(string(qrDataJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return &models.PaymentRequestResponse{
		PaymentRequest: paymentRequest,
		QRCodeBase64:   qrCodeBase64,
		PaymentURL:     paymentLink,
	}, nil
}

// GetPaymentRequest gets a payment request by ID
func (s *MerchantPaymentService) GetPaymentRequest(paymentID string) (*models.PaymentRequest, error) {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Check expiration
	if payment.ExpiresAt != nil && !payment.NeverExpires && time.Now().After(*payment.ExpiresAt) {
		if payment.Status == models.PaymentStatusPending {
			payment.Status = models.PaymentStatusExpired
			s.paymentRepo.UpdateStatus(paymentID, string(models.PaymentStatusExpired))
		}
	}

	return payment, nil
}

// PayPaymentRequest processes a payment for a payment request
func (s *MerchantPaymentService) PayPaymentRequest(customerID string, req *models.PayPaymentRequestDTO) (*models.PaymentHistory, error) {
	// Get payment request
	payment, err := s.GetPaymentRequest(req.PaymentRequestID)
	if err != nil {
		return nil, fmt.Errorf("payment request not found: %w", err)
	}

	// Validate status
	if payment.Status == models.PaymentStatusExpired {
		return nil, fmt.Errorf("payment request has expired")
	}
	if payment.Status == models.PaymentStatusCancelled {
		return nil, fmt.Errorf("payment request has been cancelled")
	}
	if payment.Status == models.PaymentStatusPaid && !payment.Reusable {
		return nil, fmt.Errorf("payment request has already been paid")
	}

	// Determine amount to pay
	var amountToPay float64
	if payment.Type == models.PaymentTypeVariable {
		if req.Amount <= 0 {
			return nil, fmt.Errorf("amount is required for variable payment")
		}
		if payment.MinAmount != nil && req.Amount < *payment.MinAmount {
			return nil, fmt.Errorf("amount is below minimum: %.2f", *payment.MinAmount)
		}
		if payment.MaxAmount != nil && req.Amount > *payment.MaxAmount {
			return nil, fmt.Errorf("amount exceeds maximum: %.2f", *payment.MaxAmount)
		}
		amountToPay = req.Amount
	} else {
		if payment.Amount == nil {
			return nil, fmt.Errorf("payment amount not set")
		}
		amountToPay = *payment.Amount
	}

	// Calculate merchant fee
	fee, err := s.feeService.CalculateFee("merchant_payment", amountToPay)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate merchant fee: %w", err)
	}
	
	netAmount := amountToPay - fee

	// Perform wallet transfer
	transferReq := &models.TransferRequest{
		FromWalletID: req.FromWalletID,
		ToWalletID:   payment.WalletID,
		Amount:       amountToPay,
		Currency:     payment.Currency,
		Description:  fmt.Sprintf("Payment: %s", payment.Title),
	}

	transaction, err := s.walletService.Transfer(customerID, transferReq)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	// Update payment request
	now := time.Now()
	payment.PaidAmount = &amountToPay
	payment.PaidBy = &customerID
	payment.PaidAt = &now
	payment.TransactionID = &transaction.ID
	payment.TimesUsed++
	
	if !payment.Reusable {
		payment.Status = models.PaymentStatusPaid
	}
	
	s.paymentRepo.Update(payment)

	// Create payment history
	history := &models.PaymentHistory{
		ID:               fmt.Sprintf("hist_%d", time.Now().UnixNano()),
		PaymentRequestID: payment.ID,
		MerchantID:       payment.MerchantID,
		CustomerID:       customerID,
		Amount:           amountToPay,
		Fee:              fee,
		NetAmount:        netAmount,
		Currency:         payment.Currency,
		TransactionID:    transaction.ID,
		PaidAt:           now,
	}

	s.paymentRepo.CreateHistory(history)

	// Publish event for notifications
	s.publishPaymentEvent("payment.completed", payment, customerID, amountToPay)

	return history, nil
}

// GetMerchantPayments gets all payment requests for a merchant
func (s *MerchantPaymentService) GetMerchantPayments(merchantID string, limit, offset int) ([]models.PaymentRequest, error) {
	return s.paymentRepo.GetByMerchantID(merchantID, limit, offset)
}

// GetPaymentHistory gets payment history for a merchant
func (s *MerchantPaymentService) GetPaymentHistory(merchantID string, limit, offset int) ([]models.PaymentHistory, error) {
	return s.paymentRepo.GetHistoryByMerchantID(merchantID, limit, offset)
}

// CancelPaymentRequest cancels a payment request
func (s *MerchantPaymentService) CancelPaymentRequest(merchantID, paymentID string) error {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return err
	}

	if payment.MerchantID != merchantID {
		return fmt.Errorf("not authorized to cancel this payment")
	}

	if payment.Status != models.PaymentStatusPending {
		return fmt.Errorf("cannot cancel payment with status: %s", payment.Status)
	}

	return s.paymentRepo.UpdateStatus(paymentID, string(models.PaymentStatusCancelled))
}

// RegenerateQRCode regenerates the QR code for a payment request
func (s *MerchantPaymentService) RegenerateQRCode(paymentID string) (string, error) {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return "", err
	}

	return s.generateQRCode(payment.QRCodeData)
}

// ===== PRIVATE METHODS =====

func (s *MerchantPaymentService) validateRequest(req *models.CreatePaymentRequestDTO) error {
	switch req.Type {
	case models.PaymentTypeFixed:
		if req.Amount == nil || *req.Amount <= 0 {
			return fmt.Errorf("amount is required for fixed payment type")
		}
	case models.PaymentTypeVariable:
		// Amount is optional, customer will define it
		if req.MinAmount != nil && req.MaxAmount != nil && *req.MinAmount > *req.MaxAmount {
			return fmt.Errorf("min_amount cannot be greater than max_amount")
		}
	case models.PaymentTypeInvoice:
		if len(req.Items) == 0 {
			return fmt.Errorf("items are required for invoice payment type")
		}
	default:
		return fmt.Errorf("invalid payment type: %s", req.Type)
	}

	return nil
}

func (s *MerchantPaymentService) generateQRCode(data string) (string, error) {
	// Create QR code
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	// Scale QR code to 256x256
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		return "", err
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, qrCode); err != nil {
		return "", err
	}

	// Convert to base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	
	return "data:image/png;base64," + base64Str, nil
}

func (s *MerchantPaymentService) publishPaymentEvent(eventType string, payment *models.PaymentRequest, customerID string, amount float64) {
	if s.kafkaClient == nil {
		return
	}

	// Notification for Customer (payer)
	customerEventData := map[string]interface{}{
		"user_id":     customerID,
		"payment_id":  payment.ID,
		"merchant_id": payment.MerchantID,
		"amount":      amount,
		"currency":    payment.Currency,
		"title":       payment.Title,
	}
	customerEnvelope := messaging.NewEventEnvelope("payment.sent", "wallet-service", customerEventData)
	s.kafkaClient.Publish(context.Background(), messaging.TopicWalletEvents, customerEnvelope)

	// Notification for Merchant (receiver)
	merchantEventData := map[string]interface{}{
		"user_id":     payment.MerchantID,
		"payment_id":  payment.ID,
		"customer_id": customerID,
		"amount":      amount,
		"currency":    payment.Currency,
		"title":       payment.Title,
	}
	merchantEnvelope := messaging.NewEventEnvelope("payment.received", "wallet-service", merchantEventData)
	s.kafkaClient.Publish(context.Background(), messaging.TopicWalletEvents, merchantEnvelope)
}
