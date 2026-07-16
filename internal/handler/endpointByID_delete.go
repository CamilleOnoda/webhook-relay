package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Delete a specific endpoint by ID for the authenticated user.
// Endpoint: DELETE /endpoints/{id}

func HandleDeleteEndpointByID(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		response.RespondWithError(w, http.StatusUnauthorized,
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
		response.RespondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	_, err = cfg.DB.GetEndpointByID(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "endpoint not found", err)
		return
	}

	err = cfg.DB.DeleteEndpointByIDAndUserID(r.Context(), database.DeleteEndpointByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to delete endpoint", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
