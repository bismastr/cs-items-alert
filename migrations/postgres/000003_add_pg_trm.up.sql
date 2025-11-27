CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_items_name_trgm
ON items USING gin (name gin_trgm_ops);