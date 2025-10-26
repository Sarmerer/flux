-- This migration removes the tables table from the core database
-- Tables are now stored in per-project databases

DROP TABLE IF EXISTS tables CASCADE;
