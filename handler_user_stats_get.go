package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetUserStats(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerGetUserStats called")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized,
			"Invalid user ID in token", nil)
		return
	}
	userID := uuid.NullUUID{
		UUID:  userIDFromToken,
		Valid: true,
	}

	type stats struct {
		EndpointCount            int64 `json:"endpoint_count"`
		EventsReceived           int64 `json:"events_received"`
		SuccessfulDeliveries     int64 `json:"successful_delivery_count"`
		DeadLetter               int64 `json:"failed_delivery_count"`
		RetryScheduledDeliveries int64 `json:"retry_scheduled_delivery_count"`
	}

	userStats, err := cfg.db.GetUserStatsByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to get stats", err)
		return
	}

	responseStats := stats{
		EndpointCount:            userStats.EndpointCount,
		EventsReceived:           userStats.EventCount,
		SuccessfulDeliveries:     userStats.SuccessfulDeliveryCount,
		DeadLetter:               userStats.FailedDeliveryCount,
		RetryScheduledDeliveries: userStats.RetryScheduledDeliveryCount,
	}

	respondWithJSON(w, http.StatusOK, responseStats)
}
