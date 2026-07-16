package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func HandleGetDeadLetters(cfg *Config, w http.ResponseWriter, r *http.Request) {
	type delivery struct {
		ID                 uuid.UUID `json:"id"`
		EndpointName       string    `json:"endpoint_name"`
		TargetUrl          string    `json:"target_url"`
		Status             string    `json:"status"`
		StatusCode         *int32    `json:"status_code"`
		CreatedAt          time.Time `json:"created_at"`
		DeliveryDurationMs *int32    `json:"delivery_duration_ms"`
		UserName           string    `json:"user_name"`
		ResponseBody       string    `json:"response_body"`
		ErrorMessage       string    `json:"error_message"`
	}

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deadLetters, err := cfg.DB.ListDeadLetterDeliveries(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to list dead letter deliveries", err)
		return
	}
	responseDeadLetters := []delivery{}
	for _, deadLetter := range deadLetters {
		var statusCode *int32
		if deadLetter.StatusCode.Valid {
			value := deadLetter.StatusCode.Int32
			statusCode = &value
		}
		var deliveryDuration *int32
		if deadLetter.DeliveryDurationMs.Valid {
			value := deadLetter.DeliveryDurationMs.Int32
			deliveryDuration = &value
		}
		responseDeadLetters = append(responseDeadLetters, delivery{
			ID:                 deadLetter.ID,
			EndpointName:       deadLetter.EndpointName,
			TargetUrl:          deadLetter.TargetUrl,
			Status:             deadLetter.Status,
			StatusCode:         statusCode,
			CreatedAt:          deadLetter.CreatedAt,
			DeliveryDurationMs: deliveryDuration,
			UserName:           deadLetter.UserName,
			ResponseBody:       deadLetter.ResponseBody.String,
			ErrorMessage:       deadLetter.ErrorMessage.String,
		})
	}

	response.RespondWithJSON(w, http.StatusOK, responseDeadLetters)
}
