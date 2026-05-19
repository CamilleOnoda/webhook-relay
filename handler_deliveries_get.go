package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

func (cfg *apiConfig) handlerListDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var deliveries []database.Delivery
	deliveries, err := cfg.db.ListDeliveries(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list deliveries", err)
		return
	}
	respondWithJSON(w, http.StatusOK, deliveries)
}
