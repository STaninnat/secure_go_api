-- name: CreateUserRefreshToken :one
INSERT INTO users_token (id, created_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
--

-- name: UpdateUserRefreshToken :one
UPDATE users_token
SET refresh_token = $1, refresh_token_expires_at = $2
WHERE user_id = $3
RETURNING id, refresh_token;
--

-- name: GetRefreshTokenByUserID :one
SELECT * FROM users_token WHERE user_id = $1
LIMIT 1;
--

-- name: GetUserByRefreshToken :one
SELECT * FROM users_token WHERE refresh_token = $1 
LIMIT 1;
--