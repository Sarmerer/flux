-- Add description column to tables
ALTER TABLE tables ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';
