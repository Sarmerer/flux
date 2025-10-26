-- This migration removes the workflows table from the core database
-- Workflows are now stored in per-project databases

DROP TABLE IF EXISTS workflows CASCADE;
