package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itujun/project-ecommerce-go-next/internal/handler"
)

// NewRouter menginisialisasi router Chi dan mendaftarkan rute dasar.
func NewRouter(authHandler *handler.AuthHandler) *chi.Mux {
    r := chi.NewRouter()

    // Contoh rute: GET /health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("OK"))
    })
    // Grup rute auth
    r.Route("/auth", func(r chi.Router) {
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
    })
    return r
}
