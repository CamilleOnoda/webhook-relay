package handler

import (
	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

type Config struct {
	DB          *database.Queries
	Environment string
	BaseURL     string
	AuthConfig  *auth.Config
}
