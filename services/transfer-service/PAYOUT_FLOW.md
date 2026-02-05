# 💸 Système de Virement Externe (Payout/Withdrawal Flow)

Ce document décrit le flux complet de virement externe (retrait) vers Mobile Money, compte bancaire, ou PayPal via les agrégateurs de paiement.

## 📋 Vue d'ensemble

Le système de payout permet aux utilisateurs de retirer des fonds de leur portefeuille vers des comptes externes (Mobile Money, banque, PayPal).

### Différence entre Deposit et Payout

| Opération | Direction | API Type | Hot Wallet |
|-----------|-----------|----------|------------|
| **Deposit** (Recharge) | Externe → User | Collection API | Reçoit les fonds |
| **Payout** (Retrait) | User → Externe | Disbursement/Transfer API | Envoie les fonds |

### Flux Principal

```
┌─────────┐     ┌──────────────┐     ┌─────────────────┐     ┌─────────────┐
│ Frontend│────▶│Transfer-Svc  │────▶│  Agrégateur     │────▶│   Webhook   │
│         │     │              │     │  (Flutterwave,  │     │   Callback  │
│ Initiate│     │ Debit User   │     │   Stripe...)    │     │             │
│ Payout  │     │ Credit Hot   │     │                 │     │             │
│         │     │ Call Provider│────▶│  Disbursement   │     │             │
│         │◀────│              │     │  API            │────▶│  Confirm TX │
│         │     │              │◀────────────────────────────│             │
│ Pending │────▶│              │     │                 │     │             │
│         │     │ Update Status│────▶│                 │     │             │
│ Success │◀────│              │     │                 │     │             │
└─────────┘     └──────────────┘     └─────────────────┘     └─────────────┘
```

## 🔄 États des Transactions

| État | Description |
|------|-------------|
| `pending` | Transaction initiée, en attente de traitement |
| `processing` | Envoi en cours chez le provider |
| `completed` | Virement réussi, fonds envoyés |
| `failed` | Échec - fonds remboursés à l'utilisateur |
| `cancelled` | Annulé par l'utilisateur - fonds remboursés |
| `expired` | Délai expiré (24h par défaut) |

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
     (timeout)    (provider OK)   (user cancel)
          │             │              │
          ▼             ▼              ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐
    │ expired  │  │processing│  │cancelled │
    └──────────┘  └────┬─────┘  └──────────┘
                       │              │
                  (webhook)      (refund)
                       │              │
          ┌────────────┴────────────┐
          │                         │
          ▼                         ▼
    ┌──────────┐              ┌──────────┐
    │completed │              │  failed  │
    └──────────┘              └──────────┘
                                   │
                              (auto refund)
```

## 🔌 Endpoints API

### Obtenir un devis (Quote)

```http
POST /api/v1/payouts/quote
Authorization: Bearer {token}
Content-Type: application/json

