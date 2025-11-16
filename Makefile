# CryptoBank Application Makefile

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE := docker-compose
PROJECT_NAME := crypto-bank
FRONTEND_URL := http://localhost:3000
API_URL := http://localhost:8080

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)CryptoBank Application - Available Commands$(NC)"
	@echo "=============================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
.PHONY: install
install: ## Install dependencies and setup environment
	@echo "$(BLUE)Setting up CryptoBank development environment...$(NC)"
	@cp .env.example .env
	@chmod +x start.sh test.sh
	@echo "$(GREEN)✅ Environment setup complete!$(NC)"
	@echo "$(YELLOW)Please edit .env file with your configuration$(NC)"

.PHONY: start
start: ## Start all services
	@echo "$(BLUE)Starting CryptoBank application...$(NC)"
	@./start.sh

.PHONY: stop
stop: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✅ All services stopped$(NC)"

.PHONY: restart
restart: stop start ## Restart all services

.PHONY: dev
dev: ## Start in development mode with hot reload
	@echo "$(BLUE)Starting in development mode...$(NC)"
	@$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.development.yml up -d
	@echo "$(GREEN)✅ Development environment started$(NC)"

# Service management
.PHONY: start-infrastructure
start-infrastructure: ## Start only infrastructure (DB, Redis, RabbitMQ)
	@echo "$(BLUE)Starting infrastructure services...$(NC)"
	@$(DOCKER_COMPOSE) up -d postgres redis rabbitmq
	@echo "$(GREEN)✅ Infrastructure services started$(NC)"

.PHONY: start-backend
start-backend: ## Start only backend services
	@echo "$(BLUE)Starting backend services...$(NC)"
	@$(DOCKER_COMPOSE) up -d api-gateway auth-service wallet-service transfer-service
	@echo "$(GREEN)✅ Backend services started$(NC)"

.PHONY: start-frontend
start-frontend: ## Start only frontend
	@echo "$(BLUE)Starting frontend...$(NC)"
	@$(DOCKER_COMPOSE) up -d frontend
	@echo "$(GREEN)✅ Frontend started$(NC)"

.PHONY: start-monitoring
start-monitoring: ## Start monitoring services
	@echo "$(BLUE)Starting monitoring services...$(NC)"
	@$(DOCKER_COMPOSE) up -d prometheus grafana
	@echo "$(GREEN)✅ Monitoring services started$(NC)"

# Logs and debugging
.PHONY: logs
logs: ## Show logs for all services
	@$(DOCKER_COMPOSE) logs -f

.PHONY: logs-api
logs-api: ## Show API Gateway logs
	@$(DOCKER_COMPOSE) logs -f api-gateway

.PHONY: logs-auth
logs-auth: ## Show Auth Service logs
	@$(DOCKER_COMPOSE) logs -f auth-service

.PHONY: logs-wallet
logs-wallet: ## Show Wallet Service logs
	@$(DOCKER_COMPOSE) logs -f wallet-service

.PHONY: logs-transfer
logs-transfer: ## Show Transfer Service logs
	@$(DOCKER_COMPOSE) logs -f transfer-service

.PHONY: logs-frontend
logs-frontend: ## Show Frontend logs
	@$(DOCKER_COMPOSE) logs -f frontend

.PHONY: logs-db
logs-db: ## Show database logs
	@$(DOCKER_COMPOSE) logs -f postgres

# Testing
.PHONY: test
test: ## Run comprehensive test suite
	@echo "$(BLUE)Running test suite...$(NC)"
	@chmod +x test.sh
	@./test.sh

.PHONY: test-api
test-api: ## Test API endpoints only
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	@curl -f $(API_URL)/health && echo "$(GREEN)✅ API Gateway healthy$(NC)" || echo "$(RED)❌ API Gateway not responding$(NC)"
	@curl -f http://localhost:8081/health && echo "$(GREEN)✅ Auth Service healthy$(NC)" || echo "$(RED)❌ Auth Service not responding$(NC)"
	@curl -f http://localhost:8083/health && echo "$(GREEN)✅ Wallet Service healthy$(NC)" || echo "$(RED)❌ Wallet Service not responding$(NC)"

