# 🏦 Zekora - Guide Complet d'Utilisation

## 🎯 Application Bancaire Crypto Complète

Zekora est une plateforme bancaire numérique moderne qui combine les services bancaires traditionnels avec les cryptomonnaies. Vous pouvez acheter/vendre des cryptos, utiliser des cartes prépayées, faire des transferts internationaux et bien plus.

## 🚀 Démarrage Rapide

### Installation
```bash
# Télécharger et configurer
git clone <votre-repo>
cd crypto-bank-app

# Configuration automatique
make install

# Démarrer tous les services
make start
# OU utiliser le docker-compose complet
docker-compose -f docker-compose-complete.yml up -d

# Vérifier que tout fonctionne
make test
```

### Accès à l'application
- **Application Web**: http://localhost:3000
- **API Documentation**: http://localhost:8080/docs
- **Monitoring**: http://localhost:3001 (Grafana)

## 💰 Acheter et Vendre des Cryptomonnaies

### 1. Acheter des Cryptos (Buy Orders)

#### Via l'interface web :
1. Connectez-vous à http://localhost:3000
2. Allez dans **"Exchange" → "Acheter Crypto"**
3. Sélectionnez la crypto à acheter (BTC, ETH, etc.)
4. Choisissez la devise de paiement (USD, EUR)
5. Entrez le montant
6. Choisissez le type d'ordre :
   - **Market** : Achat immédiat au prix actuel
   - **Limit** : Achat à un prix spécifique

#### Via API :
```bash
# Obtenir un devis
curl -X POST http://localhost:8080/api/v1/exchange/quote \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USD",
    "to_currency": "BTC",
    "from_amount": 1000.00
  }'

# Acheter des BTC avec USD (Ordre Market)
curl -X POST http://localhost:8080/api/v1/trading/buy \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "pay_currency": "USD",
    "amount": 0.02,
    "order_type": "market"
  }'

# Ordre d'achat limite (Limit Order)
curl -X POST http://localhost:8080/api/v1/trading/buy \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "pay_currency": "USD", 
    "amount": 0.02,
    "order_type": "limit",
    "limit_price": 45000.00
  }'
```

### 2. Vendre des Cryptos (Sell Orders)

#### Via l'interface web :
1. Allez dans **"Exchange" → "Vendre Crypto"**
2. Sélectionnez la crypto à vendre
3. Choisissez la devise à recevoir
4. Entrez le montant à vendre
5. Confirmez la vente

#### Via API :
```bash
# Vendre des BTC pour USD (Market)
curl -X POST http://localhost:8080/api/v1/trading/sell \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "receive_currency": "USD",
    "amount": 0.01,
    "order_type": "market"
  }'

# Ordre de vente limite
curl -X POST http://localhost:8080/api/v1/trading/sell \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "receive_currency": "USD",
    "amount": 0.01,
    "order_type": "limit",
    "limit_price": 50000.00
  }'
```

### 3. Voir vos ordres et historique

```bash
# Voir tous vos ordres
curl -X GET http://localhost:8080/api/v1/trading/orders \
  -H "Authorization: Bearer YOUR_TOKEN"

# Voir les ordres actifs
curl -X GET http://localhost:8080/api/v1/trading/orders/active \
  -H "Authorization: Bearer YOUR_TOKEN"

# Historique des échanges
curl -X GET http://localhost:8080/api/v1/exchange/history \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 💳 Cartes Prépayées Crypto/Fiat

### 1. Créer une Carte Virtuelle

#### Interface Web :
1. Allez dans **"Cartes" → "Nouvelle Carte"**
2. Choisissez **"Carte Virtuelle"**
3. Sélectionnez la devise (USD, EUR, BTC, ETH)
4. Définissez le montant initial
5. Configurez les limites

#### Via API :
```bash
# Créer une carte virtuelle USD
curl -X POST http://localhost:8080/api/v1/cards/virtual \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "USD",
    "cardholder_name": "John Doe",
    "initial_amount": 500.00,
    "source_wallet_id": "your_usd_wallet_id",
    "purpose": "online_shopping",
    "validity_months": 24
  }'

# Créer une carte crypto BTC
curl -X POST http://localhost:8080/api/v1/cards/virtual \
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

### 2. Commander une Carte Physique

```bash
curl -X POST http://localhost:8080/api/v1/cards/physical \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "USD",
    "cardholder_name": "John Doe",
    "initial_amount": 1000.00,
    "source_wallet_id": "your_wallet_id",
    "shipping_address": {
      "full_name": "John Doe",
      "address_line1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA",
      "phone_number": "+1234567890"
    },
    "express_shipping": true,
    "card_design": "classic"
  }'
```

### 3. Recharger une Carte

#### Depuis un portefeuille crypto ou fiat :
```bash
curl -X POST http://localhost:8080/api/v1/cards/CARD_ID/load \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 200.00,
    "source_wallet_id": "your_wallet_id",
    "description": "Rechargement carte"
  }'
```

