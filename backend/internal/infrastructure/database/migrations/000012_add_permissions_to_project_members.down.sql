-- Remove permissions column from project_members table
DROP INDEX IF EXISTS idx_project_members_permissions;
ALTER TABLE project_members DROP COLUMN IF EXISTS permissions;
