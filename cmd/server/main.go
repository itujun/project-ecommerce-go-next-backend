package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/itujun/project-ecommerce-go-next/internal/config"
	"github.com/itujun/project-ecommerce-go-next/internal/database"
	"github.com/itujun/project-ecommerce-go-next/internal/handler"
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

	// Inisialisasi router dengan handler auth
	router := routes.NewRouter(authHandler)

	// Jalankan server HTTP
	logger.Info("✅server dijalankan", zap.String("port", cfg.AppPort))
	if err := http.ListenAndServe(cfg.AppPort, router); err != nil {
		logger.Fatal("❌gagal menjalankan server", zap.Error(err))
	}
}