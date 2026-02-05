-- =====================================================
-- COMPLETE FIX SCRIPT FOR AGGREGATOR SYSTEM
-- This script fixes ALL issues with aggregators, instances, and country mappings
--
-- Run: docker exec -i postgres psql -U admin -d crypto_bank < COMPLETE_FIX.sql
-- Or after docker compose down -v && docker compose up -d
-- =====================================================

BEGIN;

\echo '=============================================='
\echo 'STEP 1: Fixing table structure'
\echo '=============================================='

-- Drop dependent views first
DROP VIEW IF EXISTS aggregator_instances_with_details CASCADE;
DROP VIEW IF EXISTS aggregator_instance_wallets_available CASCADE;

-- Fix hot_wallet_id type from UUID to VARCHAR(36)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'aggregator_instance_wallets'
        AND column_name = 'hot_wallet_id'
        AND data_type = 'uuid'
    ) THEN
        ALTER TABLE aggregator_instance_wallets
            ALTER COLUMN hot_wallet_id TYPE VARCHAR(36) USING hot_wallet_id::text;
        RAISE NOTICE '✅ Changed hot_wallet_id from UUID to VARCHAR(36)';
    ELSE
        RAISE NOTICE '✅ hot_wallet_id type already correct';
    END IF;
END $$;

\echo '=============================================='
\echo 'STEP 2: Ensuring all payment providers exist'
\echo '=============================================='

-- Insert/Update all payment providers
INSERT INTO payment_providers (name, display_name, provider_type, api_base_url, is_active, is_demo_mode, logo_url, supported_currencies, config_json, deposit_enabled, withdraw_enabled) VALUES
('demo', 'Mode Démo', 'all', NULL, true, true, '/icons/demo.svg',
 '["XOF", "XAF", "NGN", "GHS", "KES", "USD", "EUR"]'::jsonb,
 '{"description": "Mode test - crédite directement le compte sans paiement réel"}'::jsonb, true, true),
('flutterwave', 'Flutterwave', 'all', 'https://api.flutterwave.com/v3', true, true, '/icons/flutterwave.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "USD", "EUR", "GBP"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true}'::jsonb, true, true),
('cinetpay', 'CinetPay', 'mobile_money', 'https://api-checkout.cinetpay.com/v2', true, true, '/icons/cinetpay.svg',
 '["XOF", "XAF", "GNF"]'::jsonb,
 '{"supports_mobile_money": true, "operators": ["orange_money", "mtn_momo", "moov_money", "wave"]}'::jsonb, true, true),
('paystack', 'Paystack', 'all', 'https://api.paystack.co', true, true, '/icons/paystack.svg',
 '["NGN", "GHS", "ZAR", "KES"]'::jsonb,
 '{"supports_mobile_money": true, "supports_card": true, "supports_bank": true}'::jsonb, true, true),
('orange_money', 'Orange Money', 'mobile_money', 'https://api.orange.com/orange-money-webpay', true, true, '/icons/orange.svg',
 '["XOF", "XAF"]'::jsonb,
 '{"countries": ["CI", "SN", "ML", "BF", "CM", "GN"]}'::jsonb, true, true),
('mtn_momo', 'MTN Mobile Money', 'mobile_money', 'https://sandbox.momodeveloper.mtn.com', true, true, '/icons/mtn.svg',
 '["XOF", "XAF", "GHS", "UGX", "RWF"]'::jsonb,
 '{"countries": ["CI", "CM", "GH", "UG", "RW", "BJ", "SN", "BF"]}'::jsonb, true, true),
('wave', 'Wave', 'mobile_money', 'https://api.wave.com/v1', true, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["SN", "CI"]}'::jsonb, true, true),
('wave_ci', 'Wave CI', 'mobile_money', 'https://api.wave.com/v1', true, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["CI"]}'::jsonb, true, true),
('stripe', 'Stripe / Carte Bancaire', 'card', 'https://api.stripe.com/v1', true, true, '/icons/stripe.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD", "XOF", "XAF", "NGN"]'::jsonb,
 '{"supports_card": true}'::jsonb, true, true),
('paypal', 'PayPal', 'wallet', 'https://api-m.paypal.com', true, true, '/icons/paypal.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD"]'::jsonb,
 '{"supports_wallet": true, "supports_card": true}'::jsonb, true, true),
('lygos', 'Lygos', 'mobile_money', 'https://api.lygosapp.com/v1', true, true, '/icons/lygos.svg',
 '["XOF", "XAF", "GNF", "CDF", "RWF", "KES", "UGX", "NGN"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "SN", "BF", "ML", "GN", "CM", "CD", "RW", "KE", "UG", "NG", "TG", "BJ", "NE", "LR"]}'::jsonb, true, true),
