-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, name, password, api_key, api_key_expires_at)
VALUES (?, ?, ?, ?, ?, ?, ?);
--

-- name: GetUser :one
SELECT * FROM users WHERE api_key = ?;
--

-- name: GetUserByName :one
SELECT * FROM users WHERE name = ?;
--

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;
--

-- name: UpdateUser :exec
UPDATE users
SET updated_at = ?, api_key = ?, api_key_expires_at = ?
WHERE id = ?;
--
