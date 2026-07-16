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

type apiConfig struct {
	db          *database.Queries
	environment string
	baseURL     string
	authConfig  *auth.Config
}

func (cfg *apiConfig) handlerConfig() *handler.Config {
	return &handler.Config{
		DB:          cfg.db,
		Environment: cfg.environment,
		BaseURL:     cfg.baseURL,
		AuthConfig:  cfg.authConfig,
	}
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
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEndpoints(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("GET /api/endpoints/{id}",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEndpointByID(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("GET /api/events",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleListEventsByUser(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("GET /api/deliveries",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleListDeliveriesByUser(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("POST /api/endpoints",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleCreateEndpoint(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("GET /api/recent-activity",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUserRecentActivity(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("GET /api/stats",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUserStats(cfg.handlerConfig(), w, r)
			})))
	mux.Handle("DELETE /api/endpoints/{id}",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleDeleteEndpointByID(cfg.handlerConfig(), w, r)
			})))
	mux.HandleFunc("POST /webhooks/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleReceiveWebhook(cfg.handlerConfig(), w, r)
		})
	mux.HandleFunc("POST /api/users",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleUsersCreate(cfg.handlerConfig(), w, r)
		})
	mux.HandleFunc("POST /api/login",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleLogin(cfg.handlerConfig(), w, r)
		})
	mux.HandleFunc("POST /api/refresh",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleRefreshToken(cfg.handlerConfig(), w, r)
		})
	mux.HandleFunc("POST /api/revoke",
		func(w http.ResponseWriter, r *http.Request) {
			handler.HandleRevoke(cfg.handlerConfig(), w, r)
		})

	// admin endpoints //
	mux.Handle("DELETE /admin/users/{id}",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleDeleteUserByID(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/users",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetUsers(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/stats",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAdminStats(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/endpoints",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAllEndpoints(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/recent-activity",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetRecentActivity(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/events",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetEvents(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/deliveries",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetDeliveries(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("POST /admin/dead-letter/{id}/replay",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleReplayDeadLetter(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/deliveries/dead-letter",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetDeadLetters(cfg.handlerConfig(), w, r)
			}), cfg.db)))
	mux.Handle("GET /admin/deliveries/{id}",
		middleware.Auth(string(cfg.authConfig.AccessTokenSecret),
			middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.HandleGetAdminDeliveryDetails(cfg.handlerConfig(), w, r)
			}), cfg.db)))

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