('yellowcard', 'Yellow Card', 'crypto_ramp', 'https://api.yellowcard.io/v1', true, true, '/icons/yellowcard.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "UGX", "TZS", "RWF"]'::jsonb,
 '{"supports_crypto_ramp": true, "crypto_supported": ["BTC", "ETH", "USDT", "USDC"]}'::jsonb, true, true),
('fedapay', 'FedaPay', 'mobile_money', 'https://api.fedapay.com/v1', true, true, '/icons/fedapay.svg',
 '["XOF"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["BJ", "TG", "NE", "CI", "CM"]}'::jsonb, true, true),
('moov_money', 'Moov Money', 'mobile_money', 'https://api.moov-africa.bj/v1', true, true, '/icons/moov.svg',
 '["XOF", "XAF"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "BJ", "TG", "BF", "NE", "CM", "GA", "CG"]}'::jsonb, true, true)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    api_base_url = EXCLUDED.api_base_url,
    supported_currencies = EXCLUDED.supported_currencies,
    config_json = EXCLUDED.config_json,
    is_active = true,
    deposit_enabled = true,
    withdraw_enabled = true,
    updated_at = NOW();

\echo '=============================================='
\echo 'STEP 3: Creating aggregator_settings'
\echo '=============================================='

-- Sync aggregator_settings from payment_providers
INSERT INTO aggregator_settings (provider_code, provider_name, payment_provider_id, api_base_url, is_enabled, is_demo_mode, supports_deposit, supports_withdrawal, config)
SELECT
    pp.name,
    pp.display_name,
    pp.id,
    pp.api_base_url,
    pp.is_active,
    pp.is_demo_mode,
    COALESCE(pp.deposit_enabled, true),
    COALESCE(pp.withdraw_enabled, true),
    pp.config_json
FROM payment_providers pp
WHERE pp.name IN ('demo', 'flutterwave', 'cinetpay', 'paystack', 'orange_money', 'mtn_momo', 'wave', 'wave_ci', 'stripe', 'paypal', 'lygos', 'yellowcard', 'fedapay', 'moov_money')
ON CONFLICT (provider_code) DO UPDATE SET
    is_enabled = true,
    is_demo_mode = EXCLUDED.is_demo_mode,
    supports_deposit = true,
    supports_withdrawal = true,
    updated_at = NOW();

\echo '=============================================='
\echo 'STEP 4: Ensuring platform_accounts (hot wallets) exist'
\echo '=============================================='

-- Create hot wallets if they don't exist
INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'XOF', 'operations', 'Opérations XOF', 'Hot wallet XOF', 100, 100000000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'XOF' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'XAF', 'operations', 'Opérations XAF', 'Hot wallet XAF', 100, 100000000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'XAF' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'NGN', 'operations', 'Opérations NGN', 'Hot wallet NGN', 100, 500000000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'NGN' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'USD', 'operations', 'Opérations USD', 'Hot wallet USD', 100, 100000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'USD' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'EUR', 'operations', 'Opérations EUR', 'Hot wallet EUR', 100, 100000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'EUR' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'GHS', 'operations', 'Opérations GHS', 'Hot wallet GHS', 100, 10000000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'GHS' AND account_type = 'operations');

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active)
SELECT gen_random_uuid()::text, 'KES', 'operations', 'Opérations KES', 'Hot wallet KES', 100, 50000000, true
WHERE NOT EXISTS (SELECT 1 FROM platform_accounts WHERE currency = 'KES' AND account_type = 'operations');

\echo '=============================================='
\echo 'STEP 5: Creating aggregator instances'
\echo '=============================================='

-- Create instances for all aggregators
DO $$
DECLARE
    v_agg RECORD;
    v_inst_id UUID;
BEGIN
    FOR v_agg IN SELECT id, provider_code, provider_name FROM aggregator_settings LOOP
        -- Check if instance exists
        SELECT id INTO v_inst_id
        FROM aggregator_instances
        WHERE aggregator_id = v_agg.id
        LIMIT 1;

        IF v_inst_id IS NULL THEN
            INSERT INTO aggregator_instances (
                aggregator_id, instance_name, vault_secret_path,
                enabled, is_global, is_paused, priority, health_status, is_test_mode
            ) VALUES (
                v_agg.id, 'Instance Principale', 'secret/aggregators/' || v_agg.provider_code || '/default',
                true, true, false, 100, 'active', true
            )
            RETURNING id INTO v_inst_id;
            RAISE NOTICE '✅ Created instance for: %', v_agg.provider_code;
        ELSE
            -- Enable existing instance
            UPDATE aggregator_instances
            SET enabled = true, is_global = true, is_paused = false, health_status = 'active'
            WHERE id = v_inst_id;
            RAISE NOTICE '✅ Updated instance for: %', v_agg.provider_code;
        END IF;
    END LOOP;
