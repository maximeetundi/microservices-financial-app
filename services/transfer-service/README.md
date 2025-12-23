# Transfer Service

## Description
Service de transferts internationaux multi-partenaires avec **crypto rails invisibles**.

## Port
`8083`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Message Broker**: RabbitMQ

## Architecture Multi-Partenaires

```
┌─────────────────────────────────────────────────────────────────┐
│                    FULL TRANSFER SERVICE                         │
├─────────────────────────────────────────────────────────────────┤
│  1. Collection (collecte fonds)                                  │
│  2. Internal Pool ou Crypto Rails                                │
│  3. Payout (livraison destinataire)                             │
└─────────────────────────────────────────────────────────────────┘
```

## Providers Implémentés

### Payout Providers
| Provider | Zone | Méthodes |
|----------|------|----------|
| **Flutterwave** | Afrique | Mobile Money, Bank Transfer |
| **Thunes** | Global | E-wallets, Cash Pickup |
| **Stripe** | EU/US | SEPA, ACH, Wire |

### Crypto Rails
| Provider | Fonction |
|----------|----------|
| **Circle** | Mint/Burn USDC, Treasury |
| **Binance** | Rates conversion |

## Fonctionnalités

| Feature | Endpoint | Description |
|---------|----------|-------------|
| Créer transfert | `POST /api/v1/transfers` | Nouveau transfert |
| Lister | `GET /api/v1/transfers` | Mes transferts |
| Détails | `GET /api/v1/transfers/:id` | Info transfert |
| Banques | `GET /api/v1/transfers/banks` | Banques par pays |
| Opérateurs mobile | `GET /api/v1/transfers/mobile-operators` | Opérateurs MM |
| Frais | `GET /api/v1/transfers/fees` | Calcul frais |
| Valider | `POST /api/v1/transfers/validate-recipient` | Vérifier destinataire |

## Variables d'Environnement

```bash
PORT=8083
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
RABBITMQ_URL=amqp://admin:pass@localhost:5672/

# Providers
FLUTTERWAVE_SECRET_KEY=FLWSECK_xxx
FLUTTERWAVE_PUBLIC_KEY=FLWPUBK_xxx
THUNES_API_KEY=xxx
THUNES_API_SECRET=xxx
STRIPE_SECRET_KEY=sk_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# Crypto Rails
CIRCLE_API_KEY=xxx
CRYPTO_POOL_THRESHOLD=2500    # USD - seuil pour crypto rails
```

## Flux de Transfert

### Petits montants (< $2500)
```
Expéditeur → Pool Interne → Destinataire
              (instantané)
```

### Gros montants (> $2500)
```
Expéditeur → Lock Fonds → USDC Mint → Conversion → Payout Provider → Destinataire
                         (Circle)    (Binance)   (Flutterwave/etc)
```

## Structure

```
transfer-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   └── transfer_handler.go
    ├── models/
    │   └── transfer.go
    ├── providers/
    │   ├── provider.go           # Interface PayoutProvider
    │   ├── flutterwave.go        # Afrique
    │   ├── thunes.go             # Global
    │   ├── stripe.go             # EU/US
    │   ├── router.go             # ZoneRouter
    │   ├── config.go             # ProvidersConfig
    │   ├── collection.go         # CollectionProvider
    │   ├── crypto_rails.go       # CryptoRailsProvider (Circle)
    │   ├── internal_transfer.go  # InternalTransferService
    │   ├── full_service.go       # FullTransferService
    │   └── orchestrator.go       # Orchestration
    ├── repository/
    │   └── transfer_repository.go
    └── services/
        └── transfer_service.go
```

## Interface PayoutProvider

```go
type PayoutProvider interface {
    GetName() string
    GetSupportedCountries() []string
    GetAvailableMethods(ctx context.Context, country string) ([]AvailableMethod, error)
    GetBanks(ctx context.Context, country string) ([]Bank, error)
    GetMobileOperators(ctx context.Context, country string) ([]MobileOperator, error)
    CreatePayout(ctx context.Context, req *PayoutRequest) (*PayoutResponse, error)
    GetPayoutStatus(ctx context.Context, referenceID string) (*PayoutStatusResponse, error)
}
```

## Zones de Routage

| Zone | Pays | Provider |
|------|------|----------|
| AFRICA | NG, GH, KE, SN, CI, ... | Flutterwave |
| EUROPE | FR, DE, ES, IT, ... | Stripe |
| USA | US, CA | Stripe |
| ASIA | IN, PH, VN, ... | Thunes |
| LATAM | MX, BR, CO, ... | Thunes |

## Événements RabbitMQ

| Exchange | Routing Key | Description |
|----------|-------------|-------------|
| `transfer.events` | `transfer.initiated` | Transfert créé |
| `transfer.events` | `transfer.processing` | En cours |
| `transfer.events` | `transfer.completed` | Terminé |
| `transfer.events` | `transfer.failed` | Échec |

---
*Zekora Transfer Service - v2.0*
