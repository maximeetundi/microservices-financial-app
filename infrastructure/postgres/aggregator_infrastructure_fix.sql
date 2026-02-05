-- =====================================================
-- AGGREGATOR INFRASTRUCTURE COMPLETE FIX
-- Adds missing tables, views, aggregators for French-speaking Africa
-- Run this in the MAIN database (crypto_bank / postgres)
-- =====================================================

-- =====================================================
-- 1. ADD MISSING PAYMENT PROVIDERS (Aggregators)
-- =====================================================

-- Lygos - Popular in French-speaking Africa
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
('lygos', 'Lygos', 'mobile_money', 'https://api.lygosapp.com/v1', true, true, '/icons/lygos.svg',
 '["XOF", "XAF", "GNF", "CDF", "RWF", "KES", "UGX", "NGN"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "SN", "BF", "ML", "GN", "CM", "CD", "RW", "KE", "UG", "NG", "TG", "BJ", "NE"], "webhook_path": "/webhooks/lygos"}'::jsonb)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true;

-- YellowCard - Crypto ramp for Africa
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
('yellowcard', 'YellowCard', 'crypto_ramp', 'https://api.yellowcard.io/v1', true, true, '/icons/yellowcard.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "UGX", "TZS", "RWF", "BWP"]'::jsonb,
 '{"supports_crypto_ramp": true, "crypto_supported": ["BTC", "ETH", "USDT", "USDC"], "webhook_path": "/webhooks/yellowcard"}'::jsonb)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true;

-- FedaPay - Benin, Togo, Niger, Côte d'Ivoire
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
('fedapay', 'FedaPay', 'mobile_money', 'https://api.fedapay.com/v1', true, true, '/icons/fedapay.svg',
 '["XOF"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["BJ", "TG", "NE", "CI"], "operators": ["mtn", "moov"], "webhook_path": "/webhooks/fedapay"}'::jsonb)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true;

-- Moov Money - Côte d'Ivoire, Benin, Togo
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
('moov_money', 'Moov Money', 'mobile_money', 'https://api.moov-africa.bj/v1', true, true, '/icons/moov.svg',
 '["XOF", "XAF"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "BJ", "TG", "BF", "NE", "CM", "GA", "CG"], "webhook_path": "/webhooks/moov"}'::jsonb)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true;

-- =====================================================
-- 2. ADD COUNTRY MAPPINGS FOR NEW PROVIDERS
-- =====================================================

-- LYGOS Country Mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) VALUES
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'SN', 'Sénégal', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'BF', 'Burkina Faso', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'ML', 'Mali', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'TG', 'Togo', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'BJ', 'Bénin', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'NE', 'Niger', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'GN', 'Guinée', 'GNF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'CM', 'Cameroun', 'XAF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'CD', 'RD Congo', 'CDF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'RW', 'Rwanda', 'RWF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'KE', 'Kenya', 'KES', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'UG', 'Ouganda', 'UGX', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'lygos'), 'NG', 'Nigeria', 'NGN', 2, 1.5)
ON CONFLICT (provider_id, country_code) DO NOTHING;

-- YELLOWCARD Country Mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) VALUES
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'NG', 'Nigeria', 'NGN', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'GH', 'Ghana', 'GHS', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'KE', 'Kenya', 'KES', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'ZA', 'Afrique du Sud', 'ZAR', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'CI', 'Côte d''Ivoire', 'XOF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'SN', 'Sénégal', 'XOF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'CM', 'Cameroun', 'XAF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'UG', 'Ouganda', 'UGX', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'TZ', 'Tanzanie', 'TZS', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'RW', 'Rwanda', 'RWF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'BW', 'Botswana', 'BWP', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'BJ', 'Bénin', 'XOF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'TG', 'Togo', 'XOF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'ML', 'Mali', 'XOF', 2, 2.0),
((SELECT id FROM payment_providers WHERE name = 'yellowcard'), 'BF', 'Burkina Faso', 'XOF', 2, 2.0)
ON CONFLICT (provider_id, country_code) DO NOTHING;

-- FEDAPAY Country Mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) VALUES
((SELECT id FROM payment_providers WHERE name = 'fedapay'), 'BJ', 'Bénin', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'fedapay'), 'TG', 'Togo', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'fedapay'), 'NE', 'Niger', 'XOF', 2, 1.5),
((SELECT id FROM payment_providers WHERE name = 'fedapay'), 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5)
ON CONFLICT (provider_id, country_code) DO NOTHING;

