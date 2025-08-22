package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
	"gorm.io/gorm"
)

// productRepository adalah implementasi ProductRepository menggunakan GORM.
type productRepository struct {
    db *gorm.DB
}

// NewProductRepository membuat instance baru ProductRepository.
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

// CreateProduct menyimpan produk baru ke database.
func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}

// GetProductByID mengambil produk berdasarkan ID.
func (r *productRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
    var product domain.Product
    err := r.db.WithContext(ctx).Preload("Seller").First(&product, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}

// GetProductBySlug mengambil produk berdasarkan slug.
func (r *productRepository) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
    var product domain.Product
    err := r.db.WithContext(ctx).Preload("Seller").Where("slug = ?", slug).First(&product).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}

// ListProducts mengambil seluruh daftar produk.
func (r *productRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
    var products []domain.Product
    err := r.db.WithContext(ctx).Preload("Seller").Find(&products).Error
    return products, err
}

// UpdateProduct memperbarui data produk.
func (r *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
    return r.db.WithContext(ctx).Save(product).Error
}

// DeleteProduct menghapus (soft delete) produk berdasarkan ID.
func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
    // GORM akan mengisi kolom deleted_at sehingga data tidak benar-benar dihapus.
    return r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
}

// Penjelasan singkat:
// - Preload("Seller") digunakan untuk memuat relasi penjual ketika mengambil produk.
// - DeleteProduct menggunakan soft delete; data akan ditandai terhapus tetapi tetap ada di database