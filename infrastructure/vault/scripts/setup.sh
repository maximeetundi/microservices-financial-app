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
apk add --no-cache jq > /dev/null 2>&1

# Read keys
if [ -f "/vault/file/init-keys.json" ]; then
    UNSEAL_KEY=$(jq -r ".unseal_keys_b64[0]" /vault/file/init-keys.json)
    ROOT_TOKEN=$(jq -r ".root_token" /vault/file/init-keys.json)
else
    echo "Error: init-keys.json not found!"
    exit 1
fi

if [ -z "$UNSEAL_KEY" ] || [ "$UNSEAL_KEY" = "null" ]; then
    echo "Error: Failed to parse UNSEAL_KEY. Check init-keys.json content."
    cat /vault/file/init-keys.json
    exit 1
fi

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
fi

# Create a dev token for services
echo "Creating dev token..."
if vault token lookup dev-token-secure-2024 > /dev/null 2>&1; then
    echo "Dev token already exists"
else
    vault token create -id="dev-token-secure-2024" -policy=root || echo "Failed to create token"
    echo "Dev token created: dev-token-secure-2024"
fi

echo "Vault setup complete!"
echo "Root Token: $ROOT_TOKEN"
echo "Unseal Key: $UNSEAL_KEY"
echo "Service Token: dev-token-secure-2024"
