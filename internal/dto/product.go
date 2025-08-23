package dto

// CreateProductRequest mendefinisikan payload untuk membuat produk.
type CreateProductRequest struct {
	Name		string 	`json:"name" validate:"required,min=3,max=100"`
	Description string 	`json:"description"`
	Price		float64 `json:"price" validate:"required,gt=0"`
	Stock		int 	`json:"stock" validate:"required,gt=0"`
	Image		string 	`json:"image" validate:"required"` // nama file gambar
}

// UpdateProductRequest mendefinisikan payload untuk memperbarui produk.
type UpdateProductRequest struct {
	Name		string 	`json:"name" validate:"required,min=3,max=100"`
	Description string 	`json:"description"`
	Price		float64 `json:"price" validate:"required,gt=0"`
	Stock		int 	`json:"stock" validate:"required,gt=0"`
	Image		string 	`json:"image" validate:"required"`
}

// ProductResponse merepresentasikan data produk dalam response.
type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Image       string  `json:"image"`
	SellerID    string  `json:"seller_id"`
}