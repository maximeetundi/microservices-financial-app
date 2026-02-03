# 🌍 World Currencies Integration Guide

## Vue d'ensemble

Ce fichier SQL ajoute le support de **150+ devises mondiales** (ISO 4217) à votre application financière.

## Devises Supportées

### 💵 Devises Fiat Majeures
- **Amérique**: USD, CAD, MXN, BRL, ARS, CLP, COP, PEN
- **Europe**: EUR, GBP, CHF, NOK, SEK, DKK, PLN, CZK, HUF, RON
- **Asie**: JPY, CNY, INR, KRW, SGD, HKD, THB, MYR, PHP, IDR, VND
- **Moyen-Orient**: AED, SAR, QAR, KWD, BHD, ILS, TRY
- **Océanie**: AUD, NZD

### 🌍 Devises Africaines (Complètes)
- **Afrique de l'Ouest**: XOF (FCFA), GHS, NGN, GMD, GNF, LRD, SLL
- **Afrique Centrale**: XAF (FCFA), CDF
- **Afrique de l'Est**: KES, TZS, UGX, RWF, ETB, SOS
- **Afrique Australe**: ZAR, ZMW, BWP, NAD, LSL, SZL
- **Afrique du Nord**: MAD, EGP, TND, DZD, LYD

### 💎 Cryptomonnaies
BTC, ETH, USDT, USDC, BNB, XRP, ADA, SOL, DOT, DOGE, MATIC, LTC, TRX, AVAX, UNI

## Installation

### 1. Exécuter le script SQL

```bash
# Via Docker
docker exec -i postgres psql -U postgres -d crypto_bank < infrastructure/postgres/world_currencies.sql

# Ou via psql direct
psql -U postgres -d crypto_bank -f infrastructure/postgres/world_currencies.sql
```

### 2. Vérifier l'installation

```sql
-- Vérifier les devises supportées par Demo
SELECT name, supported_currencies 
FROM payment_providers 
WHERE name = 'demo';

-- Compter les paires de taux de change
SELECT COUNT(*) as total_pairs 
FROM exchange_rates;
```

## Providers Mis à Jour

| Provider | Devises Supportées |
|----------|-------------------|
| **Demo** | 150+ devises (toutes) |
| **Flutterwave** | 15 devises africaines + majeures |
| **CinetPay** | XOF, XAF, GNF, USD, EUR |
| **Paystack** | NGN, GHS, ZAR, KES + majeures |
| **Stripe** | 40+ devises internationales |
| **Orange Money** | XOF, XAF, EUR |
| **MTN MoMo** | XOF, XAF, GHS, UGX, RWF, ZMW, EUR |
| **Wave** | XOF, EUR, USD |

## Taux de Change Ajoutés

Le script ajoute **80+ paires de taux de change** couvrant:
- Toutes les conversions USD vers devises majeures
- Conversions entre devises africaines
- Conversions crypto-fiat (BTC/ETH vers USD/EUR/GBP)
- Spread et volume 24h pour chaque paire

## Utilisation

### Créer un wallet dans n'importe quelle devise

```go
// Backend - wallet-service
wallet := &Wallet{
    UserID:   userID,
    Currency: "MAD", // Dirham marocain
    WalletType: "fiat",
    Balance: 0,
}
```

### Frontend - Sélection de devise

```typescript
// Toutes les devises sont maintenant disponibles
const currencies = [
  { code: 'USD', name: 'Dollar américain', symbol: '$' },
  { code: 'EUR', name: 'Euro', symbol: '€' },
  { code: 'XOF', name: 'Franc CFA (BCEAO)', symbol: 'CFA' },
  { code: 'MAD', name: 'Dirham marocain', symbol: 'MAD' },
  { code: 'NGN', name: 'Naira nigérian', symbol: '₦' },
  // ... 145+ autres devises
]
```

## Notes Importantes

### Taux de Change
- Les taux sont **indicatifs** et doivent être mis à jour régulièrement
- Utilisez une API externe (Fixer.io, CurrencyLayer, etc.) pour des taux réels
- Le spread est configuré pour chaque paire

### Mode Demo
- Le provider **Demo** supporte TOUTES les devises
- Parfait pour les tests sans intégration réelle
- Crédite directement sans appel API

### Production
- Activez les providers réels selon vos besoins
- Configurez les clés API dans les variables d'environnement
- Testez chaque provider avec sa devise native

## Maintenance

### Ajouter une nouvelle devise

```sql
-- 1. Ajouter au provider Demo
UPDATE payment_providers 
SET supported_currencies = jsonb_insert(
    supported_currencies, 
    '{-1}', 
    '"NEW"'
)
WHERE name = 'demo';

-- 2. Ajouter les taux de change
INSERT INTO exchange_rates (from_currency, to_currency, rate, ...) 
VALUES ('USD', 'NEW', 1.0, ...);
```

### Mettre à jour les taux

```sql
UPDATE exchange_rates 
SET rate = 1.23, 
    last_updated = NOW() 
WHERE from_currency = 'USD' 
  AND to_currency = 'EUR';
```

## Ressources

- [ISO 4217 Currency Codes](https://www.iso.org/iso-4217-currency-codes.html)
- [Exchange Rate APIs](https://exchangeratesapi.io/)
- [Crypto Price APIs](https://www.coingecko.com/en/api)
