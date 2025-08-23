package dto

// OrderItemRequest merepresentasikan item yang diorder.
type OrderItemRequest struct {
	ProductID 	string	`json:"priduct_id" validate:"required"`		// ID produk dalam UUID
	Quantity 	int 	`json:"quantity" validate:"required,gt=0"`
}

// CreateOrderRequest merepresentasikan payload pembuatan pesanan/order.
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,dive"`	// Daftar item yang diorder. minimal satu item
}

// OrderItemResponse untuk mengembalikan data item pesanan/order.
type OrderItemResponse struct {
	ID			string	`json:"id"`
	ProductID	string	`json:"product_id"`
	Quantity	int		`json:"quantity"`
	Price		float64	`json:"price"` // Harga saat pembelian
	Name		string	`json:"name"`  // nama produk
}

// OrderResponse untuk mengembalikan detail pesanan.
type OrderResponse struct {
	ID			string				`json:"id"`
	BuyerID		string				`json:"buyer_id"`
	OrderDate	string				`json:"order_date"`
	Total		float64				`json:"total"`
	Status		string				`json:"status"`
	Items		[]OrderItemResponse	`json:"items"`
}

// Penjelasan:
// - OrderItemRequest berisi ProductID dan Quantity. Gunakan tag uuid4 untuk validasi UUID.
// - CreateOrderRequest memuat array item pesanan; minimal harus ada satu item.
// - OrderItemResponse menambahkan Name untuk menampilkan nama produk dan Price untuk harga saat pembelian.
// - OrderResponse menampilkan keseluruhan pesanan.