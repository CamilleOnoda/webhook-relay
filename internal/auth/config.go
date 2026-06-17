package auth

import (
	"errors"
	"net/http"
	"os"
	"time"
)

type Config struct {
	AccessTokenSecret []byte
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
	Issuer            string
	RefreshCookieName string
	CookieSecure      bool
	CookieSameSite    http.SameSite
}

func NewConfig() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 characters")
	}
	env := os.Getenv("ENV")
	return &Config{
		AccessTokenSecret: []byte(jwtSecret),
		AccessTokenTTL:    12 * time.Hour,
		RefreshTokenTTL:   30 * 24 * time.Hour,
		Issuer:            "webhook-relay",
		RefreshCookieName: "refresh_token",
		CookieSecure:      env == "staging" || env == "production",
		CookieSameSite:    http.SameSiteLaxMode,
	}, nil
}
