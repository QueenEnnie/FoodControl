package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"food-control/internal/models"

	"github.com/jackc/pgx/v5"
)

const recipeAvailabilityQuery = `
WITH ingredient_stock AS (
    SELECT 
        ri.recipe_id,
        ri.product_name,
        ri.required_quantity::float8,
        ri.unit,
        COALESCE(SUM(p.quantity) FILTER (
            WHERE (p.expiration_date IS NULL OR p.expiration_date >= CURRENT_DATE)
        ), 0)::float8 AS available_quantity
    FROM recipe_ingredients ri
    LEFT JOIN products p 
        ON lower(p.name) = lower(ri.product_name) 
       AND lower(p.unit) = lower(ri.unit)
    GROUP BY ri.recipe_id, ri.id, ri.product_name, ri.required_quantity, ri.unit
),
recipe_calculations AS (
    SELECT
        r.id AS recipe_id,
        r.name AS recipe_name,
        COALESCE(BOOL_AND(i.available_quantity >= i.required_quantity), TRUE) AS can_cook,
        COUNT(CASE WHEN i.available_quantity < i.required_quantity THEN 1 END) AS missing_count,
        COALESCE(
            JSON_AGG(
                JSON_BUILD_OBJECT(
                    'product_name', i.product_name,
                    'required_quantity', i.required_quantity,
                    'available_quantity', i.available_quantity,
                    'unit', i.unit,
                    'is_enough', i.available_quantity >= i.required_quantity
                )
            ) FILTER (WHERE i.product_name IS NOT NULL), '[]'::json
        ) AS ingredients_json
    FROM recipes r
    LEFT JOIN ingredient_stock i ON i.recipe_id = r.id
    GROUP BY r.id, r.name
)
SELECT recipe_id, recipe_name, can_cook, missing_count, ingredients_json
FROM recipe_calculations
`

func (r *Repository) ListRecipes(ctx context.Context) ([]models.Recipe, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.name, COALESCE(r.description, ''), ri.id, ri.product_name,
		       ri.required_quantity::float8, ri.unit
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON ri.recipe_id = r.id
		ORDER BY r.name, ri.product_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipesByID := make(map[int64]*models.Recipe)
	var order []int64
	for rows.Next() {
		var recipeID int64
		var recipeName, description string
		var ingredientID *int64
		var productName, unit *string
		var requiredQuantity *float64
		if err := rows.Scan(&recipeID, &recipeName, &description, &ingredientID, &productName, &requiredQuantity, &unit); err != nil {
			return nil, err
		}
		recipe, ok := recipesByID[recipeID]
		if !ok {
			recipe = &models.Recipe{ID: recipeID, Name: recipeName, Description: description}
			recipesByID[recipeID] = recipe
			order = append(order, recipeID)
		}
		if ingredientID != nil {
			recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
				ID:               *ingredientID,
				RecipeID:         recipeID,
				ProductName:      *productName,
				RequiredQuantity: *requiredQuantity,
				Unit:             *unit,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	recipes := make([]models.Recipe, 0, len(order))
	for _, id := range order {
		recipes = append(recipes, *recipesByID[id])
	}
	return recipes, nil
}

func (r *Repository) CreateRecipe(ctx context.Context, input models.CreateRecipeRequest) (models.Recipe, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Recipe{}, err
	}
	defer tx.Rollback(ctx)

	var recipe models.Recipe
	err = tx.QueryRow(ctx, `
		INSERT INTO recipes (name, description)
		VALUES ($1, NULLIF($2, ''))
		RETURNING id, name, COALESCE(description, '')`,
		input.Name, input.Description,
	).Scan(&recipe.ID, &recipe.Name, &recipe.Description)
	if err != nil {
		return models.Recipe{}, err
	}

	for _, ingredient := range input.Ingredients {
		var created models.RecipeIngredient
		err := tx.QueryRow(ctx, `
			INSERT INTO recipe_ingredients (recipe_id, product_name, required_quantity, unit)
			VALUES ($1, $2, $3, $4)
			RETURNING id, recipe_id, product_name, required_quantity::float8, unit`,
			recipe.ID, ingredient.ProductName, ingredient.RequiredQuantity, ingredient.Unit,
		).Scan(&created.ID, &created.RecipeID, &created.ProductName, &created.RequiredQuantity, &created.Unit)
		if err != nil {
			return models.Recipe{}, err
		}
		recipe.Ingredients = append(recipe.Ingredients, created)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Recipe{}, err
	}
	return recipe, nil
}

func (r *Repository) GetRecipeAvailability(ctx context.Context, recipeID int64) (models.RecipeAvailability, error) {
	query := fmt.Sprintf(recipeAvailabilityQuery, "WHERE r.id = $1")

	var availability models.RecipeAvailability
	var ingredientsJSON []byte

	err := r.db.QueryRow(ctx, query, recipeID).Scan(
		&availability.RecipeID,
		&availability.RecipeName,
		&availability.CanCook,
		&availability.MissingCount,
		&ingredientsJSON,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.RecipeAvailability{}, ErrNotFound
	}
	if err != nil {
		return models.RecipeAvailability{}, err
	}

	if err := json.Unmarshal(ingredientsJSON, &availability.Ingredients); err != nil {
		return models.RecipeAvailability{}, err
	}

	return availability, nil
}

func (r *Repository) ListRecipeSuggestions(ctx context.Context) ([]models.RecipeAvailability, error) {
	query := fmt.Sprintf(recipeAvailabilityQuery, "") + " ORDER BY can_cook DESC, missing_count ASC, recipe_name"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.RecipeAvailability
	for rows.Next() {
		var availability models.RecipeAvailability
		var ingredientsJSON []byte

		err := rows.Scan(
			&availability.RecipeID,
			&availability.RecipeName,
			&availability.CanCook,
			&availability.MissingCount,
			&ingredientsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(ingredientsJSON, &availability.Ingredients); err != nil {
			return nil, err
		}

		suggestions = append(suggestions, availability)
	}

	return suggestions, rows.Err()
}

func (r *Repository) recipeExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM recipes WHERE id = $1)`, id).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return exists, err
}
