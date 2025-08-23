package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware menampung secret yang dipakai untuk memverifikasi token.
// Dengan cara ini, kita bisa membuat instance baru ketika aplikasi dijalankan.
type JWTMiddleware struct {
	secret string
}

// NewJWTMiddleware mengembalikan instance JWTMiddleware baru.
func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{secret: secret}
}

// Middleware adalah fungsi actual yang akan dipasang di router.
// Ia memeriksa token, memverifikasi, kemuadian menaruh data user di context.
func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Baca header Authorization: harus berupa "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: token tidak ditemukan", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: format token salah", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		
		// Parse dan verifikasi token menggunakan secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Pastikan algoritma HS256 digunakan
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(m.secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: token tidak valid", http.StatusUnauthorized)
			return
		}

		// Ambil klaim dari token; kita asumsikan menggunakan MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: klaim tidak valid", http.StatusUnauthorized)
			return
		}

		// Ekstrak data user dari klaim
		userID, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		// Masukkan data ke context agar bisa diakses di handler/service
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "role", role)

		// Lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Penjelasan kode
// - JWTMiddleware menyimpan secret dari konfigurasi (JWT_SECRET). Ini dibutuhkan untuk memverifikasi token.
// - Fungsi Middleware adalah middleware actual. Ia membaca header Authorization, memeriksa format Bearer <token>, lalu memverifikasi tanda tangan dengan secret.
// - Jika token valid, klaim diambil dan disimpan di context (user_id, email, role). Handler berikutnya dapat mengambil nilai ini menggunakan r.Context().Value("user_id").
// - Jika token tidak valid atau tidak ada, middleware merespons 401 Unauthorized.