.PHONY: test-frontend
test-frontend: ## Test frontend accessibility
	@echo "$(BLUE)Testing frontend...$(NC)"
	@curl -f $(FRONTEND_URL) > /dev/null && echo "$(GREEN)✅ Frontend accessible$(NC)" || echo "$(RED)❌ Frontend not accessible$(NC)"

# Database operations
.PHONY: db-connect
db-connect: ## Connect to PostgreSQL database
	@echo "$(BLUE)Connecting to database...$(NC)"
	@$(DOCKER_COMPOSE) exec postgres psql -U admin crypto_bank

.PHONY: db-backup
db-backup: ## Backup database
	@echo "$(BLUE)Creating database backup...$(NC)"
	@mkdir -p backups
	@$(DOCKER_COMPOSE) exec postgres pg_dump -U admin crypto_bank > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ Database backup created in backups/ directory$(NC)"

.PHONY: db-restore
db-restore: ## Restore database (requires BACKUP_FILE variable)
	@if [ -z "$(BACKUP_FILE)" ]; then echo "$(RED)❌ Please specify BACKUP_FILE=path/to/backup.sql$(NC)"; exit 1; fi
	@echo "$(BLUE)Restoring database from $(BACKUP_FILE)...$(NC)"
	@$(DOCKER_COMPOSE) exec -T postgres psql -U admin crypto_bank < $(BACKUP_FILE)
	@echo "$(GREEN)✅ Database restored$(NC)"

.PHONY: db-reset
db-reset: ## Reset database (WARNING: destroys all data)
	@echo "$(RED)WARNING: This will destroy all data!$(NC)"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d postgres
	@echo "$(GREEN)✅ Database reset complete$(NC)"

# Cache operations
.PHONY: redis-cli
redis-cli: ## Connect to Redis CLI
	@echo "$(BLUE)Connecting to Redis...$(NC)"
	@$(DOCKER_COMPOSE) exec redis redis-cli

.PHONY: cache-clear
cache-clear: ## Clear Redis cache
	@echo "$(BLUE)Clearing Redis cache...$(NC)"
	@$(DOCKER_COMPOSE) exec redis redis-cli FLUSHALL
	@echo "$(GREEN)✅ Cache cleared$(NC)"

# Service status and monitoring
.PHONY: status
status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	@$(DOCKER_COMPOSE) ps

.PHONY: health
health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo "API Gateway:" && curl -f $(API_URL)/health > /dev/null && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Unhealthy$(NC)"
	@echo "Auth Service:" && curl -f http://localhost:8081/health > /dev/null && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Unhealthy$(NC)"
	@echo "Wallet Service:" && curl -f http://localhost:8083/health > /dev/null && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Unhealthy$(NC)"
	@echo "Transfer Service:" && curl -f http://localhost:8084/health > /dev/null && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Unhealthy$(NC)"
	@echo "Frontend:" && curl -f $(FRONTEND_URL) > /dev/null && echo "$(GREEN)✅ Healthy$(NC)" || echo "$(RED)❌ Unhealthy$(NC)"

.PHONY: stats
stats: ## Show resource usage statistics
	@echo "$(BLUE)Resource Usage:$(NC)"
	@$(DOCKER_COMPOSE) top

# Cleanup operations
.PHONY: clean
clean: ## Remove containers and networks
	@echo "$(YELLOW)Cleaning up containers and networks...$(NC)"
	@$(DOCKER_COMPOSE) down --remove-orphans
	@echo "$(GREEN)✅ Cleanup complete$(NC)"

.PHONY: clean-volumes
clean-volumes: ## Remove containers, networks, and volumes (WARNING: destroys data)
	@echo "$(RED)WARNING: This will destroy all data including databases!$(NC)"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "$(GREEN)✅ Complete cleanup finished$(NC)"

.PHONY: clean-images
clean-images: ## Remove all project Docker images
	@echo "$(YELLOW)Removing project Docker images...$(NC)"
	@docker images | grep crypto-bank | awk '{print $$3}' | xargs -r docker rmi
	@echo "$(GREEN)✅ Images removed$(NC)"

