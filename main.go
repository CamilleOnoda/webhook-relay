package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db          *database.Queries
	environment string
	baseURL     string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	env := os.Getenv("ENV")
	baseURL := os.Getenv("BASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	cfg := &apiConfig{
		db:          database.New(db),
		environment: env,
		baseURL:     baseURL,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(
		http.Dir("./internal/static/"))))

	mux.HandleFunc("GET /api/health", handlerReadiness)
	mux.HandleFunc("GET /api/endpoints", cfg.handlerListEndpoints)
	mux.HandleFunc("GET /api/endpoints/{id}", cfg.handlerGetEndpointByID)
	mux.HandleFunc("POST /api/endpoints", cfg.handlerCreateEndpoint)
	mux.HandleFunc("POST /webhooks/{id}", cfg.handlerReceiveWebhook)
	mux.HandleFunc("DELETE /admin/endpoints/delete", cfg.handlerDeleteAllEndpoints)
	mux.HandleFunc("DELETE /api/endpoints/{id}", cfg.handlerDeleteEndpointByID)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
