package httpapi

import (
	"log/slog"
	"net/http"
)

type Server struct {
	repo    Repository
	logger  *slog.Logger
	mux     *http.ServeMux
	handler http.Handler
}

func New(repo Repository, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	server := &Server{
		repo:   repo,
		logger: logger,
		mux:    http.NewServeMux(),
	}
	server.routes()
	server.handler = server.withMiddleware(server.mux)
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /products", s.listProducts)
	s.mux.HandleFunc("GET /products/expiring", s.listExpiringProducts)
	s.mux.HandleFunc("GET /products/{id}", s.getProduct)
	s.mux.HandleFunc("POST /products", s.createProduct)
	s.mux.HandleFunc("PUT /products/{id}", s.updateProduct)
	s.mux.HandleFunc("DELETE /products/{id}", s.deleteProduct)
	s.mux.HandleFunc("GET /recipes", s.listRecipes)
	s.mux.HandleFunc("POST /recipes", s.createRecipe)
	s.mux.HandleFunc("POST /recipes/{id}/cook", s.cookRecipe)
	s.mux.HandleFunc("GET /recipes/{id}/availability", s.recipeAvailability)
	s.mux.HandleFunc("GET /suggestions", s.recipeSuggestions)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
