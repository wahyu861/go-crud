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
		&models.Toko{},
		&models.Alamat{},
		&models.Category{},
		&models.Produk{},
		&models.FotoProduk{},
		&models.LogProduk{},
		&models.Trx{},
		&models.DetailTrx{},
	)

	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	println("All tables migrated successfully!")

}