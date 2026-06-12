package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Delivery struct {
	ID                 uuid.UUID `json:"id"`
	EndpointName       string    `json:"endpoint_name"`
	TargetUrl          string    `json:"target_url"`
	Status             string    `json:"status"`
	StatusCode         *int32    `json:"status_code"`
	CreatedAt          time.Time `json:"created_at"`
	DeliveryDurationMs *int32    `json:"delivery_duration_ms"`
	UserName           string    `json:"user_name"`
}

func (cfg *apiConfig) handlerListDeliveriesByUser(w http.ResponseWriter, r *http.Request) {
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

	deliveries, err := cfg.db.ListDeliveriesByUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to list deliveries", err)
		return
	}

	responseDeliveries := []Delivery{}

	for _, delivery := range deliveries {
		var statusCode *int32
		if delivery.StatusCode.Valid {
			value := delivery.StatusCode.Int32
			statusCode = &value
		}
		var duration *int32
		if delivery.DeliveryDurationMs.Valid {
			value := delivery.DeliveryDurationMs.Int32
			duration = &value
		}
		responseDeliveries = append(responseDeliveries, Delivery{
			ID:                 delivery.ID,
			EndpointName:       delivery.EndpointName,
			TargetUrl:          delivery.TargetUrl,
			Status:             delivery.Status,
			StatusCode:         statusCode,
			CreatedAt:          delivery.CreatedAt,
			DeliveryDurationMs: duration,
		})
	}

	respondWithJSON(w, http.StatusOK, responseDeliveries)
}
