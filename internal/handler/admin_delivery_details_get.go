package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func HandleGetAdminDeliveryDetails(cfg *Config, w http.ResponseWriter, r *http.Request) {
	type deliveryDetail struct {
		ID           uuid.UUID  `json:"id"`
		EndpointName string     `json:"endpoint_name"`
		TargetURL    string     `json:"target_url"`
		Status       string     `json:"status"`
		AttemptCount int32      `json:"attempt_count"`
		StatusCode   *int32     `json:"status_code"`
		ErrorMessage *string    `json:"error_message"`
		NextRetryAt  *time.Time `json:"next_retry_at"`
		CreatedAt    time.Time  `json:"created_at"`
		AttemptedAt  *time.Time `json:"attempted_at"`
		DeliveredAt  *time.Time `json:"delivered_at"`
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	ID := r.PathValue("id")
	deliveryID, err := uuid.Parse(ID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"Invalid uuid format", err)
		return
	}

	dbDeliveryDetail, err := cfg.DB.GetAdminDeliveryDetails(r.Context(), deliveryID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to get admin details from database", err)
		return
	}

	var statusCode *int32
	if dbDeliveryDetail.StatusCode.Valid {
		statusCode = &dbDeliveryDetail.StatusCode.Int32
	}
	var errorMessage *string
	if dbDeliveryDetail.ErrorMessage.Valid {
		errorMessage = &dbDeliveryDetail.ErrorMessage.String
	}
	var nextRetryAt *time.Time
	if dbDeliveryDetail.NextRetryAt.Valid {
		nextRetryAt = &dbDeliveryDetail.NextRetryAt.Time
	}
	var attemptedAt *time.Time
	if dbDeliveryDetail.AttemptedAt.Valid {
		attemptedAt = &dbDeliveryDetail.AttemptedAt.Time
	}
	var deliveredAt *time.Time
	if dbDeliveryDetail.DeliveredAt.Valid {
		deliveredAt = &dbDeliveryDetail.DeliveredAt.Time
	}

	responseDeliveryDetail := deliveryDetail{
		ID:           dbDeliveryDetail.ID,
		EndpointName: dbDeliveryDetail.EndpointName,
		TargetURL:    dbDeliveryDetail.TargetUrl,
		Status:       dbDeliveryDetail.Status,
		AttemptCount: dbDeliveryDetail.AttemptCount,
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
		NextRetryAt:  nextRetryAt,
		CreatedAt:    dbDeliveryDetail.CreatedAt,
		AttemptedAt:  attemptedAt,
		DeliveredAt:  deliveredAt,
	}

	response.RespondWithJSON(w, http.StatusOK, responseDeliveryDetail)
}
