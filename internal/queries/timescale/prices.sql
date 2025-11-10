-- name: InsertPrice :exec
INSERT INTO prices (
    item_id,
    sell_price,
    sell_listings
) VALUES ($1, $2, $3);