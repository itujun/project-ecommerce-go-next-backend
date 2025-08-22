package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"gorm.io/gorm"
)

// orderRepository adalah implementasi OrderRepository menggunakan GORM.
type orderRepository struct {
    db *gorm.DB
}

// NewOrderRepository membuat instance OrderRepository.
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
    return &orderRepository{db: db}
}

// CreateOrder menyimpan pesanan baru beserta itemnya.
func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
    return r.db.WithContext(ctx).Create(order).Error
}

// GetOrderByID mengambil pesanan berdasarkan ID, lengkap dengan pembeli dan item.
func (r *orderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
    var order domain.Order
    err := r.db.WithContext(ctx).
        Preload("Buyer").
        Preload("Items").
        Preload("Items.Product").
        First(&order, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &order, nil
}

// ListOrdersByBuyer mengembalikan daftar pesanan milik pembeli tertentu.
func (r *orderRepository) ListOrdersByBuyer(ctx context.Context, buyerID uuid.UUID) ([]domain.Order, error) {
    var orders []domain.Order
    err := r.db.WithContext(ctx).
        Preload("Buyer").
        Preload("Items").
        Preload("Items.Product").
        Where("buyer_id = ?", buyerID).
        Find(&orders).Error
    return orders, err
}

// ListAllOrders mengembalikan semua pesanan (berguna untuk admin).
func (r *orderRepository) ListAllOrders(ctx context.Context) ([]domain.Order, error) {
    var orders []domain.Order
    err := r.db.WithContext(ctx).
        Preload("Buyer").
        Preload("Items").
        Preload("Items.Product").
        Find(&orders).Error
    return orders, err
}

// Catatan:
// - Preload("Items.Product") memuat produk di dalam setiap item, sehingga data pesanan lengkap terisi.