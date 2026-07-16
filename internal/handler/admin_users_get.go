package handler

import (
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
)

func HandleGetUsers(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}

	users, err := cfg.DB.GetAllUsers(r.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"failed to list users", err)
		return
	}
	responseUsers := []User{}
	for _, user := range users {
		responseUsers = append(responseUsers, User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsAdmin:   user.IsAdmin,
		})
	}

	response.RespondWithJSON(w, http.StatusOK, responseUsers)
}
