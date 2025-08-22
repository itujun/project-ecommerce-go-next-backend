package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan entitas pengguna dalam aplikasi.
type User struct {
    ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    Name     string    `gorm:"size:255;not null" json:"name"`
    Email    string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
    Password string    `gorm:"size:255;not null" json:"-"`         // password hash, tidak diekspor ke JSON
    RoleID   uuid.UUID `gorm:"type:char(36);not null" json:"role_id"`
    Role     Role      `gorm:"foreignKey:RoleID" json:"role"`      // relasi ke Role
    gorm.Model // menyertakan CreatedAt, UpdatedAt, DeletedAt (untuk soft delete)
}
