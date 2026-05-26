package main

import "net/http"

func (cfg *apiConfig) handlerDeleteUsers(w http.ResponseWriter, r *http.Request) {
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

	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to delete users", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
