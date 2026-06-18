package main

import (
	"errors"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type response struct {
		AccessToken string `json:"access_token"`
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			respondWithError(w, http.StatusBadRequest, "cookie not found", err)
		default:
			respondWithError(w, http.StatusInternalServerError, "server error", err)
		}
		return
	}
	refreshToken := auth.HashRefreshToken(cookie.Value)
	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"failed to fetch refresh token from database", err)
		return
	}
	newAccess, err := auth.MakeJWT(
		dbUser.ID,
		string(cfg.authConfig.AccessTokenSecret),
		cfg.authConfig.AccessTokenTTL)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"failed to create a new access token", err)
		return
	}
	tokenResponse := response{
		AccessToken: newAccess,
	}
	respondWithJSON(w, http.StatusOK, tokenResponse)
}
