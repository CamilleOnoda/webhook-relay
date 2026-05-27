package main

import (
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID     `json:"id"`
	EndpointName string        `json:"endpoint_name"`
	EventType    *string       `json:"event_type"`
	ReceivedAt   time.Time     `json:"received_at"`
	UserID       uuid.NullUUID `json:"user_id"`
}

// List all webhook events for the authenticated user.
// This is a read-only operation that does not modify any data.
func (cfg *apiConfig) handlerListEventsByUser(w http.ResponseWriter, r *http.Request) {
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

	var events []database.ListEventsByUserRow
	events, err := cfg.db.ListEventsByUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to list events", err)
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
			UserID:       userID,
		})
	}

	respondWithJSON(w, http.StatusOK, responseEvents)
}
