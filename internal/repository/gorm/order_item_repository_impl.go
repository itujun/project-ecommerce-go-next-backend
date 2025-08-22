package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"gorm.io/gorm"
)

// orderItemRepository adalah implementasi OrderItemRepository menggunakan GORM.
type orderItemRepository struct {
    db *gorm.DB
}

// NewOrderItemRepository membuat instance repository.
func NewOrderItemRepository(db *gorm.DB) repository.OrderItemRepository {
    return &orderItemRepository{db: db}
}

// CreateOrderItem menyimpan item pesanan ke database.
func (r *orderItemRepository) CreateOrderItem(ctx context.Context, item *domain.OrderItem) error {
    return r.db.WithContext(ctx).Create(item).Error
}

// GetItemsByOrderID mengambil semua item untuk order tertentu.
func (r *orderItemRepository) GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
    var items []domain.OrderItem
    err := r.db.WithContext(ctx).
        Preload("Product").
        Where("order_id = ?", orderID).
        Find(&items).Error
    return items, err
}