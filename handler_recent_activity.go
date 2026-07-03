package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetUserRecentActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	type recentActivity struct {
		EventID                  uuid.UUID `json:"event_id"`
		DeliveryID               uuid.UUID `json:"delivery_id"`
		ReceivedAt               time.Time `json:"received_at"`
		EndpointName             string    `json:"endpoint_name"`
		EventType                *string   `json:"event_type"`
		LatestDeliveryStatus     string    `json:"latest_delivery_status"`
		LatestDeliveryStatusCode *int32    `json:"latest_delivery_status_code"`
		AttemptCount             int32     `json:"attempt_count"`
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

	activities, err := cfg.db.GetUserRecentActvity(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"failed to get recent activity", err)
		return
	}

	responseActivity := []recentActivity{}
	for _, recentAct := range activities {
		var eventType *string
		if recentAct.EventType.Valid {
			eventType = &recentAct.EventType.String
		}
		var statusCode *int32
		if recentAct.LatestDeliveryStatusCode.Valid {
			statusCode = &recentAct.LatestDeliveryStatusCode.Int32
		}
		responseActivity = append(responseActivity, recentActivity{
			EventID:                  recentAct.EventID,
			DeliveryID:               recentAct.DeliveryID,
			ReceivedAt:               recentAct.ReceivedAt,
			EndpointName:             recentAct.EndpointName,
			EventType:                eventType,
			AttemptCount:             recentAct.AttemptCount,
			LatestDeliveryStatus:     recentAct.LatestDeliveryStatus,
			LatestDeliveryStatusCode: statusCode,
		})
	}
	respondWithJSON(w, http.StatusOK, responseActivity)
}
