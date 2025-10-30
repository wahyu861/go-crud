package models

import "time"

type TransactionStatus string

const (
	StatusPending		TransactionStatus = "pending"
	StatusPaid			TransactionStatus = "paid"
	StatusShipped		TransactionStatus = "shipped"
	StatusCompleted		TransactionStatus = "completed"
	StatusCancelled		TransactionStatus = "cancelled"
)

type Transaction struct {
	ID         uint64             `gorm:"primaryKey;autoIncrement"`
	BuyerID    uint64             `gorm:"not null;index"`
	StoreID    uint64             `gorm:"not null;index"`
	TotalPrice float64            `gorm:"type:decimal(15,2);not null"`
	Status     TransactionStatus  `gorm:"type:enum('pending','paid','shipped','completed','cancelled');default:'pending'"`
	CreatedAt  time.Time          `gorm:"autoCreateTime"`
	UpdatedAt  time.Time          `gorm:"autoUpdateTime"`

	// Relasi
	Buyer 			User `gorm:"foreignKey:BuyerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Store 			Store `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionItems []TransactionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

