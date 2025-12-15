# Wallet Service

## Description
Service de gestion des portefeuilles, transactions et **paiements marchands QR**.

## Port
`8082`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Broker**: RabbitMQ
- **QR Code**: boombuler/barcode

## Fonctionnalités

### Portefeuilles
| Feature | Endpoint | Description |
|---------|----------|-------------|
| Liste | `GET /api/v1/wallets` | Portefeuilles utilisateur |
| Détails | `GET /api/v1/wallets/:id` | Info portefeuille |
| Créer | `POST /api/v1/wallets` | Nouveau portefeuille |
| Balance | `GET /api/v1/wallets/:id/balance` | Solde actuel |
| Transactions | `GET /api/v1/wallets/:id/transactions` | Historique |

### Paiements Marchands 🆕
| Feature | Endpoint | Description |
|---------|----------|-------------|
| Créer demande | `POST /api/v1/merchant/payments` | Nouvelle demande QR |
| Paiement rapide | `POST /api/v1/merchant/quick-pay` | QR simplifié |
| Lister demandes | `GET /api/v1/merchant/payments` | Mes demandes |
| Historique | `GET /api/v1/merchant/payments/history` | Paiements reçus |
| Détails (public) | `GET /api/v1/pay/:id` | Info paiement |
| Obtenir QR | `GET /api/v1/payments/:id/qr` | Image QR Base64 |
| Payer | `POST /api/v1/payments/:id/pay` | Effectuer paiement |

## Variables d'Environnement

```bash
PORT=8082
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://admin:pass@localhost:5672/
JWT_SECRET=your_secret_key
BASE_URL=https://app.cryptobank.com    # Pour liens de paiement
```

## Types de Paiements Marchands

| Type | Champs | Cas d'usage |
|------|--------|-------------|
| `fixed` | `amount` obligatoire | Produit, service fixe |
| `variable` | `min_amount`, `max_amount` optionnels | Pourboire, donation |
| `invoice` | `items[]` avec détails | Facture détaillée |

### Exemple de Création
```json
POST /api/v1/merchant/payments
{
  "type": "fixed",
  "wallet_id": "wallet_123",
  "amount": 25.00,
  "currency": "EUR",
  "title": "Café + Croissant",
  "description": "Merci pour votre achat!",
  "expires_in_minutes": 60,
  "reusable": false
}
```

### Réponse
```json
{
  "payment_request": {
    "id": "pay_abc123",
    "payment_link": "https://app.cryptobank.com/pay/pay_abc123",
    "qr_code_data": "...",
    "expires_at": "2024-12-15T20:00:00Z"
  },
  "qr_code_base64": "data:image/png;base64,..."
}
```

## Structure

```
wallet-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   ├── wallet_handler.go
    │   └── merchant_payment_handler.go  🆕
    ├── models/
    │   ├── wallet.go
    │   └── payment_request.go           🆕
    ├── repository/
    │   ├── wallet_repository.go
    │   ├── transaction_repository.go
    │   └── payment_request_repository.go 🆕
    └── services/
        ├── wallet_service.go
        ├── balance_service.go
        ├── crypto_service.go
        └── merchant_payment_service.go   🆕
```

## Frais Marchands

| Type | Pourcentage |
|------|-------------|
| Paiement reçu | 1.5% du montant |

## Événements RabbitMQ

| Exchange | Routing Key | Description |
|----------|-------------|-------------|
| `wallet.events` | `wallet.created` | Nouveau portefeuille |
| `wallet.events` | `transaction.completed` | Transaction complète |
| `payment.events` | `payment.completed` | Paiement marchand reçu |

---
*CryptoBank Wallet Service - v2.0*
