# 🏦 Zekora - Plateforme Bancaire Numérique Complète

Une application bancaire moderne et sécurisée supportant les cryptomonnaies et les devises fiduciaires, avec des fonctionnalités de transfert mobile money et de cartes prépayées.

## 🌟 Fonctionnalités

### 💰 Gestion des Portefeuilles
- ✅ Portefeuilles crypto (BTC, ETH, BSC)
- ✅ Portefeuilles fiat (USD, EUR, etc.)
- ✅ Génération d'adresses crypto sécurisées
- ✅ Surveillance des transactions blockchain
- ✅ Gel/dégel de portefeuilles

### 🔄 Transferts d'Argent
- ✅ Transferts domestiques instantanés
- ✅ Transferts internationaux (SWIFT)
- ✅ Mobile Money (MTN, Airtel, M-Pesa)
- ✅ Transferts crypto peer-to-peer
- ✅ Transferts groupés pour entreprises

### 💳 Cartes Prépayées
- ✅ Cartes virtuelles instantanées
- ✅ Cartes physiques
- ✅ Gestion des limites
- ✅ Rechargement automatique

### 🔐 Sécurité Avancée
- ✅ Authentification à deux facteurs (TOTP)
- ✅ Chiffrement end-to-end
- ✅ KYC/AML avec niveaux de vérification
- ✅ Détection de fraude en temps réel
- ✅ Audit trails complets

### 🌍 Conformité Mondiale
- ✅ Support multi-devises
- ✅ Conformité réglementaire
- ✅ Vérification des sanctions
- ✅ Rapports de transactions

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway    │    │   Auth Service  │
│   (Nuxt.js)     │◄──►│   (Go)           │◄──►│   (Go)          │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                       ┌────────┼────────┐
                       ▼        ▼        ▼
               ┌─────────────┬─────────────┬─────────────┐
               │   Wallet    │  Transfer   │  Exchange   │
               │   Service   │   Service   │   Service   │
               │   (Go)      │   (Go)      │   (Go)      │
               └─────────────┴─────────────┴─────────────┘
                       │        │        │
               ┌───────┼────────┼────────┼───────┐
               ▼       ▼        ▼        ▼       ▼
           ┌─────────┬─────────┬─────────┬─────────┬─────────┐
           │PostgreSQL│  Redis  │RabbitMQ │Prometheus│ Grafana │
           └─────────┴─────────┴─────────┴─────────┴─────────┘
```

### Services
- **API Gateway**: Point d'entrée unique, routage et sécurité
- **Auth Service**: Authentification, autorisation, 2FA
- **Wallet Service**: Gestion portefeuilles crypto/fiat
- **Transfer Service**: Transferts domestiques/internationaux
- **Exchange Service**: Conversion crypto/fiat
- **Card Service**: Cartes prépayées virtuelles/physiques
- **Notification Service**: Emails, SMS, push notifications

## 🚀 Démarrage Rapide

### Prérequis
- Docker et Docker Compose
- 8GB RAM minimum
- Ports disponibles: 3000, 8080-8087, 5432, 6379, 15672

### Installation

```bash
# Cloner le projet
git clone <repository-url>
cd crypto-bank-app

# Rendre le script exécutable
chmod +x start.sh

# Démarrer l'application
./start.sh
```

### Accès aux Services
- **Application Web**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Monitoring Grafana**: http://localhost:3001 (admin/admin)
- **RabbitMQ Management**: http://localhost:15672 (admin/secure_password)

## 📋 Configuration

### Variables d'Environnement

Créez un fichier `.env` dans le répertoire racine :

```bash
# Database
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password

# JWT
JWT_SECRET=your_ultra_secure_jwt_secret_minimum_32_chars

# Email (optionnel)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# SMS (optionnel)
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
TWILIO_FROM_NUMBER=+1234567890

# Crypto APIs (optionnel)
INFURA_API_KEY=your_infura_key
ALCHEMY_API_KEY=your_alchemy_key

# Mobile Money (optionnel)
MTN_API_KEY=your_mtn_api_key
AIRTEL_API_KEY=your_airtel_api_key
MPESA_API_KEY=your_mpesa_api_key
```

### Limites de Transaction par Défaut

```
Niveau KYC 1: 1,000 USD/jour, 5,000 USD/mois
Niveau KYC 2: 10,000 USD/jour, 50,000 USD/mois  
Niveau KYC 3: 100,000 USD/jour, 500,000 USD/mois
```

## 📱 Guide d'Utilisation

### 1. Inscription et KYC
1. Créez un compte sur http://localhost:3000
2. Vérifiez votre email
3. Complétez votre profil KYC
4. Activez l'authentification 2FA

### 2. Création de Portefeuilles
```javascript
// Créer un portefeuille fiat
POST /api/v1/wallets
{
  "currency": "USD",
  "wallet_type": "fiat",
  "name": "Mon Portefeuille USD"
}

// Créer un portefeuille crypto
POST /api/v1/wallets
{
  "currency": "BTC", 
  "wallet_type": "crypto",
  "name": "Mon Portefeuille Bitcoin"
}
```

### 3. Transferts
```javascript
// Transfert domestique
POST /api/v1/transfers
{
  "from_wallet_id": "wallet-id",
  "to_email": "destinataire@email.com",
  "amount": 100.00,
  "currency": "USD",
  "description": "Transfert test"
}

