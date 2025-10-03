DO $$ BEGIN
  CREATE TYPE rel_kind AS ENUM ('o2o','o2m','m2o','m2m');
EXCEPTION WHEN duplicate_object THEN
  -- type already exists, no-op
  NULL;
END $$;

CREATE TABLE IF NOT EXISTS relationships (
    id BIGSERIAL PRIMARY KEY,
    source_table_id BIGINT NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    target_table_id BIGINT NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    kind rel_kind NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_rel UNIQUE (source_table_id, target_table_id, kind)
);


