package services

import (
	"context"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// EventPublisher publishes events to Kafka for other services
type EventPublisher struct {
	kafkaClient *messaging.KafkaClient
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(kafkaClient *messaging.KafkaClient) *EventPublisher {
	return &EventPublisher{kafkaClient: kafkaClient}
}

// PublishKYCSubmitted publishes an event when a user submits KYC documents
func (p *EventPublisher) PublishKYCSubmitted(userID, userEmail, userName, docType string) {
	if p.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"user_id":    userID,
		"user_email": userEmail,
		"user_name":  userName,
		"doc_type":   docType,
	}

	p.publish("user.events", "kyc.submitted", userID, eventData)
}

// PublishPINBlocked publishes an event when a user's PIN is blocked
func (p *EventPublisher) PublishPINBlocked(userID, userEmail, userName, reason string) {
	if p.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"user_id":    userID,
		"user_email": userEmail,
		"user_name":  userName,
		"reason":     reason,
	}

	p.publish("user.events", "pin.blocked", userID, eventData)
}

// PublishPINUnlocked publishes an event when a user's PIN is unlocked
func (p *EventPublisher) PublishPINUnlocked(userID, userEmail, userName string) {
	if p.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"user_id":    userID,
		"user_email": userEmail,
		"user_name":  userName,
	}

	p.publish("user.events", "pin.unlocked", userID, eventData)
}

// publish sends an event to Kafka
func (p *EventPublisher) publish(topic, eventType, source string, data interface{}) {
	event := messaging.NewEventEnvelope(eventType, source, data)

	// Use background context for async publishing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.kafkaClient.Publish(ctx, topic, event)
	if err != nil {
		log.Printf("Failed to publish event %s to topic %s: %v", eventType, topic, err)
	} else {
		log.Printf("Published event: %s to %s", eventType, topic)
	}
}
