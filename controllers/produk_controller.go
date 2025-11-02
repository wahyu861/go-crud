package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// GET /api/products
func GetAllProducts(c echo.Context) error {
	var products []models.Produk
	if err := config.DB.Preload("Toko").Preload("Category").Preload("Fotos").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", products))
}

// GET /api/products/:id
func GetProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid product ID", []string{err.Error()}))
	}

	var product models.Produk
	if err := config.DB.Preload("Toko").Preload("Category").Preload("Fotos").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Product not found", []string{"Produk tidak ditemukan"}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", product))
}

// POST /api/products (pemilik toko)
func CreateProduct(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	// Ambil toko milik user login
	var store models.Toko
	if err := config.DB.Where("id_user = ?", authUser.ID).First(&store).Error; err != nil {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("You don't have a store", []string{"User belum memiliki toko"}))
	}

	var req models.Produk
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{err.Error()}))
	}

	req.IDToko = store.ID
	req.Slug = strings.ToLower(strings.ReplaceAll(req.NamaProduk, " ", "-"))

	if err := config.DB.Create(&req).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create product", []string{err.Error()}))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Product created successfully", req))
}

// PUT /api/products/:id (pemilik toko)
func UpdateProduct(c echo.Context) error {
	authUser, _ := getAuthUser(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var product models.Produk
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Product not found", []string{"Produk tidak ditemukan"}))
	}

	// Cek apakah produk milik toko user
	var store models.Toko
	if err := config.DB.First(&store, product.IDToko).Error; err != nil {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Unauthorized", []string{"Toko tidak ditemukan"}))
	}
	if store.IDUser != authUser.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Tidak dapat mengubah produk milik toko lain"}))
	}

	// Bind input sementara
	var req struct {
		NamaProduk    string  `json:"nama_produk" form:"nama_produk"`
		HargaKonsumen float64 `json:"harga_konsumen" form:"harga_konsumen"`
		HargaReseller float64 `json:"harga_reseller" form:"harga_reseller"`
		Stok          int     `json:"stok" form:"stok"`
		Deskripsi     string  `json:"deskripsi" form:"deskripsi"`
		IDCategory    uint64  `json:"id_category" form:"id_category"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{err.Error()}))
	}

	// Hanya update field yang dikirim
	updates := map[string]interface{}{}
	if req.NamaProduk != "" {
		updates["nama_produk"] = req.NamaProduk
		updates["slug"] = strings.ToLower(strings.ReplaceAll(req.NamaProduk, " ", "-"))
	}
	if req.HargaKonsumen != 0 {
		updates["harga_konsumen"] = req.HargaKonsumen
	}
	if req.HargaReseller != 0 {
		updates["harga_reseller"] = req.HargaReseller
	}
	if req.Stok != 0 {
		updates["stok"] = req.Stok
	}
	if req.Deskripsi != "" {
		updates["deskripsi"] = req.Deskripsi
	}
	if req.IDCategory != 0 {
		updates["id_category"] = req.IDCategory
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("No data to update", []string{"Tidak ada data yang diubah"}))
	}

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update product", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to UPDATE data", "Update produk succeed"))
}


// DELETE /api/products/:id (pemilik toko)
func DeleteProduct(c echo.Context) error {
	authUser, _ := getAuthUser(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var product models.Produk
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Product not found", []string{"Produk tidak ditemukan"}))
	}

	var store models.Toko
	if err := config.DB.First(&store, product.IDToko).Error; err != nil {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Unauthorized", []string{"Toko tidak ditemukan"}))
	}
	if store.IDUser != authUser.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Tidak dapat menghapus produk milik toko lain"}))
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete product", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Product deleted successfully", nil))
}
