package routes

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/itujun/project-ecommerce-go-next/internal/handler"
	"github.com/itujun/project-ecommerce-go-next/internal/middleware"
)

// NewRouter menginisialisasi router Chi dan mendaftarkan rute dasar.
// NewRouter menerima authHandler, productHandler, JWTMiddleware, dan enforcer Casbin.
func NewRouter(authHandler *handler.AuthHandler, productHandler *handler.ProductHandler, jwtMiddleware *middleware.JWTMiddleware, enforcer *casbin.Enforcer) *chi.Mux {
    r := chi.NewRouter()

    // Contoh rute: GET /health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("OK"))
    })
    // Authentication routes / Grup rute auth
    r.Route("/auth", func(r chi.Router) {
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
    })
    // Product routes / Grup rute product
    r.Route("/products", func(r chi.Router) {
        r.Get("/", productHandler.ListProducts)     // publik
        r.Get("/{id}", productHandler.GetProduct)   // publik
        // Endpoints di bawah ini dilindungi JWT dan Casbin.
        r.Group(func(r chi.Router)  {
            r.Use(jwtMiddleware.Middleware)                             // parse token
            r.Use(middleware.Authorize(enforcer, "product", "create"))  // role cek
            r.Post("/", productHandler.CreateProduct)
        })
        r.Group(func(r chi.Router)  {
            r.Use(jwtMiddleware.Middleware)                             // parse token
            r.Use(middleware.Authorize(enforcer, "product", "update"))  // role cek
            r.Put("/{id}", productHandler.UpdateProduct)
        })
        r.Group(func(r chi.Router)  {
            r.Use(jwtMiddleware.Middleware)                             // parse token
            r.Use(middleware.Authorize(enforcer, "product", "delete"))  // role cek
            r.Delete("/{id}", productHandler.DeleteProduct)
        })
        // Di sini, Authorize membutuhkan dua parameter: nama resource (product) dan action (create, update, delete). Peran (role) pengguna diambil dari token, kemudian dicek terhadap policy Casbin.
    })
    return r
}
