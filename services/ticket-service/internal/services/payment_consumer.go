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

	// Update ticket by transaction ID (which is the RequestID in our case)
	err = c.ticketRepo.UpdateStatusByTransactionID(statusEvent.RequestID, newStatus)
	if err != nil {
		log.Printf("Failed to update ticket status for TX %s: %v", statusEvent.RequestID, err)
	} else {
		log.Printf("[Kafka] Updated ticket status for TX %s to %s", statusEvent.RequestID, newStatus)
	}

	return nil
}
