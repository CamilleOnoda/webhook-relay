package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	events, err := cfg.db.GetAllEvents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
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

	respondWithJSON(w, http.StatusOK, responseEvents)
}
