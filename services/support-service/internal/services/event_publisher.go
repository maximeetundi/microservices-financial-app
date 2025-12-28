package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// EventPublisher publishes support events to RabbitMQ
type EventPublisher struct {
	channel *amqp.Channel
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(rabbitMQURL string) (*EventPublisher, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange for support events
	err = channel.ExchangeDeclare(
		"support.events", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

	log.Println("Support event publisher initialized")
	return &EventPublisher{channel: channel}, nil
}

// Event represents a support event
type Event struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// PublishTicketCreated publishes an event when a new support ticket is created
func (p *EventPublisher) PublishTicketCreated(conversationID, userID, userName, subject string) {
	if p.channel == nil {
		return
	}

	event := Event{
		Type:      "ticket.created",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         userID,
			"user_name":       userName,
			"subject":         subject,
		},
	}

	p.publish("support.events", "ticket.created", event)
}

// PublishUserMessage publishes an event when a user sends a message
func (p *EventPublisher) PublishUserMessage(conversationID, userID, userName, messagePreview string) {
	if p.channel == nil {
		return
	}

	// Truncate message preview
	if len(messagePreview) > 100 {
		messagePreview = messagePreview[:100] + "..."
	}

	event := Event{
		Type:      "ticket.message",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         userID,
			"user_name":       userName,
			"message_preview": messagePreview,
		},
	}

	p.publish("support.events", "ticket.message", event)
}

// PublishEscalation publishes an event when a conversation is escalated to human
func (p *EventPublisher) PublishEscalation(conversationID, userID, userName string) {
	if p.channel == nil {
		return
	}

	event := Event{
		Type:      "ticket.escalated",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         userID,
			"user_name":       userName,
		},
	}

	p.publish("support.events", "ticket.escalated", event)
}

func (p *EventPublisher) publish(exchange, routingKey string, event Event) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal support event: %v", err)
		return
	}

	err = p.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish support event %s: %v", routingKey, err)
	} else {
		log.Printf("Published support event: %s", routingKey)
	}
}
