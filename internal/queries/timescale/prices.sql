-- name: InsertPrice :exec
INSERT INTO prices (
    item_id,
    sell_price,
    sell_listings,
    item_name
) VALUES ($1, $2, $3, $4);

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

-- name: GetPriceChangesByItemIDs :many
SELECT 
    item_id::integer,
      item_name::text,
    bucket::timestamptz,
    open_price::integer,
    close_price::integer,
    sell_listings::integer,
    change_pct::float
FROM price_changes_24h
WHERE bucket = DATE_TRUNC('day', NOW() - INTERVAL '1 day')
  AND item_id = ANY(sqlc.arg(item_ids)::int[])
LIMIT sqlc.arg(max_results);

-- name: GetAllPriceChanges :many
SELECT 
        item_id::integer,
        item_name::text,
        bucket::timestamptz,
        open_price::integer,
        close_price::integer,
        sell_listings::integer,
        change_pct::float
FROM price_changes_24h
WHERE bucket = DATE_TRUNC('day', NOW() - INTERVAL '1 day')
ORDER BY 
    CASE WHEN sqlc.arg(sort_by)::text = 'gainers' THEN change_pct END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'losers' THEN change_pct END ASC,
    change_pct DESC
LIMIT $1 OFFSET $2;

-- name: SearchPriceChangesByName :many
SELECT 
    item_id::integer,
    item_name::text,
    bucket::timestamptz,
    open_price::integer,
    close_price::integer,
    sell_listings::integer,
    change_pct::float
FROM price_changes_24h
WHERE bucket = DATE_TRUNC('day', NOW() - INTERVAL '1 day')
    AND similarity(item_name, sqlc.arg(query)) > 0.3
ORDER BY 
    CASE WHEN sqlc.arg(sort_by)::text = 'gainers' THEN change_pct END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'losers' THEN change_pct END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'relevance' OR sqlc.arg(sort_by)::text = '' THEN similarity(item_name, sqlc.arg(query)) END DESC,
    change_pct DESC
LIMIT $1 OFFSET $2;

-- name: CountSearchPriceChangesByName :one
SELECT 
    COUNT(*) as count
FROM price_changes_24h
WHERE bucket = DATE_TRUNC('day', NOW() - INTERVAL '1 day')
    AND similarity(item_name, sqlc.arg(query)) > 0.3;

-- name: GetItemSparklineWeekly :many
SELECT 
    item_id::integer,
    array_agg(close_price ORDER BY bucket ASC)::int[] AS sparkline
FROM price_changes_1h
WHERE item_id = ANY(sqlc.arg(item_id)::int[])
  AND bucket >= NOW() - INTERVAL '7 days'
GROUP BY item_id;

-- name: GetItemPriceChartByDay :many
SELECT 
    bucket::timestamptz,
    open_price::integer,
    close_price::integer,
    sell_listings::integer,
    change_pct::float
FROM price_changes_24h
WHERE item_id = $1
  AND bucket >= NOW() - $2::interval
ORDER BY bucket ASC;

-- name: GetItemPriceChartByHour :many
SELECT 
    bucket::timestamptz,
    open_price::integer,
    close_price::integer,
    sell_listings::integer,
    change_pct::float
FROM price_changes_1h
WHERE item_id = $1
  AND bucket >= NOW() - $2::interval
ORDER BY bucket ASC;