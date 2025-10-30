package models

import "time"

type Category struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relasi
	Products []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}