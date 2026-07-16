package handler

import (
	"encoding/json"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessToken string    `json:"access_token"`
	IsAdmin     bool      `json:"is_admin"`
}

func HandleUsersCreate(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"invalid request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest,
			"failed to hash the password", err)
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"failed to create user", err)
		return
	}

	responseUser := User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	response.RespondWithJSON(w, http.StatusCreated, responseUser)
}
