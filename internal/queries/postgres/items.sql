-- name: CreateItem :one
INSERT INTO items (
    name,
    hash_name
) VALUES ($1, $2)
ON CONFLICT (hash_name) 
DO UPDATE SET 
    name = EXCLUDED.name,
    updated_at = NOW()
RETURNING *;

-- name: GetItemByHashName :one
SELECT
    id,
    name,
    hash_name,
    created_at,
    updated_at
FROM items
WHERE hash_name = $1;

-- name: GetItems :many
SELECT
    id,
    name,
    hash_name,
    created_at,
    updated_at
FROM items;

-- name: GetItemByID :many
SELECT
    id,
    name,
    hash_name,
    created_at,
    updated_at
FROM items
WHERE id = ANY(sqlc.arg(ids)::int[]);

-- name: SearchItemsByName :many
WITH score AS (
    SELECT 
        id,
        name, 
        similarity(name, sqlc.arg(name)) AS sim_score
    FROM items
)
SELECT 
    id,
    name,
    sim_score
FROM score
WHERE sim_score > 0.1
ORDER BY sim_score DESC
LIMIT $1 OFFSET $2;

-- name: SearchItemsCount :one
WITH score AS (
    SELECT 
        similarity(name, sqlc.arg(name)) AS sim_score
    FROM items
)
SELECT 
    COUNT(*) as count
FROM score
WHERE sim_score > 0.1;

-- name: GetAllItemsCount :one
SELECT 
    COUNT(*) as count
FROM items;

