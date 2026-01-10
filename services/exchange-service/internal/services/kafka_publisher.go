package services

import (
	"context"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// KafkaPublisher provides a simplified interface for publishing Kafka events
// This wrapper makes it easier to migrate from RabbitMQ and maintain the codebase
type KafkaPublisher struct {
	client *messaging.KafkaClient
}

// NewKafkaPublisher creates a new Kafka publisher wrapper
func NewKafkaPublisher(client *messaging.KafkaClient) *KafkaPublisher {
	return &KafkaPublisher{client: client}
}

// PublishPaymentRequest publishes a payment request event
func (p *KafkaPublisher) PublishPaymentRequest(req *messaging.PaymentRequestEvent) error {
	if p.client == nil {
		log.Println("[Kafka] Warning: client is nil, skipping publish")
		return nil
	}

	envelope := messaging.NewEventEnvelope(messaging.EventPaymentRequest, "exchange-service", req)
	return p.client.Publish(context.Background(), messaging.TopicPaymentEvents, envelope)
}

// PublishExchangeEvent publishes an exchange-related event
func (p *KafkaPublisher) PublishExchangeEvent(eventType string, data map[string]interface{}) error {
	if p.client == nil {
		return nil
	}

	envelope := messaging.NewEventEnvelope(eventType, "exchange-service", data)
	return p.client.Publish(context.Background(), messaging.TopicExchangeEvents, envelope)
}

// PublishTradingEvent publishes a trading-related event
func (p *KafkaPublisher) PublishTradingEvent(eventType string, data map[string]interface{}) error {
	if p.client == nil {
		return nil
	}

	envelope := messaging.NewEventEnvelope(eventType, "exchange-service", data)
	return p.client.Publish(context.Background(), messaging.TopicExchangeEvents, envelope)
}

// IsConnected checks if the Kafka client is available
func (p *KafkaPublisher) IsConnected() bool {
	return p.client != nil
}

// GetClient returns the underlying Kafka client for advanced operations
func (p *KafkaPublisher) GetClient() *messaging.KafkaClient {
	return p.client
}
