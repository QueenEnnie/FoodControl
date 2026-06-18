package postgres

import (
	"food-control/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.Store = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
