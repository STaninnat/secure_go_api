-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, password, api_key, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
--

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;
--

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;
--

-- name: GetuserByID :one
SELECT * FROM users WHERE id = $1;
--