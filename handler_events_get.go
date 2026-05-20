package main

import (
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `json:"id"`
	EndpointName string    `json:"endpoint_name"`
	EventType    *string   `json:"event_type"`
	ReceivedAt   time.Time `json:"received_at"`
}

func (cfg *apiConfig) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	var events []database.ListEventsRow
	events, err := cfg.db.ListEvents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to list events", err)
		return
	}

	responseEvents := []Event{}

	for _, event := range events {
		var eventType *string
		if event.EventType.Valid {
			value := event.EventType.String
			eventType = &value
		}
		responseEvents = append(responseEvents, Event{
			ID:           event.ID,
			EndpointName: event.EndpointName,
			EventType:    eventType,
			ReceivedAt:   event.ReceivedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, responseEvents)
}
