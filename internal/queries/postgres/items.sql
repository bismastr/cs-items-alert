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

