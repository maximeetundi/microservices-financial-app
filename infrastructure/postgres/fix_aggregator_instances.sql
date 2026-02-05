-- =====================================================
-- FIX AGGREGATOR INSTANCES INITIALIZATION
-- This script fixes the aggregator instances and wallet linking
-- Run this after docker compose down -v && docker compose up -d
-- or execute manually: docker exec -i postgres psql -U crypto_bank -d crypto_bank < fix_aggregator_instances.sql
-- =====================================================

-- Start transaction
BEGIN;

-- =====================================================
-- 1. FIX TYPE COMPATIBILITY ISSUE
-- Change hot_wallet_id from UUID to VARCHAR(36) to match platform_accounts.id
-- =====================================================

-- First, drop dependent views
DROP VIEW IF EXISTS aggregator_instances_with_details CASCADE;
DROP VIEW IF EXISTS aggregator_instance_wallets_available CASCADE;

-- Drop existing constraints and recreate table with correct type
ALTER TABLE aggregator_instance_wallets
    DROP CONSTRAINT IF EXISTS aggregator_instance_wallets_hot_wallet_id_fkey;

-- Change column type if needed (if table exists with wrong type)
DO $$
BEGIN
    -- Check if column exists and change type
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'aggregator_instance_wallets'
        AND column_name = 'hot_wallet_id'
        AND data_type = 'uuid'
    ) THEN
        ALTER TABLE aggregator_instance_wallets
            ALTER COLUMN hot_wallet_id TYPE VARCHAR(36) USING hot_wallet_id::text;
        RAISE NOTICE 'Changed hot_wallet_id type from UUID to VARCHAR(36)';
    END IF;
END $$;

-- Drop old unique constraint and recreate
ALTER TABLE aggregator_instance_wallets
    DROP CONSTRAINT IF EXISTS aggregator_instance_wallets_instance_id_hot_wallet_id_curren_key;

-- Recreate unique constraint
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'aggregator_instance_wallets_unique_combo'
    ) THEN
        ALTER TABLE aggregator_instance_wallets
            ADD CONSTRAINT aggregator_instance_wallets_unique_combo
            UNIQUE(instance_id, hot_wallet_id, currency);
    END IF;
EXCEPTION WHEN duplicate_object THEN
    NULL;
END $$;

-- =====================================================
-- 2. ENSURE ALL PAYMENT PROVIDERS EXIST
-- =====================================================

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
 '{"countries": ["CI", "CM", "GH", "UG", "RW", "BJ"]}'::jsonb, true, true),
('wave', 'Wave', 'mobile_money', 'https://api.wave.com/v1', true, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["SN", "CI"]}'::jsonb, true, true),
('wave_ci', 'Wave CI', 'mobile_money', 'https://api.wave.com/v1', true, true, '/icons/wave.svg',
 '["XOF"]'::jsonb,
 '{"countries": ["CI"]}'::jsonb, true, true),
('stripe', 'Stripe / Carte Bancaire', 'card', 'https://api.stripe.com/v1', true, true, '/icons/stripe.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD"]'::jsonb,
 '{"supports_card": true}'::jsonb, true, true),
('paypal', 'PayPal', 'wallet', 'https://api-m.paypal.com', true, true, '/icons/paypal.svg',
 '["USD", "EUR", "GBP", "CAD", "AUD", "BRL", "MXN", "ILS"]'::jsonb,
 '{"supports_wallet": true, "supports_card": true}'::jsonb, true, true),
('lygos', 'Lygos', 'mobile_money', 'https://api.lygosapp.com/v1', true, true, '/icons/lygos.svg',
 '["XOF", "XAF", "GNF", "CDF", "RWF", "KES", "UGX", "NGN"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["CI", "SN", "BF", "ML", "GN", "CM", "CD", "RW", "KE", "UG", "NG", "TG", "BJ", "NE"]}'::jsonb, true, true),
('yellowcard', 'Yellow Card', 'crypto_ramp', 'https://api.yellowcard.io/v1', true, true, '/icons/yellowcard.svg',
 '["NGN", "GHS", "KES", "ZAR", "XOF", "XAF", "UGX", "TZS", "RWF", "BWP"]'::jsonb,
 '{"supports_crypto_ramp": true, "crypto_supported": ["BTC", "ETH", "USDT", "USDC"]}'::jsonb, true, true),
