package main

import "net/http"

func (cfg *apiConfig) handlerGetDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deliveries, err := cfg.db.GetAlldeliveries(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
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
			EndpointName: delivery.EndpointName,
			Status:       delivery.Status,
			StatusCode:   statusCode,
			TargetUrl:    delivery.TargetUrl,
			CreatedAt:    delivery.CreatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, responseDeliveries)
}
