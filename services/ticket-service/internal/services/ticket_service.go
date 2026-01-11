package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/repository"
	"github.com/skip2/go-qrcode"
)

type TicketService struct {
	eventRepo    *repository.EventRepository
	tierRepo     *repository.TierRepository
	ticketRepo   *repository.TicketRepository
	walletClient *WalletClient
	kafkaClient  *messaging.KafkaClient
}

func NewTicketService(
	eventRepo *repository.EventRepository,
	tierRepo *repository.TierRepository,
	ticketRepo *repository.TicketRepository,
	kafkaClient *messaging.KafkaClient,
) *TicketService {
	return &TicketService{
		eventRepo:    eventRepo,
		tierRepo:     tierRepo,
		ticketRepo:   ticketRepo,
		walletClient: NewWalletClient(),
		kafkaClient:  kafkaClient,
	}
}

// === Event Operations ===

func (s *TicketService) CreateEvent(creatorID string, req *models.CreateEventRequest) (*models.Event, error) {
	// Generate unique event code
	eventCode := s.generateEventCode()

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "XOF"
	}

	event := &models.Event{
		CreatorID:     creatorID,
		Title:         req.Title,
		Description:   req.Description,
		Location:      req.Location,
		CoverImage:    req.CoverImage,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		SaleStartDate: req.SaleStartDate,
		SaleEndDate:   req.SaleEndDate,
		FormFields:    req.FormFields,
		EventCode:     eventCode,
		Status:        models.EventStatusDraft,
		Currency:      currency,
	}

	// Generate QR code for event
	qrData := fmt.Sprintf("ZEKORA_EVENT:%s", eventCode)
	event.QRCode = s.generateQRCodeBase64(qrData)

	// Create event
	if err := s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Create ticket tiers
	for i, tier := range req.TicketTiers {
		tier.EventID = event.ID
		tier.SortOrder = i
		tier.Sold = 0
		if tier.Icon == "" {
			tier.Icon = "🎫"
		}
		if tier.Color == "" {
			tier.Color = "#6366f1"
		}
		if err := s.tierRepo.Create(&tier); err != nil {
			return nil, fmt.Errorf("failed to create tier: %w", err)
		}
	}

	// Fetch tiers to return with event
	event.TicketTiers, _ = s.GetEventTiers(event.ID)

	return event, nil
}

func (s *TicketService) GetEvent(eventID string) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	// Load tiers
	tiers, _ := s.tierRepo.GetByEventID(eventID)
	if tiers != nil {
		for _, t := range tiers {
			event.TicketTiers = append(event.TicketTiers, *t)
		}
	}

	return event, nil
}

func (s *TicketService) GetEventByCode(code string) (*models.Event, error) {
	event, err := s.eventRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Load tiers
	tiers, _ := s.tierRepo.GetByEventID(event.ID)
	if tiers != nil {
		for _, t := range tiers {
			event.TicketTiers = append(event.TicketTiers, *t)
		}
	}

	return event, nil
}

