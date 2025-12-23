# 🏦 Zekora v2.0 - Application Bancaire Crypto Complète

## 🎯 Nouvelle Architecture avec Routes Préfixées

L'API Gateway utilise maintenant des routes préfixées par service pour éviter les conflits d'URL :

```
/gateway/{service-name}/{endpoint}
```

### 📍 Nouvelle Structure des URLs

| Service | Base URL | Description |
|---------|----------|-------------|
| **Auth** | `/gateway/auth-service/` | Authentification, 2FA, sessions |
| **Wallets** | `/gateway/wallet-service/` | Portefeuilles crypto/fiat |
| **Exchange** | `/gateway/exchange-service/` | Trading, conversions |
| **Cards** | `/gateway/card-service/` | Cartes prépayées |
| **Transfers** | `/gateway/transfer-service/` | Transferts d'argent |
| **Users** | `/gateway/user-service/` | Profils, KYC |
| **Notifications** | `/gateway/notification-service/` | Alertes |

## 🚀 Nouvelles Fonctionnalités v2.0

### ✨ **Conversion Monnaies Fiduciaires**
- **USD ↔ EUR ↔ GBP ↔ JPY ↔ CAD ↔ AUD**
- Taux de change en temps réel
- Frais réduits (0.15-0.25% vs 3.5% banques)
- Conversion instantanée
- Support de 20+ devises

### 💳 **Cartes Prépayées Complètes**
- **Cartes virtuelles** : Instantanées pour achats en ligne
- **Cartes physiques** : Livrées à domicile
- **Support crypto/fiat** : BTC, ETH, USD, EUR
- **Cartes cadeaux** : Envoi par email/SMS
- **Auto-rechargement** : Rechargement automatique
- **Gestion limites** : Journalières, mensuelles, ATM

### 🔄 **Trading Crypto Avancé**
- **Ordres Market** : Achat/vente instantané
- **Ordres Limit** : Prix spécifique
- **Stop-Loss** : Protection des pertes
- **Portfolio tracking** : Suivi performance
- **P2P Trading** : Trading peer-to-peer

## 📋 **Comment Utiliser les Nouvelles Fonctionnalités**

### 🚀 Démarrage Rapide

```bash
# 1. Cloner et configurer
git clone <votre-repo>
cd crypto-bank-app

# 2. Démarrer avec docker-compose v2
docker-compose -f docker-compose-complete.yml up -d

# 3. Accéder à l'application
open http://localhost:3000
```

### 💱 **Conversion Fiat (Ex: USD → EUR)**

```bash
# 1. Obtenir un devis
curl -X POST http://localhost:8080/gateway/exchange-service/fiat/quote \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USD",
    "to_currency": "EUR",
    "amount": 1000.00
  }'

# 2. Exécuter la conversion
curl -X POST http://localhost:8080/gateway/exchange-service/fiat/execute \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_wallet_id": "your_usd_wallet",
    "to_wallet_id": "your_eur_wallet",
    "from_currency": "USD",
    "to_currency": "EUR",
    "amount": 1000.00
  }'

# 3. Convertisseur simple (calculatrice)
curl "http://localhost:8080/gateway/exchange-service/fiat/convert?from=USD&to=EUR&amount=1000"
```

### ₿ **Acheter/Vendre des Cryptos**

```bash
# Acheter du Bitcoin avec USD
curl -X POST http://localhost:8080/gateway/exchange-service/trading/buy \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "pay_currency": "USD",
    "amount": 0.02,
    "order_type": "market"
  }'

# Vendre du Bitcoin pour EUR
curl -X POST http://localhost:8080/gateway/exchange-service/trading/sell \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "receive_currency": "EUR", 
    "amount": 0.01,
    "order_type": "market"
  }'

# Ordre limite d'achat
curl -X POST http://localhost:8080/gateway/exchange-service/trading/limit-order \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USD",
    "to_currency": "BTC",
    "amount": 1000.00,
    "limit_price": 42000.00,
    "order_type": "buy"
  }'
```

### 💳 **Cartes Prépayées**

#### Créer une Carte Virtuelle
```bash
# Carte virtuelle USD
curl -X POST http://localhost:8080/gateway/card-service/virtual/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "USD",
    "cardholder_name": "John Doe",
    "initial_amount": 500.00,
    "source_wallet_id": "your_usd_wallet_id",
    "purpose": "online_shopping"
  }'

# Carte virtuelle Bitcoin
curl -X POST http://localhost:8080/gateway/card-service/virtual/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "cardholder_name": "John Doe",
    "initial_amount": 0.01,
    "source_wallet_id": "your_btc_wallet_id",
    "purpose": "travel"
  }'
```

