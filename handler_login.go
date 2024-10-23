package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (apicfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		ApiKey string `json:"apikey"`
	}

	params := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	user, err := apicfg.DB.GetUser(r.Context(), params.ApiKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	jwtExpiresAt := time.Now().Add(15 * time.Minute).Unix()
	jwtExpiresAtTime := time.Unix(jwtExpiresAt, 0)
	tokenString, err := generateJWTToken(user.ID, user.Name, apicfg.JWTSecret, jwtExpiresAtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate token")
		return
	}

	response := map[string]interface{}{
		"token": tokenString,
	}

	respondWithJSON(w, http.StatusOK, response)

}
