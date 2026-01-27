package messaging

// Topic constants for Kafka
// Centralized topic definitions for consistency across services
const (
	// User topics
	TopicUserEvents = "user.events"

	// Wallet topics
	TopicWalletEvents = "wallet.events"

	// Transfer topics
	TopicTransferEvents = "transfer.events"

	// Exchange topics
	TopicExchangeEvents     = "exchange.events"
	TopicFiatExchangeEvents = "fiat_exchange.events"

	// Payment topics (for inter-service communication)
	TopicPaymentEvents = "payment.events"

	// Card topics
	TopicCardEvents = "card.events"

	// Notification topics
	TopicNotificationEvents = "notification.events"

	// Transaction topics
	TopicTransactionEvents = "transaction.events"

	// Trading topics
	TopicTradingEvents = "trading.events"
	TopicRateUpdates   = "rate.updates"

	// Audit topics
	TopicAuditEvents = "audit.events"
)

// Event type constants
const (
	// User event types
	EventUserRegistered = "user.registered"
	EventUserUpdated    = "user.updated"

	// Wallet event types
	EventWalletCreated        = "wallet.created"
	EventWalletBalanceUpdated = "wallet.balance_updated"
	EventWalletDebited        = "wallet.debited"
	EventWalletCredited       = "wallet.credited"

	// Transfer event types
	EventTransferInitiated = "transfer.initiated"
	EventTransferCompleted = "transfer.completed"
	EventTransferFailed    = "transfer.failed"

	// Exchange event types
	EventExchangeInitiated = "exchange.initiated"
	EventExchangeCompleted = "exchange.completed"
	EventExchangeFailed    = "exchange.failed"
	EventFiatExchangeRequest = "fiat_exchange.request"

	// Payment event types
	EventPaymentRequest = "payment.request"
	EventPaymentSuccess = "payment.status.success"
	EventPaymentFailed  = "payment.status.failed"

	// Card event types
	EventCardCreated     = "card.created"
	EventCardLoaded      = "card.loaded"
	EventCardTransaction = "card.transaction"
	EventCardBlocked     = "card.blocked"

	// Notification event types
	EventNotificationCreated = "notification.created"
)

// Consumer group IDs
const (
	GroupWalletService       = "wallet-service-group"
	GroupTransferService     = "transfer-service-group"
	GroupExchangeService     = "exchange-service-group"
	GroupNotificationService = "notification-service-group"
	GroupCardService         = "card-service-group"
	GroupTicketService       = "ticket-service-group"
	GroupAdminService        = "admin-service-group"
)
