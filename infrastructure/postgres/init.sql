-- Crypto Bank Database Schema
-- Optimized for reliability: Schema first, Data second

-- Create admin database
CREATE DATABASE crypto_bank_admin;

-- =====================================================
-- 1. TABLES & SCHEMA DEFINITIONS
-- =====================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    country VARCHAR(3) NOT NULL, -- ISO country code
    kyc_status VARCHAR(20) DEFAULT 'none', -- none, pending, verified, rejected
    kyc_level INTEGER DEFAULT 1, -- 1, 2, 3 (verification levels)
    role VARCHAR(20) DEFAULT 'user', -- user, admin, support
    is_active BOOLEAN DEFAULT true,
    two_fa_enabled BOOLEAN DEFAULT false,
    two_fa_secret VARCHAR(64),
    email_verified BOOLEAN DEFAULT false,
    phone_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP,
    failed_attempts INTEGER DEFAULT 0,
    total_balance_usd DECIMAL(20,8) DEFAULT 0,
    locked_until TIMESTAMP,
    -- PIN fields (5-digit transaction security PIN)
    pin_hash VARCHAR(255),
    pin_set_at TIMESTAMP,
    pin_failed_attempts INTEGER DEFAULT 0,
    pin_locked_until TIMESTAMP,
    pin_permanently_locked BOOLEAN DEFAULT FALSE,
    pin_temp_lock_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Backup codes for 2FA recovery
CREATE TABLE IF NOT EXISTS backup_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL, -- Hashed backup code
    used BOOLEAN DEFAULT false,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verification tokens for email/phone/password reset
