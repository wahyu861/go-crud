package models

type TransactionItem struct {
	ID            uint64  `gorm:"primaryKey;autoIncrement"`
	TransactionID uint64  `gorm:"not null;index"`
	ProductID     uint64  `gorm:"not null;index"`
	Quantity      int     `gorm:"not null"`
	Price         float64 `gorm:"type:decimal(15,2);not null"`
	Subtotal      float64 `gorm:"type:decimal(15,2);not null"`

	//Relasi
	Transaction Transaction `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product     Product     `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}