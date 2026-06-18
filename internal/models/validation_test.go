package models

import (
	"strings"
	"testing"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateProductRequest
		wantErr string
	}{
		{
			name: "valid product",
			input: CreateProductRequest{
				Name:     "milk",
				Quantity: 1,
				Unit:     "l",
			},
		},
		{
			name: "zero quantity is valid",
			input: CreateProductRequest{
				Name:     "milk",
				Quantity: 0,
				Unit:     "l",
			},
		},
		{
			name: "empty name",
			input: CreateProductRequest{
				Name:     " ",
				Quantity: 1,
				Unit:     "l",
			},
			wantErr: "name is required",
		},
		{
			name: "negative quantity",
			input: CreateProductRequest{
				Name:     "milk",
				Quantity: -1,
				Unit:     "l",
			},
			wantErr: "quantity must be greater than or equal to zero",
		},
		{
			name: "empty unit",
			input: CreateProductRequest{
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
		input   CreateRecipeRequest
		wantErr string
	}{
		{
			name: "valid recipe",
			input: CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
		},
		{
			name: "empty recipe name",
			input: CreateRecipeRequest{
				Name: " ",
				Ingredients: []RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
			wantErr: "name is required",
		},
		{
			name: "ingredient without product name",
			input: CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []RecipeIngredient{
					{ProductName: " ", RequiredQuantity: 0.3, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].product_name is required",
		},
		{
			name: "ingredient with zero quantity",
			input: CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: 0, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].required_quantity must be greater than zero",
		},
		{
			name: "ingredient with negative quantity",
			input: CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []RecipeIngredient{
					{ProductName: "milk", RequiredQuantity: -0.3, Unit: "l"},
				},
			},
			wantErr: "ingredients[0].required_quantity must be greater than zero",
		},
		{
			name: "ingredient without unit",
			input: CreateRecipeRequest{
				Name: "Pancakes",
				Ingredients: []RecipeIngredient{
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
