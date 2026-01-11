package services

import (
	"context"
	"encoding/json"
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

	log.Printf("[Kafka] Payment status for Request: %s - Status: %s", statusEvent.RequestID, statusEvent.Status)

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
			ticket, err := c.ticketRepo.GetByTransactionID(statusEvent.RequestID)
			if err == nil && ticket != nil {
				// 1. Notify Organizer
				if ticket.EventCreatorID != "" {
					notifEvent := messaging.NotificationEvent{
						UserID:  ticket.EventCreatorID,
						Type:    "ticket_sold",
						Title:   "Nouveau ticket vendu ! 🎫",
						Message: "Un ticket a été acheté pour votre événement " + ticket.EventTitle,
						Data: map[string]interface{}{
							"ticket_id": ticket.ID,
							"event_id":  ticket.EventID,
							"amount":    ticket.Price,
							"currency":  ticket.Currency,
						},
					}
					
					envelope := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "ticket-service", notifEvent)
					if err := c.kafkaClient.Publish(ctx, messaging.TopicNotificationEvents, envelope); err != nil {
						log.Printf("[Kafka] Failed to publish organizer notification: %v", err)
					}
				}

				// 2. Notify Buyer
				if ticket.BuyerID != "" {
					buyerNotifEvent := messaging.NotificationEvent{
						UserID:  ticket.BuyerID,
						Type:    "ticket_purchased",
						Title:   "Ticket acheté ! 🎟️",
						Message: "Votre ticket pour " + ticket.EventTitle + " est maintenant disponible.",
						Data: map[string]interface{}{
							"ticket_id": ticket.ID,
							"event_id":  ticket.EventID,
							"amount":    ticket.Price,
							"currency":  ticket.Currency,
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
