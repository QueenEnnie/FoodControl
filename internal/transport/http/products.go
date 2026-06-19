package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"food-control/internal/domain"
)

const defaultExpiringDays = 3

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.repo.ListProducts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list products")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) listExpiringProducts(w http.ResponseWriter, r *http.Request) {
	days := defaultExpiringDays
	if rawDays := r.URL.Query().Get("days"); rawDays != "" {
		parsedDays, err := strconv.Atoi(rawDays)
		if err != nil || parsedDays < 0 || parsedDays > 365 {
			writeError(w, http.StatusBadRequest, "days must be an integer between 0 and 365")
			return
		}
		days = parsedDays
	}

	products, err := s.repo.ListExpiringProducts(r.Context(), days)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list expiring products")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := s.repo.GetProduct(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get product")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateProductRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := domain.ValidateProduct(input); err != nil {
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

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var input domain.CreateProductRequest

	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := domain.ValidateProduct(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := s.repo.UpdateProduct(r.Context(), id, input)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.repo.DeleteProduct(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
