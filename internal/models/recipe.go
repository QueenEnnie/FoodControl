package models

type Recipe struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Ingredients []RecipeIngredient `json:"ingredients,omitempty"`
}

type RecipeIngredient struct {
	ID               int64   `json:"id"`
	RecipeID         int64   `json:"recipe_id,omitempty"`
	ProductName      string  `json:"product_name"`
	RequiredQuantity float64 `json:"required_quantity"`
	Unit             string  `json:"unit"`
}

type CreateRecipeRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Ingredients []RecipeIngredient `json:"ingredients"`
}

type IngredientAvailability struct {
	ProductName       string  `json:"product_name"`
	RequiredQuantity  float64 `json:"required_quantity"`
	AvailableQuantity float64 `json:"available_quantity"`
	Unit              string  `json:"unit"`
	IsEnough          bool    `json:"is_enough"`
}

type RecipeAvailability struct {
	RecipeID     int64                    `json:"recipe_id"`
	RecipeName   string                   `json:"recipe_name"`
	CanCook      bool                     `json:"can_cook"`
	MissingCount int                      `json:"missing_count"`
	Ingredients  []IngredientAvailability `json:"ingredients"`
}
