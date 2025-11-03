package models

import "time"

type FotoProduk struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	IDProduk  uint64     `gorm:"not null;index" json:"id_produk"`
	URL       string     `gorm:"type:varchar(255);not null" json:"url"` 
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Produk *Produk `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
}
