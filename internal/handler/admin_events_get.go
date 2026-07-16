package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
)

func HandleGetEvents(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	events, err := cfg.DB.GetAllEvents(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"failed to list events", err)
		return
	}

	responseEvents := []Event{}
	for _, event := range events {
		var eventType *string
		if event.EventType.Valid {
			eventType = &event.EventType.String
		}
		responseEvents = append(responseEvents, Event{
			EndpointName: event.EndpointName,
			EventType:    eventType,
			ReceivedAt:   event.ReceivedAt,
			UserName:     event.UserName,
		})
	}

	response.RespondWithJSON(w, http.StatusOK, responseEvents)
}