#### Commander une Carte Physique
```bash
curl -X POST http://localhost:8080/gateway/card-service/physical/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "EUR",
    "cardholder_name": "John Doe",
    "initial_amount": 1000.00,
    "source_wallet_id": "your_eur_wallet_id",
    "shipping_address": {
      "full_name": "John Doe",
      "address_line1": "123 Rue de la Paix",
      "city": "Paris",
      "postal_code": "75001", 
      "country": "FRA",
      "phone_number": "+33123456789"
    },
    "express_shipping": true
  }'
```

#### Recharger une Carte
```bash
# Depuis portefeuille crypto/fiat
curl -X POST http://localhost:8080/gateway/card-service/cards/CARD_ID/load \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 200.00,
    "source_wallet_id": "your_wallet_id",
    "description": "Rechargement vacances"
  }'

# Auto-rechargement
curl -X POST http://localhost:8080/gateway/card-service/cards/CARD_ID/auto-load \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reload_amount": 100.00,
    "reload_threshold": 20.00,
    "source_wallet_id": "your_wallet_id"
  }'
```

### 🎁 **Cartes Cadeaux**

```bash
# Créer une carte cadeau
curl -X POST http://localhost:8080/gateway/card-service/gift/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "currency": "USD",
    "recipient_email": "ami@email.com",
    "message": "Joyeux anniversaire !",
    "design": "birthday",
    "source_wallet_id": "your_wallet_id"
  }'

# Utiliser une carte cadeau
curl -X POST http://localhost:8080/gateway/card-service/gift/redeem \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ABCD-EFGH-IJKL-MNOP",
    "target_wallet_id": "your_wallet_id"
  }'
```

## 🌍 **Support Multi-Devises Complet**

### Cryptomonnaies Supportées
- **Bitcoin (BTC)** - La référence crypto
- **Ethereum (ETH)** - Smart contracts et DeFi
- **Binance Coin (BNB)** - Exchange leader
- **USDT/USDC** - Stablecoins sécurisés
- **Cardano (ADA)** - Blockchain proof-of-stake
- **Et plus** : XRP, DOT, LTC, LINK, BCH, XLM

### Devises Fiduciaires Supportées
- **Majeures** : USD, EUR, GBP, JPY, CAD, AUD
- **Européennes** : CHF, SEK, NOK, DKK, PLN
- **Asiatiques** : CNY, INR, KRW, SGD, HKD
- **Autres** : BRL, MXN, ZAR, TRY

## 📊 **Interface Web Complète**

### Pages Principales
- **Dashboard** (`/dashboard`) - Vue d'ensemble
- **Exchange Fiat** (`/exchange/fiat`) - Conversion devises
- **Exchange Crypto** (`/exchange/crypto`) - Trading crypto  
- **Cartes** (`/cards`) - Gestion cartes prépayées
- **Portefeuilles** (`/wallets`) - Gestion portefeuilles
- **Transferts** (`/transfers`) - Envois d'argent

### Fonctionnalités Interface
- **Convertisseurs temps réel** sur le dashboard
- **Graphiques de marché** crypto et fiat
- **Gestion visuelle des cartes** avec design moderne
- **Historique des transactions** complet
- **Comparaison des frais** avec banques traditionnelles

## 🔧 **Configuration et Déploiement**

### Variables d'Environnement Importantes

```bash
# APIs Crypto
INFURA_API_KEY=your_infura_key
ALCHEMY_API_KEY=your_alchemy_key
COINBASE_API_KEY=your_coinbase_key
BINANCE_API_KEY=your_binance_key

# Cartes (Marqeta ou similaire)
MARQETA_API_KEY=your_marqeta_key
MARQETA_API_SECRET=your_marqeta_secret

# Mobile Money
MTN_API_KEY=your_mtn_key
AIRTEL_API_KEY=your_airtel_key
MPESA_API_KEY=your_mpesa_key

# Transferts internationaux
WISE_API_KEY=your_wise_key
REMITLY_API_KEY=your_remitly_key

# Taux de change
FIXER_API_KEY=your_fixer_key
EXCHANGERATE_API_KEY=your_exchangerate_key
```

### Déploiement Production

