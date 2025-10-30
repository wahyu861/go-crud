package models

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	Phone     string    `gorm:"size:20;unique"`
	Password  string    `gorm:"size:255;not null"`
	IsAdmin   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relasi
	Store        *Store         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Addresses    []Address     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Transactions []Transaction `gorm:"foreignKey:BuyerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}