-- MOOV MONEY Country Mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) VALUES
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'CI', 'Côte d''Ivoire', 'XOF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'BJ', 'Bénin', 'XOF', 2, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'TG', 'Togo', 'XOF', 2, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'BF', 'Burkina Faso', 'XOF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'NE', 'Niger', 'XOF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'CM', 'Cameroun', 'XAF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'GA', 'Gabon', 'XAF', 3, 1.0),
((SELECT id FROM payment_providers WHERE name = 'moov_money'), 'CG', 'Congo', 'XAF', 3, 1.0)
ON CONFLICT (provider_id, country_code) DO NOTHING;

-- =====================================================
-- 3. AGGREGATOR SETTINGS TABLE (for transfer-service)
-- =====================================================

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

-- =====================================================
-- 4. AGGREGATOR INSTANCES TABLE (multi-key support)
-- =====================================================

CREATE TABLE IF NOT EXISTS aggregator_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregator_id UUID REFERENCES aggregator_settings(id) ON DELETE CASCADE,
    instance_name VARCHAR(100) NOT NULL,
    
    -- API Credentials (reference to Vault or encrypted storage)
    api_credentials JSONB DEFAULT '{}',
    vault_secret_path VARCHAR(255),
    
    -- Status and configuration
    enabled BOOLEAN DEFAULT true,
    is_paused BOOLEAN DEFAULT false,
    is_global BOOLEAN DEFAULT false,
    pause_reason TEXT,
    paused_at TIMESTAMP,
    priority INT DEFAULT 50,
    health_status VARCHAR(20) DEFAULT 'unknown',
    
    -- Limits
    daily_limit DECIMAL(20, 8),
    monthly_limit DECIMAL(20, 8),
    daily_usage DECIMAL(20, 8) DEFAULT 0,
    monthly_usage DECIMAL(20, 8) DEFAULT 0,
    usage_reset_date DATE DEFAULT CURRENT_DATE,
    
    -- Country restrictions (NULL = all countries)
    restricted_countries TEXT[],
    
    -- Mode
    is_test_mode BOOLEAN DEFAULT true,
    
    -- Statistics
    total_transactions INT DEFAULT 0,
    total_volume DECIMAL(20, 8) DEFAULT 0,
    last_used_at TIMESTAMP,
    last_error TEXT,
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID
);

-- =====================================================
-- 5. AGGREGATOR INSTANCE WALLETS (multi-wallet per instance)
-- =====================================================

CREATE TABLE IF NOT EXISTS aggregator_instance_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES aggregator_instances(id) ON DELETE CASCADE,
    hot_wallet_id UUID NOT NULL,
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

