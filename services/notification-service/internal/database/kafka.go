package database

import (
	"log"
	"strings"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// InitializeKafka creates a new Kafka client for messaging
func InitializeKafka(brokers string, groupID string) (*messaging.KafkaClient, error) {
	brokerList := strings.Split(brokers, ",")
	
	client := messaging.NewKafkaClient(brokerList, groupID)
	
	log.Printf("[Kafka] Notification-service connected to brokers: %s with group: %s", brokers, groupID)
	return client, nil
}
