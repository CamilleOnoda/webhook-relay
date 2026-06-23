package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Get a specific endpoint by its ID for the authenticated user.
// This is a read-only operation that does not modify any data.
func (cfg *apiConfig) handlerGetEndpointByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized,
			"Invalid user ID in token", nil)
		return
	}
	userID := uuid.NullUUID{
		UUID:  userIDFromToken,
		Valid: true,
	}

	endpointID := r.PathValue("id")
	id, err := uuid.Parse(endpointID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Invalid uuid format", err)
		return
	}

	endpoint, err := cfg.db.GetEndpointByIDAndUserID(r.Context(), database.GetEndpointByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "endpoint not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, endpoint)
}
