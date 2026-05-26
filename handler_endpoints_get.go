package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Get all endpoints for the authenticated user.
// This is a read-only operation that does not modify any data.
func (cfg *apiConfig) handlerGetEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
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

	var endpoints []database.WebhookEndpoint
	endpoints, err = cfg.db.GetEndpointsByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list endpoints", err)
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
