package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
)

// OrderRepository mendefinisikan operasi dasar terhadap entitas Order.
type OrderRepository interface {
    CreateOrder(ctx context.Context, order *domain.Order) error
    GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
    ListOrdersByBuyer(ctx context.Context, buyerID uuid.UUID) ([]domain.Order, error)
    ListAllOrders(ctx context.Context) ([]domain.Order, error)
}