END $$;

\echo '=============================================='
\echo 'STEP 6: Linking hot wallets to instances'
\echo '=============================================='

-- Link all operation wallets to all instances
DO $$
DECLARE
    v_inst RECORD;
    v_wallet RECORD;
    v_link_count INT := 0;
BEGIN
    FOR v_inst IN
        SELECT ai.id AS instance_id, agg.provider_code
        FROM aggregator_instances ai
        JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
        WHERE ai.enabled = true
    LOOP
        FOR v_wallet IN
            SELECT id, currency FROM platform_accounts
            WHERE account_type = 'operations' AND is_active = true
        LOOP
            INSERT INTO aggregator_instance_wallets (
                instance_id, hot_wallet_id, currency, is_primary, priority, min_balance, enabled
            ) VALUES (
                v_inst.instance_id, v_wallet.id, v_wallet.currency,
                v_wallet.currency = 'XOF', -- XOF is primary
                CASE WHEN v_wallet.currency = 'XOF' THEN 100
                     WHEN v_wallet.currency = 'XAF' THEN 95
                     WHEN v_wallet.currency = 'NGN' THEN 90
                     ELSE 50 END,
                0, true
            )
            ON CONFLICT (instance_id, hot_wallet_id, currency) DO UPDATE SET
                enabled = true, is_primary = EXCLUDED.is_primary;
            v_link_count := v_link_count + 1;
        END LOOP;
    END LOOP;
    RAISE NOTICE '✅ Created/Updated % wallet links', v_link_count;
END $$;

\echo '=============================================='
\echo 'STEP 7: Seeding provider_countries mappings'
\echo '=============================================='

-- DEMO (Global)
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 1, 0, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('CM','Cameroun','XAF'),('BJ','Bénin','XOF'),
        ('TG','Togo','XOF'),('BF','Burkina Faso','XOF'),('ML','Mali','XOF'),('NE','Niger','XOF'),
        ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('KE','Kenya','KES')) AS c(code,name,curr)
WHERE payment_providers.name = 'demo' ON CONFLICT DO NOTHING;

-- LYGOS
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.5, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('BF','Burkina Faso','XOF'),('ML','Mali','XOF'),
        ('TG','Togo','XOF'),('BJ','Bénin','XOF'),('NE','Niger','XOF'),('CM','Cameroun','XAF'),
        ('CD','RD Congo','CDF'),('GN','Guinée','GNF'),('LR','Liberia','LRD')) AS c(code,name,curr)
WHERE payment_providers.name = 'lygos' ON CONFLICT DO NOTHING;

-- CINETPAY
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.5, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('CM','Cameroun','XAF'),('BF','Burkina Faso','XOF'),
        ('ML','Mali','XOF'),('TG','Togo','XOF'),('BJ','Bénin','XOF'),('NE','Niger','XOF'),('GN','Guinée','GNF')) AS c(code,name,curr)
WHERE payment_providers.name = 'cinetpay' ON CONFLICT DO NOTHING;

-- FLUTTERWAVE
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.4, true FROM payment_providers,
(VALUES ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('KE','Kenya','KES'),('ZA','Afrique du Sud','ZAR'),
        ('UG','Uganda','UGX'),('TZ','Tanzania','TZS'),('RW','Rwanda','RWF'),('CI','Côte d''Ivoire','XOF'),('CM','Cameroun','XAF')) AS c(code,name,curr)
WHERE payment_providers.name = 'flutterwave' ON CONFLICT DO NOTHING;

-- PAYSTACK
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.5, true FROM payment_providers,
(VALUES ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('ZA','Afrique du Sud','ZAR'),('KE','Kenya','KES')) AS c(code,name,curr)
WHERE payment_providers.name = 'paystack' ON CONFLICT DO NOTHING;

-- ORANGE MONEY
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 1, 1.0, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('ML','Mali','XOF'),('BF','Burkina Faso','XOF'),
        ('CM','Cameroun','XAF'),('GN','Guinée','GNF')) AS c(code,name,curr)
WHERE payment_providers.name = 'orange_money' ON CONFLICT DO NOTHING;

-- MTN MOMO
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 1, 1.0, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('CM','Cameroun','XAF'),('SN','Sénégal','XOF'),('BJ','Bénin','XOF'),
        ('GH','Ghana','GHS'),('UG','Uganda','UGX'),('RW','Rwanda','RWF'),('BF','Burkina Faso','XOF')) AS c(code,name,curr)
WHERE payment_providers.name = 'mtn_momo' ON CONFLICT DO NOTHING;

