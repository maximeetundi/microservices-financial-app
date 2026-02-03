package services

import (
	"context"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// EventPublisher publishes support events to Kafka
type EventPublisher struct {
	kafkaClient *messaging.KafkaClient
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(kafkaClient *messaging.KafkaClient) *EventPublisher {
	return &EventPublisher{kafkaClient: kafkaClient}
}

// PublishTicketCreated publishes an event when a new support ticket is created
func (p *EventPublisher) PublishTicketCreated(conversationID, userID, userName, subject string) {
	if p.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"conversation_id": conversationID,
		"user_id":         userID,
		"user_name":       userName,
		"subject":         subject,
	}

	p.publish("support.events", "ticket.created", userID, eventData)
}

// PublishUserMessage publishes an event when a user sends a message
func (p *EventPublisher) PublishUserMessage(conversationID, userID, userName, messagePreview string) {
	if p.kafkaClient == nil {
		return
	}

	// Truncate message preview
	if len(messagePreview) > 100 {
		messagePreview = messagePreview[:100] + "..."
	}

	eventData := map[string]interface{}{
		"conversation_id": conversationID,
		"user_id":         userID,
		"user_name":       userName,
		"message_preview": messagePreview,
	}

	p.publish("support.events", "ticket.message", userID, eventData)
}

// PublishEscalation publishes an event when a conversation is escalated to human
func (p *EventPublisher) PublishEscalation(conversationID, userID, userName string) {
	if p.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"conversation_id": conversationID,
		"user_id":         userID,
		"user_name":       userName,
	}

	p.publish("support.events", "ticket.escalated", userID, eventData)
}

func (p *EventPublisher) publish(topic, eventType, source string, data interface{}) {
	event := messaging.NewEventEnvelope(eventType, source, data)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.kafkaClient.Publish(ctx, topic, event)
	if err != nil {
		log.Printf("Failed to publish support event %s: %v", eventType, err)
	} else {
		log.Printf("Published support event: %s", eventType)
	}
}
