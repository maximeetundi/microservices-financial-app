#!/bin/sh

set -e

export VAULT_ADDR=http://vault:8200

# Wait for Vault to start
echo "Waiting for Vault to start..."
until wget -q --spider http://vault:8200/v1/sys/seal-status; do
  echo "Valut is not ready yet..."
  sleep 2
done

echo "Vault is reachable."

# Check if initialized
if wget -q -O- http://vault:8200/v1/sys/health | grep -q '"initialized":true'; then
    echo "Vault is already initialized."
else
    echo "Initializing Vault..."
    # Initialize and save keys to file (SECURE THIS FILE IN PROD!)
    vault operator init -key-shares=1 -key-threshold=1 -format=json > /vault/file/init-keys.json
    echo "Vault initialized. Keys saved to /vault/file/init-keys.json"
fi

# Install jquery for JSON parsing
echo "Installing jq..."
apk add --no-cache jq

# Read keys
if [ -f "/vault/file/init-keys.json" ]; then
    echo "Found init-keys.json. Content:"
    cat /vault/file/init-keys.json

    echo "Attempting to parse keys..."
    UNSEAL_KEY=$(jq -r ".unseal_keys_b64[0]" /vault/file/init-keys.json)
    ROOT_TOKEN=$(jq -r ".root_token" /vault/file/init-keys.json)

    # Fallback if jq failed or returned null
    if [ -z "$UNSEAL_KEY" ] || [ "$UNSEAL_KEY" = "null" ]; then
        echo "jq failed. Trying fallback grep/sed..."
        UNSEAL_KEY=$(grep "unseal_keys_b64" /vault/file/init-keys.json | sed -E 's/.*"unseal_keys_b64":\["([^"]+)".*/\1/')
        ROOT_TOKEN=$(grep "root_token" /vault/file/init-keys.json | sed -E 's/.*"root_token":"([^"]+)".*/\1/')
    fi
else
    echo "Error: init-keys.json not found!"
    exit 1
fi

if [ -z "$UNSEAL_KEY" ] || [ "$UNSEAL_KEY" = "null" ]; then
    echo "CRITICAL ERROR: Failed to extract UNSEAL_KEY. Aborting."
    exit 1
fi

echo "Unseal Key found: ${UNSEAL_KEY:0:5}..." # Print only first 5 chars for verification

# Unseal
echo "Unsealing Vault..."
vault operator unseal $UNSEAL_KEY || echo "Already unsealed or error"

# Login
echo "Logging in..."
export VAULT_TOKEN=$ROOT_TOKEN
vault login $ROOT_TOKEN

# Enable KV engine if not enabled
echo "Enabling KV secrets engine..."
vault secrets enable -path=secret kv || echo "KV engine already enabled or error"

# Write secrets (Example: Tatum Key)
echo "Writing secrets..."
# Note: In a real script we would pass these via ENV vars to this script
if [ -n "$TATUM_API_KEY" ]; then
    vault kv put secret/wallet-service tatum_api_key="$TATUM_API_KEY"
    echo "Tatum API Key written to secret/wallet-service"
fi

if [ -n "$FIXER_API_KEY" ]; then
    vault kv put secret/exchange-service fixer_api_key="$FIXER_API_KEY" currency_layer_api_key="$CURRENCYLAYER_API_KEY"
    echo "Exchange keys written to secret/exchange-service"
    echo "Exchange keys written to secret/exchange-service"
fi

# Generate WALLET_MASTER_KEY if not exists
echo "Checking for WALLET_MASTER_KEY..."
EXISTING_KEY=$(vault kv get -field=wallet_master_key secret/wallet-service 2>/dev/null || echo "")