-- WAVE
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 1, 1.0, true FROM payment_providers,
(VALUES ('SN','Sénégal','XOF'),('CI','Côte d''Ivoire','XOF')) AS c(code,name,curr)
WHERE payment_providers.name = 'wave' ON CONFLICT DO NOTHING;

INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'wave_ci' ON CONFLICT DO NOTHING;

-- MOOV MONEY
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.0, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('BJ','Bénin','XOF'),('TG','Togo','XOF'),('BF','Burkina Faso','XOF'),
        ('NE','Niger','XOF'),('CM','Cameroun','XAF'),('GA','Gabon','XAF'),('CG','Congo','XAF')) AS c(code,name,curr)
WHERE payment_providers.name = 'moov_money' ON CONFLICT DO NOTHING;

-- FEDAPAY
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 1.5, true FROM payment_providers,
(VALUES ('BJ','Bénin','XOF'),('TG','Togo','XOF'),('NE','Niger','XOF'),('CI','Côte d''Ivoire','XOF'),('CM','Cameroun','XAF')) AS c(code,name,curr)
WHERE payment_providers.name = 'fedapay' ON CONFLICT DO NOTHING;

-- YELLOWCARD
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 2, 2.0, true FROM payment_providers,
(VALUES ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),
        ('KE','Kenya','KES'),('ZA','Afrique du Sud','ZAR'),('CM','Cameroun','XAF'),('UG','Uganda','UGX'),
        ('TZ','Tanzania','TZS'),('RW','Rwanda','RWF')) AS c(code,name,curr)
WHERE payment_providers.name = 'yellowcard' ON CONFLICT DO NOTHING;

-- STRIPE (Global card payments)
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 5, 2.9, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('CM','Cameroun','XAF'),
        ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('KE','Kenya','KES')) AS c(code,name,curr)
WHERE payment_providers.name = 'stripe' ON CONFLICT DO NOTHING;

-- PAYPAL (Global)
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, code, name, curr, 5, 3.9, true FROM payment_providers,
(VALUES ('CI','Côte d''Ivoire','XOF'),('SN','Sénégal','XOF'),('CM','Cameroun','XAF'),
        ('NG','Nigeria','NGN'),('GH','Ghana','GHS'),('KE','Kenya','KES')) AS c(code,name,curr)
WHERE payment_providers.name = 'paypal' ON CONFLICT DO NOTHING;

\echo '=============================================='
\echo 'STEP 8: Recreating the view'
\echo '=============================================='

-- Recreate view with correct column order for Go code
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

\echo '=============================================='
\echo 'STEP 9: Verification'
\echo '=============================================='

-- Show aggregator status
SELECT
    agg.provider_code,
    agg.is_enabled AS enabled,
    COUNT(DISTINCT ai.id) AS instances,
    COUNT(DISTINCT aiw.id) AS wallet_links,
    COUNT(DISTINCT pc.id) AS country_mappings,
    CASE
        WHEN COUNT(DISTINCT ai.id) = 0 THEN '❌ NO INSTANCE'
        WHEN COUNT(DISTINCT aiw.id) = 0 THEN '⚠️ NO WALLETS'
        WHEN COUNT(DISTINCT pc.id) = 0 THEN '⚠️ NO COUNTRIES'
        ELSE '✅ OK'
    END AS status
FROM aggregator_settings agg
LEFT JOIN aggregator_instances ai ON agg.id = ai.aggregator_id AND ai.enabled = true
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id AND aiw.enabled = true
LEFT JOIN payment_providers pp ON agg.payment_provider_id = pp.id
LEFT JOIN provider_countries pc ON pp.id = pc.provider_id AND pc.is_active = true
GROUP BY agg.provider_code, agg.is_enabled
ORDER BY agg.provider_code;

-- Show available instances
SELECT provider_code, instance_name, hot_wallet_currency, hot_wallet_balance, availability_status
FROM aggregator_instances_with_details
WHERE availability_status = 'available'
ORDER BY provider_code;

-- Show country mappings count
SELECT
    pp.name as provider,
    COUNT(pc.id) as countries,
    string_agg(pc.country_code, ',' ORDER BY pc.country_code) as codes
FROM payment_providers pp
LEFT JOIN provider_countries pc ON pp.id = pc.provider_id AND pc.is_active = true
GROUP BY pp.name
HAVING COUNT(pc.id) > 0
ORDER BY pp.name;

COMMIT;

\echo ''
\echo '=============================================='
\echo '✅ COMPLETE FIX APPLIED SUCCESSFULLY!'
\echo '=============================================='
\echo ''
\echo 'Next steps:'
\echo '1. Restart transfer-service: docker restart transfer-service'
\echo '2. Restart admin-service: docker restart admin-service'
\echo '3. Test deposit from frontend'
\echo ''