('fedapay', 'FedaPay', 'mobile_money', 'https://api.fedapay.com/v1', true, true, '/icons/fedapay.svg',
 '["XOF"]'::jsonb,
 '{"supports_mobile_money": true, "countries": ["BJ", "TG", "NE", "CI"]}'::jsonb, true, true),
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

-- =====================================================
-- 3. ENSURE AGGREGATOR_SETTINGS EXISTS FOR ALL PROVIDERS
-- =====================================================

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

-- =====================================================
-- 4. ENSURE PLATFORM_ACCOUNTS (HOT WALLETS) EXIST
-- =====================================================

INSERT INTO platform_accounts (id, currency, account_type, name, description, priority, balance, is_active) VALUES
-- XOF Wallets
(gen_random_uuid()::text, 'XOF', 'operations', 'Opérations XOF', 'Hot wallet principal XOF pour dépôts/retraits', 100, 100000000, true),
(gen_random_uuid()::text, 'XOF', 'reserve', 'Réserve XOF', 'Réserve XOF', 90, 1000000000, true),
(gen_random_uuid()::text, 'XOF', 'fees', 'Frais XOF', 'Frais collectés XOF', 50, 0, true),
-- XAF Wallets
(gen_random_uuid()::text, 'XAF', 'operations', 'Opérations XAF', 'Hot wallet principal XAF pour dépôts/retraits', 100, 100000000, true),
(gen_random_uuid()::text, 'XAF', 'reserve', 'Réserve XAF', 'Réserve XAF', 90, 1000000000, true),
-- NGN Wallets
(gen_random_uuid()::text, 'NGN', 'operations', 'Opérations NGN', 'Hot wallet principal NGN pour dépôts/retraits', 100, 500000000, true),
(gen_random_uuid()::text, 'NGN', 'reserve', 'Réserve NGN', 'Réserve NGN', 90, 1000000000, true),
-- USD Wallets
(gen_random_uuid()::text, 'USD', 'operations', 'Opérations USD', 'Hot wallet principal USD', 100, 100000, true),
(gen_random_uuid()::text, 'USD', 'reserve', 'Réserve USD', 'Réserve USD', 90, 1000000, true),
-- EUR Wallets
(gen_random_uuid()::text, 'EUR', 'operations', 'Opérations EUR', 'Hot wallet principal EUR', 100, 100000, true),
(gen_random_uuid()::text, 'EUR', 'reserve', 'Réserve EUR', 'Réserve EUR', 90, 1000000, true),
-- GHS Wallets
(gen_random_uuid()::text, 'GHS', 'operations', 'Opérations GHS', 'Hot wallet principal GHS', 100, 10000000, true),
-- KES Wallets
(gen_random_uuid()::text, 'KES', 'operations', 'Opérations KES', 'Hot wallet principal KES', 100, 50000000, true),
-- GNF Wallets
(gen_random_uuid()::text, 'GNF', 'operations', 'Opérations GNF', 'Hot wallet principal GNF', 100, 500000000, true)
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- 5. CREATE INSTANCES FOR ALL AGGREGATORS
-- =====================================================

DO $$
DECLARE
    v_agg RECORD;
    v_inst_id UUID;
BEGIN
    RAISE NOTICE 'Creating aggregator instances...';

    FOR v_agg IN SELECT id, provider_code, provider_name FROM aggregator_settings LOOP
        -- Check if instance exists
        SELECT id INTO v_inst_id
        FROM aggregator_instances
        WHERE aggregator_id = v_agg.id
        AND instance_name = 'Instance Principale'
        LIMIT 1;

        IF v_inst_id IS NULL THEN
            -- Create new instance
            INSERT INTO aggregator_instances (
                aggregator_id,
                instance_name,
                vault_secret_path,
                enabled,
                is_global,
                is_paused,
                priority,
                health_status,
                is_test_mode
            ) VALUES (
                v_agg.id,
                'Instance Principale',
                'secret/aggregators/' || v_agg.provider_code || '/default',
                true,
                true,  -- is_global = true means available for all countries
                false, -- not paused
                100,   -- high priority
                'active',
                true   -- test mode
            )
            RETURNING id INTO v_inst_id;

            RAISE NOTICE 'Created instance for aggregator: %', v_agg.provider_code;
        ELSE
            -- Update existing instance to ensure it's enabled
            UPDATE aggregator_instances
            SET enabled = true,
                is_global = true,
                is_paused = false,
                health_status = 'active',
                updated_at = NOW()
            WHERE id = v_inst_id;

            RAISE NOTICE 'Updated instance for aggregator: %', v_agg.provider_code;
        END IF;
    END LOOP;
