package database

import (
	"fmt"

	"github.com/itujun/project-ecommerce-go-next/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQL membuka koneksi databse MySQL menggunakan GORM.
// Mengembalikan pointer *gorm.DB untuk digunakan oleh repository
func NewMySQL(cfg *config.Config) (*gorm.DB, error)  {
	// Format DSN (Data Source Name) sesuai MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBCharset,
    )
    // Membuka koneksi menggunakan GORM
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
    }
	return db, nil
}