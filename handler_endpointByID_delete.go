package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Delete a specific endpoint by ID for the authenticated user.
// Endpoint: DELETE /endpoints/{id}
func (cfg *apiConfig) handlerDeleteEndpointByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
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
		respondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	_, err = cfg.db.GetEndpointByIDAndUserID(r.Context(), database.GetEndpointByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "endpoint not found", err)
		return
	}

	err = cfg.db.DeleteEndpointByIDAndUserID(r.Context(), database.DeleteEndpointByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete endpoint", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
