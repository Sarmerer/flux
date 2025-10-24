-- Add permissions column to project_members table
ALTER TABLE project_members ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '[]'::jsonb;

-- Create index on permissions for faster lookups
CREATE INDEX IF NOT EXISTS idx_project_members_permissions ON project_members USING gin(permissions);
