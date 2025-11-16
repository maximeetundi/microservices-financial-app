#!/bin/bash

# CryptoBank Application Test Suite

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080/api/v1"
FRONTEND_URL="http://localhost:3000"

# Test results
TESTS_PASSED=0
TESTS_FAILED=0
ACCESS_TOKEN=""
USER_ID=""
WALLET_ID=""

echo -e "${BLUE}🧪 CryptoBank Application Test Suite${NC}"
echo "========================================"

# Function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local auth_header=""
    
    if [[ -n "$ACCESS_TOKEN" ]]; then
        auth_header="-H \"Authorization: Bearer $ACCESS_TOKEN\""
    fi
    
    if [[ "$method" == "GET" ]]; then
        response=$(eval curl -s -w "HTTPSTATUS:%{http_code}" -X GET "$API_BASE$endpoint" $auth_header)
    else
        response=$(eval curl -s -w "HTTPSTATUS:%{http_code}" -X $method "$API_BASE$endpoint" -H "Content-Type: application/json" $auth_header -d "'$data'")
    fi
    
    http_status=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    echo "$http_status|$body"
}

# Function to run test
run_test() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    
    echo -e "\n${YELLOW}Testing: $test_name${NC}"
    
    result=$(api_call "$method" "$endpoint" "$data" "$expected_status")
    status=$(echo "$result" | cut -d'|' -f1)
    body=$(echo "$result" | cut -d'|' -f2)
    
    if [[ "$status" == "$expected_status" ]]; then
        echo -e "${GREEN}✅ PASS - Status: $status${NC}"
        ((TESTS_PASSED++))
        echo "$body"
        return 0
    else
        echo -e "${RED}❌ FAIL - Expected: $expected_status, Got: $status${NC}"
        echo -e "${RED}Response: $body${NC}"
        ((TESTS_FAILED++))
        return 1
    fi
}

# Wait for services to be ready
wait_for_services() {
    echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "$API_BASE/../health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ API Gateway is ready${NC}"
            break
        fi
        echo -e "${YELLOW}Attempt $attempt/$max_attempts - waiting for API Gateway...${NC}"
        sleep 2
        ((attempt++))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        echo -e "${RED}❌ API Gateway is not responding${NC}"
        exit 1
    fi
}

# Test 1: Health Check
test_health_check() {
    echo -e "\n${BLUE}=== Health Check Tests ===${NC}"
    
    run_test "API Gateway Health" "GET" "/../health" "" "200"
    
    # Test individual service health
    services=("auth-service:8081" "wallet-service:8083" "transfer-service:8084")
    for service in "${services[@]}"; do
        service_name=$(echo $service | cut -d':' -f1)
        port=$(echo $service | cut -d':' -f2)
        
        if curl -f -s "http://localhost:$port/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $service_name is healthy${NC}"
        else
            echo -e "${RED}❌ $service_name is not responding${NC}"
        fi
    done
}

# Test 2: User Registration and Authentication
test_authentication() {
    echo -e "\n${BLUE}=== Authentication Tests ===${NC}"
    
    # Register a new user
    register_data='{
        "email": "testuser@example.com",
        "phone": "+1234567890",
        "password": "SecurePassword123!",
        "first_name": "Test",
        "last_name": "User",
        "country": "USA",
        "date_of_birth": "1990-01-01"
    }'
    
    if run_test "User Registration" "POST" "/auth/register" "$register_data" "201"; then
        echo -e "${GREEN}Registration successful${NC}"
    fi
    
    # Login
    login_data='{
        "email": "testuser@example.com",
        "password": "SecurePassword123!"
    }'
    
    if run_test "User Login" "POST" "/auth/login" "$login_data" "200"; then
        # Extract access token from response
        ACCESS_TOKEN=$(echo "$body" | grep -o '"access_token":"[^"]*"' | sed 's/"access_token":"\([^"]*\)"/\1/')
        USER_ID=$(echo "$body" | grep -o '"id":"[^"]*"' | sed 's/"id":"\([^"]*\)"/\1/')
        echo -e "${GREEN}Login successful - Token acquired${NC}"
    fi
    
    # Test invalid login
    invalid_login_data='{
        "email": "testuser@example.com",
        "password": "wrongpassword"
    }'
    
    run_test "Invalid Login" "POST" "/auth/login" "$invalid_login_data" "401"
}

# Test 3: Wallet Operations
test_wallet_operations() {
    echo -e "\n${BLUE}=== Wallet Operations Tests ===${NC}"
    
    if [[ -z "$ACCESS_TOKEN" ]]; then
        echo -e "${RED}❌ No access token - skipping wallet tests${NC}"
        return
    fi
    
    # Create USD wallet
    usd_wallet_data='{
        "currency": "USD",
        "wallet_type": "fiat",
        "name": "Test USD Wallet"
    }'
    
    if run_test "Create USD Wallet" "POST" "/wallets" "$usd_wallet_data" "201"; then
        WALLET_ID=$(echo "$body" | grep -o '"id":"[^"]*"' | sed 's/"id":"\([^"]*\)"/\1/')
        echo -e "${GREEN}USD Wallet created: $WALLET_ID${NC}"
    fi
    
    # Create BTC wallet
    btc_wallet_data='{
        "currency": "BTC",
        "wallet_type": "crypto",
        "name": "Test BTC Wallet"
    }'
    
    run_test "Create BTC Wallet" "POST" "/wallets" "$btc_wallet_data" "201"
    
    # Get user wallets
    run_test "Get User Wallets" "GET" "/wallets" "" "200"
    
    # Get wallet balance
    if [[ -n "$WALLET_ID" ]]; then
        run_test "Get Wallet Balance" "GET" "/wallets/$WALLET_ID/balance" "" "200"
    fi
}

