-- Migration: Add missing columns to payment_providers table
-- Run this on the existing database to add the new columns

-- Add deposit_enabled column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS deposit_enabled BOOLEAN DEFAULT true;

-- Add withdraw_enabled column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS withdraw_enabled BOOLEAN DEFAULT true;

-- Add fee_percentage column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS fee_percentage DECIMAL(5,2) DEFAULT 0;

-- Add fee_fixed column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS fee_fixed DECIMAL(20,2) DEFAULT 0;

-- Add min_transaction column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS min_transaction DECIMAL(20,2) DEFAULT 100;

-- Add max_transaction column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS max_transaction DECIMAL(20,2) DEFAULT 10000000;

-- Add daily_limit column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS daily_limit DECIMAL(20,2) DEFAULT 50000000;

-- Add priority column
ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS priority INT DEFAULT 1;

-- Update demo provider to have proper defaults
UPDATE payment_providers 
SET deposit_enabled = true, withdraw_enabled = true, priority = 1
WHERE name = 'demo';

-- Update other providers
UPDATE payment_providers 
SET deposit_enabled = true, withdraw_enabled = true, priority = 2
WHERE name != 'demo';

-- Verify columns were added
SELECT column_name, data_type, column_default 
FROM information_schema.columns 
WHERE table_name = 'payment_providers' 
ORDER BY ordinal_position;
