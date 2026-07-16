package handler

import (
	"errors"
	"net/http"

	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func HandleRefreshToken(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type tokenResponsePayload struct {
		AccessToken string `json:"access_token"`
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			response.RespondWithError(w, http.StatusBadRequest, "cookie not found", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "server error", err)
		}
		return
	}
	refreshToken := auth.HashRefreshToken(cookie.Value)
	dbUser, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized,
			"failed to fetch refresh token from database", err)
		return
	}
	newAccess, err := auth.MakeJWT(
		dbUser.ID,
		string(cfg.AuthConfig.AccessTokenSecret),
		cfg.AuthConfig.AccessTokenTTL)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized,
			"failed to create a new access token", err)
		return
	}
	payload := tokenResponsePayload{
		AccessToken: newAccess,
	}
	response.RespondWithJSON(w, http.StatusOK, payload)
}
