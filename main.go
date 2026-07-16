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
	"github.com/CamilleOnoda/webhook-relay.git/internal/handler"
	"github.com/CamilleOnoda/webhook-relay.git/internal/middleware"
	"github.com/CamilleOnoda/webhook-relay.git/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

	cfg := &handler.Config{
		DB:          database.New(db),
		Environment: os.Getenv("ENV"),
		BaseURL:     os.Getenv("BASE_URL"),
		AuthConfig:  authConfig,
	}

	httpDelivery := service.NewDeliveryService()
	processor := service.NewDatabaseDeliveryProcessor(cfg.DB, httpDelivery, 5)
	go processor.Start(ctx, 10*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/health", handlerReadiness)
	mux.Handle("GET /api/endpoints",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEndpoints(cfg, w, r)
			})))
	mux.Handle("GET /api/endpoints/{id}",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEndpointByID(cfg, w, r)
			})))
	mux.Handle("GET /api/events",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleListEventsByUser(cfg, w, r)
			})))
	mux.Handle("GET /api/deliveries",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleListDeliveriesByUser(cfg, w, r)
			})))
	mux.Handle("POST /api/endpoints",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleCreateEndpoint(cfg, w, r)
			})))
	mux.Handle("GET /api/recent-activity",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUserRecentActivity(cfg, w, r)
			})))
	mux.Handle("GET /api/stats",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUserStats(cfg, w, r)
			})))
	mux.Handle("DELETE /api/endpoints/{id}",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleDeleteEndpointByID(cfg, w, r)
			})))
	mux.HandleFunc("POST /webhooks/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleReceiveWebhook(cfg, w, r)
		})
	mux.HandleFunc("POST /api/users",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleUsersCreate(cfg, w, r)
		})
	mux.HandleFunc("POST /api/login",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleLogin(cfg, w, r)
		})
	mux.HandleFunc("POST /api/refresh",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleRefreshToken(cfg, w, r)
		})
	mux.HandleFunc("POST /api/revoke",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleRevoke(cfg, w, r)
		})

	// admin endpoints //
	mux.Handle("DELETE /admin/users/{id}",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleDeleteUserByID(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/users",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUsers(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/stats",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAdminStats(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/endpoints",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAllEndpoints(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/recent-activity",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetRecentActivity(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/events",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEvents(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/deliveries",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetDeliveries(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("POST /admin/dead-letter/{id}/replay",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleReplayDeadLetter(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/deliveries/dead-letter",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetDeadLetters(cfg, w, r)
			}), cfg.DB)))
	mux.Handle("GET /admin/deliveries/{id}",
		middleware.Auth(string(cfg.AuthConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAdminDeliveryDetails(cfg, w, r)
			}), cfg.DB)))

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
