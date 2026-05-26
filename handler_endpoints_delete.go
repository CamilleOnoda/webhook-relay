package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerDeleteAllEndpoints(w http.ResponseWriter, r *http.Request) {
	if cfg.environment != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	if err := cfg.db.DeleteAllEndpoints(r.Context()); err != nil {
		respondWithError(w, http.StatusNotFound, "Not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
