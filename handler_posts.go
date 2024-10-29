package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
)

func (apicfg *apiConfig) handlerPostsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Post string `json:"post"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if strings.TrimSpace(params.Post) == "" {
		respondWithError(w, http.StatusBadRequest, "post content cannot be empty")
		return
	}

	id := uuid.New()
	err = apicfg.DB.CreatePost(r.Context(), database.CreatePostParams{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Post:      params.Post,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create post")
		return
	}

	post, err := apicfg.DB.GetPost(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get post")
		return
	}

	postResp, err := databasePostToPost(post)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't convert note")
		return
	}

	respondWithJSON(w, http.StatusCreated, postResp)
}

func (apicfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apicfg.DB.GetPostsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get posts for user")
		return
	}

	postsResp, err := databasePostsToPosts(posts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't convert posts")
		return
	}

	respondWithJSON(w, http.StatusOK, postsResp)
}
