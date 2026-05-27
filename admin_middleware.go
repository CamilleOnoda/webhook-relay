package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
		if !ok {
			respondWithError(w, http.StatusUnauthorized,
				"Invalid user ID in token", nil)
			return
		}
		isAdmin, err := cfg.db.IsUserAdmin(r.Context(), userIDFromToken)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				"Failed to check admin status", err)
			return
		}
		if isAdmin != true {
			respondWithError(w, http.StatusForbidden,
				"Admin access required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