CREATE TABLE IF NOT EXISTS verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL, -- email_verification, phone_verification, password_reset
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency VARCHAR(10) NOT NULL, -- USD, EUR, BTC, ETH, etc.
    wallet_type VARCHAR(20) NOT NULL, -- fiat, crypto
    balance DECIMAL(20,8) DEFAULT 0,
    frozen_balance DECIMAL(20,8) DEFAULT 0,
    wallet_address VARCHAR(255), -- For crypto wallets
    private_key_encrypted TEXT, -- Encrypted private key
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, currency)
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_wallet_id UUID REFERENCES wallets(id),
    to_wallet_id UUID REFERENCES wallets(id),
    transaction_type VARCHAR(20) NOT NULL, -- transfer, exchange, deposit, withdrawal
    amount DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) DEFAULT 0,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed, cancelled
    blockchain_tx_hash VARCHAR(255), -- For crypto transactions
    reference_id VARCHAR(100), -- External reference
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Deposit Transactions table (for tracking deposit flow with aggregators)
CREATE TABLE IF NOT EXISTS deposit_transactions (
    id VARCHAR(100) PRIMARY KEY, -- dep_xxxx_timestamp format
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Amount details
    amount DECIMAL(20,8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    fee DECIMAL(20,8) DEFAULT 0,
    net_amount DECIMAL(20,8), -- Amount after fees

    -- Provider/Aggregator info
    provider_code VARCHAR(50) NOT NULL, -- lygos, flutterwave, stripe, etc.
    aggregator_instance_id UUID, -- Reference to aggregator_instances
    hot_wallet_id VARCHAR(36), -- Reference to platform_accounts

    -- Payment details
    payment_url TEXT, -- URL where user is redirected to pay
    provider_reference VARCHAR(255), -- Provider's transaction ID

    -- User info for payment
    user_email VARCHAR(255),
    user_phone VARCHAR(50),
    user_country VARCHAR(3),

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed, cancelled, expired
    status_message TEXT,

    -- Webhook data
    webhook_received_at TIMESTAMP,
    webhook_data JSONB, -- Raw webhook payload
    webhook_verified BOOLEAN DEFAULT false,

    -- Retry & timeout handling
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    expires_at TIMESTAMP, -- Transaction expires after this time

    -- Completion
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    failed_at TIMESTAMP,
    failure_reason TEXT,

    -- User wallet credited
    user_wallet_id UUID REFERENCES wallets(id),
    wallet_credited BOOLEAN DEFAULT false,
    wallet_credited_at TIMESTAMP,

    -- Return URLs
    return_url TEXT,
    cancel_url TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',
    ip_address VARCHAR(45),
    user_agent TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for deposit_transactions
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_user ON deposit_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_status ON deposit_transactions(status);
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_provider ON deposit_transactions(provider_code);
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_provider_ref ON deposit_transactions(provider_reference);
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_expires ON deposit_transactions(expires_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_deposit_transactions_created ON deposit_transactions(created_at DESC);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_deposit_transaction_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS deposit_transaction_updated_at ON deposit_transactions;
CREATE TRIGGER deposit_transaction_updated_at
    BEFORE UPDATE ON deposit_transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_deposit_transaction_timestamp();

-- Function to expire old pending transactions (call this via cron job or scheduled task)
CREATE OR REPLACE FUNCTION expire_pending_deposits() RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    UPDATE deposit_transactions
    SET status = 'expired',
        status_message = 'Transaction expired due to timeout',
        failed_at = CURRENT_TIMESTAMP
    WHERE status = 'pending'
      AND expires_at IS NOT NULL
      AND expires_at < CURRENT_TIMESTAMP;

    GET DIAGNOSTICS expired_count = ROW_COUNT;
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- transfer, transaction, security, card, kyc, etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB, -- Additional notification data
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Exchange rates table
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    bid_price DECIMAL(20,8) DEFAULT 0,
    ask_price DECIMAL(20,8) DEFAULT 0,
    spread DECIMAL(20,8) DEFAULT 0,
    source VARCHAR(50) NOT NULL, -- coinbase, binance, etc.
    volume_24h DECIMAL(20,8) DEFAULT 0,
    change_24h DECIMAL(10,4) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_currency, to_currency)
);

-- Cards table
CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_number VARCHAR(19) NOT NULL, -- Masked: ****-****-****-1234
    card_number_full VARCHAR(255), -- Encrypted full number
    card_type VARCHAR(20) NOT NULL, -- prepaid, virtual, gift
    card_category VARCHAR(20) DEFAULT 'personal', -- personal, business
    currency VARCHAR(10) NOT NULL,
    cardholder_name VARCHAR(100),
    balance DECIMAL(20,8) DEFAULT 0,
    available_balance DECIMAL(20,8) DEFAULT 0,
    expiry_month INTEGER,
    expiry_year INTEGER,
    cvv VARCHAR(255), -- Encrypted
    pin_hash VARCHAR(255),
    status VARCHAR(20) DEFAULT 'inactive', -- active, inactive, blocked, expired
    is_virtual BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    activated_at TIMESTAMP,
    expires_at TIMESTAMP,

    -- Limits
    daily_limit DECIMAL(20,8) DEFAULT 1000,
    monthly_limit DECIMAL(20,8) DEFAULT 10000,
    single_tx_limit DECIMAL(20,8) DEFAULT 500,
    atm_daily_limit DECIMAL(20,8) DEFAULT 300,
    online_tx_limit DECIMAL(20,8) DEFAULT 2000,

    -- Current usage
    daily_spent DECIMAL(20,8) DEFAULT 0,
    monthly_spent DECIMAL(20,8) DEFAULT 0,
    atm_daily_spent DECIMAL(20,8) DEFAULT 0,

    -- Settings
    allow_atm BOOLEAN DEFAULT true,
    allow_online BOOLEAN DEFAULT true,
    allow_international BOOLEAN DEFAULT false,
    allow_contactless BOOLEAN DEFAULT true,

    -- Auto-reload
    auto_reload_enabled BOOLEAN DEFAULT false,
    auto_reload_amount DECIMAL(20,8),
    auto_reload_threshold DECIMAL(20,8),
    reload_wallet_id UUID,

    -- Physical card shipping
    shipping_address TEXT,
    shipping_status VARCHAR(20),
    tracking_number VARCHAR(100),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,

    -- External processor
    external_card_id VARCHAR(100),
    issuer_id VARCHAR(50) DEFAULT 'internal',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gift cards table
CREATE TABLE IF NOT EXISTS gift_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) UNIQUE NOT NULL,
    sender_id UUID NOT NULL REFERENCES users(id),
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(20),
    amount DECIMAL(20,8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    message TEXT,
    design VARCHAR(50) DEFAULT 'default', -- birthday, christmas, etc.
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, redeemed, expired
    redeemed_by UUID REFERENCES users(id),
    redeemed_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Card transactions table
CREATE TABLE IF NOT EXISTS card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id),
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_type VARCHAR(20) NOT NULL, -- purchase, withdrawal, load, refund
    amount DECIMAL(20,8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    original_amount DECIMAL(20,8),
    original_currency VARCHAR(10),
    exchange_rate DECIMAL(20,8),
    fee DECIMAL(20,8) DEFAULT 0,
    merchant_name VARCHAR(255),
    merchant_category VARCHAR(100),
    merchant_city VARCHAR(100),
    merchant_country VARCHAR(3),
    authorization_code VARCHAR(50),
    reference_number VARCHAR(100),
    external_transaction_id VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, declined, reversed
    decline_reason TEXT,
    is_online BOOLEAN DEFAULT false,
    is_international BOOLEAN DEFAULT false,
    is_contactless BOOLEAN DEFAULT false,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- KYC documents table
CREATE TABLE IF NOT EXISTS kyc_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- passport, id_card, driver_license, utility_bill
    identity_sub_type VARCHAR(50),       -- cni, passport, permis (for identity documents)
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    file_size BIGINT,
    mime_type VARCHAR(100),
    document_number VARCHAR(100),  -- ID/Passport/License number
    expiry_date DATE,              -- Document expiration date
    verification_status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    status VARCHAR(20) DEFAULT 'pending', -- alias for verification_status
    verified_by UUID REFERENCES users(id),
    reviewed_by VARCHAR(100),
    reviewed_at TIMESTAMP,
    rejection_reason TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User preferences table
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'dark', -- dark, light, system
    language VARCHAR(10) DEFAULT 'fr', -- fr, en, es, ar
    currency VARCHAR(10) DEFAULT 'XOF', -- USD, EUR, XOF, etc.
    timezone VARCHAR(50) DEFAULT 'Europe/Paris',
    number_format VARCHAR(10) DEFAULT 'fr', -- fr, en
    date_format VARCHAR(20) DEFAULT 'DD/MM/YYYY',
    hide_balances BOOLEAN DEFAULT false,
    analytics_enabled BOOLEAN DEFAULT true,
    auto_lock_minutes INTEGER DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notification preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    push_enabled BOOLEAN DEFAULT true,
    transfer_received BOOLEAN DEFAULT true,
    transfer_sent BOOLEAN DEFAULT true,
    card_payment BOOLEAN DEFAULT true,
    low_balance BOOLEAN DEFAULT true,
    new_login BOOLEAN DEFAULT true,
    password_change BOOLEAN DEFAULT true,
    otp_via_sms BOOLEAN DEFAULT true,
    weekly_report BOOLEAN DEFAULT false,
    newsletter BOOLEAN DEFAULT false,
    promotions BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token TEXT NOT NULL UNIQUE,
    refresh_token TEXT NOT NULL UNIQUE,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Compliance checks table
CREATE TABLE IF NOT EXISTS compliance_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_id UUID REFERENCES transactions(id),
    check_type VARCHAR(50) NOT NULL, -- aml, sanctions, pep
    status VARCHAR(20) DEFAULT 'pending', -- pending, passed, failed
    risk_score INTEGER,
    details JSONB,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trading orders table
CREATE TABLE IF NOT EXISTS trading_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pair VARCHAR(20) NOT NULL, -- BTC/USD, ETH/EUR, etc.
    order_type VARCHAR(20) NOT NULL, -- market, limit, stop_loss, take_profit
    side VARCHAR(10) NOT NULL, -- buy, sell
    amount DECIMAL(20,8) NOT NULL,
    price DECIMAL(20,8), -- NULL for market orders
    filled_amount DECIMAL(20,8) DEFAULT 0,
    average_price DECIMAL(20,8),
    status VARCHAR(20) DEFAULT 'pending', -- pending, partial, filled, cancelled, expired
    stop_price DECIMAL(20,8), -- For stop-loss orders
    expires_at TIMESTAMP,
    filled_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Exchanges/Conversions table
CREATE TABLE IF NOT EXISTS exchanges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    from_amount DECIMAL(20,8) NOT NULL,
    to_amount DECIMAL(20,8) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) DEFAULT 0,
    fee_currency VARCHAR(10),
    status VARCHAR(20) DEFAULT 'completed', -- pending, completed, failed, refunded
    from_wallet_id UUID,
    to_wallet_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- P2P trading offers table
