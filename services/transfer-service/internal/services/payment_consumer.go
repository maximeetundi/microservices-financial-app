package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/repository"
)

// PaymentRequestConsumer handles payment request events from Kafka
type PaymentRequestConsumer struct {
	kafkaClient  *messaging.KafkaClient
	walletClient *WalletClient
	walletRepo   *repository.WalletRepository
}

// NewPaymentRequestConsumer creates a new payment request consumer
func NewPaymentRequestConsumer(kafkaClient *messaging.KafkaClient, walletClient *WalletClient, walletRepo *repository.WalletRepository) *PaymentRequestConsumer {
	return &PaymentRequestConsumer{
		kafkaClient:  kafkaClient,
		walletClient: walletClient,
		walletRepo:   walletRepo,
	}
}

// Start begins consuming payment events from Kafka
func (c *PaymentRequestConsumer) Start() error {
	// Subscribe to payment events
	if err := c.kafkaClient.Subscribe(messaging.TopicPaymentEvents, c.handlePaymentEvent); err != nil {
		return err
	}

	log.Println("[Kafka] PaymentRequestConsumer started")
	return nil
}

// handlePaymentEvent processes payment.request events
func (c *PaymentRequestConsumer) handlePaymentEvent(ctx context.Context, event *messaging.EventEnvelope) error {
	log.Printf("[Kafka] Processing payment event: %s", event.Type)

	if event.Type != messaging.EventPaymentRequest {
		return nil
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var paymentReq messaging.PaymentRequestEvent
	if err := json.Unmarshal(dataBytes, &paymentReq); err != nil {
		log.Printf("Failed to unmarshal payment request: %v", err)
		return err
	}

	log.Printf("[Kafka] Processing payment request: RequestID=%s, Type=%s, Amount=%.2f",
		paymentReq.RequestID, paymentReq.Type, paymentReq.DebitAmount)

	// Process debit operation via wallet client
	err = c.walletClient.ProcessTransaction(paymentReq.FromWalletID, paymentReq.DebitAmount, paymentReq.Currency, "debit", paymentReq.ReferenceID)
	if err != nil {
		log.Printf("[Kafka] Debit failed: %v", err)
		c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", err.Error())
		return err
	}

	// If there's a credit operation (to wallet)
	if paymentReq.ToWalletID != "" && paymentReq.CreditAmount > 0 {
		err = c.walletClient.ProcessTransaction(paymentReq.ToWalletID, paymentReq.CreditAmount, paymentReq.Currency, "credit", paymentReq.ReferenceID)
		if err != nil {
			log.Printf("[Kafka] Credit failed: %v", err)
			// Rollback debit
			c.walletClient.ProcessTransaction(paymentReq.FromWalletID, paymentReq.DebitAmount, paymentReq.Currency, "credit", paymentReq.ReferenceID+"_rollback")
			c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "failed", err.Error())
			return err
		}
	}

	log.Printf("[Kafka] Payment request %s processed successfully", paymentReq.RequestID)
	c.publishPaymentStatus(paymentReq.RequestID, paymentReq.ReferenceID, paymentReq.Type, "success", "")

	return nil
}

// publishPaymentStatus publishes status updates back to Kafka
func (c *PaymentRequestConsumer) publishPaymentStatus(requestID, referenceID, reqType, status, errMsg string) {
	statusEvent := messaging.PaymentStatusEvent{
		RequestID:   requestID,
		ReferenceID: referenceID,
		Type:        reqType,
		Status:      status,
		Error:       errMsg,
	}

	eventType := messaging.EventPaymentSuccess
	if status == "failed" {
		eventType = messaging.EventPaymentFailed
	}

	envelope := messaging.NewEventEnvelope(eventType, "transfer-service", statusEvent)
	c.kafkaClient.Publish(context.Background(), messaging.TopicPaymentEvents, envelope)
}
