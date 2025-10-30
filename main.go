package main

import (
	"fmt"
	"go-crud/config"
	"go-crud/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// 🔹 1. Koneksi ke database
	config.ConnectDatabase()
	fmt.Println("✅ Berhasil konek ke database")

	// 🔹 2. Buat instance Echo
	e := echo.New()

	// 🔹 3. Load semua route dari folder routes
	routes.InitRoutes(e)

	// 🔹 4. Jalankan server
	e.Logger.Fatal(e.Start(":8080"))
}

