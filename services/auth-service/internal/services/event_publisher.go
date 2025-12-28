package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// EventPublisher publishes events to RabbitMQ for other services
type EventPublisher struct {
	channel *amqp.Channel
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(channel *amqp.Channel) *EventPublisher {
	return &EventPublisher{channel: channel}
}

// Event represents an event to be published
type Event struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// PublishKYCSubmitted publishes an event when a user submits KYC documents
func (p *EventPublisher) PublishKYCSubmitted(userID, userEmail, userName, docType string) {
	if p.channel == nil {
		return
	}

	event := Event{
		Type:      "kyc.submitted",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"user_id":    userID,
			"user_email": userEmail,
			"user_name":  userName,
			"doc_type":   docType,
		},
	}

	p.publish("user.events", "kyc.submitted", event)
}

// PublishPINBlocked publishes an event when a user's PIN is blocked
func (p *EventPublisher) PublishPINBlocked(userID, userEmail, userName, reason string) {
	if p.channel == nil {
		return
	}

	event := Event{
		Type:      "pin.blocked",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"user_id":    userID,
			"user_email": userEmail,
			"user_name":  userName,
			"reason":     reason,
		},
	}

	p.publish("user.events", "pin.blocked", event)
}

// PublishPINUnlocked publishes an event when a user's PIN is unlocked
func (p *EventPublisher) PublishPINUnlocked(userID, userEmail, userName string) {
	if p.channel == nil {
		return
	}

	event := Event{
		Type:      "pin.unlocked",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"user_id":    userID,
			"user_email": userEmail,
			"user_name":  userName,
		},
	}

	p.publish("user.events", "pin.unlocked", event)
}

// publish sends an event to RabbitMQ
func (p *EventPublisher) publish(exchange, routingKey string, event Event) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = p.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish event %s: %v", routingKey, err)
	} else {
		log.Printf("Published event: %s", routingKey)
	}
}
