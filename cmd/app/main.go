package main

import (
	"github.com/stanislavCasciuc/atom-fit-go/internal/api"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
	"github.com/stanislavCasciuc/atom-fit-go/internal/database"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.Envs

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))

	db, err := database.New(cfg.DbCfg)
	if err != nil {
		log.Error("cannot to connect to db", err)
	}

	if err := db.Ping(); err != nil {
		log.Error("cannot to ping to db", err)
	}

	log.Info("database successfully connected")

	server := api.NewServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Error("cannot to run api server ", err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	}

	return log
}