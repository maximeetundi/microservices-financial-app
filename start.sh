#!/bin/bash

# CryptoBank Application Startup Script

echo "🚀 Starting CryptoBank Application..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if service is ready
check_service_health() {
    local service_name=$1
    local port=$2
    local max_attempts=30
    local attempt=1

    echo -e "${YELLOW}Checking ${service_name} health...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s http://localhost:${port}/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ ${service_name} is ready!${NC}"
            return 0
        fi
        echo -e "${YELLOW}⏳ Waiting for ${service_name} (attempt $attempt/$max_attempts)...${NC}"
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ ${service_name} failed to start properly${NC}"
    return 1
}

# Function to display service URLs
show_service_urls() {
    echo -e "\n${BLUE}🌐 Service URLs:${NC}"
    echo -e "${GREEN}Frontend (Web App):${NC}     http://localhost:3000"
    echo -e "${GREEN}API Gateway:${NC}            http://localhost:8080"
    echo -e "${GREEN}Auth Service:${NC}           http://localhost:8081"
    echo -e "${GREEN}User Service:${NC}           http://localhost:8082"
    echo -e "${GREEN}Wallet Service:${NC}         http://localhost:8083"
    echo -e "${GREEN}Transfer Service:${NC}       http://localhost:8084"
    echo -e "${GREEN}Exchange Service:${NC}       http://localhost:8085"
    echo -e "${GREEN}Card Service:${NC}           http://localhost:8086"
    echo -e "${GREEN}Notification Service:${NC}   http://localhost:8087"
    echo ""
    echo -e "${BLUE}📊 Monitoring:${NC}"
    echo -e "${GREEN}Prometheus:${NC}             http://localhost:9090"
    echo -e "${GREEN}Grafana:${NC}                http://localhost:3001 (admin/admin)"
    echo -e "${GREEN}RabbitMQ Management:${NC}    http://localhost:15672 (admin/secure_password)"
    echo ""
    echo -e "${BLUE}💾 Databases:${NC}"
    echo -e "${GREEN}PostgreSQL:${NC}             localhost:5432 (admin/secure_password)"
    echo -e "${GREEN}Redis:${NC}                  localhost:6379"
}

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose >/dev/null 2>&1; then
    echo -e "${RED}❌ docker-compose is not installed. Please install it and try again.${NC}"
    exit 1
fi

echo -e "${BLUE}🐳 Starting infrastructure services...${NC}"

# Start infrastructure first (databases, message queue)
docker-compose up -d postgres redis rabbitmq

echo -e "${YELLOW}⏳ Waiting for databases to be ready...${NC}"
sleep 10

# Start core services
echo -e "${BLUE}🔧 Starting core services...${NC}"
docker-compose up -d auth-service wallet-service transfer-service

echo -e "${YELLOW}⏳ Waiting for core services...${NC}"
sleep 15

# Start API Gateway
echo -e "${BLUE}🌉 Starting API Gateway...${NC}"
docker-compose up -d api-gateway

echo -e "${YELLOW}⏳ Waiting for API Gateway...${NC}"
sleep 10

# Start frontend
echo -e "${BLUE}🎨 Starting Frontend...${NC}"
docker-compose up -d frontend

# Start monitoring services
echo -e "${BLUE}📊 Starting monitoring services...${NC}"
docker-compose up -d prometheus grafana

echo -e "\n${YELLOW}🔍 Performing health checks...${NC}"

# Health checks for core services
check_service_health "Auth Service" "8081"
check_service_health "Wallet Service" "8083" 
check_service_health "Transfer Service" "8084"
check_service_health "API Gateway" "8080"

# Check if frontend is responding
if curl -f -s http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is ready!${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend might still be loading...${NC}"
fi

echo -e "\n${GREEN}🎉 CryptoBank Application is now running!${NC}"

show_service_urls

echo -e "${BLUE}📝 Quick Start:${NC}"
echo "1. Open your browser to http://localhost:3000"
echo "2. Register a new account or login"
echo "3. Complete KYC verification"
echo "4. Create wallets and start transferring!"
echo ""
echo -e "${YELLOW}📋 To stop all services: docker-compose down${NC}"
echo -e "${YELLOW}📋 To view logs: docker-compose logs -f [service-name]${NC}"
echo -e "${YELLOW}📋 To restart a service: docker-compose restart [service-name]${NC}"

# Optional: Open browser automatically
if command -v open >/dev/null 2>&1; then
    echo -e "\n${BLUE}🌐 Opening browser...${NC}"
    sleep 3
    open http://localhost:3000
elif command -v xdg-open >/dev/null 2>&1; then
    echo -e "\n${BLUE}🌐 Opening browser...${NC}"
    sleep 3
    xdg-open http://localhost:3000
fi