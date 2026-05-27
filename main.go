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
	mux.HandleFunc("GET /admin/health", handlerReadiness)

	mux.Handle("GET /api/endpoints",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerGetEndpoints)))

	mux.Handle("GET /api/endpoints/{id}",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerGetEndpointByID)))

	mux.Handle("GET /api/events",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerListEventsByUser)))

	mux.Handle("GET /api/deliveries",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerListDeliveriesByUser)))

	mux.Handle("/api/endpoints",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerCreateEndpoint)))

	mux.HandleFunc("POST /webhooks/{id}", cfg.handlerReceiveWebhook)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	mux.HandleFunc("/api/login", cfg.handlerLogin)

	mux.Handle("DELETE /admin/users/delete",
		cfg.authMiddleware(
			cfg.adminMiddleware(
				http.HandlerFunc(cfg.handlerDeleteUsers))))

	mux.Handle("DELETE /admin/endpoints/delete",
		cfg.authMiddleware(
			cfg.adminMiddleware(
				http.HandlerFunc(cfg.handlerDeleteAllEndpoints))))

	mux.Handle("DELETE /api/endpoints/{id}",
		cfg.authMiddleware(
			http.HandlerFunc(cfg.handlerDeleteEndpointByID)))

	fileServer := http.FileServer(http.Dir("./internal/static"))
	mux.Handle("/", fileServer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
