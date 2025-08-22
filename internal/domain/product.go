package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product merepresentasikan produk di toko.
type Product struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    Name        string    `gorm:"size:255;not null" json:"name"`
    Slug        string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
    Description string    `gorm:"type:text" json:"description"`
    Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
    Image       string    `gorm:"size:255" json:"image"`
    Stock       int       `gorm:"not null" json:"stock"`
    SellerID    uuid.UUID `gorm:"type:char(36);not null" json:"seller_id"`
    Seller      User      `gorm:"foreignKey:SellerID" json:"seller"`
    gorm.Model
}
