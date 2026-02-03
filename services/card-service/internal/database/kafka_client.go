package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	writer *kafka.Writer
}

func InitializeKafka(brokers string) (*KafkaClient, error) {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		// Batch settings
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		// Retry settings
		MaxAttempts:  10,
		WriteTimeout: 10 * time.Second,
	}

	return &KafkaClient{writer: w}, nil
}

func (k *KafkaClient) Close() {
	if k.writer != nil {
		k.writer.Close()
	}
}

func (k *KafkaClient) Publish(topic string, key string, message interface{}) error {
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	log.Printf("Publishing to topic %s: %s", topic, string(value))

	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
			Time:  time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

// Event payload structure
type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

func (k *KafkaClient) PublishEvent(topic string, eventType string, payload interface{}) error {
	event := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	return k.Publish(topic, eventType, event)
}
