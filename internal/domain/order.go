package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order menyimpan data pesanan pembeli.
type Order struct {
    ID        uuid.UUID   `gorm:"type:char(36);primaryKey" json:"id"`
    BuyerID   uuid.UUID   `gorm:"type:char(36);not null" json:"buyer_id"`
    Buyer     User        `gorm:"foreignKey:BuyerID" json:"buyer"`
    OrderDate time.Time   `gorm:"not null" json:"order_date"`
    Total     float64     `gorm:"type:decimal(10,2);not null" json:"total"`
    Status    string      `gorm:"size:50;not null" json:"status"`
    Items     []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
    gorm.Model
}