# Build operations
.PHONY: build
build: ## Build all Docker images
	@echo "$(BLUE)Building all Docker images...$(NC)"
	@$(DOCKER_COMPOSE) build
	@echo "$(GREEN)✅ All images built successfully$(NC)"

.PHONY: build-no-cache
build-no-cache: ## Build all Docker images without cache
	@echo "$(BLUE)Building all Docker images without cache...$(NC)"
	@$(DOCKER_COMPOSE) build --no-cache
	@echo "$(GREEN)✅ All images built successfully$(NC)"

.PHONY: pull
pull: ## Pull latest base images
	@echo "$(BLUE)Pulling latest base images...$(NC)"
	@$(DOCKER_COMPOSE) pull
	@echo "$(GREEN)✅ Base images updated$(NC)"

# Security operations
.PHONY: security-scan
security-scan: ## Run security scan on Docker images
	@echo "$(BLUE)Running security scan...$(NC)"
	@command -v trivy >/dev/null 2>&1 || (echo "$(RED)❌ Trivy not installed. Install with: curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/master/contrib/install.sh | sh$(NC)" && exit 1)
	@trivy image crypto-bank-app_api-gateway:latest
	@trivy image crypto-bank-app_auth-service:latest
	@trivy image crypto-bank-app_wallet-service:latest

# Production operations
.PHONY: prod-build
prod-build: ## Build images for production
	@echo "$(BLUE)Building production images...$(NC)"
	@$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.prod.yml build
	@echo "$(GREEN)✅ Production images built$(NC)"

.PHONY: prod-deploy
prod-deploy: ## Deploy to production (requires proper configuration)
	@echo "$(RED)WARNING: This will deploy to production!$(NC)"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "$(BLUE)Deploying to production...$(NC)"
	@$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "$(GREEN)✅ Production deployment complete$(NC)"

# Utility commands
.PHONY: open
open: ## Open application URLs in browser
	@echo "$(BLUE)Opening CryptoBank in browser...$(NC)"
	@command -v open >/dev/null 2>&1 && open $(FRONTEND_URL) || \
	 command -v xdg-open >/dev/null 2>&1 && xdg-open $(FRONTEND_URL) || \
	 echo "$(YELLOW)Please open $(FRONTEND_URL) in your browser$(NC)"

.PHONY: docs
docs: ## Show important URLs and documentation
	@echo "$(BLUE)📚 CryptoBank Documentation$(NC)"
	@echo "================================"
	@echo "$(GREEN)Frontend:$(NC)          $(FRONTEND_URL)"
	@echo "$(GREEN)API Gateway:$(NC)       $(API_URL)"
	@echo "$(GREEN)API Documentation:$(NC) $(API_URL)/docs"
	@echo "$(GREEN)Grafana:$(NC)           http://localhost:3001"
	@echo "$(GREEN)Prometheus:$(NC)        http://localhost:9090"
	@echo "$(GREEN)RabbitMQ:$(NC)          http://localhost:15672"
	@echo ""
	@echo "$(BLUE)📖 Quick Commands:$(NC)"
	@echo "$(YELLOW)make start$(NC)     - Start all services"
	@echo "$(YELLOW)make test$(NC)      - Run test suite"
	@echo "$(YELLOW)make logs$(NC)      - View all logs"
	@echo "$(YELLOW)make stop$(NC)      - Stop all services"

.PHONY: env-check
env-check: ## Check environment configuration
	@echo "$(BLUE)Checking environment configuration...$(NC)"
	@if [ -f .env ]; then \
		echo "$(GREEN)✅ .env file exists$(NC)"; \
		grep -q "POSTGRES_PASSWORD=secure_password" .env && echo "$(RED)⚠️  Using default database password$(NC)"; \
		grep -q "JWT_SECRET=ultra_secure" .env && echo "$(RED)⚠️  Using default JWT secret$(NC)"; \
	else \
		echo "$(RED)❌ .env file not found. Run 'make install' first$(NC)"; \
	fi