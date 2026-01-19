CREATE MATERIALIZED VIEW price_changes_1h 
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', time) AS bucket,
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

SELECT add_continuous_aggregate_policy(
    'price_changes_1h',
    start_offset => INTERVAL '30 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour'
);