# Architecture des Wallets Fiat

## Vue d'Ensemble

Les wallets fiat (USD, EUR, XAF, NGN, etc.) gèrent les devises traditionnelles. Contrairement aux wallets crypto qui ont des adresses blockchain, les wallets fiat sont purement des **enregistrements en base de données** avec des intégrations vers des partenaires de paiement externes.

## Types de Transactions Fiat

| Type | Source/Destination | Blockchain TX | Méthode |
|------|-------------------|---------------|---------|
| **Dépôt Externe** | Carte/Bank → User | Non | Stripe, Flutterwave |
| **Dépôt Mobile Money** | M-Pesa/Orange → User | Non | Flutterwave |
| **Retrait Externe** | User → Bank/MoMo | Non | Stripe, Flutterwave |
| **Transfert Interne** | User A → User B | Non | DB Debit/Credit |
| **Conversion FX** | EUR → XAF | Non | Crypto Rails (stablecoin) |
| **International** | France → Côte d'Ivoire | Non | Thunes |

---

## 1. Modèle Wallet Fiat

```go
type Wallet struct {
    ID          string    // UUID
    UserID      string    // Propriétaire
    Currency    string    // USD, EUR, XAF, NGN...
    WalletType  string    // "fiat"
    Balance     float64   // Solde disponible
    LockedBal   float64   // Fonds bloqués (en attente)
    
    // Pas d'adresse blockchain pour fiat
    WalletAddress *string  // NULL pour fiat
    PrivateKey    *string  // NULL pour fiat
}
```

---

## 2. Dépôts Externes (Fiat → Wallet)

