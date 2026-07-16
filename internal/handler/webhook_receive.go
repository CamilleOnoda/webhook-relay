package handler

import (
	"database/sql"
	"encoding/json"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"io"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

// Receive incoming webhooks for a specific endpoint and forward them to the target URL.
// Endpoint: POST /webhooks/{id}

func HandleReceiveWebhook(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodPost {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	endpointID := r.PathValue("id")
	id, err := uuid.Parse(endpointID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"Invalid endpoint ID format", err)
		return
	}

	endpoint, err := cfg.DB.GetEndpointByID(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Endpoint not found", err)
		return
	}

	rawStream := http.MaxBytesReader(w, r.Body, 1024*1024) // limit incoming payload to max 1MB
	defer rawStream.Close()

	eventPayload, err := io.ReadAll(rawStream)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"Failed to read request body", err)
		return
	}

	marshaledHeaders, err := json.Marshal(r.Header)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to marshal headers", err)
		return
	}

	eventType := r.Header.Get("X-Event-Type")
	event, err := cfg.DB.CreateEvent(r.Context(), database.CreateEventParams{
		EndpointID: id,
		EventType:  sql.NullString{String: eventType, Valid: eventType != ""},
		Payload:    eventPayload,
		Headers:    marshaledHeaders,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to store event", err)
		return
	}

	delivery, err := cfg.DB.CreateDelivery(r.Context(), database.CreateDeliveryParams{
		EventID:   event.ID,
		TargetUrl: endpoint.TargetUrl,
		Status:    "pending",
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create delivery", err)
		return
	}

	response.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"message":     "Event received and delivery queued",
		"delivery_id": delivery.ID.String(),
	})

}
