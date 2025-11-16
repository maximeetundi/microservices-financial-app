# Crypto Bank Application

Une application bancaire complète permettant:
- Transferts d'argent mobile
- Gestion des cryptomonnaies
- Conversion crypto/fiat
- Paiement par cartes prépayées
- Sécurité avancée
- Scalabilité mondiale

## Architecture

### Frontend
- **Web**: Nuxt.js avec TypeScript
- **Mobile**: À développer après le web

### Backend
- **API Gateway**: Go avec Gin/Fiber
- **Microservices**: Go
- **Base de données**: PostgreSQL + Redis
- **Blockchain**: Intégration Web3
- **Sécurité**: JWT, 2FA, encryption

### Services principaux
1. **Auth Service** - Authentification et autorisation
2. **User Service** - Gestion des utilisateurs
3. **Wallet Service** - Portefeuilles crypto/fiat
4. **Transfer Service** - Transferts d'argent
5. **Exchange Service** - Conversion crypto/fiat
6. **Card Service** - Cartes prépayées
7. **Compliance Service** - KYC/AML
8. **Notification Service** - Alertes et notifications

## Sécurité
- Chiffrement end-to-end
- 2FA obligatoire
- KYC/AML compliance
- Audit trails
- Rate limiting
- DDoS protection