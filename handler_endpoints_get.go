package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

func (cfg *apiConfig) handlerListEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var endpoints []database.WebhookEndpoint
	endpoints, err := cfg.db.ListEndpoints(r.Context())
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
		})
	}

	respondWithJSON(w, http.StatusOK, responseEndpoints)
}
