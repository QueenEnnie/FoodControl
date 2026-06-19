package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"food-control/internal/config"
	"food-control/internal/storage/postgres"
	"food-control/internal/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	ctx := context.Background()

	pool, err := connectDatabase(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("failed to connect database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.New(pool)
	api := httpapi.New(repo, logger)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           api,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("FoodControl API listening", slog.String("addr", cfg.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		return
	}
	logger.Info("FoodControl API stopped")
}

func connectDatabase(ctx context.Context, databaseURL string, logger *slog.Logger) (*pgxpool.Pool, error) {
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		pool, err := postgres.NewPool(ctx, databaseURL)
		if err == nil {
			return pool, nil
		}
		logger.Warn("database is not ready yet, retrying",
			slog.Int("attempt", attempt),
			slog.Int("max_attempts", 10),
			slog.Any("error", err),
		)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
