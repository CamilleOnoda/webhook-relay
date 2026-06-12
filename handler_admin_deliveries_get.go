package main

import "net/http"

func (cfg *apiConfig) handlerGetDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	deliveries, err := cfg.db.GetAllDeliveries(r.Context())
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
			CreatedAt:    delivery.CreatedAt,
			Status:       delivery.Status,
			StatusCode:   statusCode,
			EndpointName: delivery.EndpointName,
			UserName:     delivery.UserName,
		})
	}
	respondWithJSON(w, http.StatusOK, responseDeliveries)
}