{
  "amount": 50000,
  "currency": "XOF",
  "provider": "flutterwave",
  "country": "CI",
  "payout_method": "mobile_money"
}
```

**Réponse:**
```json
{
  "amount": 50000,
  "currency": "XOF",
  "fee": 750,
  "fee_type": "percentage",
  "amount_received": 49250,
  "exchange_rate": 1.0,
  "received_currency": "XOF",
  "estimated_minutes": 5,
  "provider": "flutterwave",
  "payout_method": "mobile_money",
  "min_amount": 100,
  "max_amount": 5000000
}
```

### Initier un virement

```http
POST /api/v1/payouts/initiate
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": "uuid",
  "wallet_id": "uuid",
  "amount": 50000,
  "currency": "XOF",
  "provider": "flutterwave",
  "country": "CI",
  "payout_method": "mobile_money",
  "recipient_name": "Jean Dupont",
  "recipient_phone": "+2250701234567",
  "mobile_operator": "MTN",
  "mobile_number": "0701234567",
  "narration": "Retrait vers MTN",
  "pin": "1234"
}
```

**Réponse:**
```json
{
  "transaction_id": "pay_abc123_1738765432",
  "status": "processing",
  "provider": "flutterwave",
  "payout_method": "mobile_money",
  "amount": 50000,
  "currency": "XOF",
  "fee": 750,
  "amount_received": 49250,
  "recipient_name": "Jean Dupont",
  "recipient_account": "MTN ***4567",
  "message": "Payout initiated successfully",
  "aggregator_instance_id": "uuid",
  "hot_wallet_id": "uuid",
  "estimated_delivery": "2025-02-05T16:10:00Z",
  "new_balance": 150000
}
```

### Types de Payout Supportés

#### 1. Mobile Money

```json
{
  "payout_method": "mobile_money",
  "mobile_operator": "MTN",
  "mobile_number": "0701234567",
  "recipient_name": "Jean Dupont"
}
```

#### 2. Virement Bancaire

```json
{
  "payout_method": "bank_transfer",
  "bank_code": "012",
  "bank_name": "SGBCI",
  "account_number": "CI1234567890123456",
  "recipient_name": "Jean Dupont"
}
```

#### 3. PayPal

```json
{
  "payout_method": "paypal",
  "paypal_email": "jean.dupont@email.com",
  "recipient_name": "Jean Dupont"
}
```

#### 4. Virement International (IBAN/SWIFT)

```json
{
  "payout_method": "bank_transfer",
  "iban": "FR7612345678901234567890123",
  "swift_code": "BNPAFRPP",
  "bank_name": "BNP Paribas",
  "recipient_name": "Jean Dupont"
}
```

### Vérifier le statut

```http
GET /api/v1/payouts/{transaction_id}/status
Authorization: Bearer {token}
```

**Réponse:**
```json
{
  "transaction_id": "pay_abc123_1738765432",
  "status": "completed",
  "provider_reference": "FLW-MOCK-xxx",
  "amount": 50000,
  "currency": "XOF",
  "fee": 750,
  "amount_received": 49250,
  "recipient_name": "Jean Dupont",
  "recipient_account": "MTN ***4567",
  "payout_method": "mobile_money",
  "provider": "flutterwave",
  "status_message": "",
  "created_at": "2025-02-05T16:00:00Z",
  "completed_at": "2025-02-05T16:05:23Z"
}
```

### Annuler un virement

```http
POST /api/v1/payouts/{transaction_id}/cancel
Authorization: Bearer {token}
```

### Historique des virements

```http
GET /api/v1/payouts/user/{user_id}
Authorization: Bearer {token}
```

### Liste des banques

```http
GET /api/v1/payouts/banks?country=CI
```

### Liste des opérateurs Mobile Money

```http
GET /api/v1/payouts/mobile-operators?country=CI
```

## 🪝 Webhooks

### URL Format
```
POST /api/v1/payouts/webhook/{provider}
```

### Providers supportés

| Provider | Webhook URL | Événements |
|----------|-------------|------------|
| Flutterwave | `/webhook/flutterwave` | transfer.completed, transfer.failed |
| Stripe | `/webhook/stripe` | payout.paid, payout.failed |
| Paystack | `/webhook/paystack` | transfer.success, transfer.failed |
| PayPal | `/webhook/paypal` | PAYMENT.PAYOUTSBATCH.SUCCESS |
| MTN MoMo | `/webhook/mtn_momo` | disbursement callback |
| Orange Money | `/webhook/orange_money` | transfer callback |
| Wave | `/webhook/wave` | payout webhook |
| Thunes | `/webhook/thunes` | transfer status |

### Exemple Flutterwave Webhook (Payout)

```json
{
  "event": "transfer.completed",
  "data": {
    "id": 1234567,
    "reference": "pay_abc123_1738765432",
    "status": "SUCCESSFUL",
    "amount": 49250,
    "currency": "XOF",
    "complete_message": "Transfer successful"
  }
}
```

## 🏗️ Architecture Interne

### Flux de Mouvement de Fonds

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                     FLUX: USER WALLET → HOT WALLET → EXTERNE                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] VÉRIFICATION & DÉBIT USER                                               ║
║  ═════════════════════════════                                               ║
║                                                                              ║
║      User demande retrait de 50,000 XOF                                      ║
║                    ↓                                                         ║
║      Verify PIN (auth-service)                                               ║
║                    ↓                                                         ║
║      Check Balance: 200,000 XOF ✓                                            ║
║                    ↓                                                         ║
║      DB Update:                                                              ║
║      ┌─────────────────────────────────────────────────────────┐             ║
║      │ BEGIN TRANSACTION                                       │             ║
║      │                                                         │             ║
║      │ -- Débit User Wallet                                    │             ║
║      │ UPDATE wallets SET                                      │             ║
║      │   balance = balance - 50000                             │             ║
║      │ WHERE id = 'user-wallet-id'                             │             ║
║      │                                                         │             ║
║      │ -- Crédit Hot Wallet (staging)                          │             ║
║      │ UPDATE platform_accounts SET                            │             ║
║      │   balance = balance + 50000                             │             ║
║      │ WHERE id = 'hot-wallet-xof'                             │             ║
║      │                                                         │             ║
║      │ -- Log Transaction                                      │             ║
║      │ INSERT INTO transactions (...) VALUES (...)             │             ║
║      │                                                         │             ║
║      │ COMMIT                                                  │             ║
║      └─────────────────────────────────────────────────────────┘             ║
║                                                                              ║
║      ⭐ FONDS SÉCURISÉS DANS HOT WALLET                                       ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] APPEL API AGRÉGATEUR                                                    ║
║  ═════════════════════════                                                   ║
║                                                                              ║
║      Flutterwave POST /v3/transfers                                          ║
║      {                                                                       ║
║          "account_bank": "MTN",                                              ║
║          "account_number": "0701234567",                                     ║
║          "amount": 49250,  // Montant - Frais                                ║
║          "currency": "XOF",                                                  ║
║          "reference": "pay_abc123_1738765432",                               ║
║          "beneficiary_name": "Jean Dupont"                                   ║
║      }                                                                       ║
║                    ↓                                                         ║
║      Response: { "status": "success", "data": { "id": 123, ... } }           ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] WEBHOOK CONFIRMATION                                                    ║
║  ═════════════════════════                                                   ║
║                                                                              ║
║      Flutterwave POST /webhooks/flutterwave                                  ║
║      { "event": "transfer.completed", "data": { ... } }                      ║
║                    ↓                                                         ║
║      IF status == "SUCCESSFUL":                                              ║
║          Mark payout as completed                                            ║
║          Debit Hot Wallet (real payout happened)                             ║
║          ✅ VIREMENT RÉUSSI                                                   ║
║                    ↓                                                         ║
║      ELSE IF status == "FAILED":                                             ║
║          Mark payout as failed                                               ║
║          REVERSE: Credit User Wallet back                                    ║
║          REVERSE: Debit Hot Wallet back                                      ║
║          ❌ ÉCHEC - FONDS REMBOURSÉS                                          ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Double-Entry Accounting

Chaque payout génère des écritures comptables :

| Étape | Débit | Crédit | Description |
|-------|-------|--------|-------------|
| Initiation | User Wallet | Hot Wallet | Staging des fonds |
| Succès | Hot Wallet | (Externe) | Payout réel effectué |
| Échec | Hot Wallet | User Wallet | Remboursement auto |

## 💰 Frais par Méthode

| Méthode | Frais | Min | Max | Délai Estimé |
|---------|-------|-----|-----|--------------|
| Mobile Money | 1.5% | 50 XOF | - | 1-5 min |
| Virement Bancaire Local | 2% | 500 XOF | - | 1-24h |
| PayPal | 2.9% + 0.30 USD | - | - | 1h |
| Virement International | 3% | 10 USD | - | 1-5 jours |

## 🌍 Opérateurs Mobile Money par Pays

| Pays | Opérateurs |
|------|------------|
| 🇨🇮 Côte d'Ivoire | MTN, Orange, Moov, Wave |
| 🇸🇳 Sénégal | Orange, Free, Wave |
| 🇳🇬 Nigeria | MTN, Airtel, Glo |
| 🇬🇭 Ghana | MTN, Vodafone, AirtelTigo |
| 🇰🇪 Kenya | M-Pesa, Airtel |
| 🇨🇲 Cameroun | MTN, Orange |

## 🔐 Sécurité

### Vérifications Avant Payout

1. **Authentification JWT** - Token valide requis
2. **Vérification PIN** - PIN utilisateur vérifié via auth-service
3. **Vérification Solde** - Balance suffisante
4. **Limites** - Respect des limites quotidiennes/mensuelles
5. **KYC** - Niveau KYC suffisant pour le montant

### Vérification des Webhooks

```go
// Paystack: HMAC-SHA512
func verifyPaystackSignature(body []byte, signature, secret string) bool {
    mac := hmac.New(sha512.New, []byte(secret))
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expected))
}
```

### Variables d'environnement requises

```bash
# Webhook secrets
FLUTTERWAVE_WEBHOOK_SECRET=xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
PAYSTACK_WEBHOOK_SECRET=xxx
PAYPAL_WEBHOOK_SECRET=xxx