CREATE TABLE IF NOT EXISTS p2p_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    offer_type VARCHAR(10) NOT NULL, -- buy, sell
    currency VARCHAR(10) NOT NULL,
    fiat_currency VARCHAR(10) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    min_amount DECIMAL(20,8),
    max_amount DECIMAL(20,8),
    price DECIMAL(20,8) NOT NULL,
    payment_methods JSONB,
    terms TEXT,
    status VARCHAR(20) DEFAULT 'active', -- active, paused, completed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Support tickets table
CREATE TABLE IF NOT EXISTS support_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50) NOT NULL, -- account, transaction, card, technical, other
    priority VARCHAR(20) DEFAULT 'normal', -- low, normal, high, urgent
    status VARCHAR(20) DEFAULT 'open', -- open, in_progress, resolved, closed
    assigned_to UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    sender_type VARCHAR(20) NOT NULL, -- user, agent
    message TEXT NOT NULL,
    attachments JSONB,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Payment Providers (Flutterwave, CinetPay, Paystack, etc.)
CREATE TABLE IF NOT EXISTS payment_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,           -- flutterwave, cinetpay, paystack, demo
    display_name VARCHAR(100) NOT NULL,         -- Flutterwave, CinetPay
    provider_type VARCHAR(30) NOT NULL,         -- mobile_money, card, bank_transfer, all
    api_base_url VARCHAR(255),
    api_key_encrypted VARCHAR(500),
    api_secret_encrypted VARCHAR(500),
    public_key_encrypted VARCHAR(500),
    webhook_secret_encrypted VARCHAR(500),
    is_active BOOLEAN DEFAULT false,
    is_demo_mode BOOLEAN DEFAULT true,          -- Mode démo = crédit direct sans API
    deposit_enabled BOOLEAN DEFAULT true,       -- Dépôts activés
    withdraw_enabled BOOLEAN DEFAULT true,      -- Retraits activés
    logo_url VARCHAR(255),
    supported_currencies JSONB DEFAULT '[]',    -- ["XOF", "NGN", "GHS"]
    config_json JSONB DEFAULT '{}',             -- Config spécifique au provider
    fee_percentage DECIMAL(5,2) DEFAULT 0,      -- Frais en pourcentage
    fee_fixed DECIMAL(20,2) DEFAULT 0,          -- Frais fixes
    min_transaction DECIMAL(20,2) DEFAULT 100,  -- Montant minimum par transaction
    max_transaction DECIMAL(20,2) DEFAULT 10000000,  -- Montant maximum par transaction
    daily_limit DECIMAL(20,2) DEFAULT 50000000, -- Limite journalière
    priority INT DEFAULT 1,                     -- Priorité d'affichage (1 = haute)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Provider-Country Mapping
CREATE TABLE IF NOT EXISTS provider_countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES payment_providers(id) ON DELETE CASCADE,
    country_code VARCHAR(3) NOT NULL,           -- CI, SN, NG, GH, CM
    country_name VARCHAR(100) NOT NULL,         -- Côte d'Ivoire, Sénégal
    currency VARCHAR(10) NOT NULL,              -- XOF, NGN, GHS
    is_active BOOLEAN DEFAULT true,
    priority INT DEFAULT 1,                     -- Pour ordonner les providers (1 = priorité haute)
    min_amount DECIMAL(20,2) DEFAULT 100,
    max_amount DECIMAL(20,2) DEFAULT 10000000,
    fee_percentage DECIMAL(5,2) DEFAULT 1.5,
    fee_fixed DECIMAL(20,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider_id, country_code)
);

