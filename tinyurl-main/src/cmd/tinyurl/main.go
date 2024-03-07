package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinyurl/internal/api"
	"tinyurl/internal/config"
	"tinyurl/internal/db"
	"tinyurl/internal/shorten"
	"tinyurl/internal/storage"
)

const (
	envLocal       = "local"
	envDevelopment = "development"
	envProduction  = "production"
)

func main() {

	cfg := config.Load()
	logger := setupLogger(cfg.Env)
	slog.SetDefault(logger)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), time.Duration(cfg.DbConTimeout)*time.Second)
	defer dbCancel()

	postgresClient, err := db.Connect(dbCtx, cfg.Dsn())
	if err != nil {
		slog.Error("error connecting to db: %v", err)
		os.Exit(1)
	}

	// Это всего лишь пул подключений!
	postgresDb := postgresClient.Client()

	linkStorage := storage.NewPostgresDB(postgresDb)
	dwarfService := shorten.NewService(linkStorage)

	srv := api.New(dwarfService, cfg.BaseUrl())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(cfg.ListenAddr(), srv); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error running api: %v", err)
			os.Exit(1)
		}
	}()

	slog.Info(fmt.Sprintf("api started on %s", cfg.ListenAddr()))
	<-quit
	slog.Info("api shutting down")
}

func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProduction:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envDevelopment:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}
