package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"gorm.io/gorm"
)

// userRepository adalah implementasi UserRepository menggunakan GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance repository.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepository{db: db}
}

// CreateUser menyimpan user baru ke database.
func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByEmail mencari user berdasarkan email.
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// GetUserByID mencari user berdasarkan ID.
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ListUsers mengambil semua user.
func (r *userRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
    var users []domain.User
    err := r.db.WithContext(ctx).Preload("Role").Find(&users).Error
    return users, err
}