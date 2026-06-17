package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
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

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		string(cfg.authConfig.AccessTokenSecret),
		cfg.authConfig.AccessTokenTTL,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to create token", err)
		return
	}

	rawToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to make refresh token", err)
		return
	}
	hashedToken := auth.HashRefreshToken(rawToken)
	expiresAt := time.Now().UTC().Add(cfg.authConfig.RefreshTokenTTL)
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     hashedToken,
		UserID:    dbUser.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Failed to store refresh token", err)
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.authConfig.RefreshCookieName,
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.authConfig.CookieSecure,
		SameSite: cfg.authConfig.CookieSameSite,
		Expires:  expiresAt,
	}

	http.SetCookie(w, cookie)

	responseUser := User{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		AccessToken: accessToken,
		IsAdmin:     dbUser.IsAdmin,
	}

	respondWithJSON(w, http.StatusOK, responseUser)

}
