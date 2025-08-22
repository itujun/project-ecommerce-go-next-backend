package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderItem menyimpan rincian produk dalam pesanan.
type OrderItem struct {
    ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    OrderID   uuid.UUID `gorm:"type:char(36);not null" json:"order_id"`
    Order     Order     `gorm:"foreignKey:OrderID" json:"-"`       // tidak diserialisasi untuk menghindari loop
    ProductID uuid.UUID `gorm:"type:char(36);not null" json:"product_id"`
    Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
    Quantity  int       `gorm:"not null" json:"quantity"`
    Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
    gorm.Model
}
