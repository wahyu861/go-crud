package models

import "time"

type Alamat struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser        uint64     `gorm:"not null;index" json:"id_user"`
	JudulAlamat   string     `gorm:"type:varchar(100);not null" json:"judul_alamat"`
	NamaPenerima  string     `gorm:"type:varchar(100);not null" json:"nama_penerima"`
	NoTelp        string     `gorm:"type:varchar(20);not null" json:"no_telp"`
	DetailAlamat  string     `gorm:"type:text;not null" json:"detail_alamat"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	User *User `gorm:"foreignKey:IDUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}
