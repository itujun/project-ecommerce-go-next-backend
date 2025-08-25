package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, rt *domain.RefreshToken) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
	Update(ctx context.Context, rt *domain.RefreshToken) error
}
