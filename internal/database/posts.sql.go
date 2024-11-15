// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
)

const createPost = `-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, post, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, post, user_id
`

type CreatePostParams struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Post      string
	UserID    string
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Post,
		arg.UserID,
	)
	return err
}

const getPost = `-- name: GetPost :one

SELECT id, created_at, updated_at, post, user_id FROM posts WHERE id = $1
`

func (q *Queries) GetPost(ctx context.Context, id string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Post,
		&i.UserID,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many

SELECT id, created_at, updated_at, post, user_id FROM posts WHERE user_id = $1
`

func (q *Queries) GetPostsForUser(ctx context.Context, userID string) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Post,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
