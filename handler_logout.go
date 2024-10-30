package main

import (
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
)

func (apicfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	existingToken, err := apicfg.DB.GetUserByRfKey(r.Context(), user.ID.String())
	if err != nil || existingToken.RefreshTokenExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "invalid or mismatched refresh token")
		return
	}

	newRefreshTokenExpiredAt := time.Now().Add(-time.Hour).Unix()
	newRefreshTokenExpiredAtTime := time.Unix(newRefreshTokenExpiredAt, 0)
	_, err = apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		RefreshToken:          "",
		RefreshTokenExpiresAt: newRefreshTokenExpiredAtTime,
		UserID:                existingToken.UserID,
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
