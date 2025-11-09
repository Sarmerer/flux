ALTER TABLE project_members
ADD CONSTRAINT check_valid_role
CHECK (role IN ('owner', 'admin', 'editor', 'viewer'));
