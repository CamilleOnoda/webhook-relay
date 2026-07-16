package handler

import (
	"encoding/json"
	response "github.com/CamilleOnoda/webhook-relay.git/internal/response"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

func HandleLogin(cfg *Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to find user", err)
		return
	}

	correctPassword, err := auth.CheckPassword(req.Password, dbUser.HashedPassword)
	if err != nil || !correctPassword {
		response.RespondWithError(w, http.StatusUnauthorized,
			"Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		string(cfg.AuthConfig.AccessTokenSecret),
		cfg.AuthConfig.AccessTokenTTL,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to create token", err)
		return
	}

	rawToken, err := auth.MakeRefreshToken()
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to make refresh token", err)
		return
	}
	hashedToken := auth.HashRefreshToken(rawToken)
	expiresAt := time.Now().UTC().Add(cfg.AuthConfig.RefreshTokenTTL)
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     hashedToken,
		UserID:    dbUser.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError,
			"Failed to store refresh token", err)
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.AuthConfig.RefreshCookieName,
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.AuthConfig.CookieSecure,
		SameSite: cfg.AuthConfig.CookieSameSite,
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

	response.RespondWithJSON(w, http.StatusOK, responseUser)

}
