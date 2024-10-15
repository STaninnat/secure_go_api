package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
)

func (apicfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	apiKey, err := generateRandomSHA256HASH()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate apikey")
		return
	}

	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ApiKey:    apiKey,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusCreated, userResp)
}

func (apicfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}

func generateRandomSHA256HASH() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)
	hashString := hex.EncodeToString(hash[:])
	return hashString, nil
}
