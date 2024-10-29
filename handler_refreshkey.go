package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
)

func (apicfg *apiConfig) handlerRefreshKey(w http.ResponseWriter, r *http.Request) {
	type refreshParams struct {
		RefreshToken string `json:"refresh_token"`
	}

	params := refreshParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
		return
	}

	user, err := apicfg.DB.GetUserByRfKey(r.Context(), params.RefreshToken)
	if err != nil || user.RefreshTokenExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	_, newHashedApiKey, err := generateAndHashAPIKey()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate new apikey")
		return
	}
	newApiKeyExpiresAt := time.Now().UTC().Add(365 * 24 * time.Hour)

	newAccessTokenExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	newAccessToken, err := generateJWTToken(user.UserID, apicfg.JWTSecret, newAccessTokenExpiresAt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate new access token")
		return
	}

	err = apicfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ApiKey:          newHashedApiKey,
		ApiKeyExpiresAt: newApiKeyExpiresAt,
		ID:              user.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update apikey")
		return
	}

	newRefreshTokenExpiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	_, err = apicfg.DB.UpdateUserRfKey(r.Context(), database.UpdateUserRfKeyParams{
		RefreshToken:          params.RefreshToken,
		RefreshTokenExpiresAt: newRefreshTokenExpiresAt,
		UserID:                user.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update refresh token")
		return
	}

	userResp := map[string]interface{}{
		"access_token":             newAccessToken,
		"access_token_expires_at":  newAccessTokenExpiresAt,
		"refresh_token":            params.RefreshToken,
		"refresh_token_expires_at": newRefreshTokenExpiresAt,
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
