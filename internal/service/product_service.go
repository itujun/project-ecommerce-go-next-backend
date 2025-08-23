package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
)

// ProductService menampung dependensi yang dibutuhkan.
type ProductService struct {
	productRepo repository.ProductRepository
	userRepo	repository.UserRepository
	validator	*validator.Validate
}

// NewProductService membuat instance ProductService baru.
func NewProductService(productRepo repository.ProductRepository, userRepo repository.UserRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		userRepo: userRepo,
		validator: validator.New(),
	}
}

// CreateProduct membuat produk baru.
func (s *ProductService) CreateProduct(ctx context.Context, sellerID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// validasi request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	// Periksa apakah penjual ada
	user, err := s.userRepo.GetUserByID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("seller tidak ditemukan")
	}
	// Hanya penjual atau admin yang boleh membuat produk.
	if user.Role.Name != "seller" && user.Role.Name != "admin" {
		return nil, fmt.Errorf("anda tidak memiliki izin untuk membuat produk")
	}
	// Generate slug unik
	prodSlug := slug.Make(req.Name)
	// Pastikan slug belum ada; jika ada, tambahkan suffix
	counter := 1
	for {
		existing, _ := s.productRepo.GetProductBySlug(ctx, prodSlug)
		if existing == nil {
			break
		}
		counter++
		prodSlug = fmt.Sprintf("%s-%d", slug.Make(req.Name), counter)
	}
	product := &domain.Product{
		ID:				uuid.New(),
		Name:			req.Name,
		Slug:			prodSlug,
		Description:	req.Description,
		Price:			req.Price,
		Image:			req.Image,
		Stock:			req.Stock,
		SellerID:		user.ID,
	}
	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	return &dto.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		SellerID:    product.SellerID.String(),
	}, nil
}

// GetProductByID mengembalikan detail produk.
func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		SellerID:    product.SellerID.String(),
	}, nil
}

// ListProducts mengembalikan semua produk.
func (s *ProductService) ListProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	var result []dto.ProductResponse
	for _, product := range products {
		result = append(result, dto.ProductResponse{
			ID:          product.ID.String(),
			Name:        product.Name,
			Slug:        product.Slug,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			SellerID:    product.SellerID.String(),
		})
	}
	return result, nil
}

// UpdateProduct memperbarui data produk.
func (s *ProductService) UpdateProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("produk tidak ditemukan")
	}
	// Pastikan hanya seller yang membuat produk ini atau admin yang bisa edit.
	seller, _ := s.userRepo.GetUserByID(ctx, sellerID)
	if seller.Role.Name != "admin" && product.SellerID != seller.ID {
		return nil, fmt.Errorf("anda tidak memiliki izin untuk mengubah produk ini")
	}
	product.Name = req.Name
	product.Slug = slug.Make(req.Name)
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Image = req.Image
	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}
	return &dto.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		SellerID:    product.SellerID.String(),
	}, nil
}

// DeleteProduct melakukan soft delete produk.
func (s *ProductService) DeleteProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID) error {
	prod, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("produk tidak ditemukan")
	} 
	seller, _ := s.userRepo.GetUserByID(ctx, sellerID)
	if seller.Role.Name != "admin" && prod.SellerID != seller.ID {
		return fmt.Errorf("anda tidak memiliki izin untuk menghapus produk ini")
	}
	return s.productRepo.DeleteProduct(ctx, id)
}