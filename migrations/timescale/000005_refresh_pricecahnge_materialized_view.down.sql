DROP INDEX IF EXISTS idx_price_name_trgm;
ALTER TABLE price DROP COLUMN IF EXISTS item_name;