### 2.1 Flux de Dépôt par Carte Bancaire (Stripe)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                  FLUX: DÉPÔT CARTE BANCAIRE (STRIPE)                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] INITIATION SUR LE FRONTEND                                              ║
║  ════════════════════════════════                                            ║
║                                                                              ║
║      User → "Je veux déposer 100 USD par carte"                              ║
║                    ↓                                                         ║
║      Frontend affiche le formulaire Stripe Elements                          ║
║      (Secure Card Input - PCI DSS Compliant)                                 ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] CRÉATION DU PAYMENT INTENT (BACKEND)                                    ║
║  ════════════════════════════════════════                                    ║
║                                                                              ║
║      transfer-service reçoit: POST /api/v1/deposits/card                     ║
║      {                                                                       ║
║          "wallet_id": "wallet-abc",                                          ║
║          "amount": 100,                                                      ║
║          "currency": "USD"                                                   ║
║      }                                                                       ║
║                    ↓                                                         ║
║      stripeProvider.CreatePaymentIntent(amount, currency, metadata)          ║
║      → Retourne: { client_secret: "pi_xxx_secret_xxx" }                      ║
║      → Retourne: { client_secret: "pi_xxx_secret_xxx" }                      ║
║                                                                              ║
║      DB: transactions table                                                  ║
║      ┌─────────────────────────────────────────────────────────┐             ║
║      │ id: "tx-123"                                            │             ║
║      │ wallet_id: "wallet-abc"                                 │             ║
║      │ type: "deposit"                                         │             ║
║      │ amount: 100                                             │             ║
║      │ status: "pending"                                       │             ║
║      │ provider: "stripe"                                      │             ║
║      │ provider_ref: "pi_xxx"                                  │             ║
║      └─────────────────────────────────────────────────────────┘             ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] PAIEMENT CÔTÉ CLIENT                                                    ║
║  ═════════════════════════                                                   ║
║                                                                              ║
║      Frontend: stripe.confirmCardPayment(client_secret, { card })            ║
║                    ↓                                                         ║
║      Stripe traite le paiement (3DS si nécessaire)                           ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [4] WEBHOOK STRIPE → CRÉDIT WALLET                                          ║
║  ══════════════════════════════════                                          ║
║                                                                              ║
║      Stripe envoie: POST /webhooks/stripe                                    ║
║      { event: "payment_intent.succeeded", ... }                              ║
║                    ↓                                                         ║
║      stripeHandler.HandleWebhook(event)                                      ║
║                    ↓                                                         ║
║      DB Update:                                                              ║
║      ┌─────────────────────────────────────────────────────────┐             ║
║      │ UPDATE wallets SET                                      │             ║
║      │   balance = balance + 100                               │             ║
║      │ WHERE id = 'wallet-abc'                                 │             ║
║      │                                                         │             ║
║      │ UPDATE transactions SET                                 │             ║
║      │   status = 'completed'                                  │             ║
║      │ WHERE provider_ref = 'pi_xxx'                           │             ║
║      └─────────────────────────────────────────────────────────┘             ║
║                                                                              ║
║     ⭐ SOLDE USER CRÉDITÉ!                                                    ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### 2.2 Flux de Dépôt Mobile Money (Flutterwave)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║              FLUX: DÉPÔT MOBILE MONEY (FLUTTERWAVE)                          ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] INITIATION                                                              ║
║  ════════════════                                                            ║
║                                                                              ║
║      User → "Je veux déposer 5000 XAF via Orange Money"                      ║
║                    ↓                                                         ║
║      POST /api/v1/deposits/mobile-money                                      ║
║      {                                                                       ║
║          "wallet_id": "wallet-def",                                          ║
║          "amount": 5000,                                                     ║
║          "currency": "XAF",                                                  ║
║          "provider": "orange_money",                                         ║
║          "phone": "+237612345678"                                            ║
║      }                                                                       ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] INITIATION STK PUSH                                                     ║
║  ═══════════════════════                                                     ║
║                                                                              ║
║      flutterwaveProvider.InitiateCollection(request)                         ║
║      → API Flutterwave → STK Push vers le téléphone user                     ║
║                    ↓                                                         ║
║      User reçoit notification sur son téléphone:                             ║
║      "Confirmer paiement de 5000 XAF vers Zekora"                            ║
║                    ↓                                                         ║
║      User entre son code PIN Mobile Money                                    ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] WEBHOOK FLUTTERWAVE → CRÉDIT                                            ║
║  ═════════════════════════════════                                           ║
║                                                                              ║
║      Flutterwave envoie: POST /webhooks/flutterwave                          ║
║      { event: "charge.completed", status: "successful" }                     ║
║                    ↓                                                         ║
║      flutterwaveHandler.HandleWebhook(event)                                 ║
║                    ↓                                                         ║
║      UPDATE wallets SET balance = balance + 5000                             ║
║      WHERE id = 'wallet-def'                                                 ║
║                                                                              ║
║     ⭐ SOLDE USER CRÉDITÉ (souvent en < 30 secondes)                          ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 3. Transferts Internes (P2P)

