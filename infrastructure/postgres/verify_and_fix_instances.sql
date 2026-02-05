-- =====================================================
-- QUICK VERIFICATION AND FIX SCRIPT FOR AGGREGATOR INSTANCES
-- Run with: docker exec -i postgres psql -U crypto_bank -d crypto_bank < verify_and_fix_instances.sql
-- =====================================================

\echo '=============================================='
\echo '1. CHECKING PAYMENT PROVIDERS'
\echo '=============================================='

SELECT name, display_name, is_active, deposit_enabled, withdraw_enabled
FROM payment_providers
ORDER BY name;

\echo ''
\echo '=============================================='
\echo '2. CHECKING AGGREGATOR SETTINGS'
\echo '=============================================='

SELECT provider_code, provider_name, is_enabled, is_demo_mode, supports_deposit, supports_withdrawal
FROM aggregator_settings
ORDER BY provider_code;

\echo ''
\echo '=============================================='
\echo '3. CHECKING AGGREGATOR INSTANCES'
\echo '=============================================='

SELECT
    ai.id,
    agg.provider_code,
    ai.instance_name,
    ai.enabled,
    ai.is_paused,
    ai.is_global,
    ai.health_status
FROM aggregator_instances ai
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
ORDER BY agg.provider_code;

\echo ''
\echo '=============================================='
\echo '4. CHECKING PLATFORM ACCOUNTS (HOT WALLETS)'
\echo '=============================================='

SELECT id, currency, account_type, name, balance, is_active
FROM platform_accounts
WHERE account_type = 'operations'
ORDER BY currency;

\echo ''
\echo '=============================================='
\echo '5. CHECKING INSTANCE-WALLET LINKS'
\echo '=============================================='

SELECT
    agg.provider_code,
    aiw.currency,
    aiw.is_primary,
    aiw.enabled,
    pa.balance AS wallet_balance
FROM aggregator_instance_wallets aiw
JOIN aggregator_instances ai ON aiw.instance_id = ai.id
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
LEFT JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id
ORDER BY agg.provider_code, aiw.currency;

\echo ''
\echo '=============================================='
\echo '6. CHECKING VIEW DATA'
\echo '=============================================='

SELECT
    provider_code,
    instance_name,
    hot_wallet_currency,
    hot_wallet_balance,
    availability_status
FROM aggregator_instances_with_details
ORDER BY provider_code;

\echo ''
\echo '=============================================='
\echo '7. DIAGNOSING MISSING INSTANCES'
\echo '=============================================='

-- Show aggregators without instances
SELECT
    agg.provider_code,
    'MISSING INSTANCE' AS issue
FROM aggregator_settings agg
LEFT JOIN aggregator_instances ai ON agg.id = ai.aggregator_id
WHERE ai.id IS NULL;

-- Show instances without wallet links
SELECT
    agg.provider_code,
    ai.instance_name,
    'NO WALLET LINKED' AS issue
FROM aggregator_instances ai
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id
WHERE aiw.id IS NULL;

\echo ''
\echo '=============================================='
\echo '8. APPLYING FIXES'
\echo '=============================================='

-- Fix: Change hot_wallet_id type if needed
DO $$
BEGIN
    -- Check if column needs type change
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'aggregator_instance_wallets'
        AND column_name = 'hot_wallet_id'
        AND data_type = 'uuid'
    ) THEN
        RAISE NOTICE 'Changing hot_wallet_id from UUID to VARCHAR(36)...';
        ALTER TABLE aggregator_instance_wallets
            ALTER COLUMN hot_wallet_id TYPE VARCHAR(36) USING hot_wallet_id::text;
    ELSE
        RAISE NOTICE 'hot_wallet_id type is already correct';
    END IF;
END $$;

-- Fix: Create missing instances
DO $$
DECLARE
    v_agg RECORD;
    v_inst_id UUID;
    v_count INT := 0;
BEGIN
    FOR v_agg IN SELECT id, provider_code FROM aggregator_settings LOOP
        -- Check if instance exists
        SELECT id INTO v_inst_id
        FROM aggregator_instances
        WHERE aggregator_id = v_agg.id
        LIMIT 1;

        IF v_inst_id IS NULL THEN
            INSERT INTO aggregator_instances (
                aggregator_id, instance_name, enabled, is_global, is_paused, priority, health_status, is_test_mode
            ) VALUES (
                v_agg.id, 'Instance Principale', true, true, false, 100, 'active', true
            )
            RETURNING id INTO v_inst_id;
            v_count := v_count + 1;
            RAISE NOTICE 'Created instance for: %', v_agg.provider_code;
        END IF;
    END LOOP;
    RAISE NOTICE 'Created % new instances', v_count;
END $$;

-- Fix: Link wallets to instances that don't have links
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
        LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id
        WHERE aiw.id IS NULL
    LOOP
        FOR v_wallet IN
            SELECT id, currency FROM platform_accounts
            WHERE account_type = 'operations' AND is_active = true
        LOOP
            INSERT INTO aggregator_instance_wallets (
                instance_id, hot_wallet_id, currency, is_primary, priority, enabled
            ) VALUES (
                v_inst.instance_id, v_wallet.id, v_wallet.currency,
                v_wallet.currency = 'XOF', 100, true
            )
            ON CONFLICT DO NOTHING;
            v_link_count := v_link_count + 1;
        END LOOP;
        RAISE NOTICE 'Linked wallets to: %', v_inst.provider_code;
    END LOOP;
    RAISE NOTICE 'Created % wallet links', v_link_count;
END $$;

-- Fix: Enable all instances and aggregators
UPDATE aggregator_settings SET is_enabled = true WHERE is_enabled = false;
UPDATE aggregator_instances SET enabled = true, is_paused = false WHERE enabled = false OR is_paused = true;
UPDATE aggregator_instance_wallets SET enabled = true WHERE enabled = false;

\echo ''
\echo '=============================================='
\echo '9. FINAL VERIFICATION'
\echo '=============================================='

-- Show final status
SELECT
    agg.provider_code,
    COUNT(DISTINCT ai.id) AS instances,
    COUNT(DISTINCT aiw.id) AS wallet_links,
    CASE
        WHEN COUNT(DISTINCT ai.id) = 0 THEN '❌ NO INSTANCE'
        WHEN COUNT(DISTINCT aiw.id) = 0 THEN '⚠️ NO WALLETS'
        ELSE '✅ OK'
    END AS status
FROM aggregator_settings agg
LEFT JOIN aggregator_instances ai ON agg.id = ai.aggregator_id AND ai.enabled = true
LEFT JOIN aggregator_instance_wallets aiw ON ai.id = aiw.instance_id AND aiw.enabled = true
GROUP BY agg.provider_code
ORDER BY agg.provider_code;

\echo ''
\echo '=============================================='
\echo 'VERIFICATION COMPLETE'
\echo '=============================================='
