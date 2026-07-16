package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
)

func HandleDeleteAllEndpoints(cfg *Config, w http.ResponseWriter, r *http.Request) {
	if cfg.Environment != "dev" {
		response.RespondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodDelete {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	if err := cfg.DB.DeleteAllEndpoints(r.Context()); err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
