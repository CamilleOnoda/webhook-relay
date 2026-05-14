package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestEndpointCreate(t *testing.T) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpassword"

	tests := []struct {
		name         string
		endpointName string
		targetURL    string
		wantStatus   int
	}{
		{
			name:         "successful endpoint creation returns 201",
			endpointName: "test-endpoint",
			targetURL:    "https://example.com/webhook",
			wantStatus:   http.StatusCreated,
		},
		{
			name:         "missing name returns 400",
			endpointName: "",
			targetURL:    "https://example.com/webhook",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "missing target URL returns 400",
			endpointName: "test-endpoint",
			targetURL:    "",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "invalid URL scheme returns 400",
			endpointName: "test-endpoint",
			targetURL:    "ftp://example.com/webhook",
			wantStatus:   http.StatusBadRequest,
		},
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	if err := runMigrations(connStr); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/endpoints",
				bytes.NewBufferString(
					fmt.Sprintf(
						`{"name": "%s", "target_url": "%s"}`,
						test.endpointName, test.targetURL)))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			queries := database.New(db)
			apiCfg := apiConfig{
				db: queries,
			}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /api/endpoints", apiCfg.handlerCreateEndpoint)
			mux.ServeHTTP(rec, req)

			if rec.Code != test.wantStatus {
				t.Errorf("expected status %d, got %d", test.wantStatus, rec.Code)
			}
		})
	}
}

func runMigrations(connStr string) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, "./sql/schema")
}
