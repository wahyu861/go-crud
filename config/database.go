package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectDatabase() {
	dbUser := "root"
	dbPass := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "crud_go"

	// Format DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Gagal koneksi ke database: %v", err))
	}

	fmt.Println("Berhasil konek ke database")
}