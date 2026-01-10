package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// EventConsumer consumes events from Kafka for admin notifications
type EventConsumer struct {
	kafkaClient *messaging.KafkaClient
	db          *sql.DB
}

// Event represents an incoming event from Kafka
type Event struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

func NewEventConsumer(kafkaClient *messaging.KafkaClient, db *sql.DB) *EventConsumer {
	return &EventConsumer{
		kafkaClient: kafkaClient,
		db:          db,
	}
}

func (c *EventConsumer) StartConsuming() error {
	// Subscribe to user events
	if err := c.kafkaClient.Subscribe(messaging.TopicUserEvents, c.handleUserEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to user events: %v", err)
	}

	// Subscribe to wallet events
	if err := c.kafkaClient.Subscribe(messaging.TopicWalletEvents, c.handleWalletEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to wallet events: %v", err)
	}

	// Subscribe to transfer events
	if err := c.kafkaClient.Subscribe(messaging.TopicTransferEvents, c.handleTransferEvent); err != nil {
		log.Printf("Warning: Failed to subscribe to transfer events: %v", err)
	}

	log.Println("[Kafka] Admin event consumer started")
	return nil
}

func (c *EventConsumer) handleUserEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Admin received user event: %s", event.Type)
	return c.storeNotification(event.Type, event.Data)
}

func (c *EventConsumer) handleWalletEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Admin received wallet event: %s", event.Type)
	return c.storeNotification(event.Type, event.Data)
}

func (c *EventConsumer) handleTransferEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Admin received transfer event: %s", event.Type)
	return c.storeNotification(event.Type, event.Data)
}

func (c *EventConsumer) storeNotification(eventType string, data interface{}) error {
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)

	_, err := c.db.Exec(`
		INSERT INTO admin_notifications (id, type, title, message, data, is_read, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, false, NOW())
	`, eventType, eventType, eventType, dataStr)

	if err != nil {
		log.Printf("[Kafka] Failed to store notification: %v", err)
		return err
	}

	return nil
}
