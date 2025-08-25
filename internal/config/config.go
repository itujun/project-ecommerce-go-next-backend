package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config menampung seluruh konfigurasi alikasi yang dibaca dari file .env
type Config struct {
	AppPort		string // port aplikasi HTTP, contoh ":8080"
	DBUser		string // username db MySQL
	DBPassword	string // password db MySQL
	DBHost		string // host db, contoh "localhost"
	DBPort		string // port db, contoh "3306"
	DBName		string // nama db 
	DBCharset	string // character set db
	JWTSecret   string // Secret untuk menandatangani token JWT
	JWTAccessSecret   	string 			// secret HMAC untuk AT
	JWTRefreshSecret	string 			// secret HMAC untuk RT
	AccessTTL			time.Duration 	// durasi AT, mis. 15m
	RefreshTTL			time.Duration 	// durasi RT, mis. 168h (7d)
}

// LoadConfig membaca konfigurasi file .env dan environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")				// tentukan file konfigurasi
	viper.SetDefault("APP_PORT", ":8080")	// nilai default jika tidak diset
	viper.SetDefault("DB_CHARSET", "utf8mb4")
	viper.SetDefault("JWT_SECRET", "supersecret")
	viper.SetDefault("JWT_ACCESS_SECRET", "super-at-secret")
	viper.SetDefault("JWT_REFRESH_SECRET", "super-rt-secret")
	viper.SetDefault("JWT_ACCESS_TTL", "30m")
	viper.SetDefault("JWT_REFRESH_TTL", "72h") // 7 hari

	// Membaca file .env (jika ada)
	if err := viper.ReadInConfig(); err != nil {
		// Tidak menemukan file .env bukanlah error fatal; variabel environment tetap bisa digunakan.
		fmt.Printf("tidak menemunkan file .env: %v\n", err)
	}

	viper.AutomaticEnv() // override dengan environtment variable

	accessTTL, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TTL"))
	if err != nil { return nil, err}
	refreshTTL, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TTL"))
	if err != nil { return nil, err}

	cfg := &Config{
		AppPort: 	viper.GetString("APP_PORT"),
		DBUser: 	viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
        DBHost:     viper.GetString("DB_HOST"),
        DBPort:     viper.GetString("DB_PORT"),
        DBName:     viper.GetString("DB_NAME"),
        DBCharset:  viper.GetString("DB_CHARSET"),
        JWTSecret:  viper.GetString("JWT_SECRET"),
		JWTAccessSecret: viper.GetString("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
		AccessTTL: accessTTL,
		RefreshTTL: refreshTTL,
	}
	return cfg,nil
}