package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Get all endpoints for the authenticated user.
// This is a read-only operation that does not modify any data.
func (cfg *apiConfig) handlerGetEndpoints(w http.ResponseWriter, r *http.Request) {
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

	var endpoints []database.WebhookEndpoint
	endpoints, err := cfg.db.GetEndpointsByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to list endpoints", err)
		return
	}

	responseEndpoints := []Endpoint{}

	for _, endpoint := range endpoints {
		responseEndpoints = append(responseEndpoints, Endpoint{
			ID:           endpoint.ID,
			Name:         endpoint.Name,
			TargetUrl:    endpoint.TargetUrl,
			IsActive:     endpoint.IsActive,
			CreatedAt:    endpoint.CreatedAt,
			UpdatedAt:    endpoint.UpdatedAt,
			GeneratedURL: cfg.baseURL + "/webhooks/" + endpoint.ID.String(),
			Description:  &endpoint.Description.String,
			UserID:       endpoint.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, responseEndpoints)
}
