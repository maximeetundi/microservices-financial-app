#!/bin/sh
# ============================================================================
#  Vault Payment Provider Secrets Initialization
# ============================================================================
#  This script initializes credentials for all payment aggregators in Vault.
#
#  IMPORTANT: The paths used here MUST match what's configured in:
#    - Admin Panel (admin-dashboard): vault_secret_path field
#    - Transfer Service: instance_loader.go
#    - Admin Service: database.go (SeedProviderInstances)
#
#  Path Format: secret/aggregators/{provider_code}/default
#  Example:     secret/aggregators/cinetpay/default
#
#  To configure REAL credentials:
#    1. Set environment variables before running this script
#    2. OR use the Admin Panel -> Agrégateurs -> [Provider] -> Credentials API
# ============================================================================

set -e

export VAULT_ADDR=${VAULT_ADDR:-"http://vault:8200"}
export VAULT_TOKEN=${VAULT_TOKEN:-"dev-token-secure-2024"}

echo "═══════════════════════════════════════════════════════════════════════"
echo "  🔐 Payment Provider Secrets Initialization"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""
echo "  Vault Address: $VAULT_ADDR"
echo "  Using paths:   secret/aggregators/{provider}/default"
echo ""

# Wait for Vault to be ready
echo "⏳ Waiting for Vault to be ready..."
until vault status > /dev/null 2>&1; do
    sleep 1
done
echo "✅ Vault is ready!"
echo ""

# Enable KV secrets engine v2 if not already enabled
echo "📦 Ensuring KV v2 secrets engine is enabled..."
vault secrets enable -path=secret -version=2 kv 2>/dev/null || true