-- Payment Transactions Log
CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID REFERENCES payment_providers(id),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id UUID NOT NULL,
    external_reference VARCHAR(255),            -- Reference du provider
    internal_reference VARCHAR(255) NOT NULL,   -- Notre reference
    amount DECIMAL(20,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    fee DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(30) DEFAULT 'pending',       -- pending, processing, completed, failed, cancelled
    payment_method VARCHAR(50),                 -- orange_money, mtn_momo, card, bank
    phone_number VARCHAR(20),
    metadata JSONB DEFAULT '{}',
    provider_response JSONB,
    error_message TEXT,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Platform Accounts (Hot Wallets & Reserves)
CREATE TABLE IF NOT EXISTS platform_accounts (
    id VARCHAR(36) PRIMARY KEY, -- JSON/Go UUID string
    currency VARCHAR(10) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- reserve, operational, fees
    name VARCHAR(100) NOT NULL,
    balance DECIMAL(20, 8) DEFAULT 0,
    min_balance DECIMAL(20, 8) DEFAULT 0,
    max_balance DECIMAL(20, 8) DEFAULT 0,
    priority INT DEFAULT 50,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Aggregator Settings table (for transfer-service routing)
CREATE TABLE IF NOT EXISTS aggregator_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_code VARCHAR(50) NOT NULL UNIQUE,
    provider_name VARCHAR(100) NOT NULL,
    payment_provider_id UUID REFERENCES payment_providers(id),
    api_base_url VARCHAR(255),
    is_enabled BOOLEAN DEFAULT true,
    is_demo_mode BOOLEAN DEFAULT true,
    supports_deposit BOOLEAN DEFAULT true,
    supports_withdrawal BOOLEAN DEFAULT true,
    min_amount DECIMAL(20, 2) DEFAULT 100,
    max_amount DECIMAL(20, 2) DEFAULT 10000000,
    fee_percentage DECIMAL(5, 4) DEFAULT 0.015,
    fee_fixed DECIMAL(20, 2) DEFAULT 0,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Aggregator Instances table (multi-key support)
CREATE TABLE IF NOT EXISTS aggregator_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregator_id UUID REFERENCES aggregator_settings(id) ON DELETE CASCADE,
    instance_name VARCHAR(100) NOT NULL,
    api_credentials JSONB DEFAULT '{}',
    vault_secret_path VARCHAR(255),
    enabled BOOLEAN DEFAULT true,
    is_paused BOOLEAN DEFAULT false,
    is_global BOOLEAN DEFAULT false,
    pause_reason TEXT,
    paused_at TIMESTAMP,
    priority INT DEFAULT 50,
    health_status VARCHAR(20) DEFAULT 'active',
    daily_limit DECIMAL(20, 8),
    monthly_limit DECIMAL(20, 8),
    daily_usage DECIMAL(20, 8) DEFAULT 0,
    monthly_usage DECIMAL(20, 8) DEFAULT 0,
    usage_reset_date DATE DEFAULT CURRENT_DATE,
    restricted_countries TEXT[],
    is_test_mode BOOLEAN DEFAULT true,
    total_transactions INT DEFAULT 0,
    total_volume DECIMAL(20, 8) DEFAULT 0,
    last_used_at TIMESTAMP,
    last_error TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID
);

-- Aggregator Instance Wallets table (multi-wallet per instance)
CREATE TABLE IF NOT EXISTS aggregator_instance_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES aggregator_instances(id) ON DELETE CASCADE,
    hot_wallet_id VARCHAR(36) NOT NULL,  -- VARCHAR to match platform_accounts.id type
    currency VARCHAR(10) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    priority INT DEFAULT 50,
    min_balance DECIMAL(20, 8) DEFAULT 0,
    max_balance DECIMAL(20, 8),
    enabled BOOLEAN DEFAULT true,
    total_deposits DECIMAL(20, 8) DEFAULT 0,
    total_withdrawals DECIMAL(20, 8) DEFAULT 0,
    transaction_count INT DEFAULT 0,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(instance_id, hot_wallet_id, currency)
);

-- Association Tables
CREATE TABLE IF NOT EXISTS associations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'tontine',
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency VARCHAR(10) DEFAULT 'XOF',
    contribution_amount DECIMAL(20,2) DEFAULT 0,
    contribution_frequency VARCHAR(20) DEFAULT 'monthly',
    treasury_balance DECIMAL(20,2) DEFAULT 0,
    total_members INTEGER DEFAULT 1,
    rules JSONB DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS association_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name VARCHAR(255),
    user_email VARCHAR(255),
    role VARCHAR(50) DEFAULT 'member',
    contributions_paid DECIMAL(20,2) DEFAULT 0,
    contributions_count INTEGER DEFAULT 0,
    loans_received DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(association_id, user_id)
);

CREATE TABLE IF NOT EXISTS association_treasury (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    amount DECIMAL(20,2) NOT NULL,
    from_member_id UUID REFERENCES association_members(id),
    to_member_id UUID REFERENCES association_members(id),
    description TEXT,
    status VARCHAR(20) DEFAULT 'completed',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS association_loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
    borrower_id UUID NOT NULL REFERENCES association_members(id),
    amount DECIMAL(20,2) NOT NULL,
    interest_rate DECIMAL(5,2) DEFAULT 0,
    duration INTEGER DEFAULT 3,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    repayments JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'pending',
    approved_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS association_meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    location VARCHAR(500),
    attendees JSONB DEFAULT '[]',
    minutes TEXT,
    status VARCHAR(20) DEFAULT 'scheduled',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admin Tables
CREATE TABLE IF NOT EXISTS admin_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB DEFAULT '[]',
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role_id UUID REFERENCES admin_roles(id),
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES admins(id)
);

