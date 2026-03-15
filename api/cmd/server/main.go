package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rknruben56/juwanna-fantasy/api/internal/config"
	"github.com/rknruben56/juwanna-fantasy/api/internal/handler"
	"github.com/rknruben56/juwanna-fantasy/api/internal/repository/postgres"
	"github.com/rknruben56/juwanna-fantasy/api/internal/service"
	"github.com/rknruben56/juwanna-fantasy/api/internal/sleeper"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Configure logging
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	if level == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Configure database pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse database config")
	}
	poolConfig.MaxConns = int32(cfg.DatabaseMaxConns)
	poolConfig.MinConns = int32(cfg.DatabaseMinConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database pool")
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("database ping failed")
	}
	log.Info().Msg("database connected")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.CORSOrigins(),
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// Wire up repositories
	ownerRepo := postgres.NewOwnerRepo(pool)
	seasonRepo := postgres.NewSeasonRepo(pool)
	beltRepo := postgres.NewBeltRepo(pool)
	mappingRepo := postgres.NewMappingRepo(pool)
	leagueConfigRepo := postgres.NewLeagueConfigRepo(pool)

	// Historical service and handler
	historicalSvc := service.NewHistoricalService(ownerRepo, seasonRepo, beltRepo)
	h := handler.New(historicalSvc)
	h.RegisterRoutes(r)

	// Sleeper integration
	cache := sleeper.NewCache()
	sleeperClient := sleeper.NewClient(cache)
	sleeperSvc := service.NewSleeperService(sleeperClient, mappingRepo, leagueConfigRepo, beltRepo, cfg.SleeperLeagueID)

	liveHandler := handler.NewLiveHandler(sleeperSvc)
	liveHandler.RegisterRoutes(r)

	adminHandler := handler.NewAdminHandler(sleeperSvc)
	adminHandler.RegisterRoutes(r)

	// Analytics (rivalries, power rankings, projections)
	rivalryRepo := postgres.NewRivalryRepo(pool)
	analyticsSvc := service.NewAnalyticsService(rivalryRepo, sleeperClient, sleeperSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)
	analyticsHandler.RegisterRoutes(r)

	// Weekly digest
	digestSvc := service.NewDigestService(sleeperSvc, analyticsSvc, rivalryRepo, sleeperClient)
	digestHandler := handler.NewDigestHandler(digestSvc)
	digestHandler.RegisterRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","db":"down"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","db":"up"}`))
	})

	readTimeout := time.Duration(cfg.ReadTimeout) * time.Second
	writeTimeout := time.Duration(cfg.WriteTimeout) * time.Second
	shutdownTimeout := time.Duration(cfg.ShutdownTimeout) * time.Second

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		log.Info().Msgf("server listening on :%d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server shutdown failed")
	}
	log.Info().Msg("server stopped")
}
