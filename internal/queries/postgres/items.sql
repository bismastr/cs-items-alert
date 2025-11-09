-- name: CreateItem :one
INSERT INTO items (
    name,
    hash_name
) VALUES ($1, $2)
ON CONFLICT (hash_name) DO NOTHING
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

