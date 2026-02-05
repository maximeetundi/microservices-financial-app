-- =====================================================
-- MIGRATION: Add pause/global fields to aggregator_instances
-- And add tracking fields to transfers table
-- =====================================================

-- 1. Add pause and global fields to aggregator_instances
ALTER TABLE aggregator_instances 
ADD COLUMN IF NOT EXISTS is_paused BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS is_global BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS pause_reason TEXT,
ADD COLUMN IF NOT EXISTS paused_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS paused_by UUID REFERENCES admins(id);

-- Index for quick pause status lookup
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_paused 
ON aggregator_instances(is_paused) WHERE is_paused = true;

-- 2. Add tracking fields to transfers (OPTIONAL - can be NULL for crypto/internal transfers)
ALTER TABLE transfers
ADD COLUMN IF NOT EXISTS aggregator_instance_id UUID REFERENCES aggregator_instances(id),
ADD COLUMN IF NOT EXISTS hot_wallet_id UUID,
ADD COLUMN IF NOT EXISTS provider_code VARCHAR(50);

-- Index for finding transactions by instance (useful for debugging/auditing)
CREATE INDEX IF NOT EXISTS idx_transfers_instance 
ON transfers(aggregator_instance_id) WHERE aggregator_instance_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_transfers_hot_wallet 
ON transfers(hot_wallet_id) WHERE hot_wallet_id IS NOT NULL;

-- 3. Update the availability view to include pause status
DROP VIEW IF EXISTS aggregator_instance_wallets_available;

CREATE VIEW aggregator_instance_wallets_available AS
SELECT 
    aiw.id,
    aiw.instance_id,
    aiw.hot_wallet_id,
    aiw.is_primary,
    aiw.priority,
    aiw.min_balance,
    aiw.max_balance,
    aiw.enabled,
    -- Instance info
    ai.aggregator_id,
    ai.instance_name,
    ai.enabled AS instance_enabled,
    ai.is_paused,
    ai.is_global,
    ai.pause_reason,
    agg.provider_code,
    agg.provider_name,
    agg.is_enabled AS aggregator_enabled,
    -- Wallet info
    pa.currency AS wallet_currency,
    pa.balance AS wallet_balance,
    pa.account_type,
    -- Disponibilité (includes pause check)
    CASE 
        WHEN NOT aiw.enabled THEN 'wallet_disabled'
        WHEN ai.is_paused THEN 'instance_paused'
        WHEN NOT ai.enabled THEN 'instance_disabled'
        WHEN NOT agg.is_enabled THEN 'aggregator_disabled'
        WHEN pa.balance < COALESCE(aiw.min_balance, 0) THEN 'insufficient_balance'
        WHEN aiw.max_balance IS NOT NULL AND pa.balance > aiw.max_balance THEN 'balance_too_high'
        ELSE 'available'
    END AS availability_status
FROM aggregator_instance_wallets aiw
JOIN aggregator_instances ai ON aiw.instance_id = ai.id
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id;

-- 4. Update the instance details view to include new fields
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
    -- Aggregator info
    ai.aggregator_id,
    agg.provider_code,
    agg.provider_name,
    agg.logo_url AS provider_logo,
    agg.is_enabled AS aggregator_enabled,
    -- Primary hot wallet info (from first linked wallet)
    pw.hot_wallet_id,
    pa.account_type,
    pa.currency AS hot_wallet_currency,
    pa.balance AS hot_wallet_balance,
    pw.min_balance,
    pw.max_balance,
    -- Availability status
    CASE 
        WHEN ai.is_paused THEN 'instance_paused'
        WHEN NOT ai.enabled THEN 'instance_disabled'
        WHEN NOT agg.is_enabled THEN 'aggregator_disabled'
        WHEN pa.balance IS NOT NULL AND pa.balance < COALESCE(pw.min_balance, 0) THEN 'insufficient_balance'
        ELSE 'available'
    END AS availability_status
FROM aggregator_instances ai
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
LEFT JOIN LATERAL (
    SELECT * FROM aggregator_instance_wallets 
    WHERE instance_id = ai.id 
    ORDER BY is_primary DESC, priority DESC 
    LIMIT 1
) pw ON true
LEFT JOIN platform_accounts pa ON pw.hot_wallet_id = pa.id;

-- 5. Function to pause/resume an instance
CREATE OR REPLACE FUNCTION toggle_instance_pause(
    p_instance_id UUID,
    p_paused BOOLEAN,
    p_reason TEXT DEFAULT NULL,
    p_admin_id UUID DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    UPDATE aggregator_instances
    SET 
        is_paused = p_paused,
        pause_reason = CASE WHEN p_paused THEN p_reason ELSE NULL END,
        paused_at = CASE WHEN p_paused THEN CURRENT_TIMESTAMP ELSE NULL END,
        paused_by = CASE WHEN p_paused THEN p_admin_id ELSE NULL END,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_instance_id;
END;
$$ LANGUAGE plpgsql;

-- 6. Comment for documentation
COMMENT ON COLUMN aggregator_instances.is_paused IS 'Admin can pause instance to block all transactions';
COMMENT ON COLUMN aggregator_instances.is_global IS 'If true, accepts ALL countries (ignores restricted_countries)';
COMMENT ON COLUMN aggregator_instances.pause_reason IS 'Reason shown to users when transactions are blocked';
COMMENT ON COLUMN transfers.aggregator_instance_id IS 'Optional: Which instance processed this transaction (NULL for crypto/internal)';
COMMENT ON COLUMN transfers.hot_wallet_id IS 'Optional: Which hot wallet was used (NULL for crypto/internal)';

-- Done!
