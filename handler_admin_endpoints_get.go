package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	endpoints, err := cfg.db.GetAllEndpoints(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to list endpoints", err)
		return
	}

	responseEndpoints := []Endpoint{}
	for _, endpoint := range endpoints {
		responseEndpoints = append(responseEndpoints, Endpoint{
			Name:      endpoint.Name,
			IsActive:  endpoint.IsActive,
			CreatedAt: endpoint.CreatedAt,
			UserID:    endpoint.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, responseEndpoints)
}
