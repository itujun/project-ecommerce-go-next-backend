package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter menginisialisasi router Chi dan mendaftarkan rute dasar.
func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // Contoh rute: GET /health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("OK"))
    })
    return r
}
