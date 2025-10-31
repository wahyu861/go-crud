package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ========================== RESPONSE STRUCT ==========================

type SafeUserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokoResponse struct {
	ID          uint64            `json:"id"`
	NamaToko    string            `json:"nama_toko"`
	UrlFoto     string            `json:"url_foto"`
	IDUser      uint64            `json:"user_id"`
	User        *SafeUserResponse `json:"user,omitempty"`
}

// Ambil user dari context JWT
func getAuthUser(c echo.Context) (*models.User, error) {
	authUserRaw := c.Get("authUser")
	if authUserRaw == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User tidak ditemukan dalam context (token tidak valid?)")
	}

	authUser, ok := authUserRaw.(models.User)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Gagal membaca data user dari context")
	}
	return &authUser, nil
}

// ========================== HANDLER ===============================

// GET /api/toko (Admin only)
func GetAllToko(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to GET data", []string{"Hanya admin yang dapat melihat semua toko"}))
	}

	var toko []models.Toko
	if err := config.DB.Preload("User").Find(&toko).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}

	var result []TokoResponse
	for _, s := range toko {
		resp := TokoResponse{
			ID:       s.ID,
			NamaToko: s.NamaToko,
			UrlFoto:  *s.UrlFoto,
			IDUser:   s.IDUser,
		}
		if s.User.ID != 0 {
			resp.User = &SafeUserResponse{
				ID:    s.User.ID,
				Name:  s.User.Nama,
				Email: s.User.Email,
			}
		}
		result = append(result, resp)
	}

	data := map[string]interface{}{
		"page":  1,
		"limit": 10,
		"data":  result,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", data))
}

// GET /api/toko/my
func GetMyToko(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}

	var toko models.Toko
	if err := config.DB.Preload("User").Where("id_user = ?", authUser.ID).First(&toko).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to GET data", []string{"Toko tidak ditemukan"}))
	}

	resp := TokoResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  *toko.UrlFoto,
		IDUser:   toko.IDUser,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", resp))
}

// GET /api/toko/:id
func GetTokoByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to GET data", []string{"ID toko tidak valid"}))
	}

	var toko models.Toko
	if err := config.DB.Preload("User").First(&toko, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to GET data", []string{"Toko tidak ditemukan"}))
	}

	resp := TokoResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  *toko.UrlFoto,
		IDUser:   toko.IDUser,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", resp))
}

// PUT /api/toko/:id
func UpdateToko(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"ID toko tidak valid"}))
	}

	var toko models.Toko
	if err := config.DB.First(&toko, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to UPDATE data", []string{"Toko tidak ditemukan"}))
	}

	if toko.IDUser != authUser.ID && !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to UPDATE data", []string{"Tidak bisa mengubah toko milik orang lain"}))
	}

	var input struct {
		NamaToko    string `json:"nama_toko" form:"nama_toko"`
		UrlFoto     string `json:"url_foto" form:"url_foto"`
		Description string `json:"description" form:"description"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"Input tidak valid"}))
	}

	updates := map[string]interface{}{}
	if input.NamaToko != "" {
		updates["name"] = input.NamaToko
	}
	if input.UrlFoto != "" {
		updates["url_foto"] = input.UrlFoto
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to UPDATE data", []string{"Tidak ada data yang diubah"}))
	}

	if err := config.DB.Model(&toko).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to UPDATE data", "Update toko succeed"))
}

// DELETE /api/toko/:id (admin only - nonaktifkan toko)
func DeleteToko(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Hanya admin yang dapat menonaktifkan toko"}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to DELETE data", []string{"ID toko tidak valid"}))
	}

	var toko models.Toko
	if err := config.DB.First(&toko, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to DELETE data", []string{"Toko tidak ditemukan"}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to DELETE data", "Toko berhasil dinonaktifkan"))
}
