package models

import "time"

type Address struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	UserID     uint64    `gorm:"not null;index"`
	Province   string    `gorm:"size:100"`
	City       string    `gorm:"size:100"`
	District   string    `gorm:"size:100"`
	PostalCode string    `gorm:"size:10"`
	Detail     string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Relasi
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}