package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itujun/project-ecommerce-go-next/internal/domain"
	"github.com/itujun/project-ecommerce-go-next/internal/dto"
	"github.com/itujun/project-ecommerce-go-next/internal/repository"
)

// OrderService menangani logika bisnis untuk pesanan.
type OrderService struct {
	orderRepo		repository.OrderRepository
	orderItemRepo	repository.OrderItemRepository
	productRepo		repository.ProductRepository
	userRepo		repository.UserRepository
	validator		*validator.Validate
}

// NewOrderService mengembalikan instance baru OrderService.
func NewOrderService(orderRepo repository.OrderRepository, orderItemRepo repository.OrderItemRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		orderItemRepo: orderItemRepo,
		productRepo: productRepo,
		userRepo: userRepo,
		validator: validator.New(),
	}
}

// CreateOrder membuat pesanan baru untuk pembeli.
func (s *OrderService) CreateOrder(ctx context.Context, buyerID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// validasi request
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	// Ambil data pembeli untuk verifikasi role(buyer)
	buyer, err := s.userRepo.GetUserByID(ctx, buyerID)
	if err != nil {
		return nil, fmt.Errorf("pembeli tidak ditemukan")
	}
	if buyer.Role.Name != "buyer" {
		return nil, fmt.Errorf("hanya pembeli yang dapat membuat pesanan")
	}

	// Hitung total harga dan siapkan items
	var total float64
	var items []domain.OrderItem
	for _, it := range req.Items {
		prodID, _ := uuid.Parse(it.ProductID)
		prod, err := s.productRepo.GetProductByID(ctx, prodID)
		if err != nil {
			return nil, fmt.Errorf("produk dengan ID %s todal ditemukan", it.ProductID)
		}
		if it.Quantity > prod.Stock {
			return nil, fmt.Errorf("stok produk %s tidak mencukupi", prod.Name)
		}
		//  Kurangi stok produk (optional: lakukan update stock di repository)
		prod.Stock -= it.Quantity
		if err := s.productRepo.UpdateProduct(ctx, prod); err != nil {
			return nil, fmt.Errorf("gagal memperbarui stok produk")
		}
		// Tambahkan ke item pesanan
		items = append(items, domain.OrderItem{
			ID:       	uuid.New(),
			ProductID:	prod.ID,
			Quantity: 	it.Quantity,
			Price:		prod.Price,
		})
		total += prod.Price * float64(it.Quantity)
	}

	// Buat Pesanan
	order := &domain.Order{
		ID:        uuid.New(),
		BuyerID:   buyer.ID,
		OrderDate: time.Now(),
		Total:     total,
		Status:    "pending",
	}
	// Simpan order utama
	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	// Hubungkan items dengan order_id baru
	for i := range items {
		items[i].OrderID = order.ID
		if err := s.orderItemRepo.CreateOrderItem(ctx, &items[i]); err != nil {
			return nil, err
		}
	}
	// Kembalikan response
	var respItems []dto.OrderItemResponse
	for _, it := range items {
		// Ambil nama produk
		prod, _ := s.productRepo.GetProductByID(ctx, it.ProductID)
		respItems = append(respItems, dto.OrderItemResponse{
			ID:			it.ID.String(),
			ProductID:	it.ProductID.String(),
			Quantity:	it.Quantity,
			Price:		it.Price,
			Name:		prod.Name,
		})
	}
	return &dto.OrderResponse{
		ID:			order.ID.String(),
		BuyerID:	buyer.ID.String(),
		OrderDate:	order.OrderDate.Format(time.RFC3339),
		Total:		order.Total,
		Status:		order.Status,
		Items:		respItems,
	}, nil
}

// Keterangan penting:
// - CreateOrder memvalidasi input, memeriksa role pembeli, menghitung total, mengurangi stok produk, lalu menyimpan order dan item ke database.
// - ListOrdersForBuyer mengembalikan pesanan milik pembeli tertentu.
// - ListAllOrdersAdminSeller mengembalikan semua pesanan; hanya dipanggil oleh admin/seller.
// - Anda dapat menambahkan metode untuk memperbarui status pesanan (dikirim, selesai, dll.) jika diperlukan.