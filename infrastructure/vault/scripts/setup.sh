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
    # Generate 32 bytes (64 hex chars)
    NEW_KEY=$(hexdump -vn32 -e'4/4 "%08X" 1 "\n"' /dev/urandom | tr -d ' \n' | tr '[:upper:]' '[:lower:]')
    
    # If hexdump fails (some alpine versions), try od or openssl
    if [ -z "$NEW_KEY" ] || [ ${#NEW_KEY} -lt 64 ]; then
        echo "hexdump failed or insufficient, trying od..."
        NEW_KEY=$(od -N 32 -t x1 /dev/urandom | head -n 2 | awk '{print $2$3$4$5$6$7$8$9$10$11$12$13$14$15$16}' | tr -d '\n ')
    fi
    
    # Ensure we preserve existing secrets when writing new one
    # Note: vault kv put replaces all keys in the path by default for v1, or adds version for v2.
    # We should read existing, merge, and write back. 
    # But since we just wrote tatum_api_key above, we can just write both.
    
    # Better approach: Use patch if available or just re-write with both
    # For now, let's just write/patch.
    
    # Actually, let's just use "vault kv patch" if available (v2) or "vault kv put" merging data.
    # The image is hashicorp/vault:latest, likely supports kv put directly.
    # Re-writing tatum_key if needed is safer.
    
    vault kv put secret/wallet-service \
        tatum_api_key="${TATUM_API_KEY}" \
        wallet_master_key="${NEW_KEY}"
        
    echo "Generated and injected WALLET_MASTER_KEY: ${NEW_KEY:0:10}..."
else
    echo "WALLET_MASTER_KEY already exists in Vault."
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
