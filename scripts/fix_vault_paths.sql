-- ============================================================================
-- Fix Vault Paths in provider_instances table
-- ============================================================================
-- This script removes the "/default" suffix from vault_secret_path
-- to match the actual Vault structure: secret/aggregators/{provider}
--
-- Run this script against the admin database:
--   psql -h localhost -U admin -d crypto_bank_admin -f fix_vault_paths.sql
-- ============================================================================

-- Show current paths before update
SELECT
    pi.id,
    pp.name as provider_code,
    pp.display_name,
    pi.name as instance_name,
    pi.vault_secret_path as old_path
FROM provider_instances pi
JOIN payment_providers pp ON pi.provider_id = pp.id
WHERE pi.vault_secret_path LIKE '%/default';

-- Update paths: remove /default suffix
UPDATE provider_instances
SET vault_secret_path = REPLACE(vault_secret_path, '/default', '')
WHERE vault_secret_path LIKE '%/default';

-- Also update any paths that might have been set incorrectly with double slashes
UPDATE provider_instances
SET vault_secret_path = REPLACE(vault_secret_path, '//', '/')
WHERE vault_secret_path LIKE '%//%';

-- Show updated paths
SELECT
    pi.id,
    pp.name as provider_code,
    pp.display_name,
    pi.name as instance_name,
    pi.vault_secret_path as new_path
FROM provider_instances pi
JOIN payment_providers pp ON pi.provider_id = pp.id
ORDER BY pp.name;

-- Verify all paths now match the expected format
SELECT
    COUNT(*) as total_instances,
    COUNT(*) FILTER (WHERE vault_secret_path LIKE 'secret/aggregators/%') as correct_format,
    COUNT(*) FILTER (WHERE vault_secret_path LIKE '%/default') as still_has_default
FROM provider_instances;

-- ============================================================================
-- Expected Vault paths after fix:
-- ============================================================================
-- Provider        | Expected Vault Path
-- ----------------|----------------------------------
-- cinetpay        | secret/aggregators/cinetpay
-- paypal          | secret/aggregators/paypal
-- stripe          | secret/aggregators/stripe
-- flutterwave     | secret/aggregators/flutterwave
-- wave_money      | secret/aggregators/wave_money
-- mtn_money       | secret/aggregators/mtn_money
-- orange_money    | secret/aggregators/orange_money
-- moov_money      | secret/aggregators/moov_money
-- paystack        | secret/aggregators/paystack
-- ============================================================================

COMMIT;
