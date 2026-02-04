#!/bin/sh
# Initialize payment provider secrets in Vault

set -e

export VAULT_ADDR=http://vault:8200
export VAULT_TOKEN=${VAULT_TOKEN:-"dev-token-secure-2024"}

echo "========================================="
echo " Payment Provider Secrets Initialization"
echo "========================================="

# Flutterwave
echo "Configuring Flutterwave..."
vault kv put secret/payment/flutterwave \
    public_key="${FLUTTERWAVE_PUBLIC_KEY:-FLWPUBK_TEST-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-X}" \
    secret_key="${FLUTTERWAVE_SECRET_KEY:-FLWSECK_TEST-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-X}" \
    encryption_key="${FLUTTERWAVE_ENCRYPTION_KEY:-FLWSECK_TEST}" \
    webhook_secret="${FLUTTERWAVE_WEBHOOK_SECRET:-}" \
    base_url="${FLUTTERWAVE_BASE_URL:-https://api.flutterwave.com/v3}"

# CinetPay
echo "Configuring CinetPay..."
vault kv put secret/payment/cinetpay \
    api_key="${CINETPAY_API_KEY:-}" \
    site_id="${CINETPAY_SITE_ID:-}" \
    base_url="${CINETPAY_BASE_URL:-https://api-checkout.cinetpay.com/v2}"

# Paystack
echo "Configuring Paystack..."
vault kv put secret/payment/paystack \
    public_key="${PAYSTACK_PUBLIC_KEY:-pk_test_xxxxxxxxxxxxxxxxxxxx}" \
    secret_key="${PAYSTACK_SECRET_KEY:-sk_test_xxxxxxxxxxxxxxxxxxxx}" \
    base_url="${PAYSTACK_BASE_URL:-https://api.paystack.co}"

# Orange Money
echo "Configuring Orange Money..."
vault kv put secret/payment/orange_money \
    client_id="${ORANGE_MONEY_CLIENT_ID:-}" \
    client_secret="${ORANGE_MONEY_CLIENT_SECRET:-}" \
    merchant_key="${ORANGE_MONEY_MERCHANT_KEY:-}" \
    base_url="${ORANGE_MONEY_BASE_URL:-https://api.orange.com/orange-money-webpay}"

# MTN MoMo
echo "Configuring MTN MoMo..."
vault kv put secret/payment/mtn_momo \
    api_user="${MTN_MOMO_API_USER:-}" \
    api_key="${MTN_MOMO_API_KEY:-}" \
    subscription_key="${MTN_MOMO_SUBSCRIPTION_KEY:-}" \
    base_url="${MTN_MOMO_BASE_URL:-https://sandbox.momodeveloper.mtn.com}"

# Wave
echo "Configuring Wave..."
vault kv put secret/payment/wave \
    api_key="${WAVE_API_KEY:-wave_sn_prod_xxxxxxxxxxxxxxxx}" \
    base_url="${WAVE_BASE_URL:-https://api.wave.com/v1}"

# Stripe
echo "Configuring Stripe..."
vault kv put secret/payment/stripe \
    public_key="${STRIPE_PUBLIC_KEY:-pk_test_xxxxxxxxxxxxxxxxxxxx}" \
    secret_key="${STRIPE_SECRET_KEY:-sk_test_xxxxxxxxxxxxxxxxxxxx}" \
    webhook_secret="${STRIPE_WEBHOOK_SECRET:-whsec_xxxxxxxxxxxxxxxxxxxxxx}" \
    base_url="${STRIPE_BASE_URL:-https://api.stripe.com/v1}"

# PayPal
echo "Configuring PayPal..."
vault kv put secret/payment/paypal \
    client_id="${PAYPAL_CLIENT_ID:-}" \
    client_secret="${PAYPAL_CLIENT_SECRET:-}" \
    mode="${PAYPAL_MODE:-sandbox}" \
    base_url="${PAYPAL_BASE_URL:-https://api-m.sandbox.paypal.com}"

echo ""
echo "========================================="
echo " Payment Secrets Initialized Successfully"
echo "========================================="
echo ""
echo "Providers configured:"
echo "  ✓ Flutterwave (https://flutterwave.com)"
echo "  ✓ CinetPay (https://cinetpay.com)"
echo "  ✓ Paystack (https://paystack.com)"
echo "  ✓ Orange Money (https://orange.com)"
echo "  ✓ MTN MoMo (https://momodeveloper.mtn.com)"
echo "  ✓ Wave (https://wave.com)"
echo "  ✓ Stripe (https://stripe.com)"
echo "  ✓ PayPal (https://paypal.com)"
echo ""
echo "To verify secrets are stored, run:"
echo "  vault kv list secret/payment"
echo "  vault kv get secret/payment/flutterwave"
echo "========================================="
