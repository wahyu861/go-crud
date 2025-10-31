package models

import "time"

type Trx struct {
	ID               uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser           uint64      `gorm:"not null" json:"id_user"`
	AlamatPengiriman *uint64     `gorm:"index" json:"alamat_pengiriman,omitempty"`
	HargaTotal       int         `gorm:"not null;default:0" json:"harga_total"`
	KodeInvoice      string      `gorm:"type:varchar(50);unique;not null" json:"kode_invoice"`
	MethodBayar      *string     `gorm:"type:varchar(50)" json:"method_bayar,omitempty"`
	CreatedAt        time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	User      *User        `gorm:"foreignKey:IDUser" json:"user,omitempty"`
	Alamat    *Alamat      `gorm:"foreignKey:AlamatPengiriman" json:"alamat,omitempty"`
	DetailTrx []DetailTrx  `gorm:"foreignKey:IDTrx" json:"detail_trx,omitempty"`
}
