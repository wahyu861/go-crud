package models

import "time"

type User struct {
	ID           uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama         string      `gorm:"type:varchar(100);not null" json:"nama"`
	KataSandi    string      `gorm:"type:varchar(255);not null" json:"kata_sandi"`
	NoTelp       *string     `gorm:"type:varchar(20);unique" json:"no_telp"`     
	TanggalLahir *time.Time  `json:"tanggal_lahir"`                             
	JenisKelamin *string     `gorm:"type:enum('Laki-laki','Perempuan')" json:"jenis_kelamin"` 
	Tentang      *string     `gorm:"type:text" json:"tentang"`                  
	Pekerjaan    *string     `gorm:"type:varchar(100)" json:"pekerjaan"`        
	Email        string      `gorm:"type:varchar(100);unique;not null" json:"email"`
	IDProvinsi   *string     `gorm:"type:varchar(10)" json:"id_provinsi"`       
	IDKota       *string     `gorm:"type:varchar(10)" json:"id_kota"`           
	IsAdmin      bool        `gorm:"not null;default:false" json:"is_admin"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Toko     *Toko     `gorm:"foreignKey:IDUser" json:"toko,omitempty"`
	Alamat   []Alamat  `gorm:"foreignKey:IDUser" json:"alamat,omitempty"`
	Trx      []Trx     `gorm:"foreignKey:IDUser" json:"trx,omitempty"`
}
