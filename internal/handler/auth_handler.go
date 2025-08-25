package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/service"
	"github.com/itujun/project-ecommerce-go-next/internal/utils"
)

// AuthHandler menagani request registrasi dan login.
type AuthHandler struct {
	userService *service.UserService
	jwtService  *service.JWTService // <-- tambahkan JWT service
}

// NewAuthHandler membuat instance baru AuthHandler
func NewAuthHandler(userService *service.UserService, jwtService *service.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// helper: tulis JSON respons
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Register menangani POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	res, err := h.userService.RegisterUser(context.Background(), req)
	if err != nil {
		// cek apakah err adalah ValidatonErrors
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// terjemahkan ke map
			fieldErrors := utils.ValidationErrorsToMap(ve)
			writeJSON(w, http.StatusBadRequest, fieldErrors)
			return
		}
		// error lain
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

// Login menangani POST /auth/login.
// Login menangani POST /auth/login.
// Flow:
// 1) Validasi kredensial via userService.LoginUser
// 2) Generate AT & RT via jwtService
// 3) Simpan RT (hash) ke DB via userService.SaveRefreshToken
// 4) Set cookie HttpOnly untuk AT & RT
// 5) Kembalikan data user (tanpa token)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
	// Dekode body JSON ke struct LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

	// Panggil service untuk login (cek email/password)
    res, err := h.userService.LoginUser(context.Background(), req)
    if err != nil {
        // Jika error adalah validasi field
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            // Terjemahkan validation errors ke map[field]pesan
            fieldErrors := utils.ValidationErrorsToMap(ve)
            writeJSON(w, http.StatusBadRequest, fieldErrors)
            return
        }
        // Jika error bukan validasi (contoh: email atau password salah)
        writeJSON(w, http.StatusUnauthorized, map[string]string{
			"general": err.Error(), // misalnya: "email atau password salah"
		})
        return
    }

	// Ambil ID & Role dari DTO
	uid, err := uuid.Parse(res.User.ID) // asumsi dto.UserResponse.ID berupa string UUID
	if err != nil { http.Error(w, "invalid user id", http.StatusInternalServerError); return }

	var rid uuid.UUID
	if res.User.Role != "" {
		rid, err = uuid.Parse(res.User.Role)
		if err != nil { http.Error(w, "invalid role id", http.StatusInternalServerError); return }
	} else {
		// Jika DTO tidak memuat RoleID, ambil dari DB
		uDomain, err := h.userService.GetUserByID(r.Context(), uid)
		if err != nil { http.Error(w, "user not found", http.StatusInternalServerError); return }
		rid = uDomain.RoleID
	}

    // --- Generate Access Token & Refresh Token ---
	atStr, atExp, err := h.jwtService.GenerateAccessToken(uid, rid)
	if err != nil { http.Error(w, "cannot issue access token", http.StatusInternalServerError); return }

	rtStr, rtExp, err := h.jwtService.GenerateRefreshToken(uid)
	if err != nil { http.Error(w, "cannot issue refresh token", http.StatusInternalServerError); return }

	// Ambil klaim RT (untuk issuedAt, expiresAt, jti)
	// Pars RT untuk ambil jti/issuedAt/expiresAt → simpan hash RT di DB
	rtClaims, err := h.jwtService.VerifyRefreshToken(rtStr)
	if err != nil { 
		http.Error(w, "invalid refresh token claims", http.StatusInternalServerError); return 
	}

	issuedAt := time.Unix(rtClaims.IssuedAt.Unix(), 0)
	expiresAt := time.Unix(rtClaims.ExpiresAt.Unix(), 0)
	jti := rtClaims.ID

	// Simpan hash RT ke DB (best practice: jangan simpan plaintext)
	if err := h.userService.SaveRefreshToken(r.Context(), uid, rtStr, issuedAt, expiresAt, jti); err != nil {
    	http.Error(w, "cannot persist refresh token", http.StatusInternalServerError)
    	return
	}

	// --- Set cookie HttpOnly untuk AT & RT ---
	// NOTE: Secure=false untuk dev HTTP lokal; set true saat produksi (HTTPS)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    atStr,
		Path:     "/",
		Expires:  atExp,
		MaxAge:   int(time.Until(atExp).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true di produksi (HTTPS)
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rtStr,
		Path:     "/auth/refresh", // cakupan sempit; bisa "/" jika diinginkan
		Expires:  rtExp,
		MaxAge:   int(time.Until(rtExp).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // true di produksi (HTTPS)
	})

	// Kembalikan info user (tanpa token di body)
	writeJSON(w, http.StatusOK, map[string]any{
		"user":    res.User,
		"message": "login success",
	})
}