### Flux: User A → User B (Même Plateforme)

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                  FLUX: TRANSFERT INTERNE (DB-ONLY)                           ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] INITIATION                                                              ║
║  ════════════════                                                            ║
║                                                                              ║
║      User A → "Envoyer 50 USD à user-b@email.com"                            ║
║                    ↓                                                         ║
║      POST /api/v1/transfers                                                  ║
║      {                                                                       ║
║          "from_wallet_id": "wallet-a",                                       ║
║          "to_email": "user-b@email.com",                                     ║
║          "amount": 50,                                                       ║
║          "currency": "USD"                                                   ║
║      }                                                                       ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] RÉSOLUTION DESTINATAIRE                                                 ║
║  ════════════════════════════                                                ║
║                                                                              ║
║      service.resolveOrCreateRecipientWallet(email, currency)                 ║
║                    ↓                                                         ║
║      walletRepo.GetUserIDByEmail("user-b@email.com") → user-b-id             ║
║                    ↓                                                         ║
║      walletRepo.GetWalletIDByUserAndCurrency(user-b-id, "USD")               ║
║      → wallet-b (ou création si n'existe pas)                                ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] CALCUL FRAIS                                                            ║
║  ════════════════                                                            ║
║                                                                              ║
║      feeService.CalculateFee("transfer_domestic", 50) → 0.25 USD (0.5%)      ║
║      Total débit: 50 + 0.25 = 50.25 USD                                      ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [4] EXÉCUTION (TRANSACTION ATOMIQUE)                                        ║
║  ═════════════════════════════════════                                       ║
║                                                                              ║
║      BEGIN TRANSACTION;                                                      ║
║                                                                              ║
║      -- Débit User A (montant + frais)                                       ║
║      UPDATE wallets SET balance = balance - 50.25                            ║
║      WHERE id = 'wallet-a' AND balance >= 50.25;                             ║
║                                                                              ║
║      -- Crédit User B (montant seulement)                                    ║
║      UPDATE wallets SET balance = balance + 50                               ║
║      WHERE id = 'wallet-b';                                                  ║
║                                                                              ║
║      -- Crédit Compte Frais Plateforme                                       ║
║      UPDATE platform_accounts SET balance = balance + 0.25                   ║
║      WHERE type = 'fees' AND currency = 'USD';                               ║
║                                                                              ║
║      -- Enregistrement Transfer                                              ║
║      INSERT INTO transfers (from_wallet_id, to_wallet_id, amount,            ║
║          fee, status) VALUES ('wallet-a', 'wallet-b', 50, 0.25, 'completed');║
║                                                                              ║
║      COMMIT;                                                                 ║
║                                                                              ║
║     ⭐ TRANSFERT INSTANTANÉ - AUCUN APPEL EXTERNE!                            ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 4. Retraits Externes (Wallet → Bank/MoMo)

### 4.1 Flux de Retrait Mobile Money

```
╔══════════════════════════════════════════════════════════════════════════════╗
║              FLUX: RETRAIT VERS MOBILE MONEY (PAYOUT)                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] INITIATION RETRAIT                                                      ║
║  ═══════════════════                                                         ║
║                                                                              ║
║      User → "Retirer 10000 XAF vers mon Orange Money"                        ║
║                    ↓                                                         ║
║      POST /api/v1/withdrawals/mobile-money                                   ║
║      {                                                                       ║
║          "wallet_id": "wallet-xyz",                                          ║
║          "amount": 10000,                                                    ║
║          "provider": "orange_money",                                         ║
║          "phone": "+237612345678"                                            ║
║      }                                                                       ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] VÉRIFICATION & DÉBIT IMMÉDIAT                                           ║
║  ══════════════════════════════════                                          ║
║                                                                              ║
║      1. Vérifier solde: wallet.balance >= 10000 + frais                      ║
║      2. Calculer frais: feeService.CalculateFee("withdrawal_momo") → 150 XAF ║
║                    ↓                                                         ║
║      UPDATE wallets SET                                                      ║
║          balance = balance - 10150,                                          ║
║          locked_balance = locked_balance + 10150                             ║
║      WHERE id = 'wallet-xyz'                                                 ║
║                                                                              ║
║      → Fonds BLOQUÉS jusqu'à confirmation du payout                          ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] EXÉCUTION PAYOUT                                                        ║
║  ════════════════════                                                        ║
║                                                                              ║
║      flutterwaveProvider.InitiatePayout({                                    ║
║          amount: 10000,                                                      ║
║          currency: "XAF",                                                    ║
║          recipient: {                                                        ║
║              type: "mobile_money_franco",                                    ║
║              account: "+237612345678",                                       ║
║              bank_code: "FRA"                                                ║
║          }                                                                   ║
║      })                                                                      ║
║      → Flutterwave envoie les fonds au numéro                                ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [4] WEBHOOK CONFIRMATION → DÉBLOCAGE                                        ║
║  ═════════════════════════════════════                                       ║
║                                                                              ║
║      Flutterwave webhook: { status: "successful" }                           ║
║                    ↓                                                         ║
║      UPDATE wallets SET                                                      ║
║          locked_balance = locked_balance - 10150                             ║
║      WHERE id = 'wallet-xyz'                                                 ║
║                                                                              ║
║      -- Les 150 XAF de frais vont au compte fees                             ║
║      UPDATE platform_accounts SET balance = balance + 150                    ║
║      WHERE type = 'fees' AND currency = 'XAF'                                ║
║                                                                              ║
║     ⭐ FONDS REÇUS SUR LE TÉLÉPHONE USER                                      ║
║                                                                              ║
║  [ÉCHEC] Si webhook: { status: "failed" }                                    ║
║          → Débloquer et rembourser:                                          ║
║          UPDATE wallets SET                                                  ║
║              balance = balance + 10150,                                      ║
║              locked_balance = locked_balance - 10150                         ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

---

## 5. Conversion de Devises (FX via Crypto Rails)

Quand un utilisateur envoie EUR vers quelqu'un qui veut recevoir XAF, le système utilise les "crypto rails" (stablecoins) pour la conversion.

```
┌─────────────────────────────────────────────────────────────────┐
│              CONVERSION FX VIA CRYPTO RAILS                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   EUR (User A) → USDC/USDT (Pool) → XAF (User B)                │
│                                                                 │
│   Étapes:                                                       │
│   1. Débiter EUR de User A                                      │
│   2. Convertir EUR → USDC (taux spot)                           │
│   3. Conversion interne USDC (pool plateforme)                  │
│   4. Convertir USDC → XAF (taux spot)                           │
│   5. Créditer XAF à User B                                      │
│                                                                 │
│   Avantages:                                                    │
│   - Taux compétitifs (crypto markets)                           │
│   - Settlement rapide (pas SWIFT)                               │
│   - Traçabilité blockchain si nécessaire                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Code InternalTransferService

```go
// internal_transfer.go
func (s *InternalTransferService) processInstantTransfer(ctx context.Context, transfer, req) {
    // Cross-currency? Use crypto rails
    if req.SourceCurrency != req.TargetCurrency {
        convResult, _ := s.cryptoRails.ConvertViaStablecoin(ctx, &ConversionRequest{
            SourceAmount:   req.Amount,
            SourceCurrency: req.SourceCurrency,  // EUR
            TargetCurrency: req.TargetCurrency,  // XAF
        })
        
        transfer.RecipientAmount = convResult.TargetAmount
        transfer.ExchangeRate = convResult.Rate
        transfer.ConversionFee = convResult.ConversionFee
    }
    
    // Debit sender, Credit recipient (DB)
    s.walletService.DebitWallet(ctx, req.SenderWalletID, transfer.SenderAmount, ref)
    s.walletService.CreditWallet(ctx, req.RecipientWalletID, transfer.RecipientAmount, ref)
}
```

---

## 6. Comptes Plateforme Fiat (Équivalent Hot/Cold)

> [!NOTE]
> Contrairement aux crypto qui ont des wallets **Hot** et **Cold** sur la blockchain, les comptes fiat sont purement **virtuels en base de données**. Ils représentent les fonds que la plateforme détient chez les agrégateurs (Stripe, Flutterwave) et les banques.

### 6.1 Types de Comptes Fiat Plateforme

| Type | Équivalent Crypto | Fonction |
|------|-------------------|----------|
| **Reserve** | Hot Wallet | Fonds liquides pour créditer les users lors des dépôts |
| **Operations** | Hot Wallet 2 | Buffer pour retraits/payouts vers agrégateurs |
| **Fees** | Fee Account | Accumulation des frais perçus |
| **Pending** | - | Fonds en attente de confirmation |

### 6.2 Modèle de Données

```go
// platform_account.go
type PlatformAccount struct {
    ID          string    // UUID
    Currency    string    // FCFA, EUR, USD
    AccountType string    // reserve, fees, operations, pending
    Name        string    // "Réserve FCFA Principal"
    Balance     float64   // Solde comptable
    MinBalance  float64   // Ne pas descendre en dessous
    MaxBalance  float64   // Limite de sécurité (0 = illimité)
    Priority    int       // 1-100, utilisé pour sélection intelligente
    IsActive    bool      // Compte actif/désactivé
}
```

### 6.3 Flux Détaillé: Dépôt Agrégateur → Crédit User

Quand un utilisateur fait un dépôt via Stripe ou Flutterwave :

```
╔══════════════════════════════════════════════════════════════════════════════╗
║       FLUX: WEBHOOK AGRÉGATEUR → CRÉDIT BALANCE USER                         ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] PAIEMENT CÔTÉ AGRÉGATEUR (Stripe/Flutterwave)                           ║
║  ══════════════════════════════════════════════════                          ║
║                                                                              ║
║      User paie 100 EUR via sa carte sur Stripe                               ║
║      Stripe: "Payment Intent Succeeded" → pi_xxx                             ║
║      L'ARGENT RÉEL est maintenant dans le compte Stripe de la plateforme     ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] WEBHOOK REÇU PAR NOTRE BACKEND                                          ║
║  ═════════════════════════════════                                           ║
║                                                                              ║
║      POST /webhooks/stripe                                                   ║
║      { event: "payment_intent.succeeded", amount: 10000, currency: "eur" }   ║
║                    ↓                                                         ║
║      stripeHandler.HandleWebhook() valide la signature                       ║
║                    ↓                                                         ║
║      Trouve le deposit pending: SELECT * FROM deposits WHERE stripe_ref = ?  ║
║      → deposit { wallet_id: "user-wallet-123", amount: 100, currency: "EUR" }║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] OPÉRATIONS COMPTABLES (Double Entry)                                    ║
║  ════════════════════════════════════════                                    ║
║                                                                              ║
║      Step A: Sélection intelligente du compte reserve                        ║
║      ─────────────────────────────────────────────                           ║
║      platformService.SelectBestAccountForCredit("EUR", "reserve", 100)       ║
║                    ↓                                                         ║
║      SQL:                                                                    ║
║      SELECT * FROM platform_accounts                                         ║
║      WHERE currency = 'EUR' AND account_type = 'reserve' AND is_active = true║
║        AND (max_balance = 0 OR balance + 100 <= max_balance)                 ║
║      ORDER BY priority DESC                                                  ║
║      LIMIT 1                                                                 ║
║      → Retourne: "Réserve EUR Principal" (balance: 950,000)                  ║
║                                                                              ║
║      Step B: Crédit Compte Réserve Plateforme                                ║
║      ────────────────────────────────────────                                ║
║      ⭐ Le compte reserve reflète les fonds que nous avons chez Stripe       ║
║      UPDATE platform_accounts SET balance = balance + 100                    ║
║      WHERE id = 'reserve-eur-1'                                              ║
║      → Nouveau solde: 950,100 EUR                                            ║
║                                                                              ║
║      Step C: Crédit Wallet Utilisateur                                       ║
║      ─────────────────────────────────                                       ║
║      UPDATE wallets SET balance = balance + 100                              ║
║      WHERE id = 'user-wallet-123'                                            ║
║      → User peut maintenant utiliser ses 100 EUR                             ║
║                                                                              ║
║      Step D: Enregistrement Transaction Ledger                               ║
║      ────────────────────────────────────────────                            ║
║      INSERT INTO platform_transactions (                                     ║
║          debit_account_id: 'external',           -- Source: Stripe           ║
║          debit_account_type: 'external',                                     ║
║          credit_account_id: 'reserve-eur-1',     -- Réserve plateforme       ║
║          credit_account_type: 'platform_fiat',                               ║
║          amount: 100,                                                        ║
║          currency: 'EUR',                                                    ║
║          operation_type: 'deposit',                                          ║
║          reference_type: 'stripe_intent',                                    ║
║          reference_id: 'pi_xxx'                                              ║
║      )                                                                       ║
║                                                                              ║
║     ⭐ À CE STADE:                                                            ║
║        - User a +100 EUR dans son wallet                                     ║
║        - Compte Reserve EUR a +100 EUR (reflète les fonds chez Stripe)       ║
║        - L'argent RÉEL est sur le compte Stripe/bancaire de la plateforme    ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### 6.4 Flux: Retrait User → Payout Agrégateur

```
╔══════════════════════════════════════════════════════════════════════════════╗
║       FLUX: RETRAIT USER → DÉBIT RÉSERVE → PAYOUT                            ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [1] INITIATION RETRAIT                                                      ║
║  ═════════════════════                                                       ║
║                                                                              ║
║      User: "Retirer 50 EUR vers mon compte bancaire"                         ║
║      POST /api/v1/withdrawals/bank { wallet_id, amount: 50, iban: "FR76..." }║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [2] DÉBIT WALLET USER (Immédiat)                                            ║
║  ═══════════════════════════════                                             ║
║                                                                              ║
║      UPDATE wallets SET balance = balance - 50                               ║
║      WHERE id = 'user-wallet-123' AND balance >= 50                          ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [3] SÉLECTION COMPTE OPERATIONS                                             ║
║  ═══════════════════════════════                                             ║
║                                                                              ║
║      platformService.SelectBestAccountForDebit("EUR", "operations", 50)      ║
║      → Retourne le compte operations avec solde suffisant                    ║
║                                                                              ║
║      ⭐ IMPORTANT: On débite le compte "operations" (pas "reserve")           ║
║         - Reserve = fonds des dépôts clients (couverture)                    ║
║         - Operations = fonds pour les payouts (buffer opérationnel)          ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [4] EXÉCUTION PAYOUT VIA AGRÉGATEUR                                         ║
║  ════════════════════════════════════                                        ║
║                                                                              ║
║      stripeProvider.CreatePayout({                                           ║
║          amount: 5000,  // centimes                                          ║
║          currency: "eur",                                                    ║
║          destination: "ba_xxx"  // External bank account                     ║
║      })                                                                      ║
║                                                                              ║
║      ⭐ Stripe va transférer 50€ de notre compte Stripe → Banque user        ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  [5] CONFIRMATION (Webhook ou Poll)                                          ║
║  ═══════════════════════════════════                                         ║
║                                                                              ║
║      Stripe: "payout.paid" → Le virement a été effectué                      ║
║                    ↓                                                         ║
║      UPDATE platform_accounts SET balance = balance - 50                     ║
║      WHERE id = 'operations-eur'                                             ║
║                                                                              ║
║      INSERT INTO platform_transactions (                                     ║
║          debit_account_id: 'operations-eur',                                 ║
║          credit_account_id: 'external',                                      ║
║          operation_type: 'withdrawal'                                        ║
║      )                                                                       ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### 6.5 Réconciliation Plateforme ↔ Agrégateur

> [!IMPORTANT]
> Les balances `platform_accounts` sont des **balances comptables internes**.
> Elles **DOIVENT** être réconciliées régulièrement avec les balances réelles chez Stripe/Flutterwave.

```
┌─────────────────────────────────────────────────────────────────┐
│                    RÉCONCILIATION                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Compte "Réserve EUR Principal" DB = 950,100 EUR               │
│                         ↕                                       │
│   Stripe Dashboard Balance           = 950,100 EUR   ✅ OK      │
│                                                                 │
│   Si décalage:                                                  │
│   - Vérifier transactions pendentes                             │
│   - Vérifier refunds non traités                                │
│   - Audit trail via platform_transactions                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.6 Sélection Intelligente des Comptes

Le système choisit automatiquement le meilleur compte pour chaque opération.

```go
// Pour CRÉDIT (dépôt reçu):
// 1. Compte actif avec capacité (balance + amount <= max_balance)
// 2. Plus haute priorité d'abord
SelectAccountForCredit(currency, "reserve", amount)

// Pour DÉBIT (retrait/payout):
// 1. Compte actif avec solde suffisant (balance - amount >= min_balance)
// 2. Plus haute priorité d'abord, puis plus gros solde
SelectAccountForDebit(currency, "operations", amount)
```

### 6.7 Configuration par Défaut (Seeded)

| Devise | Type | Nom | Solde Initial | Priorité |
|--------|------|-----|---------------|----------|
| FCFA | reserve | Réserve FCFA Principal | 1 milliard | 100 |
| FCFA | fees | Frais collectés FCFA | 0 | 100 |
| FCFA | operations | Opérations FCFA | 100 millions | 80 |
| EUR | reserve | Réserve EUR | 1 milliard | 100 |
| EUR | fees | Frais EUR | 0 | 100 |
| USD | reserve | Réserve USD | 1 milliard | 100 |
| USD | fees | Frais USD | 0 | 100 |

---

## 8. Annulation & Remboursement

### Règles d'Annulation

| Durée | Initiateur | Condition | Action |
|-------|------------|-----------|--------|
| < 5 min | Émetteur | Bénéficiaire n'a pas dépensé | Annulation complète |
| < 7 jours | Bénéficiaire | Veut retourner les fonds | Remboursement |
| > 7 jours | - | - | Aucune annulation possible |

```go
// transfer.go
func (s *TransferService) CancelTransfer(id, requesterID, reason string) {
    // 1. Vérifier que l'émetteur demande l'annulation
    // 2. Vérifier délai < 5 minutes
    // 3. Vérifier que le bénéficiaire n'a pas dépensé
    // 4. Si OK: Débiter bénéficiaire, Créditer émetteur
    // 5. Marquer transfer.status = "cancelled"
}
```

---

## 9. Administration & Gestion Agrégateurs

### 9.1 API Endpoints (Admin Only)

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| **GET** | `/admin/aggregators` | Liste tous les agrégateurs et leurs statuts |
| **PATCH** | `/admin/aggregators/:code` | Modifier configuration (fees, limits) |
| **POST** | `/admin/aggregators/:code/enable` | Activer globalement un agrégateur |
| **POST** | `/admin/aggregators/:code/disable` | Désactiver globalement |
| **POST** | `/admin/aggregators/:code/toggle-deposit` | Activer/Désactiver dépôts uniquement |
| **POST** | `/admin/aggregators/:code/toggle-withdraw` | Activer/Désactiver retraits uniquement |
| **POST** | `/admin/aggregators/:code/maintenance` | Définir mode maintenance (message custom) |

### 9.2 API Endpoints (Mode Test)

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| **POST** | `/admin/test-mode/credit` | Créditer manuellement un wallet (Hot Wallet → User) |
| **POST** | `/admin/test-mode/quick-credit` | Crédit rapide avec presets (Small/Medium/Large) |
| **POST** | `/admin/test-mode/simulate-webhook` | Simuler réception webhook (test integration) |
| **GET** | `/admin/test-mode/logs` | Audit logs des opérations de test |

### 9.3 Modèle de Données (AggregatorSetting)

```go
type AggregatorSetting struct {
    ProviderCode    string   // "stripe", "flutterwave"
    ProviderName    string   // "Stripe Payments"
    IsEnabled       bool     // Master switch
    DepositEnabled  bool     // Feature switch
    WithdrawEnabled bool     // Feature switch
    SupportedRegions []string // ["US", "EU", "CM", "CI"]
    Priority        int      // 1-100 (pour affichage frontend)
    
    // Frais & Limites
    MinAmount       float64
    MaxAmount       float64
    FeePercent      float64  // ex: 2.9%
    FeeFixed        float64  // ex: 0.30
    FeeCurrency     string   // "USD"
    
    // Maintenance
    MaintenanceMode bool
    MaintenanceMsg  string
}
```