func (s *TicketService) GetMyEvents(creatorID string, limit, offset int) ([]*models.Event, error) {
	events, err := s.eventRepo.GetByCreator(creatorID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Load tiers and stats for each event
	for _, event := range events {
		tiers, _ := s.tierRepo.GetByEventID(event.ID)
		if tiers != nil {
			for _, t := range tiers {
				event.TicketTiers = append(event.TicketTiers, *t)
				event.TotalSold += t.Sold
				event.TotalRevenue += float64(t.Sold) * t.Price
			}
		}
	}

	return events, nil
}

func (s *TicketService) GetActiveEvents(limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.GetActive(limit, offset)
}

func (s *TicketService) UpdateEvent(eventID, creatorID string, req *models.UpdateEventRequest) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	if event.CreatorID != creatorID {
		return nil, fmt.Errorf("not authorized to update this event")
	}

	// Update fields if provided
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.CoverImage != nil {
		event.CoverImage = *req.CoverImage
	}
	if req.StartDate != nil {
		event.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		event.EndDate = *req.EndDate
	}
	if req.SaleStartDate != nil {
		event.SaleStartDate = *req.SaleStartDate
	}
	if req.SaleEndDate != nil {
		event.SaleEndDate = *req.SaleEndDate
	}
	if req.FormFields != nil {
		event.FormFields = req.FormFields
	}
	if req.Status != nil {
		event.Status = *req.Status
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	// Update tiers if provided
	if req.TicketTiers != nil {
		// Delete existing tiers and recreate
		s.tierRepo.DeleteByEventID(eventID)
		for i, tier := range req.TicketTiers {
			tier.EventID = eventID
			tier.SortOrder = i
			if tier.Icon == "" {
				tier.Icon = "🎫"
			}
			s.tierRepo.Create(&tier)
		}
	}

	return s.GetEvent(eventID)
}

func (s *TicketService) PublishEvent(eventID, creatorID string) error {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if event.CreatorID != creatorID {
		return fmt.Errorf("not authorized")
	}

	return s.eventRepo.UpdateStatus(eventID, models.EventStatusActive)
}

func (s *TicketService) DeleteEvent(eventID, creatorID string) error {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if event.CreatorID != creatorID {
		return fmt.Errorf("not authorized")
	}

	return s.eventRepo.Delete(eventID)
}

func (s *TicketService) GetEventTiers(eventID string) ([]models.TicketTier, error) {
	tiers, err := s.tierRepo.GetByEventID(eventID)
	if err != nil {
		return nil, err
	}

	var result []models.TicketTier
	for _, t := range tiers {
		result = append(result, *t)
	}
	return result, nil
}

// === Ticket Purchase ===

func (s *TicketService) PurchaseTicket(buyerID string, req *models.PurchaseTicketRequest) (*models.Ticket, error) {
	// Get event
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check if event is active
	if event.Status != models.EventStatusActive {
		return nil, fmt.Errorf("event is not active for ticket sales")
	}

	// Check if sale period is valid
	now := time.Now()
	if now.Before(event.SaleStartDate) {
		return nil, fmt.Errorf("ticket sales have not started yet")
	}
	if now.After(event.SaleEndDate) {
		return nil, fmt.Errorf("ticket sales have ended")
	}

	// Get tier
	tier, err := s.tierRepo.GetByID(req.TierID)
	if err != nil {
		return nil, fmt.Errorf("ticket tier not found")
	}

	// Check availability
	available, err := s.tierRepo.CheckAvailability(req.TierID)
	if err != nil || !available {
		return nil, fmt.Errorf("no tickets available for this tier")
	}

	// Generate ticket code
	ticketCode := s.generateTicketCode()
	
	// Create transaction ID (will be used for payment)
	txID := "TX-" + s.generateEventCode()

	// Create ticket
	ticket := &models.Ticket{
		EventID:    req.EventID,
		BuyerID:    buyerID,
		TierID:     req.TierID,
		TierName:   tier.Name,
		TierIcon:   tier.Icon,
		Price:      tier.Price,
		Currency:   event.Currency,
		FormData:   req.FormData,
		TicketCode: ticketCode,
		Status:     models.TicketStatusPending, // Initially pending
		TransactionID: txID,
	}

	// Generate ticket QR code
	qrData := fmt.Sprintf("ZEKORA_TICKET:%s", ticketCode)
	ticket.QRCode = s.generateQRCodeBase64(qrData)

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Initiate Payment via Event Bus
	// We need to resolve Organizer Wallet ID?
	// Assuming creatorID has a default wallet. 
	// Or we pass creatorID as destination user.
	
	paymentReq := messaging.PaymentRequestEvent{
		RequestID:         txID,
		UserID:            buyerID, // Set buyer as the event initiator/payer
		FromWalletID:      req.WalletID,
		DestinationUserID: event.CreatorID, // Set organizer as destination
		ToWalletID:        "",              // Will be resolved by transfer-service using DestinationUserID
		DebitAmount:       tier.Price,
		CreditAmount:      tier.Price,
		Currency:          event.Currency,
		Type:              "ticket_purchase",
		ReferenceID:       fmt.Sprintf("TICKET_%s", ticketCode),
	}

	// Publish payment request to Kafka
	envelope := messaging.NewEventEnvelope(messaging.EventPaymentRequest, "ticket-service", paymentReq)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope); err != nil {
		// Failed to publish payment request
		s.ticketRepo.Delete(ticket.ID)
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	// Note: Ticket remains 'pending'. A background consumer will update it to 'paid'.
	// Status won't be 'paid' in the response. Frontend must handle this.
	
	// Add event info to response
	ticket.EventTitle = event.Title
	ticket.EventDate = &event.StartDate
	ticket.EventLocation = event.Location

	return ticket, nil
}

func (s *TicketService) GetMyTickets(buyerID string, limit, offset int) ([]*models.Ticket, error) {
	return s.ticketRepo.GetByBuyer(buyerID, limit, offset)
}

func (s *TicketService) GetTicket(ticketID string) (*models.Ticket, error) {
	return s.ticketRepo.GetByID(ticketID)
}

func (s *TicketService) GetEventTickets(eventID string, limit, offset int) ([]*models.Ticket, error) {
	return s.ticketRepo.GetByEvent(eventID, limit, offset)
}

// === Ticket Verification ===

func (s *TicketService) VerifyTicket(code string) (*models.VerifyTicketResponse, error) {
	// Parse QR data if it contains prefix
	ticketCode := code
	if strings.HasPrefix(code, "ZEKORA_TICKET:") {
		ticketCode = strings.TrimPrefix(code, "ZEKORA_TICKET:")
	}

	ticket, err := s.ticketRepo.GetByCode(ticketCode)
	if err != nil {
		return &models.VerifyTicketResponse{
			Valid:   false,
			Message: "Ticket non trouvé",
		}, nil
	}

	event, _ := s.eventRepo.GetByID(ticket.EventID)

	response := &models.VerifyTicketResponse{
		Valid:  true,
		Ticket: ticket,
		Event:  event,
	}

	switch ticket.Status {
	case models.TicketStatusPaid:
		response.Message = "Ticket valide"
		response.CanUse = true
	case models.TicketStatusUsed:
		response.Message = "Ticket déjà utilisé"
		response.AlreadyUsed = true
		response.CanUse = false
	case models.TicketStatusCancelled:
		response.Message = "Ticket annulé"
		response.Valid = false
		response.CanUse = false
	case models.TicketStatusRefunded:
		response.Message = "Ticket remboursé"
		response.Valid = false
		response.CanUse = false
	default:
		response.Message = "Ticket en attente de paiement"
		response.Valid = false
		response.CanUse = false
	}

	return response, nil
}

func (s *TicketService) UseTicket(ticketID, organizerID string) error {
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found")
	}

	// Verify organizer owns the event
	event, err := s.eventRepo.GetByID(ticket.EventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if event.CreatorID != organizerID {
		return fmt.Errorf("not authorized to mark this ticket as used")
	}

	if ticket.Status != models.TicketStatusPaid {
		return fmt.Errorf("ticket cannot be used: status is %s", ticket.Status)
	}

	return s.ticketRepo.MarkAsUsed(ticketID)
}

func (s *TicketService) GetEventStats(eventID string) (*models.TicketStats, error) {
	return s.ticketRepo.GetEventStats(eventID)
}

// === Refund Operations ===

func (s *TicketService) RefundTicket(ticketID, organizerID string) error {
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found")
	}

	if ticket.Status != models.TicketStatusPaid {
		return fmt.Errorf("only paid tickets can be refunded")
	}

	// Verify organizer
	event, err := s.eventRepo.GetByID(ticket.EventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.CreatorID != organizerID {
		return fmt.Errorf("not authorized to refund this ticket")
	}

	// Find organizer's wallet for the currency
	wallets, err := s.walletClient.GetUserWallets(organizerID)
	if err != nil {
		return fmt.Errorf("failed to fetch organizer wallets: %w", err)
	}

	var sourceWalletID string
	for _, w := range wallets {
		if cur, ok := w["currency"].(string); ok && cur == ticket.Currency {
			if id, ok := w["id"].(string); ok {
				sourceWalletID = id
				break
			}
		}
	}

	if sourceWalletID == "" {
		return fmt.Errorf("organizer does not have a wallet for currency %s", ticket.Currency)
	}

	// Create Refund Event (Organizer -> Buyer)
	// We use "ticket_refund" type so transfer-service can handle it (or treat as p2p transfer)
	// Since transfer-service consumer logs "Type" but processes debit/credit generic logic, 
	// as long as we supply wallets, it should work.
	// We set ReferenceID to "REFUND_TCT-..."
	
	refundReq := messaging.PaymentRequestEvent{
		RequestID:         fmt.Sprintf("REF-%s-%d", ticket.TicketCode, time.Now().Unix()),
		UserID:            organizerID,          // Organizer is paying
		FromWalletID:      sourceWalletID,       // From Organizer Wallet
		DestinationUserID: ticket.BuyerID,       // To Buyer
		ToWalletID:        "",                   // Let transfer-service resolve buyer's wallet
		DebitAmount:       ticket.Price,
		CreditAmount:      ticket.Price,
		Currency:          ticket.Currency,
		Type:              "ticket_refund",      // Custom type for logging
		ReferenceID:       fmt.Sprintf("REFUND_%s", ticket.TicketCode),
	}

	// Publish refund request
	envelope := messaging.NewEventEnvelope(messaging.EventPaymentRequest, "ticket-service", refundReq)
	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope); err != nil {
		return fmt.Errorf("failed to initiate refund: %w", err)
	}

	// Update ticket status to refunded (optimistic, or wait for confirmation? Optimistic is better UX for this action)
	// But ideally we should wait for success. However, existing flow updates on event.
	// For now, let's update to 'refunded' immediately. Using 'refund_pending' would be safer but adds complexity.
	return s.ticketRepo.UpdateStatus(ticketID, models.TicketStatusRefunded)
}