// Refresh menangani POST /auth/refresh.
// Flow:
// 1) Ambil refresh_token dari cookie
// 2) Verifikasi JWT RT + cek DB (not revoked, not expired)
// 3) Generate AT baru + RT baru (rotasi RT)
// 4) Revoke RT lama, simpan RT baru
// 5) Set cookie baru
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	rtCookie, err := r.Cookie("refresh_token")
	if err != nil || rtCookie.Value == "" {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	// Verifikasi RT di JWT & DB (dapatkan user & RT-ID lama)
	user, oldRTID, err := h.userService.VerifyRefreshTokenDB(r.Context(), rtCookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Buat AT & RT baru
	atStr, atExp, err := h.jwtService.GenerateAccessToken(user.ID, user.RoleID)
	if err != nil {
		http.Error(w, "cannot issue access token", http.StatusInternalServerError)
		return
	}
	newRTStr, _, err := h.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "cannot issue refresh token", http.StatusInternalServerError)
		return
	}

	// Ambil klaim RT baru (issuedAt, expiresAt, jti)
	newRTClaims, err := h.jwtService.VerifyRefreshToken(newRTStr)
	if err != nil {
		http.Error(w, "invalid new refresh token claims", http.StatusInternalServerError)
		return
	}
	newIssuedAt := time.Unix(newRTClaims.IssuedAt.Unix(), 0)
	newExpiresAt := time.Unix(newRTClaims.ExpiresAt.Unix(), 0)
	newJTI := newRTClaims.ID

	// Rotasi RT: revoke RT lama, simpan RT baru
	if err := h.userService.RotateRefreshToken(r.Context(), user, *oldRTID, newRTStr, newIssuedAt, newExpiresAt, newJTI); err != nil {
		http.Error(w, "cannot rotate refresh token", http.StatusInternalServerError)
		return
	}

	// Set cookie baru
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    atStr,
		Path:     "/",
		Expires:  atExp,
		MaxAge:   int(time.Until(atExp).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true di produksi
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRTStr,
		Path:     "/auth/refresh",
		Expires:  newExpiresAt,
		MaxAge:   int(time.Until(newExpiresAt).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // true di produksi
	})

	// Tidak perlu body; 204 cukup
	w.WriteHeader(http.StatusNoContent)
}

// Logout menangani POST /auth/logout.
// Flow:
// - Hapus cookie AT & RT (expire).
// - (Opsional) Revoke semua RT milik user di DB agar logout dari semua sesi.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Coba identifikasi user dari AT untuk revoke semua RT (opsional)
	if atCookie, err := r.Cookie("access_token"); err == nil && atCookie.Value != "" {
		if claims, err := h.jwtService.VerifyAccessToken(atCookie.Value); err == nil {
			_ = h.userService.RevokeAllUserTokens(r.Context(), claims.UserID)
		}
	}

	// Hapus cookie dengan MaxAge negatif
	del := func(name, path string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false, // true di produksi
		})
	}
	del("access_token", "/")
	del("refresh_token", "/auth/refresh")

	w.WriteHeader(http.StatusNoContent)
}

// Me menangani GET /auth/me.
// Verifikasi AT dari cookie lalu kembalikan info user.
// Berguna untuk FE mengecek status login tanpa menyentuh cookie secara langsung.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	atCookie, err := r.Cookie("access_token")
	if err != nil || atCookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	claims, err := h.jwtService.VerifyAccessToken(atCookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}