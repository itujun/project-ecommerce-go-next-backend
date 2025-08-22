package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
)

// RoleRepository mendefinisikan operasi untuk entitas Role.
type RoleRepository interface {
    CreateRole(ctx context.Context, role *domain.Role) error
    GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
    GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
    ListRoles(ctx context.Context) ([]domain.Role, error)
}