CREATE TABLE IF NOT EXISTS admin_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID REFERENCES admins(id),
    admin_email VARCHAR(255),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id VARCHAR(100),
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 2. INDEX CREATION
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = false;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_wallet ON transactions(from_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_wallet ON transactions(to_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_backup_codes_user_id ON backup_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_token ON verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_cards_user_id ON cards(user_id);
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status);
CREATE INDEX IF NOT EXISTS idx_gift_cards_code ON gift_cards(code);
CREATE INDEX IF NOT EXISTS idx_gift_cards_sender_id ON gift_cards(sender_id);
CREATE INDEX IF NOT EXISTS idx_card_transactions_card_id ON card_transactions(card_id);
CREATE INDEX IF NOT EXISTS idx_card_transactions_user_id ON card_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_card_transactions_created_at ON card_transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_trading_orders_user_id ON trading_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_trading_orders_status ON trading_orders(status);
CREATE INDEX IF NOT EXISTS idx_trading_orders_pair ON trading_orders(pair);
CREATE INDEX IF NOT EXISTS idx_exchanges_user_id ON exchanges(user_id);
CREATE INDEX IF NOT EXISTS idx_exchanges_status ON exchanges(status);
CREATE INDEX IF NOT EXISTS idx_p2p_offers_user_id ON p2p_offers(user_id);
CREATE INDEX IF NOT EXISTS idx_p2p_offers_status ON p2p_offers(status);
CREATE INDEX IF NOT EXISTS idx_support_tickets_user_id ON support_tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_support_tickets_status ON support_tickets(status);
CREATE INDEX IF NOT EXISTS idx_support_tickets_assigned_to ON support_tickets(assigned_to);
CREATE INDEX IF NOT EXISTS idx_chat_messages_ticket_id ON chat_messages(ticket_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_sender_id ON chat_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_payment_providers_active ON payment_providers(is_active);
CREATE INDEX IF NOT EXISTS idx_payment_providers_name ON payment_providers(name);
CREATE INDEX IF NOT EXISTS idx_provider_countries_country ON provider_countries(country_code);
CREATE INDEX IF NOT EXISTS idx_provider_countries_provider ON provider_countries(provider_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_user ON payment_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_ref ON payment_transactions(external_reference);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_prefs_user_id ON notification_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_platform_accounts_type ON platform_accounts(account_type);
CREATE INDEX IF NOT EXISTS idx_platform_accounts_currency ON platform_accounts(currency);
CREATE INDEX IF NOT EXISTS idx_aggregator_settings_code ON aggregator_settings(provider_code);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_aggregator ON aggregator_instances(aggregator_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_enabled ON aggregator_instances(enabled);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_paused ON aggregator_instances(is_paused);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_instance ON aggregator_instance_wallets(instance_id);
CREATE INDEX IF NOT EXISTS idx_associations_creator ON associations(creator_id);
CREATE INDEX IF NOT EXISTS idx_associations_status ON associations(status);
CREATE INDEX IF NOT EXISTS idx_association_members_assoc ON association_members(association_id);
CREATE INDEX IF NOT EXISTS idx_association_members_user ON association_members(user_id);
CREATE INDEX IF NOT EXISTS idx_association_treasury_assoc ON association_treasury(association_id);
CREATE INDEX IF NOT EXISTS idx_association_loans_assoc ON association_loans(association_id);
CREATE INDEX IF NOT EXISTS idx_association_loans_borrower ON association_loans(borrower_id);
CREATE INDEX IF NOT EXISTS idx_association_meetings_assoc ON association_meetings(association_id);
CREATE INDEX IF NOT EXISTS idx_admins_email ON admins(email);
CREATE INDEX IF NOT EXISTS idx_admins_role ON admins(role_id);
CREATE INDEX IF NOT EXISTS idx_admin_audit_admin ON admin_audit_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_audit_action ON admin_audit_logs(action);


-- =====================================================
-- 3. VIEWS DEFINITION
-- =====================================================

-- View for transfer-service to query instances with details
-- Column order MUST match Go code's scanInstanceWithDetails function exactly!
CREATE OR REPLACE VIEW aggregator_instances_with_details AS
SELECT
    -- Columns 1-9: Instance base info
    ai.id,
    ai.instance_name,
    ai.enabled,
    ai.priority,
    COALESCE(ai.is_test_mode, true) AS is_test_mode,
    COALESCE(ai.is_paused, false) AS is_paused,
    COALESCE(ai.is_global, true) AS is_global,
    ai.pause_reason,
    ai.paused_at,
    -- Column 10: restricted_countries (array)
    ai.restricted_countries,
    -- Columns 11-17: Limits and usage
    ai.daily_limit,
    ai.monthly_limit,
    COALESCE(ai.daily_usage, 0) AS daily_usage,
    COALESCE(ai.monthly_usage, 0) AS monthly_usage,
    COALESCE(ai.total_transactions, 0) AS total_transactions,
    COALESCE(ai.total_volume, 0) AS total_volume,
    ai.last_used_at,
    -- Columns 18-19: Timestamps
    ai.created_at,
    ai.updated_at,
    -- Column 20: aggregator_id
    ai.aggregator_id,
    -- Columns 21-24: Aggregator info
    agg.provider_code,
    agg.provider_name,
    COALESCE(pp.logo_url, '/icons/aggregators/' || agg.provider_code || '.svg') AS provider_logo,
    agg.is_enabled AS aggregator_enabled,
    -- Columns 25-27: Hot wallet info
    aiw.hot_wallet_id,
    'operations'::text AS account_type,
    COALESCE(aiw.currency, 'XOF') AS hot_wallet_currency,
    -- Column 28: Hot wallet balance
    COALESCE(pa.balance, 0) AS hot_wallet_balance,
    -- Columns 29-30: Wallet limits
    COALESCE(aiw.min_balance, 0) AS min_balance,
    aiw.max_balance,
    -- Column 31: Availability status
    CASE
        WHEN NOT ai.enabled THEN 'instance_disabled'
        WHEN COALESCE(ai.is_paused, false) THEN 'paused'
        WHEN NOT agg.is_enabled THEN 'aggregator_disabled'
        WHEN aiw.id IS NULL THEN 'no_wallet'
        WHEN pa.id IS NULL THEN 'wallet_not_found'
        WHEN pa.balance < COALESCE(aiw.min_balance, 0) THEN 'insufficient_balance'
        ELSE 'available'
    END AS availability_status
FROM aggregator_instances ai
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
LEFT JOIN payment_providers pp ON agg.payment_provider_id = pp.id
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id AND aiw.is_primary = true AND aiw.enabled = true
LEFT JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id AND pa.is_active = true;


-- =====================================================
-- 4. DATA SEEDING
-- =====================================================

-- Supported currencies
INSERT INTO exchange_rates (from_currency, to_currency, rate, bid_price, ask_price, spread, source, volume_24h, change_24h, valid_until) VALUES
('USD', 'EUR', 0.85, 0.8498, 0.8502, 0.0004, 'system', 5000000.00, 0.1, NOW() + INTERVAL '1 day'),
('EUR', 'USD', 1.18, 1.1798, 1.1802, 0.0004, 'system', 5000000.00, -0.1, NOW() + INTERVAL '1 day'),
('BTC', 'USD', 45000.00, 44990.00, 45010.00, 20.00, 'binance', 1234567.89, 2.5, NOW() + INTERVAL '1 day'),
('ETH', 'USD', 3000.00, 2998.00, 3002.00, 4.00, 'binance', 987654.32, 1.8, NOW() + INTERVAL '1 day'),
('USD', 'BTC', 0.000022, 0.0000219, 0.0000221, 0.0000002, 'binance', 1234567.89, -2.5, NOW() + INTERVAL '1 day'),
('USD', 'ETH', 0.000333, 0.000332, 0.000334, 0.000002, 'binance', 987654.32, -1.8, NOW() + INTERVAL '1 day'),
('SOL', 'USD', 100.00, 99.95, 100.05, 0.10, 'binance', 543210.11, 3.2, NOW() + INTERVAL '1 day'),
('XRP', 'USD', 0.63, 0.6298, 0.6302, 0.0004, 'binance', 432109.87, 1.7, NOW() + INTERVAL '1 day')
ON CONFLICT (from_currency, to_currency) DO UPDATE SET
    rate = EXCLUDED.rate,
    bid_price = EXCLUDED.bid_price,
    ask_price = EXCLUDED.ask_price,
    spread = EXCLUDED.spread,
    volume_24h = EXCLUDED.volume_24h,
    change_24h = EXCLUDED.change_24h,
    last_updated = NOW();

-- Insert default payment providers
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
-- Demo Mode
('demo', 'Mode Démo', 'all', NULL, true, true, '/icons/demo.svg',
 '["XOF", "XAF", "NGN", "GHS", "KES", "USD", "EUR"]'::jsonb,
 '{"description": "Mode test - crédite directement le compte sans paiement réel"}'::jsonb),
-- Flutterwave
('flutterwave', 'Flutterwave', 'all', 'https://api.flutterwave.com/v3', true, true, '/icons/flutterwave.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "USD", "EUR", "GBP"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true, "webhook_path": "/webhooks/flutterwave"}'::jsonb),
-- CinetPay
('cinetpay', 'CinetPay', 'mobile_money', 'https://api-checkout.cinetpay.com/v2', true, true, '/icons/cinetpay.svg',
 '["XOF", "XAF", "GNF"]'::jsonb,
 '{"supports_mobile_money": true, "operators": ["orange_money", "mtn_momo", "moov_money", "wave"], "webhook_path": "/webhooks/cinetpay"}'::jsonb),
-- Paystack
('paystack', 'Paystack', 'all', 'https://api.paystack.co', true, true, '/icons/paystack.svg',
 '["NGN", "GHS", "ZAR", "KES"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true, "webhook_path": "/webhooks/paystack"}'::jsonb),
-- Orange Money
('orange_money', 'Orange Money', 'mobile_money', 'https://api.orange.com/orange-money-webpay', true, true, '/icons/orange.svg',
 '["XOF", "XAF"]'::jsonb,
 '{"countries": ["CI", "SN", "ML", "BF", "CM", "GN"], "webhook_path": "/webhooks/orange"}'::jsonb),
-- MTN MoMo
('mtn_momo', 'MTN Mobile Money', 'mobile_money', 'https://sandbox.momodeveloper.mtn.com', true, true, '/icons/mtn.svg',
 '["XOF", "XAF", "GHS", "UGX", "RWF"]'::jsonb,
 '{"countries": ["CI", "CM", "GH", "UG", "RW", "BJ"], "webhook_path": "/webhooks/mtn"}'::jsonb),
-- Wave
('wave', 'Wave', 'mobile_money', 'https://api.wave.com/v1', true, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["SN", "CI"], "webhook_path": "/webhooks/wave"}'::jsonb),
-- Stripe
('stripe', 'Stripe', 'card', 'https://api.stripe.com/v1', true, true, '/icons/stripe.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD"]'::jsonb,
 '{"supports_card": true, "webhook_path": "/webhooks/stripe"}'::jsonb),
-- PayPal
('paypal', 'PayPal', 'wallet', 'https://api-m.paypal.com', true, true, '/icons/paypal.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD", "BRL", "MXN", "ILS"]'::jsonb,
 '{"supports_wallet": true, "supports_card": true, "webhook_path": "/webhooks/paypal"}'::jsonb),
-- Lygos
('lygos', 'Lygos', 'mobile_money', 'https://api.lygosapp.com/v1', true, true, '/icons/lygos.svg',
 '["XOF", "XAF", "GNF", "CDF", "RWF", "KES", "UGX", "NGN"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "SN", "BF", "ML", "GN", "CM", "CD", "RW", "KE", "UG", "NG", "TG", "BJ", "NE"]}'::jsonb),
-- YellowCard
('yellowcard', 'YellowCard', 'crypto_ramp', 'https://api.yellowcard.io/v1', true, true, '/icons/yellowcard.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "UGX", "TZS", "RWF", "BWP"]'::jsonb,
 '{"supports_crypto_ramp": true, "crypto_supported": ["BTC", "ETH", "USDT", "USDC"]}'::jsonb),
-- FedaPay
('fedapay', 'FedaPay', 'mobile_money', 'https://api.fedapay.com/v1', true, true, '/icons/fedapay.svg',
 '["XOF"]'::jsonb, '{"supports_mobile_money": true, "countries": ["BJ", "TG", "NE", "CI"]}'::jsonb),
-- Moov Money
('moov_money', 'Moov Money', 'mobile_money', 'https://api.moov-africa.bj/v1', true, true, '/icons/moov.svg',
 '["XOF", "XAF"]'::jsonb, '{"supports_mobile_money": true, "countries": ["CI", "BJ", "TG", "BF", "NE", "CM", "GA", "CG"]}'::jsonb)

ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    api_base_url = EXCLUDED.api_base_url,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true,
    updated_at = NOW();

-- Insert default country mappings (simplified for brevity, ensuring distinct selects)
-- DEMO
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 0 FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage)
SELECT id, 'SN', 'Sénégal', 'XOF', 1, 0 FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
-- (You would include all other mappings here, keeping the long list if needed, or trusting the dynamic seeding if it's external.
-- For now, I'll include the critical ones for testing).

-- CinetPay
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;

-- Seed Platform Accounts
INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance) VALUES
-- FCFA
(gen_random_uuid()::text, 'FCFA', 'reserve', 'Réserve FCFA Principal', 'Réserve principale en FCFA', 100, 1000000000),
(gen_random_uuid()::text, 'FCFA', 'fees', 'Frais collectés FCFA', 'Frais de transaction collectés', 100, 0),
(gen_random_uuid()::text, 'FCFA', 'operations', 'Opérations FCFA', 'Compte opérationnel pour retraits/dépôts', 80, 100000000),
-- XOF
(gen_random_uuid()::text, 'XOF', 'reserve', 'Réserve XOF', 'Réserve principale en XOF', 100, 1000000000),
(gen_random_uuid()::text, 'XOF', 'fees', 'Frais collectés XOF', 'Frais de transaction en XOF', 100, 0),
(gen_random_uuid()::text, 'XOF', 'operations', 'Opérations XOF', 'Compte opérationnel pour retraits/dépôts XOF', 80, 100000000),
-- USD
(gen_random_uuid()::text, 'USD', 'reserve', 'Réserve USD', 'Réserve principale en USD', 100, 1000000),
(gen_random_uuid()::text, 'USD', 'fees', 'Frais collectés USD', 'Frais de transaction en USD', 100, 0),
(gen_random_uuid()::text, 'USD', 'operations', 'Opérations USD', 'Compte opérationnel USD', 80, 50000),
-- EUR
(gen_random_uuid()::text, 'EUR', 'reserve', 'Réserve EUR', 'Réserve principale en EUR', 100, 1000000),
(gen_random_uuid()::text, 'EUR', 'fees', 'Frais collectés EUR', 'Frais de transaction en EUR', 100, 0),
(gen_random_uuid()::text, 'EUR', 'operations', 'Opérations EUR', 'Compte opérationnel EUR', 80, 50000),
-- GHS
(gen_random_uuid()::text, 'GHS', 'reserve', 'Réserve GHS', 'Réserve principale en GHS', 100, 1000000000),
(gen_random_uuid()::text, 'GHS', 'operations', 'Opérations GHS', 'Compte opérationnel pour retraits/dépôts GHS', 80, 10000000),
-- NGN
(gen_random_uuid()::text, 'NGN', 'reserve', 'Réserve NGN', 'Réserve principale en NGN', 100, 1000000000),
(gen_random_uuid()::text, 'NGN', 'operations', 'Opérations NGN', 'Compte opérationnel pour retraits/dépôts NGN', 80, 500000000)
ON CONFLICT DO NOTHING;

-- Admin Roles
INSERT INTO admin_roles (name, description, permissions, is_system) VALUES
('super_admin', 'Super Administrateur - Tous les droits', '["users.view", "users.create", "users.update", "users.block", "users.delete", "kyc.view", "kyc.approve", "kyc.reject", "transactions.view", "transactions.block", "transactions.refund", "cards.view", "cards.freeze", "cards.block", "wallets.view", "wallets.freeze", "wallets.adjust", "system.view", "system.logs", "system.settings", "admins.view", "admins.create", "admins.update", "admins.delete", "admins.roles", "analytics.view", "analytics.export"]'::jsonb, true),
('admin', 'Administrateur', '["users.view", "users.update", "users.block", "kyc.view", "kyc.approve", "kyc.reject", "transactions.view", "cards.view", "cards.freeze", "wallets.view", "wallets.freeze", "system.view", "system.logs", "analytics.view"]'::jsonb, true),
('support', 'Support Client', '["users.view", "kyc.view", "transactions.view", "cards.view", "wallets.view"]'::jsonb, true)
ON CONFLICT (name) DO NOTHING;

-- Super Admin
INSERT INTO admins (email, password_hash, first_name, last_name, role_id, is_active) VALUES
('admin@crypto-bank.com', '$2a$10$rQZ7nAQ.L3dNLqvHQQQQwekHqg4BzKqvBFHqKqH9.K3XGMGGe0Gey', 'Super', 'Admin', (SELECT id FROM admin_roles WHERE name = 'super_admin'), true)
ON CONFLICT (email) DO NOTHING;

-- Seed aggregator_settings from payment_providers
INSERT INTO aggregator_settings (provider_code, provider_name, payment_provider_id, api_base_url, is_enabled, is_demo_mode, supports_deposit, supports_withdrawal, config)
SELECT name, display_name, id, api_base_url, is_active, is_demo_mode, deposit_enabled, withdraw_enabled, config_json
FROM payment_providers
WHERE name IN ('demo', 'flutterwave', 'cinetpay', 'paystack', 'orange_money', 'mtn_momo', 'wave', 'stripe', 'paypal', 'lygos', 'yellowcard', 'fedapay', 'moov_money')
ON CONFLICT (provider_code) DO UPDATE SET is_enabled = EXCLUDED.is_enabled, is_demo_mode = EXCLUDED.is_demo_mode;

-- Seed default instances for each aggregator
DO $$
DECLARE v_agg RECORD; v_inst_id UUID; v_wallet RECORD;
BEGIN
    FOR v_agg IN SELECT id, provider_code, provider_name FROM aggregator_settings LOOP
        SELECT id INTO v_inst_id FROM aggregator_instances WHERE aggregator_id = v_agg.id AND instance_name = 'Instance Principale';
        IF v_inst_id IS NULL THEN
            INSERT INTO aggregator_instances (aggregator_id, instance_name, api_credentials, vault_secret_path, enabled, is_global, priority, health_status, is_test_mode)
            VALUES (
                v_agg.id,
                'Instance Principale',
                CASE v_agg.provider_code
                    WHEN 'paypal' THEN jsonb_build_object(
                        'client_id', 'REPLACE_ME',
                        'client_secret', 'REPLACE_ME',
                        'mode', 'sandbox'
                    )
                    WHEN 'stripe' THEN jsonb_build_object(
                        'api_key', 'REPLACE_ME',
                        'public_key', 'REPLACE_ME',
                        'webhook_secret', 'REPLACE_ME'
                    )
                    WHEN 'flutterwave' THEN jsonb_build_object(
                        'public_key', 'REPLACE_ME',
                        'secret_key', 'REPLACE_ME',
                        'encryption_key', 'REPLACE_ME'
                    )
                    WHEN 'cinetpay' THEN jsonb_build_object(
                        'api_key', 'REPLACE_ME',
                        'site_id', 'REPLACE_ME',
                        'secret_key', 'REPLACE_ME'
                    )
                    WHEN 'paystack' THEN jsonb_build_object(
                        'public_key', 'REPLACE_ME',
                        'secret_key', 'REPLACE_ME'
                    )
                    WHEN 'orange_money' THEN jsonb_build_object(
                        'client_id', 'REPLACE_ME',
                        'client_secret', 'REPLACE_ME',
                        'merchant_key', 'REPLACE_ME'
                    )
                    WHEN 'mtn_momo' THEN jsonb_build_object(
                        'api_user', 'REPLACE_ME',
                        'api_key', 'REPLACE_ME',
                        'subscription_key', 'REPLACE_ME',
                        'environment', 'sandbox'
                    )
                    WHEN 'wave' THEN jsonb_build_object(
                        'api_key', 'REPLACE_ME'
                    )
                    WHEN 'lygos' THEN jsonb_build_object(
                        'api_key', 'REPLACE_ME',
                        'secret_key', 'REPLACE_ME'
                    )
                    ELSE '{}'::jsonb
                END,
                'secret/aggregators/' || v_agg.provider_code || '/default',
                true,
                true,
                100,
                'active',
                true
            )
            RETURNING id INTO v_inst_id;
        END IF;
        IF v_inst_id IS NOT NULL THEN
            FOR v_wallet IN SELECT id, currency FROM platform_accounts WHERE account_type = 'operations' AND is_active = true LOOP
                INSERT INTO aggregator_instance_wallets (instance_id, hot_wallet_id, currency, is_primary, priority, enabled)
                VALUES (v_inst_id, v_wallet.id, v_wallet.currency, v_wallet.currency IN ('XOF', 'NGN', 'USD'), CASE WHEN v_wallet.currency IN ('XOF', 'NGN') THEN 100 ELSE 50 END, true)
                ON CONFLICT (instance_id, hot_wallet_id, currency) DO NOTHING;
            END LOOP;
        END IF;
    END LOOP;
END $$;
