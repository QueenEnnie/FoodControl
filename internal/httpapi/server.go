package httpapi

import (
	"net/http"

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
	s.mux.HandleFunc("GET /products/{id}", s.getProduct)
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
