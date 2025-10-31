package models

import "time"

type Toko struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaToko  string     `gorm:"not null" json:"nama_toko"`
	UrlFoto   *string    `json:"url_foto"`             
	IDUser    uint64     `gorm:"not null" json:"id_user"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	User   *User     `gorm:"foreignKey:IDUser" json:"user,omitempty"`
	Produk []Produk  `gorm:"foreignKey:IDToko" json:"produk,omitempty"`
}
