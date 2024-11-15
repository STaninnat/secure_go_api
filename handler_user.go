package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	if params.Name == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if !isValidUserName(params.Name) {
		respondWithError(w, http.StatusBadRequest, "invalid username format")
		return
	}

	_, err = apicfg.DB.GetUserByName(r.Context(), params.Name)
	if err == nil {
		respondWithError(w, http.StatusBadRequest, "username already exists")
		return
	}

	if len(params.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "password must be at least 8 ")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	_, hashedApiKey, err := generateAndHashAPIKey()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate apikey")
		return
	}

	apiKeyExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)

	err = apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:              uuid.New().String(),
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Name:            params.Name,
		Password:        string(hashedPassword),
		ApiKey:          hashedApiKey,
		ApiKeyExpiresAt: apiKeyExpiresAt,
	})
	if err != nil {
		log.Printf("Error while creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	jwtExpiresAt := time.Now().UTC().Add(15 * time.Minute)

	user, err := apicfg.DB.GetUser(r.Context(), hashedApiKey)
	if err != nil {
		log.Printf("Error while getting user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get user")
		return
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "invalid user ID")
		return
	}

	tokenString, err := generateJWTToken(userID, apicfg.JWTSecret, jwtExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate access token")
		return
	}

	refreshExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	refreshToken, err := generateJWTToken(userID, apicfg.RefreshSecret, refreshExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
		return
	}

	err = apicfg.DB.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
		ID:                    uuid.New().String(),
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             time.Now().UTC(),
		AccessTokenExpiresAt:  jwtExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
		UserID:                user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create new refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  jwtExpiresAt,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  refreshExpiresAt,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	userResp := map[string]string{
		"message": "User created successfully",
	}

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

func isValidUserName(name string) bool {
	var usernameRegex = `^[a-zA-Z0-9]+([-._]?[a-zA-Z0-9]+)*$`

	re := regexp.MustCompile(usernameRegex)

	return len(name) >= 3 && len(name) <= 30 && re.MatchString(name)
}
