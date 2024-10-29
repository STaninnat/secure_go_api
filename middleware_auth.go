package main

import (
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/auth"
	"github.com/STaninnat/capstone_project/internal/database"
)

type authhandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg apiConfig) middlewareAuth(handler authhandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "couldn't find token")
			return
		}

		claims, err := validateJWTToken(tokenString, apicfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := apicfg.DB.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "couldn't get user")
			return
		}

		if user.ApiKeyExpiresAt.Before(time.Now().UTC()) {
			respondWithError(w, http.StatusUnauthorized, "API key expired")
			return
		}

		handler(w, r, user)
	}
}
