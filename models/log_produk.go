package models

import "time"

type LogProduk struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	IDProduk      uint64     `gorm:"not null;index" json:"id_produk"`                        
	NamaProduk    string     `gorm:"type:varchar(150);not null" json:"nama_produk"`
	Slug          string     `gorm:"type:varchar(200);not null" json:"slug"`
	HargaReseller int        `gorm:"not null;default:0" json:"harga_reseller"`
	HargaKonsumen int        `gorm:"not null;default:0" json:"harga_konsumen"`
	Deskripsi     *string    `gorm:"type:text" json:"deskripsi,omitempty"`
	IDToko        uint64     `gorm:"not null;index" json:"id_toko"`
	IDCategory    *uint64    `gorm:"index" json:"id_category,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Toko     *Toko     `gorm:"foreignKey:IDToko" json:"toko,omitempty"`
	Category *Category `gorm:"foreignKey:IDCategory" json:"category,omitempty"`
}
