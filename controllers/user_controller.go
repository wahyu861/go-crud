package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// ===================================================
// 🔹 GET /users (Hanya Admin)
// ===================================================
func GetAllUsers(c echo.Context) error {
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User tidak ditemukan dalam context", []string{"unauthorized"}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Akses ditolak. Hanya admin yang dapat melihat semua user.", []string{"forbidden"}))
	}

	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal mengambil data user", []string{err.Error()}))
	}

	var enrichedUsers []map[string]interface{}
	for _, user := range users {
		user.KataSandi = ""

		var provinsi, kota interface{}
		if user.IDProvinsi != nil {
			provinsi, _ = utils.GetProvinceByID(*user.IDProvinsi)
		}
		if user.IDKota != nil {
			kota, _ = utils.GetCityByID(*user.IDKota)
		}

		enrichedUsers = append(enrichedUsers, map[string]interface{}{
			"user":        user,
			"id_provinsi": provinsi,
			"id_kota":     kota,
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET all users", enrichedUsers))
}

// ===================================================
// 🔹 GET /users/:id (User Lihat Profil Sendiri)
// ===================================================
func GetUserByID(c echo.Context) error {
	idParam := c.Param("id")

	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User tidak ditemukan dalam context", []string{"unauthorized"}))
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", idParam).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User tidak ditemukan", []string{"user_not_found"}))
	}

	if !authUser.IsAdmin && authUser.ID != user.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Tidak boleh melihat data user lain", []string{"forbidden"}))
	}

	user.KataSandi = ""

	var provinsi, kota interface{}
	if user.IDProvinsi != nil {
		provinsi, _ = utils.GetProvinceByID(*user.IDProvinsi)
	}
	if user.IDKota != nil {
		kota, _ = utils.GetCityByID(*user.IDKota)
	}

	data := map[string]interface{}{
		"user":        user,
		"id_provinsi": provinsi,
		"id_kota":     kota,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET user", data))
}

// ===================================================
// 🔹 PUT /users/:id (Update Data User)
// ===================================================
func UpdateUser(c echo.Context) error {
	id := c.Param("id")

	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User tidak ditemukan dalam context", []string{"unauthorized"}))
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User tidak ditemukan", []string{"user_not_found"}))
	}

	if !authUser.IsAdmin && authUser.ID != user.ID {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Tidak bisa mengedit user lain", []string{"forbidden"}))
	}

	type UpdateUserRequest struct {
		Nama          *string    `json:"nama"`
		KataSandi     *string    `json:"kata_sandi"`
		NoTelp        *string    `json:"no_telp"`
		TanggalLahir  *time.Time `json:"tanggal_lahir"`
		JenisKelamin  *string    `json:"jenis_kelamin"`
		Tentang       *string    `json:"tentang"`
		Pekerjaan     *string    `json:"pekerjaan"`
		Email         *string    `json:"email"`
		IDProvinsi    *string    `json:"id_provinsi"`
		IDKota        *string    `json:"id_kota"`
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Input tidak valid", []string{err.Error()}))
	}

	if req.Nama != nil {
		user.Nama = *req.Nama
	}
	if req.KataSandi != nil && *req.KataSandi != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.KataSandi), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal hash password", []string{err.Error()}))
		}
		user.KataSandi = string(hashed)
	}
	if req.NoTelp != nil {
		user.NoTelp = req.NoTelp
	}
	if req.TanggalLahir != nil {
		user.TanggalLahir = req.TanggalLahir
	}
	if req.JenisKelamin != nil {
		user.JenisKelamin = req.JenisKelamin
	}
	if req.Tentang != nil {
		user.Tentang = req.Tentang
	}
	if req.Pekerjaan != nil {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.IDProvinsi != nil {
		user.IDProvinsi = req.IDProvinsi
	}
	if req.IDKota != nil {
		user.IDKota = req.IDKota
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan perubahan", []string{err.Error()}))
	}

	user.KataSandi = ""

	var provinsi, kota interface{}
	if user.IDProvinsi != nil {
		provinsi, _ = utils.GetProvinceByID(*user.IDProvinsi)
	}
	if user.IDKota != nil {
		kota, _ = utils.GetCityByID(*user.IDKota)
	}

	data := map[string]interface{}{
		"user":        user,
		"id_provinsi": provinsi,
		"id_kota":     kota,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("User berhasil diperbarui", data))
}

// ===================================================
// 🔹 DELETE /users/:id (Hanya Admin)
// ===================================================
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User tidak ditemukan dalam context", []string{"unauthorized"}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Hanya admin yang bisa menghapus user", []string{"forbidden"}))
	}

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menghapus user", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("User berhasil dihapus", nil))
}
