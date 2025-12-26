-- Crypto Bank Database Schema

-- Create admin database
CREATE DATABASE crypto_bank_admin;

-- Users table
CREATE TABLE users (
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
CREATE TABLE backup_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL, -- Hashed backup code
    used BOOLEAN DEFAULT false,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verification tokens for email/phone/password reset
CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL, -- email_verification, phone_verification, password_reset
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Wallets table
CREATE TABLE wallets (
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
CREATE TABLE transactions (
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

-- Notifications table
CREATE TABLE notifications (
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

-- Index for faster notification queries
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = false;

-- Exchange rates table
CREATE TABLE exchange_rates (
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

-- Cards table (enhanced for card-service)
CREATE TABLE cards (
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
CREATE TABLE gift_cards (
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
CREATE TABLE card_transactions (
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
CREATE TABLE kyc_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- passport, id_card, driver_license, utility_bill
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    file_size BIGINT,
    mime_type VARCHAR(100),
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
CREATE TABLE user_preferences (
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
CREATE TABLE notification_preferences (
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

-- Indexes for preferences
CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX idx_notification_prefs_user_id ON notification_preferences(user_id);

-- Audit logs table
CREATE TABLE audit_logs (
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
CREATE TABLE user_sessions (
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

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- email, sms, push
    subject VARCHAR(255),
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, failed
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Compliance checks table
CREATE TABLE compliance_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_id UUID REFERENCES transactions(id),
    check_type VARCHAR(50) NOT NULL, -- aml, sanctions, pep
    status VARCHAR(20) DEFAULT 'pending', -- pending, passed, failed
    risk_score INTEGER,
    details JSONB,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_from_wallet ON transactions(from_wallet_id);
CREATE INDEX idx_transactions_to_wallet ON transactions(to_wallet_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_backup_codes_user_id ON backup_codes(user_id);
CREATE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE INDEX idx_verification_tokens_token ON verification_tokens(token);
CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_gift_cards_code ON gift_cards(code);
CREATE INDEX idx_gift_cards_sender_id ON gift_cards(sender_id);
CREATE INDEX idx_card_transactions_card_id ON card_transactions(card_id);
CREATE INDEX idx_card_transactions_user_id ON card_transactions(user_id);
CREATE INDEX idx_card_transactions_created_at ON card_transactions(created_at);

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

-- Trading orders table
CREATE TABLE trading_orders (
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
CREATE TABLE exchanges (
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
CREATE TABLE p2p_offers (
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

-- Indexes for trading tables
CREATE INDEX idx_trading_orders_user_id ON trading_orders(user_id);
CREATE INDEX idx_trading_orders_status ON trading_orders(status);
CREATE INDEX idx_trading_orders_pair ON trading_orders(pair);
CREATE INDEX idx_exchanges_user_id ON exchanges(user_id);
CREATE INDEX idx_exchanges_status ON exchanges(status);
CREATE INDEX idx_p2p_offers_user_id ON p2p_offers(user_id);
CREATE INDEX idx_p2p_offers_status ON p2p_offers(status);

-- Support tickets table
CREATE TABLE support_tickets (
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
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    sender_type VARCHAR(20) NOT NULL, -- user, agent
    message TEXT NOT NULL,
    attachments JSONB,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for support tables
CREATE INDEX idx_support_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);
CREATE INDEX idx_support_tickets_assigned_to ON support_tickets(assigned_to);
CREATE INDEX idx_chat_messages_ticket_id ON chat_messages(ticket_id);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);

-- =====================================================
-- PAYMENT PROVIDERS / AGGREGATORS
-- =====================================================

-- Payment Providers (Flutterwave, CinetPay, Paystack, etc.)
CREATE TABLE payment_providers (
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
    logo_url VARCHAR(255),
    supported_currencies JSONB DEFAULT '[]',    -- ["XOF", "NGN", "GHS"]
    config_json JSONB DEFAULT '{}',             -- Config spécifique au provider
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Provider-Country Mapping
CREATE TABLE provider_countries (
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
CREATE TABLE payment_transactions (
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

-- Indexes for payment tables
CREATE INDEX idx_payment_providers_active ON payment_providers(is_active);
CREATE INDEX idx_payment_providers_name ON payment_providers(name);
CREATE INDEX idx_provider_countries_country ON provider_countries(country_code);
CREATE INDEX idx_provider_countries_provider ON provider_countries(provider_id);
CREATE INDEX idx_payment_transactions_user ON payment_transactions(user_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX idx_payment_transactions_ref ON payment_transactions(external_reference);

-- Insert default payment providers
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
-- Demo Mode (always available for testing)
('demo', 'Mode Démo', 'all', NULL, true, true, '/icons/demo.svg', 
 '["XOF", "XAF", "NGN", "GHS", "KES", "USD", "EUR"]'::jsonb,
 '{"description": "Mode test - crédite directement le compte sans paiement réel"}'::jsonb),

-- Flutterwave
('flutterwave', 'Flutterwave', 'all', 'https://api.flutterwave.com/v3', false, true, '/icons/flutterwave.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "USD", "EUR", "GBP"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true, "webhook_path": "/webhooks/flutterwave"}'::jsonb),

-- CinetPay
('cinetpay', 'CinetPay', 'mobile_money', 'https://api-checkout.cinetpay.com/v2', false, true, '/icons/cinetpay.svg',
 '["XOF", "XAF", "GNF"]'::jsonb,
 '{"supports_mobile_money": true, "operators": ["orange_money", "mtn_momo", "moov_money", "wave"], "webhook_path": "/webhooks/cinetpay"}'::jsonb),

-- Paystack
('paystack', 'Paystack', 'all', 'https://api.paystack.co', false, true, '/icons/paystack.svg',
 '["NGN", "GHS", "ZAR", "KES"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true, "webhook_path": "/webhooks/paystack"}'::jsonb),

-- Orange Money Direct API
('orange_money', 'Orange Money', 'mobile_money', 'https://api.orange.com/orange-money-webpay', false, true, '/icons/orange.svg',
 '["XOF", "XAF"]'::jsonb,
 '{"countries": ["CI", "SN", "ML", "BF", "CM", "GN"], "webhook_path": "/webhooks/orange"}'::jsonb),

-- MTN MoMo
('mtn_momo', 'MTN Mobile Money', 'mobile_money', 'https://sandbox.momodeveloper.mtn.com', false, true, '/icons/mtn.svg',
 '["XOF", "XAF", "GHS", "UGX", "RWF"]'::jsonb,
 '{"countries": ["CI", "CM", "GH", "UG", "RW", "BJ"], "webhook_path": "/webhooks/mtn"}'::jsonb),

-- Wave (Sénégal, CI)
('wave', 'Wave', 'mobile_money', 'https://api.wave.com/v1', false, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["SN", "CI"], "webhook_path": "/webhooks/wave"}'::jsonb),

-- Stripe (International)
('stripe', 'Stripe', 'card', 'https://api.stripe.com/v1', false, true, '/icons/stripe.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD"]'::jsonb,
 '{"supports_card": true, "webhook_path": "/webhooks/stripe"}'::jsonb)

ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    api_base_url = EXCLUDED.api_base_url,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    updated_at = NOW();

-- Insert default country mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) VALUES
-- Demo mode for all countries
((SELECT id FROM payment_providers WHERE name = 'demo'), 'CI', 'Côte d''Ivoire', 'XOF', 1, 0),
((SELECT id FROM payment_providers WHERE name = 'demo'), 'SN', 'Sénégal', 'XOF', 1, 0),
((SELECT id FROM payment_providers WHERE name = 'demo'), 'NG', 'Nigeria', 'NGN', 1, 0),
((SELECT id FROM payment_providers WHERE name = 'demo'), 'GH', 'Ghana', 'GHS', 1, 0),
((SELECT id FROM payment_providers WHERE name = 'demo'), 'CM', 'Cameroun', 'XAF', 1, 0),

-- CinetPay
((SELECT id FROM payment_providers WHERE name = 'cinetpay'), 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'cinetpay'), 'SN', 'Sénégal', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'cinetpay'), 'CM', 'Cameroun', 'XAF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'cinetpay'), 'BF', 'Burkina Faso', 'XOF', 2, 1.5),

-- Flutterwave
((SELECT id FROM payment_providers WHERE name = 'flutterwave'), 'NG', 'Nigeria', 'NGN', 2, 1.4),
((SELECT id FROM payment_providers WHERE name = 'flutterwave'), 'GH', 'Ghana', 'GHS', 2, 1.4),
((SELECT id FROM payment_providers WHERE name = 'flutterwave'), 'KE', 'Kenya', 'KES', 2, 1.4),
((SELECT id FROM payment_providers WHERE name = 'flutterwave'), 'CI', 'Côte d''Ivoire', 'XOF', 3, 1.4),

-- Paystack
((SELECT id FROM payment_providers WHERE name = 'paystack'), 'NG', 'Nigeria', 'NGN', 3, 1.5),
((SELECT id FROM payment_providers WHERE name = 'paystack'), 'GH', 'Ghana', 'GHS', 3, 1.5),

-- Orange Money
((SELECT id FROM payment_providers WHERE name = 'orange_money'), 'CI', 'Côte d''Ivoire', 'XOF', 4, 1.0),
((SELECT id FROM payment_providers WHERE name = 'orange_money'), 'SN', 'Sénégal', 'XOF', 4, 1.0),

-- MTN MoMo
((SELECT id FROM payment_providers WHERE name = 'mtn_momo'), 'CI', 'Côte d''Ivoire', 'XOF', 5, 1.0),
((SELECT id FROM payment_providers WHERE name = 'mtn_momo'), 'CM', 'Cameroun', 'XAF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'mtn_momo'), 'GH', 'Ghana', 'GHS', 4, 1.0),

-- Wave
((SELECT id FROM payment_providers WHERE name = 'wave'), 'SN', 'Sénégal', 'XOF', 2, 1.0),
((SELECT id FROM payment_providers WHERE name = 'wave'), 'CI', 'Côte d''Ivoire', 'XOF', 5, 1.0)

ON CONFLICT (provider_id, country_code) DO UPDATE SET
    priority = EXCLUDED.priority,
    fee_percentage = EXCLUDED.fee_percentage;

-- =====================================================
-- ADMIN USERS & PERMISSIONS
-- =====================================================

-- Admin Roles
CREATE TABLE admin_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB DEFAULT '[]',
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admin Users
CREATE TABLE admins (
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

-- Admin Audit Logs
CREATE TABLE admin_audit_logs (
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

-- Indexes
CREATE INDEX idx_admins_email ON admins(email);
CREATE INDEX idx_admins_role ON admins(role_id);
CREATE INDEX idx_admin_audit_admin ON admin_audit_logs(admin_id);
CREATE INDEX idx_admin_audit_action ON admin_audit_logs(action);

-- Insert default admin roles
INSERT INTO admin_roles (name, description, permissions, is_system) VALUES
('super_admin', 'Super Administrateur - Tous les droits', 
 '["users.view", "users.create", "users.update", "users.block", "users.delete", "kyc.view", "kyc.approve", "kyc.reject", "transactions.view", "transactions.block", "transactions.refund", "cards.view", "cards.freeze", "cards.block", "wallets.view", "wallets.freeze", "wallets.adjust", "system.view", "system.logs", "system.settings", "admins.view", "admins.create", "admins.update", "admins.delete", "admins.roles", "analytics.view", "analytics.export"]'::jsonb, 
 true),
('admin', 'Administrateur', 
 '["users.view", "users.update", "users.block", "kyc.view", "kyc.approve", "kyc.reject", "transactions.view", "cards.view", "cards.freeze", "wallets.view", "wallets.freeze", "system.view", "system.logs", "analytics.view"]'::jsonb, 
 true),
('support', 'Support Client', 
 '["users.view", "kyc.view", "transactions.view", "cards.view", "wallets.view"]'::jsonb, 
 true),
('viewer', 'Lecteur seul', 
 '["users.view", "transactions.view", "analytics.view"]'::jsonb, 
 true)
ON CONFLICT (name) DO NOTHING;

-- Insert default super admin
-- Password: Admin123! (bcrypt hash)
INSERT INTO admins (email, password_hash, first_name, last_name, role_id, is_active) VALUES
('admin@crypto-bank.com', 
 '$2a$10$rQZ7nAQ.L3dNLqvHQQQQwekHqg4BzKqvBFHqKqH9.K3XGMGGe0Gey',
 'Super', 
 'Admin',
 (SELECT id FROM admin_roles WHERE name = 'super_admin'),
 true)
ON CONFLICT (email) DO NOTHING;

-- =====================================================
-- PIN PROGRESSIVE LOCK SYSTEM
-- Add columns for progressive PIN lockout:
-- - pin_permanently_locked: TRUE when locked by admin only
-- - pin_temp_lock_count: Number of times temp locked (0 or 1)
-- =====================================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin_permanently_locked BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin_temp_lock_count INTEGER DEFAULT 0;