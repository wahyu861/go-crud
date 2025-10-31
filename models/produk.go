package models

import "time"

type Produk struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaProduk     string     `gorm:"type:varchar(150);not null;index" json:"nama_produk"` 
	Slug           string     `gorm:"type:varchar(200);unique;not null" json:"slug"`       
	HargaReseller  int        `gorm:"not null;default:0" json:"harga_reseller"`
	HargaKonsumen  int        `gorm:"not null;default:0" json:"harga_konsumen"`
	Stok           int        `gorm:"not null;default:0" json:"stok"`
	Deskripsi      *string    `gorm:"type:text" json:"deskripsi,omitempty"`                
	IDToko         uint64     `gorm:"not null;index" json:"id_toko"`                       
	IDCategory     *uint64    `gorm:"index" json:"id_category,omitempty"`                  
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Toko        *Toko         `gorm:"foreignKey:IDToko" json:"toko,omitempty"`
	Category    *Category     `gorm:"foreignKey:IDCategory" json:"category,omitempty"`
	FotoProduk  []FotoProduk  `gorm:"foreignKey:IDProduk" json:"foto_produk,omitempty"`
	LogProduk   []LogProduk   `gorm:"foreignKey:IDProduk" json:"log_produk,omitempty"`
}
