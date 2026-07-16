package handler

import (
	"database/sql"
	"errors"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"

	"github.com/google/uuid"
)

func HandleReplayDeadLetter(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deliveryIDstr := r.PathValue("id")
	deliveryID, err := uuid.Parse(deliveryIDstr)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid delivery ID", err)
		return
	}

	err = cfg.DB.ReplayDeadLetterDelivery(r.Context(), deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithError(w, http.StatusNotFound, "Delivery not found or not in dead_letter status", err)
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to replay dead letter delivery", err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message":     "Delivery successfully moved back to pending queue",
		"delivery_id": deliveryID.String(),
	})
}
