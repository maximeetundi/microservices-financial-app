package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/database"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/repository"
)

type PaymentStatusConsumer struct {
	mqClient   *database.RabbitMQClient
	ticketRepo *repository.TicketRepository
}

func NewPaymentStatusConsumer(mqClient *database.RabbitMQClient, ticketRepo *repository.TicketRepository) *PaymentStatusConsumer {
	return &PaymentStatusConsumer{
		mqClient:   mqClient,
		ticketRepo: ticketRepo,
	}
}

func (c *PaymentStatusConsumer) Start() error {
	msgs, err := c.mqClient.Consume("ticket.payment.updates")
	if err != nil {
		return fmt.Errorf("failed to start consuming ticket payment updates queue: %w", err)
	}

	go func() {
		for d := range msgs {
			var event models.PaymentStatusEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshalling payment status event: %v", err)
				d.Ack(false)
				continue
			}

			log.Printf("Received payment status for TX: %s - Status: %s", event.TransactionID, event.Status)

			// Logic to check if this transaction belongs to ticket service
			// We can check prefix or just try update.
			// "TX-EVT-..." is generic but likely unique enough.
			// Actually TicketService sets TransactionID = "TX-" + EventCode (EVT-...) 
			// So "TX-EVT-..." prefix.
			
			newStatus := "paid"
			if event.Status != "completed" {
				newStatus = "failed"
				// Optionally "cancelled" if we distinguish.
			}

			err := c.ticketRepo.UpdateStatusByTransactionID(event.TransactionID, newStatus)
			if err != nil {
				log.Printf("Failed to update ticket status for TX %s: %v", event.TransactionID, err)
				// Nack? Or just log.
				// If error is sql.ErrNoRows, then it's not a ticket transaction or invalid ID.
				// In that case, we should Ack to ignore.
				// If db error, maybe Nack.
			} else {
				log.Printf("Updated ticket status for TX %s to %s", event.TransactionID, newStatus)
			}

			d.Ack(false)
		}
	}()

	log.Println("Payment Status Consumer started")
	return nil
}
