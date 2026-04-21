package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mediflow/backend/internal/alert"
	"github.com/mediflow/backend/internal/analytics"
	"github.com/mediflow/backend/internal/auth"
	"github.com/mediflow/backend/internal/equipment"
	"github.com/mediflow/backend/internal/request"
	"github.com/mediflow/backend/internal/shared/db"
	"github.com/mediflow/backend/internal/shared/jobs"
	"github.com/mediflow/backend/internal/shared/redis"
	"github.com/mediflow/backend/internal/shared/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration (simplified for Phase 1)
	dbCfg := db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}
	if dbCfg.Host == "" {
		dbCfg.Host = "localhost"
		dbCfg.Port = "5432"
		dbCfg.User = "mediflow"
		dbCfg.Password = "mediflow_secret"
		dbCfg.DBName = "mediflow_db"
	}

	redisCfg := redis.Config{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}
	if redisCfg.Host == "" {
		redisCfg.Host = "localhost"
		redisCfg.Port = "6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_key"
	}

	// 1. Connect to Database
	database, err := db.Connect(dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// 2. Connect to Redis
	redisClient, err := redis.Connect(redisCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// 3. Initialize WebSocket Hub & Bridge
	wsHub := websocket.NewHub()
	go wsHub.Run()
	go websocket.StartBridge(context.Background(), wsHub, redisClient)

	// 4. Initialize Services & Handlers
	jwtManager := auth.NewJWTManager(jwtSecret, 15*time.Minute)
	
	equipRepo := equipment.NewRepository(database)
	availMgr := equipment.NewAvailabilityManager(equipRepo, redisClient)
	equipService := equipment.NewService(equipRepo, redisClient, availMgr)
	equipHandler := equipment.NewHandler(equipService)
	
	reqRepo := request.NewRepository(database)
	reqService := request.NewService(reqRepo, equipService)
	reqHandler := request.NewHandler(reqService)

	alertService := alert.NewService(database, redisClient)
	
	analyticsRepo := analytics.NewRepository(database)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)
	
	worker := jobs.NewWorker(reqService, equipService)
	go worker.Start(context.Background())

	authHandler := auth.NewHandler(jwtManager)

	// 4. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", authHandler.Login)
		
		// WebSocket endpoint
		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			claims, err := jwtManager.Verify(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			websocket.ServeWs(wsHub, w, r, claims.TenantID.String(), claims.UserID.String())
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtManager))
			equipHandler.RegisterRoutes(r)
			reqHandler.RegisterRoutes(r)
			analyticsHandler.RegisterRoutes(r)
		})
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info().Msgf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

// Note: GetClaimsMiddleware is a small helper to be added to middleware
