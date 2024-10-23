package main

import (
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

func databaseUserToUser(user database.User) (User, error) {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
		ExpiresAt: user.ExpiresAt,
	}, nil
}

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