#### Auto-rechargement automatique :
```bash
curl -X POST http://localhost:8080/api/v1/cards/CARD_ID/auto-load \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reload_amount": 100.00,
    "reload_threshold": 20.00,
    "source_wallet_id": "your_wallet_id"
  }'
```

### 4. Gestion des Cartes

```bash
# Voir toutes vos cartes
curl -X GET http://localhost:8080/api/v1/cards \
  -H "Authorization: Bearer YOUR_TOKEN"

# Détails d'une carte (avec numéro démasqué)
curl -X GET http://localhost:8080/api/v1/cards/CARD_ID/details \
  -H "Authorization: Bearer YOUR_TOKEN"

# Bloquer temporairement
curl -X POST http://localhost:8080/api/v1/cards/CARD_ID/freeze \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Voyage à l'étranger"}'

# Débloquer
curl -X POST http://localhost:8080/api/v1/cards/CARD_ID/unfreeze \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🎁 Cartes Cadeaux (Gift Cards)

### 1. Créer une Carte Cadeau

```bash
# Carte cadeau en USD
curl -X POST http://localhost:8080/api/v1/cards/gift \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "currency": "USD",
    "recipient_email": "friend@example.com",
    "message": "Joyeux anniversaire !",
    "design": "birthday",
    "validity_days": 365,
    "source_wallet_id": "your_wallet_id"
  }'

# Carte cadeau en crypto
curl -X POST http://localhost:8080/api/v1/cards/gift \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 0.001,
    "currency": "BTC",
    "recipient_phone": "+1234567890",
    "message": "Voici du Bitcoin pour toi !",
    "design": "crypto",
    "source_wallet_id": "your_btc_wallet_id"
  }'
```

### 2. Utiliser une Carte Cadeau

```bash
curl -X POST http://localhost:8080/api/v1/cards/gift/redeem \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ABCD-EFGH-IJKL-MNOP",
    "target_wallet_id": "your_wallet_id"
  }'
```

## 🔄 Conversion Crypto ↔ Fiat

### 1. Conversion Instantanée

```bash
# Convertir BTC → USD
curl -X POST http://localhost:8080/api/v1/exchange/execute \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_wallet_id": "your_btc_wallet",
    "to_wallet_id": "your_usd_wallet",
    "quote_id": "quote_from_previous_step"
  }'

# Convertir USD → ETH
curl -X POST http://localhost:8080/api/v1/exchange/execute \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_wallet_id": "your_usd_wallet", 
    "to_wallet_id": "your_eth_wallet",
    "quote_id": "quote_id"
  }'
```

### 2. Taux de Change en Temps Réel

```bash
# Voir tous les taux
curl -X GET http://localhost:8080/api/v1/exchange/rates

# Taux spécifique BTC/USD
curl -X GET http://localhost:8080/api/v1/exchange/rates/BTC/USD

# Marchés disponibles
curl -X GET http://localhost:8080/api/v1/exchange/markets
```

## 📊 Trading Avancé

### 1. Ordres Complexes

```bash
# Ordre Limite (Limit Order)
curl -X POST http://localhost:8080/api/v1/trading/limit-order \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USD",
    "to_currency": "BTC", 
    "amount": 1000.00,
    "limit_price": 42000.00,
    "order_type": "buy"
  }'

# Stop Loss (À venir)
curl -X POST http://localhost:8080/api/v1/trading/stop-loss \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "BTC",
    "to_currency": "USD",
    "amount": 0.1, 
    "stop_price": 40000.00
  }'
```

### 2. Portfolio et Performance

```bash
# Voir votre portfolio
curl -X GET http://localhost:8080/api/v1/trading/portfolio \
  -H "Authorization: Bearer YOUR_TOKEN"

# Métriques de performance
curl -X GET http://localhost:8080/api/v1/trading/performance \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🌍 Fonctionnalités Multi-Devises

### Devises Supportées

#### Cryptomonnaies :
- **Bitcoin (BTC)** - La référence
- **Ethereum (ETH)** - Smart contracts 
- **USDT/USDC** - Stablecoins
- **Binance Coin (BNB)** - Exchange token
- **Et plus** : ADA, XRP, DOT, LTC, LINK

#### Devises Fiduciaires :
- **USD** - Dollar américain
- **EUR** - Euro
- **GBP** - Livre sterling  
- **CAD** - Dollar canadien
- **AUD** - Dollar australien
- **JPY** - Yen japonais

### Conversion Multi-Devises
```bash
# EUR → BTC
curl -X POST http://localhost:8080/api/v1/exchange/quote \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "EUR",
    "to_currency": "BTC",
    "from_amount": 850.00
  }'

# BTC → JPY  
curl -X POST http://localhost:8080/api/v1/exchange/quote \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "BTC", 
    "to_currency": "JPY",
    "from_amount": 0.02
  }'
```

## 🔧 Configuration et Personnalisation

