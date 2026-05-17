package main

import (
	"context"
	"errors"
	"food-control/internal/config"
	"food-control/internal/db"
	"food-control/internal/httpapi"
	"food-control/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := connectDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	api := httpapi.New(repo)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           api,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("FoodControl API listening on %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func connectDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		pool, err := db.New(ctx, databaseURL)
		if err == nil {
			return pool, nil
		}
		log.Printf("database is not ready yet, retrying (%d/10): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
