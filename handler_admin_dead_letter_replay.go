package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerReplayDeadLetter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deliveryIDstr := r.PathValue("id")
	deliveryID, err := uuid.Parse(deliveryIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid delivery ID", err)
		return
	}

	err = cfg.db.ReplayDeadLetterDelivery(r.Context(), deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Delivery not found or not in dead_letter status", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to replay dead letter delivery", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":     "Delivery successfully moved back to pending queue",
		"delivery_id": deliveryID.String(),
	})
}
