// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users_key.sql

package database

import (
	"context"
)

const createUserRfKey = `-- name: CreateUserRfKey :one
INSERT INTO users_key (id, created_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

type CreateUserRfKeyParams struct {
	ID                    string
	CreatedAt             string
	AccessTokenExpiresAt  string
	RefreshToken          string
	RefreshTokenExpiresAt string
	UserID                string
}

func (q *Queries) CreateUserRfKey(ctx context.Context, arg CreateUserRfKeyParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createUserRfKey,
		arg.ID,
		arg.CreatedAt,
		arg.AccessTokenExpiresAt,
		arg.RefreshToken,
		arg.RefreshTokenExpiresAt,
		arg.UserID,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const getRfKeyByUserID = `-- name: GetRfKeyByUserID :one

SELECT id, created_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id FROM users_key WHERE user_id = $1
LIMIT 1
`

func (q *Queries) GetRfKeyByUserID(ctx context.Context, userID string) (UsersKey, error) {
	row := q.db.QueryRowContext(ctx, getRfKeyByUserID, userID)
	var i UsersKey
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.AccessTokenExpiresAt,
		&i.RefreshToken,
		&i.RefreshTokenExpiresAt,
		&i.UserID,
	)
	return i, err
}

const getUserByRfKey = `-- name: GetUserByRfKey :one

SELECT id, created_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id FROM users_key WHERE refresh_token = $1 
LIMIT 1
`

func (q *Queries) GetUserByRfKey(ctx context.Context, refreshToken string) (UsersKey, error) {
	row := q.db.QueryRowContext(ctx, getUserByRfKey, refreshToken)
	var i UsersKey
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.AccessTokenExpiresAt,
		&i.RefreshToken,
		&i.RefreshTokenExpiresAt,
		&i.UserID,
	)
	return i, err
}

const updateUserRfKey = `-- name: UpdateUserRfKey :one

UPDATE users_key
SET access_token_expires_at = $1, refresh_token = $2, refresh_token_expires_at = $3
WHERE user_id = $4
RETURNING id, refresh_token
`

type UpdateUserRfKeyParams struct {
	AccessTokenExpiresAt  string
	RefreshToken          string
	RefreshTokenExpiresAt string
	UserID                string
}

type UpdateUserRfKeyRow struct {
	ID           string
	RefreshToken string
}

func (q *Queries) UpdateUserRfKey(ctx context.Context, arg UpdateUserRfKeyParams) (UpdateUserRfKeyRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserRfKey,
		arg.AccessTokenExpiresAt,
		arg.RefreshToken,
		arg.RefreshTokenExpiresAt,
		arg.UserID,
	)
	var i UpdateUserRfKeyRow
	err := row.Scan(&i.ID, &i.RefreshToken)
	return i, err
}
