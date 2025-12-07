DROP MATERIALIZED VIEW IF EXISTS price_changes_24h CASCADE;

-- +migrate StatementBegin
CREATE MATERIALIZED VIEW price_changes_24h 
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('24 hour', time) AS bucket,
    item_id,
    item_name,
    FIRST(sell_price, time) AS open_price,
    LAST(sell_price, time) AS close_price,
    LAST(sell_listings, time) AS sell_listings,
    CASE
        WHEN FIRST(sell_price, time) = 0 THEN 0
        ELSE (
            (LAST(sell_price, time) - FIRST(sell_price, time))::float / FIRST(sell_price, time) * 100
        )
    END AS change_pct
FROM prices
GROUP BY bucket, item_id, item_name
WITH NO DATA;
-- +migrate StatementEnd

CREATE INDEX idx_price_changes_24h_name_trgm 
ON price_changes_24h USING gin (item_name gin_trgm_ops);

CREATE INDEX idx_price_changes_24h_item_id 
ON price_changes_24h (item_id);