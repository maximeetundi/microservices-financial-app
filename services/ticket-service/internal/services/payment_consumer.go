package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/repository"
)

type PaymentStatusConsumer struct {
	kafkaClient *messaging.KafkaClient
	ticketRepo  *repository.TicketRepository
}

func NewPaymentStatusConsumer(kafkaClient *messaging.KafkaClient, ticketRepo *repository.TicketRepository) *PaymentStatusConsumer {
	return &PaymentStatusConsumer{
		kafkaClient: kafkaClient,
		ticketRepo:  ticketRepo,
	}
}

func (c *PaymentStatusConsumer) Start() error {
	// Subscribe to payment events
	if err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent); err != nil {
		return err
	}

	log.Println("[Kafka] Payment Status Consumer started for ticket-service")
	return nil
}

func (c *PaymentStatusConsumer) handlePaymentEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	// Only process payment success/failed events
	if event.Type != messaging.EventPaymentSuccess && event.Type != messaging.EventPaymentFailed {
		return nil
	}

	log.Printf("[Kafka] Received payment event: %s", event.Type)

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var statusEvent messaging.PaymentStatusEvent
	if err := json.Unmarshal(dataBytes, &statusEvent); err != nil {
		log.Printf("Error unmarshalling payment status event: %v", err)
		return err
	}

	log.Printf("[Kafka] Payment status for Request: %s - Status: %s - Error: %s", statusEvent.RequestID, statusEvent.Status, statusEvent.Error)

	// Determine new ticket status
	newStatus := "paid"
	if statusEvent.Status != "success" {
		newStatus = "failed"
	}

	// Update ticket by transaction ID
	err = c.ticketRepo.UpdateStatusByTransactionID(statusEvent.RequestID, newStatus)
	if err != nil {
		log.Printf("[Kafka] Failed to update ticket status for TX %s: %v", statusEvent.RequestID, err)
	} else {
		log.Printf("[Kafka] Updated ticket status for TX %s to %s", statusEvent.RequestID, newStatus)


		// If payment successful, notify the organizer and the buyer
		if newStatus == "paid" {
			tickets, err := c.ticketRepo.GetListByTransactionID(statusEvent.RequestID)
			if err == nil && len(tickets) > 0 {
				ticket := tickets[0] // Representative ticket for common details
				count := len(tickets)

				// 1. Notify Organizer
				if ticket.EventCreatorID != "" {
					msg := "Un ticket a été acheté pour votre événement " + ticket.EventTitle
					if count > 1 {
						msg = fmt.Sprintf("%d tickets ont été achetés pour votre événement %s", count, ticket.EventTitle)
					}
					
					notifEvent := messaging.NotificationEvent{
						UserID:  ticket.EventCreatorID,
						Type:    "ticket_sold",
						Title:   "Nouveau ticket vendu ! 🎫",
						Message: msg,
						Data: map[string]interface{}{
							"ticket_id": ticket.ID, // Just pass one ID or maybe none if batch
							"event_id":  ticket.EventID,
							"amount":     ticket.Price * float64(count),
							"currency":   ticket.Currency,
							"quantity":   count,
						},
					}
					
					envelope := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "ticket-service", notifEvent)
					if err := c.kafkaClient.Publish(ctx, messaging.TopicNotificationEvents, envelope); err != nil {
						log.Printf("[Kafka] Failed to publish organizer notification: %v", err)
					}
				}

				// 2. Notify Buyer
				if ticket.BuyerID != "" {
					msg := "Votre ticket " + ticket.TierName + " pour " + ticket.EventTitle + " a été payé avec succès."
					if count > 1 {
						msg = fmt.Sprintf("Vos %d tickets %s pour %s ont été payés avec succès.", count, ticket.TierName, ticket.EventTitle)
					}

					buyerNotifEvent := messaging.NotificationEvent{
						UserID:  ticket.BuyerID,
						Type:    "ticket_purchased",
						Title:   "Paiement réussi ✅",
						Message: msg,
						Data: map[string]interface{}{
							"ticket_id": ticket.ID,
							"event_id":  ticket.EventID,
							"amount":    ticket.Price * float64(count),
							"currency":  ticket.Currency,
							"route":     "/my-tickets",
							"quantity":  count,
						},
					}
					
					buyerEnvelope := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "ticket-service", buyerNotifEvent)
					if err := c.kafkaClient.Publish(ctx, messaging.TopicNotificationEvents, buyerEnvelope); err != nil {
						log.Printf("[Kafka] Failed to publish buyer notification: %v", err)
					}
				}
			}
		}
	}

	return nil
}
