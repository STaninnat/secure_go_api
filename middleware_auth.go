package main

import (
	"net/http"

	"github.com/STaninnat/capstone_project/internal/auth"
	"github.com/STaninnat/capstone_project/internal/database"
)

type authhandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg apiConfig) middlewareAuth(handler authhandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find apikey")
			return
		}

		user, err := apicfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get user")
			return
		}

		handler(w, r, user)
	}
}
