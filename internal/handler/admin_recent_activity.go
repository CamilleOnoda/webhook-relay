package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func HandleGetRecentActivity(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	type recentActivity struct {
		EventID                  uuid.UUID `json:"event_id"`
		DeliveryID               uuid.UUID `json:"delivery_id"`
		ReceivedAt               time.Time `json:"received_at"`
		EndpointName             string    `json:"endpoint_name"`
		UserName                 string    `json:"user_name"`
		EventType                *string   `json:"event_type"`
		LatestDeliveryStatus     string    `json:"latest_delivery_status"`
		LatestDeliveryStatusCode *int32    `json:"latest_delivery_status_code"`
		AttemptCount             int32     `json:"attempt_count"`
	}

	activities, err := cfg.DB.GetAdminRecentActivity(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
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
			UserName:                 recentAct.UserName,
			EventType:                eventType,
			AttemptCount:             recentAct.AttemptCount,
			LatestDeliveryStatus:     recentAct.LatestDeliveryStatus,
			LatestDeliveryStatusCode: statusCode,
		})
	}
	response.RespondWithJSON(w, http.StatusOK, responseActivity)
}
