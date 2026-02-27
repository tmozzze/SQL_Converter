package main

import (
	"log/slog"
	"os"

	"github.com/tmozzze/SQL_Converter/internal/config"
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

	// Init Repos

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
