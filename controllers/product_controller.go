package controllers

import (
	"fmt"
	"go-crud/config"
	"go-crud/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ✅ List Produk
func ProductList(c echo.Context) error {
	var products []models.Product
	result := config.DB.Find(&products)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "Gagal mengambil data produk")
	}

	tmpl, err := template.ParseFiles("views/layout.html", "views/product_list.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Template error: %v", err))
	}

	data := map[string]interface{}{
		"Title":    "Daftar Produk",
		"Products": products,
	}
	return tmpl.ExecuteTemplate(c.Response(), "layout.html", data)
}

// ✅ Form Tambah Produk
func ProductForm(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/layout.html", "views/product_form.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Template error: %v", err))
	}

	data := map[string]interface{}{
		"Title":  "Tambah Produk",
		"Action": "/products/add",
	}
	return tmpl.ExecuteTemplate(c.Response(), "layout.html", data)
}

// ✅ Simpan Produk Baru
func ProductAdd(c echo.Context) error {
	name := c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))

	product := models.Product{Name: name, Price: price, Stock: stock}
	if err := config.DB.Create(&product).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Gagal menambahkan produk")
	}

	return c.Redirect(http.StatusSeeOther, "/products")
}

// ✅ Form Edit Produk
func ProductEditForm(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Produk tidak ditemukan")
	}

	tmpl, err := template.ParseFiles("views/layout.html", "views/product_form.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Template error: %v", err))
	}

	data := map[string]interface{}{
		"Title":   "Edit Produk",
		"Action":  "/products/edit/" + id,
		"Product": product,
	}
	return tmpl.ExecuteTemplate(c.Response(), "layout.html", data)
}

// ✅ Update Produk
func ProductUpdate(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Produk tidak ditemukan")
	}

	product.Name = c.FormValue("name")
	product.Price, _ = strconv.ParseFloat(c.FormValue("price"), 64)
	product.Stock, _ = strconv.Atoi(c.FormValue("stock")) // ❗ huruf kecil 'stock' (bukan 'Stock')

	if err := config.DB.Save(&product).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Gagal update produk")
	}

	return c.Redirect(http.StatusSeeOther, "/products")
}

// ✅ Hapus Produk
func ProductDelete(c echo.Context) error {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Gagal menghapus produk")
	}

	return c.Redirect(http.StatusSeeOther, "/products")
}
