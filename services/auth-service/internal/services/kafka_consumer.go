package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

type KafkaConsumer struct {
	client   *messaging.KafkaClient
	userRepo *repository.UserRepository
}

func NewKafkaConsumer(client *messaging.KafkaClient, userRepo *repository.UserRepository) *KafkaConsumer {
	return &KafkaConsumer{
		client:   client,
		userRepo: userRepo,
	}
}

func (c *KafkaConsumer) Start() error {
    log.Println("Starting Auth Service Kafka consumer...")
	return c.client.Subscribe(messaging.TopicUserEvents, c.handleUserEvent)
}

func (c *KafkaConsumer) handleUserEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	if event.Type == "user.balance_updated" {
		dataBytes, err := json.Marshal(event.Data)
		if err != nil {
            return err
        }
		var data struct {
			UserID string  `json:"user_id"`
			Total  float64 `json:"total_balance_usd"`
		}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
            log.Printf("Failed to unmarshal update event: %v", err)
            return err
        }

		if data.UserID != "" {
			log.Printf("Updating Total Balance for user %s: %.2f USD", data.UserID, data.Total)
			return c.userRepo.UpdateTotalBalanceUSD(data.UserID, data.Total)
		}
	}
	return nil
}
