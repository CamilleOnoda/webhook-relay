package main

import (
	"context"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized,
				"failed to get token", err)
			return
		}
		userID, err := auth.ValidateJWT(token, string(cfg.authConfig.AccessTokenSecret))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized,
				"failed to validate token", err)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