END $$;

-- =====================================================
-- 6. LINK HOT WALLETS TO ALL INSTANCES
-- =====================================================

DO $$
DECLARE
    v_inst RECORD;
    v_wallet RECORD;
    v_link_count INT := 0;
BEGIN
    RAISE NOTICE 'Linking hot wallets to instances...';

    -- For each instance, link all operation wallets
    FOR v_inst IN
        SELECT ai.id AS instance_id, agg.provider_code
        FROM aggregator_instances ai
        JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
        WHERE ai.enabled = true
    LOOP
        FOR v_wallet IN
            SELECT id, currency
            FROM platform_accounts
            WHERE account_type = 'operations'
            AND is_active = true
        LOOP
            -- Insert wallet link (using VARCHAR type for hot_wallet_id)
            INSERT INTO aggregator_instance_wallets (
                instance_id,
                hot_wallet_id,
                currency,
                is_primary,
                priority,
                min_balance,
                enabled
            ) VALUES (
                v_inst.instance_id,
                v_wallet.id,  -- VARCHAR(36)
                v_wallet.currency,
                -- Set XOF as primary for most African providers
                v_wallet.currency = 'XOF',
                CASE
                    WHEN v_wallet.currency = 'XOF' THEN 100
                    WHEN v_wallet.currency = 'XAF' THEN 95
                    WHEN v_wallet.currency = 'NGN' THEN 90
                    WHEN v_wallet.currency = 'USD' THEN 80
                    ELSE 50
                END,
                0,  -- min_balance
                true
            )
            ON CONFLICT (instance_id, hot_wallet_id, currency) DO UPDATE SET
                enabled = true,
                is_primary = EXCLUDED.is_primary,
                priority = EXCLUDED.priority;

            v_link_count := v_link_count + 1;
        END LOOP;
    END LOOP;

    RAISE NOTICE 'Created/Updated % wallet links', v_link_count;
END $$;

-- =====================================================
-- 7. RECREATE THE VIEW WITH CORRECT JOINS
-- Column order MUST match Go code's scanInstanceWithDetails function exactly!
-- =====================================================

DROP VIEW IF EXISTS aggregator_instances_with_details CASCADE;

CREATE VIEW aggregator_instances_with_details AS
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
    -- Columns 11-18: Limits and usage
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
-- 8. VERIFICATION QUERIES
-- =====================================================

-- Show all aggregators and their instance count
DO $$
DECLARE
    v_result RECORD;
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============ AGGREGATOR STATUS ============';

    FOR v_result IN
        SELECT
            agg.provider_code,
            agg.is_enabled,
            COUNT(DISTINCT ai.id) AS instance_count,
            COUNT(DISTINCT aiw.id) AS wallet_links
        FROM aggregator_settings agg
        LEFT JOIN aggregator_instances ai ON agg.id = ai.aggregator_id AND ai.enabled = true
        LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id AND aiw.enabled = true
        GROUP BY agg.provider_code, agg.is_enabled
        ORDER BY agg.provider_code
    LOOP
        RAISE NOTICE 'Provider: %, Enabled: %, Instances: %, Wallet Links: %',
            v_result.provider_code,
            v_result.is_enabled,
            v_result.instance_count,
            v_result.wallet_links;
    END LOOP;

    RAISE NOTICE '';
    RAISE NOTICE '============ AVAILABLE INSTANCES ============';

    FOR v_result IN
        SELECT provider_code, instance_name, hot_wallet_currency, hot_wallet_balance, availability_status
        FROM aggregator_instances_with_details
        ORDER BY provider_code
    LOOP
        RAISE NOTICE '% | % | Currency: % | Balance: % | Status: %',
            v_result.provider_code,
            v_result.instance_name,
            COALESCE(v_result.hot_wallet_currency, 'N/A'),
            COALESCE(v_result.hot_wallet_balance::text, '0'),
            v_result.availability_status;
    END LOOP;
END $$;

-- =====================================================
-- 9. GRANT PERMISSIONS (if using separate DB users)
-- =====================================================

-- Uncomment if needed:
-- GRANT SELECT ON aggregator_instances_with_details TO transfer_service;
-- GRANT SELECT, INSERT, UPDATE ON aggregator_instances TO transfer_service;
-- GRANT SELECT, INSERT, UPDATE ON aggregator_instance_wallets TO transfer_service;

COMMIT;

-- Final message
DO $$ BEGIN RAISE NOTICE '✅ Aggregator instances fix completed successfully!'; END $$;
