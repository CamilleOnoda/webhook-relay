package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
)

func HandleGetAllEndpoints(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	endpoints, err := cfg.DB.GetAllEndpoints(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"failed to list endpoints", err)
		return
	}

	responseEndpoints := []Endpoint{}
	for _, endpoint := range endpoints {
		responseEndpoints = append(responseEndpoints, Endpoint{
			Name:      endpoint.Name,
			IsActive:  endpoint.IsActive,
			CreatedAt: endpoint.CreatedAt,
			UserID:    endpoint.UserID,
			UserName:  endpoint.UserName,
		})
	}
	response.RespondWithJSON(w, http.StatusOK, responseEndpoints)
}
