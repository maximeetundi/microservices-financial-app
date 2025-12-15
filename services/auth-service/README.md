# Auth Service

## Description
Service d'authentification et gestion des sessions pour CryptoBank.

## Port
`8081`

## Technologies
- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis (sessions)
- **Message Broker**: RabbitMQ

## Fonctionnalités

| Feature | Endpoint | Description |
|---------|----------|-------------|
| Registration | `POST /api/v1/register` | Inscription utilisateur |
| Login | `POST /api/v1/login` | Connexion + JWT |
| Refresh Token | `POST /api/v1/refresh` | Renouveler access token |
| Logout | `POST /api/v1/logout` | Déconnexion |
| 2FA Setup | `POST /api/v1/enable-2fa` | Activer TOTP |
| 2FA Verify | `POST /api/v1/verify-2fa` | Valider 2FA |
| Sessions | `GET /api/v1/sessions` | Lister sessions actives |

## Variables d'Environnement

```bash
PORT=8081
DB_URL=postgres://user:pass@localhost:5432/crypto_bank
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://admin:pass@localhost:5672/
JWT_SECRET=your_secret_key
```

## Dépendances Go

```go
github.com/gin-gonic/gin
github.com/golang-jwt/jwt/v5
github.com/lib/pq
github.com/go-redis/redis/v8
github.com/pquerna/otp          // TOTP 2FA
github.com/streadway/amqp       // RabbitMQ
golang.org/x/crypto             // bcrypt
```

## Structure

```
auth-service/
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go         # PostgreSQL, Redis, RabbitMQ
    ├── handlers/
    │   └── auth_handler.go
    ├── middleware/
    │   └── auth.go             # JWT middleware
    ├── models/
    │   └── user.go
    ├── repository/
    │   ├── user_repository.go
    │   └── session_repository.go
    └── services/
        ├── auth_service.go
        ├── email_service.go
        ├── sms_service.go
        ├── totp_service.go
        └── audit_service.go
```

## Sécurité

- JWT access token: 15 minutes
- JWT refresh token: 7 jours
- Password: bcrypt (cost 12)
- Sessions: stockées Redis avec TTL
- Rate limiting: 100 req/min

## Événements RabbitMQ

| Exchange | Routing Key | Description |
|----------|-------------|-------------|
| `auth.events` | `user.registered` | Nouvelle inscription |
| `auth.events` | `user.login` | Connexion réussie |
| `auth.events` | `user.logout` | Déconnexion |
| `audit.events` | `security.login.failed` | Tentative échouée |

---
*CryptoBank Auth Service - v2.0*
