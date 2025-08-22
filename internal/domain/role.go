package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role menampung jenis peran (super_admin, admin, penjual, pembeli).
type Role struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
    Description string    `gorm:"size:255" json:"description"`
    gorm.Model
    Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"` // satu role dimiliki banyak user
}
