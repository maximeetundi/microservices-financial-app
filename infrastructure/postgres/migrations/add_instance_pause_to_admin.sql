-- Add is_paused and is_global columns to provider_instances table
-- This migration adds pause functionality for temporarily blocking transactions on an instance

-- Add new columns
ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS is_paused BOOLEAN DEFAULT FALSE;
ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS is_global BOOLEAN DEFAULT FALSE;
ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS pause_reason TEXT;
ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS paused_at TIMESTAMP;

-- Add index for efficient querying of paused instances
CREATE INDEX IF NOT EXISTS idx_provider_instances_paused ON provider_instances(is_paused) WHERE is_paused = TRUE;
CREATE INDEX IF NOT EXISTS idx_provider_instances_global ON provider_instances(is_global) WHERE is_global = TRUE;

-- Update existing records to have default values
UPDATE provider_instances SET is_paused = FALSE WHERE is_paused IS NULL;
UPDATE provider_instances SET is_global = FALSE WHERE is_global IS NULL;

COMMENT ON COLUMN provider_instances.is_paused IS 'When TRUE, no transactions are processed via this instance';
COMMENT ON COLUMN provider_instances.is_global IS 'When TRUE, this instance accepts all countries (no geo restrictions)';
COMMENT ON COLUMN provider_instances.pause_reason IS 'Optional reason for why the instance was paused';
COMMENT ON COLUMN provider_instances.paused_at IS 'Timestamp when the instance was paused';
