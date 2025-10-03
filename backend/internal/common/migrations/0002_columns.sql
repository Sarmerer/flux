CREATE TABLE IF NOT EXISTS columns (
    id BIGSERIAL PRIMARY KEY,
    table_id BIGINT NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('text','integer','boolean','timestamp')),
    required BOOLEAN NOT NULL DEFAULT FALSE,
    unique_col BOOLEAN NOT NULL DEFAULT FALSE,
    default_value TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT columns_table_name_unique UNIQUE (table_id, name)
);


