package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteEndpointByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"failed to get token", err)
		return
	}
	userIDFromToken, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"failed to validate token", err)
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