-- =====================================================
-- 6. INDEXES
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_aggregator_settings_code ON aggregator_settings(provider_code);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_aggregator ON aggregator_instances(aggregator_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_enabled ON aggregator_instances(enabled);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_paused ON aggregator_instances(is_paused);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_instance ON aggregator_instance_wallets(instance_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_wallet ON aggregator_instance_wallets(hot_wallet_id);

-- =====================================================
-- 7. VIEW: aggregator_instances_with_details
-- This is the VIEW required by transfer-service repository
-- =====================================================

DROP VIEW IF EXISTS aggregator_instances_with_details;

CREATE VIEW aggregator_instances_with_details AS
SELECT 
    ai.id,
    ai.aggregator_id,
    ai.instance_name,
    ai.api_credentials,
    ai.vault_secret_path,
    ai.enabled,
    ai.is_paused,
    ai.is_global,
    ai.pause_reason,
    ai.paused_at,
    ai.priority,
    ai.health_status,
    ai.daily_limit,
    ai.monthly_limit,
    ai.daily_usage,
    ai.monthly_usage,
    ai.restricted_countries,
    ai.is_test_mode,
    ai.total_transactions,
    ai.total_volume,
    ai.last_used_at,
    ai.last_error,
    ai.created_at,
    ai.updated_at,
    -- Aggregator info
    agg.provider_code,
    agg.provider_name,
    agg.is_enabled AS aggregator_enabled,
    agg.is_demo_mode,
    agg.supports_deposit,
    agg.supports_withdrawal,
    agg.min_amount AS aggregator_min_amount,
    agg.max_amount AS aggregator_max_amount,
    agg.config AS aggregator_config,
    -- Primary wallet info (if linked)
    aiw.hot_wallet_id,
    aiw.currency AS wallet_currency,
    aiw.min_balance,
    aiw.max_balance,
    COALESCE(pa.balance, 0) AS hot_wallet_balance,
    -- Availability status
    CASE 
        WHEN NOT ai.enabled THEN 'instance_disabled'
        WHEN ai.is_paused THEN 'paused'
        WHEN NOT agg.is_enabled THEN 'aggregator_disabled'
        WHEN aiw.id IS NULL THEN 'no_wallet'
        WHEN pa.balance IS NULL THEN 'wallet_not_found'
        WHEN pa.balance < COALESCE(aiw.min_balance, 0) THEN 'insufficient_balance'
        ELSE 'available'
    END AS availability_status
FROM aggregator_instances ai
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id AND aiw.is_primary = true
LEFT JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id;

-- =====================================================
-- 8. SEED AGGREGATOR SETTINGS FROM PAYMENT PROVIDERS
-- =====================================================

INSERT INTO aggregator_settings (provider_code, provider_name, payment_provider_id, api_base_url, is_enabled, is_demo_mode, supports_deposit, supports_withdrawal, config)
SELECT 
    name,
    display_name,
    id,
    api_base_url,
    is_active,
    is_demo_mode,
    deposit_enabled,
    withdraw_enabled,
    config_json
FROM payment_providers
WHERE name IN ('demo', 'flutterwave', 'cinetpay', 'paystack', 'orange_money', 'mtn_momo', 'wave', 'stripe', 'paypal', 'lygos', 'yellowcard', 'fedapay', 'moov_money')
ON CONFLICT (provider_code) DO UPDATE SET
    provider_name = EXCLUDED.provider_name,
    api_base_url = EXCLUDED.api_base_url,
    is_enabled = EXCLUDED.is_enabled,
    is_demo_mode = EXCLUDED.is_demo_mode;

-- =====================================================
-- 9. SEED DEFAULT INSTANCES FOR EACH AGGREGATOR
-- =====================================================

DO $$
DECLARE
    v_aggregator RECORD;
    v_instance_id UUID;
    v_hot_wallet RECORD;
BEGIN
    -- For each aggregator, create a default instance
    FOR v_aggregator IN SELECT id, provider_code, provider_name FROM aggregator_settings LOOP
        -- Check if instance already exists
        SELECT id INTO v_instance_id 
        FROM aggregator_instances 
        WHERE aggregator_id = v_aggregator.id AND instance_name = 'Instance Principale';
        
        IF v_instance_id IS NULL THEN
            INSERT INTO aggregator_instances (
                aggregator_id,
                instance_name,
                vault_secret_path,
                enabled,
                is_global,
                priority,
                health_status,
                is_test_mode
            ) VALUES (
                v_aggregator.id,
                'Instance Principale',
                'secret/aggregators/' || v_aggregator.provider_code || '/default',
                true,
                true,
                100,
                'active',
                true
            )
            RETURNING id INTO v_instance_id;
            
            RAISE NOTICE 'Created instance for: %', v_aggregator.provider_name;
        END IF;
        
        -- Link to operations hot wallets
        IF v_instance_id IS NOT NULL THEN
            FOR v_hot_wallet IN 
                SELECT id, currency 
                FROM platform_accounts 
                WHERE account_type = 'operations' AND is_active = true
            LOOP
                INSERT INTO aggregator_instance_wallets (
                    instance_id,
                    hot_wallet_id,
                    currency,
                    is_primary,
                    priority,
                    enabled
                ) VALUES (
                    v_instance_id,
                    v_hot_wallet.id,
                    v_hot_wallet.currency,
                    v_hot_wallet.currency IN ('XOF', 'NGN', 'USD'),
                    CASE 
                        WHEN v_hot_wallet.currency IN ('XOF', 'NGN') THEN 100
                        ELSE 50
                    END,
                    true
                )
                ON CONFLICT (instance_id, hot_wallet_id, currency) DO NOTHING;
            END LOOP;
        END IF;
    END LOOP;
END $$;

-- =====================================================
-- 10. GRANT PERMISSIONS
-- =====================================================

-- Grant permissions if needed for transfer-service user
-- GRANT SELECT, INSERT, UPDATE ON aggregator_settings TO transfer_service;
-- GRANT SELECT, INSERT, UPDATE ON aggregator_instances TO transfer_service;
-- GRANT SELECT, INSERT, UPDATE ON aggregator_instance_wallets TO transfer_service;
-- GRANT SELECT ON aggregator_instances_with_details TO transfer_service;

RAISE NOTICE 'Aggregator infrastructure setup complete!';
