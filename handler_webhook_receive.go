package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/CamilleOnoda/webhook-relay.git/internal/service"
	"github.com/google/uuid"
)

// Receive incoming webhooks for a specific endpoint and forward them to the target URL.
// Endpoint: POST /webhooks/{id}
func (cfg *apiConfig) handlerReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	endpointID := r.PathValue("id")
	id, err := uuid.Parse(endpointID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Invalid endpoint ID format", err)
		return
	}

	endpoint, err := cfg.db.GetEndpointByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Endpoint not found", err)
		return
	}

	rawStream := http.MaxBytesReader(w, r.Body, 1024*1024) // limit incoming payload to max 1MB
	defer rawStream.Close()

	eventPayload, err := io.ReadAll(rawStream)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Failed to read request body", err)
		return
	}

	marshaledHeaders, err := json.Marshal(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to marshal headers", err)
		return
	}

	eventType := r.Header.Get("X-Event-Type")
	event, err := cfg.db.CreateEvent(r.Context(), database.CreateEventParams{
		EndpointID: id,
		EventType:  sql.NullString{String: eventType, Valid: eventType != ""},
		Payload:    eventPayload,
		Headers:    marshaledHeaders,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to store event", err)
		return
	}

	delivery, err := cfg.db.CreateDelivery(r.Context(), database.CreateDeliveryParams{
		EventID:   event.ID,
		TargetUrl: endpoint.TargetUrl,
		Status:    "pending",
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to create delivery", err)
		return
	}

	attemptResult, err := service.AttemptDelivery(r.Context(), event, endpoint.TargetUrl)
	if err != nil {
		log.Printf("failed to forward event_id=%s target_url=%s error=%v",
			event.ID, endpoint.TargetUrl, err)
	}

	status := "success"
	if err != nil || attemptResult.ErrorMessage.Valid {
		status = "retry_scheduled"
	}

	if err := cfg.db.UpdateDelivery(r.Context(), database.UpdateDeliveryParams{
		ID:                 delivery.ID,
		Status:             status,
		StatusCode:         attemptResult.StatusCode,
		ResponseBody:       attemptResult.ResponseBody,
		ErrorMessage:       attemptResult.ErrorMessage,
		DeliveryDurationMs: attemptResult.DeliveryDurationMs,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to update delivery record", err)
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{
		"message": "Event received and delivery attempted",
	})

}
