-- +migrate StatementBegin
CREATE MATERIALIZED VIEW price_changes_24h 
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('24 hour', time) AS bucket,
    item_id,
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
GROUP BY bucket, item_id
WITH NO DATA;
-- +migrate StatementEnd

SELECT add_continuous_aggregate_policy(
    'price_changes_24h',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour'
);

CREATE INDEX idx_price_changes_cursor 
ON price_changes_24h (bucket DESC, item_id ASC);

CREATE INDEX idx_price_changes_pct 
ON price_changes_24h (change_pct DESC);

CREATE INDEX idx_price_changes_item 
ON price_changes_24h (item_id);