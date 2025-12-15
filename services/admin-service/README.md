# Admin Service

## Description
Service d'administration de la plateforme CryptoBank.

## Port
`8088`

## Technologies
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL (séparée)
- **Cache**: Redis
- **Message Broker**: RabbitMQ

## Fonctionnalités

| Feature | Endpoint | Description |
|---------|----------|-------------|
| Dashboard | `GET /api/v1/admin/dashboard` | Statistiques |
| Utilisateurs | `GET /api/v1/admin/users` | Liste users |
| Détail user | `GET /api/v1/admin/users/:id` | Info user |
| KYC | `POST /api/v1/admin/users/:id/kyc` | Valider KYC |
| Bloquer | `POST /api/v1/admin/users/:id/block` | Bloquer user |
| Transactions | `GET /api/v1/admin/transactions` | Toutes transactions |
| Admins | `GET /api/v1/admin/admins` | Liste admins |
| Créer admin | `POST /api/v1/admin/admins` | Nouvel admin |
| Rôles | `GET /api/v1/admin/roles` | Rôles disponibles |
| Logs | `GET /api/v1/admin/audit-logs` | Logs d'audit |

## Variables d'Environnement

```bash
PORT=8088
ADMIN_DB_URL=postgres://user:pass@localhost:5432/crypto_bank_admin
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://admin:pass@localhost:5672/
JWT_SECRET=admin_jwt_secret
```

## Admin par Défaut

Au démarrage, un super admin est créé:
- **Email**: `admin@cryptobank.com`
- **Password**: `Admin123!`

## Rôles

| Rôle | Permissions |
|------|-------------|
| `super_admin` | Toutes |
| `admin` | Users, Transactions, KYC |
| `support` | Lecture seule, Tickets |
| `compliance` | KYC, AML |

## Structure

```
admin-service/
├── main.go
├── go.mod
├── Dockerfile
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── database.go
    ├── handlers/
    │   ├── admin_handler.go
    │   ├── user_handler.go
    │   └── dashboard_handler.go
    ├── middleware/
    │   └── admin_auth.go
    ├── models/
    │   ├── admin.go
    │   └── audit_log.go
    ├── repository/
    │   ├── admin_repository.go
    │   └── audit_repository.go
    └── services/
        ├── admin_service.go
        └── audit_service.go
```

## Audit Logs

Toutes les actions admin sont loguées:
```go
type AuditLog struct {
    ID        string
    AdminID   string
    Action    string    // user.blocked, kyc.approved, etc.
    TargetID  string
    Details   map[string]interface{}
    IPAddress string
    CreatedAt time.Time
}
```

---
*CryptoBank Admin Service - v2.0*
