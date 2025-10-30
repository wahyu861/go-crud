package models

import "time"

type Store struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	UserID      uint64    `gorm:"not null;index"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	UrlFoto     string    `json:"url_foto" gorm:"type:varchar(255)"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relasi
	User         User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Products     []Product      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Transactions []Transaction  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}