-- This migration removes the logs table from the core database
-- Logs are now stored in per-project databases

DROP TABLE IF EXISTS logs CASCADE;
