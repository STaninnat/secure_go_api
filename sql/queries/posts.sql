-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, post, user_id)
VALUES (?, ?, ?, ?, ?);
--

-- name: GetPost :one
SELECT * FROM posts WHERE id = ?;
--

-- name: GetPostsForUser :many
SELECT * FROM posts WHERE user_id = ?;
--