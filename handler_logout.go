package main

import (
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
)

func (apicfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	newRefreshTokenExpiredAt := time.Now().UTC().Add(-24 * time.Hour).Unix()
	newRefreshTokenExpiredAtTime := time.Unix(newRefreshTokenExpiredAt, 0)

	_, err := apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		RefreshToken:          "",
		RefreshTokenExpiresAt: newRefreshTokenExpiredAtTime,
		UserID:                user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().UTC().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	resp := map[string]interface{}{
		"message": "logged out sucessfully",
	}

	respondWithJSON(w, http.StatusOK, resp)
}
