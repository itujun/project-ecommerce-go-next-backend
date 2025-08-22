package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
)

// ProductRepository mendefinisikan operasi CRUD untuk entitas Product.
type ProductRepository interface {
    CreateProduct(ctx context.Context, product *domain.Product) error
    GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
    GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error)
    ListProducts(ctx context.Context) ([]domain.Product, error)
    UpdateProduct(ctx context.Context, product *domain.Product) error
    DeleteProduct(ctx context.Context, id uuid.UUID) error
}