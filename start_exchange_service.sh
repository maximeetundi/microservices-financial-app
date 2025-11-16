#!/bin/bash

# Script de démarrage complet pour le service d'échange
# Ce script configure et démarre tous les composants nécessaires

echo "🚀 Démarrage du Service d'Échange Crypto Bank"
echo "=============================================="

# Couleurs pour les messages
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérification des prérequis
check_requirements() {
    log_info "Vérification des prérequis..."
    
    # Vérifier Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker n'est pas installé"
        exit 1
    fi
    
    # Vérifier Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose n'est pas installé"
        exit 1
    fi
    
    # Vérifier Go (pour le build local)
    if ! command -v go &> /dev/null; then
        log_warn "Go n'est pas installé (requis pour le build local uniquement)"
    fi
    
    log_info "Prérequis validés ✓"
}

# Créer le réseau Docker s'il n'existe pas
setup_network() {
    log_info "Configuration du réseau Docker..."
    
    if ! docker network ls | grep -q "crypto-bank-network"; then
        docker network create crypto-bank-network
        log_info "Réseau crypto-bank-network créé ✓"
    else
        log_info "Réseau crypto-bank-network existe déjà ✓"
    fi
}

# Créer les répertoires nécessaires
setup_directories() {
    log_info "Création des répertoires..."
    
    mkdir -p services/exchange-service/logs
    mkdir -p data/postgres
    mkdir -p data/redis
    mkdir -p data/rabbitmq
    
    log_info "Répertoires créés ✓"
}

# Configuration de l'environnement
setup_environment() {
    log_info "Configuration de l'environnement..."
    
    # Copier le fichier d'exemple s'il n'existe pas
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            log_info "Fichier .env créé à partir de .env.example"
        else
            # Créer un fichier .env basique
            cat > .env << EOF
# Environment
ENVIRONMENT=development

# Database
POSTGRES_DB=crypto_bank_exchange
POSTGRES_USER=user
POSTGRES_PASSWORD=password
DATABASE_URL=postgres://user:password@localhost:5432/crypto_bank_exchange?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Services
WALLET_SERVICE_URL=http://localhost:8084
API_GATEWAY_URL=http://localhost:8080

# Exchange Fees (percentage)
CRYPTO_TO_CRYPTO_FEE=0.5
CRYPTO_TO_FIAT_FEE=0.75
FIAT_TO_CRYPTO_FEE=0.75
FIAT_TO_FIAT_FEE=0.25

# Rate Update
RATE_UPDATE_INTERVAL=30

# Ports
EXCHANGE_SERVICE_PORT=8083
POSTGRES_PORT=5432
REDIS_PORT=6379
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672
EOF
            log_info "Fichier .env créé avec les valeurs par défaut"
        fi
    else
        log_info "Fichier .env existe déjà ✓"
    fi
}

# Build du service Go
build_service() {
    log_info "Build du service d'échange..."
    
    cd services/exchange-service
    
    # Vérifier si go.mod existe
    if [ ! -f go.mod ]; then
        log_info "Initialisation du module Go..."
        go mod init github.com/crypto-bank/exchange-service
    fi
    
    # Télécharger les dépendances
    log_info "Téléchargement des dépendances Go..."
    go mod tidy
    
    # Build l'application
    log_info "Compilation de l'application..."
    go build -o bin/exchange-service main_updated.go
    
    if [ $? -eq 0 ]; then
        log_info "Build réussi ✓"
    else
        log_error "Échec du build"
        exit 1
    fi
    
    cd ../..
}

# Démarrer les services infrastructure
start_infrastructure() {
    log_info "Démarrage de l'infrastructure..."
    
    cd services/exchange-service
    
    # Démarrer PostgreSQL, Redis et RabbitMQ
    docker-compose up -d postgres redis rabbitmq
    
    log_info "Attente du démarrage des services..."
    sleep 30
    
    # Vérifier que PostgreSQL est prêt
    log_info "Vérification de PostgreSQL..."
    for i in {1..30}; do
        if docker-compose exec -T postgres pg_isready -U user -d crypto_bank_exchange; then
            log_info "PostgreSQL prêt ✓"
            break
        fi
        if [ $i -eq 30 ]; then
            log_error "PostgreSQL n'est pas prêt"
            exit 1
        fi
        sleep 2
    done
    
    cd ../..
}

