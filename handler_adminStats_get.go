package main

import (
	"net/http"
)

type Stats struct {
	Users               int64 `json:"users"`
	EventsReceived      int64 `json:"events_received"`
	SucessfulDeliveries int64 `json:"successful_deliveries"`
	FailedDeliveries    int64 `json:"failed_deliveries"`
}

func (cfg *apiConfig) handlerGetAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	stats, err := cfg.db.GetAdminStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to get stats", err)
		return
	}

	responseStats := Stats{
		Users:               stats.Users,
		EventsReceived:      stats.EventsReceived,
		SucessfulDeliveries: stats.SuccessfulDeliveries,
		FailedDeliveries:    stats.FailedDeliveries,
	}

	respondWithJSON(w, http.StatusOK, responseStats)
}
