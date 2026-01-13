package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

type KafkaConsumer struct {
	client      *messaging.KafkaClient
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewKafkaConsumer(client *messaging.KafkaClient, userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *KafkaConsumer {
	return &KafkaConsumer{
		client:      client,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (c *KafkaConsumer) Start() error {
	log.Println("Starting Auth Service Kafka consumer...")
	return c.client.Subscribe(messaging.TopicUserEvents, c.handleUserEvent)
}

func (c *KafkaConsumer) handleUserEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	switch event.Type {
	case "user.balance_updated":
		return c.handleBalanceUpdated(event)
	case "user.block":
		return c.handleUserBlocked(event)
	}
	return nil
}

func (c *KafkaConsumer) handleBalanceUpdated(event *messaging.EventEnvelope) error {
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
	return nil
}

func (c *KafkaConsumer) handleUserBlocked(event *messaging.EventEnvelope) error {
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	var data struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		log.Printf("Failed to unmarshal block event: %v", err)
		return err
	}

	if data.UserID != "" {
		log.Printf("Revoking sessions for blocked user %s. Reason: %s", data.UserID, data.Reason)
		// Revoke all sessions immediately
		if err := c.sessionRepo.RevokeAllUserSessions(data.UserID); err != nil {
			log.Printf("Failed to revoke sessions for blocked user %s: %v", data.UserID, err)
			return err
		}
		log.Printf("Successfully revoked all sessions for user %s", data.UserID)
	}
	return nil
}
