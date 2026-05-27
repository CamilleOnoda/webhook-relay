package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to find user", err)
		return
	}

	correctPassword, err := auth.CheckPassword(req.Password, dbUser.HashedPassword)
	if err != nil || !correctPassword {
		respondWithError(w, http.StatusUnauthorized,
			"Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwt_secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to create token", err)
		return
	}

	responseUser := User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Token:     token,
		IsAdmin:   dbUser.IsAdmin,
	}

	respondWithJSON(w, http.StatusOK, responseUser)

}
