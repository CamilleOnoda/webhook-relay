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
	jwt_secret  string
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	cfg := &apiConfig{
		db:          database.New(db),
		environment: os.Getenv("ENV"),
		baseURL:     os.Getenv("BASE_URL"),
		jwt_secret:  os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./internal/static"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("GET /api/health", handlerReadiness)
	mux.HandleFunc("GET /api/endpoints", cfg.handlerListEndpoints)
	mux.HandleFunc("GET /api/endpoints/{id}", cfg.handlerGetEndpointByID)
	mux.HandleFunc("GET /api/events", cfg.handlerListEvents)
	mux.HandleFunc("GET /api/deliveries", cfg.handlerListDeliveries)
	mux.HandleFunc("POST /api/endpoints", cfg.handlerCreateEndpoint)
	mux.HandleFunc("POST /webhooks/{id}", cfg.handlerReceiveWebhook)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("DELETE /admin/endpoints/delete", cfg.handlerDeleteAllEndpoints)
	mux.HandleFunc("DELETE /api/endpoints/{id}", cfg.handlerDeleteEndpointByID)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
