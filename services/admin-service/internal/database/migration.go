package database

import (
	"database/sql"
	"log"
)

// MigrateAggregatorSchema executes the necessary SQL to ensure the aggregator infrastructure exists.
// This includes tables, views, and default data required for the transfer-service.
func MigrateAggregatorSchema(db *sql.DB) error {
	log.Println("[Database] 🔄 Starting aggregator schema migration...")

	query := `
-- =====================================================
-- AGGREGATOR INFRASTRUCTURE (for transfer-service)
-- Tables, views, seed data for payment aggregator instances
-- =====================================================

-- Add missing payment providers (Lygos, YellowCard, FedaPay, Moov)
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json) VALUES
('lygos', 'Lygos', 'mobile_money', 'https://api.lygosapp.com/v1', true, true, '/icons/lygos.svg',
 '["XOF", "XAF", "GNF", "CDF", "RWF", "KES", "UGX", "NGN"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "SN", "BF", "ML", "GN", "CM", "CD", "RW", "KE", "UG", "NG", "TG", "BJ", "NE"]}'::jsonb),
('yellowcard', 'YellowCard', 'crypto_ramp', 'https://api.yellowcard.io/v1', true, true, '/icons/yellowcard.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "UGX", "TZS", "RWF", "BWP"]'::jsonb,
 '{"supports_crypto_ramp": true, "crypto_supported": ["BTC", "ETH", "USDT", "USDC"]}'::jsonb),
('fedapay', 'FedaPay', 'mobile_money', 'https://api.fedapay.com/v1', true, true, '/icons/fedapay.svg',
 '["XOF"]'::jsonb, '{"supports_mobile_money": true, "countries": ["BJ", "TG", "NE", "CI"]}'::jsonb),
('moov_money', 'Moov Money', 'mobile_money', 'https://api.moov-africa.bj/v1', true, true, '/icons/moov.svg',
 '["XOF", "XAF"]'::jsonb, '{"supports_mobile_money": true, "countries": ["CI", "BJ", "TG", "BF", "NE", "CM", "GA", "CG"]}'::jsonb)
ON CONFLICT (name) DO UPDATE SET is_active = true;

-- Country mappings
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'BF', 'Burkina Faso', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'ML', 'Mali', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'TG', 'Togo', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CM', 'Cameroun', 'XAF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CD', 'RD Congo', 'CDF', 2, 1.5 FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;

INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'NG', 'Nigeria', 'NGN', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'GH', 'Ghana', 'GHS', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'KE', 'Kenya', 'KES', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'ZA', 'Afrique du Sud', 'ZAR', 2, 2.0 FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;

INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'TG', 'Togo', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5 FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;

INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 3, 1.0 FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.0 FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'TG', 'Togo', 'XOF', 2, 1.0 FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'BF', 'Burkina Faso', 'XOF', 3, 1.0 FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage) 
SELECT id, 'CM', 'Cameroun', 'XAF', 3, 1.0 FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;

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

-- Indexes
CREATE INDEX IF NOT EXISTS idx_aggregator_settings_code ON aggregator_settings(provider_code);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_aggregator ON aggregator_instances(aggregator_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_enabled ON aggregator_instances(enabled);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_paused ON aggregator_instances(is_paused);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_instance ON aggregator_instance_wallets(instance_id);

-- View for transfer-service to query instances with details
-- We use CREATE OR REPLACE to ensure it's up to date
DROP VIEW IF EXISTS aggregator_instances_with_details;
CREATE VIEW aggregator_instances_with_details AS
SELECT 
    ai.id,
    ai.instance_name,
    ai.enabled,
    ai.priority,
    ai.is_test_mode,
    ai.is_paused,
    ai.is_global,
    ai.pause_reason,
    ai.paused_at,
    ai.restricted_countries,
    ai.daily_limit,
    ai.monthly_limit,
    ai.daily_usage,
    ai.monthly_usage,
    ai.total_transactions,
    ai.total_volume,
    ai.last_used_at,
    ai.created_at,
    ai.updated_at,
    ai.aggregator_id,
    agg.provider_code,
    agg.provider_name,
    '' AS provider_logo,
    agg.is_enabled AS aggregator_enabled,
    aiw.hot_wallet_id::text AS hot_wallet_id,
    'operations' AS account_type,
    aiw.currency AS hot_wallet_currency,
    COALESCE(pa.balance, 0) AS hot_wallet_balance,
    aiw.min_balance,
    aiw.max_balance,
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
LEFT JOIN platform_accounts pa ON aiw.hot_wallet_id::text = pa.id;

-- Seed aggregator_settings from payment_providers
INSERT INTO aggregator_settings (provider_code, provider_name, payment_provider_id, api_base_url, is_enabled, is_demo_mode, supports_deposit, supports_withdrawal, config)
SELECT name, display_name, id, api_base_url, is_active, is_demo_mode, is_active, is_active, config_json::jsonb
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
            INSERT INTO aggregator_instances (aggregator_id, instance_name, vault_secret_path, enabled, is_global, priority, health_status, is_test_mode)
            VALUES (v_agg.id, 'Instance Principale', 'secret/aggregators/' || v_agg.provider_code || '/default', true, true, 100, 'active', true)
            RETURNING id INTO v_inst_id;
        END IF;
        IF v_inst_id IS NOT NULL THEN
            FOR v_wallet IN SELECT id, currency FROM platform_accounts WHERE account_type = 'operations' AND is_active = true LOOP
                INSERT INTO aggregator_instance_wallets (instance_id, hot_wallet_id, currency, is_primary, priority, enabled)
                VALUES (v_inst_id, v_wallet.id::uuid, v_wallet.currency, v_wallet.currency IN ('XOF', 'NGN', 'USD'), CASE WHEN v_wallet.currency IN ('XOF', 'NGN') THEN 100 ELSE 50 END, true)
                ON CONFLICT (instance_id, hot_wallet_id, currency) DO NOTHING;
            END LOOP;
        END IF;
    END LOOP;
END $$;
	`

	if _, err := db.Exec(query); err != nil {
		log.Printf("[Database] ❌ Aggregator schema migration failed: %v", err)
		return err
	}

	log.Println("[Database] ✅ Aggregator schema migration completed successfully")
	return nil
}
