package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/service"
)

// OrderHandler menampung OrderService
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler mengembalikan instance baru OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder menangani POST /orders (pembeli membuat pesanan).
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// Ambil ID user dari context (setelah melewati JWT middleware)
	buyerIDStr, _ := r.Context().Value("user_id").(string)
	buyerID, err := uuid.Parse(buyerIDStr)
	if err != nil {
		http.Error(w, "invalid buyer ID", http.StatusUnauthorized)
		return
	}
	res, err := h.orderService.CreateOrder(r.Context(), buyerID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

// ListOrders menangani GET /orders untuk pembeli (menampilkan pesanan dirinya).
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value("role").(string)
	userIDStr, _ := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var res []dto.OrderResponse
	var err error
	if role == "buyer" {
		res, err = h.orderService.ListOrdersForBuyer(r.Context(), userID)
	}else {
		// admin atau seller melihat semua pesanan
		res, err = h.orderService.ListAllOrdersAdminSeller(r.Context())
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}