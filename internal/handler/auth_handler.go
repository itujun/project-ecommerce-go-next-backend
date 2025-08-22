package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/service"
)

// AuthHandler menagani request registrasi dan login.
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register menangani POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res, err := h.userService.RegisterUser(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

// Login menangani POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    res, err := h.userService.LoginUser(context.Background(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(res)
}