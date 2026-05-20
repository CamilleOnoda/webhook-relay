package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Delivery struct {
	ID                 uuid.UUID     `json:"id"`
	TargetUrl          string        `json:"target_url"`
	Status             string        `json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	DeliveryDurationMs sql.NullInt32 `json:"delivery_duration_ms"`
}

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

	responseDeliveries := []Delivery{}

	for _, delivery := range deliveries {
		responseDeliveries = append(responseDeliveries, Delivery{
			ID:                 delivery.ID,
			TargetUrl:          delivery.TargetUrl,
			Status:             delivery.Status,
			CreatedAt:          delivery.CreatedAt,
			DeliveryDurationMs: delivery.DeliveryDurationMs,
		})
	}

	respondWithJSON(w, http.StatusOK, responseDeliveries)
}
