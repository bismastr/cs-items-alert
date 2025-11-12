-- name: InsertPrice :exec
INSERT INTO prices (
    item_id,
    sell_price,
    sell_listings
) VALUES ($1, $2, $3);

-- name: Get24HourPricesChanges :many
WITH latest_sell_price AS (
    SELECT DISTINCT ON (item_id)
        item_id,
        sell_price as latest_sell_price
    FROM prices
    WHERE time >= NOW() - INTERVAL '2 hours'
    ORDER BY item_id, time DESC 
),
old_sell_price AS (
    SELECT DISTINCT ON (item_id)
        item_id,
        sell_price as old_sell_price
    FROM prices
    WHERE time BETWEEN NOW() - INTERVAL '26 hours' AND NOW() - INTERVAL '22 hours'
    ORDER BY item_id, time DESC
)
SELECT
    l.item_id,
    l.latest_sell_price,
    o.old_sell_price
FROM latest_sell_price l
JOIN old_sell_price o ON l.item_id = o.item_id;
