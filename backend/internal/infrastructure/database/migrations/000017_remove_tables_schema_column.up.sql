-- Remove schema column from tables
-- Schema is now derived from information_schema/pg_catalog
ALTER TABLE tables DROP COLUMN IF EXISTS schema;
