package repository

import (
	"context"
	"errors"

	"food-control/internal/models"

	"github.com/jackc/pgx/v5"
)

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
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.name, ri.product_name, ri.required_quantity::float8,
		       COALESCE(SUM(p.quantity) FILTER (
		       	WHERE lower(p.name) = lower(ri.product_name)
		       	  AND lower(p.unit) = lower(ri.unit)
		       	  AND (p.expiration_date IS NULL OR p.expiration_date >= CURRENT_DATE)
		       ), 0)::float8 AS available_quantity,
		       ri.unit
		FROM recipes r
		JOIN recipe_ingredients ri ON ri.recipe_id = r.id
		LEFT JOIN products p ON lower(p.name) = lower(ri.product_name) AND lower(p.unit) = lower(ri.unit)
		WHERE r.id = $1
		GROUP BY r.id, r.name, ri.id, ri.product_name, ri.required_quantity, ri.unit
		ORDER BY ri.product_name`, recipeID)
	if err != nil {
		return models.RecipeAvailability{}, err
	}
	defer rows.Close()

	availability := models.RecipeAvailability{RecipeID: recipeID, CanCook: true}
	found := false
	for rows.Next() {
		found = true
		var ingredient models.IngredientAvailability
		if err := rows.Scan(&availability.RecipeID, &availability.RecipeName, &ingredient.ProductName, &ingredient.RequiredQuantity, &ingredient.AvailableQuantity, &ingredient.Unit); err != nil {
			return models.RecipeAvailability{}, err
		}
		ingredient.IsEnough = ingredient.AvailableQuantity >= ingredient.RequiredQuantity
		if !ingredient.IsEnough {
			availability.CanCook = false
			availability.MissingCount++
		}
		availability.Ingredients = append(availability.Ingredients, ingredient)
	}
	if err := rows.Err(); err != nil {
		return models.RecipeAvailability{}, err
	}
	if !found {
		exists, err := r.recipeExists(ctx, recipeID)
		if err != nil {
			return models.RecipeAvailability{}, err
		}
		if !exists {
			return models.RecipeAvailability{}, ErrNotFound
		}
		availability.CanCook = true
	}
	return availability, nil
}

func (r *Repository) ListRecipeSuggestions(ctx context.Context) ([]models.RecipeAvailability, error) {
	recipes, err := r.ListRecipes(ctx)
	if err != nil {
		return nil, err
	}

	suggestions := make([]models.RecipeAvailability, 0, len(recipes))
	for _, recipe := range recipes {
		availability, err := r.GetRecipeAvailability(ctx, recipe.ID)
		if err != nil {
			return nil, err
		}
		if len(recipe.Ingredients) == 0 {
			availability.RecipeID = recipe.ID
			availability.RecipeName = recipe.Name
			availability.CanCook = true
		}
		suggestions = append(suggestions, availability)
	}
	return suggestions, nil
}

func (r *Repository) recipeExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM recipes WHERE id = $1)`, id).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return exists, err
}
