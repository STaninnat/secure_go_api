package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apicfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	defer r.Body.Close()
	params := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	user, err := apicfg.DB.GetUserByName(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if user.ApiKeyExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "apikey expired")
		return
	}

	jwtExpiresAt := time.Now().Add(15 * time.Minute).Unix()
	jwtExpiresAtTime := time.Unix(jwtExpiresAt, 0)
	tokenString, err := generateJWTToken(user.ID, apicfg.JWTSecret, jwtExpiresAtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate token")
		return
	}

	existingToken, err := apicfg.DB.GetRfKeyByUserID(r.Context(), user.ID)
	if err != nil && existingToken.RefreshTokenExpiresAt.After(time.Now()) {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  jwtExpiresAtTime,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    existingToken.RefreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  existingToken.AccessTokenExpiresAt,
		})

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Login Successful"})
		return
	}

	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour).Unix()
	refreshExpiresAtTime := time.Unix(refreshExpiresAt, 0)
	refreshToken, err := generateJWTToken(user.ID, apicfg.RefreshSecret, refreshExpiresAtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
		return
	}

	_, err = apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAtTime,
		UserID:                user.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = apicfg.DB.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
				ID:                    uuid.New(),
				CreatedAt:             time.Now().UTC(),
				AccessTokenExpiresAt:  jwtExpiresAtTime,
				RefreshToken:          refreshToken,
				RefreshTokenExpiresAt: refreshExpiresAtTime,
				UserID:                user.ID,
			})
			if err != nil {
				// log.Printf("error creating new refresh token: %v\n", err)
				respondWithError(w, http.StatusInternalServerError, "failed to create new refresh token")
				return
			}
		} else {
			// log.Printf("error updating refresh token: %v\n", err)
			respondWithError(w, http.StatusInternalServerError, "failed to update new refresh token")
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  jwtExpiresAtTime,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  refreshExpiresAtTime,
	})

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Login Successful"})
}
