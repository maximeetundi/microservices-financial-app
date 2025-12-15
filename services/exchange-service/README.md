# Exchange Service

## Description
Service de conversion de devises et **trading crypto via Binance**.

## Port
`8084`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis (rates cache)
- **Message Broker**: RabbitMQ

## Fonctionnalités

### Conversion Fiat
| Feature | Endpoint | Description |
|---------|----------|-------------|
| Tous les taux | `GET /api/v1/exchange/rates` | Taux actuels |
| Taux spécifique | `GET /api/v1/exchange/rate` | from → to |
| Convertir | `POST /api/v1/exchange/convert` | Conversion fiat |
| Historique | `GET /api/v1/exchange/history` | Mes conversions |

### Trading Crypto 🆕
| Feature | Endpoint | Description |
|---------|----------|-------------|
| Taux crypto | `GET /api/v1/exchange/crypto/rates` | BTC, ETH, etc. |
| Acheter | `POST /api/v1/exchange/crypto/buy` | Achat crypto |
| Vendre | `POST /api/v1/exchange/crypto/sell` | Vente crypto |

## Variables d'Environnement

```bash
PORT=8084
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://admin:pass@localhost:5672/

# Binance
BINANCE_API_KEY=xxx
BINANCE_API_SECRET=xxx
BINANCE_BASE_URL=https://api.binance.com
BINANCE_TEST_MODE=false
```

## Provider Binance Implémenté

```go
type BinanceProvider struct {
    apiKey    string
    apiSecret string
    baseURL   string
    testMode  bool
}

// Méthodes
func (b *BinanceProvider) GetPrice(symbol string) (*PriceResponse, error)
func (b *BinanceProvider) ExecuteTrade(req *TradeRequest) (*TradeResponse, error)
func (b *BinanceProvider) GetOrderStatus(orderID string) (*OrderStatus, error)
func (b *BinanceProvider) CancelOrder(orderID string) error
func (b *BinanceProvider) GetAccountBalances() (map[string]Balance, error)
func (b *BinanceProvider) GetConvertQuote(from, to string, amount float64) (*ConvertQuote, error)
```

## Structure

```
exchange-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   └── exchange_handler.go
    ├── models/
    │   └── exchange.go
    ├── repository/
    │   └── exchange_repository.go
    └── services/
        ├── exchange_service.go
        ├── rates_service.go
        └── binance_provider.go    🆕
```

## Paires Supportées

### Fiat
EUR, USD, GBP, XOF, XAF, NGN, GHS, KES, ZAR

### Crypto
BTC, ETH, USDT, USDC, BNB, SOL, XRP

## Cache des Taux

- Taux fiat: cachés 5 minutes
- Taux crypto: cachés 30 secondes
- Source: APIs externes + Binance

## Événements RabbitMQ

| Exchange | Routing Key | Description |
|----------|-------------|-------------|
| `exchange.events` | `conversion.completed` | Conversion terminée |
| `exchange.events` | `trade.executed` | Trade crypto exécuté |

---
*CryptoBank Exchange Service - v2.0*