package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Endpoint struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	TargetUrl    string     `json:"target_url"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	GeneratedURL string     `json:"generated_url"`
	Description  *string    `json:"description,omitempty"`
	UserID       *uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req struct {
		Name        string     `json:"name"`
		TargetUrl   string     `json:"target_url"`
		Description *string    `json:"description,omitempty"`
		UserID      *uuid.UUID `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Name == "" || req.TargetUrl == "" {
		respondWithError(w, http.StatusBadRequest, "Name and TargetUrl are required", nil)
		return
	}

	validURL, err := url.Parse(req.TargetUrl)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid target url format", err)
		return
	}
	if validURL.Scheme != "https" {
		respondWithError(w, http.StatusBadRequest, "Target url must use HTTPS scheme", nil)
		return
	}

	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	} else {
		description = sql.NullString{Valid: false}
	}

	var userID uuid.NullUUID
	if req.UserID != nil {
		userID = uuid.NullUUID{UUID: *req.UserID, Valid: true}
	} else {
		userID = uuid.NullUUID{Valid: false}
	}

	dbEndpoint, err := cfg.db.CreateEndpoint(r.Context(), database.CreateEndpointParams{
		Name:        req.Name,
		TargetUrl:   validURL.String(),
		Description: description,
		UserID:      userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create endpoint in database", err)
		return
	}

	var responseDescription *string
	if req.Description != nil {
		responseDescription = &dbEndpoint.Description.String
	} else {
		responseDescription = nil
	}

	responseEndpoint := Endpoint{
		ID:           dbEndpoint.ID,
		Name:         dbEndpoint.Name,
		TargetUrl:    dbEndpoint.TargetUrl,
		IsActive:     dbEndpoint.IsActive,
		CreatedAt:    dbEndpoint.CreatedAt,
		UpdatedAt:    dbEndpoint.UpdatedAt,
		GeneratedURL: "",
		Description:  responseDescription,
	}

	generatedURL := cfg.baseURL + "/webhooks" + "/" + dbEndpoint.ID.String()
	responseEndpoint.GeneratedURL = generatedURL

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseEndpoint)

}
