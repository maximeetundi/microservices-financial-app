# Card Service

## Description
Service de gestion des cartes bancaires virtuelles et physiques.

## Port
`8085`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Message Broker**: RabbitMQ

## Fonctionnalités

| Feature | Endpoint | Description |
|---------|----------|-------------|
| Lister cartes | `GET /api/v1/cards` | Mes cartes |
| Détails | `GET /api/v1/cards/:id` | Info carte |
| Créer | `POST /api/v1/cards` | Nouvelle carte |
| Activer | `POST /api/v1/cards/:id/activate` | Activer carte |
| Geler | `POST /api/v1/cards/:id/freeze` | Bloquer temporairement |
| Dégeler | `POST /api/v1/cards/:id/unfreeze` | Débloquer |
| Définir PIN | `POST /api/v1/cards/:id/pin` | Changer PIN |
| Limites | `POST /api/v1/cards/:id/limits` | Modifier limites |
| Transactions | `GET /api/v1/cards/:id/transactions` | Historique |
| Recharger | `POST /api/v1/cards/:id/topup` | Ajouter fonds |

## Variables d'Environnement

```bash
PORT=8085
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
RABBITMQ_URL=amqp://admin:pass@localhost:5672/
JWT_SECRET=your_secret_key

# Card Provider (à configurer selon provider)
CARD_PROVIDER_API_KEY=xxx
CARD_PROVIDER_BASE_URL=https://api.cardprovider.com
```

## Types de Cartes

| Type | Description | Fonctionnalités |
|------|-------------|-----------------|
| `virtual` | Carte dématérialisée | Paiements en ligne |
| `physical` | Carte plastique | POS, DAB, Online |

## Structure

```
card-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   └── card_handler.go
    ├── models/
    │   └── card.go
    ├── repository/
    │   └── card_repository.go
    └── services/
        └── card_service.go
```

## Modèle Carte

```go
type Card struct {
    ID              string
    UserID          string
    WalletID        string
    CardNumber      string    // Chiffré
    CVV             string    // Chiffré
    ExpiryMonth     int
    ExpiryYear      int
    CardType        string    // virtual, physical
    Status          string    // pending, active, frozen, cancelled
    SpendingLimit   float64
    DailyLimit      float64
    MonthlyLimit    float64
    Currency        string
}
```

## Limites par Défaut

| Type | Quotidien | Mensuel |
|------|-----------|---------|
| Virtual | €1,000 | €5,000 |
| Physical | €2,500 | €10,000 |

## Événements RabbitMQ

| Exchange | Routing Key | Description |
|----------|-------------|-------------|
| `card.events` | `card.created` | Nouvelle carte |
| `card.events` | `card.activated` | Carte activée |
| `card.events` | `card.frozen` | Carte gelée |
| `card.events` | `card.transaction` | Transaction carte |

## Sécurité

- Numéros de carte chiffrés (AES-256)
- CVV non stocké après activation
- PIN hashé (bcrypt)
- 3 tentatives PIN max

---
*CryptoBank Card Service - v2.0*
