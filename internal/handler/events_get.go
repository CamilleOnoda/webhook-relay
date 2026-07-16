package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID     `json:"id"`
	EndpointName string        `json:"endpoint_name"`
	EventType    *string       `json:"event_type"`
	ReceivedAt   time.Time     `json:"received_at"`
	UserID       uuid.NullUUID `json:"user_id"`
	UserName     string        `json:"user_name"`
}

// List all webhook events for the authenticated user.
// This is a read-only operation that does not modify any data.

func HandleListEventsByUser(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	userIDFromToken, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		response.RespondWithError(w, http.StatusUnauthorized,
			"Invalid user ID in token", nil)
		return
	}
	userID := uuid.NullUUID{
		UUID:  userIDFromToken,
		Valid: true,
	}

	events, err := cfg.DB.ListEventsByUser(r.Context(), userID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
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

	response.RespondWithJSON(w, http.StatusOK, responseEvents)
}
