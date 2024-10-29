-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, password, api_key, api_key_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
--

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;
--

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
--

-- name: UpdateUserApiKey :exec
UPDATE users
SET api_key = $1, api_key_expires_at = $2
WHERE id = $3;
--