// Transfert Mobile Money
POST /api/v1/transfers/mobile/send
{
  "from_wallet_id": "wallet-id",
  "to_phone": "+233241234567",
  "provider": "mtn",
  "amount": 50.00,
  "currency": "USD",
  "country": "GHA"
}
```

### 4. Exchange Crypto/Fiat
```javascript
// Obtenir un devis
POST /api/v1/exchange/quote
{
  "from_currency": "BTC",
  "to_currency": "USD", 
  "amount": 0.01
}

// Exécuter l'échange
POST /api/v1/exchange/execute
{
  "from_wallet_id": "btc-wallet-id",
  "to_wallet_id": "usd-wallet-id",
  "amount": 0.01,
  "quote_id": "quote-id"
}
```

## 🔧 Développement

### Structure du Projet
```
crypto-bank-app/
├── services/
│   ├── api-gateway/       # Point d'entrée API
│   ├── auth-service/      # Authentification
│   ├── wallet-service/    # Gestion portefeuilles
│   ├── transfer-service/  # Transferts d'argent
│   └── exchange-service/  # Exchange crypto/fiat
├── frontend/              # Interface Nuxt.js
├── infrastructure/        # Configuration DB
└── docker-compose.yml     # Orchestration
```

### Commandes de Développement
```bash
# Démarrer en mode développement
docker-compose -f docker-compose.development.yml up

# Voir les logs d'un service
docker-compose logs -f auth-service

# Redémarrer un service
docker-compose restart wallet-service

# Accéder à la base de données
docker-compose exec postgres psql -U admin crypto_bank

# Monitoring des services
docker-compose logs -f | grep ERROR
```

### Tests
```bash
# Tests d'intégration
cd services/auth-service
go test ./...

# Tests de charge
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

## 🛡️ Sécurité

### Mesures Implémentées
- Chiffrement AES-256 pour les clés privées
- Rate limiting (100 req/min par IP)
- Headers de sécurité HTTP
- Validation stricte des entrées
- Audit trails complets
- Isolation des services
- Secrets management

### Bonnes Pratiques
```bash
# Changer les mots de passe par défaut
# Activer HTTPS en production
# Configurer les firewalls
# Monitorer les logs de sécurité
# Sauvegardes régulières chiffrées
```

## 📊 Monitoring

### Métriques Disponibles
- Transactions par seconde
- Temps de réponse des APIs
- Taux d'erreur par service
- Utilisation des ressources
- Soldes des portefeuilles
- Activité utilisateur

### Alertes Configurées
- Service indisponible
- Pic de transactions échouées
- Utilisation mémoire/CPU élevée
- Tentatives de connexion suspectes
- Transactions importantes

## 🌍 Déploiement Production

### Prérequis Production
```yaml
# Infrastructure recommandée
CPU: 8 cores minimum
RAM: 32GB minimum  
Storage: 500GB SSD
Network: 1Gbps
OS: Ubuntu 22.04 LTS
```

### Kubernetes (Recommandé)
```bash
# Utiliser les charts Helm fournis
cd helm
helm install crypto-bank ./crypto-bank-chart
```

### Configuration Production
```bash
# Variables d'environnement critiques
ENVIRONMENT=production
JWT_SECRET=<256-bit-random-key>
ENCRYPTION_KEY=<256-bit-random-key>
DB_URL=<production-db-url>

# TLS/SSL obligatoire
ENABLE_TLS=true
TLS_CERT_PATH=/path/to/cert.pem
TLS_KEY_PATH=/path/to/key.pem
```

## 🤝 Support et Contribution

### Support Technique
- 📧 Email: support@cryptobank.com
- 📱 Téléphone: +1-800-CRYPTO-BANK
- 💬 Discord: [Serveur Communauté]
- 📚 Documentation: https://docs.cryptobank.com

### Contribution
1. Fork le projet
2. Créez une branche feature (`git checkout -b feature/nouvelle-fonctionnalite`)
3. Commitez vos changements (`git commit -am 'Ajout nouvelle fonctionnalité'`)
4. Poussez la branche (`git push origin feature/nouvelle-fonctionnalite`)
5. Créez une Pull Request

## 📜 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🔄 Changelog

### Version 1.0.0 (2024-01-15)
- ✅ Lancement initial
- ✅ Support crypto (BTC, ETH, BSC)
- ✅ Transferts internationaux
- ✅ Mobile Money Afrique
- ✅ Interface web responsive
- ✅ APIs REST complètes

### Prochaines Versions
- 📱 Application mobile (iOS/Android)
- 🔄 DeFi integration (staking, yield farming)
- 💎 NFT marketplace
- 🤖 Trading algorithmique
- 🌐 Support multi-langues

---

**⚠️ Avertissement**: Cette application est fournie à des fins éducatives et de démonstration. Pour un usage en production, assurez-vous de respecter toutes les réglementations financières locales et d'effectuer un audit de sécurité complet.