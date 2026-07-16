package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"

	"github.com/google/uuid"
)

func HandleDeleteUserByID(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodDelete {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	userID := r.PathValue("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid UUID format", err)
		return
	}

	_, err = cfg.DB.GetUserByID(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	if err := cfg.DB.DeleteUserByID(r.Context(), id); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to delete users", err)
		return
	}

	response.RespondWithJSON(w, http.StatusNoContent, nil)
}
