# Expositions des Services (Ports & Domaines)

Ce fichier récapitule les ports internes utilisés par chaque microservice et les URLs/domaines publics pour y accéder.

## 🌍 Domaines Publics

| Application | URL / Domaine |
| :--- | :--- |
| **Application Web (Frontend)** | `https://app.maximeetundi.store` |
| **Admin Dashboard** | `https://admin.maximeetundi.store` |
| **API Gateway (Publique)** | `https://api.app.maximeetundi.store` |
| **API Admin (Gateway)** | `https://api.admin.maximeetundi.store` |
| **CDN (Assets)** | `https://cdn.maximeetundi.store` |

---

## 🔌 Services & Ports Internes

Voici la liste des services, leurs ports Docker internes, et leurs routes via l'API Gateway (Kong).

| Service | Port Interne | Route API Gateway |
| :--- | :---: | :--- |
| **API Gateway (Kong)** | `8080` | `/` |
| **Kong Admin** | `8001` | *(Non exposé publiquement)* |
| **Frontend** | `3000` | - |
| **Admin Dashboard** | `3002` | - |
| **Auth Service** | `8081` | `/auth-service` |
| **Wallet Service** | `8083` | `/wallet-service` |
| **Transfer Service** | `8084` | `/transfer-service` |
| **Exchange Service** | `8085` | `/exchange-service` |
| **Card Service** | `8086` | `/card-service` |
| **Notification Service** | `8087` | `/notification-service` |
| **Admin Service** (Backend) | `8088` | `/admin-service` |
| **Support Service** | `8089` | `/support-service` |
| **Ticket Service** | `8090` | `/ticket-service` |
| **Messaging Service** | `8095` | `/messaging-service` |
| **Donation Service** | `8096` | `/donation-service` |
| **Enterprise Service** | `8097` | `/enterprise-service` |
| **Shop Service** | `8098` | `/shop-service` |

---

## 🗄️ Infrastructure & Bases de Données

| Service | Port Interne | Notes |
| :--- | :---: | :--- |
| **PostgreSQL** | `5432` | Base de données principale |
| **MongoDB** | `27017` | NoSQL (Messages, Tickets, Shops) |
| **Redis** | `6379` | Cache & Sessions |
| **RabbitMQ** | `5672` | `15672` (Management UI) |
| **Kafka** | `9092` | Event Bus |
| **MinIO API** | `9000` | S3 Compatible Storage |
| **MinIO Console** | `9001` | Interface d'administration |
| **Prometheus** | `9090` | Monitoring |
| **Grafana** | `3001` | Visualisation (`:3000` interne) |

> **Note :** Les routes API sont accessibles via `https://api.app.maximeetundi.store/<route>`. Exemple : `https://api.app.maximeetundi.store/wallet-service/health`
