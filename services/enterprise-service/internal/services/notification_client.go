package services

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
)

// NotificationClient sends notifications via Kafka to notification-service
type NotificationClient struct {
	kafka *messaging.KafkaClient
}

func NewNotificationClient(kafka *messaging.KafkaClient) *NotificationClient {
	return &NotificationClient{kafka: kafka}
}

// === Enterprise Notification Types ===

type EnterpriseNotificationType string

const (
	// Employee Notifications
	NotifyEmployeeInvited    EnterpriseNotificationType = "enterprise.employee.invited"
	NotifyEmployeeAccepted   EnterpriseNotificationType = "enterprise.employee.accepted"
	NotifyEmployeeTerminated EnterpriseNotificationType = "enterprise.employee.terminated"
	NotifyEmployeePromoted   EnterpriseNotificationType = "enterprise.employee.promoted"
	
	// Payroll Notifications
	NotifyPayrollExecuted    EnterpriseNotificationType = "enterprise.payroll.executed"
	NotifySalaryReceived     EnterpriseNotificationType = "enterprise.salary.received"
	
	// Invoice/Subscription Notifications
	NotifyNewSubscription    EnterpriseNotificationType = "enterprise.subscription.created"
	NotifyInvoiceCreated     EnterpriseNotificationType = "enterprise.invoice.created"
	NotifyInvoiceDue         EnterpriseNotificationType = "enterprise.invoice.due"
	NotifyPaymentReceived    EnterpriseNotificationType = "enterprise.payment.received"
)

// NotifyUser sends a notification to a specific user
func (n *NotificationClient) NotifyUser(ctx context.Context, userID string, notifType string, title, message string, data map[string]interface{}) error {
	if n.kafka == nil {
		return fmt.Errorf("kafka client not initialized")
	}
	
	event := messaging.NotificationEvent{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    data,
	}
	
	envelope := messaging.NewEventEnvelope(messaging.EventNotificationCreated, "enterprise-service", event)
	return n.kafka.Publish(ctx, messaging.TopicNotificationEvents, envelope)
}

// === Convenience Methods ===

// NotifyEmployeeInvitation sends invitation notification to employee
func (n *NotificationClient) NotifyEmployeeInvitation(ctx context.Context, employeeUserID, enterpriseName string) error {
	return n.NotifyUser(ctx, employeeUserID, string(NotifyEmployeeInvited),
		"Invitation d'emploi",
		fmt.Sprintf("Vous avez reçu une invitation pour rejoindre %s. Ouvrez l'app pour accepter.", enterpriseName),
		map[string]interface{}{
			"enterprise_name": enterpriseName,
			"action":          "accept_invitation",
		},
	)
}

// NotifyEmployeeAcceptance sends confirmation to enterprise owner when employee accepts
func (n *NotificationClient) NotifyEmployeeAcceptance(ctx context.Context, ownerUserID, employeeName, enterpriseName string) error {
	return n.NotifyUser(ctx, ownerUserID, string(NotifyEmployeeAccepted),
		"Employé accepté",
		fmt.Sprintf("%s a accepté l'invitation pour %s.", employeeName, enterpriseName),
		map[string]interface{}{
			"employee_name":   employeeName,
			"enterprise_name": enterpriseName,
		},
	)
}

// NotifyPayrollExecution sends notification to enterprise owner after payroll
func (n *NotificationClient) NotifyPayrollExecution(ctx context.Context, ownerUserID string, totalAmount float64, employeeCount int, period string) error {
	return n.NotifyUser(ctx, ownerUserID, string(NotifyPayrollExecuted),
		"Paie exécutée",
		fmt.Sprintf("La paie de %s a été exécutée. %d employés, %.0f XOF au total.", period, employeeCount, totalAmount),
		map[string]interface{}{
			"total_amount":   totalAmount,
			"employee_count": employeeCount,
			"period":         period,
		},
	)
}

// NotifySalaryPayment sends notification to employee when salary is paid
func (n *NotificationClient) NotifySalaryPayment(ctx context.Context, employeeUserID, enterpriseName string, netAmount float64, period string) error {
	return n.NotifyUser(ctx, employeeUserID, string(NotifySalaryReceived),
		"Salaire reçu",
		fmt.Sprintf("Vous avez reçu votre salaire de %.0f XOF de %s pour %s.", netAmount, enterpriseName, period),
		map[string]interface{}{
			"amount":          netAmount,
			"enterprise_name": enterpriseName,
			"period":          period,
		},
	)
}

// NotifySubscriptionCreated sends notification to subscriber
func (n *NotificationClient) NotifySubscriptionCreated(ctx context.Context, clientUserID, enterpriseName, serviceName string) error {
	return n.NotifyUser(ctx, clientUserID, string(NotifyNewSubscription),
		"Nouvel abonnement",
		fmt.Sprintf("Vous êtes maintenant abonné à %s chez %s.", serviceName, enterpriseName),
		map[string]interface{}{
			"enterprise_name": enterpriseName,
			"service_name":    serviceName,
		},
	)
}

// NotifyInvoice sends invoice notification to client
func (n *NotificationClient) NotifyInvoice(ctx context.Context, clientUserID, enterpriseName, serviceName string, amount float64, dueDate string) error {
	return n.NotifyUser(ctx, clientUserID, string(NotifyInvoiceCreated),
		"Nouvelle facture",
		fmt.Sprintf("%s: Facture de %.0f XOF pour %s. À payer avant le %s.", enterpriseName, amount, serviceName, dueDate),
		map[string]interface{}{
			"enterprise_name": enterpriseName,
			"service_name":    serviceName,
			"amount":          amount,
			"due_date":        dueDate,
			"action":          "pay_invoice",
		},
	)
}

// NotifyPaymentReceived sends confirmation to enterprise when payment is received
func (n *NotificationClient) NotifyPaymentReceived(ctx context.Context, ownerUserID, clientName string, amount float64, serviceName string) error {
	return n.NotifyUser(ctx, ownerUserID, string(NotifyPaymentReceived),
		"Paiement reçu",
		fmt.Sprintf("%s a payé %.0f XOF pour %s.", clientName, amount, serviceName),
		map[string]interface{}{
			"client_name":  clientName,
			"amount":       amount,
			"service_name": serviceName,
		},
	)
}
