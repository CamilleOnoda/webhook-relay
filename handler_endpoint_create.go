package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Endpoint struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	TargetUrl string    `json:"target_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req struct {
		Name      string `json:"name"`
		TargetUrl string `json:"target_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Name == "" || req.TargetUrl == "" {
		respondWithError(w, http.StatusBadRequest, "Name and TargetUrl are required", nil)
		return
	}

	validURL, err := url.ParseRequestURI(req.TargetUrl)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid target url format", err)
		return
	}
	if validURL.Scheme != "https" {
		respondWithError(w, http.StatusBadRequest, "Target url must use https scheme", nil)
	}

	dbEndpoint, err := cfg.db.CreateEndpoint(r.Context(), database.CreateEndpointParams{
		Name:      req.Name,
		TargetUrl: validURL.String(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create endpoint in database", err)
		return
	}

	responseEndpoint := Endpoint{
		ID:        dbEndpoint.ID,
		Name:      dbEndpoint.Name,
		TargetUrl: dbEndpoint.TargetUrl,
		IsActive:  dbEndpoint.IsActive,
		CreatedAt: dbEndpoint.CreatedAt,
		UpdatedAt: dbEndpoint.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseEndpoint)

}
