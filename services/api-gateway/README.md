# API Gateway

## Description
Point d'entrée unique pour tous les microservices. Gère le routage, l'authentification et le rate limiting.

## Port
`8080`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Auth**: JWT validation

## Fonctionnalités

| Feature | Description |
|---------|-------------|
| Routing | Proxy vers microservices |
| Auth | Validation JWT |
| Rate Limiting | 100 req/min par IP |
| CORS | Cross-Origin configuré |
| Health Check | `/health` |

## Variables d'Environnement

```bash
PORT=8080
JWT_SECRET=your_secret_key

# Services URLs
AUTH_SERVICE_URL=http://auth-service:8081
WALLET_SERVICE_URL=http://wallet-service:8082
TRANSFER_SERVICE_URL=http://transfer-service:8083
EXCHANGE_SERVICE_URL=http://exchange-service:8084
CARD_SERVICE_URL=http://card-service:8085
```

## Routes Proxifiées

| Prefix | Service |
|--------|---------|
| `/api/v1/auth/*` | Auth Service |
| `/api/v1/wallets/*` | Wallet Service |
| `/api/v1/merchant/*` | Wallet Service |
| `/api/v1/payments/*` | Wallet Service |
| `/api/v1/pay/*` | Wallet Service (public) |
| `/api/v1/transfers/*` | Transfer Service |
| `/api/v1/exchange/*` | Exchange Service |
| `/api/v1/cards/*` | Card Service |

## Structure

```
api-gateway/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── middleware/
    │   ├── auth.go
    │   ├── cors.go
    │   └── rate_limiter.go
    └── proxy/
        └── proxy.go
```

## Health Check

```bash
GET /health
Response: {"status": "ok", "service": "api-gateway"}
```

---
*CryptoBank API Gateway - v2.0*
