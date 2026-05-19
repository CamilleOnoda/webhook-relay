package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetEndpointByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	endpointID := r.PathValue("id")
	id, err := uuid.Parse(endpointID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	endpoint, err := cfg.db.GetEndpointByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "endpoint not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, endpoint)
}
