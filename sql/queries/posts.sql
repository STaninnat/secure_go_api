-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, post, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
--

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;
--

-- name: GetPostsForUser :many
SELECT * FROM posts WHERE user_id = $1;
--