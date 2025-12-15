# Notification Service

## Description
Service de notifications **multi-canal** (Email, SMS, Push) avec sélection intelligente.

## Port
`8086`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Message Broker**: RabbitMQ (consumer)

## Canaux Implémentés

| Canal | Provider | Config |
|-------|----------|--------|
| **Email** | SMTP (TLS) | Host, User, Password |
| **SMS** | Twilio | Account SID, Auth Token |
| **Push** | Firebase FCM | Server Key |

## Variables d'Environnement

```bash
PORT=8086
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
RABBITMQ_URL=amqp://admin:pass@localhost:5672/

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@cryptobank.com
SMTP_PASSWORD=app_password

# SMS Twilio
TWILIO_ACCOUNT_SID=ACxxxxxx
TWILIO_AUTH_TOKEN=xxxxxx
TWILIO_FROM_NUMBER=+1234567890

# Push FCM
FCM_SERVER_KEY=xxxxxx
```

## Logique de Sélection des Canaux

### Par Priorité
| Priorité | Canaux Activés | Exemples |
|----------|----------------|----------|
| **Critical** | SMS + Email + Push | Wallet gelé, Login suspect |
| **High** | Email + Push | Transfert échoué, Argent reçu |
| **Normal** | Email + Push | Transfert complété |
| **Low** | Push uniquement | Limite de taux atteinte |

### Règles SMS (Canal Critique)
SMS réservé pour:
- `wallet.frozen` - Portefeuille gelé
- `security.suspicious_login` - Connexion suspecte
- `transfer.failed` - Transfert échoué
- `card.blocked` - Carte bloquée
- `money.received` - Argent reçu (montant significatif)
- `money.sent` - Argent envoyé (confirmation)

## Structure

```
notification-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   └── notification_handler.go
    ├── models/
    │   └── notification.go
    ├── repository/
    │   └── notification_repository.go
    └── services/
        └── notification_service.go
```

## Modèle de Notification

```go
type Notification struct {
    ID        string
    UserID    string
    Type      string              // transfer.completed, wallet.frozen, etc.
    Title     string
    Message   string
    Channels  []string            // ["email", "sms", "push"]
    Priority  NotificationPriority // Low, Normal, High, Critical
    Data      map[string]interface{}
    Status    string              // pending, sent, failed
    SentAt    *time.Time
}
```

## Événements RabbitMQ Consommés

| Exchange | Routing Key | Action |
|----------|-------------|--------|
| `notification.events` | `notification.send` | Envoyer notification |
| `transfer.events` | `transfer.*` | Notifier expéditeur |
| `wallet.events` | `wallet.frozen` | SMS + Email critique |
| `auth.events` | `security.*` | Alertes sécurité |
| `payment.events` | `payment.completed` | Notifier marchand |

## Templates Email

```
Subject: {title}

Bonjour,

{message}

---
CryptoBank - Votre banque digitale
```

## Limitations SMS

- Message tronqué à 160 caractères
- Tarification Twilio par message
- Réservé aux événements critiques uniquement

---
*CryptoBank Notification Service - v2.0*
