package main

import (
	"errors"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cfg.authConfig.RefreshCookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			w.WriteHeader(http.StatusNoContent)
		default:
			respondWithError(w, http.StatusInternalServerError, "server error", err)
		}
		return
	}

	refreshToken := auth.HashRefreshToken(cookie.Value)
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "server error", err)
		return
	}

	cfg.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.authConfig.RefreshCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.authConfig.CookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
