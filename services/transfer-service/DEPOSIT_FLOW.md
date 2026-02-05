# 💰 Système de Recharge (Deposit Flow)

Ce document décrit le flux complet de recharge de wallet via les agrégateurs de paiement.

## 📋 Vue d'ensemble

Le système de recharge permet aux utilisateurs de créditer leur portefeuille en utilisant différents moyens de paiement (Mobile Money, Carte Bancaire, PayPal, etc.) via des agrégateurs tiers.

### Flux Principal

```
┌─────────┐     ┌──────────────┐     ┌─────────────────┐     ┌─────────────┐
│ Frontend│────▶│Transfer-Svc  │────▶│  Agrégateur     │────▶│   Webhook   │
│         │     │              │     │  (Flutterwave,  │     │   Callback  │
│ Initiate│     │ Create TX    │     │   Stripe...)    │     │             │
│ Deposit │     │ Get Payment  │     │                 │     │             │
│         │◀────│ URL          │◀────│  Payment Page   │     │             │
│         │     │              │     │                 │────▶│             │
│ Redirect│────▶│              │     │  User Pays      │     │  Confirm TX │
│         │     │              │◀────────────────────────────│             │
│         │     │ Process      │     │                 │     │             │
│         │     │ Credit User  │────▶│                 │     │             │
│ Success │◀────│ Wallet       │     │                 │     │             │
└─────────┘     └──────────────┘     └─────────────────┘     └─────────────┘
```

## 🔄 États des Transactions

| État | Description |
|------|-------------|
| `pending` | Transaction initiée, en attente de paiement |
| `processing` | Paiement en cours de traitement |
| `completed` | Paiement réussi, wallet crédité |
| `failed` | Paiement échoué |
| `cancelled` | Annulé par l'utilisateur |
| `expired` | Délai expiré (2h par défaut) |

### Diagramme d'états

```
                    ┌──────────┐
                    │  START   │
                    └────┬─────┘
                         │
                         ▼
                    ┌──────────┐
          ┌────────│ pending  │────────┐
          │        └────┬─────┘        │
          │             │              │
     (timeout)    (webhook OK)    (user cancel)
          │             │              │
          ▼             ▼              ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐
    │ expired  │  │completed │  │cancelled │
    └──────────┘  └──────────┘  └──────────┘
          │             │
     (webhook fail)     │
          │             │
          ▼             │
    ┌──────────┐        │
    │  failed  │◀───────┘
    └──────────┘
```

## 🔌 Endpoints API

### Initier un dépôt

```http
POST /api/v1/deposits/initiate
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": "uuid",
  "wallet_id": "uuid",  // Optionnel
  "amount": 5000,
  "currency": "XOF",
  "provider": "flutterwave",  // lygos, stripe, paystack, cinetpay, etc.
  "country": "CI",
  "email": "user@example.com",
  "phone": "+2250701234567",
  "return_url": "https://app.tech-afm.com/wallet?deposit_callback=true",
  "cancel_url": "https://app.tech-afm.com/wallet?deposit_cancelled=true"
}
```

**Réponse:**
```json
{
  "transaction_id": "dep_abc123_1738765432",
  "status": "pending",
  "payment_url": "https://checkout.flutterwave.com/pay/xyz",
  "provider": "flutterwave",
  "amount": 5000,
  "currency": "XOF",
  "fee": 75,
  "expires_at": "2025-02-05T18:00:00Z",
  "sdk_config": {
    "public_key": "FLWPUBK_TEST-xxx",
    "environment": "test",
    "currency": "XOF",
    "country": "CI"
  }
}
```

### Vérifier le statut

```http
GET /api/v1/deposits/{transaction_id}/status
Authorization: Bearer {token}
```

**Réponse:**
```json
{
  "transaction_id": "dep_abc123_1738765432",
  "status": "completed",
  "amount": 5000,
  "currency": "XOF",
  "fee": 75,
  "wallet_credited": true,
  "completed_at": "2025-02-05T16:05:23Z"
}
```

### Annuler un dépôt

```http
POST /api/v1/deposits/{transaction_id}/cancel
Authorization: Bearer {token}
```

### Historique des dépôts

```http
GET /api/v1/deposits/user/{user_id}
Authorization: Bearer {token}
```

## 🪝 Webhooks

### URL Format
```
POST /api/v1/deposits/webhook/{provider}
```

### Providers supportés

| Provider | Webhook URL | Signature Header |
|----------|-------------|------------------|
| Flutterwave | `/webhook/flutterwave` | `verif-hash` |
| Stripe | `/webhook/stripe` | `Stripe-Signature` |
| Paystack | `/webhook/paystack` | `X-Paystack-Signature` |
| CinetPay | `/webhook/cinetpay` | N/A (IP whitelist) |
| Lygos | `/webhook/lygos` | Custom |
| Orange Money | `/webhook/orange_money` | Token |
| MTN MoMo | `/webhook/mtn_momo` | Subscription Key |
| Wave | `/webhook/wave` | Custom |
| PayPal | `/webhook/paypal` | PAYPAL-TRANSMISSION-SIG |
| FedaPay | `/webhook/fedapay` | X-FedaPay-Signature |
| Moov Money | `/webhook/moov_money` | N/A |
| YellowCard | `/webhook/yellowcard` | X-YellowCard-Signature |

### Exemple Flutterwave Webhook

