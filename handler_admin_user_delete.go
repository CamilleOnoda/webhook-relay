package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	userID := r.PathValue("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	if err := cfg.db.DeleteUserByID(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to delete users", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
