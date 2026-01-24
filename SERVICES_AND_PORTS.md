# Expositions des Services (Ports & Domaines)

Ce fichier récapitule les ports internes utilisés par chaque microservice et les URLs/domaines publics pour y accéder.

## 🌍 Domaines Publics

| Application | URL / Domaine | Port Interne |
| :--- | :--- | :---: |
| **Application Web (Frontend)** | `https://app.tech-afm.com` | `3000` |
| **Admin Dashboard** | `https://admin.tech-afm.com` | `3002` |
| **API Gateway (Publique)** | `https://api.app.tech-afm.com` | `8080` |
| **API Admin (Backend)** | `https://api.admin.tech-afm.com` | `8088` |
| **CDN (Assets)** | `https://cdn.tech-afm.com` | `9000` |

---

## 🔌 Services & Ports Internes

Voici la liste des services, leurs ports Docker internes, et leurs URLs publiques via l'API Gateway.

| Service | Port Local (Docker) | URL Publique |
| :--- | :---: | :--- |
| **API Gateway (Kong)** | `8080` | `https://api.app.tech-afm.com` |
| **Kong Admin** | `8001` | *(Non exposé publiquement)* |
| **Frontend** | `3000` | `https://app.tech-afm.com` |
| **Admin Dashboard** | `3002` | `https://admin.tech-afm.com` |
| **Auth Service** | `8081` | `https://api.app.tech-afm.com/auth-service` |
| **Wallet Service** | `8083` | `https://api.app.tech-afm.com/wallet-service` |
| **Transfer Service** | `8084` | `https://api.app.tech-afm.com/transfer-service` |
| **Exchange Service** | `8085` | `https://api.app.tech-afm.com/exchange-service` |
| **Card Service** | `8086` | `https://api.app.tech-afm.com/card-service` |
| **Notification Service** | `8087` | `https://api.app.tech-afm.com/notification-service` |
| **Admin Service** (Backend) | `8088` | `https://api.admin.tech-afm.com/admin-service` |
| **Support Service** | `8089` | `https://api.app.tech-afm.com/support-service` |
| **Ticket Service** | `8090` | `https://api.app.tech-afm.com/ticket-service` |
| **Messaging Service** | `8095` | `https://api.app.tech-afm.com/messaging-service` |
| **Donation Service** | `8096` | `https://api.app.tech-afm.com/donation-service` |
| **Enterprise Service** | `8097` | `https://api.app.tech-afm.com/enterprise-service` |
| **Shop Service** | `8098` | `https://api.app.tech-afm.com/shop-service` |

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

> **Note :** Les routes API sont accessibles via `https://api.app.tech-afm.com/<route>`. Exemple : `https://api.app.tech-afm.com/wallet-service/health`
