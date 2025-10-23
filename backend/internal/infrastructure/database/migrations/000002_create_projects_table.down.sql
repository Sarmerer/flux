-- Rollback projects table
DROP INDEX IF EXISTS idx_projects_api_key;
DROP INDEX IF EXISTS idx_projects_owner_id;
DROP TABLE IF EXISTS projects CASCADE;
