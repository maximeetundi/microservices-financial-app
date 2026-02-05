# 🔧 Fix Agrégateurs & Mappings Pays - Guide Complet

## 📋 Problèmes Identifiés

### Problème 1: "Provider not available: no available instance for provider X"

**Symptôme**: Lors d'une tentative de recharge wallet, l'erreur suivante apparaît:
```
Provider not available: no available instance for provider lygos: sql: no rows in result set
```

**Cause**: La table `aggregator_instances` n'a pas d'entrée pour le provider demandé, OU la liaison entre l'instance et les hot wallets est manquante.

**Cause technique**: Incompatibilité de types entre:
- `platform_accounts.id` = `VARCHAR(36)`
- `aggregator_instance_wallets.hot_wallet_id` = `UUID` (devrait être VARCHAR)

### Problème 2: Les agrégateurs ne sont pas filtrés par pays

**Symptôme**: Sur le frontend utilisateur, tous les agrégateurs ne s'affichent pas pour le pays de l'utilisateur, même s'ils sont activés dans le panel admin.

**Cause**: La table `provider_countries` ne contient pas les mappings pays ↔ providers pour tous les agrégateurs. Seuls `demo` et `cinetpay` avaient des mappings dans le script `init.sql`.

---

## 🏗️ Architecture du Système

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         BASE PRINCIPALE (crypto_bank)                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   payment_providers        →  Liste des providers (Lygos, Stripe, etc.)     │
│         │                                                                    │
│         ▼                                                                    │
│   provider_countries       →  Mapping provider ↔ pays (CI, SN, CM...)       │
│         │                                                                    │
│         ▼                                                                    │
│   aggregator_settings      →  Config agrégateurs (enabled, fees, etc.)      │
│         │                                                                    │
│         ▼                                                                    │
│   aggregator_instances     →  Instances de chaque agrégateur                │
│         │                                                                    │
│         ▼                                                                    │
│   aggregator_instance_wallets  →  Liaison instance ↔ hot wallet             │
│         │                                                                    │
│         ▼                                                                    │
│   platform_accounts        →  Hot wallets (XOF, XAF, NGN, USD, EUR...)      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Flux de données

1. **Panel Admin** (`/dashboard/aggregators`):
   - Appelle `GET /api/v1/admin/payment-providers`
   - Retourne TOUS les providers avec leurs pays configurés
   - ✅ Fonctionne car il lit directement `payment_providers`

2. **Frontend User** (recharge wallet):
   - Appelle `GET /aggregators/deposit?country=CI`
   - Filtre par `provider_countries.country_code = 'CI'`
   - ❌ Ne trouve rien si pas de mapping dans `provider_countries`

3. **Initiation dépôt** (`POST /deposits/initiate`):
   - Cherche une instance dans `aggregator_instances_with_details`
   - ❌ Erreur "no rows" si pas d'instance ou pas de wallet lié

---

## 🚀 Solutions

### Option 1: Script Complet (Recommandé)

```bash
# Sur le VPS, exécuter le script complet
docker exec -i postgres psql -U admin -d crypto_bank < infrastructure/postgres/COMPLETE_FIX.sql
```

Ce script fait TOUT:
- Corrige le type de `hot_wallet_id`
- Crée tous les payment_providers
- Crée les aggregator_settings
- Crée les aggregator_instances
- Lie les hot wallets aux instances
- Ajoute TOUS les mappings pays ↔ providers
- Recrée la vue avec le bon ordre de colonnes

### Option 2: Reset Complet (Mode Dev)

```bash
cd /chemin/vers/microservices-financial-app

# Arrêter et supprimer les volumes
docker compose down -v

# Relancer (les scripts seront exécutés automatiquement)
docker compose up -d

# Vérifier les logs
docker logs postgres 2>&1 | grep -E "(NOTICE|ERROR|✅)"
```

### Option 3: Scripts Individuels

```bash
# 1. Fix instances et wallets
docker exec -i postgres psql -U admin -d crypto_bank < infrastructure/postgres/fix_aggregator_instances.sql

# 2. Ajouter les mappings pays
docker exec -i postgres psql -U admin -d crypto_bank < infrastructure/postgres/seed_provider_countries.sql

# 3. Vérifier
docker exec -i postgres psql -U admin -d crypto_bank < infrastructure/postgres/verify_and_fix_instances.sql
```

---

## ✅ Vérification

### 1. Vérifier les instances créées

```bash
docker exec -it postgres psql -U admin -d crypto_bank -c "
SELECT 
    agg.provider_code, 
    COUNT(DISTINCT ai.id) AS instances,
    COUNT(DISTINCT aiw.id) AS wallet_links,
    COUNT(DISTINCT pc.id) AS country_mappings
FROM aggregator_settings agg
LEFT JOIN aggregator_instances ai ON agg.id = ai.aggregator_id
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id
LEFT JOIN payment_providers pp ON agg.payment_provider_id = pp.id
LEFT JOIN provider_countries pc ON pp.id = pc.provider_id
GROUP BY agg.provider_code
ORDER BY agg.provider_code;
"
```

