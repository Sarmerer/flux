-- Restore schema column to tables
-- Note: This is for rollback only - data cannot be recovered
ALTER TABLE tables ADD COLUMN IF NOT EXISTS schema JSONB;