```bash
# 1. Configuration production
cp .env.example .env.production
# Éditez .env.production avec vos clés API

# 2. Build et déploiement
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 3. Vérification
curl http://your-domain.com/gateway/services
```

## 🎯 **Cas d'Usage Pratiques**

### 1. **Voyageur International**
```bash
# Convertir BTC → EUR avant voyage
POST /gateway/exchange-service/trading/sell
# Créer carte EUR physique
POST /gateway/card-service/physical/
# Activer pour l'international
PUT /gateway/card-service/cards/:id/limits
```

### 2. **E-commerce Business** 
```bash
# Cartes virtuelles à usage unique
POST /gateway/card-service/virtual/ {"purpose": "single_use"}
# Conversion automatique profits crypto → fiat
POST /gateway/exchange-service/fiat/execute
# Transferts groupés employés
POST /gateway/transfer-service/bulk/
```

### 3. **Trader Crypto**
```bash
# Ordres limite multiples
POST /gateway/exchange-service/trading/limit-order
# Stop-loss protection
POST /gateway/exchange-service/trading/stop-loss
# Portfolio tracking
GET /gateway/exchange-service/trading/portfolio
```

### 4. **Envoi d'Argent Famille**
```bash
# Conversion crypto → fiat local
POST /gateway/exchange-service/fiat/execute
# Transfert mobile money
POST /gateway/transfer-service/mobile/send
# Carte cadeau crypto
POST /gateway/card-service/gift/
```

## 📈 **Avantages Compétitifs**

### 💰 **Frais Réduits**
| Service | Zekora | Banques Trad. | Économies |
|---------|------------|---------------|-----------|
| **Conversion fiat** | 0.15-0.25% | 3.5% + 15€ | ~90% |
| **Trading crypto** | 0.1-0.5% | 1-3% | ~80% |
| **Transfert international** | 1-2% | 3.5% + frais | ~60% |
| **Cartes prépayées** | 1€/mois | 5€/mois + frais | ~70% |

### ⚡ **Rapidité**
- **Conversion fiat** : Instantané
- **Trading crypto** : < 15 minutes
- **Cartes virtuelles** : Immédiat
- **Transferts** : 2-24h (vs 3-7 jours)

### 🌍 **Disponibilité**
- **24/7** : Service non-stop
- **200+ pays** : Couverture mondiale
- **Multi-devises** : 50+ devises supportées
- **Mobile-first** : APIs prêtes pour mobile

## 🔍 **Tests et Monitoring**

### Health Check Complet
```bash
# Gateway principal
curl http://localhost:8080/health

# Services individuels
curl http://localhost:8081/health  # Auth
curl http://localhost:8083/health  # Wallet
curl http://localhost:8085/health  # Exchange
curl http://localhost:8086/health  # Cards
```

### Test Suite Automatique
```bash
# Lancer tous les tests
make test

# Tests spécifiques
make test-fiat-conversion
make test-crypto-trading
make test-card-creation
```

### Monitoring Production
- **Prometheus** : Métriques temps réel
- **Grafana** : Dashboards visuels
- **Alertes** : Notifications incidents
- **Logs centralisés** : ELK Stack

## 🚀 **Roadmap v3.0**

### 🔜 Prochaines Fonctionnalités
- **DeFi Integration** : Yield farming, staking
- **NFT Marketplace** : Achat/vente NFT
- **Prêts crypto** : Prêts garantis par crypto
- **App mobile native** : iOS/Android
- **Trading algorithmique** : Bots automatisés
- **DAO Governance** : Gouvernance décentralisée

### 🌟 **Zekora : La Banque du Futur**

Avec Zekora v2.0, vous avez accès à l'écosystème financier le plus complet :

✅ **Trading crypto** professionnel avec ordres avancés  
✅ **Conversion fiat** instantanée 20+ devises  
✅ **Cartes prépayées** crypto/fiat avec livraison  
✅ **Transferts internationaux** rapides et économiques  
✅ **Mobile Money** Afrique (MTN, Airtel, M-Pesa)  
✅ **Gift Cards** numériques  
✅ **APIs complètes** pour développeurs  
✅ **Interface moderne** Nuxt.js  
✅ **Sécurité bancaire** avec 2FA et chiffrement  

## 📞 **Support et Documentation**

- **Documentation API** : http://localhost:8080/gateway/docs
- **Interface Web** : http://localhost:3000
- **Monitoring** : http://localhost:3001
- **Support** : support@cryptobank.com

**Commencez maintenant** et révolutionnez votre expérience bancaire ! 🚀