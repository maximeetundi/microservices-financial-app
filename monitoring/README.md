# Monitoring Stack - Guide d'Installation

## 📦 Prérequis

1. Docker et Docker Compose installés
2. Services principaux (`docker-compose.yml`) démarrés d'abord

## 🚀 Démarrage Rapide

```bash
# 1. Démarrer les services principaux
docker compose up -d

# 2. Attendre que les services soient prêts
sleep 30

# 3. Démarrer le monitoring
docker compose -f docker-compose-monitoring.yml up -d
```

## 🔗 Accès aux Interfaces

| Service | URL | Credentials |
|---------|-----|-------------|
| **Kibana** (Logs) | http://localhost:5601 | - |
| **Grafana** (Metrics) | http://localhost:3001 | admin / Zekora2024! |
| **Prometheus** | http://localhost:9090 | - |
| **Jaeger** (Tracing) | http://localhost:16686 | - |
| **Alertmanager** | http://localhost:9093 | - |
| **Uptime Kuma** | http://localhost:3003 | Créer compte au 1er lancement |
| **Elasticsearch** | http://localhost:9200 | - |

## 📊 Architecture

```
                                ┌─────────────────┐
                                │   Grafana       │ ← Dashboards
                                │   :3001         │
                                └────────┬────────┘
                                         │
              ┌──────────────────────────┼──────────────────────────┐
              │                          │                          │
    ┌─────────▼────────┐       ┌────────▼────────┐       ┌─────────▼─────────┐
    │   Prometheus     │       │  Elasticsearch  │       │     Jaeger        │
    │   :9090          │       │  :9200          │       │     :16686        │
    │   (Metrics)      │       │  (Logs)         │       │     (Traces)      │
    └─────────┬────────┘       └────────┬────────┘       └───────────────────┘
              │                          │
    ┌─────────┼─────────┐       ┌────────┴────────┐
    │         │         │       │    Logstash     │
    ▼         ▼         ▼       │    :5044        │
Exporters  Services  cAdvisor   └────────┬────────┘
                                         │
                                ┌────────┴────────┐
                                │    Filebeat     │
                                │  (Log shipper)  │
                                └─────────────────┘
```

## 🔧 Configuration des Microservices

Pour que les microservices exposent leurs métriques à Prometheus, chaque service Go doit:

### 1. Ajouter la dépendance (go.mod)
```go
github.com/prometheus/client_golang v1.18.0
```

### 2. Ajouter le middleware (main.go)
```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
        []string{"method", "path", "status"},
    )
)

func prometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.URL.Path == "/metrics" {
            c.Next()
            return
        }
        start := time.Now()
        c.Next()
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        path := c.FullPath()
        if path == "" {
            path = c.Request.URL.Path
        }
        httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, path, status).Observe(duration)
    }
}

// Dans main():
router.Use(prometheusMiddleware())
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## 📈 Dashboards Disponibles

### Grafana
1. **Platform Overview** - Vue d'ensemble de tous les services
2. **Service Details** - Métriques détaillées par service
3. **Database** - PostgreSQL, Redis
4. **Infrastructure** - CPU, Memory, Disk

### Kibana
1. Créer un Index Pattern: `zekora-logs-*`
2. Explorer les logs dans Discover
3. Créer des dashboards personnalisés

## 🔔 Alertes Configurées

| Alert | Sévérité | Description |
|-------|----------|-------------|
| ServiceDown | Critical | Un service ne répond plus |
| HighCPUUsage | Warning | CPU > 80% pendant 5min |
| HighMemoryUsage | Warning | Mémoire > 85% |
| LowDiskSpace | Critical | Espace disque < 15% |
| PostgreSQLDown | Critical | Base de données indisponible |
| RedisDown | Critical | Cache indisponible |
| HighErrorRate | Critical | Taux d'erreur > 5% |
| HighRequestLatency | Warning | P95 latence > 2s |

## 🐛 Troubleshooting

### Elasticsearch ne démarre pas
```bash
# Sur Linux, augmenter vm.max_map_count
sudo sysctl -w vm.max_map_count=262144

# Rendre permanent
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### Filebeat ne collecte pas les logs
```bash
# Vérifier les permissions
docker compose -f docker-compose-monitoring.yml logs filebeat
```

### Prometheus n'atteint pas les services
```bash
# Vérifier les targets
curl http://localhost:9090/api/v1/targets
```

## 📁 Structure des Fichiers

```
monitoring/
├── alertmanager/
│   └── alertmanager.yml      # Config alertes
├── filebeat/
│   └── filebeat.yml          # Config collecte logs
├── grafana/
│   ├── dashboards/           # JSON dashboards
│   └── provisioning/
│       ├── dashboards/       # Provisioning dashboards
│       └── datasources/      # Provisioning datasources
├── logstash/
│   ├── config/
│   │   └── logstash.yml      # Config Logstash
│   └── pipeline/
│       └── logstash.conf     # Pipeline processing
└── prometheus/
    ├── alerts/
    │   └── alerts.yml        # Règles d'alerte
    └── prometheus.yml        # Config scraping
```

## 🔒 Sécurité (Production)

En production, il faut:

1. **Activer l'authentification Elasticsearch**
2. **Configurer HTTPS pour Kibana**
3. **Changer les mots de passe par défaut**
4. **Restreindre l'accès réseau**
5. **Configurer les alertes email/Slack réelles**
