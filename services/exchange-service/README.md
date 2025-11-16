# Exchange Service

Service complet d'échange pour la banque crypto, gérant les échanges entre cryptomonnaies et monnaies fiduciaires.

## Fonctionnalités

### Échanges de Cryptomonnaies
- ✅ Échange instantané entre cryptomonnaies (BTC, ETH, LTC, ADA, DOT, XRP, etc.)
- ✅ Taux de change en temps réel avec mise à jour automatique
- ✅ Calcul automatique des frais basé sur le volume
- ✅ Historique complet des échanges

### Échanges de Monnaies Fiduciaires
- ✅ Échange entre devises principales (USD, EUR, GBP, JPY, CAD, AUD)
- ✅ Taux interbancaires compétitifs
- ✅ Frais transparents et réduits
- ✅ Exécution instantanée

### Trading Avancé
- ✅ Ordres au marché (Market Orders)
- ✅ Ordres à cours limité (Limit Orders)
- ✅ Ordres stop-loss
- ✅ Carnet d'ordres en temps réel
- ✅ Gestion de portfolio
- ✅ Analyses de performance

### Cotations et Taux
- ✅ API publique pour les taux de change
- ✅ Historique des taux
- ✅ Convertisseur de devises
- ✅ Comparaison des frais bancaires

## Architecture

```
exchange-service/
├── cmd/                    # Points d'entrée de l'application
├── internal/
│   ├── config/            # Configuration
│   ├── database/          # Connexions base de données
│   ├── handlers/          # Handlers HTTP
│   │   ├── exchange_handler.go
│   │   ├── fiat_handler.go
│   │   └── trading_handler.go
│   ├── middleware/        # Middlewares
│   ├── models/           # Modèles de données
│   ├── repository/       # Couche d'accès aux données
│   └── services/         # Logique métier
├── docker-compose.yml    # Configuration Docker
├── Dockerfile           # Image Docker
└── init.sql            # Schema de base de données
```

## API Endpoints

### Endpoints Publics

#### Taux de Change Crypto
```
GET /api/v1/rates                           # Tous les taux
GET /api/v1/rates/{from}/{to}               # Taux spécifique
GET /api/v1/rates/{from}/{to}/history       # Historique
GET /api/v1/convert                         # Convertisseur
GET /api/v1/supported-currencies            # Devises supportées
```

#### Taux de Change Fiat
```
GET /api/v1/fiat/rates                      # Tous les taux fiat
GET /api/v1/fiat/rates/{from}/{to}          # Taux fiat spécifique
GET /api/v1/fiat/rates/{from}/{to}/history  # Historique fiat
GET /api/v1/fiat/convert                    # Convertisseur fiat
GET /api/v1/fiat/currencies                 # Devises fiat supportées
GET /api/v1/fiat/fees/compare               # Comparaison frais
```

#### Trading Public
```
GET /api/v1/trading/tickers                 # Tickers du marché
GET /api/v1/trading/orderbook/{pair}        # Carnet d'ordres
```

### Endpoints Protégés (Authentification requise)

#### Opérations d'Échange
```
POST /api/v1/exchange/quote                 # Demander un devis
POST /api/v1/exchange/execute               # Exécuter un échange
GET  /api/v1/exchange/history               # Historique des échanges
GET  /api/v1/exchange/{id}                  # Détails d'un échange
```

#### Opérations Fiat
```
POST /api/v1/fiat/quote                     # Devis fiat
POST /api/v1/fiat/execute                   # Exécution fiat
```

#### Trading
```
POST /api/v1/trading/market-order           # Ordre au marché
POST /api/v1/trading/limit-order            # Ordre à cours limité
POST /api/v1/trading/stop-order             # Ordre stop-loss
GET  /api/v1/trading/orders                 # Ordres utilisateur
POST /api/v1/trading/orders/{id}/cancel     # Annuler un ordre
GET  /api/v1/trading/portfolio              # Portfolio
```

