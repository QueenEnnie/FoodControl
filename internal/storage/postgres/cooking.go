package postgres

import (
	"context"
	"errors"
	"math"

	"food-control/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const maxCookAttempts = 3

type recipeIngredient struct {
	productName string
	quantity    float64
	unit        string
}

type productStock struct {
	id       int64
	quantity float64
}

func (r *Repository) CookRecipe(ctx context.Context, recipeID int64) (domain.CookRecipeResult, error) {
	var result domain.CookRecipeResult
	var err error

	for attempt := 1; attempt <= maxCookAttempts; attempt++ {
		result, err = r.cookRecipeOnce(ctx, recipeID)
		if err == nil || !isRetryableTransactionError(err) {
			return result, err
		}
	}

	return domain.CookRecipeResult{}, err
}

func (r *Repository) cookRecipeOnce(ctx context.Context, recipeID int64) (domain.CookRecipeResult, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return domain.CookRecipeResult{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var recipeName string
	err = tx.QueryRow(ctx, `SELECT name FROM recipes WHERE id = $1`, recipeID).Scan(&recipeName)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CookRecipeResult{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.CookRecipeResult{}, err
	}

	ingredients, err := listRecipeIngredients(ctx, tx, recipeID)
	if err != nil {
		return domain.CookRecipeResult{}, err
	}

	result := domain.CookRecipeResult{
		RecipeID:   recipeID,
		RecipeName: recipeName,
		Status:     "cooked",
		Consumed:   make([]domain.ConsumedIngredient, 0, len(ingredients)),
	}

	for _, ingredient := range ingredients {
		stocks, err := lockProductStocks(ctx, tx, ingredient)
		if err != nil {
			return domain.CookRecipeResult{}, err
		}

		available := 0.0
		for _, stock := range stocks {
			available += stock.quantity
		}
		if available+1e-9 < ingredient.quantity {
			return domain.CookRecipeResult{}, domain.ErrInsufficientIngredients
		}

		remaining := ingredient.quantity
		for _, stock := range stocks {
			if remaining <= 1e-9 {
				break
			}

			consumed := math.Min(stock.quantity, remaining)
			if _, err := tx.Exec(ctx, `
				UPDATE products
				SET quantity = quantity - $1
				WHERE id = $2`,
				consumed, stock.id,
			); err != nil {
				return domain.CookRecipeResult{}, err
			}
			remaining -= consumed
		}

		result.Consumed = append(result.Consumed, domain.ConsumedIngredient{
			ProductName: ingredient.productName,
			Quantity:    ingredient.quantity,
			Unit:        ingredient.unit,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.CookRecipeResult{}, err
	}
	return result, nil
}

func listRecipeIngredients(ctx context.Context, tx pgx.Tx, recipeID int64) ([]recipeIngredient, error) {
	rows, err := tx.Query(ctx, `
		SELECT product_name, required_quantity::float8, unit
		FROM recipe_ingredients
		WHERE recipe_id = $1
		ORDER BY lower(product_name), lower(unit), id`,
		recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []recipeIngredient
	for rows.Next() {
		var ingredient recipeIngredient
		if err := rows.Scan(&ingredient.productName, &ingredient.quantity, &ingredient.unit); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, rows.Err()
}

func lockProductStocks(ctx context.Context, tx pgx.Tx, ingredient recipeIngredient) ([]productStock, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, quantity::float8
		FROM products
		WHERE lower(name) = lower($1)
		  AND lower(unit) = lower($2)
		  AND deleted_at IS NULL
		  AND quantity > 0
		  AND (expiration_date IS NULL OR expiration_date >= CURRENT_DATE)
		ORDER BY expiration_date NULLS LAST, id
		FOR UPDATE`,
		ingredient.productName, ingredient.unit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []productStock
	for rows.Next() {
		var stock productStock
		if err := rows.Scan(&stock.id, &stock.quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, rows.Err()
}

func isRetryableTransactionError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "40001" || pgErr.Code == "40P01"
}