# Service URLs
WALLET_SERVICE_URL=http://wallet-service:8083
AUTH_SERVICE_URL=http://auth-service:8081
```

## 📊 Monitoring

### Métriques Prometheus

```
# Payouts initiés
payouts_initiated_total{provider="flutterwave",method="mobile_money"}

# Payouts complétés
payouts_completed_total{provider="flutterwave",status="success"}

# Volume de payouts
payouts_volume_total{provider="flutterwave",currency="XOF"}

# Durée de traitement
payout_processing_duration_seconds{provider="flutterwave"}
```

### Logs

```
[PayoutHandler] ✅ Payout initiated: pay_abc123_1738765432 | Provider: flutterwave | Amount: 50000.00 XOF | Recipient: MTN ***4567
[Flutterwave Payout Webhook] transfer.completed: pay_abc123_1738765432 -> SUCCESSFUL
[PayoutHandler] ✅ Payout completed: pay_abc123_1738765432 | Provider Ref: 1234567
```

## 🧪 Mode Démo

Le provider `demo` permet de tester sans vraie transaction:

```json
{
  "provider": "demo",
  "amount": 50000,
  "currency": "XOF",
  "payout_method": "mobile_money"
}
```

**Comportement:**
- Débit instantané du wallet utilisateur
- Pas d'appel externe réel
- Statut `completed` immédiat
- Parfait pour les tests

## 🚨 Gestion des Erreurs

| Code | Message | Action |
|------|---------|--------|
| 400 | Invalid PIN | Vérifier le code PIN |
| 400 | Insufficient balance | Vérifier le solde |
| 400 | Provider not available | Provider non disponible pour ce pays |
| 400 | Amount below minimum | Montant trop faible |
| 400 | Amount above maximum | Montant trop élevé |
| 400 | Withdrawals not enabled | Retraits désactivés pour ce provider |
| 500 | Provider error | Erreur chez l'agrégateur |

### Remboursement Automatique

En cas d'échec du payout après le débit utilisateur :

1. Le webhook signale l'échec
2. Le système reverse automatiquement les fonds
3. L'utilisateur récupère son solde
4. Une notification est envoyée

```
[PayoutHandler] ❌ Payout failed: pay_abc123_1738765432 | Reason: Invalid recipient (funds returned to user)
```

## 📝 Checklist Déploiement

- [ ] Configurer les clés API des agrégateurs dans Vault
- [ ] Configurer les webhook secrets
- [ ] Enregistrer les URLs de webhook chez chaque provider
- [ ] Créer les hot wallets pour chaque devise
- [ ] Configurer les mappings pays ↔ providers
- [ ] Configurer les limites de retrait
- [ ] Tester avec le provider `demo`
- [ ] Tester avec un provider réel en mode sandbox
- [ ] Activer le mode production

## 🔗 Liens Utiles

- [Flutterwave Transfers API](https://developer.flutterwave.com/reference/create-a-transfer)
- [Paystack Transfers API](https://paystack.com/docs/transfers)
- [MTN MoMo Disbursement API](https://momodeveloper.mtn.com/api-documentation/disbursement)
- [PayPal Payouts API](https://developer.paypal.com/docs/api/payments.payouts-batch/v1/)