```json
{
  "event": "charge.completed",
  "data": {
    "id": 1234567,
    "tx_ref": "dep_abc123_1738765432",
    "flw_ref": "FLW-MOCK-xxx",
    "amount": 5000,
    "currency": "XOF",
    "charged_amount": 5075,
    "app_fee": 75,
    "status": "successful",
    "customer": {
      "email": "user@example.com"
    },
    "meta": {
      "wallet_id": "uuid",
      "user_id": "uuid"
    }
  }
}
```

## 🎨 Intégration Frontend

### Composant Vue.js

```vue
<template>
  <DepositModal 
    :is-open="showDepositModal"
    :wallet-id="selectedWallet.id"
    :currency="selectedWallet.currency"
    @close="showDepositModal = false"
    @success="onDepositSuccess"
  />
</template>
```

### SDK JavaScript

Le composant charge automatiquement les SDK nécessaires:

- **Flutterwave**: `https://checkout.flutterwave.com/v3.js`
- **Paystack**: `https://js.paystack.co/v1/inline.js`
- **Stripe**: Redirect vers Checkout

### Exemple d'utilisation SDK Flutterwave

```javascript
FlutterwaveCheckout({
  public_key: sdkConfig.public_key,
  tx_ref: transactionId,
  amount: amount,
  currency: 'XOF',
  payment_options: 'card,mobilemoney,ussd',
  customer: {
    email: user.email,
    phone_number: phone,
    name: `${user.first_name} ${user.last_name}`
  },
  customizations: {
    title: 'Zekora - Recharge',
    description: 'Recharge de 5000 XOF',
    logo: '/logo.png'
  },
  callback: (response) => {
    if (response.status === 'successful') {
      // Success!
    }
  },
  onclose: () => {
    // User closed modal
  }
});
```

## ⏰ Expiration Automatique

Un service en arrière-plan vérifie les transactions en attente et les marque comme expirées après le délai configuré (2h par défaut).

```go
// Configuration
depositExpiryService := service.NewDepositExpiryService(
    depositRepo, 
    5*time.Minute,  // Intervalle de vérification
)
depositExpiryService.Start()
```

### Requête SQL d'expiration

```sql
UPDATE deposit_transactions
SET status = 'expired',
    status_message = 'Transaction expired due to timeout',
    failed_at = CURRENT_TIMESTAMP
WHERE status = 'pending'
  AND expires_at IS NOT NULL
  AND expires_at < CURRENT_TIMESTAMP;
```

## 🏗️ Architecture des Tables

### deposit_transactions

```sql
CREATE TABLE deposit_transactions (
    id VARCHAR(100) PRIMARY KEY,
    user_id UUID NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    fee DECIMAL(20,8) DEFAULT 0,
    provider_code VARCHAR(50) NOT NULL,
    aggregator_instance_id UUID,
    hot_wallet_id VARCHAR(36),
    payment_url TEXT,
    provider_reference VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending',
    webhook_data JSONB,
    expires_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔐 Sécurité

### Vérification des Webhooks

Chaque provider utilise un mécanisme de signature différent:

```go
// Paystack: HMAC-SHA512
func verifyHMACSHA512(body []byte, signature, secret string) bool {
    mac := hmac.New(sha512.New, []byte(secret))
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expected))
}

// Stripe: HMAC-SHA256 avec timestamp
func verifyStripeSignature(body []byte, signature, secret string) bool {
    // Parse: t=timestamp,v1=signature
    // Compute: HMAC-SHA256(timestamp + "." + body, secret)
}
```

### Variables d'environnement requises

```bash
# Webhook secrets
FLUTTERWAVE_WEBHOOK_SECRET=xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
PAYSTACK_WEBHOOK_SECRET=xxx
CINETPAY_WEBHOOK_SECRET=xxx
PAYPAL_WEBHOOK_SECRET=xxx

# Service URLs
WALLET_SERVICE_URL=http://wallet-service:8083
```

## 📊 Monitoring

### Métriques Prometheus

```
# Dépôts initiés
deposit_initiated_total{provider="flutterwave",currency="XOF"}

# Dépôts complétés
deposit_completed_total{provider="flutterwave",status="success"}

# Durée de traitement
deposit_processing_duration_seconds{provider="flutterwave"}
```

### Logs

```
[DepositHandler] ✅ Deposit initiated: dep_abc123_1738765432 | Provider: flutterwave | Amount: 5000.00 XOF
[Flutterwave Webhook] Processed: dep_abc123_1738765432 -> successful
[DepositHandler] ✅ Deposit completed: dep_abc123_1738765432 | Amount: 5000.00 XOF | User: user-uuid
```

## 🧪 Mode Démo

Le provider `demo` permet de tester sans vraie transaction:

```json
{
  "provider": "demo",
  "amount": 5000,
  "currency": "XOF"
}
```

**Comportement:**
- Crédit instantané du wallet
- Pas de redirection
- Statut `completed` immédiat

## 🚨 Gestion des Erreurs

| Code | Message | Action |
|------|---------|--------|
| 400 | Provider not available | Vérifier les instances d'agrégateurs |
| 400 | No available wallet | Créer un hot wallet pour la devise |
| 500 | Failed to initiate payment | Vérifier les clés API du provider |
| 408 | Transaction expired | Réinitier une nouvelle transaction |

## 📝 Checklist Déploiement

- [ ] Configurer les webhook secrets dans les variables d'environnement
- [ ] Enregistrer les URLs de webhook chez chaque provider
- [ ] Vérifier que les hot wallets existent pour chaque devise
- [ ] Créer les instances d'agrégateurs dans la DB
- [ ] Configurer les mappings pays ↔ providers
- [ ] Tester avec le provider `demo`
- [ ] Tester avec un provider réel en mode sandbox
- [ ] Activer le mode production