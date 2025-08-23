package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/service"
)

// ProductHandler menampung ProductService.
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler membuat instance handler baru.
func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: service}
}

// ListProducts menangani GET /products.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	res, err := h.productService.ListProducts(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// GetProduct menangani GET /products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        http.Error(w, "invalid product id", http.StatusBadRequest)
        return
    }
    res, err := h.productService.GetProductByID(context.Background(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(res)
}

// CreateProduct menangani POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    // Dapatkan sellerID dari JWT (middleware menaruhnya di context)
    sellerID := r.Context().Value("user_id").(string)
    uid, _ := uuid.Parse(sellerID)
    res, err := h.productService.CreateProduct(r.Context(), uid, req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(res)
}

// UpdateProduct menangani PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
    var req dto.UpdateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    idParam := chi.URLParam(r, "id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        http.Error(w, "invalid product id", http.StatusBadRequest)
        return
    }
    sellerID := r.Context().Value("user_id").(string)
    uid, _ := uuid.Parse(sellerID)
    res, err := h.productService.UpdateProduct(r.Context(), uid, id, req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(res)
}

// DeleteProduct menangani DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        http.Error(w, "invalid product id", http.StatusBadRequest)
        return
    }
    sellerID := r.Context().Value("user_id").(string)
    uid, _ := uuid.Parse(sellerID)
    if err := h.productService.DeleteProduct(r.Context(), uid, id); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Endpoint yang memodifikasi produk (POST, PUT, DELETE) memerlukan ID pengguna yang diambil dari context (set oleh middleware JWT) dan memeriksa peran di service.