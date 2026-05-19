package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

func (cfg *apiConfig) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var events []database.WebhookEvent
	events, err := cfg.db.ListEvents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list events", err)
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}
