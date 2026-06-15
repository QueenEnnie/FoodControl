package repository

import (
	"strings"
	"testing"

	"food-control/internal/models"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateProductRequest
		wantErr string
	}{
		{
			name: "valid product",
			input: models.CreateProductRequest{
				Name:     "milk",
				Quantity: 1,
				Unit:     "l",
			},
		},
		{
			name: "zero quantity is valid",
			input: models.CreateProductRequest{
				Name:     "milk",
				Quantity: 0,
				Unit:     "l",
			},
		},
		{
			name: "empty name",
			input: models.CreateProductRequest{
				Name:     " ",
				Quantity: 1,
				Unit:     "l",
			},
			wantErr: "name is required",
		},
		{
			name: "negative quantity",
			input: models.CreateProductRequest{
				Name:     "milk",
				Quantity: -1,
				Unit:     "l",
			},
			wantErr: "quantity must be greater than or equal to zero",
		},
		{
			name: "empty unit",
			input: models.CreateProductRequest{
				Name:     "milk",
				Quantity: 1,
				Unit:     " ",
			},
			wantErr: "unit is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProduct(tt.input)
			assertValidationError(t, err, tt.wantErr)
		})
	}
}

func TestValidateRecipe(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateRecipeRequest
		wantErr string
	}{
		{
			name: "valid recipe",
			input: models.CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []models.RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
		},
		{
			name: "empty recipe name",
			input: models.CreateRecipeRequest{
				Name: " ",
				Ingredients: []models.RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
			wantErr: "name is required",
		},
		{
			name: "ingredient without product name",
			input: models.CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []models.RecipeIngredient{
					{ProductName: " ", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].product_name is required",
		},
		{
			name: "ingredient with zero quantity",
			input: models.CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []models.RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].required_quantity must be greater than zero",
		},
		{
			name: "ingredient with negative quantity",
			input: models.CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []models.RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: -0.3, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].required_quantity must be greater than zero",
		},
		{
			name: "ingredient without unit",
			input: models.CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []models.RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0.3, Unit: " "},
				},
			},
			wantErr: "ingredients[0].unit is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRecipe(tt.input)
			assertValidationError(t, err, tt.wantErr)
		})
	}
}

func assertValidationError(t *testing.T, err error, want string) {
	t.Helper()

	if want == "" {
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("expected error containing %q, got nil", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected error containing %q, got %q", want, err.Error())
	}
}
