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
    kyc_status VARCHAR(20) DEFAULT 'pending', -- pending, verified, rejected
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

-- Exchange rates table
CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    source VARCHAR(50) NOT NULL, -- coinbase, binance, etc.
    valid_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_currency, to_currency, source)
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
    verification_status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    verified_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
INSERT INTO exchange_rates (from_currency, to_currency, rate, source, valid_until) VALUES
('USD', 'EUR', 0.85, 'system', NOW() + INTERVAL '1 day'),
('EUR', 'USD', 1.18, 'system', NOW() + INTERVAL '1 day'),
('BTC', 'USD', 45000.00, 'system', NOW() + INTERVAL '1 day'),
('ETH', 'USD', 3000.00, 'system', NOW() + INTERVAL '1 day'),
('USD', 'BTC', 0.000022, 'system', NOW() + INTERVAL '1 day'),
('USD', 'ETH', 0.000333, 'system', NOW() + INTERVAL '1 day');