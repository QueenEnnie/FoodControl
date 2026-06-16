package repository

import (
	"context"
	"errors"

	"food-control/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) ListProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), quantity::float8, unit, expiration_date, created_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY expiration_date NULLS LAST, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.Unit, &product.ExpirationDate, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *Repository) GetProduct(ctx context.Context, id int64) (models.Product, error) {
	var product models.Product

	err := r.db.QueryRow(ctx, `
		SELECT id, name, COALESCE(description, ''), quantity::float8, unit, expiration_date, created_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Quantity,
		&product.Unit,
		&product.ExpirationDate,
		&product.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, ErrNotFound
	}

	return product, err
}

func (r *Repository) CreateProduct(ctx context.Context, input models.CreateProductRequest) (models.Product, error) {
	var product models.Product
	err := r.db.QueryRow(ctx, `
		INSERT INTO products (name, description, quantity, unit, expiration_date)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5)
		RETURNING id, name, COALESCE(description, ''), quantity::float8, unit, expiration_date, created_at`,
		input.Name, input.Description, input.Quantity, input.Unit, input.ExpirationDate,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.Unit, &product.ExpirationDate, &product.CreatedAt)
	return product, err
}

func (r *Repository) UpdateProduct(ctx context.Context, id int64, input models.CreateProductRequest) (models.Product, error) {
	var product models.Product

	err := r.db.QueryRow(ctx, `
		UPDATE products
		SET name = $1,
		    description = NULLIF($2, ''),
		    quantity = $3,
		    unit = $4,
		    expiration_date = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, COALESCE(description, ''), quantity::float8, unit, expiration_date, created_at`,
		input.Name, input.Description, input.Quantity, input.Unit, input.ExpirationDate, id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Quantity,
		&product.Unit,
		&product.ExpirationDate,
		&product.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, ErrNotFound
	}

	return product, err
}

func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
