package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"food-control/internal/models"
	"food-control/internal/repository"
)

type Server struct {
	repo *repository.Repository
	mux  *http.ServeMux
}

func New(repo *repository.Repository) *Server {
	server := &Server{
		repo: repo,
		mux:  http.NewServeMux(),
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /products", s.listProducts)
	s.mux.HandleFunc("POST /products", s.createProduct)
	s.mux.HandleFunc("PUT /products/{id}", s.updateProduct)
	s.mux.HandleFunc("DELETE /products/{id}", s.deleteProduct)
	s.mux.HandleFunc("GET /recipes", s.listRecipes)
	s.mux.HandleFunc("POST /recipes", s.createRecipe)
	s.mux.HandleFunc("GET /recipes/{id}/availability", s.recipeAvailability)
	s.mux.HandleFunc("GET /suggestions", s.recipeSuggestions)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.repo.ListProducts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list products")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var input models.CreateProductRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := repository.ValidateProduct(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := s.repo.CreateProduct(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create product")
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.repo.DeleteProduct(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var input models.CreateProductRequest

	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := repository.ValidateProduct(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := s.repo.UpdateProduct(r.Context(), id, input)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) listRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := s.repo.ListRecipes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list recipes")
		return
	}
	writeJSON(w, http.StatusOK, recipes)
}

func (s *Server) createRecipe(w http.ResponseWriter, r *http.Request) {
	var input models.CreateRecipeRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := repository.ValidateRecipe(input); err != nil {
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

func (s *Server) recipeAvailability(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	availability, err := s.repo.GetRecipeAvailability(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
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

func parseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}

func readJSON(r *http.Request, dst any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") && r.ContentLength != 0 {
		return errors.New("content type must be application/json")
	}
	if r.Body == nil {
		return errors.New("request body is required")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return errors.New("invalid JSON body")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
