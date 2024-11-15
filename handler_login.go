package main

import (
	"database/sql"
	"encoding/json"
	"log"
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
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusBadRequest, "username not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "error retrieving user")
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect password")
		return
	}

	if user.ApiKeyExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "apikey expired")
		return
	}

	jwtExpiresAt := time.Now().UTC().Add(1 * time.Hour).Unix()
	jwtExpiresAtTime := time.Unix(jwtExpiresAt, 0)

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "invalid user ID")
		return
	}

	tokenString, err := generateJWTToken(userID, apicfg.JWTSecret, jwtExpiresAtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate access token")
		return
	}

	tx, err := apicfg.DBConn.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to start transaction")
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Failed to commit transaction: %v", err)
				respondWithError(w, http.StatusInternalServerError, "failed to commit transaction")
				return
			}
		}
	}()

	queriesTx := database.New(tx)

	refreshExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour).Unix()
	refreshExpiresAtTime := time.Unix(refreshExpiresAt, 0)

	refreshToken, err := generateJWTToken(userID, apicfg.RefreshSecret, refreshExpiresAtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
		return
	}

	err = queriesTx.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		UpdatedAt:             time.Now().UTC(),
		AccessTokenExpiresAt:  jwtExpiresAtTime,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAtTime,
		UserID:                user.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			err = queriesTx.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
				ID:                    uuid.New().String(),
				CreatedAt:             time.Now().UTC(),
				UpdatedAt:             time.Now().UTC(),
				AccessTokenExpiresAt:  jwtExpiresAtTime,
				RefreshToken:          refreshToken,
				RefreshTokenExpiresAt: refreshExpiresAtTime,
				UserID:                user.ID,
			})
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "failed to create new refresh token")
				return
			}
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to update new refresh token")
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Expires:  jwtExpiresAtTime,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExpiresAtTime,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	userResp := map[string]string{
		"message": "Login Successful",
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
