package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
)

// OrderItemRepository mendefinisikan operasi terhadap entitas OrderItem.
type OrderItemRepository interface {
    CreateOrderItem(ctx context.Context, item *domain.OrderItem) error
    GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error)
}
