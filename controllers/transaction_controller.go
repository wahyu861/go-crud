package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ===================== HANDLERS ======================

// POST /api/transactions
func CreateTransaction(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	var req struct {
		MethodBayar      string `json:"method_bayar"`
		AlamatPengiriman uint64 `json:"alamat_pengiriman"`
		DetailTrx        []struct {
			IDProduk  uint64 `json:"id_produk"`
			Kuantitas int    `json:"kuantitas"`
		} `json:"detail_trx"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{"Gagal membaca input transaksi"}))
	}

	if len(req.DetailTrx) == 0 {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{"Tidak ada produk yang dibeli"}))
	}

	// Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	trx := models.Trx{
		IDUser:           authUser.ID,
		AlamatPengiriman: &req.AlamatPengiriman,
		KodeInvoice:      "INV-" + strconv.FormatInt(time.Now().Unix(), 10),
		MethodBayar:      &req.MethodBayar,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal membuat transaksi", []string{err.Error()}))
	}

	totalHarga := 0

	for _, item := range req.DetailTrx {
		if item.IDProduk == 0 || item.Kuantitas <= 0 {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{"Produk tidak valid"}))
		}

		var product models.Produk
		if err := tx.First(&product, item.IDProduk).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusNotFound, utils.ErrorResponse("Produk tidak ditemukan", []string{err.Error()}))
		}

		// Cek stok
		if product.Stok < item.Kuantitas {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Stok tidak mencukupi", []string{product.NamaProduk}))
		}

		// Kurangi stok
		product.Stok -= item.Kuantitas
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal mengurangi stok", []string{err.Error()}))
		}

		// Hitung subtotal
		subtotal := item.Kuantitas * product.HargaKonsumen
		totalHarga += subtotal

		// Simpan log produk
		log := models.LogProduk{
			IDProduk:      product.ID,
			NamaProduk:    product.NamaProduk,
			Slug:          strings.ToLower(strings.ReplaceAll(product.NamaProduk, " ", "-")),
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Deskripsi:     product.Deskripsi,
			IDToko:        product.IDToko,
			IDCategory:    product.IDCategory,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan log produk", []string{err.Error()}))
		}

		// Simpan detail transaksi
		detail := models.DetailTrx{
			IDTrx:       trx.ID,
			IDLogProduk: log.ID,
			IDToko:      product.IDToko,
			Kuantitas:   item.Kuantitas,
			HargaTotal:  subtotal,
		}
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan detail transaksi", []string{err.Error()}))
		}
	}

	// Simpan total harga ke transaksi
	trx.HargaTotal = totalHarga
	if err := tx.Save(&trx).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan total harga", []string{err.Error()}))
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Commit gagal", []string{err.Error()}))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Transaksi berhasil dibuat", map[string]interface{}{
		"id":           trx.ID,
		"kode_invoice": trx.KodeInvoice,
		"harga_total":  trx.HargaTotal,
		"method_bayar": req.MethodBayar,
	}))
}

// GET /api/transactions (Admin only)
func GetAllTransactions(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Hanya admin yang dapat melihat semua transaksi"}))
	}

	var trans []models.Trx
	if err := config.DB.Preload("DetailTrx.LogProduk").
		Order("created_at desc").
		Find(&trans).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal mengambil data transaksi", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", trans))
}

// GET /api/transactions/:id
func GetTransactionByID(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", []string{"ID transaksi tidak valid"}))
	}

	var trx models.Trx
	if err := config.DB.
		Preload("Alamat").
		Preload("DetailTrx.LogProduk.Toko").
		Preload("DetailTrx.LogProduk.Category").
		Preload("DetailTrx.LogProduk.Photos").
		First(&trx, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Transaksi tidak ditemukan", []string{err.Error()}))
	}

	// hanya pemilik transaksi atau admin yang boleh lihat
	if trx.IDUser != authUser.ID && !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Anda tidak memiliki akses ke transaksi ini"}))
	}

	// ===============================
	// ✅ bentuk ulang sesuai contoh
	// ===============================

	response := map[string]interface{}{
		"id":           trx.ID,
		"harga_total":  trx.HargaTotal,
		"kode_invoice": trx.KodeInvoice,
		"method_bayar": trx.MethodBayar,
		"alamat_kirim": map[string]interface{}{
			"id":             trx.Alamat.ID,
			"judul_alamat":   trx.Alamat.JudulAlamat,
			"nama_penerima":  trx.Alamat.NamaPenerima,
			"no_telp":        trx.Alamat.NoTelp,
			"detail_alamat":  trx.Alamat.DetailAlamat,
		},
	}

	var details []map[string]interface{}
	for _, d := range trx.DetailTrx {
		p := d.LogProduk

		// ambil semua foto produk
		var photos []map[string]interface{}
		for _, f := range p.Photos {
			photos = append(photos, map[string]interface{}{
				"id":         f.ID,
				"product_id": f.IDProduk,
				"url":        f.URL,
			})
		}

		details = append(details, map[string]interface{}{
			"product": map[string]interface{}{
				"id":              p.IDProduk,
				"nama_produk":     p.NamaProduk,
				"slug":            p.Slug,
				"harga_reseler":   p.HargaReseller,
				"harga_konsumen":  p.HargaKonsumen,
				"deskripsi":       p.Deskripsi,
				"toko": map[string]interface{}{
					"nama_toko": p.Toko.NamaToko,
					"url_foto":  p.Toko.UrlFoto,
				},
				"category": map[string]interface{}{
					"id":            p.Category.ID,
					"nama_category": p.Category.NamaCategory,
				},
				"photos": photos,
			},
			"toko": map[string]interface{}{
				"id":         p.Toko.ID,
				"nama_toko":  p.Toko.NamaToko,
				"url_foto":   p.Toko.UrlFoto,
			},
			"kuantitas":   d.Kuantitas,
			"harga_total": d.HargaTotal,
		})
	}

	response["detail_trx"] = details

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", response))
}

