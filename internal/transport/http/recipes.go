package httpapi

import (
	"errors"
	"net/http"

	"food-control/internal/domain"
)

func (s *Server) listRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := s.repo.ListRecipes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list recipes")
		return
	}
	writeJSON(w, http.StatusOK, recipes)
}

func (s *Server) createRecipe(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateRecipeRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := domain.ValidateRecipe(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := s.repo.CreateRecipe(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create recipe")
		return
	}
	writeJSON(w, http.StatusCreated, recipe)
}

func (s *Server) cookRecipe(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.repo.CookRecipe(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}
	if errors.Is(err, domain.ErrInsufficientIngredients) {
		writeError(w, http.StatusConflict, "not enough ingredients to cook recipe")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to cook recipe")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) recipeAvailability(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	availability, err := s.repo.GetRecipeAvailability(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check recipe availability")
		return
	}
	writeJSON(w, http.StatusOK, availability)
}

func (s *Server) recipeSuggestions(w http.ResponseWriter, r *http.Request) {
	suggestions, err := s.repo.ListRecipeSuggestions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to build suggestions")
		return
	}
	writeJSON(w, http.StatusOK, suggestions)
}