### Limites et Frais par Défaut

#### Frais de Trading :
- **Crypto → Fiat** : 0.5%
- **Fiat → Crypto** : 0.75% 
- **Crypto → Crypto** : 0.25%
- **Ordres Limit** : -0.1% (maker rebate)

#### Limites de Cartes :
- **Carte Virtuelle** : 5,000 USD/jour
- **Carte Physique** : 10,000 USD/jour  
- **ATM** : 500 USD/jour
- **International** : Selon KYC

### Variables d'Environnement Importantes

```bash
# Crypto APIs
INFURA_API_KEY=your_infura_key
ALCHEMY_API_KEY=your_alchemy_key
COINBASE_API_KEY=your_coinbase_key

# Cartes
MARQETA_API_KEY=your_marqeta_key
MARQETA_API_SECRET=your_marqeta_secret

# Mobile Money
MTN_API_KEY=your_mtn_key
AIRTEL_API_KEY=your_airtel_key
```

## 📱 Utilisation Mobile (API)

Toutes les fonctionnalités sont disponibles via API REST pour une future app mobile :

### Headers Requis
```bash
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json
X-API-Version: v1
```

### Endpoints Principaux
```
# Authentification
POST /api/v1/auth/login
POST /api/v1/auth/register

# Portefeuilles
GET /api/v1/wallets
POST /api/v1/wallets

# Trading/Exchange
GET /api/v1/exchange/rates
POST /api/v1/trading/buy
POST /api/v1/trading/sell

# Cartes
GET /api/v1/cards
POST /api/v1/cards
POST /api/v1/cards/:id/load

# Transferts
POST /api/v1/transfers
POST /api/v1/transfers/mobile/send
```

## 🛡️ Sécurité et Conformité

### Niveaux KYC
- **Niveau 1** : Email + Téléphone → 1,000 USD/jour
- **Niveau 2** : ID + Adresse → 10,000 USD/jour  
- **Niveau 3** : Vérification avancée → 100,000 USD/jour

### Sécurité
- **Chiffrement** : AES-256 pour les clés privées
- **2FA** : TOTP obligatoire pour montants élevés
- **Biométrie** : Support Face ID/Touch ID (mobile)
- **Surveillance** : Détection de fraude en temps réel

## 🚨 Résolution de Problèmes

### Problèmes Courants

#### 1. Ordre d'achat échoué
```bash
# Vérifier le solde du portefeuille
curl -X GET http://localhost:8080/api/v1/wallets/WALLET_ID/balance \
  -H "Authorization: Bearer YOUR_TOKEN"

# Vérifier les limites KYC
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 2. Carte bloquée
```bash
# Vérifier le statut
curl -X GET http://localhost:8080/api/v1/cards/CARD_ID \
  -H "Authorization: Bearer YOUR_TOKEN"

# Débloquer si nécessaire
curl -X POST http://localhost:8080/api/v1/cards/CARD_ID/unfreeze \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 3. Conversion échouée
```bash
# Vérifier si le devis est encore valide
curl -X GET http://localhost:8080/api/v1/exchange/quote/QUOTE_ID \
  -H "Authorization: Bearer YOUR_TOKEN"

# Obtenir un nouveau devis
curl -X POST http://localhost:8080/api/v1/exchange/quote \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"from_currency":"BTC","to_currency":"USD","from_amount":0.01}'
```

### Support et Logs
```bash
# Voir les logs en temps réel
make logs

# Logs spécifiques 
make logs-exchange  # Service d'échange
make logs-cards     # Service de cartes

# Vérifier la santé des services
make health
```

## 🎯 Cas d'Usage Pratiques

### 1. Voyageur International
1. Créer une carte EUR depuis BTC
2. Configurer auto-rechargement  
3. Bloquer temporairement pour sécurité
4. Utiliser à l'étranger sans frais

### 2. Trader Crypto
1. Surveiller les taux en temps réel
2. Placer des ordres limite
3. Convertir profits en stablecoin
4. Retirer via carte prépayée

### 3. E-commerce
1. Carte virtuelle pour achats en ligne
2. Limites personnalisées par vendeur
3. Rechargement automatique
4. Historique détaillé

### 4. Envoi d'Argent
1. Convertir crypto → fiat local
2. Envoyer via Mobile Money
3. Gift cards pour la famille
4. Suivi en temps réel

---

## 🏆 Zekora : La Banque du Futur

Avec Zekora, vous avez accès à tous les outils financiers modernes en un seul endroit :
- ✅ **Trading crypto** professionnel
- ✅ **Cartes prépayées** crypto et fiat
- ✅ **Conversion** instantanée multi-devises
- ✅ **Transferts** internationaux
- ✅ **Mobile Money** Afrique
- ✅ **Gift Cards** numériques
- ✅ **APIs** complètes pour développeurs

**Commencez dès maintenant** : http://localhost:3000

Pour plus d'aide : support@cryptobank.com