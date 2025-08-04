-- name: CreateSystemUser :one
INSERT INTO system_users (id, created_at, form_data, needed_min_btu, needed_max_btu)
VALUES ($1, NOW(), $2, $3, $4)
RETURNING *;

-- name: GetSystemUser :one
SELECT * FROM system_users
WHERE id = $1;