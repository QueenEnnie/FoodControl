package db

import (
    "context"
    "log"
    "github.com/jackc/pgx/v5/pgxpool"
)

func New() *pgxpool.Pool {
    connStr := "postgres://fcuser:fcpass@localhost:5433/fooddb"

    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        log.Fatal(err)
    }

    return pool
}