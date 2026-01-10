package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// KafkaConfig holds configuration for Kafka connection
type KafkaConfig struct {
	Brokers       []string
	GroupID       string
	ClientID      string
	RetryAttempts int
	RetryDelay    time.Duration
}

// KafkaClient provides a unified interface for Kafka operations
// Implements Event Bus pattern with Consumer Groups support
type KafkaClient struct {
	config    KafkaConfig
	writers   map[string]*kafka.Writer
	readers   map[string]*kafka.Reader
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewKafkaClient creates a new Kafka client with the given configuration
func NewKafkaClient(brokers []string, groupID string) *KafkaClient {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &KafkaClient{
		config: KafkaConfig{
			Brokers:       brokers,
			GroupID:       groupID,
			ClientID:      groupID,
			RetryAttempts: 3,
			RetryDelay:    time.Second,
		},
		writers: make(map[string]*kafka.Writer),
		readers: make(map[string]*kafka.Reader),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// getWriter returns a writer for the given topic, creating one if needed
func (k *KafkaClient) getWriter(topic string) *kafka.Writer {
	k.mu.RLock()
	writer, exists := k.writers[topic]
	k.mu.RUnlock()
	
	if exists {
		return writer
	}
	
	k.mu.Lock()
	defer k.mu.Unlock()
	
	// Double-check after acquiring write lock
	if writer, exists = k.writers[topic]; exists {
		return writer
	}
	
	writer = &kafka.Writer{
		Addr:         kafka.TCP(k.config.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        false, // Synchronous for reliability
		RequiredAcks: kafka.RequireOne,
	}
	k.writers[topic] = writer
	
	return writer
}

// Publish sends an event to the specified topic with retry logic
func (k *KafkaClient) Publish(ctx context.Context, topic string, event *EventEnvelope) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.Version == "" {
		event.Version = "1.0"
	}
	
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "source", Value: []byte(event.Source)},
			{Key: "correlation_id", Value: []byte(event.CorrelationID)},
		},
	}
	
	writer := k.getWriter(topic)
	
	var lastErr error
	for attempt := 0; attempt < k.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(k.config.RetryDelay * time.Duration(attempt))
			log.Printf("[Kafka] Retry attempt %d for topic %s", attempt+1, topic)
		}
		
		if err := writer.WriteMessages(ctx, msg); err != nil {
			lastErr = err
			log.Printf("[Kafka] Publish failed (attempt %d): %v", attempt+1, err)
			continue
		}
		
		log.Printf("[Kafka] Published event %s to topic %s", event.Type, topic)
		return nil
	}
	
	return fmt.Errorf("failed to publish after %d attempts: %w", k.config.RetryAttempts, lastErr)
}

// PublishRaw sends raw JSON data to a topic (for backward compatibility)
func (k *KafkaClient) PublishRaw(ctx context.Context, topic string, data []byte) error {
	msg := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: data,
	}
	
	writer := k.getWriter(topic)
	return writer.WriteMessages(ctx, msg)
}

// EventHandler is a function that processes an event
type EventHandler func(ctx context.Context, event *EventEnvelope) error

// Subscribe starts consuming messages from a topic with the given handler
// Uses Consumer Group pattern for horizontal scaling
func (k *KafkaClient) Subscribe(topic string, handler EventHandler) error {
	k.mu.Lock()
	
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.config.Brokers,
		Topic:          topic,
		GroupID:        k.config.GroupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        500 * time.Millisecond,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second,
	})
	k.readers[topic] = reader
	k.mu.Unlock()
	
	go k.consumeLoop(reader, topic, handler)
	
	log.Printf("[Kafka] Subscribed to topic %s with group %s", topic, k.config.GroupID)
	return nil
}

// consumeLoop runs the main consumption loop for a topic
func (k *KafkaClient) consumeLoop(reader *kafka.Reader, topic string, handler EventHandler) {
	for {
		select {
		case <-k.ctx.Done():
			log.Printf("[Kafka] Stopping consumer for topic %s", topic)
			return
		default:
			msg, err := reader.ReadMessage(k.ctx)
			if err != nil {
				if k.ctx.Err() != nil {
					return // Context cancelled
				}
				log.Printf("[Kafka] Error reading message from %s: %v", topic, err)
				time.Sleep(time.Second)
				continue
			}
			
			var event EventEnvelope
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("[Kafka] Failed to unmarshal message from %s: %v", topic, err)
				// Send to DLQ in production
				continue
			}
			
			if err := handler(k.ctx, &event); err != nil {
				log.Printf("[Kafka] Handler error for event %s: %v", event.Type, err)
				// Implement retry logic or DLQ here
			}
		}
	}
}

// SubscribeRaw subscribes with a raw message handler (for backward compatibility)
type RawHandler func(ctx context.Context, data []byte) error

func (k *KafkaClient) SubscribeRaw(topic string, handler RawHandler) error {
	return k.Subscribe(topic, func(ctx context.Context, event *EventEnvelope) error {
		data, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}
		return handler(ctx, data)
	})
}

// Close gracefully shuts down all writers and readers
func (k *KafkaClient) Close() error {
	k.cancel()
	
	k.mu.Lock()
	defer k.mu.Unlock()
	
	var errs []error
	
	for topic, writer := range k.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close writer for %s: %w", topic, err))
		}
	}
	
	for topic, reader := range k.readers {
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close reader for %s: %w", topic, err))
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("errors closing Kafka client: %v", errs)
	}
	
	log.Println("[Kafka] Client closed successfully")
	return nil
}

// Health checks if Kafka connection is healthy
func (k *KafkaClient) Health(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", k.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()
	
	_, err = conn.Brokers()
	if err != nil {
		return fmt.Errorf("failed to get brokers: %w", err)
	}
	
	return nil
}
