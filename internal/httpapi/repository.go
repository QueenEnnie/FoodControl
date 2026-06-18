package httpapi

import (
	"context"

	"food-control/internal/models"
)

type ProductRepository interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	ListExpiringProducts(ctx context.Context, days int) ([]models.Product, error)
	GetProduct(ctx context.Context, id int64) (models.Product, error)
	CreateProduct(ctx context.Context, input models.CreateProductRequest) (models.Product, error)
	UpdateProduct(ctx context.Context, id int64, input models.CreateProductRequest) (models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type RecipeRepository interface {
	ListRecipes(ctx context.Context) ([]models.Recipe, error)
	CreateRecipe(ctx context.Context, input models.CreateRecipeRequest) (models.Recipe, error)
	GetRecipeAvailability(ctx context.Context, recipeID int64) (models.RecipeAvailability, error)
	ListRecipeSuggestions(ctx context.Context) ([]models.RecipeAvailability, error)
	CookRecipe(ctx context.Context, recipeID int64) (models.CookRecipeResult, error)
}

type Repository interface {
	ProductRepository
	RecipeRepository
}
