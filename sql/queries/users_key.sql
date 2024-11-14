-- name: CreateUserRfKey :one
INSERT INTO users_key (id, created_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
--

-- name: UpdateUserRfKey :one
UPDATE users_key
SET access_token_expires_at = $1, refresh_token = $2, refresh_token_expires_at = $3
WHERE user_id = $4
RETURNING id, refresh_token;
--

-- name: GetRfKeyByUserID :one
SELECT * FROM users_key WHERE user_id = $1
LIMIT 1;
--

-- name: GetUserByRfKey :one
SELECT * FROM users_key WHERE refresh_token = $1 
LIMIT 1;
--