package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
)

// Authorize menerima enforcer, nama resource, dan action.
// Ia mengembalikan middleware yang memeriksa role user dari context.
// kemudian memanggil enforcer untuk mengecek izin.
func Authorize(enforcer *casbin.Enforcer, obj string, act string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil role user dari context (di-set oleh JWT middleware)
			roleVal := r.Context().Value("user_role")
			role, ok := roleVal.(string)
			if !ok || role == "" {
				http.Error(w, "Forbidden: role tidak ditemukan", http.StatusForbidden)
				return
			}

			// Panggil enforcer untuk memeriksa apakah role boleh melakukan act pada obj
			allowed, err := enforcer.Enforce(role, obj, act)
			if err != nil {
				http.Error(w, "Internal Server Error: gagal memeriksa izin", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "Forbidden: Anda tidak memiliki izin", http.StatusForbidden)
				return
			}

			// Jika diizinkan, lanjutkan ke handler selanjutnya
			next.ServeHTTP(w, r)
		})
	}
}

// Penjelasan kode
// - Fungsi Authorize mengembalikan middleware dinamis berdasarkan obj (resource) dan act (action). Parameter pertama adalah enforcer yang sudah diinisialisasi.
// - Middleware mengambil role dari context (di-set oleh JWT middleware).
// - Fungsi enforcer.Enforce(subject, object, action) akan mengembalikan true jika izin ada di file policy.
// - Jika tidak ada izin, middleware mengembalikan 403 Forbidden.