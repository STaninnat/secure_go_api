package main

import (
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/google/uuid"
)

func (apicfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request, user database.User) {
	newTokenExpiredAt := time.Now().UTC().Add(-24 * time.Hour).Unix()
	newTokenExpiredAtTime := time.Unix(newTokenExpiredAt, 0)

	newExpiredToken := "expired-" + uuid.New().String()[:28]

	err := apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		AccessTokenExpiresAt:  newTokenExpiredAtTime,
		RefreshToken:          newExpiredToken,
		RefreshTokenExpiresAt: newTokenExpiredAtTime,
		UserID:                user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  newTokenExpiredAtTime,
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  newTokenExpiredAtTime,
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	resp := map[string]interface{}{
		"message": "logged out sucessfully",
	}

	respondWithJSON(w, http.StatusOK, resp)
}
