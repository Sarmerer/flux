-- Restore logs table to core database (for rollback only)
CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY,
    level VARCHAR(10) NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    database_id UUID,
    table_id UUID,
    workflow_id UUID,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    request_id VARCHAR(255),
    operation VARCHAR(255),
    fields JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_logs_project_id ON logs(project_id);
CREATE INDEX IF NOT EXISTS idx_logs_database_id ON logs(database_id);
CREATE INDEX IF NOT EXISTS idx_logs_table_id ON logs(table_id);
CREATE INDEX IF NOT EXISTS idx_logs_workflow_id ON logs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
CREATE INDEX IF NOT EXISTS idx_logs_operation ON logs(operation);
CREATE INDEX IF NOT EXISTS idx_logs_project_timestamp ON logs(project_id, timestamp DESC);
