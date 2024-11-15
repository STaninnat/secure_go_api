package main

import (
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type authhandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg apiConfig) middlewareAuth(handler authhandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("access_token")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
			return
		}

		claims, err := validateJWTToken(tokenString.Value, apicfg.JWTSecret)
		if err != nil {
			log.Printf("Token validation error: %v\n", err)
			if err == jwt.ErrTokenExpired {
				respondWithError(w, http.StatusUnauthorized, "token expired")
				return
			}
			respondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := apicfg.DB.GetUserByID(r.Context(), claims.UserID.String())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get user")
			return
		}

		if isAPIKeyExpired(user) {
			respondWithError(w, http.StatusUnauthorized, "api key expired")
			return
		}

		handler(w, r, user)
	}
}

func isAPIKeyExpired(user database.User) bool {
	apiKeyExpiresAt, err := time.Parse(time.RFC3339, user.ApiKeyExpiresAt)
	if err != nil {
		log.Printf("Error parsing API key expiration time: %v", err)
		return true
	}

	return apiKeyExpiresAt.Before(time.Now().UTC())
}