# Test 4: Transfer Operations
test_transfer_operations() {
    echo -e "\n${BLUE}=== Transfer Operations Tests ===${NC}"
    
    if [[ -z "$ACCESS_TOKEN" || -z "$WALLET_ID" ]]; then
        echo -e "${RED}❌ Missing access token or wallet ID - skipping transfer tests${NC}"
        return
    fi
    
    # Test transfer creation (should fail due to insufficient balance)
    transfer_data='{
        "from_wallet_id": "'$WALLET_ID'",
        "to_email": "recipient@example.com",
        "amount": 100.00,
        "currency": "USD",
        "description": "Test transfer"
    }'
    
    # This should fail with insufficient balance
    run_test "Create Transfer (Insufficient Balance)" "POST" "/transfers" "$transfer_data" "400"
    
    # Get transfer history
    run_test "Get Transfer History" "GET" "/transfers" "" "200"
    
    # Get mobile money providers
    run_test "Get Mobile Money Providers" "GET" "/transfers/mobile/providers" "" "200"
}

# Test 5: Exchange Operations
test_exchange_operations() {
    echo -e "\n${BLUE}=== Exchange Operations Tests ===${NC}"
    
    # Get exchange rates (should work without auth)
    run_test "Get Exchange Rates" "GET" "/exchange/rates" "" "200"
    
    # Get specific rate
    run_test "Get BTC/USD Rate" "GET" "/exchange/rates/BTC/USD" "" "200"
    
    if [[ -n "$ACCESS_TOKEN" ]]; then
        # Get quote
        quote_data='{
            "from_currency": "BTC",
            "to_currency": "USD",
            "amount": 0.01
        }'
        
        run_test "Get Exchange Quote" "POST" "/exchange/quote" "$quote_data" "200"
    fi
}

# Test 6: Frontend Accessibility
test_frontend() {
    echo -e "\n${BLUE}=== Frontend Tests ===${NC}"
    
    if curl -f -s "$FRONTEND_URL" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend is accessible${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ Frontend is not accessible${NC}"
        ((TESTS_FAILED++))
    fi
    
    # Test static assets
    if curl -f -s "$FRONTEND_URL/_nuxt/" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Static assets are loading${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${YELLOW}⚠️  Static assets may not be ready yet${NC}"
    fi
}

# Test 7: Security Tests
test_security() {
    echo -e "\n${BLUE}=== Security Tests ===${NC}"
    
    # Test access without token
    run_test "Access Protected Route Without Token" "GET" "/wallets" "" "401"
    
    # Test SQL injection attempt
    sql_injection_data='{
        "email": "test@example.com",
        "password": "password OR 1=1"
    }'
    
    run_test "SQL Injection Attempt" "POST" "/auth/login" "$sql_injection_data" "401"
    
    # Test XSS attempt
    xss_data='{
        "email": "test@example.com",
        "password": "<script>alert(\"xss\")</script>"
    }'
    
    run_test "XSS Attempt" "POST" "/auth/login" "$xss_data" "401"
}

# Performance test
test_performance() {
    echo -e "\n${BLUE}=== Performance Tests ===${NC}"
    
    echo -e "${YELLOW}Testing API response time...${NC}"
    
    start_time=$(date +%s%N)
    curl -s "$API_BASE/../health" > /dev/null
    end_time=$(date +%s%N)
    
    duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if [[ $duration -lt 1000 ]]; then
        echo -e "${GREEN}✅ API response time: ${duration}ms (Good)${NC}"
        ((TESTS_PASSED++))
    elif [[ $duration -lt 3000 ]]; then
        echo -e "${YELLOW}⚠️  API response time: ${duration}ms (Acceptable)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ API response time: ${duration}ms (Too slow)${NC}"
        ((TESTS_FAILED++))
    fi
}

# Main test execution
main() {
    echo -e "${BLUE}Starting comprehensive test suite...${NC}"
    
    wait_for_services
    
    test_health_check
    test_authentication
    test_wallet_operations
    test_transfer_operations
    test_exchange_operations
    test_frontend
    test_security
    test_performance
    
    # Test summary
    echo -e "\n${BLUE}=== Test Summary ===${NC}"
    echo "================================"
    echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
    
    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "\n${GREEN}🎉 All tests passed! CryptoBank is working correctly.${NC}"
        exit 0
    else
        echo -e "\n${RED}⚠️  Some tests failed. Please check the output above.${NC}"
        exit 1
    fi
}

# Run tests
main