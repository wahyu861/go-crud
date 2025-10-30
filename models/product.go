package models

import "time"

type Product struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	StoreID     uint64    `gorm:"not null;index"`
	CategoryID  *uint64   `gorm:"index"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(15,2);not null"`
	Stock       int       `gorm:"not null"`
	Image       string    `gorm:"size:255"` // opsional, bisa kosong
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relasi
	Store    Store    `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionItems []TransactionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}