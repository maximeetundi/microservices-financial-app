-- =====================================================
-- SEED PROVIDER COUNTRIES - COMPLETE MAPPINGS
-- This script adds all country mappings for payment providers
-- Run: docker exec -i postgres psql -U admin -d crypto_bank < seed_provider_countries.sql
-- =====================================================

BEGIN;

-- =====================================================
-- 1. DEMO PROVIDER (Global - all countries)
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TG', 'Togo', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ML', 'Mali', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NE', 'Niger', 'XOF', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 1, 0, true FROM payment_providers WHERE name = 'demo' ON CONFLICT DO NOTHING;

-- =====================================================
-- 2. LYGOS - Mobile Money Africa
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ML', 'Mali', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TG', 'Togo', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CD', 'RD Congo', 'CDF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GN', 'Guinée', 'GNF', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'LR', 'Liberia', 'LRD', 2, 1.5, true FROM payment_providers WHERE name = 'lygos' ON CONFLICT DO NOTHING;

-- =====================================================
-- 3. CINETPAY - Mobile Money UEMOA/CEMAC
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ML', 'Mali', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TG', 'Togo', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GN', 'Guinée', 'GNF', 2, 1.5, true FROM payment_providers WHERE name = 'cinetpay' ON CONFLICT DO NOTHING;

-- =====================================================
-- 4. FLUTTERWAVE - Pan-African
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ZA', 'Afrique du Sud', 'ZAR', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'UG', 'Uganda', 'UGX', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TZ', 'Tanzania', 'TZS', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'RW', 'Rwanda', 'RWF', 2, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 3, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 3, 1.4, true FROM payment_providers WHERE name = 'flutterwave' ON CONFLICT DO NOTHING;

-- =====================================================
-- 5. PAYSTACK - Nigeria/Ghana Focus
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 1, 1.5, true FROM payment_providers WHERE name = 'paystack' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 2, 1.5, true FROM payment_providers WHERE name = 'paystack' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ZA', 'Afrique du Sud', 'ZAR', 2, 1.5, true FROM payment_providers WHERE name = 'paystack' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 3, 1.5, true FROM payment_providers WHERE name = 'paystack' ON CONFLICT DO NOTHING;

-- =====================================================
-- 6. ORANGE MONEY
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ML', 'Mali', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GN', 'Guinée', 'GNF', 1, 1.0, true FROM payment_providers WHERE name = 'orange_money' ON CONFLICT DO NOTHING;

-- =====================================================
-- 7. MTN MOMO
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 1, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 1, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'UG', 'Uganda', 'UGX', 1, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'RW', 'Rwanda', 'RWF', 1, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'mtn_momo' ON CONFLICT DO NOTHING;

-- =====================================================
-- 8. WAVE (Sénégal & Côte d'Ivoire)
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'wave' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'wave' ON CONFLICT DO NOTHING;

-- Wave CI spécifique
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'wave_ci' ON CONFLICT DO NOTHING;

-- =====================================================
-- 9. MOOV MONEY
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TG', 'Togo', 'XOF', 1, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BF', 'Burkina Faso', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GA', 'Gabon', 'XAF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CG', 'Congo', 'XAF', 2, 1.0, true FROM payment_providers WHERE name = 'moov_money' ON CONFLICT DO NOTHING;

-- =====================================================
-- 10. FEDAPAY
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'BJ', 'Bénin', 'XOF', 1, 1.5, true FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TG', 'Togo', 'XOF', 1, 1.5, true FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NE', 'Niger', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 1.5, true FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 3, 1.5, true FROM payment_providers WHERE name = 'fedapay' ON CONFLICT DO NOTHING;

-- =====================================================
-- 11. YELLOWCARD (Crypto Ramp)
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'ZA', 'Afrique du Sud', 'ZAR', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'UG', 'Uganda', 'UGX', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'TZ', 'Tanzania', 'TZS', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'RW', 'Rwanda', 'RWF', 2, 2.0, true FROM payment_providers WHERE name = 'yellowcard' ON CONFLICT DO NOTHING;

-- =====================================================
-- 12. STRIPE (Global - Cards)
-- These are global providers, but we add some mappings for visibility
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 5, 2.9, true FROM payment_providers WHERE name = 'stripe' ON CONFLICT DO NOTHING;

-- =====================================================
-- 13. PAYPAL (Global)
-- =====================================================
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CI', 'Côte d''Ivoire', 'XOF', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'SN', 'Sénégal', 'XOF', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'CM', 'Cameroun', 'XAF', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'NG', 'Nigeria', 'NGN', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'GH', 'Ghana', 'GHS', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;
INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, fee_percentage, is_active)
SELECT id, 'KE', 'Kenya', 'KES', 5, 3.9, true FROM payment_providers WHERE name = 'paypal' ON CONFLICT DO NOTHING;

-- =====================================================
-- VERIFICATION
-- =====================================================

DO $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count FROM provider_countries;
    RAISE NOTICE '✅ Total provider_countries mappings: %', v_count;

    -- Show summary by provider
    FOR v_count IN
        SELECT pp.name, COUNT(pc.id) as country_count
        FROM payment_providers pp
        LEFT JOIN provider_countries pc ON pp.id = pc.provider_id
        GROUP BY pp.name
        ORDER BY pp.name
    LOOP
        RAISE NOTICE 'Provider mappings created';
    END LOOP;
END $$;

-- Show summary
SELECT
    pp.name as provider,
    COUNT(pc.id) as countries_mapped,
    string_agg(pc.country_code, ', ' ORDER BY pc.country_code) as country_codes
FROM payment_providers pp
LEFT JOIN provider_countries pc ON pp.id = pc.provider_id AND pc.is_active = true
GROUP BY pp.name
ORDER BY pp.name;

COMMIT;

-- Final message
DO $$ BEGIN RAISE NOTICE '✅ Provider countries seeding completed!'; END $$;
