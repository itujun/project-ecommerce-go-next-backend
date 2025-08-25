package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService menyediakan logika bisnis terkait pengguna.
type UserService struct {
	userRepo	repository.UserRepository
	roleRepo 	repository.RoleRepository
	validator 	*validator.Validate
	jwtSecret	string
	jwtSvc *JWTService
	rtRepo repository.RefreshTokenRepository
}

// NewUserService membuat instance UserService baru.
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtSecret string) *UserService {
    return &UserService{
        userRepo:  userRepo,
        roleRepo:  roleRepo,
        validator: validator.New(),
        jwtSecret: jwtSecret,
    }
}

// RegisterUser membuat pengguna baru setelah validasi dan hasing password
func (s *UserService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error) {
	// Validasi input
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	// Pastikan email belum terdaftar
	if existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email); existing != nil {
		return nil, fmt.Errorf("email %s sudah terdaftar", req.Email)
	}
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat hash password: %w", err)
	}
	// Tentukan role default, misalnya "buyer"
	role, err := s.roleRepo.GetRoleByName(ctx, "buyer")
	if err != nil {
		return nil, fmt.Errorf("role default tidak ditemukan: %w", err)
	}
	user := &domain.User{
		ID:			uuid.New(),
		Name:		req.Name,
		Email:		req.Email,
		Password:	string(hashed),
		RoleID:		role.ID,
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	
	return &dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.Name,
	}, nil
}

// LoginUser memverifikasi kredensial dan mengembalikan token JWT.
func (s *UserService) LoginUser(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("email tidak terdaftar")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("email atau password salah")
	}
	// Buat klaim JWT
	claims := jwt.MapClaims{
		"user_id": 	user.ID.String(),
		"email":	user.Email,
		"role":		user.Role.Name,
		"exp":		time.Now().Add(24 * time.Hour).Unix(), // token berlaku 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("gagal menandatangani token: %w", err)
	}

	return &dto.LoginResponse{
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.Name,
		},
		Token: signed,
	},nil
}

func (s *UserService) SaveRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, issuedAt, expiresAt time.Time, jti string) error {
	// Simpan hash dari RT
	hash := s.jwtSvc.HashRefreshToken(refreshToken)
	rt := &domain.RefreshToken{
		ID:        uuid.MustParse(jti),
		UserID:    userID,
		TokenHash: hash,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	return s.rtRepo.Save(ctx, rt)
}

// VerifyRefreshTokenDB: validasi RT berdasar klaim JWT + cek DB (revoked/expired)
func (s *UserService) VerifyRefreshTokenDB(ctx context.Context, tokenStr string) (*domain.User, *uuid.UUID, error) {
	claims, err := s.jwtSvc.VerifyRefreshToken(tokenStr)
	if err != nil {
		return nil, nil, err
	}

	// Ambil jti & userID dari klaim
	jti := claims.ID
	if jti == "" {
		return nil, nil, errors.New("missing jti")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, nil, errors.New("invalid subject")
	}

	// Cek record RT di DB
	rtModel, err := s.rtRepo.FindByID(ctx, uuid.MustParse(jti))
	if err != nil {
		return nil, nil, errors.New("refresh token not found")
	}
	if rtModel.Revoked || time.Now().After(rtModel.ExpiresAt) {
		return nil, nil, errors.New("refresh token revoked or expired")
	}

	// Validasi hash (opsional—kalau ingin bind token ke DB)
	// hash := s.jwtSvc.HashRefreshToken(tokenStr)
	// if hash != rtModel.TokenHash { return nil, nil, errors.New("token mismatch") }

	// Ambil user (pastikan repository user Anda tersedia)
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	rtID := uuid.MustParse(jti)
	return user, &rtID, nil
}

// RotateRefreshToken: revoke RT lama & simpan RT baru
func (s *UserService) RotateRefreshToken(ctx context.Context, user *domain.User, oldRTID uuid.UUID, newRT string, issuedAt, expiresAt time.Time, newJTI string) error {
	// Revoke RT lama
	if err := s.rtRepo.Revoke(ctx, oldRTID); err != nil {
		return err
	}
	// Simpan RT baru
	return s.SaveRefreshToken(ctx, user.ID, newRT, issuedAt, expiresAt, newJTI)
}

// RevokeAllUserTokens: logout dari semua sesi
func (s *UserService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return s.rtRepo.RevokeAllByUser(ctx, userID)
}