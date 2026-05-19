package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteEndpointByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	endpointID := r.PathValue("id")
	id, err := uuid.Parse(endpointID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}
	_, err = cfg.db.GetEndpointByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("endpoint not found"))
		return
	}
	err = cfg.db.DeleteEndpointByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete endpoint", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
