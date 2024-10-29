package main

import (
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	Password        string    `json:"password"`
	ApiKey          string    `json:"api_key"`
	ApiKeyExpiresAt time.Time `json:"api_key_expires_at"`
}

func databaseUserToUser(user database.User) (User, error) {
	return User{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		Password:        user.Password,
		ApiKey:          user.ApiKey,
		ApiKeyExpiresAt: user.ApiKeyExpiresAt,
	}, nil
}

// type UsersToken struct {
// 	ID                    uuid.UUID `json:"id"`
// 	CreatedAt             time.Time `json:"created_at"`
// 	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
// 	RefreshToken          string    `json:"refresh_token"`
// 	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
// 	UserID                uuid.UUID `json:"user_id"`
// }

// func databaseUsersTokenToUsersToken(userstoken database.UsersToken) (UsersToken, error) {
// 	return UsersToken{
// 		ID:                    userstoken.ID,
// 		CreatedAt:             userstoken.CreatedAt,
// 		AccessTokenExpiresAt:  userstoken.AccessTokenExpiresAt,
// 		RefreshToken:          userstoken.RefreshToken,
// 		RefreshTokenExpiresAt: userstoken.RefreshTokenExpiresAt,
// 		UserID:                userstoken.UserID,
// 	}, nil
// }

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Post      string    `json:"post"`
	UserID    uuid.UUID `json:"user_id"`
}

func databasePostToPost(post database.Post) (Post, error) {
	return Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Post:      post.Post,
		UserID:    post.UserID,
	}, nil
}

func databasePostsToPosts(posts []database.Post) ([]Post, error) {
	postsResult := make([]Post, len(posts))
	for i, post := range posts {
		var err error
		postsResult[i], err = databasePostToPost(post)
		if err != nil {
			return nil, err
		}
	}

	return postsResult, nil
}
