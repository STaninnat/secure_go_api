package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
)

func (apicfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	type logoutParams struct {
		RefreshToken string `json:"refresh_token"`
	}

	defer r.Body.Close()
	param := logoutParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	userByRefreshKey, err := apicfg.DB.GetUserByRfKey(r.Context(), param.RefreshToken)
	if err != nil || userByRefreshKey.UserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "invalid or mismatched refresh token")
		return
	}

	newRefreshTokenExpiredAt := time.Now().Add(-time.Hour).Unix()
	newRefreshTokenExpiredAtTime := time.Unix(newRefreshTokenExpiredAt, 0)
	_, err = apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		RefreshToken:          "",
		RefreshTokenExpiresAt: newRefreshTokenExpiredAtTime,
		UserID:                userByRefreshKey.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't logout")
		return
	}

	userResp := map[string]interface{}{
		"message": "logged out successfully",
		"action":  "remove access token from client",
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
