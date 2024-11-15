package main

import (
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
)

type User struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	Password        string    `json:"password"`
	ApiKey          string    `json:"api_key"`
	ApiKeyExpiresAt time.Time `json:"api_key_expires_at"`
}

func databaseUserToUser(user database.User) (User, error) {
	// createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
	// if err != nil {
	// 	return User{}, err
	// }

	// updatedAt, err := time.Parse(time.RFC3339, user.UpdatedAt)
	// if err != nil {
	// 	return User{}, err
	// }

	// apiKeyExpiresAt, err := time.Parse(time.RFC3339, user.ApiKeyExpiresAt)
	// if err != nil {
	// 	return User{}, err
	// }

	return User{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		ApiKeyExpiresAt: user.ApiKeyExpiresAt,
	}, nil
}

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Post      string    `json:"post"`
	UserID    string    `json:"user_id"`
}

func databasePostToPost(post database.Post) (Post, error) {
	// createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
	// if err != nil {
	// 	return Post{}, err
	// }

	// updatedAt, err := time.Parse(time.RFC3339, post.UpdatedAt)
	// if err != nil {
	// 	return Post{}, err
	// }
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
