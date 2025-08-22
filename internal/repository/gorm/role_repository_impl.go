package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"gorm.io/gorm"
)

// roleRepository adalah implementasi RoleRepository.
type roleRepository struct {
    db *gorm.DB
}

// NewRoleRepository membuat instance repository.
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
    return &roleRepository{db: db}
}

// CreateRole menyimpan role baru.
func (r *roleRepository) CreateRole(ctx context.Context, role *domain.Role) error {
    return r.db.WithContext(ctx).Create(role).Error
}

// GetRoleByID mencari role berdasarkan ID.
func (r *roleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
    var role domain.Role
    err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &role, nil
}

// GetRoleByName mencari role berdasarkan nama.
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
    var role domain.Role
    err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
    if err != nil {
        return nil, err
    }
    return &role, nil
}

// ListRoles mengambil semua role.
func (r *roleRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
    var roles []domain.Role
    err := r.db.WithContext(ctx).Find(&roles).Error
    return roles, err
}