if [ -z "$EXISTING_KEY" ]; then
    echo "Generating new WALLET_MASTER_KEY..."

    # Method 1: Try openssl (most reliable)
    if command -v openssl > /dev/null 2>&1; then
        NEW_KEY=$(openssl rand -hex 32)
        echo "Generated key using openssl"
    # Method 2: Try /dev/urandom with dd
    elif [ -r /dev/urandom ]; then
        NEW_KEY=$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | xxd -p | tr -d '\n ')
        if [ -z "$NEW_KEY" ] || [ ${#NEW_KEY} -lt 64 ]; then
            # xxd not available, try od
            NEW_KEY=$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | od -An -tx1 | tr -d ' \n')
        fi
        echo "Generated key using /dev/urandom"
    # Method 3: Fallback to hexdump
    else
        NEW_KEY=$(hexdump -vn32 -e'4/4 "%08X" 1 "\n"' /dev/urandom | tr -d ' \n' | tr '[:upper:]' '[:lower:]')
        echo "Generated key using hexdump fallback"
    fi

    # Validate key length
    if [ -z "$NEW_KEY" ] || [ ${#NEW_KEY} -lt 64 ]; then
        echo "ERROR: Failed to generate valid WALLET_MASTER_KEY (got ${#NEW_KEY} chars, need 64)"
        exit 1
    fi

    # Truncate to exactly 64 chars if longer
    NEW_KEY=$(echo "$NEW_KEY" | head -c 64)

    # Also generate WALLET_SECRET and WALLET_SALT for backward compatibility
    if command -v openssl > /dev/null 2>&1; then
        NEW_SECRET=$(openssl rand -base64 32)
        NEW_SALT=$(openssl rand -base64 16)
    else
        NEW_SECRET=$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | base64)
        NEW_SALT=$(dd if=/dev/urandom bs=16 count=1 2>/dev/null | base64)
    fi

    # Write all secrets to Vault
    vault kv put secret/wallet-service \
        tatum_api_key="${TATUM_API_KEY}" \
        wallet_master_key="${NEW_KEY}" \
        wallet_secret="${NEW_SECRET}" \
        wallet_salt="${NEW_SALT}"

    echo "========================================"
    echo "  SECURITY KEYS GENERATED SUCCESSFULLY"
    echo "========================================"
    echo "WALLET_MASTER_KEY: ${NEW_KEY:0:16}...${NEW_KEY:48:16}"
    echo "WALLET_SECRET: ${NEW_SECRET:0:10}..."
    echo "WALLET_SALT: ${NEW_SALT:0:10}..."
    echo ""
    echo "Keys stored in HashiCorp Vault at: secret/wallet-service"
    echo "========================================"
else
    echo "WALLET_MASTER_KEY already exists in Vault (${EXISTING_KEY:0:8}...)"
fi

# Create a dev token for services
echo "Creating dev token..."
if vault token lookup dev-token-secure-2024 > /dev/null 2>&1; then
    echo "Dev token already exists"
else
    vault token create -id="dev-token-secure-2024" -policy=root || echo "Failed to create token"
    echo "Dev token created: dev-token-secure-2024"
fi

# ============================================================================
# Initialize Payment Provider Secrets
# ============================================================================
echo ""
echo "=========================================="
echo "  Initializing Payment Provider Secrets"
echo "=========================================="

# CinetPay - West Africa (UEMOA)
echo "🔧 Configuring CinetPay..."
vault kv put secret/aggregators/cinetpay \
    api_key="${CINETPAY_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    site_id="${CINETPAY_SITE_ID:-REPLACE_WITH_YOUR_SITE_ID}" \
    secret_key="${CINETPAY_SECRET_KEY:-}" \
    base_url="https://api-checkout.cinetpay.com/v2" \
    mode="${CINETPAY_MODE:-sandbox}" || echo "CinetPay config failed"

# Wave - Senegal, Côte d'Ivoire
echo "🔧 Configuring Wave..."
vault kv put secret/aggregators/wave_money \
    api_key="${WAVE_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    secret_key="${WAVE_SECRET_KEY:-}" \
    webhook_secret="${WAVE_WEBHOOK_SECRET:-}" \
    base_url="https://api.wave.com/v1" \
    environment="${WAVE_ENVIRONMENT:-sandbox}" || echo "Wave config failed"

# MTN MoMo
echo "🔧 Configuring MTN MoMo..."
vault kv put secret/aggregators/mtn_money \
    api_user="${MTN_MOMO_API_USER:-REPLACE_WITH_YOUR_API_USER}" \
    api_key="${MTN_MOMO_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    subscription_key="${MTN_MOMO_SUBSCRIPTION_KEY:-REPLACE_WITH_YOUR_SUBSCRIPTION_KEY}" \
    base_url="https://sandbox.momodeveloper.mtn.com" \
    environment="${MTN_MOMO_ENVIRONMENT:-sandbox}" || echo "MTN MoMo config failed"

# Orange Money
echo "🔧 Configuring Orange Money..."
vault kv put secret/aggregators/orange_money \
    client_id="${ORANGE_MONEY_CLIENT_ID:-REPLACE_WITH_YOUR_CLIENT_ID}" \
    client_secret="${ORANGE_MONEY_CLIENT_SECRET:-REPLACE_WITH_YOUR_CLIENT_SECRET}" \
    merchant_key="${ORANGE_MONEY_MERCHANT_KEY:-REPLACE_WITH_YOUR_MERCHANT_KEY}" \
    base_url="https://api.orange.com/orange-money-webpay" \
    environment="${ORANGE_MONEY_ENVIRONMENT:-sandbox}" || echo "Orange Money config failed"

# Flutterwave
echo "🔧 Configuring Flutterwave..."
vault kv put secret/aggregators/flutterwave \
    public_key="${FLUTTERWAVE_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${FLUTTERWAVE_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    encryption_key="${FLUTTERWAVE_ENCRYPTION_KEY:-}" \
    webhook_secret="${FLUTTERWAVE_WEBHOOK_SECRET:-}" \
    base_url="https://api.flutterwave.com/v3" || echo "Flutterwave config failed"

# Paystack
echo "🔧 Configuring Paystack..."
vault kv put secret/aggregators/paystack \
    public_key="${PAYSTACK_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${PAYSTACK_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    webhook_secret="${PAYSTACK_WEBHOOK_SECRET:-}" \
    base_url="https://api.paystack.co" || echo "Paystack config failed"

# PayPal
echo "🔧 Configuring PayPal..."
vault kv put secret/aggregators/paypal \
    client_id="${PAYPAL_CLIENT_ID:-REPLACE_WITH_YOUR_CLIENT_ID}" \
    client_secret="${PAYPAL_CLIENT_SECRET:-REPLACE_WITH_YOUR_CLIENT_SECRET}" \
    webhook_id="${PAYPAL_WEBHOOK_ID:-}" \
    mode="${PAYPAL_MODE:-sandbox}" \
    base_url="https://api-m.sandbox.paypal.com" || echo "PayPal config failed"

# Stripe
echo "🔧 Configuring Stripe..."
vault kv put secret/aggregators/stripe \
    api_key="${STRIPE_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    public_key="${STRIPE_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLISHABLE_KEY}" \
    webhook_secret="${STRIPE_WEBHOOK_SECRET:-}" \
    base_url="https://api.stripe.com/v1" || echo "Stripe config failed"

# Demo Provider (always works)
echo "🔧 Configuring Demo Provider..."
vault kv put secret/aggregators/demo \
    api_key="demo_api_key_always_works" \
    mode="demo" || echo "Demo config failed"

echo ""
echo "✅ Payment provider secrets initialized!"
echo ""
echo "=========================================="

echo "Vault setup complete!"
echo "Root Token: $ROOT_TOKEN"
echo "Unseal Key: $UNSEAL_KEY"
echo "Service Token: dev-token-secure-2024"
echo ""
echo "📋 Vault paths for aggregators:"
echo "   secret/aggregators/cinetpay"
echo "   secret/aggregators/wave_money"
echo "   secret/aggregators/mtn_money"
echo "   secret/aggregators/orange_money"
echo "   secret/aggregators/flutterwave"
echo "   secret/aggregators/paystack"
echo "   secret/aggregators/paypal"
echo "   secret/aggregators/stripe"
echo "   secret/aggregators/demo"
echo ""
echo "⚠️  Replace PLACEHOLDER credentials with real API keys!"
