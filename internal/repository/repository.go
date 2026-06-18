package repository

import (
	"context"
	"errors"

	"food-control/internal/models"
)

var ErrNotFound = errors.New("not found")

var ErrInsufficientIngredients = errors.New("insufficient ingredients")

type Store interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	ListExpiringProducts(ctx context.Context, days int) ([]models.Product, error)
	GetProduct(ctx context.Context, id int64) (models.Product, error)
	CreateProduct(ctx context.Context, input models.CreateProductRequest) (models.Product, error)
	UpdateProduct(ctx context.Context, id int64, input models.CreateProductRequest) (models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListRecipes(ctx context.Context) ([]models.Recipe, error)
	CreateRecipe(ctx context.Context, input models.CreateRecipeRequest) (models.Recipe, error)
	GetRecipeAvailability(ctx context.Context, recipeID int64) (models.RecipeAvailability, error)
	ListRecipeSuggestions(ctx context.Context) ([]models.RecipeAvailability, error)
	CookRecipe(ctx context.Context, recipeID int64) (models.CookRecipeResult, error)
}