## Configuration

### Variables d'Environnement

```bash
# Base de données
DATABASE_URL=postgres://user:password@localhost/crypto_bank_exchange?sslmode=disable
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Sécurité
JWT_SECRET=your-secret-key

# Services
WALLET_SERVICE_URL=http://localhost:8084

# Frais (en pourcentage)
CRYPTO_TO_CRYPTO_FEE=0.5    # 0.5%
CRYPTO_TO_FIAT_FEE=0.75     # 0.75%
FIAT_TO_CRYPTO_FEE=0.75     # 0.75%
FIAT_TO_FIAT_FEE=0.25       # 0.25%

# Mise à jour des taux
RATE_UPDATE_INTERVAL=30     # 30 secondes
```

### Frais de Transaction

| Type d'Échange | Frais Base | Réductions Volume |
|----------------|------------|-------------------|
| Crypto → Crypto | 0.50% | Jusqu'à 50% |
| Crypto → Fiat | 0.75% | Jusqu'à 50% |
| Fiat → Crypto | 0.75% | Jusqu'à 50% |
| Fiat → Fiat | 0.25% | Jusqu'à 50% |

#### Réductions de Volume
- $100,000+ : -50% de frais
- $50,000+ : -40% de frais
- $10,000+ : -30% de frais
- $1,000+ : -15% de frais

## Déploiement

### Docker Compose

```bash
# Démarrer tous les services
docker-compose up -d

# Voir les logs
docker-compose logs -f exchange-service

# Arrêter les services
docker-compose down
```

### Build Manuel

```bash
# Installer les dépendances
go mod download

# Build l'application
go build -o bin/exchange-service main.go

# Lancer le service
./bin/exchange-service
```

## Tests

```bash
# Lancer tous les tests
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests d'intégration
go test -tags=integration ./...
```

## Monitoring

### Health Check
```
GET /health
```

Réponse :
```json
{
  "status": "healthy",
  "service": "exchange-service",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Métriques
- Taux de réussite des échanges
- Volume des transactions
- Latence des opérations
- Statut des connexions externes

## Sécurité

### Authentification
- JWT tokens requis pour toutes les opérations
- Validation des rôles utilisateur
- Rate limiting par IP

### Validation
- Validation stricte des montants
- Vérification des soldes wallet
- Prévention double dépense

### Audit
- Log de toutes les transactions
- Traçabilité complète
- Alertes automatiques

## Développement

### Structure du Code
```go
// Exemple de handler
func (h *ExchangeHandler) ExecuteExchange(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var req ExchangeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    exchange, err := h.exchangeService.ExecuteExchange(
        userID, req.FromWalletID, req.ToWalletID,
        req.FromCurrency, req.ToCurrency, req.Amount)
    
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, exchange)
}
```

### Ajout de Nouvelles Devises

1. Mettre à jour `isCryptoCurrency()` ou `isFiatCurrency()`
2. Ajouter les taux dans `simulateCryptoRate()` ou `simulateFiatRate()`
3. Mettre à jour la base de données avec les nouveaux taux
4. Tester les conversions

## Support

### Logs
Les logs sont disponibles dans :
- Container : `/app/logs/`
- Host : `./logs/`

### Dépannage

1. **Service ne démarre pas**
   - Vérifier les variables d'environnement
   - Contrôler la connectivité base de données
   - Vérifier les ports disponibles

2. **Taux non mis à jour**
   - Contrôler les logs du rate updater
   - Vérifier la connectivité Redis
   - Redémarrer le service

3. **Échanges échouent**
   - Vérifier les soldes wallet
   - Contrôler les logs de transaction
   - Valider les paramètres d'échange

## Roadmap

- [ ] P2P Trading
- [ ] Options et dérivés
- [ ] Staking intégré
- [ ] API GraphQL
- [ ] Websockets pour le trading en temps réel
- [ ] Interface mobile native