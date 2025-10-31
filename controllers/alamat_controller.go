package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ===================================================
// GET /api/alamat/my
// ===================================================
func GetMyAlamat(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	// ambil semua kolom dulu
	var alamat []models.Alamat
	if err := config.DB.Where("id_user = ?", authUser.ID).Find(&alamat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}

	// mapping ke struct untuk response clean
	type AlamatResponse struct {
		ID            uint64   `json:"id"`
		JudulAlamat   string `json:"judul_alamat"`
		NamaPenerima  string `json:"nama_penerima"`
		NoTelp        string `json:"no_telp"`
		DetailAlamat  string `json:"detail_alamat"`
	}

	var result []AlamatResponse
	for _, a := range alamat {
		result = append(result, AlamatResponse{
			ID:           a.ID,
			JudulAlamat:  a.JudulAlamat,
			NamaPenerima: a.NamaPenerima,
			NoTelp:       a.NoTelp,
			DetailAlamat: a.DetailAlamat,
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", result))
}



// ===================================================
// GET /api/alamat/:id (Ambil satu alamat milik user login)
// ===================================================
func GetAlamatByID(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to GET data", []string{"ID tidak valid"}))
	}

	var alamat models.Alamat
	if err := config.DB.First(&alamat, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to GET data", []string{"Alamat tidak ditemukan"}))
	}

	// Pastikan user hanya bisa lihat alamat miliknya
	if alamat.IDUser != authUser.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to GET data", []string{"Tidak dapat melihat alamat orang lain"}))
	}

	// Ambil data provinsi & kota dari API EMSIFA
	var provinsi, kota interface{}
	if alamat.User.IDProvinsi != nil && *alamat.User.IDProvinsi != "" {
		provinsi, _ = utils.GetProvinceByID(*alamat.User.IDProvinsi)
	}
	if alamat.User.IDKota != nil && *alamat.User.IDKota != "" {
		kota, _ = utils.GetCityByID(*alamat.User.IDKota)
	}

	data := map[string]interface{}{
		"alamat":   alamat,
		"provinsi": provinsi,
		"kota":     kota,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", data))
}

// ===================================================
// POST /api/alamat
// ===================================================
func CreateAlamat(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	var input models.Alamat
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to CREATE data", []string{"Input tidak valid"}))
	}

	input.IDUser = authUser.ID

	if err := config.DB.Create(&input).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to CREATE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to CREATE data", input))
}

// ===================================================
// PUT /api/alamat/:id
// ===================================================
func UpdateAlamat(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"ID tidak valid"}))
	}

	var alamat models.Alamat
	if err := config.DB.First(&alamat, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to UPDATE data", []string{"Alamat tidak ditemukan"}))
	}

	if alamat.IDUser != authUser.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to UPDATE data", []string{"Tidak dapat mengubah alamat orang lain"}))
	}

	var input models.Alamat
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"Input tidak valid"}))
	}

	// 🔒 Siapkan field yang akan diupdate
	updates := map[string]interface{}{}

	// Cek apakah field User tidak nil sebelum akses
	if input.User != nil {
		if input.User.IDProvinsi != nil && *input.User.IDProvinsi != "" {
			updates["id_provinsi"] = *input.User.IDProvinsi
		}
		if input.User.IDKota != nil && *input.User.IDKota != "" {
			updates["id_kota"] = *input.User.IDKota
		}
	}

	// Field lain
	if input.DetailAlamat != "" {
		updates["detail_alamat"] = input.DetailAlamat
	}
	if input.JudulAlamat != "" {
		updates["judul_alamat"] = input.JudulAlamat
	}
	if input.NamaPenerima != "" {
		updates["nama_penerima"] = input.NamaPenerima
	}
	if input.NoTelp != "" {
		updates["no_telp"] = input.NoTelp
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"Tidak ada data yang diubah"}))
	}

	if err := config.DB.Model(&alamat).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to UPDATE data", ""))
}

// ===================================================
//  DELETE /api/alamat/:id
// ===================================================
func DeleteAlamat(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to DELETE data", []string{"ID tidak valid"}))
	}

	var alamat models.Alamat
	if err := config.DB.First(&alamat, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to DELETE data", []string{"Alamat tidak ditemukan"}))
	}

	if alamat.IDUser != authUser.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to DELETE data", []string{"Tidak dapat menghapus alamat orang lain"}))
	}

	if err := config.DB.Delete(&alamat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to DELETE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to DELETE data", "Alamat berhasil dihapus"))
}
