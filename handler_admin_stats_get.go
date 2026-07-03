package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	type stats struct {
		Users                    int64 `json:"users"`
		EventsReceived           int64 `json:"events_received"`
		SucessfulDeliveries      int64 `json:"successful_deliveries"`
		DeadLetter               int64 `json:"dead_letter"`
		RetryScheduledDeliveries int64 `json:"retry_scheduled_deliveries"`
	}

	adminStats, err := cfg.db.GetAdminStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to get stats", err)
		return
	}

	responseStats := stats{
		Users:                    adminStats.Users,
		EventsReceived:           adminStats.EventsReceived,
		SucessfulDeliveries:      adminStats.SuccessfulDeliveries,
		DeadLetter:               adminStats.DeadLetter,
		RetryScheduledDeliveries: adminStats.RetryScheduledDeliveries,
	}

	respondWithJSON(w, http.StatusOK, responseStats)
}
