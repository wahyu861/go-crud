package models

import "time"

type DetailTrx struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	IDTrx       uint64     `gorm:"not null;index" json:"id_trx"`           
	IDLogProduk uint64     `gorm:"not null;index" json:"id_log_produk"`    
	IDToko      uint64     `gorm:"not null;index" json:"id_toko"`          
	Kuantitas   int        `gorm:"not null;default:1" json:"kuantitas"`
	HargaTotal  int        `gorm:"not null;default:0" json:"harga_total"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Trx       *Trx       `gorm:"foreignKey:IDTrx" json:"trx,omitempty"`
	LogProduk *LogProduk `gorm:"foreignKey:IDLogProduk" json:"log_produk,omitempty"`
	Toko      *Toko      `gorm:"foreignKey:IDToko" json:"toko,omitempty"`
}
