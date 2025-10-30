package main

import (
	"go-crud/config"
	"go-crud/models"
	"log"
)

func main() {
	config.ConnectDatabase()

	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
	)

	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	println("All tables migrated successfully!")

}