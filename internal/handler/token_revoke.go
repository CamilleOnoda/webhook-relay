package handler

import (
	"errors"
	"net/http"

	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func HandleRevoke(cfg *Config, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cfg.AuthConfig.RefreshCookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			w.WriteHeader(http.StatusNoContent)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "server error", err)
		}
		return
	}

	refreshToken := auth.HashRefreshToken(cookie.Value)
	_, err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "server error", err)
		return
	}

	cfg.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *Config) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.AuthConfig.RefreshCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.AuthConfig.CookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