# Démarrer le service d'échange
start_exchange_service() {
    log_info "Démarrage du service d'échange..."
    
    cd services/exchange-service
    
    # Option 1: Docker Compose (recommandé)
    if [ "$1" = "--docker" ]; then
        docker-compose up -d exchange-service
        log_info "Service d'échange démarré avec Docker ✓"
    else
        # Option 2: Binaire local
        ./bin/exchange-service &
        EXCHANGE_PID=$!
        echo $EXCHANGE_PID > exchange-service.pid
        log_info "Service d'échange démarré localement (PID: $EXCHANGE_PID) ✓"
    fi
    
    cd ../..
}

# Vérifier que tous les services fonctionnent
health_check() {
    log_info "Vérification de l'état des services..."
    
    # Attendre que le service soit prêt
    sleep 10
    
    # Test de santé du service d'échange
    if curl -f http://localhost:8083/health &>/dev/null; then
        log_info "Service d'échange: ✓ Opérationnel"
    else
        log_error "Service d'échange: ✗ Non disponible"
    fi
    
    # Test PostgreSQL
    if docker exec $(docker-compose -f services/exchange-service/docker-compose.yml ps -q postgres) pg_isready -U user &>/dev/null; then
        log_info "PostgreSQL: ✓ Opérationnel"
    else
        log_error "PostgreSQL: ✗ Non disponible"
    fi
    
    # Test Redis
    if docker exec $(docker-compose -f services/exchange-service/docker-compose.yml ps -q redis) redis-cli ping &>/dev/null; then
        log_info "Redis: ✓ Opérationnel"
    else
        log_error "Redis: ✗ Non disponible"
    fi
    
    # Test RabbitMQ
    if curl -f http://localhost:15673 &>/dev/null; then
        log_info "RabbitMQ Management: ✓ Opérationnel"
    else
        log_warn "RabbitMQ Management: Interface non accessible"
    fi
}

# Afficher le statut final
show_status() {
    echo ""
    echo "🎉 Service d'échange démarré avec succès!"
    echo "========================================="
    echo ""
    echo "📋 Services disponibles:"
    echo "  • Service d'échange:     http://localhost:8083"
    echo "  • Health Check:          http://localhost:8083/health"
    echo "  • API Documentation:     http://localhost:8083/docs"
    echo "  • PostgreSQL:            localhost:5432"
    echo "  • Redis:                 localhost:6379"
    echo "  • RabbitMQ Management:   http://localhost:15672"
    echo ""
    echo "🔑 Identifiants par défaut:"
    echo "  • PostgreSQL: user/password"
    echo "  • RabbitMQ:   guest/guest"
    echo ""
    echo "📖 Endpoints principaux:"
    echo "  • GET  /api/v1/rates                    - Taux de change"
    echo "  • GET  /api/v1/fiat/rates               - Taux fiat"
    echo "  • POST /api/v1/exchange/quote           - Devis d'échange"
    echo "  • POST /api/v1/exchange/execute         - Exécuter échange"
    echo "  • GET  /api/v1/trading/tickers          - Données de marché"
    echo ""
    echo "🛑 Pour arrêter:"
    echo "  docker-compose -f services/exchange-service/docker-compose.yml down"
    echo ""
}

# Fonction d'arrêt
stop_services() {
    log_info "Arrêt des services..."
    
    # Arrêter le service local s'il existe
    if [ -f services/exchange-service/exchange-service.pid ]; then
        PID=$(cat services/exchange-service/exchange-service.pid)
        kill $PID 2>/dev/null
        rm services/exchange-service/exchange-service.pid
        log_info "Service local arrêté"
    fi
    
    # Arrêter Docker Compose
    cd services/exchange-service
    docker-compose down
    cd ../..
    
    log_info "Tous les services arrêtés ✓"
}

# Gestion des signaux pour un arrêt propre
trap stop_services EXIT

# Menu principal
case "$1" in
    "start")
        check_requirements
        setup_network
        setup_directories
        setup_environment
        build_service
        start_infrastructure
        start_exchange_service $2
        health_check
        show_status
        ;;
    "stop")
        stop_services
        ;;
    "restart")
        stop_services
        sleep 5
        $0 start $2
        ;;
    "status")
        health_check
        ;;
    "logs")
        cd services/exchange-service
        docker-compose logs -f exchange-service
        cd ../..
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs} [--docker]"
        echo ""
        echo "Commandes:"
        echo "  start [--docker]  Démarrer tous les services"
        echo "  stop             Arrêter tous les services"
        echo "  restart          Redémarrer tous les services"
        echo "  status           Vérifier l'état des services"
        echo "  logs             Afficher les logs du service"
        echo ""
        echo "Options:"
        echo "  --docker         Utiliser Docker pour le service d'échange"
        echo ""
        exit 1
        ;;
esac