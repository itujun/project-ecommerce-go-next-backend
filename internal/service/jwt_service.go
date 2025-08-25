package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/config"
)

// CustomClaims menyimpan data user minimal + jti
type CustomClaims struct {
	UserID uuid.UUID `json:"uid"`
	RoleID uuid.UUID `json:"rid"`
	jwt.RegisteredClaims
}

// JWTService menyediakan util untuk generate/verify token
type JWTService struct {
	cfg *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{cfg: cfg}
}

// GenerateAccessToken membuat AT (durasi pendek) untuk akses API
// semula: func (s *JWTService) GenerateAccessToken(u *domain.User) (string, time.Time, error)
// ganti jadi menerima primitive:
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, roleID uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.cfg.AccessTTL)

	claims := CustomClaims{
		 UserID: userID,
        RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecommerce-go",
			Subject:   userID.String(),
			Audience:  []string{"ecommerce-client"},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(), // jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(s.cfg.JWTAccessSecret))
	return str, exp, err
}

// GenerateRefreshToken membuat RT (durasi lebih panjang) untuk refresh AT
// semula: func (s *JWTService) GenerateRefreshToken(u *domain.User) (string, time.Time, error)
// ganti jadi menerima hanya userID (role tidak perlu untuk RT)
func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.cfg.RefreshTTL)

	claims := jwt.RegisteredClaims{
		Issuer:    "ecommerce-go",
		Subject:   userID.String(),
		Audience:  []string{"ecommerce-client"},
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ID:        uuid.NewString(), // jti untuk lacak RT di DB
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(s.cfg.JWTRefreshSecret))
	return str, exp, err
}

// VerifyAccessToken memverifikasi & mengembalikan claims
func (s *JWTService) VerifyAccessToken(tokenStr string) (*CustomClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid access token")
	}
	claims, ok := tok.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

// VerifyRefreshToken mengembalikan klaim RegisteredClaims
func (s *JWTService) VerifyRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTRefreshSecret), nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid refresh token")
	}
	claims, ok := tok.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid refresh claims")
	}
	return claims, nil
}

// HashRefreshToken meng‑hash RT untuk disimpan di DB (jangan simpan plaintext)
func (s *JWTService) HashRefreshToken(rt string) string {
	sum := sha256.Sum256([]byte(rt))
	return hex.EncodeToString(sum[:])
}
