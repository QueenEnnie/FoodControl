package httpapi

import (
	"context"

	"food-control/internal/domain"
)

type ProductRepository interface {
	ListProducts(ctx context.Context) ([]domain.Product, error)
	ListExpiringProducts(ctx context.Context, days int) ([]domain.Product, error)
	GetProduct(ctx context.Context, id int64) (domain.Product, error)
	CreateProduct(ctx context.Context, input domain.CreateProductRequest) (domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, input domain.CreateProductRequest) (domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type RecipeRepository interface {
	ListRecipes(ctx context.Context) ([]domain.Recipe, error)
	CreateRecipe(ctx context.Context, input domain.CreateRecipeRequest) (domain.Recipe, error)
	GetRecipeAvailability(ctx context.Context, recipeID int64) (domain.RecipeAvailability, error)
	ListRecipeSuggestions(ctx context.Context) ([]domain.RecipeAvailability, error)
	CookRecipe(ctx context.Context, recipeID int64) (domain.CookRecipeResult, error)
}

type Repository interface {
	ProductRepository
	RecipeRepository
}