**Résultat attendu**:
```
 provider_code | instances | wallet_links | country_mappings
---------------+-----------+--------------+------------------
 cinetpay      |         1 |            5 |                9
 demo          |         1 |            5 |               11
 fedapay       |         1 |            5 |                5
 flutterwave   |         1 |            5 |                9
 lygos         |         1 |            5 |               11
 moov_money    |         1 |            5 |                8
 mtn_momo      |         1 |            5 |                8
 orange_money  |         1 |            5 |                6
 paypal        |         1 |            5 |                6
 paystack      |         1 |            5 |                4
 stripe        |         1 |            5 |                6
 wave          |         1 |            5 |                2
 yellowcard    |         1 |            5 |               10
```

### 2. Vérifier la disponibilité des instances

```bash
docker exec -it postgres psql -U admin -d crypto_bank -c "
SELECT provider_code, instance_name, hot_wallet_currency, availability_status 
FROM aggregator_instances_with_details
WHERE availability_status = 'available'
ORDER BY provider_code;
"
```

### 3. Vérifier les mappings pays pour un provider

```bash
docker exec -it postgres psql -U admin -d crypto_bank -c "
SELECT pp.name, pc.country_code, pc.currency, pc.is_active
FROM payment_providers pp
JOIN provider_countries pc ON pp.id = pc.provider_id
WHERE pp.name = 'lygos'
ORDER BY pc.country_code;
"
```

### 4. Test depuis le frontend

1. Connectez-vous à l'application
2. Allez dans "Recharger Compte"
3. Vérifiez que les agrégateurs correspondant à votre pays s'affichent
4. Testez un dépôt avec le provider "demo"

---

## 📁 Fichiers Créés/Modifiés

| Fichier | Description |
|---------|-------------|
| `init.sql` | Type `hot_wallet_id` corrigé, vue mise à jour |
| `fix_aggregator_instances.sql` | Script de correction des instances |
| `seed_provider_countries.sql` | Script de seeding des mappings pays |
| `verify_and_fix_instances.sql` | Script de diagnostic rapide |
| `COMPLETE_FIX.sql` | Script tout-en-un (recommandé) |
| `docker-compose.yml` | Scripts ajoutés à l'initialisation |

---

## 🔍 Debugging

### Voir les logs du transfer-service
```bash
docker logs transfer-service 2>&1 | grep -E "(ERROR|instance|provider)"
```

### Voir les logs du admin-service
```bash
docker logs admin-service 2>&1 | grep -E "(ERROR|provider|country)"
```

### Tester l'API directement
```bash
# Liste des agrégateurs pour CI (Côte d'Ivoire)
curl "https://api.app.tech-afm.com/transfer-service/api/v1/aggregators/deposit?country=CI"

# Liste des méthodes de paiement
curl "https://api.admin.tech-afm.com/api/v1/admin/payment-methods?country=CI"
```

---

## 🌍 Pays Supportés par Provider

| Provider | Pays |
|----------|------|
| **demo** | CI, SN, CM, BJ, TG, BF, ML, NE, NG, GH, KE |
| **lygos** | CI, SN, BF, ML, TG, BJ, NE, CM, CD, GN, LR |
| **cinetpay** | CI, SN, CM, BF, ML, TG, BJ, NE, GN |
| **flutterwave** | NG, GH, KE, ZA, UG, TZ, RW, CI, CM |
| **paystack** | NG, GH, ZA, KE |
| **orange_money** | CI, SN, ML, BF, CM, GN |
| **mtn_momo** | CI, CM, SN, BJ, GH, UG, RW, BF |
| **wave** | SN, CI |
| **moov_money** | CI, BJ, TG, BF, NE, CM, GA, CG |
| **fedapay** | BJ, TG, NE, CI, CM |
| **yellowcard** | NG, GH, CI, SN, KE, ZA, CM, UG, TZ, RW |
| **stripe** | Global (CI, SN, CM, NG, GH, KE) |
| **paypal** | Global (CI, SN, CM, NG, GH, KE) |

---

## ⚠️ Notes Importantes

1. **Après avoir appliqué le fix**, redémarrez les services:
   ```bash
   docker restart transfer-service admin-service
   ```

2. **Le cache de l'admin-client** dans transfer-service expire après 5 minutes. Si les changements ne sont pas visibles immédiatement, attendez ou redémarrez.

3. **En production**, testez d'abord sur un environnement de staging.

4. **Les hot wallets** doivent avoir un solde suffisant pour que l'instance soit considérée comme "available".