# ============================================================================
#  CINETPAY - Popular in West Africa (UEMOA)
#  https://cinetpay.com - Get credentials at: https://app.cinetpay.com
#  DB Provider Code: cinetpay
# ============================================================================
echo "🔧 Configuring CinetPay..."
vault kv put secret/aggregators/cinetpay/default \
    api_key="${CINETPAY_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    site_id="${CINETPAY_SITE_ID:-REPLACE_WITH_YOUR_SITE_ID}" \
    secret_key="${CINETPAY_SECRET_KEY:-}" \
    base_url="${CINETPAY_BASE_URL:-https://api-checkout.cinetpay.com/v2}" \
    mode="${CINETPAY_MODE:-sandbox}"

# ============================================================================
#  WAVE - Senegal, Côte d'Ivoire, Mali, Burkina Faso
#  https://wave.com - Contact Wave for API access
#  DB Provider Codes: wave_money, wave_ci, wave_sn
# ============================================================================
echo "🔧 Configuring Wave (wave_money)..."
vault kv put secret/aggregators/wave_money/default \
    api_key="${WAVE_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    secret_key="${WAVE_SECRET_KEY:-}" \
    webhook_secret="${WAVE_WEBHOOK_SECRET:-}" \
    base_url="${WAVE_BASE_URL:-https://api.wave.com/v1}" \
    environment="${WAVE_ENVIRONMENT:-sandbox}"

# Also create aliases for wave_ci and wave_sn
vault kv put secret/aggregators/wave_ci/default \
    api_key="${WAVE_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    secret_key="${WAVE_SECRET_KEY:-}" \
    webhook_secret="${WAVE_WEBHOOK_SECRET:-}" \
    base_url="${WAVE_BASE_URL:-https://api.wave.com/v1}" \
    environment="${WAVE_ENVIRONMENT:-sandbox}"

vault kv put secret/aggregators/wave_sn/default \
    api_key="${WAVE_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    secret_key="${WAVE_SECRET_KEY:-}" \
    webhook_secret="${WAVE_WEBHOOK_SECRET:-}" \
    base_url="${WAVE_BASE_URL:-https://api.wave.com/v1}" \
    environment="${WAVE_ENVIRONMENT:-sandbox}"

# ============================================================================
#  MTN MOMO - MTN Mobile Money across Africa
#  https://momodeveloper.mtn.com - Register for API access
#  DB Provider Code: mtn_money
# ============================================================================
echo "🔧 Configuring MTN MoMo (mtn_money)..."
vault kv put secret/aggregators/mtn_money/default \
    api_user="${MTN_MOMO_API_USER:-REPLACE_WITH_YOUR_API_USER}" \
    api_key="${MTN_MOMO_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    subscription_key="${MTN_MOMO_SUBSCRIPTION_KEY:-REPLACE_WITH_YOUR_SUBSCRIPTION_KEY}" \
    base_url="${MTN_MOMO_BASE_URL:-https://sandbox.momodeveloper.mtn.com}" \
    environment="${MTN_MOMO_ENVIRONMENT:-sandbox}"

# ============================================================================
#  ORANGE MONEY - Orange Money across Africa
#  https://developer.orange.com - Register for API access
#  DB Provider Code: orange_money
# ============================================================================
echo "🔧 Configuring Orange Money (orange_money)..."
vault kv put secret/aggregators/orange_money/default \
    client_id="${ORANGE_MONEY_CLIENT_ID:-REPLACE_WITH_YOUR_CLIENT_ID}" \
    client_secret="${ORANGE_MONEY_CLIENT_SECRET:-REPLACE_WITH_YOUR_CLIENT_SECRET}" \
    merchant_key="${ORANGE_MONEY_MERCHANT_KEY:-REPLACE_WITH_YOUR_MERCHANT_KEY}" \
    base_url="${ORANGE_MONEY_BASE_URL:-https://api.orange.com/orange-money-webpay}" \
    environment="${ORANGE_MONEY_ENVIRONMENT:-sandbox}"

# ============================================================================
#  MOOV MONEY - Moov Africa Mobile Money
#  Contact Moov for API access
#  DB Provider Code: moov_money
# ============================================================================
echo "🔧 Configuring Moov Money (moov_money)..."
vault kv put secret/aggregators/moov_money/default \
    api_key="${MOOV_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    secret_key="${MOOV_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    merchant_key="${MOOV_MERCHANT_ID:-REPLACE_WITH_YOUR_MERCHANT_ID}" \
    base_url="${MOOV_BASE_URL:-https://api.moov-africa.com}"

# ============================================================================
#  FLUTTERWAVE - Pan-African payment gateway
#  https://flutterwave.com - Get credentials at: https://dashboard.flutterwave.com
#  DB Provider Code: flutterwave
# ============================================================================
echo "🔧 Configuring Flutterwave..."
vault kv put secret/aggregators/flutterwave/default \
    public_key="${FLUTTERWAVE_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${FLUTTERWAVE_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    encryption_key="${FLUTTERWAVE_ENCRYPTION_KEY:-REPLACE_WITH_YOUR_ENCRYPTION_KEY}" \
    webhook_secret="${FLUTTERWAVE_WEBHOOK_SECRET:-}" \
    base_url="${FLUTTERWAVE_BASE_URL:-https://api.flutterwave.com/v3}"

# ============================================================================
#  PAYSTACK - Nigeria, Ghana, South Africa, Kenya
#  https://paystack.com - Get credentials at: https://dashboard.paystack.com
#  DB Provider Code: paystack
# ============================================================================
echo "🔧 Configuring Paystack..."
vault kv put secret/aggregators/paystack/default \
    public_key="${PAYSTACK_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${PAYSTACK_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    webhook_secret="${PAYSTACK_WEBHOOK_SECRET:-}" \
    base_url="${PAYSTACK_BASE_URL:-https://api.paystack.co}"

# ============================================================================
#  PAWAPAY - Multi-country mobile money aggregator
#  https://pawapay.io - Contact for API access
#  DB Provider Code: pawapay
# ============================================================================
echo "🔧 Configuring PawaPay..."
vault kv put secret/aggregators/pawapay/default \
    api_key="${PAWAPAY_API_KEY:-REPLACE_WITH_YOUR_API_TOKEN}" \
    webhook_secret="${PAWAPAY_WEBHOOK_SECRET:-}" \
    base_url="${PAWAPAY_BASE_URL:-https://api.sandbox.pawapay.cloud}"

# ============================================================================
#  HUBTEL - Ghana payment gateway
#  https://hubtel.com - Get credentials at: https://developers.hubtel.com
#  DB Provider Code: hubtel
# ============================================================================
echo "🔧 Configuring Hubtel..."
vault kv put secret/aggregators/hubtel/default \
    client_id="${HUBTEL_CLIENT_ID:-REPLACE_WITH_YOUR_CLIENT_ID}" \
    client_secret="${HUBTEL_CLIENT_SECRET:-REPLACE_WITH_YOUR_CLIENT_SECRET}" \
    merchant_key="${HUBTEL_MERCHANT_ACCOUNT:-REPLACE_WITH_YOUR_MERCHANT_ACCOUNT}" \
    base_url="${HUBTEL_BASE_URL:-https://api.hubtel.com}"

# ============================================================================
#  PAYPAL - International payments
#  https://developer.paypal.com - Get credentials at: https://developer.paypal.com/dashboard
#  DB Provider Code: paypal
# ============================================================================
echo "🔧 Configuring PayPal..."
vault kv put secret/aggregators/paypal/default \
    client_id="${PAYPAL_CLIENT_ID:-REPLACE_WITH_YOUR_CLIENT_ID}" \
    client_secret="${PAYPAL_CLIENT_SECRET:-REPLACE_WITH_YOUR_CLIENT_SECRET}" \
    webhook_id="${PAYPAL_WEBHOOK_ID:-}" \
    mode="${PAYPAL_MODE:-sandbox}" \
    base_url="${PAYPAL_BASE_URL:-https://api-m.sandbox.paypal.com}"

# ============================================================================
#  STRIPE - International payments
#  https://stripe.com - Get credentials at: https://dashboard.stripe.com/apikeys
#  DB Provider Code: stripe
# ============================================================================
echo "🔧 Configuring Stripe..."
vault kv put secret/aggregators/stripe/default \
    api_key="${STRIPE_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    public_key="${STRIPE_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLISHABLE_KEY}" \
    webhook_secret="${STRIPE_WEBHOOK_SECRET:-}" \
    base_url="${STRIPE_BASE_URL:-https://api.stripe.com/v1}"

# ============================================================================
#  DEMO PROVIDER - For testing (always works)
# ============================================================================
echo "🔧 Configuring Demo Provider..."
vault kv put secret/aggregators/demo/default \
    api_key="demo_api_key_always_works" \
    mode="demo"

# ============================================================================
#  LEGACY PATHS (backward compatibility with older code)
#  These mirror the new paths for any code still using old paths
# ============================================================================
echo ""
echo "🔄 Setting up legacy paths (secret/payment/*)..."

vault kv put secret/payment/cinetpay \
    api_key="${CINETPAY_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    site_id="${CINETPAY_SITE_ID:-REPLACE_WITH_YOUR_SITE_ID}" \
    base_url="${CINETPAY_BASE_URL:-https://api-checkout.cinetpay.com/v2}"

vault kv put secret/payment/wave \
    api_key="${WAVE_API_KEY:-REPLACE_WITH_YOUR_API_KEY}" \
    base_url="${WAVE_BASE_URL:-https://api.wave.com/v1}"

vault kv put secret/payment/flutterwave \
    public_key="${FLUTTERWAVE_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${FLUTTERWAVE_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    encryption_key="${FLUTTERWAVE_ENCRYPTION_KEY:-}" \
    base_url="${FLUTTERWAVE_BASE_URL:-https://api.flutterwave.com/v3}"

vault kv put secret/payment/paystack \
    public_key="${PAYSTACK_PUBLIC_KEY:-REPLACE_WITH_YOUR_PUBLIC_KEY}" \
    secret_key="${PAYSTACK_SECRET_KEY:-REPLACE_WITH_YOUR_SECRET_KEY}" \
    base_url="${PAYSTACK_BASE_URL:-https://api.paystack.co}"

vault kv put secret/payment/orange_money \
    client_id="${ORANGE_MONEY_CLIENT_ID:-}" \
    client_secret="${ORANGE_MONEY_CLIENT_SECRET:-}" \
    merchant_key="${ORANGE_MONEY_MERCHANT_KEY:-}" \
    base_url="${ORANGE_MONEY_BASE_URL:-https://api.orange.com/orange-money-webpay}"

vault kv put secret/payment/mtn_momo \
    api_user="${MTN_MOMO_API_USER:-}" \
    api_key="${MTN_MOMO_API_KEY:-}" \
    subscription_key="${MTN_MOMO_SUBSCRIPTION_KEY:-}" \
    base_url="${MTN_MOMO_BASE_URL:-https://sandbox.momodeveloper.mtn.com}"

vault kv put secret/payment/paypal \
    client_id="${PAYPAL_CLIENT_ID:-}" \
    client_secret="${PAYPAL_CLIENT_SECRET:-}" \
    mode="${PAYPAL_MODE:-sandbox}" \
    base_url="${PAYPAL_BASE_URL:-https://api-m.sandbox.paypal.com}"

vault kv put secret/payment/stripe \
    public_key="${STRIPE_PUBLIC_KEY:-}" \
    secret_key="${STRIPE_SECRET_KEY:-}" \
    webhook_secret="${STRIPE_WEBHOOK_SECRET:-}" \
    base_url="${STRIPE_BASE_URL:-https://api.stripe.com/v1}"

# ============================================================================
#  VERIFICATION
# ============================================================================
echo ""
echo "═══════════════════════════════════════════════════════════════════════"
echo "  ✅ Payment Secrets Initialized Successfully!"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""
echo "  📋 Configured Providers:"
echo "  ─────────────────────────────────────────────────────────────────────"
echo "  │ Provider      │ Path                                    │ Region  │"
echo "  ─────────────────────────────────────────────────────────────────────"
echo "  │ CinetPay      │ secret/aggregators/cinetpay/default     │ UEMOA   │"
echo "  │ Wave          │ secret/aggregators/wave_money/default   │ W.Africa│"
echo "  │ MTN MoMo      │ secret/aggregators/mtn_money/default    │ Africa  │"
echo "  │ Orange Money  │ secret/aggregators/orange_money/default │ Africa  │"
echo "  │ Moov Money    │ secret/aggregators/moov_money/default   │ Africa  │"
echo "  │ Flutterwave   │ secret/aggregators/flutterwave/default  │ Africa  │"
echo "  │ Paystack      │ secret/aggregators/paystack/default     │ Africa  │"
echo "  │ PawaPay       │ secret/aggregators/pawapay/default      │ Africa  │"
echo "  │ Hubtel        │ secret/aggregators/hubtel/default       │ Ghana   │"
echo "  │ PayPal        │ secret/aggregators/paypal/default       │ Global  │"
echo "  │ Stripe        │ secret/aggregators/stripe/default       │ Global  │"
echo "  │ Demo          │ secret/aggregators/demo/default         │ Test    │"
echo "  ─────────────────────────────────────────────────────────────────────"
echo ""
echo "  ⚠️  IMPORTANT: Most credentials are PLACEHOLDERS!"
echo ""
echo "  🔧 TO CONFIGURE REAL CREDENTIALS:"
echo ""
echo "  Option 1: Environment Variables (before running this script)"
echo "    export CINETPAY_API_KEY='your_real_api_key'"
echo "    export CINETPAY_SITE_ID='your_real_site_id'"
echo "    ./setup-payment-secrets.sh"
echo ""
echo "  Option 2: Admin Panel (recommended for production)"
echo "    1. Go to: https://admin.yourdomain.com/dashboard/aggregators"
echo "    2. Click on a provider (e.g., CinetPay)"
echo "    3. Click 'Modifier' on an instance"
echo "    4. Go to 'Credentials API' tab"
echo "    5. Enter your real credentials and save"
echo ""
echo "  Option 3: Direct Vault CLI"
echo "    vault kv put secret/aggregators/cinetpay/default \\"
echo "      api_key='YOUR_REAL_API_KEY' \\"
echo "      site_id='YOUR_REAL_SITE_ID'"
echo ""
echo "  📖 To verify a secret:"
echo "    vault kv get secret/aggregators/cinetpay/default"
echo ""
echo "  📋 To list all aggregator secrets:"
echo "    vault kv list secret/aggregators"
echo ""
echo "═══════════════════════════════════════════════════════════════════════"
