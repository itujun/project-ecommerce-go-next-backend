package service

import (
	"context"
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