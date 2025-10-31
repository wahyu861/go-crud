package models

import "time"

type Category struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaCategory string     `gorm:"type:varchar(100);not null;unique" json:"nama_category"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Produk []Produk `gorm:"foreignKey:IDCategory;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"produk,omitempty"`
}
