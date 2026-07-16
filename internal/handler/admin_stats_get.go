package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
)

func HandleGetAdminStats(cfg *Config, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
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

	adminStats, err := cfg.DB.GetAdminStats(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
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

	response.RespondWithJSON(w, http.StatusOK, responseStats)
}
