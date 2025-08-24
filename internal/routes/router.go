package routes

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/itujun/project-ecommerce-go-next/internal/handler"
	"github.com/itujun/project-ecommerce-go-next/internal/middleware"
)

// NewRouter menginisialisasi router Chi dan mendaftarkan rute dasar.
// NewRouter menerima authHandler, productHandler, JWTMiddleware, dan enforcer Casbin.
func NewRouter(
    authHandler *handler.AuthHandler, 
    productHandler *handler.ProductHandler, 
    orderHandler *handler.OrderHandler, 
    jwtMiddleware *middleware.JWTMiddleware, 
    enforcer *casbin.Enforcer) *chi.Mux {
    r := chi.NewRouter()

    // Konfigurasi CORS
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // domain front‑end
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    })
    r.Use(corsHandler.Handler)
    
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

    // Order routes / Grup rute order
    r.Route("/orders", func(r chi.Router)  {
        // rute untuk create order: hanya pembeli (buyer) yang diizinkan
        r.Group(func(r chi.Router)  {
            r.Use(jwtMiddleware.Middleware)    
            r.Use(middleware.Authorize(enforcer, "order", "create"))
            r.Post("/", orderHandler.CreateOrder)
        })
        // rute untuk list orders: buyer melihat pesanan sendiri, admin & seller semua
        r.Group(func(r chi.Router)  {
            r.Use(jwtMiddleware.Middleware)
            r.Use(middleware.Authorize(enforcer, "order", "read"))
            r.Get("/", orderHandler.ListOrders)
        })
    })

    return r
}
