package middleware

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"github.com/google/uuid"
)

func Admin(next http.Handler, db *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
		if !ok {
			response.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token", nil)
			return
		}
		isAdmin, err := db.IsUserAdmin(r.Context(), userIDFromToken)
		if err != nil {
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to check admin status", err)
			return
		}
		if !isAdmin {
			response.RespondWithError(w, http.StatusForbidden, "Admin access required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
