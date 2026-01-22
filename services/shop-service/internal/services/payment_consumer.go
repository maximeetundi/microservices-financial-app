package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/shop-service/internal/models"
)

type PaymentConsumer struct {
	kafkaClient  *messaging.KafkaClient
	orderService *OrderService
}

func NewPaymentConsumer(kafkaClient *messaging.KafkaClient, orderService *OrderService) *PaymentConsumer {
	return &PaymentConsumer{
		kafkaClient:  kafkaClient,
		orderService: orderService,
	}
}

func (c *PaymentConsumer) Start() {
	if c.kafkaClient == nil {
		log.Println("Kafka client not available, payment consumer not started")
		return
	}

	log.Println("Starting payment consumer...")
	
	handler := func(topic string, message []byte) error {
		log.Printf("Received message on topic %s", topic)
		
		var event models.PaymentStatusEvent
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Failed to unmarshal payment status event: %v", err)
			return nil // Don't retry on unmarshal errors
		}
		
		ctx := context.Background()
		if err := c.orderService.HandlePaymentStatus(ctx, &event); err != nil {
			log.Printf("Failed to handle payment status: %v", err)
			return err
		}
		
		log.Printf("Successfully processed payment status for transaction %s: %s", event.TransactionID, event.Status)
		return nil
	}
	
	// Subscribe to payment status topic
	if err := c.kafkaClient.Subscribe("wallet.payment.status", handler); err != nil {
		log.Printf("Failed to subscribe to payment status topic: %v", err)
		return
	}
	
	log.Println("Payment consumer started successfully")
	
	// Keep the consumer running
	select {}
}
