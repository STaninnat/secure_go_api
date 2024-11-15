-- name: CreateUserRfKey :exec
INSERT INTO users_key (id, created_at, updated_at ,access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?);
--

-- name: UpdateUserRfKey :exec
UPDATE users_key
SET updated_at = ?,access_token_expires_at = ?, refresh_token = ?, refresh_token_expires_at = ?
WHERE user_id = ?;
--

-- name: GetRfKeyByUserID :one
SELECT * FROM users_key WHERE user_id = ?
LIMIT 1;
--

-- name: GetUserByRfKey :one
SELECT * FROM users_key WHERE refresh_token = ? 
LIMIT 1;
--