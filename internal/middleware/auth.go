package middleware

import (
	"context"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
)

func Auth(accessTokenSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "failed to get token", err)
			return
		}
		userID, err := auth.ValidateJWT(token, accessTokenSecret)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "failed to validate token", err)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
