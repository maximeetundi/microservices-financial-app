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

		// If payment successful, notify the organizer
		if newStatus == "paid" {
			ticket, err := c.ticketRepo.GetByTransactionID(statusEvent.RequestID)
				// Let's assume ticket struct has what we need or we fetch event.
				// Checking ticket model in repo: populateEventDetails gets Title, Date, Location.
				// We probably need to fetch the Event to get CreatorID.
				
				// Let's try to fetch event directly if needed, or query ticketRepo events collection?
				// TicketRepo has eventColl.
				// But we are in consumer.
				// Let's check if we can get CreatorID easily. 
				// For now, let's implement the notification publishing assuming we can get the UserID.
				// Oops, I need to be sure about CreatorID.
				// Let's look at the Ticket struct in the repo/model first if I need to.
				// But waiting for view_file is better? No, I am in parallel execution manually here?
				// I will implement a helper or logic here.
				
				// Just implementing the update for now inside this block is risky without knowing fields.
				// I will hold off on this replace until I verify the Ticket model in the next step.
				// Actually, I can do it in two steps.
			}
		}
	}

	return nil
}
