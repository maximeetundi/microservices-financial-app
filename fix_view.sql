-- Create the missing view manually with corrected column order
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
