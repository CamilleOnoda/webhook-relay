package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/auth"
	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/CamilleOnoda/webhook-relay.git/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db          *database.Queries
	environment string
	baseURL     string
	authConfig  *auth.Config
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	authConfig, err := auth.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		db:          database.New(db),
		environment: os.Getenv("ENV"),
		baseURL:     os.Getenv("BASE_URL"),
		authConfig:  authConfig,
	}

	httpDelivery := service.NewDeliveryService()
	processor := service.NewDatabaseDeliveryProcessor(cfg.db, httpDelivery, 5)
	go processor.Start(ctx, 10*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/health", handlerReadiness)
	mux.Handle("GET /api/endpoints",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerGetEndpoints)))
	mux.Handle("GET /api/endpoints/{id}",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerGetEndpointByID)))
	mux.Handle("GET /api/events",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerListEventsByUser)))
	mux.Handle("GET /api/deliveries",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerListDeliveriesByUser)))
	mux.Handle("/api/endpoints",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerCreateEndpoint)))
	mux.Handle("DELETE /api/endpoints/{id}",
		cfg.authMiddleware(http.HandlerFunc(cfg.handlerDeleteEndpointByID)))
	mux.HandleFunc("POST /webhooks/{id}", cfg.handlerReceiveWebhook)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	// admin endpoints //
	mux.Handle("DELETE /admin/users/{id}",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerDeleteUserByID))))
	mux.Handle("GET /admin/users",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetUsers))))
	mux.Handle("GET /admin/stats",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetAdminStats))))
	mux.Handle("GET /admin/endpoints",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetAllEndpoints))))
	mux.Handle("GET /admin/recent-activity",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetRecentActivity))))
	mux.Handle("GET /admin/events",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetEvents))))
	mux.Handle("GET /admin/deliveries",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetDeliveries))))
	mux.Handle("POST /admin/deliveries/{id}/replay",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerReplayDeadLetter))))
	mux.Handle("GET /admin/deliveries/dead-letter",
		cfg.authMiddleware(cfg.adminMiddleware(http.HandlerFunc(cfg.handlerGetDeadLetters))))

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
