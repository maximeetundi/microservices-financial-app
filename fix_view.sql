-- Create the missing view manually
DROP VIEW IF EXISTS aggregator_instances_with_details;
CREATE VIEW aggregator_instances_with_details AS
SELECT 
    ai.id, ai.aggregator_id, ai.instance_name, ai.api_credentials, ai.vault_secret_path,
    ai.enabled, ai.is_paused, ai.is_global, ai.pause_reason, ai.paused_at,
    ai.priority, ai.health_status, ai.daily_limit, ai.monthly_limit,
    ai.daily_usage, ai.monthly_usage, ai.restricted_countries, ai.is_test_mode,
    ai.total_transactions, ai.total_volume, ai.last_used_at, ai.last_error,
    ai.created_at, ai.updated_at,
    agg.provider_code, agg.provider_name, agg.is_enabled AS aggregator_enabled,
    agg.is_demo_mode, agg.supports_deposit, agg.supports_withdrawal,
    agg.min_amount AS aggregator_min_amount, agg.max_amount AS aggregator_max_amount,
    agg.config AS aggregator_config,
    aiw.hot_wallet_id, aiw.currency AS wallet_currency,
    aiw.min_balance, aiw.max_balance,
    COALESCE(pa.balance, 0) AS hot_wallet_balance,
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
