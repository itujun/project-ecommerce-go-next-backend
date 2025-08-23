package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/itujun/project-ecommerce-go-next/internal/authorization"
	"github.com/itujun/project-ecommerce-go-next/internal/config"
	"github.com/itujun/project-ecommerce-go-next/internal/database"
	"github.com/itujun/project-ecommerce-go-next/internal/handler"
	"github.com/itujun/project-ecommerce-go-next/internal/middleware"
	"github.com/itujun/project-ecommerce-go-next/internal/repository/gorm"
	"github.com/itujun/project-ecommerce-go-next/internal/routes"
	"github.com/itujun/project-ecommerce-go-next/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Inisialisasi logger menggunakan Zap
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("❌gagal membuat logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()	// flush buffer log sebelum keluar
	}()

	// Muat konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("❌gagal memuat konfigurasi", zap.Error(err))
	}

	// Buka koneksi database
	db, err := database.NewMySQL(cfg)
	if err != nil {
		logger.Fatal("❌gagal koneksi ke database", zap.Error(err))
	}
	_ = db // variable db belum digunakan pada langkah ini; nanti akan diteruskan ke repository

	// Inisialisasi repository dan service
	userRepo	:= gorm.NewUserRepository(db)
	roleRepo	:= gorm.NewRoleRepository(db)
	userService := service.NewUserService(userRepo, roleRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(userService)

	// Inisialisasi enforcer Casbin
    enforcer, err := authorization.NewEnforcer("config/rbac_model.conf", "config/rbac_policy.csv")
    if err != nil {
        logger.Fatal("gagal inisialisasi Casbin", zap.Error(err))
    }
	
	// Inisialisasi JWT middleware
    jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)
	
	// Inisialisasi repository dan service
    productRepo 	:= gorm.NewProductRepository(db)
	orderRepo 		:= gorm.NewOrderRepository(db)
    orderItemRepo 	:= gorm.NewOrderItemRepository(db)
    productService 	:= service.NewProductService(productRepo, userRepo)
	orderService 	:= service.NewOrderService(orderRepo, orderItemRepo, productRepo, userRepo)
    productHandler 	:= handler.NewProductHandler(productService)
	orderHandler 	:= handler.NewOrderHandler(orderService)
	
	// Router dengan authHandler (dari langkah 3), productHandler, jwtMiddleware, enforcer
    router := routes.NewRouter(authHandler, productHandler, orderHandler, jwtMiddleware, enforcer)

	// Jalankan server HTTP
	logger.Info("✅server dijalankan", zap.String("port", cfg.AppPort))
	if err := http.ListenAndServe(cfg.AppPort, router); err != nil {
		logger.Fatal("❌gagal menjalankan server", zap.Error(err))
	}
}