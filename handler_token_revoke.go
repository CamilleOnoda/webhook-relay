package main

import (
	"errors"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "error revoking token", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
