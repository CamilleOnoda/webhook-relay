package handler

import (
	"net/http"

	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
)

func HandleGetDeliveries(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deliveries, err := cfg.DB.GetAllDeliveries(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"failed to list endpoints", err)
		return
	}

	responseDeliveries := []Delivery{}
	for _, delivery := range deliveries {
		var statusCode *int32
		if delivery.StatusCode.Valid {
			statusCode = &delivery.StatusCode.Int32
		}
		responseDeliveries = append(responseDeliveries, Delivery{
			CreatedAt:    delivery.CreatedAt,
			Status:       delivery.Status,
			StatusCode:   statusCode,
			EndpointName: delivery.EndpointName,
			UserName:     delivery.UserName,
		})
	}
	response.RespondWithJSON(w, http.StatusOK, responseDeliveries)
}
