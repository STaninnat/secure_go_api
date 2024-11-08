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

	jwtExpiresAt := time.Now().UTC().Add(15 * time.Minute).Unix()
	jwtExpiresAtTime := time.Unix(jwtExpiresAt, 0)
	tokenString, err := generateJWTToken(user.ID, apicfg.JWTSecret, jwtExpiresAtTime)
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
		}
	}()

	queriesTx := database.New(tx)

	existingToken, err := queriesTx.GetRfKeyByUserID(r.Context(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve refresh token")
		return
	}
	if existingToken.RefreshToken == "" || existingToken.RefreshTokenExpiresAt.Before(time.Now().UTC()) {
		refreshExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour).Unix()
		refreshExpiresAtTime := time.Unix(refreshExpiresAt, 0)
		refreshToken, err := generateJWTToken(user.ID, apicfg.RefreshSecret, refreshExpiresAtTime)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
			return
		}

		updatedToken, err := queriesTx.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshExpiresAtTime,
			UserID:                user.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = queriesTx.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
					ID:                    uuid.New(),
					CreatedAt:             time.Now().UTC(),
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
		log.Printf("Debug: Successfully updated refresh token for user ID %v: %+v", user.ID, updatedToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  refreshExpiresAtTime,
		})
	} else {
		refreshToken := existingToken.RefreshToken
		refreshExpiresAtTime := existingToken.RefreshTokenExpiresAt
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  refreshExpiresAtTime,
		})
	}
	log.Printf("Debug: Generated access token for user ID %v: %s", user.ID, tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  jwtExpiresAtTime,
	})

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Login Successful"})
}
