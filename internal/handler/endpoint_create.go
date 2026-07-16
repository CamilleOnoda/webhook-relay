package handler

import (
	"database/sql"
	"encoding/json"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"net/url"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type Endpoint struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	TargetUrl    string        `json:"target_url"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	GeneratedURL string        `json:"generated_url"`
	Description  *string       `json:"description,omitempty"`
	UserID       uuid.NullUUID `json:"user_id"`
	UserName     string        `json:"user_name"`
}

// Create a new webhook endpoint for the authenticated user.
// This is a write operation that modifies data.

func HandleCreateEndpoint(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		TargetUrl   string  `json:"target_url"`
		Description *string `json:"description,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"Invalid request payload", err)
		return
	}

	if req.Name == "" || req.TargetUrl == "" {
		response.RespondWithError(w, http.StatusBadRequest,
			"Name and TargetUrl are required", nil)
		return
	}

	validURL, err := url.Parse(req.TargetUrl)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"Invalid target url format", err)
		return
	}
	if validURL.Scheme != "https" {
		response.RespondWithError(w, http.StatusBadRequest,
			"Target url must use HTTPS scheme", nil)
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

	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	} else {
		description = sql.NullString{Valid: false}
	}

	dbEndpoint, err := cfg.DB.CreateEndpoint(r.Context(), database.CreateEndpointParams{
		Name:        req.Name,
		TargetUrl:   validURL.String(),
		Description: description,
		UserID:      userID,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to create endpoint in database", err)
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
		UserID:       userID,
	}

	generatedURL := cfg.BaseURL + "/webhooks" + "/" + dbEndpoint.ID.String()
	responseEndpoint.GeneratedURL = generatedURL

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseEndpoint)

}
