package repository

import (
	"fmt"
	"strings"

	"food-control/internal/models"
)

func ValidateProduct(input models.CreateProductRequest) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if input.Quantity < 0 {
		return fmt.Errorf("quantity must be greater than or equal to zero")
	}
	if strings.TrimSpace(input.Unit) == "" {
		return fmt.Errorf("unit is required")
	}
	return nil
}

func ValidateRecipe(input models.CreateRecipeRequest) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	for i, ingredient := range input.Ingredients {
		if strings.TrimSpace(ingredient.ProductName) == "" {
			return fmt.Errorf("ingredients[%d].product_name is required", i)
		}
		if ingredient.RequiredQuantity <= 0 {
			return fmt.Errorf("ingredients[%d].required_quantity must be greater than zero", i)
		}
		if strings.TrimSpace(ingredient.Unit) == "" {
			return fmt.Errorf("ingredients[%d].unit is required", i)
		}
	}
	return nil
}
