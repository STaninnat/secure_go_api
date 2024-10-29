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
	"golang.org/x/crypto/bcrypt"
)

func (apicfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	_, err = apicfg.DB.GetUserByName(r.Context(), params.Name)
	if err == nil {
		respondWithError(w, http.StatusBadRequest, "username already exists")
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "please enter a password")
		return
	}
	if len(params.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "password must be least 8 ")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	apiKey, hashedApiKey, err := generateAndHashAPIKey()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate apikey")
		return
	}

	apiKeyExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)

	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Name:            params.Name,
		Password:        string(hashedPassword),
		ApiKey:          hashedApiKey,
		ApiKeyExpiresAt: apiKeyExpiresAt,
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

	userResp.ApiKey = apiKey

	respondWithJSON(w, http.StatusCreated, userResp)
}

func (apicfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	userResp, err := databaseUserToUser(user)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't convert user")
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}

func generateAndHashAPIKey() (string, string, error) {
	apiKey, err := generateRandomSHA256HASH()
	if err != nil {
		return "", "", err
	}
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return apiKey, string(hashedKey), nil
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
