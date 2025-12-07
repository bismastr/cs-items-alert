ALTER TABLE prices ADD COLUMN item_name TEXT;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_price_name_trgm
ON prices USING gin (item_name gin_trgm_ops);