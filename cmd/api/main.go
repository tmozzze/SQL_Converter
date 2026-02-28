package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/tmozzze/SQL_Converter/internal/config"
	"github.com/tmozzze/SQL_Converter/internal/repository/postgres"
	"github.com/tmozzze/SQL_Converter/pkg/database"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

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
	log.Info("database is initialized")

	// Init Repos
	repo := postgres.NewRepository(db, log)
	fmt.Println(repo)

	// Init Service

	// Start Server (net/http)
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