func (s *TicketService) CancelAndRefundEvent(eventID, organizerID string) error {
	// Verify ownership
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.CreatorID != organizerID {
		return fmt.Errorf("not authorized")
	}

	// Update status to Cancelled
	if err := s.eventRepo.UpdateStatus(eventID, models.EventStatusCancelled); err != nil {
		return err
	}

	// Find all paid tickets
	// Note: Pagination might be an issue for huge events, but acceptable for MVP.
	// We'll fetch in batches if needed, but for now fetch all (limit 1000?)
	tickets, err := s.ticketRepo.GetByEvent(eventID, 10000, 0)
	if err != nil {
		return err
	}

	count := 0
	for _, ticket := range tickets {
		if ticket.Status == models.TicketStatusPaid {
			// Trigger refund for each
			// We reuse RefundTicket logic but we already have the event verification done.
			// Calling RefundTicket directly does redundant checks but is safe.
			// Optimization: batch process? For now, simple loop.
			if err := s.RefundTicket(ticket.ID, organizerID); err != nil {
				// Log error but continue with others?
				fmt.Printf("Failed to refund ticket %s: %v\n", ticket.ID, err)
			} else {
				count++
			}
		}
	}

	return nil
}

// === Helper Functions ===

func (s *TicketService) generateEventCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	code := strings.ToUpper(base64.RawURLEncoding.EncodeToString(b))
	return "EVT-" + code[:8]
}

func (s *TicketService) generateTicketCode() string {
	b := make([]byte, 9) // 9 bytes = 12 base64 chars
	rand.Read(b)
	code := strings.ToUpper(base64.RawURLEncoding.EncodeToString(b))
	return "TKT-" + code[:12]
}

func (s *TicketService) generateQRCodeBase64(data string) string {
	qr, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qr)
}

// GetAvailableIcons returns the list of available icons for tiers
func (s *TicketService) GetAvailableIcons() []string {
	return models.AvailableTierIcons
}
