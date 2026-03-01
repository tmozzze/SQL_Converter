package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/tmozzze/SQL_Converter/docs"
	"github.com/tmozzze/SQL_Converter/internal/config"
	"github.com/tmozzze/SQL_Converter/internal/http/handler"
	"github.com/tmozzze/SQL_Converter/internal/repository/postgres"
	"github.com/tmozzze/SQL_Converter/internal/service"
	"github.com/tmozzze/SQL_Converter/pkg/database"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Open http://localhost:8080/swagger/index.html.

// @title SQL Converter API
// @version 1.0
// @description Service for converting CSV/XLSX files to PostgreSQL tables.
// @host localhost:8080
// @BasePath /
func main() {

	// Init Config
	cfg := config.MustLoad()

	// Init logger (slog)
	log := setupLogger(cfg.Env)
	log.Info("starting SQL_Converter", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init DB
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Error("failed to init database", slog.Any("err", err))
		os.Exit(1)
	}
	defer db.Close()
	log.Info("connect to DB")

	// Init Repos
	repo := postgres.NewRepository(db, log)

	// Init Service
	svc := service.NewService(repo, log)

	// Init Handler
	handler := handler.NewHandler(svc, log)

	// Init Router
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Init HTTP Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Start Server (net/http)
	go func() {
		log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.Any("err", err))
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sign := <-stop
	log.Info("stopping server", slog.String("signal", sign.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("err", err))
	}

	log.Info("server exited properly")
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal: // Text Debug
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev: // JSON Debug
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd: // JSON Info
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
