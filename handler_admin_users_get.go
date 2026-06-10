package main

import (
	"net/http"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Method not allowed", nil)
		return
	}
	var users []database.GetAllUsersRow
	users, err := cfg.db.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
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

	respondWithJSON(w, http.StatusOK, responseUsers)
}
